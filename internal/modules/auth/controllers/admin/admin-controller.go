package admin

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	request "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/request/admin"
	adminservice "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/validations"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	if err := validations.AdminLoginValidation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	admin, err := adminservice.Login(req)
	if err != nil {
		if errors.Is(err, constants.InvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, constants.InactiveAccount) {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "statusCode": http.StatusUnauthorized, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong + err.Error(), "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.LoginSuccess,
		"data":       *admin,
	})
}

func Logout(c *gin.Context) {
	var req usermodel.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	if err := adminservice.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong + err.Error(), "data": gin.H{}})
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.LogoutSuccess,
		"data":       gin.H{},
	})
}

func RefreshToken(c *gin.Context) {
	var req usermodel.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	resp, err := adminservice.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "statusCode": http.StatusForbidden, "message": err.Error(), "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.RefreshSuccess,
		"data":       resp,
	})
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated, "data": gin.H{}})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated, "data": gin.H{}})
		return
	}
	profile, err := adminservice.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "statusCode": http.StatusOK, "message": constants.ProfileFetchSuccess, "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.ProfileFetchSuccess,
		"data":       profile,
	})
}
