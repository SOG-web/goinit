package routes

import (
	"github.com/gin-gonic/gin"
	sseHandler "sog.com/goinit/gin/api/protocol/sse"
	wsHandler "sog.com/goinit/gin/api/protocol/ws"
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