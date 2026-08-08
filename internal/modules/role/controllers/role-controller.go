package role

import (
	"errors"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/common"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	svc "github.com/akshit_tyagi/postgresql_project/internal/modules/role/services"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/role/validations"
	"github.com/gin-gonic/gin"
)

func GetAllWithPagination(c *gin.Context) {
	var req rolemodel.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "statusCode": http.StatusBadRequest, "message": err.Error(), "data": gin.H{}})
		return
	}
	roleList, total, lastPage, err := svc.GetAllWithPagination(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "statusCode": http.StatusInternalServerError, "message": err.Error(), "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleFetchedSuccess,
		"data":       roleList,
		"pagination": &common.Pagination{
			CurrentPage: req.Page,
			PerPage:     req.Limit,
			Total:       total,
			LastPage:    lastPage,
		},
	})
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
	role, err := svc.Create(req)
	if err != nil {
		if errors.Is(err, constants.RoleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": false, "statusCode": http.StatusConflict, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleCreatedSuccess,
		"data":       *role,
	})
}

func GetByID(c *gin.Context) {
	idStr := c.Param("id")
	role, err := svc.GetByID(idStr)
	if err != nil {
		if err == constants.RoleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleRetrievedSuccess,
		"data":       *role,
	})
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
	role, err := svc.Update(id, req)
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
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleUpdatedSuccess,
		"data":       *role,
	})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	role, err := svc.Delete(id)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.RoleNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "statusCode": http.StatusInternalServerError, "message": constants.SomethingWentWrong})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    constants.RoleDeletedSuccess,
		"data":       *role,
	})
}
