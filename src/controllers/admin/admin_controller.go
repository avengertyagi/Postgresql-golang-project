package controller

import (
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/src/constants"
	authmodel "github.com/akshit_tyagi/postgresql_project/src/models"
	adminservice "github.com/akshit_tyagi/postgresql_project/src/services"
	adminvalidation "github.com/akshit_tyagi/postgresql_project/src/validations"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req authmodel.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid request. Please provide email and password."})
		return
	}
	if err := adminvalidation.AdminLoginValidation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	admin, err := adminservice.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.INVALID_CREDENTIALS})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    constants.LOGIN_SUCCESS,
		"data":       gin.H{"admin": admin},
	})
}

func Logout(c *gin.Context) {
	var req authmodel.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid request. Please provide refresh token."})
		return
	}
	if err := adminservice.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    constants.LOGOUT_SUCCESS,
	})
}

func RefreshToken(c *gin.Context) {
	var req authmodel.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid request. Please provide refresh token."})
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
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.UNAUTHENTICATED})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.UNAUTHENTICATED})
		return
	}
	profile, err := adminservice.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": false, "statusCode": http.StatusOK, "message": constants.PROFILE_FETCH_SUCCESS})
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.PROFILE_FETCH_SUCCESS,
		"data":       gin.H{"profile": profile},
	})
}
