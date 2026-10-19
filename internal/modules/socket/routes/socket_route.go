package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/modules/socket/controllers"
	"github.com/gin-gonic/gin"
)

func SocketRoutes(router *gin.RouterGroup) {
	router.GET("/ws", controllers.HandleSocket)
	router.GET("/ws/stats", controllers.GetStats)
}
