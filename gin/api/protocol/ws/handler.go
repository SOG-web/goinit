package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mutex    sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from any origin in development
				// In production, you should check the origin
				return true
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// HandleConnection upgrades HTTP connection to WebSocket
// @Summary Establish WebSocket connection
// @Description Upgrades HTTP connection to WebSocket for real-time communication
// @Tags WebSocket
// @Accept json
// @Produce json
// @Success 101 {string} string "WebSocket connection established"
// @Router /api/ws/connect [get]
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Add client to the list
	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	log.Printf("New WebSocket connection established. Total clients: %d", len(h.clients))

	// Send welcome message
	welcomeMsg := Message{
		Type:      "welcome",
		Data:      "Connected to WebSocket server",
		Timestamp: time.Now(),
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		h.removeClient(conn)
		return
	}

	// Start a goroutine to send periodic pings
	go h.sendPeriodicPings(conn)

	// Handle incoming messages
	h.handleMessages(conn)
}

// handleMessages processes incoming WebSocket messages
func (h *WebSocketHandler) handleMessages(conn *websocket.Conn) {
	defer h.removeClient(conn)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		msg.Timestamp = time.Now()
		log.Printf("Received message: %+v", msg)

		// Echo the message back with a response
		response := Message{
			Type:      "echo",
			Data:      fmt.Sprintf("Echo: %v", msg.Data),
			Timestamp: time.Now(),
		}

		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Error sending response: %v", err)
			break
		}

		// Broadcast to all clients if it's a broadcast message
		if msg.Type == "broadcast" {
			h.broadcastMessage(msg)
		}
	}
}

// sendPeriodicPings sends ping messages to keep connection alive
func (h *WebSocketHandler) sendPeriodicPings(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteJSON(Message{
				Type:      "ping",
				Data:      "keepalive",
				Timestamp: time.Now(),
			}); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
		}
	}
}

// broadcastMessage sends a message to all connected clients
func (h *WebSocketHandler) broadcastMessage(msg Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	broadcastMsg := Message{
		Type:      "broadcast",
		Data:      msg.Data,
		Timestamp: time.Now(),
	}

	for client := range h.clients {
		if err := client.WriteJSON(broadcastMsg); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			h.removeClient(client)
		}
	}
}

// removeClient removes a client from the clients map
func (h *WebSocketHandler) removeClient(conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.clients, conn)
	conn.Close()
	log.Printf("WebSocket connection closed. Total clients: %d", len(h.clients))
}

// GetClientCount returns the number of connected clients
func (h *WebSocketHandler) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}