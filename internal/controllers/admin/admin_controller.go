package admin

import (
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
	"github.com/akshit_tyagi/postgresql_project/internal/services"
	"github.com/akshit_tyagi/postgresql_project/internal/validations"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid request. Please provide email and password."})
		return
	}
	if err := validations.AdminLoginValidation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	admin, err := services.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    constants.LoginSuccess,
		"data":       gin.H{"admin": admin},
	})
}

func Logout(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid request. Please provide refresh token."})
		return
	}
	if err := services.Logout(req.RefreshToken); err != nil {
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
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid request. Please provide refresh token."})
		return
	}
	resp, err := services.RefreshToken(req.RefreshToken)
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
	profile, err := services.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": false, "statusCode": http.StatusOK, "message": constants.ProfileFetchSuccess})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.ProfileFetchSuccess,
		"data":       gin.H{"profile": profile},
	})
}
