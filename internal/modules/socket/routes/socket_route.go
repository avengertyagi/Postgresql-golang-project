package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/modules/socket/controllers"
	"github.com/gin-gonic/gin"
)

func SocketRoutes(router *gin.RouterGroup, controller *controllers.SocketController) {
	router.GET("/ws", controller.HandleSocket)
	router.GET("/ws/stats", controller.GetStats)
}
