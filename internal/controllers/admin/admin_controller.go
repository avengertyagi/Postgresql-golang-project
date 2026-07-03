package admin

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/models/user"
	adminservice "github.com/akshit_tyagi/postgresql_project/internal/services/admin"
	"github.com/akshit_tyagi/postgresql_project/internal/validations"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req usermodel.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ApiResponse{Status: false, StatusCode: http.StatusBadRequest, Message: err.Error(), Data: gin.H{}})
		return
	}
	if err := validations.AdminLoginValidation(req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ApiResponse{Status: false, StatusCode: http.StatusBadRequest, Message: err.Error(), Data: gin.H{}})
		return
	}
	admin, err := adminservice.Login(req)
	if err != nil {
		if errors.Is(err, constants.InvalidCredentials) {
			c.JSON(http.StatusUnauthorized, helpers.ApiResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: err.Error(), Data: gin.H{}})
			return
		}
		if errors.Is(err, constants.InactiveAccount) {
			c.JSON(http.StatusUnauthorized, helpers.ApiResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: err.Error(), Data: gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, helpers.ApiResponse{Status: false, StatusCode: http.StatusInternalServerError, Message: constants.SomethingWentWrong, Data: gin.H{}})
		return
	}
	c.JSON(http.StatusOK, helpers.ApiResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.LoginSuccess,
		Data:       *admin,
	})
}

func Logout(c *gin.Context) {
	var req usermodel.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ApiResponse{Status: false, StatusCode: http.StatusBadRequest, Message: err.Error(), Data: gin.H{}})
		return
	}
	if err := adminservice.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, helpers.ApiResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: err.Error(), Data: gin.H{}})
		return
	}
	c.JSON(http.StatusOK, helpers.ApiResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.LogoutSuccess,
		Data:       nil,
	})
}

func RefreshToken(c *gin.Context) {
	var req usermodel.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ApiResponse{Status: false, StatusCode: http.StatusBadRequest, Message: err.Error(), Data: gin.H{}})
		return
	}
	resp, err := adminservice.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusForbidden, helpers.ApiResponse{Status: false, StatusCode: http.StatusForbidden, Message: err.Error(), Data: gin.H{}})
		return
	}
	c.JSON(http.StatusOK, helpers.ApiResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Token refreshed successfully",
		Data:       resp,
	})
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.ApiResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: constants.Unauthenticated, Data: gin.H{}})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, helpers.ApiResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: constants.Unauthenticated, Data: gin.H{}})
		return
	}
	profile, err := adminservice.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusOK, helpers.ApiResponse{Status: false, StatusCode: http.StatusOK, Message: constants.ProfileFetchSuccess, Data: gin.H{}})
		return
	}
	c.JSON(http.StatusOK, helpers.ApiResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.ProfileFetchSuccess,
		Data:       profile,
	})
}
