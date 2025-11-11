package protocol

import "time"

// Message types
const (
	// Client -> Server
	TypeQuery      = "query"
	TypeIntrospect = "introspect"
	TypePing       = "ping"

	// Server -> Client
	TypeResult = "result"
	TypeError  = "error"
	TypeSchema = "schema"
	TypePong   = "pong"
)

// Message is the base structure for all messages
type Message struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// ClientMessage represents messages from the client
type ClientMessage struct {
	ID      string       `json:"id"`
	Type    string       `json:"type"`
	Payload QueryPayload `json:"payload"`
}

// ServerMessage represents messages from the server
type ServerMessage struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// QueryPayload contains query execution details
type QueryPayload struct {
	SQL     string        `json:"sql"`
	Params  []interface{} `json:"params,omitempty"`
	Timeout int           `json:"timeout,omitempty"` // milliseconds
}

// ResultPayload contains query results
type ResultPayload struct {
	Rows          []map[string]interface{} `json:"rows"`
	Columns       []ColumnInfo             `json:"columns"`
	RowCount      int                      `json:"rowCount"`
	ExecutionTime int64                    `json:"executionTime"` // milliseconds
}

// ColumnInfo describes a result column
type ColumnInfo struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
	TableOID uint32 `json:"tableOid,omitempty"`
}

// ErrorPayload contains error details
type ErrorPayload struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Detail   string `json:"detail,omitempty"`
	Position int    `json:"position,omitempty"`
}

// SchemaPayload contains database schema information
type SchemaPayload struct {
	Tables    []TableInfo    `json:"tables"`
	Functions []FunctionInfo `json:"functions"`
}

// TableInfo describes a database table
type TableInfo struct {
	Schema  string       `json:"schema"`
	Name    string       `json:"name"`
	Type    string       `json:"type"` // 'r' = table, 'v' = view, 'm' = materialized view
	Columns []ColumnInfo `json:"columns"`
}

// FunctionInfo describes a database function
type FunctionInfo struct {
	Schema     string `json:"schema"`
	Name       string `json:"name"`
	ReturnType string `json:"returnType"`
}

// NewQueryResult creates a result message
func NewQueryResult(id string, rows []map[string]interface{}, columns []ColumnInfo, executionTime time.Duration) ServerMessage {
	return ServerMessage{
		ID:   id,
		Type: TypeResult,
		Payload: ResultPayload{
			Rows:          rows,
			Columns:       columns,
			RowCount:      len(rows),
			ExecutionTime: executionTime.Milliseconds(),
		},
	}
}

// NewError creates an error message
func NewError(id string, code, message, detail string) ServerMessage {
	return ServerMessage{
		ID:   id,
		Type: TypeError,
		Payload: ErrorPayload{
			Code:    code,
			Message: message,
			Detail:  detail,
		},
	}
}

// NewSchemaResult creates a schema message
func NewSchemaResult(id string, tables []TableInfo, functions []FunctionInfo) ServerMessage {
	return ServerMessage{
		ID:   id,
		Type: TypeSchema,
		Payload: SchemaPayload{
			Tables:    tables,
			Functions: functions,
		},
	}
}

// NewPong creates a pong message
func NewPong(id string) ServerMessage {
	return ServerMessage{
		ID:      id,
		Type:    TypePong,
		Payload: nil,
	}
}
