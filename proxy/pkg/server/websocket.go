package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MPJHorner/PostgresMaster/proxy/pkg/auth"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/protocol"
	"github.com/gorilla/websocket"
)

// Server represents a WebSocket server
type Server struct {
	secret   string
	upgrader websocket.Upgrader
}

// NewServer creates a new WebSocket server
func NewServer(secret string) *Server {
	return &Server{
		secret: secret,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from localhost only
				origin := r.Header.Get("Origin")
				return origin == "http://localhost:5173" ||
					origin == "http://localhost:3000" ||
					origin == "http://127.0.0.1:5173" ||
					origin == "http://127.0.0.1:3000" ||
					origin == "" // Allow non-browser clients
			},
		},
	}
}

// HandleConnection upgrades HTTP connection to WebSocket and handles messages
func (s *Server) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Extract secret from query parameter
	clientSecret := r.URL.Query().Get("secret")
	if !auth.ValidateSecret(clientSecret) || clientSecret != s.secret {
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	// Message handling loop
	for {
		var msg protocol.ClientMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle message based on type
		response := s.handleMessage(msg)

		// Send response
		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Failed to send response: %v", err)
			break
		}
	}

	log.Println("Client disconnected")
}

// handleMessage routes messages to appropriate handlers
func (s *Server) handleMessage(msg protocol.ClientMessage) protocol.ServerMessage {
	switch msg.Type {
	case protocol.TypePing:
		return protocol.NewPong(msg.ID)
	case protocol.TypeQuery:
		// TODO: Implement query execution
		return protocol.NewError(msg.ID, "NOT_IMPLEMENTED", "Query execution not yet implemented", "")
	case protocol.TypeIntrospect:
		// TODO: Implement schema introspection
		return protocol.NewError(msg.ID, "NOT_IMPLEMENTED", "Schema introspection not yet implemented", "")
	default:
		return protocol.NewError(msg.ID, "INVALID_MESSAGE_TYPE", fmt.Sprintf("Unknown message type: %s", msg.Type), "")
	}
}

// SendMessage sends a message to the client
func SendMessage(conn *websocket.Conn, msg protocol.ServerMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
