package tenant

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	tenantmodel "github.com/akshit_tyagi/postgresql_project/internal/models/tenant"
	tenantservice "github.com/akshit_tyagi/postgresql_project/internal/services/tenant"
	"github.com/akshit_tyagi/postgresql_project/internal/validations"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	var req tenantmodel.TenantListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := tenantservice.GetAll(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong + " " + err.Error()})
		return
	}
	response := tenantmodel.TenantListAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleFetchedSuccess,
		Data:       role.Data,
		Pagination: tenantmodel.Pagination{
			CurrentPage: role.CurrentPage,
			PerPage:     role.PerPage,
			Total:       role.Total,
			LastPage:    role.LastPage,
		},
	}
	c.JSON(http.StatusOK, response)
}

func Create(c *gin.Context) {
	var req tenantmodel.TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.ValidateTenant(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := tenantservice.Create(req)
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
	response := tenantmodel.TenantAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleCreatedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}

func GetByID(c *gin.Context) {
	idStr := c.Param("id")
	role, err := tenantservice.GetByID(idStr)
	if err != nil {
		if err == constants.RoleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	response := tenantmodel.TenantAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleRetrievedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}

func Update(c *gin.Context) {
	var req tenantmodel.TenantRequest
	id := c.Param("id")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.ValidateTenant(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := tenantservice.Update(id, req)
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
	response := tenantmodel.TenantAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleUpdatedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	role, err := tenantservice.Delete(id)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	response := tenantmodel.TenantAPIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    constants.RoleDeletedSuccess,
		Data:       *role,
	}
	c.JSON(http.StatusOK, response)
}
