package sse

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

// SSEHandler handles Server-Sent Events
type SSEHandler struct{}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

// StreamEvents streams server-sent events to the client
// @Summary Stream server events
// @Description Establishes an SSE connection for real-time updates
// @Tags SSE
// @Accept json
// @Produce text/event-stream
// @Success 200 {string} string "SSE stream established"
// @Router /api/sse/events [get]
func (h *SSEHandler) StreamEvents(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create a channel to send events
	eventChan := make(chan sse.Event)

	// Start a goroutine to send periodic events
	go func() {
		defer close(eventChan)
		counter := 0

		for {
			select {
			case <-c.Request.Context().Done():
				// Client disconnected
				return
			default:
				// Send an event
				event := sse.Event{
					Event: "message",
					Data:  fmt.Sprintf(`{"timestamp": "%s", "counter": %d, "message": "Hello from server!"}`, time.Now().Format(time.RFC3339), counter),
					Id:    fmt.Sprintf("%d", counter),
				}
				eventChan <- event
				counter++

				// Wait 2 seconds before sending next event
				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Stream the events to the client
	c.Stream(func(w io.Writer) bool {
		if event, ok := <-eventChan; ok {
			c.Render(-1, event)
			return true
		}
		return false
	})
}

// StreamNotifications streams notification events
// @Summary Stream notifications
// @Description Establishes an SSE connection for notifications
// @Tags SSE
// @Accept json
// @Produce text/event-stream
// @Success 200 {string} string "Notification SSE stream established"
// @Router /api/sse/notifications [get]
func (h *SSEHandler) StreamNotifications(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	eventChan := make(chan sse.Event)

	go func() {
		defer close(eventChan)
		notificationID := 0

		for {
			select {
			case <-c.Request.Context().Done():
				return
			default:
				notification := sse.Event{
					Event: "notification",
					Data:  fmt.Sprintf(`{"id": %d, "type": "info", "message": "System notification %d", "timestamp": "%s"}`, notificationID, notificationID, time.Now().Format(time.RFC3339)),
					Id:    fmt.Sprintf("notif-%d", notificationID),
				}
				eventChan <- notification
				notificationID++

				// Send notifications every 5 seconds
				time.Sleep(5 * time.Second)
			}
		}
	}()

	c.Stream(func(w io.Writer) bool {
		if event, ok := <-eventChan; ok {
			c.Render(-1, event)
			return true
		}
		return false
	})
}