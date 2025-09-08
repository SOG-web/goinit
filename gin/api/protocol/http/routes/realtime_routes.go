package routes

import (
	sseHandler "github.com/SOG-web/gin/api/protocol/sse"
	wsHandler "github.com/SOG-web/gin/api/protocol/ws"
	"github.com/gin-gonic/gin"
)

// SetupSSERoutes sets up Server-Sent Events routes
func SetupSSERoutes(router *gin.Engine) {
	sse := sseHandler.NewSSEHandler()

	sseGroup := router.Group("/api/sse")
	{
		sseGroup.GET("/events", sse.StreamEvents)
		sseGroup.GET("/notifications", sse.StreamNotifications)
	}
}

// SetupWSRoutes sets up WebSocket routes
func SetupWSRoutes(router *gin.Engine) {
	ws := wsHandler.NewWebSocketHandler()

	wsGroup := router.Group("/api/ws")
	{
		wsGroup.GET("/connect", ws.HandleConnection)
	}
}