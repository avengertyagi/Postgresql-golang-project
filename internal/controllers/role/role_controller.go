package role

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/models/role"
	roleservice "github.com/akshit_tyagi/postgresql_project/internal/services/role"
	"github.com/akshit_tyagi/postgresql_project/internal/validations"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	var req rolemodel.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := roleservice.GetAll(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong + " " + err.Error()})
		return
	}
	response := rolemodel.RoleListAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleFetchedSuccess,
		Data:       role.Data,
		Pagination: rolemodel.Pagination{
			CurrentPage: role.CurrentPage,
			PerPage:     role.PerPage,
			Total:       role.Total,
			LastPage:    role.LastPage,
		},
	}
	c.JSON(http.StatusOK, response)
}

func Create(c *gin.Context) {
	var req rolemodel.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.ValidateRole(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := roleservice.Create(req)
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
	response := rolemodel.RoleAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleCreatedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}

func GetByID(c *gin.Context) {
	idStr := c.Param("id")
	role, err := roleservice.GetByID(idStr)
	if err != nil {
		if err == constants.RoleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	response := rolemodel.RoleAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleRetrievedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}

func Update(c *gin.Context) {
	var req rolemodel.RoleRequest
	id := c.Param("id")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.ValidateRole(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := roleservice.Update(id, req)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		if errors.Is(err, constants.RoleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": constants.RoleAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	response := rolemodel.RoleAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleUpdatedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	role, err := roleservice.Delete(id)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	response := rolemodel.RoleAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleDeletedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}
