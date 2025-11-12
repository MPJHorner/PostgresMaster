/**
 * WebSocket Protocol Types
 * These types match the Go protocol types in proxy/pkg/protocol/messages.go
 */

// Message type constants
export const MessageType = {
	// Client -> Server
	QUERY: 'query',
	INTROSPECT: 'introspect',
	PING: 'ping',

	// Server -> Client
	RESULT: 'result',
	ERROR: 'error',
	SCHEMA: 'schema',
	PONG: 'pong'
} as const;

export type MessageTypeValue = (typeof MessageType)[keyof typeof MessageType];

/**
 * Base message structure for all WebSocket messages
 */
export interface Message {
	id: string;
	type: MessageTypeValue;
	payload: unknown;
}

/**
 * ClientMessage represents messages sent from the client to the server
 */
export interface ClientMessage<T = unknown> {
	id: string;
	type: MessageTypeValue;
	payload: T;
}

/**
 * ServerMessage represents messages sent from the server to the client
 */
export interface ServerMessage<T = unknown> {
	id: string;
	type: MessageTypeValue;
	payload: T;
}

/**
 * QueryPayload contains query execution details
 */
export interface QueryPayload {
	sql: string;
	params?: unknown[];
	timeout?: number; // milliseconds
}

/**
 * ColumnInfo describes a result column
 */
export interface ColumnInfo {
	name: string;
	dataType: string;
	typeOid?: number;
	nullable?: boolean;
}

/**
 * ResultPayload contains query results
 */
export interface ResultPayload {
	rows: Record<string, unknown>[];
	columns: ColumnInfo[];
	rowCount: number;
	executionTime: number; // milliseconds
}

/**
 * ErrorPayload contains error details
 */
export interface ErrorPayload {
	code: string;
	message: string;
	detail?: string;
	hint?: string;
	position?: number;
}

/**
 * FunctionInfo describes a database function
 */
export interface FunctionInfo {
	schema: string;
	name: string;
	returnType: string;
}

/**
 * TableInfo describes a database table
 */
export interface TableInfo {
	schema: string;
	name: string;
	type: string; // 'r' = table, 'v' = view, 'm' = materialized view
	columns: ColumnInfo[];
}

/**
 * SchemaPayload contains database schema information
 */
export interface SchemaPayload {
	tables: TableInfo[];
	functions: FunctionInfo[];
}

/**
 * PingPayload represents a ping request (empty)
 */
export interface PingPayload {
	// Empty payload
}

/**
 * PongPayload represents a pong response
 */
export interface PongPayload {
	timestamp: string; // ISO 8601 timestamp
}

/**
 * Helper type for query messages
 */
export type QueryMessage = ClientMessage<QueryPayload>;

/**
 * Helper type for introspect messages
 */
export type IntrospectMessage = ClientMessage<Record<string, never>>;

/**
 * Helper type for ping messages
 */
export type PingMessage = ClientMessage<PingPayload>;

/**
 * Helper type for result messages
 */
export type ResultMessage = ServerMessage<ResultPayload>;

/**
 * Helper type for error messages
 */
export type ErrorMessage = ServerMessage<ErrorPayload>;

/**
 * Helper type for schema messages
 */
export type SchemaMessage = ServerMessage<SchemaPayload>;

/**
 * Helper type for pong messages
 */
export type PongMessage = ServerMessage<PongPayload>;

/**
 * Creates a new query message
 */
export function createQueryMessage(id: string, sql: string, params?: unknown[]): QueryMessage {
	return {
		id,
		type: MessageType.QUERY,
		payload: {
			sql,
			params
		}
	};
}

/**
 * Creates a new introspect message
 */
export function createIntrospectMessage(id: string): IntrospectMessage {
	return {
		id,
		type: MessageType.INTROSPECT,
		payload: {}
	};
}

/**
 * Creates a new ping message
 */
export function createPingMessage(id: string): PingMessage {
	return {
		id,
		type: MessageType.PING,
		payload: {}
	};
}

/**
 * Type guard to check if a message is a result message
 */
export function isResultMessage(message: ServerMessage): message is ResultMessage {
	return message.type === MessageType.RESULT;
}

/**
 * Type guard to check if a message is an error message
 */
export function isErrorMessage(message: ServerMessage): message is ErrorMessage {
	return message.type === MessageType.ERROR;
}

/**
 * Type guard to check if a message is a schema message
 */
export function isSchemaMessage(message: ServerMessage): message is SchemaMessage {
	return message.type === MessageType.SCHEMA;
}

/**
 * Type guard to check if a message is a pong message
 */
export function isPongMessage(message: ServerMessage): message is PongMessage {
	return message.type === MessageType.PONG;
}
