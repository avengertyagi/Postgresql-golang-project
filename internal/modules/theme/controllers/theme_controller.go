package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/common"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/helpers"
	themeconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/dto"
	themeservice "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/theme/validations"
	"github.com/gin-gonic/gin"
)

type ThemeController struct {
	service themeservice.ThemeService
}

func NewThemeController(s themeservice.ThemeService) *ThemeController {
	return &ThemeController{service: s}
}

func (ctl *ThemeController) List(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	themeList, total, err := ctl.service.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("theme list error", "error", err)
		return
	}
	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    themeconstants.ThemeFetchedSuccess.Error(),
		"data":       themeList,
		"pagination": &common.Pagination{
			CurrentPage: query.Page,
			PerPage:     query.Limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

func (ctl *ThemeController) Create(c *gin.Context) {
	var req dto.ThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": constants.InvalidRequestBody.Error(), "data": gin.H{}})
		return
	}
	if err := validations.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	theme, err := ctl.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, themeconstants.ThemeAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("theme create error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    themeconstants.ThemeCreatedSuccess.Error(),
		"data":       *theme,
	})
}

func (ctl *ThemeController) GetByID(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	theme, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == themeconstants.ThemeNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": themeconstants.ThemeNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("theme get by id error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    themeconstants.ThemeRetrievedSuccess.Error(),
		"data":       *theme,
	})
}

func (ctl *ThemeController) Update(c *gin.Context) {
	var req dto.ThemeRequest
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
	theme, err := ctl.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, themeconstants.ThemeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": themeconstants.ThemeNotFound.Error(), "data": gin.H{}})
			return
		}
		if errors.Is(err, themeconstants.ThemeAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": themeconstants.ThemeAlreadyExists.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("theme update error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    themeconstants.ThemeUpdatedSuccess.Error(),
		"data":       *theme,
	})
}

func (ctl *ThemeController) UpdateStatus(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	theme, err := ctl.service.UpdateStatus(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, themeconstants.ThemeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": themeconstants.ThemeNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("theme update status error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    themeconstants.ThemeStatusUpdatedSuccess.Error(),
		"data":       *theme,
	})
}

func (ctl *ThemeController) Delete(c *gin.Context) {
	id, err := helpers.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	theme, err := ctl.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, themeconstants.ThemeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": themeconstants.ThemeNotFound.Error(), "data": gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong.Error(), "data": gin.H{}})
		slog.Error("theme delete error", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": http.StatusOK,
		"message":    themeconstants.ThemeDeletedSuccess.Error(),
		"data":       *theme,
	})
}
