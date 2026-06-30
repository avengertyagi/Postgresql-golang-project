package admin

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/models/user"
	adminservice "github.com/akshit_tyagi/postgresql_project/internal/services/admin"
	"github.com/akshit_tyagi/postgresql_project/internal/validations"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req usermodel.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.AdminLoginValidation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	admin, err := adminservice.Login(req)
	if err != nil {
		if errors.Is(err, constants.InvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     false,
				"statusCode": http.StatusUnauthorized,
				"message":    err.Error(),
			})
			return
		}
		if errors.Is(err, constants.InactiveAccount) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     false,
				"statusCode": http.StatusUnauthorized,
				"message":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	response := usermodel.AdminAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.LoginSuccess,
		Data:       *admin,
	}
	c.JSON(http.StatusOK, response)
}

func Logout(c *gin.Context) {
	var req usermodel.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := adminservice.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    constants.LogoutSuccess,
	})
}

func RefreshToken(c *gin.Context) {
	var req usermodel.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	resp, err := adminservice.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Token refreshed successfully",
		"data":       resp,
	})
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated})
		return
	}
	profile, err := adminservice.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": false, "statusCode": http.StatusOK, "message": constants.ProfileFetchSuccess})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.ProfileFetchSuccess,
		"data":       profile,
	})
}
