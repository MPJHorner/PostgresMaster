package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageSerialization(t *testing.T) {
	tests := []struct {
		name    string
		message ServerMessage
		want    string
	}{
		{
			name: "query result message",
			message: ServerMessage{
				ID:   "test-id-1",
				Type: TypeResult,
				Payload: ResultPayload{
					Rows: []map[string]interface{}{
						{"id": 1, "name": "Alice"},
						{"id": 2, "name": "Bob"},
					},
					Columns: []ColumnInfo{
						{Name: "id", DataType: "integer", TypeOID: 23, Nullable: false},
						{Name: "name", DataType: "text", TypeOID: 25, Nullable: true},
					},
					RowCount:      2,
					ExecutionTime: 45,
				},
			},
		},
		{
			name: "error message",
			message: ServerMessage{
				ID:   "test-id-2",
				Type: TypeError,
				Payload: ErrorPayload{
					Code:     "42P01",
					Message:  "relation \"users\" does not exist",
					Detail:   "The table users was not found in the database",
					Hint:     "Check your schema and table name",
					Position: 15,
				},
			},
		},
		{
			name: "schema message",
			message: ServerMessage{
				ID:   "test-id-3",
				Type: TypeSchema,
				Payload: SchemaPayload{
					Tables: []TableInfo{
						{
							Schema: "public",
							Name:   "users",
							Type:   "r",
							Columns: []ColumnInfo{
								{Name: "id", DataType: "integer"},
								{Name: "email", DataType: "text"},
							},
						},
					},
					Functions: []FunctionInfo{
						{
							Schema:     "public",
							Name:       "get_user",
							ReturnType: "record",
						},
					},
				},
			},
		},
		{
			name: "pong message",
			message: ServerMessage{
				ID:   "test-id-4",
				Type: TypePong,
				Payload: PongPayload{
					Timestamp: time.Date(2025, 11, 11, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize to JSON
			data, err := json.Marshal(tt.message)
			if err != nil {
				t.Fatalf("Failed to marshal message: %v", err)
			}

			// Verify it's valid JSON
			if !json.Valid(data) {
				t.Error("Produced invalid JSON")
			}

			// Deserialize back
			var result ServerMessage
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal message: %v", err)
			}

			// Verify basic fields
			if result.ID != tt.message.ID {
				t.Errorf("ID mismatch: got %s, want %s", result.ID, tt.message.ID)
			}
			if result.Type != tt.message.Type {
				t.Errorf("Type mismatch: got %s, want %s", result.Type, tt.message.Type)
			}
		})
	}
}

func TestClientMessageSerialization(t *testing.T) {
	tests := []struct {
		name    string
		message ClientMessage
	}{
		{
			name: "query message",
			message: ClientMessage{
				ID:   "client-1",
				Type: TypeQuery,
				Payload: QueryPayload{
					SQL:     "SELECT * FROM users WHERE id = $1",
					Params:  []interface{}{123},
					Timeout: 30000,
				},
			},
		},
		{
			name: "introspect message",
			message: ClientMessage{
				ID:      "client-2",
				Type:    TypeIntrospect,
				Payload: struct{}{},
			},
		},
		{
			name: "ping message",
			message: ClientMessage{
				ID:      "client-3",
				Type:    TypePing,
				Payload: PingPayload{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize to JSON
			data, err := json.Marshal(tt.message)
			if err != nil {
				t.Fatalf("Failed to marshal message: %v", err)
			}

			// Verify it's valid JSON
			if !json.Valid(data) {
				t.Error("Produced invalid JSON")
			}

			// Deserialize back
			var result ClientMessage
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal message: %v", err)
			}

			// Verify basic fields
			if result.ID != tt.message.ID {
				t.Errorf("ID mismatch: got %s, want %s", result.ID, tt.message.ID)
			}
			if result.Type != tt.message.Type {
				t.Errorf("Type mismatch: got %s, want %s", result.Type, tt.message.Type)
			}
		})
	}
}

func TestQueryPayloadSerialization(t *testing.T) {
	payload := QueryPayload{
		SQL:     "SELECT * FROM users",
		Params:  []interface{}{1, "test", true},
		Timeout: 5000,
	}

	// Serialize
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal QueryPayload: %v", err)
	}

	// Deserialize
	var result QueryPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal QueryPayload: %v", err)
	}

	// Verify fields
	if result.SQL != payload.SQL {
		t.Errorf("SQL mismatch: got %s, want %s", result.SQL, payload.SQL)
	}
	if result.Timeout != payload.Timeout {
		t.Errorf("Timeout mismatch: got %d, want %d", result.Timeout, payload.Timeout)
	}
	if len(result.Params) != len(payload.Params) {
		t.Errorf("Params length mismatch: got %d, want %d", len(result.Params), len(payload.Params))
	}
}

func TestResultPayloadSerialization(t *testing.T) {
	payload := ResultPayload{
		Rows: []map[string]interface{}{
			{"id": float64(1), "name": "Alice", "active": true},
			{"id": float64(2), "name": "Bob", "active": false},
		},
		Columns: []ColumnInfo{
			{Name: "id", DataType: "integer", TypeOID: 23, Nullable: false},
			{Name: "name", DataType: "text", TypeOID: 25, Nullable: true},
			{Name: "active", DataType: "boolean", TypeOID: 16, Nullable: false},
		},
		RowCount:      2,
		ExecutionTime: 123,
	}

	// Serialize
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal ResultPayload: %v", err)
	}

	// Deserialize
	var result ResultPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal ResultPayload: %v", err)
	}

	// Verify fields
	if result.RowCount != payload.RowCount {
		t.Errorf("RowCount mismatch: got %d, want %d", result.RowCount, payload.RowCount)
	}
	if result.ExecutionTime != payload.ExecutionTime {
		t.Errorf("ExecutionTime mismatch: got %d, want %d", result.ExecutionTime, payload.ExecutionTime)
	}
	if len(result.Rows) != len(payload.Rows) {
		t.Errorf("Rows length mismatch: got %d, want %d", len(result.Rows), len(payload.Rows))
	}
	if len(result.Columns) != len(payload.Columns) {
		t.Errorf("Columns length mismatch: got %d, want %d", len(result.Columns), len(payload.Columns))
	}

	// Verify column details
	for i, col := range result.Columns {
		if col.Name != payload.Columns[i].Name {
			t.Errorf("Column %d name mismatch: got %s, want %s", i, col.Name, payload.Columns[i].Name)
		}
		if col.DataType != payload.Columns[i].DataType {
			t.Errorf("Column %d dataType mismatch: got %s, want %s", i, col.DataType, payload.Columns[i].DataType)
		}
		if col.TypeOID != payload.Columns[i].TypeOID {
			t.Errorf("Column %d typeOID mismatch: got %d, want %d", i, col.TypeOID, payload.Columns[i].TypeOID)
		}
		if col.Nullable != payload.Columns[i].Nullable {
			t.Errorf("Column %d nullable mismatch: got %v, want %v", i, col.Nullable, payload.Columns[i].Nullable)
		}
	}
}

func TestErrorPayloadSerialization(t *testing.T) {
	payload := ErrorPayload{
		Code:     "42P01",
		Message:  "relation does not exist",
		Detail:   "Table users not found",
		Hint:     "Check the table name",
		Position: 25,
	}

	// Serialize
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorPayload: %v", err)
	}

	// Deserialize
	var result ErrorPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal ErrorPayload: %v", err)
	}

	// Verify all fields
	if result.Code != payload.Code {
		t.Errorf("Code mismatch: got %s, want %s", result.Code, payload.Code)
	}
	if result.Message != payload.Message {
		t.Errorf("Message mismatch: got %s, want %s", result.Message, payload.Message)
	}
	if result.Detail != payload.Detail {
		t.Errorf("Detail mismatch: got %s, want %s", result.Detail, payload.Detail)
	}
	if result.Hint != payload.Hint {
		t.Errorf("Hint mismatch: got %s, want %s", result.Hint, payload.Hint)
	}
	if result.Position != payload.Position {
		t.Errorf("Position mismatch: got %d, want %d", result.Position, payload.Position)
	}
}

func TestSchemaPayloadSerialization(t *testing.T) {
	payload := SchemaPayload{
		Tables: []TableInfo{
			{
				Schema: "public",
				Name:   "users",
				Type:   "r",
				Columns: []ColumnInfo{
					{Name: "id", DataType: "integer", TypeOID: 23, Nullable: false},
					{Name: "email", DataType: "text", TypeOID: 25, Nullable: true},
				},
			},
			{
				Schema: "public",
				Name:   "posts",
				Type:   "r",
				Columns: []ColumnInfo{
					{Name: "id", DataType: "integer", TypeOID: 23, Nullable: false},
					{Name: "title", DataType: "text", TypeOID: 25, Nullable: false},
				},
			},
		},
		Functions: []FunctionInfo{
			{Schema: "public", Name: "get_user", ReturnType: "record"},
			{Schema: "public", Name: "count_posts", ReturnType: "integer"},
		},
	}

	// Serialize
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal SchemaPayload: %v", err)
	}

	// Deserialize
	var result SchemaPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal SchemaPayload: %v", err)
	}

	// Verify tables
	if len(result.Tables) != len(payload.Tables) {
		t.Errorf("Tables length mismatch: got %d, want %d", len(result.Tables), len(payload.Tables))
	}
	for i, table := range result.Tables {
		if table.Name != payload.Tables[i].Name {
			t.Errorf("Table %d name mismatch: got %s, want %s", i, table.Name, payload.Tables[i].Name)
		}
		if table.Schema != payload.Tables[i].Schema {
			t.Errorf("Table %d schema mismatch: got %s, want %s", i, table.Schema, payload.Tables[i].Schema)
		}
		if len(table.Columns) != len(payload.Tables[i].Columns) {
			t.Errorf("Table %d columns length mismatch: got %d, want %d", i, len(table.Columns), len(payload.Tables[i].Columns))
		}
	}

	// Verify functions
	if len(result.Functions) != len(payload.Functions) {
		t.Errorf("Functions length mismatch: got %d, want %d", len(result.Functions), len(payload.Functions))
	}
	for i, fn := range result.Functions {
		if fn.Name != payload.Functions[i].Name {
			t.Errorf("Function %d name mismatch: got %s, want %s", i, fn.Name, payload.Functions[i].Name)
		}
		if fn.ReturnType != payload.Functions[i].ReturnType {
			t.Errorf("Function %d returnType mismatch: got %s, want %s", i, fn.ReturnType, payload.Functions[i].ReturnType)
		}
	}
}

func TestPongPayloadSerialization(t *testing.T) {
	now := time.Now()
	payload := PongPayload{
		Timestamp: now,
	}

	// Serialize
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal PongPayload: %v", err)
	}

	// Deserialize
	var result PongPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal PongPayload: %v", err)
	}

	// Verify timestamp (allowing for small differences due to JSON precision)
	if result.Timestamp.Unix() != now.Unix() {
		t.Errorf("Timestamp mismatch: got %v, want %v", result.Timestamp, now)
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("NewQueryResult", func(t *testing.T) {
		rows := []map[string]interface{}{
			{"id": 1, "name": "test"},
		}
		columns := []ColumnInfo{
			{Name: "id", DataType: "integer"},
			{Name: "name", DataType: "text"},
		}
		duration := 100 * time.Millisecond

		msg := NewQueryResult("test-id", rows, columns, duration)

		if msg.ID != "test-id" {
			t.Errorf("ID mismatch: got %s, want test-id", msg.ID)
		}
		if msg.Type != TypeResult {
			t.Errorf("Type mismatch: got %s, want %s", msg.Type, TypeResult)
		}

		payload, ok := msg.Payload.(ResultPayload)
		if !ok {
			t.Fatal("Payload is not ResultPayload")
		}
		if payload.RowCount != 1 {
			t.Errorf("RowCount mismatch: got %d, want 1", payload.RowCount)
		}
		if payload.ExecutionTime != 100 {
			t.Errorf("ExecutionTime mismatch: got %d, want 100", payload.ExecutionTime)
		}
	})

	t.Run("NewError", func(t *testing.T) {
		msg := NewError("test-id", "42P01", "table not found", "check schema")

		if msg.ID != "test-id" {
			t.Errorf("ID mismatch: got %s, want test-id", msg.ID)
		}
		if msg.Type != TypeError {
			t.Errorf("Type mismatch: got %s, want %s", msg.Type, TypeError)
		}

		payload, ok := msg.Payload.(ErrorPayload)
		if !ok {
			t.Fatal("Payload is not ErrorPayload")
		}
		if payload.Code != "42P01" {
			t.Errorf("Code mismatch: got %s, want 42P01", payload.Code)
		}
		if payload.Message != "table not found" {
			t.Errorf("Message mismatch: got %s, want 'table not found'", payload.Message)
		}
	})

	t.Run("NewSchemaResult", func(t *testing.T) {
		tables := []TableInfo{
			{Schema: "public", Name: "users", Type: "r"},
		}
		functions := []FunctionInfo{
			{Schema: "public", Name: "get_user", ReturnType: "record"},
		}

		msg := NewSchemaResult("test-id", tables, functions)

		if msg.ID != "test-id" {
			t.Errorf("ID mismatch: got %s, want test-id", msg.ID)
		}
		if msg.Type != TypeSchema {
			t.Errorf("Type mismatch: got %s, want %s", msg.Type, TypeSchema)
		}

		payload, ok := msg.Payload.(SchemaPayload)
		if !ok {
			t.Fatal("Payload is not SchemaPayload")
		}
		if len(payload.Tables) != 1 {
			t.Errorf("Tables length mismatch: got %d, want 1", len(payload.Tables))
		}
		if len(payload.Functions) != 1 {
			t.Errorf("Functions length mismatch: got %d, want 1", len(payload.Functions))
		}
	})

	t.Run("NewPong", func(t *testing.T) {
		msg := NewPong("test-id")

		if msg.ID != "test-id" {
			t.Errorf("ID mismatch: got %s, want test-id", msg.ID)
		}
		if msg.Type != TypePong {
			t.Errorf("Type mismatch: got %s, want %s", msg.Type, TypePong)
		}

		payload, ok := msg.Payload.(PongPayload)
		if !ok {
			t.Fatal("Payload is not PongPayload")
		}
		// Just verify timestamp is recent (within last second)
		if time.Since(payload.Timestamp) > time.Second {
			t.Error("Timestamp is not recent")
		}
	})
}

func TestJSONOmitEmpty(t *testing.T) {
	t.Run("QueryPayload omits empty fields", func(t *testing.T) {
		payload := QueryPayload{
			SQL: "SELECT 1",
		}
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		// Verify params and timeout are not in JSON
		jsonStr := string(data)
		if contains(jsonStr, "params") {
			t.Error("JSON should not contain 'params' field when empty")
		}
		if contains(jsonStr, "timeout") {
			t.Error("JSON should not contain 'timeout' field when zero")
		}
	})

	t.Run("ErrorPayload omits empty fields", func(t *testing.T) {
		payload := ErrorPayload{
			Code:    "ERROR",
			Message: "Something went wrong",
		}
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		// Verify optional fields are not in JSON
		jsonStr := string(data)
		if contains(jsonStr, "detail") {
			t.Error("JSON should not contain 'detail' field when empty")
		}
		if contains(jsonStr, "hint") {
			t.Error("JSON should not contain 'hint' field when empty")
		}
		if contains(jsonStr, "position") {
			t.Error("JSON should not contain 'position' field when zero")
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)+1 && containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
