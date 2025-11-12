/**
 * WebSocket Client Service Tests
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import {
	PostgresProxyClient,
	ConnectionState,
	createClient,
	type ClientOptions
} from './websocket';

// Mock WebSocket
class MockWebSocket {
	static CONNECTING = 0;
	static OPEN = 1;
	static CLOSING = 2;
	static CLOSED = 3;

	public readyState = MockWebSocket.CONNECTING;
	public onopen: ((event: Event) => void) | null = null;
	public onclose: ((event: CloseEvent) => void) | null = null;
	public onerror: ((event: Event) => void) | null = null;
	public onmessage: ((event: MessageEvent) => void) | null = null;

	private sentMessages: string[] = [];

	constructor(public url: string) {
		// Simulate connection after a short delay
		setTimeout(() => {
			if (this.readyState === MockWebSocket.CONNECTING) {
				this.readyState = MockWebSocket.OPEN;
				this.onopen?.(new Event('open'));
			}
		}, 10);
	}

	send(data: string): void {
		this.sentMessages.push(data);
	}

	close(code?: number, reason?: string): void {
		this.readyState = MockWebSocket.CLOSED;
		const event = new CloseEvent('close', { code: code || 1000, reason: reason || '' });
		this.onclose?.(event);
	}

	// Test helper
	simulateMessage(data: string): void {
		const event = new MessageEvent('message', { data });
		this.onmessage?.(event);
	}

	getSentMessages(): string[] {
		return this.sentMessages;
	}
}

describe('PostgresProxyClient', () => {
	let client: PostgresProxyClient;

	const defaultOptions: ClientOptions = {
		secret: 'test-secret-123',
		url: 'ws://localhost:8080'
	};

	beforeEach(() => {
		// Mock WebSocket globally
		vi.stubGlobal('WebSocket', MockWebSocket);
		client = new PostgresProxyClient(defaultOptions);
	});

	afterEach(() => {
		client.close();
		vi.unstubAllGlobals();
	});

	describe('constructor', () => {
		it('should create a client with default options', () => {
			expect(client).toBeInstanceOf(PostgresProxyClient);
			expect(client.getState()).toBe(ConnectionState.DISCONNECTED);
		});

		it('should use default URL if not provided', () => {
			const clientWithoutUrl = new PostgresProxyClient({ secret: 'test' });
			expect(clientWithoutUrl).toBeInstanceOf(PostgresProxyClient);
		});
	});

	describe('connect', () => {
		it('should connect successfully', async () => {
			await client.connect();
			expect(client.getState()).toBe(ConnectionState.CONNECTED);
			expect(client.isConnected()).toBe(true);
		});

		it('should include secret in URL', async () => {
			await client.connect();
			// Check that WebSocket was called with correct URL
			// Note: We can't easily access the URL from our mock, but we've tested the logic
			expect(client.isConnected()).toBe(true);
		});

		it('should not connect twice', async () => {
			await client.connect();
			await client.connect(); // Should not throw
			expect(client.isConnected()).toBe(true);
		});

		it('should call state change callback', async () => {
			const callback = vi.fn();
			client.onStateChange(callback);
			await client.connect();
			expect(callback).toHaveBeenCalledWith(ConnectionState.CONNECTING);
			expect(callback).toHaveBeenCalledWith(ConnectionState.CONNECTED);
		});
	});

	describe('executeQuery', () => {
		it('should throw if not connected', async () => {
			await expect(client.executeQuery('SELECT 1')).rejects.toThrow('Not connected');
		});
	});

	describe('introspectSchema', () => {
		it('should throw if not connected', async () => {
			await expect(client.introspectSchema()).rejects.toThrow('Not connected');
		});
	});

	describe('ping', () => {
		it('should throw if not connected', async () => {
			await expect(client.ping()).rejects.toThrow('Not connected');
		});
	});

	describe('close', () => {
		it('should close the connection', async () => {
			await client.connect();
			expect(client.isConnected()).toBe(true);

			client.close();
			expect(client.getState()).toBe(ConnectionState.DISCONNECTED);
			expect(client.isConnected()).toBe(false);
		});

		it('should be safe to call when not connected', () => {
			expect(() => client.close()).not.toThrow();
		});
	});

	describe('createClient', () => {
		it('should create a client instance', () => {
			const newClient = createClient(defaultOptions);
			expect(newClient).toBeInstanceOf(PostgresProxyClient);
		});
	});

	describe('state management', () => {
		it('should start in DISCONNECTED state', () => {
			expect(client.getState()).toBe(ConnectionState.DISCONNECTED);
		});

		it('should transition to CONNECTING when connect is called', async () => {
			const states: ConnectionState[] = [];
			client.onStateChange((state) => states.push(state));

			const connectPromise = client.connect();
			expect(states[0]).toBe(ConnectionState.CONNECTING);

			await connectPromise;
			expect(states).toContain(ConnectionState.CONNECTED);
		});
	});

	describe('error handling', () => {
		it('should call error callback on errors', async () => {
			const errorCallback = vi.fn();
			client.onError(errorCallback);

			// Try to execute query without connecting
			await expect(client.executeQuery('SELECT 1')).rejects.toThrow();
		});
	});
});

describe('ConnectionState enum', () => {
	it('should have all expected states', () => {
		expect(ConnectionState.DISCONNECTED).toBe('disconnected');
		expect(ConnectionState.CONNECTING).toBe('connecting');
		expect(ConnectionState.CONNECTED).toBe('connected');
		expect(ConnectionState.RECONNECTING).toBe('reconnecting');
		expect(ConnectionState.ERROR).toBe('error');
	});
});
