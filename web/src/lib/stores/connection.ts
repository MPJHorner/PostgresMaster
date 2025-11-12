/**
 * Connection Store
 * Manages the WebSocket connection state and client instance
 */

import { writable, derived, type Readable } from 'svelte/store';
import { PostgresProxyClient, ConnectionState } from '$lib/services/websocket';

/**
 * Connection state interface
 */
export interface ConnectionStateData {
	/** Current connection state */
	state: ConnectionState;
	/** Active client instance (null if not connected) */
	client: PostgresProxyClient | null;
	/** Last error message (null if no error) */
	error: string | null;
	/** Secret used for connection */
	secret: string | null;
}

/**
 * Initial connection state
 */
const initialState: ConnectionStateData = {
	state: ConnectionState.DISCONNECTED,
	client: null,
	error: null,
	secret: null
};

/**
 * Main connection store
 * Holds the current connection state and client instance
 */
export const connectionStore = writable<ConnectionStateData>(initialState);

/**
 * Derived store: indicates if currently connected
 */
export const isConnected: Readable<boolean> = derived(
	connectionStore,
	($connection) =>
		$connection.state === ConnectionState.CONNECTED && $connection.client !== null
);

/**
 * Derived store: provides the active client instance
 */
export const client: Readable<PostgresProxyClient | null> = derived(
	connectionStore,
	($connection) => $connection.client
);

/**
 * Derived store: indicates if currently connecting
 */
export const isConnecting: Readable<boolean> = derived(
	connectionStore,
	($connection) => $connection.state === ConnectionState.CONNECTING
);

/**
 * Derived store: indicates if currently reconnecting
 */
export const isReconnecting: Readable<boolean> = derived(
	connectionStore,
	($connection) => $connection.state === ConnectionState.RECONNECTING
);

/**
 * Derived store: indicates if in error state
 */
export const hasError: Readable<boolean> = derived(
	connectionStore,
	($connection) => $connection.state === ConnectionState.ERROR || $connection.error !== null
);

/**
 * Derived store: provides current error message
 */
export const errorMessage: Readable<string | null> = derived(
	connectionStore,
	($connection) => $connection.error
);

/**
 * Derived store: provides current connection state
 */
export const currentState: Readable<ConnectionState> = derived(
	connectionStore,
	($connection) => $connection.state
);

/**
 * Connects to the proxy server with the given secret
 * @param secret Authentication secret from URL parameter
 * @param url Optional WebSocket URL (defaults to ws://localhost:8080)
 * @returns Promise that resolves when connected
 */
export async function connect(secret: string, url?: string): Promise<void> {
	// Clear any previous error
	connectionStore.update((state) => ({
		...state,
		error: null,
		secret,
		state: ConnectionState.CONNECTING
	}));

	try {
		// Create new client instance
		const newClient = new PostgresProxyClient({
			secret,
			url: url || 'ws://localhost:8080',
			maxReconnectAttempts: 3,
			reconnectDelay: 2000,
			connectionTimeout: 10000
		});

		// Register state change callback
		newClient.onStateChange((newState) => {
			connectionStore.update((state) => ({
				...state,
				state: newState
			}));
		});

		// Register error callback
		newClient.onError((error) => {
			connectionStore.update((state) => ({
				...state,
				error: error.message,
				state: ConnectionState.ERROR
			}));
		});

		// Attempt connection
		await newClient.connect();

		// Update store with connected state
		connectionStore.update((state) => ({
			...state,
			client: newClient,
			state: ConnectionState.CONNECTED,
			error: null
		}));
	} catch (error) {
		// Handle connection error
		const errorMsg = error instanceof Error ? error.message : String(error);
		connectionStore.update((state) => ({
			...state,
			client: null,
			state: ConnectionState.ERROR,
			error: errorMsg
		}));
		throw error;
	}
}

/**
 * Disconnects from the proxy server
 */
export function disconnect(): void {
	connectionStore.update((state) => {
		// Close existing connection if present
		if (state.client) {
			state.client.close();
		}

		return {
			...initialState
		};
	});
}

/**
 * Clears the current error message
 */
export function clearError(): void {
	connectionStore.update((state) => ({
		...state,
		error: null
	}));
}

/**
 * Sets a custom error message
 * @param message Error message to display
 */
export function setError(message: string): void {
	connectionStore.update((state) => ({
		...state,
		error: message,
		state: ConnectionState.ERROR
	}));
}

/**
 * Resets the connection store to initial state
 */
export function resetConnection(): void {
	connectionStore.update((state) => {
		// Close existing connection if present
		if (state.client) {
			state.client.close();
		}

		return { ...initialState };
	});
}

/**
 * Gets the current connection state synchronously
 * Note: Prefer using the derived stores in components
 */
export function getConnectionState(): ConnectionStateData {
	let state: ConnectionStateData = initialState;
	connectionStore.subscribe((s) => {
		state = s;
	})();
	return state;
}
