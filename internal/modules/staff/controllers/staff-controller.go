package controllers

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/common"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/helpers"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/dto"
	staffservice "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/staff/validations"
	"github.com/gin-gonic/gin"
)

type StaffController struct {
	service staffservice.StaffService
}

func NewStaffController(s staffservice.StaffService) *StaffController {
	return &StaffController{service: s}
}

func (ctl *StaffController) List(c *gin.Context) {
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
		"message":    constants.StaffFetchedSuccess,
		"data":       roleList,
		"pagination": &common.Pagination{
			CurrentPage: query.Page,
			PerPage:     query.Limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

func (ctl *StaffController) Create(c *gin.Context) {
	var req dto.StaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := validations.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	staff, err := ctl.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, constants.StaffAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.StaffCreatedSuccess,
		"data":       *staff,
	})
}

func (ctl *StaffController) GetByID(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == constants.StaffNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.StaffNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.StaffRetrievedSuccess,
		"data":       *role,
	})
}

func (ctl *StaffController) Update(c *gin.Context) {
	var req dto.StaffRequest
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
		if errors.Is(err, constants.StaffNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.StaffNotFound.Error()})
			return
		}
		if errors.Is(err, constants.StaffAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": constants.StaffAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.StaffUpdatedSuccess,
		"data":       *role,
	})
}

func (ctl *StaffController) Delete(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}
	role, err := ctl.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, constants.StaffNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.StaffNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    constants.StaffDeletedSuccess,
		"data":       *role,
	})
}
