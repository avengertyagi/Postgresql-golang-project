package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
	"gorm.io/gorm"
)

type OTPService interface {
	Create(ctx context.Context, user *models.TempUser) (*dto.UserRegisterResponse, error)
	Verify(ctx context.Context, req dto.VerifyOTPRequest) (*dto.UserRegisterResponse, error)
}

type otpService struct {
	db       *gorm.DB
	tempRepo repositories.TempUserRepository
}

func NewOTPService(db *gorm.DB, tempRepo repositories.TempUserRepository) OTPService {
	return &otpService{db: db, tempRepo: tempRepo}
}

func (s *otpService) Create(ctx context.Context, user *models.TempUser) (*dto.UserRegisterResponse, error) {
	emailCode, err := generateOTP()
	if err != nil {
		return nil, err
	}
	mobileCode, err := generateOTP()
	if err != nil {
		return nil, err
	}
	otp := &models.OTP{TempUserID: user.ID, EmailCode: emailCode, MobileCode: mobileCode, ExpiresAt: time.Now().Add(10 * time.Minute)}
	if err := s.db.WithContext(ctx).Where("temp_user_id = ?", user.ID).Assign(ototpValues(otp)).FirstOrCreate(otp).Error; err != nil {
		return nil, err
	}
	slog.Info("registration OTP generated", "email", user.Email, "mobile", user.Mobile, "email_otp", emailCode, "mobile_otp", mobileCode)
	return &dto.UserRegisterResponse{ID: user.ID, Email: user.Email, Mobile: user.Mobile, EmailOTP: emailCode, MobileOTP: mobileCode}, nil
}

func ototpValues(otp *models.OTP) map[string]interface{} {
	return map[string]interface{}{"email_code": otp.EmailCode, "mobile_code": otp.MobileCode, "expires_at": otp.ExpiresAt}
}

func (s *otpService) Verify(ctx context.Context, req dto.VerifyOTPRequest) (*dto.UserRegisterResponse, error) {
	var tempUser models.TempUser
	var otp models.OTP
	if err := s.db.WithContext(ctx).First(&tempUser, req.TempUserID).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Where("temp_user_id = ?", req.TempUserID).First(&otp).Error; err != nil {
		return nil, err
	}
	if time.Now().After(otp.ExpiresAt) || otp.EmailCode != req.EmailOTP || otp.MobileCode != req.MobileOTP {
		return nil, fmt.Errorf("invalid or expired OTP")
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := &models.User{Name: tempUser.Name, Email: tempUser.Email, Mobile: tempUser.Mobile, Password: tempUser.Password, CountryCode: tempUser.CountryCode, Dob: tempUser.Dob, DeviceToken: tempUser.DeviceToken, DeviceType: tempUser.DeviceType}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		tempUser.ID = user.ID
		if err := tx.Delete(&tempUser).Error; err != nil {
			return err
		}
		return tx.Delete(&otp).Error
	})
	if err != nil {
		return nil, err
	}
	return &dto.UserRegisterResponse{ID: tempUser.ID, Name: tempUser.Name, Email: tempUser.Email, Mobile: tempUser.Mobile, CountryCode: tempUser.CountryCode, Dob: tempUser.Dob}, nil
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
