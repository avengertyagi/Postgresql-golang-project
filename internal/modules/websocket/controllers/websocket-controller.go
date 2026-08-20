package controllers

import (
	"log/slog"
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/helpers"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/websocket/services"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	hub *services.Hub
}

func NewWebSocketController(hub *services.Hub) *WebSocketController {
	return &WebSocketController{
		hub: hub,
	}
}

func (wsc *WebSocketController) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		slog.Warn("WebSocket: missing token",
			"request_id", requestid.Get(c),
			"path", c.Request.URL.Path,
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Unauthorized: missing token",
		})
		return
	}
	claims, err := helpers.ParseAccessToken(token)
	if err != nil {
		slog.Warn("WebSocket: invalid token",
			"request_id", requestid.Get(c),
			"error", err.Error(),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Unauthorized: invalid token",
		})
		return
	}
	c.Set("user_id", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("role", claims.Role)
	c.Set("guard", claims.Guard)
	c.Set("permissions", claims.Permissions)

	slog.Info("WebSocket: token verified",
		"request_id", requestid.Get(c),
		"user_id", claims.UserID,
		"email", claims.Email,
	)
	conn, err := services.GetUpgrader().Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("WebSocket upgrade error", "error", err)
		return
	}

	userID := c.GetString("user_id")
	userName := claims.Email

	client := wsc.hub.NewClient(conn, userID, userName)
	wsc.hub.Register <- client

	slog.Info("New WebSocket connection",
		"user_id", userID,
		"email", userName,
		"request_id", requestid.Get(c),
	)
	go client.ReadPump()
	go client.WritePump()
}

func (wsc *WebSocketController) GetStats(c *gin.Context) {
	connectedClients := wsc.hub.GetConnectedClients()
	c.JSON(http.StatusOK, gin.H{
		"status":            true,
		"connected_clients": connectedClients,
		"message":           "WebSocket stats",
	})
}
