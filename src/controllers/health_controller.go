package controller

import (
	"net/http"
	"time"

	"github.com/akshit_tyagi/postgresql_project/src/config"
	"github.com/gin-gonic/gin"
)

// Healthz is a liveness probe. It always returns 200 if the process is
// running; load balancers and orchestrators use it to decide when to
// restart a stuck instance.
func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"uptime":  time.Since(startedAt).String(),
		"version": "1.0.0",
	})
}

// Readyz is a readiness probe. It pings the database and returns 503
// if the connection is unavailable, so orchestrators stop routing
// traffic until the dependency is back.
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
