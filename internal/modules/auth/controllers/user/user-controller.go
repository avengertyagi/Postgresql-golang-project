package user

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	authconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	userservice "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/validations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	service userservice.TempUserService
	otp     userservice.OTPService
}

func NewUserController(s userservice.TempUserService, otp userservice.OTPService) *UserController {
	return &UserController{service: s, otp: otp}
}

func (ctl *UserController) Register(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error() + ": " + err.Error(), "data": gin.H{}})
		return
	}
	if err := validations.UserRegisterValidation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	user, err := ctl.service.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, authconstants.EmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, authconstants.MobileAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("user register error", "error", err)
		return
	}
	otpResponse, err := ctl.otp.Create(c.Request.Context(), &models.TempUser{Model: gorm.Model{ID: user.ID}, Email: user.Email, Mobile: user.Mobile})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		return
	}
	user.EmailOTP = otpResponse.EmailOTP
	user.MobileOTP = otpResponse.MobileOTP
	c.JSON(http.StatusOK, gin.H{"status": true, "statusCode": http.StatusOK, "message": authconstants.SignUpSuccess.Error(), "data": *user})
}

func (ctl *UserController) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	user, err := ctl.otp.Verify(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "statusCode": http.StatusOK, "message": authconstants.SignUpSuccess.Error(), "data": *user})
}
