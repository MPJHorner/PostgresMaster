package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MPJHorner/PostgresMaster/proxy/pkg/auth"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/postgres"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/protocol"
	"github.com/gorilla/websocket"
)

// MockPostgresClient implements the PostgresClient interface for testing
type MockPostgresClient struct {
	ExecuteQueryFunc      func(ctx context.Context, sql string, params []interface{}) (*postgres.QueryResult, error)
	IntrospectSchemaFunc  func(ctx context.Context) (*protocol.SchemaPayload, error)
}

func (m *MockPostgresClient) ExecuteQuery(ctx context.Context, sql string, params []interface{}) (*postgres.QueryResult, error) {
	if m.ExecuteQueryFunc != nil {
		return m.ExecuteQueryFunc(ctx, sql, params)
	}
	return &postgres.QueryResult{
		Rows:          []map[string]interface{}{},
		Columns:       []protocol.ColumnInfo{},
		RowCount:      0,
		ExecutionTime: 0,
	}, nil
}

func (m *MockPostgresClient) IntrospectSchema(ctx context.Context) (*protocol.SchemaPayload, error) {
	if m.IntrospectSchemaFunc != nil {
		return m.IntrospectSchemaFunc(ctx)
	}
	return &protocol.SchemaPayload{
		Tables:    []protocol.TableInfo{},
		Functions: []protocol.FunctionInfo{},
	}, nil
}

func TestNewServer(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{}
	server := NewServer(secret, mockClient)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}

	if server.secret != secret {
		t.Errorf("Expected secret %s, got %s", secret, server.secret)
	}

	if server.pgClient == nil {
		t.Error("Expected pgClient to be set")
	}
}

func TestHandleConnection_InvalidSecret(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{}
	server := NewServer(secret, mockClient)

	// Create a test HTTP request with invalid secret
	req := httptest.NewRequest("GET", "/?secret=invalidsecret", nil)
	w := httptest.NewRecorder()

	server.HandleConnection(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Invalid secret") {
		t.Errorf("Expected 'Invalid secret' in response body, got: %s", body)
	}
}

func TestHandleConnection_MissingSecret(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{}
	server := NewServer(secret, mockClient)

	// Create a test HTTP request without secret
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.HandleConnection(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandleMessage_Ping(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{}
	server := NewServer(secret, mockClient)

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    protocol.TypePing,
		Payload: protocol.PingPayload{},
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypePong {
		t.Errorf("Expected response type %s, got %s", protocol.TypePong, response.Type)
	}

	if response.ID != msg.ID {
		t.Errorf("Expected response ID %s, got %s", msg.ID, response.ID)
	}
}

func TestHandleMessage_UnknownType(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{}
	server := NewServer(secret, mockClient)

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    "unknown",
		Payload: nil,
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypeError {
		t.Errorf("Expected response type %s, got %s", protocol.TypeError, response.Type)
	}

	errorPayload, ok := response.Payload.(protocol.ErrorPayload)
	if !ok {
		t.Fatal("Expected ErrorPayload in response")
	}

	if errorPayload.Code != "INVALID_MESSAGE_TYPE" {
		t.Errorf("Expected error code INVALID_MESSAGE_TYPE, got %s", errorPayload.Code)
	}
}

func TestHandleQuery_Success(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{
		ExecuteQueryFunc: func(ctx context.Context, sql string, params []interface{}) (*postgres.QueryResult, error) {
			return &postgres.QueryResult{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "test"},
				},
				Columns: []protocol.ColumnInfo{
					{Name: "id", DataType: "int4"},
					{Name: "name", DataType: "text"},
				},
				RowCount:      1,
				ExecutionTime: 10 * time.Millisecond,
			}, nil
		},
	}
	server := NewServer(secret, mockClient)

	payload := protocol.QueryPayload{
		SQL:    "SELECT * FROM users",
		Params: []interface{}{},
	}

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    protocol.TypeQuery,
		Payload: payload,
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypeResult {
		t.Errorf("Expected response type %s, got %s", protocol.TypeResult, response.Type)
	}

	// Marshal and unmarshal to get the result payload
	payloadBytes, _ := json.Marshal(response.Payload)
	var resultPayload protocol.ResultPayload
	json.Unmarshal(payloadBytes, &resultPayload)

	if resultPayload.RowCount != 1 {
		t.Errorf("Expected row count 1, got %d", resultPayload.RowCount)
	}

	if len(resultPayload.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(resultPayload.Columns))
	}
}

func TestHandleQuery_EmptySQL(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{}
	server := NewServer(secret, mockClient)

	payload := protocol.QueryPayload{
		SQL:    "",
		Params: []interface{}{},
	}

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    protocol.TypeQuery,
		Payload: payload,
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypeError {
		t.Errorf("Expected response type %s, got %s", protocol.TypeError, response.Type)
	}

	payloadBytes, _ := json.Marshal(response.Payload)
	var errorPayload protocol.ErrorPayload
	json.Unmarshal(payloadBytes, &errorPayload)

	if errorPayload.Code != "EMPTY_QUERY" {
		t.Errorf("Expected error code EMPTY_QUERY, got %s", errorPayload.Code)
	}
}

func TestHandleQuery_DatabaseError(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{
		ExecuteQueryFunc: func(ctx context.Context, sql string, params []interface{}) (*postgres.QueryResult, error) {
			return nil, fmt.Errorf("table does not exist")
		},
	}
	server := NewServer(secret, mockClient)

	payload := protocol.QueryPayload{
		SQL:    "SELECT * FROM nonexistent",
		Params: []interface{}{},
	}

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    protocol.TypeQuery,
		Payload: payload,
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypeError {
		t.Errorf("Expected response type %s, got %s", protocol.TypeError, response.Type)
	}

	payloadBytes, _ := json.Marshal(response.Payload)
	var errorPayload protocol.ErrorPayload
	json.Unmarshal(payloadBytes, &errorPayload)

	if errorPayload.Code != "QUERY_ERROR" {
		t.Errorf("Expected error code QUERY_ERROR, got %s", errorPayload.Code)
	}

	if !strings.Contains(errorPayload.Message, "table does not exist") {
		t.Errorf("Expected error message to contain 'table does not exist', got: %s", errorPayload.Message)
	}
}

func TestHandleQuery_WithTimeout(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{
		ExecuteQueryFunc: func(ctx context.Context, sql string, params []interface{}) (*postgres.QueryResult, error) {
			// Check that context has a deadline
			if _, ok := ctx.Deadline(); !ok {
				t.Error("Expected context to have a deadline")
			}

			return &postgres.QueryResult{
				Rows:          []map[string]interface{}{},
				Columns:       []protocol.ColumnInfo{},
				RowCount:      0,
				ExecutionTime: 5 * time.Millisecond,
			}, nil
		},
	}
	server := NewServer(secret, mockClient)

	payload := protocol.QueryPayload{
		SQL:     "SELECT 1",
		Params:  []interface{}{},
		Timeout: 5000, // 5 seconds
	}

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    protocol.TypeQuery,
		Payload: payload,
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypeResult {
		t.Errorf("Expected response type %s, got %s", protocol.TypeResult, response.Type)
	}
}

func TestHandleIntrospect_Success(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{
		IntrospectSchemaFunc: func(ctx context.Context) (*protocol.SchemaPayload, error) {
			return &protocol.SchemaPayload{
				Tables: []protocol.TableInfo{
					{
						Schema: "public",
						Name:   "users",
						Type:   "table",
						Columns: []protocol.ColumnInfo{
							{Name: "id", DataType: "int4"},
							{Name: "name", DataType: "text"},
						},
					},
				},
				Functions: []protocol.FunctionInfo{
					{
						Schema:     "public",
						Name:       "get_user",
						ReturnType: "users",
					},
				},
			}, nil
		},
	}
	server := NewServer(secret, mockClient)

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    protocol.TypeIntrospect,
		Payload: nil,
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypeSchema {
		t.Errorf("Expected response type %s, got %s", protocol.TypeSchema, response.Type)
	}

	payloadBytes, _ := json.Marshal(response.Payload)
	var schemaPayload protocol.SchemaPayload
	json.Unmarshal(payloadBytes, &schemaPayload)

	if len(schemaPayload.Tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(schemaPayload.Tables))
	}

	if len(schemaPayload.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(schemaPayload.Functions))
	}
}

func TestHandleIntrospect_Error(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{
		IntrospectSchemaFunc: func(ctx context.Context) (*protocol.SchemaPayload, error) {
			return nil, fmt.Errorf("connection lost")
		},
	}
	server := NewServer(secret, mockClient)

	msg := protocol.ClientMessage{
		ID:      "test-1",
		Type:    protocol.TypeIntrospect,
		Payload: nil,
	}

	response := server.handleMessage(msg)

	if response.Type != protocol.TypeError {
		t.Errorf("Expected response type %s, got %s", protocol.TypeError, response.Type)
	}

	payloadBytes, _ := json.Marshal(response.Payload)
	var errorPayload protocol.ErrorPayload
	json.Unmarshal(payloadBytes, &errorPayload)

	if errorPayload.Code != "INTROSPECTION_ERROR" {
		t.Errorf("Expected error code INTROSPECTION_ERROR, got %s", errorPayload.Code)
	}
}

func TestHandleConnection_ValidWebSocket(t *testing.T) {
	secret, err := auth.GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	mockClient := &MockPostgresClient{}
	server := NewServer(secret, mockClient)

	// Create a test HTTP server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.HandleConnection(w, r)
	}))
	defer testServer.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http") + "?secret=" + secret

	// Connect to the WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Send a ping message
	pingMsg := protocol.ClientMessage{
		ID:      "test-ping",
		Type:    protocol.TypePing,
		Payload: protocol.PingPayload{},
	}

	if err := ws.WriteJSON(pingMsg); err != nil {
		t.Fatalf("Failed to send ping message: %v", err)
	}

	// Read the response
	var response protocol.ServerMessage
	if err := ws.ReadJSON(&response); err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if response.Type != protocol.TypePong {
		t.Errorf("Expected response type %s, got %s", protocol.TypePong, response.Type)
	}

	if response.ID != pingMsg.ID {
		t.Errorf("Expected response ID %s, got %s", pingMsg.ID, response.ID)
	}
}

func TestSendMessage(t *testing.T) {
	// Create a mock WebSocket connection using httptest
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		// Send a test message
		msg := protocol.NewPong("test-id")
		if err := SendMessage(conn, msg); err != nil {
			t.Errorf("SendMessage failed: %v", err)
		}
	}))
	defer testServer.Close()

	// Connect to the test server
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer ws.Close()

	// Read the message
	var response protocol.ServerMessage
	if err := ws.ReadJSON(&response); err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	if response.Type != protocol.TypePong {
		t.Errorf("Expected message type %s, got %s", protocol.TypePong, response.Type)
	}
}
