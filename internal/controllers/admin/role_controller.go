package admin

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
	"github.com/akshit_tyagi/postgresql_project/internal/services"
	"github.com/akshit_tyagi/postgresql_project/internal/validations"
	"github.com/gin-gonic/gin"
)

func CreateRole(c *gin.Context) {
	var req models.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.ValidateRole(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := services.CreateRole(req)
	if err != nil {
		if errors.Is(err, constants.RoleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"status":     false,
				"statusCode": http.StatusConflict,
				"message":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleCreatedSuccess,
		"data":       role,
	})
}
