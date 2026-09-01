package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	authconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	adminservice "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/validations"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	service adminservice.AdminService
}

func NewAuthController(s adminservice.AdminService) *AdminController {
	return &AdminController{service: s}
}
func (ctl *AdminController) Login(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	if err := validations.AdminLoginValidation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	admin, err := ctl.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, authconstants.InvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, authconstants.InactiveAccount) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("admin login error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "statusCode": http.StatusOK, "message": authconstants.LoginSuccess.Error(), "data": *admin})
}

func (ctl *AdminController) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	if err := ctl.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if errors.Is(err, authconstants.SessionNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, authconstants.SessionAlreadyRevoked) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, authconstants.SessionExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("admin logout error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "statusCode": http.StatusOK, "message": authconstants.LogoutSuccess.Error(), "data": gin.H{}})
}

func (ctl *AdminController) RefreshToken(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	resp, err := ctl.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, authconstants.SessionNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, authconstants.SessionAlreadyRevoked) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, authconstants.SessionExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("admin refresh token error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    authconstants.RefreshSuccess.Error(),
		"data":       resp,
	})
}

func (ctl *AdminController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated.Error(), "data": gin.H{}})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated.Error(), "data": gin.H{}})
		return
	}
	profile, err := ctl.service.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.NotFound.Error(), "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    authconstants.ProfileFetchSuccess.Error(),
		"data":       profile,
	})
}

func (ctl *AdminController) UpdateProfile(c *gin.Context) {
	var req dto.AdminUpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated.Error(), "data": gin.H{}})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated.Error(), "data": gin.H{}})
		return
	}
	var profilePictureFile interface{}
	file, err := c.FormFile("profile_picture")
	if err == nil {
		profilePictureFile = file
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid file upload", "data": gin.H{}})
		return
	}

	profile, err := ctl.service.UpdateProfile(c.Request.Context(), id, req, profilePictureFile)
	if err != nil {
		if errors.Is(err, constants.NotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.NotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": err.Error(), "data": gin.H{}})
		slog.Error("admin update profile error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    authconstants.ProfileUpdateSuccess.Error(),
		"data":       profile,
	})
}
