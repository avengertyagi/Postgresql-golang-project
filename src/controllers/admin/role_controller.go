package controller

import (
	"net/http"

	rolemodel "github.com/akshit_tyagi/postgresql_project/src/models"
	roleservice "github.com/akshit_tyagi/postgresql_project/src/services"
	rolevalidation "github.com/akshit_tyagi/postgresql_project/src/validations"
	"github.com/akshit_tyagi/postgresql_project/src/constants"
	"github.com/gin-gonic/gin"
)

func CreateRole(c *gin.Context) {
	var req rolemodel.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": "Invalid request. Please provide role name."})
		return
	}
	if err := rolevalidation.ValidateRole(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := roleservice.CreateRole(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.ROLE_CREATED_SUCCESS,
		"data":       gin.H{"role": role},
	})
}
