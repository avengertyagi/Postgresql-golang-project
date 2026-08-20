package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/modules/websocket/controllers"
	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(router *gin.RouterGroup, controller *controllers.WebSocketController) {
	router.GET("/ws", controller.HandleWebSocket)
	router.GET("/ws/stats", controller.GetStats)
}
