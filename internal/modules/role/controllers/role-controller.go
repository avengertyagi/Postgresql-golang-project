package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/common"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/helpers"
	roleconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/role/constants"
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
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("role list error", "error", err)
		return
	}
	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    roleconstants.RoleFetchedSuccess.Error(),
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
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	if err := validations.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := ctl.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, roleconstants.RoleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("role create error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    roleconstants.RoleCreatedSuccess.Error(),
		"data":       *role,
	})
}

func (ctl *RoleController) GetByID(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == roleconstants.RoleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": roleconstants.RoleNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("role get by id error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    roleconstants.RoleRetrievedSuccess.Error(),
		"data":       *role,
	})
}

func (ctl *RoleController) Update(c *gin.Context) {
	var req dto.RoleRequest
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	if err := validations.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := ctl.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, roleconstants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": roleconstants.RoleNotFound.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, roleconstants.RoleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": roleconstants.RoleAlreadyExists.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("role update error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    roleconstants.RoleUpdatedSuccess.Error(),
		"data":       *role,
	})
}

func (ctl *RoleController) UpdateStatus(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := ctl.service.UpdateStatus(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, roleconstants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": roleconstants.RoleNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("role update status error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    roleconstants.RoleStatusUpdatedSuccess.Error(),
		"data":       *role,
	})
}

func (ctl *RoleController) Delete(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := ctl.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, roleconstants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": roleconstants.RoleNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("role delete error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    roleconstants.RoleDeletedSuccess.Error(),
		"data":       *role,
	})
}
