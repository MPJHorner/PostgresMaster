/**
 * WebSocket Client Service
 * Handles communication with the Go proxy server via WebSocket protocol
 */

import { v4 as uuidv4 } from 'uuid';
import {
	type ServerMessage,
	type ResultPayload,
	type ErrorPayload,
	type SchemaPayload,
	createQueryMessage,
	createIntrospectMessage,
	createPingMessage,
	isResultMessage,
	isErrorMessage,
	isSchemaMessage,
	isPongMessage
} from './protocol';

/**
 * Connection state for the WebSocket client
 */
export enum ConnectionState {
	DISCONNECTED = 'disconnected',
	CONNECTING = 'connecting',
	CONNECTED = 'connected',
	RECONNECTING = 'reconnecting',
	ERROR = 'error'
}

/**
 * Options for PostgresProxyClient
 */
export interface ClientOptions {
	/** WebSocket URL (default: ws://localhost:8080) */
	url?: string;
	/** Authentication secret */
	secret: string;
	/** Reconnection attempts (default: 3) */
	maxReconnectAttempts?: number;
	/** Reconnection delay in ms (default: 2000) */
	reconnectDelay?: number;
	/** Connection timeout in ms (default: 10000) */
	connectionTimeout?: number;
}

/**
 * Pending request tracking
 */
interface PendingRequest {
	resolve: (value: unknown) => void;
	reject: (reason: Error) => void;
	timeout: ReturnType<typeof setTimeout>;
}

/**
 * PostgresProxyClient handles WebSocket communication with the Go proxy
 *
 * @example
 * ```typescript
 * const client = new PostgresProxyClient({ secret: 'abc123' });
 * await client.connect();
 *
 * const result = await client.executeQuery('SELECT * FROM users');
 * console.log(result.rows);
 *
 * await client.close();
 * ```
 */
export class PostgresProxyClient {
	private ws: WebSocket | null = null;
	private state: ConnectionState = ConnectionState.DISCONNECTED;
	private pendingRequests = new Map<string, PendingRequest>();
	private reconnectAttempts = 0;
	private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

	private readonly url: string;
	private readonly secret: string;
	private readonly maxReconnectAttempts: number;
	private readonly reconnectDelay: number;
	private readonly connectionTimeout: number;

	// Event callbacks
	private onStateChangeCallback?: (state: ConnectionState) => void;
	private onErrorCallback?: (error: Error) => void;

	/**
	 * Creates a new PostgresProxyClient
	 * @param options Client configuration options
	 */
	constructor(options: ClientOptions) {
		this.url = options.url || 'ws://localhost:8080';
		this.secret = options.secret;
		this.maxReconnectAttempts = options.maxReconnectAttempts ?? 3;
		this.reconnectDelay = options.reconnectDelay ?? 2000;
		this.connectionTimeout = options.connectionTimeout ?? 10000;
	}

	/**
	 * Gets the current connection state
	 */
	public getState(): ConnectionState {
		return this.state;
	}

	/**
	 * Checks if the client is connected
	 */
	public isConnected(): boolean {
		return this.state === ConnectionState.CONNECTED && this.ws?.readyState === WebSocket.OPEN;
	}

	/**
	 * Registers a callback for state changes
	 */
	public onStateChange(callback: (state: ConnectionState) => void): void {
		this.onStateChangeCallback = callback;
	}

	/**
	 * Registers a callback for errors
	 */
	public onError(callback: (error: Error) => void): void {
		this.onErrorCallback = callback;
	}

	/**
	 * Connects to the WebSocket server
	 * @throws {Error} If connection fails
	 */
	public async connect(): Promise<void> {
		if (this.state === ConnectionState.CONNECTED || this.state === ConnectionState.CONNECTING) {
			return;
		}

		this.setState(ConnectionState.CONNECTING);

		return new Promise((resolve, reject) => {
			try {
				// Add secret as URL parameter
				const wsUrl = `${this.url}?secret=${encodeURIComponent(this.secret)}`;
				this.ws = new WebSocket(wsUrl);

				// Set connection timeout
				const timeoutId = setTimeout(() => {
					if (this.state === ConnectionState.CONNECTING) {
						this.ws?.close();
						const error = new Error('Connection timeout');
						this.handleError(error);
						reject(error);
					}
				}, this.connectionTimeout);

				this.ws.onopen = () => {
					clearTimeout(timeoutId);
					this.reconnectAttempts = 0;
					this.setState(ConnectionState.CONNECTED);
					resolve();
				};

				this.ws.onmessage = (event) => {
					this.handleMessage(event.data);
				};

				this.ws.onerror = () => {
					clearTimeout(timeoutId);
					const error = new Error('WebSocket error');
					this.handleError(error);
					reject(error);
				};

				this.ws.onclose = (event) => {
					clearTimeout(timeoutId);
					this.handleClose(event);

					// Only reject if we're still in connecting state
					if (this.state === ConnectionState.CONNECTING) {
						reject(new Error(`Connection closed: ${event.code} ${event.reason}`));
					}
				};
			} catch (error) {
				this.handleError(error instanceof Error ? error : new Error(String(error)));
				reject(error);
			}
		});
	}

	/**
	 * Executes a SQL query
	 * @param sql SQL query string
	 * @param params Optional query parameters
	 * @param timeout Optional query timeout in milliseconds
	 * @returns Query result payload
	 * @throws {Error} If query execution fails
	 */
	public async executeQuery(
		sql: string,
		params?: unknown[],
		timeout?: number
	): Promise<ResultPayload> {
		if (!this.isConnected()) {
			throw new Error('Not connected to proxy server');
		}

		const id = uuidv4();
		const message = createQueryMessage(id, sql, params);

		return this.sendRequest<ResultPayload>(message, timeout || 30000);
	}

	/**
	 * Introspects the database schema
	 * @returns Schema information
	 * @throws {Error} If introspection fails
	 */
	public async introspectSchema(): Promise<SchemaPayload> {
		if (!this.isConnected()) {
			throw new Error('Not connected to proxy server');
		}

		const id = uuidv4();
		const message = createIntrospectMessage(id);

		return this.sendRequest<SchemaPayload>(message, 30000);
	}

	/**
	 * Sends a ping to the server
	 * @returns Pong response timestamp
	 * @throws {Error} If ping fails
	 */
	public async ping(): Promise<string> {
		if (!this.isConnected()) {
			throw new Error('Not connected to proxy server');
		}

		const id = uuidv4();
		const message = createPingMessage(id);

		const response = await this.sendRequest<{ timestamp: string }>(message, 5000);
		return response.timestamp;
	}

	/**
	 * Closes the WebSocket connection
	 */
	public close(): void {
		// Clear reconnection timer
		if (this.reconnectTimer) {
			clearTimeout(this.reconnectTimer);
			this.reconnectTimer = null;
		}

		// Reject all pending requests
		this.pendingRequests.forEach((request) => {
			clearTimeout(request.timeout);
			request.reject(new Error('Connection closed'));
		});
		this.pendingRequests.clear();

		// Close WebSocket
		if (this.ws) {
			this.ws.close(1000, 'Client closed connection');
			this.ws = null;
		}

		this.setState(ConnectionState.DISCONNECTED);
	}

	/**
	 * Sends a request and waits for a response
	 */
	private sendRequest<T>(message: unknown, timeout: number): Promise<T> {
		return new Promise((resolve, reject) => {
			if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
				reject(new Error('WebSocket not connected'));
				return;
			}

			const msg = message as { id: string };
			const messageId = msg.id;

			// Set up timeout
			const timeoutId = setTimeout(() => {
				this.pendingRequests.delete(messageId);
				reject(new Error(`Request timeout after ${timeout}ms`));
			}, timeout);

			// Store pending request
			this.pendingRequests.set(messageId, {
				resolve: resolve as (value: unknown) => void,
				reject,
				timeout: timeoutId
			});

			// Send message
			try {
				this.ws.send(JSON.stringify(message));
			} catch (error) {
				this.pendingRequests.delete(messageId);
				clearTimeout(timeoutId);
				reject(error);
			}
		});
	}

	/**
	 * Handles incoming WebSocket messages
	 */
	private handleMessage(data: string): void {
		try {
			const message = JSON.parse(data) as ServerMessage;

			// Find pending request
			const pending = this.pendingRequests.get(message.id);
			if (!pending) {
				console.warn('Received message with unknown ID:', message.id);
				return;
			}

			// Clear timeout
			clearTimeout(pending.timeout);
			this.pendingRequests.delete(message.id);

			// Handle message based on type
			if (isErrorMessage(message)) {
				const errorPayload = message.payload as ErrorPayload;
				const error = new Error(errorPayload.message);
				Object.assign(error, errorPayload);
				pending.reject(error);
			} else if (isResultMessage(message)) {
				pending.resolve(message.payload);
			} else if (isSchemaMessage(message)) {
				pending.resolve(message.payload);
			} else if (isPongMessage(message)) {
				pending.resolve(message.payload);
			} else {
				pending.reject(new Error(`Unknown message type: ${message.type}`));
			}
		} catch (error) {
			console.error('Failed to handle message:', error);
			this.handleError(error instanceof Error ? error : new Error(String(error)));
		}
	}

	/**
	 * Handles WebSocket close events
	 */
	private handleClose(event: CloseEvent): void {
		console.log('WebSocket closed:', event.code, event.reason);

		// Clean up
		this.ws = null;

		// Reject all pending requests
		this.pendingRequests.forEach((request) => {
			clearTimeout(request.timeout);
			request.reject(new Error('Connection closed'));
		});
		this.pendingRequests.clear();

		// Attempt reconnection if not a clean close
		if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
			this.attemptReconnect();
		} else {
			this.setState(ConnectionState.DISCONNECTED);
		}
	}

	/**
	 * Attempts to reconnect to the server
	 */
	private attemptReconnect(): void {
		this.reconnectAttempts++;
		this.setState(ConnectionState.RECONNECTING);

		console.log(
			`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`
		);

		this.reconnectTimer = setTimeout(() => {
			this.connect().catch((error) => {
				console.error('Reconnection failed:', error);
				if (this.reconnectAttempts >= this.maxReconnectAttempts) {
					this.handleError(new Error('Max reconnection attempts reached'));
					this.setState(ConnectionState.ERROR);
				}
			});
		}, this.reconnectDelay * this.reconnectAttempts);
	}

	/**
	 * Sets the connection state and notifies listeners
	 */
	private setState(state: ConnectionState): void {
		if (this.state !== state) {
			this.state = state;
			this.onStateChangeCallback?.(state);
		}
	}

	/**
	 * Handles errors
	 */
	private handleError(error: Error): void {
		console.error('PostgresProxyClient error:', error);
		this.onErrorCallback?.(error);
	}
}

/**
 * Creates a new PostgresProxyClient instance
 * @param options Client configuration options
 * @returns New client instance
 */
export function createClient(options: ClientOptions): PostgresProxyClient {
	return new PostgresProxyClient(options);
}
