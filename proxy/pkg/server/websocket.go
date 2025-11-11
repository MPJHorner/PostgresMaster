package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MPJHorner/PostgresMaster/proxy/pkg/auth"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/postgres"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/protocol"
	"github.com/gorilla/websocket"
)

// PostgresClient defines the interface for Postgres operations
type PostgresClient interface {
	ExecuteQuery(ctx context.Context, sql string, params []interface{}) (*postgres.QueryResult, error)
	IntrospectSchema(ctx context.Context) (*protocol.SchemaPayload, error)
}

// Server represents a WebSocket server
type Server struct {
	secret   string
	upgrader websocket.Upgrader
	pgClient PostgresClient
}

// NewServer creates a new WebSocket server
func NewServer(secret string, pgClient PostgresClient) *Server {
	return &Server{
		secret:   secret,
		pgClient: pgClient,
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
		return s.handleQuery(msg)
	case protocol.TypeIntrospect:
		return s.handleIntrospect(msg)
	default:
		return protocol.NewError(msg.ID, "INVALID_MESSAGE_TYPE", fmt.Sprintf("Unknown message type: %s", msg.Type), "")
	}
}

// handleQuery processes query execution requests
func (s *Server) handleQuery(msg protocol.ClientMessage) protocol.ServerMessage {
	// Parse the payload
	payloadBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return protocol.NewError(msg.ID, "INVALID_PAYLOAD", "Failed to parse payload", err.Error())
	}

	var payload protocol.QueryPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return protocol.NewError(msg.ID, "INVALID_PAYLOAD", "Failed to unmarshal query payload", err.Error())
	}

	// Validate SQL is not empty
	if payload.SQL == "" {
		return protocol.NewError(msg.ID, "EMPTY_QUERY", "SQL query cannot be empty", "")
	}

	// Create context with timeout if specified
	ctx := context.Background()
	if payload.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(payload.Timeout)*time.Millisecond)
		defer cancel()
	}

	// Execute the query
	result, err := s.pgClient.ExecuteQuery(ctx, payload.SQL, payload.Params)
	if err != nil {
		return protocol.NewError(msg.ID, "QUERY_ERROR", err.Error(), "")
	}

	// Return the result
	return protocol.NewQueryResult(msg.ID, result.Rows, result.Columns, result.ExecutionTime)
}

// handleIntrospect processes schema introspection requests
func (s *Server) handleIntrospect(msg protocol.ClientMessage) protocol.ServerMessage {
	// Create context with reasonable timeout for introspection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Introspect the schema
	schema, err := s.pgClient.IntrospectSchema(ctx)
	if err != nil {
		return protocol.NewError(msg.ID, "INTROSPECTION_ERROR", err.Error(), "")
	}

	// Return the schema
	return protocol.NewSchemaResult(msg.ID, schema.Tables, schema.Functions)
}

// SendMessage sends a message to the client
func SendMessage(conn *websocket.Conn, msg protocol.ServerMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
