package health

import (
	"net/http"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"uptime":  time.Since(startedAt).String(),
		"version": "1.0.0",
	})
}

func Readyz(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": false, "message": "database not initialized"})
		return
	}
	sqlDB, err := config.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": false, "message": "database handle unavailable"})
		return
	}
	if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": false, "message": "database ping failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "ready"})
}

var startedAt = time.Now()
