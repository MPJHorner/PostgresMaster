/**
 * Connection Store Tests
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { get } from 'svelte/store';
import {
	connectionStore,
	isConnected,
	client,
	isConnecting,
	isReconnecting,
	hasError,
	errorMessage,
	currentState,
	disconnect,
	clearError,
	setError,
	resetConnection,
	getConnectionState
} from './connection';
import { ConnectionState } from '$lib/services/websocket';

describe('Connection Store', () => {
	beforeEach(() => {
		// Reset connection store before each test
		resetConnection();
	});

	describe('Initial State', () => {
		it('should have correct initial state', () => {
			const state = get(connectionStore);
			expect(state.state).toBe(ConnectionState.DISCONNECTED);
			expect(state.client).toBeNull();
			expect(state.error).toBeNull();
			expect(state.secret).toBeNull();
		});

		it('should not be connected initially', () => {
			expect(get(isConnected)).toBe(false);
		});

		it('should not be connecting initially', () => {
			expect(get(isConnecting)).toBe(false);
		});

		it('should not have error initially', () => {
			expect(get(hasError)).toBe(false);
			expect(get(errorMessage)).toBeNull();
		});

		it('should have client as null initially', () => {
			expect(get(client)).toBeNull();
		});

		it('should have DISCONNECTED state initially', () => {
			expect(get(currentState)).toBe(ConnectionState.DISCONNECTED);
		});
	});

	describe('Derived Stores', () => {
		it('isConnected should be true when state is CONNECTED and client exists', () => {
			connectionStore.set({
				state: ConnectionState.CONNECTED,
				client: {} as any, // Mock client
				error: null,
				secret: 'test-secret'
			});

			expect(get(isConnected)).toBe(true);
		});

		it('isConnected should be false when state is CONNECTED but client is null', () => {
			connectionStore.set({
				state: ConnectionState.CONNECTED,
				client: null,
				error: null,
				secret: 'test-secret'
			});

			expect(get(isConnected)).toBe(false);
		});

		it('isConnecting should be true when state is CONNECTING', () => {
			connectionStore.set({
				state: ConnectionState.CONNECTING,
				client: null,
				error: null,
				secret: 'test-secret'
			});

			expect(get(isConnecting)).toBe(true);
		});

		it('isReconnecting should be true when state is RECONNECTING', () => {
			connectionStore.set({
				state: ConnectionState.RECONNECTING,
				client: null,
				error: null,
				secret: 'test-secret'
			});

			expect(get(isReconnecting)).toBe(true);
		});

		it('hasError should be true when state is ERROR', () => {
			connectionStore.set({
				state: ConnectionState.ERROR,
				client: null,
				error: null,
				secret: null
			});

			expect(get(hasError)).toBe(true);
		});

		it('hasError should be true when error message is present', () => {
			connectionStore.set({
				state: ConnectionState.DISCONNECTED,
				client: null,
				error: 'Test error',
				secret: null
			});

			expect(get(hasError)).toBe(true);
		});

		it('errorMessage should return the error message', () => {
			const testError = 'Connection failed';
			connectionStore.set({
				state: ConnectionState.ERROR,
				client: null,
				error: testError,
				secret: null
			});

			expect(get(errorMessage)).toBe(testError);
		});

		it('currentState should return the current state', () => {
			connectionStore.set({
				state: ConnectionState.CONNECTING,
				client: null,
				error: null,
				secret: 'test-secret'
			});

			expect(get(currentState)).toBe(ConnectionState.CONNECTING);
		});

		it('client should return the client instance', () => {
			const mockClient = {} as any;
			connectionStore.set({
				state: ConnectionState.CONNECTED,
				client: mockClient,
				error: null,
				secret: 'test-secret'
			});

			expect(get(client)).toBe(mockClient);
		});
	});

	describe('Helper Functions', () => {
		it('setError should set error message and ERROR state', () => {
			const errorMsg = 'Test error message';
			setError(errorMsg);

			const state = get(connectionStore);
			expect(state.error).toBe(errorMsg);
			expect(state.state).toBe(ConnectionState.ERROR);
		});

		it('clearError should clear the error message', () => {
			setError('Test error');
			clearError();

			const state = get(connectionStore);
			expect(state.error).toBeNull();
		});

		it('disconnect should reset to initial state', () => {
			connectionStore.set({
				state: ConnectionState.CONNECTED,
				client: {} as any,
				error: null,
				secret: 'test-secret'
			});

			disconnect();

			const state = get(connectionStore);
			expect(state.state).toBe(ConnectionState.DISCONNECTED);
			expect(state.client).toBeNull();
			expect(state.error).toBeNull();
			expect(state.secret).toBeNull();
		});

		it('resetConnection should reset to initial state', () => {
			connectionStore.set({
				state: ConnectionState.CONNECTED,
				client: {} as any,
				error: 'Some error',
				secret: 'test-secret'
			});

			resetConnection();

			const state = get(connectionStore);
			expect(state.state).toBe(ConnectionState.DISCONNECTED);
			expect(state.client).toBeNull();
			expect(state.error).toBeNull();
			expect(state.secret).toBeNull();
		});

		it('getConnectionState should return current state synchronously', () => {
			const testState = {
				state: ConnectionState.CONNECTED,
				client: {} as any,
				error: null,
				secret: 'test-secret'
			};
			connectionStore.set(testState);

			const state = getConnectionState();
			expect(state.state).toBe(ConnectionState.CONNECTED);
			expect(state.secret).toBe('test-secret');
		});
	});

	describe('Store Reactivity', () => {
		it('should update derived stores when main store changes', () => {
			const values: boolean[] = [];
			const unsubscribe = isConnected.subscribe((value) => {
				values.push(value);
			});

			connectionStore.set({
				state: ConnectionState.CONNECTED,
				client: {} as any,
				error: null,
				secret: 'test-secret'
			});

			expect(values.length).toBeGreaterThan(1);
			expect(values[values.length - 1]).toBe(true);

			unsubscribe();
		});

		it('should update multiple derived stores simultaneously', () => {
			connectionStore.set({
				state: ConnectionState.CONNECTING,
				client: null,
				error: null,
				secret: 'test-secret'
			});

			expect(get(isConnecting)).toBe(true);
			expect(get(isConnected)).toBe(false);
			expect(get(currentState)).toBe(ConnectionState.CONNECTING);
		});
	});
});
