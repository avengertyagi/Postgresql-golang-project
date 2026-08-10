package admin

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/common"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
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
		response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusBadRequest, Message: err.Error(), Data: gin.H{}}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := validations.AdminLoginValidation(req); err != nil {
		response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusBadRequest, Message: err.Error(), Data: gin.H{}}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	admin, err := ctl.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, constants.InvalidCredentials) {
			response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: err.Error(), Data: gin.H{}}
			c.JSON(http.StatusUnauthorized, response)
			return
		}
		if errors.Is(err, constants.InactiveAccount) {
			response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: err.Error(), Data: gin.H{}}
			c.JSON(http.StatusUnauthorized, response)
			return
		}
		if errors.Is(err, constants.InactiveAccount) {
			response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: err.Error(), Data: gin.H{}}
			c.JSON(http.StatusUnauthorized, response)
			return
		}
		response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusInternalServerError, Message: constants.SomethingWentWrong + err.Error(), Data: gin.H{}}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := common.ApiSuccessResponse{Status: true, StatusCode: http.StatusOK, Message: constants.LoginSuccess, Data: *admin}
	c.JSON(http.StatusOK, response)
}

func (ctl *AdminController) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusBadRequest, Message: err.Error(), Data: gin.H{}}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := ctl.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if errors.Is(err, constants.SessionAlreadyRevoked) {
			response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusUnauthorized, Message: err.Error(), Data: gin.H{}}
			c.JSON(http.StatusUnauthorized, response)
			return
		}
		response := common.ApiErrorResponse{Status: false, StatusCode: http.StatusInternalServerError, Message: constants.SomethingWentWrong + err.Error(), Data: gin.H{}}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := common.ApiSuccessResponse{Status: true, StatusCode: http.StatusOK, Message: constants.LogoutSuccess, Data: gin.H{}}
	c.JSON(http.StatusOK, response)
}

func (ctl *AdminController) RefreshToken(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	resp, err := ctl.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": err.Error(), "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.RefreshSuccess,
		"data":       resp,
	})
}

func (ctl *AdminController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated, "data": gin.H{}})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated, "data": gin.H{}})
		return
	}
	profile, err := ctl.service.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": false, "statusCode": http.StatusOK, "message": constants.ProfileFetchSuccess, "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.ProfileFetchSuccess,
		"data":       profile,
	})
}
