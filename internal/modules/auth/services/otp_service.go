package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	"gorm.io/gorm"
)

func CreateOTP(ctx context.Context, user *models.TempUser) (*dto.UserRegisterResponse, error) {
	emailCode, err := generateOTP()
	if err != nil {
		return nil, err
	}
	mobileCode := ""
	if user.CountryCode == "+91" || user.CountryCode == "91" {
		generatedMobileCode, err := generateOTP()
		if err != nil {
			return nil, err
		}
		mobileCode = generatedMobileCode
	}
	otp := &models.OTP{TempUserID: user.ID, EmailCode: emailCode, MobileCode: mobileCode, ExpiresAt: time.Now().Add(10 * time.Minute)}
	if err := config.DB.WithContext(ctx).Where("temp_user_id = ?", user.ID).Assign(ototpValues(otp)).FirstOrCreate(otp).Error; err != nil {
		return nil, err
	}
	slog.Info("registration OTP generated", "email", user.Email, "mobile", user.Mobile, "email_otp", emailCode, "mobile_otp", mobileCode)
	return &dto.UserRegisterResponse{ID: user.ID, Email: user.Email, Mobile: user.Mobile, EmailOTP: emailCode, MobileOTP: mobileCode}, nil
}

func ototpValues(otp *models.OTP) map[string]interface{} {
	return map[string]interface{}{"email_code": otp.EmailCode, "mobile_code": otp.MobileCode, "expires_at": otp.ExpiresAt}
}

func VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.UserRegisterResponse, error) {
	var tempUser models.TempUser
	var otp models.OTP
	if err := config.DB.WithContext(ctx).First(&tempUser, req.TempUserID).Error; err != nil {
		return nil, err
	}
	if err := config.DB.WithContext(ctx).Where("temp_user_id = ?", req.TempUserID).First(&otp).Error; err != nil {
		return nil, err
	}
	if time.Now().After(otp.ExpiresAt) || otp.EmailCode != req.EmailOTP || otp.MobileCode != req.MobileOTP {
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	var userRole rolemodel.Role
	if err := config.DB.WithContext(ctx).Where("name = ?", "User").FirstOrCreate(&userRole).Error; err != nil {
		return nil, fmt.Errorf("failed to ensure default user role: %w", err)
	}

	var createdUserID uint
	err := config.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := &models.User{
			RoleID:      userRole.ID,
			UserType:    constants.UserRole,
			Name:        tempUser.Name,
			Email:       tempUser.Email,
			Mobile:      tempUser.Mobile,
			Password:    tempUser.Password,
			CountryCode: tempUser.CountryCode,
			Dob:         tempUser.Dob,
			DeviceToken: tempUser.DeviceToken,
			DeviceType:  tempUser.DeviceType,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		createdUserID = user.ID
		if err := tx.Delete(&models.TempUser{}, tempUser.ID).Error; err != nil {
			return err
		}
		return tx.Delete(&otp).Error
	})
	if err != nil {
		return nil, err
	}
	return &dto.UserRegisterResponse{ID: createdUserID, Name: tempUser.Name, Email: tempUser.Email, Mobile: tempUser.Mobile, CountryCode: tempUser.CountryCode, Dob: tempUser.Dob}, nil
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
