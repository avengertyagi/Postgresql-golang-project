package controllers

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/common"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/helpers"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/role/dto"
	roleservice "github.com/akshit_tyagi/postgresql_project/internal/modules/role/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/role/validations"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	service roleservice.RoleService
}

func NewRoleController(s roleservice.RoleService) *RoleController {
	return &RoleController{service: s}
}

func (ctl *RoleController) List(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	roleList, total, err := ctl.service.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": err.Error(), "data": gin.H{}})
		return
	}
	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleFetchedSuccess,
		"data":       roleList,
		"pagination": &common.Pagination{
			CurrentPage: query.Page,
			PerPage:     query.Limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

func (ctl *RoleController) Create(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := ctl.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, constants.RoleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleCreatedSuccess,
		"data":       *role,
	})
}

func (ctl *RoleController) GetByID(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == constants.RoleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleRetrievedSuccess,
		"data":       *role,
	})
}

func (ctl *RoleController) Update(c *gin.Context) {
	var req dto.RoleRequest
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := ctl.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound.Error()})
			return
		}
		if errors.Is(err, constants.RoleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": constants.RoleAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleUpdatedSuccess,
		"data":       *role,
	})
}

func (ctl *RoleController) Delete(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := ctl.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleDeletedSuccess,
		"data":       *role,
	})
}
