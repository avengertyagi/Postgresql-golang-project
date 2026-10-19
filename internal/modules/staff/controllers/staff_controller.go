package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/common"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/helpers"
	staffconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/dto"
	staffservice "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/staff/validations"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	staffList, total, err := staffservice.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("staff list error", "error", err)
		return
	}
	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    staffconstants.StaffFetchedSuccess.Error(),
		"data":       staffList,
		"pagination": &common.Pagination{
			CurrentPage: query.Page,
			PerPage:     query.Limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

func Create(c *gin.Context) {
	var req dto.StaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	if err := validations.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	staff, err := staffservice.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, staffconstants.StaffAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, staffconstants.RoleIdNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error() + err.Error(), "data": gin.H{}})
		slog.Error("staff create error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    staffconstants.StaffCreatedSuccess.Error(),
		"data":       *staff,
	})
}

func GetByID(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := staffservice.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == staffconstants.StaffNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": staffconstants.StaffNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("staff get by id error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    staffconstants.StaffRetrievedSuccess.Error(),
		"data":       *role,
	})
}

func Update(c *gin.Context) {
	var req dto.StaffRequest
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
	role, err := staffservice.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, staffconstants.StaffNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": staffconstants.StaffNotFound.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, staffconstants.StaffAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": staffconstants.StaffAlreadyExists.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, staffconstants.RoleIdNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("staff update error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    staffconstants.StaffUpdatedSuccess.Error(),
		"data":       *role,
	})
}

func UpdateStatus(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := staffservice.UpdateStatus(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, staffconstants.StaffNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": staffconstants.StaffNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("staff update status error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    staffconstants.StaffStatusUpdatedSuccess.Error(),
		"data":       *role,
	})
}

func Delete(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	role, err := staffservice.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, staffconstants.StaffNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": staffconstants.StaffNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("staff delete error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    staffconstants.StaffDeletedSuccess.Error(),
		"data":       *role,
	})
}
