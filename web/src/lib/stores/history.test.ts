import { describe, it, expect, beforeEach } from 'vitest';
import { queryHistory, type QueryHistoryItem } from './history';
import { get } from 'svelte/store';

describe('queryHistory store', () => {
	beforeEach(() => {
		// Clear history before each test
		queryHistory.clear();
	});

	it('should start with empty history', () => {
		const history = get(queryHistory);
		expect(history).toEqual([]);
	});

	it('should add a query to history', () => {
		const query: QueryHistoryItem = {
			sql: 'SELECT 1',
			timestamp: new Date(),
			success: true,
			executionTime: 10,
			rowCount: 1
		};

		queryHistory.addQuery(query);

		const history = get(queryHistory);
		expect(history).toHaveLength(1);
		expect(history[0]).toEqual(query);
	});

	it('should add failed query to history', () => {
		const query: QueryHistoryItem = {
			sql: 'SELECT * FROM nonexistent',
			timestamp: new Date(),
			success: false,
			executionTime: 5,
			error: 'Table not found'
		};

		queryHistory.addQuery(query);

		const history = get(queryHistory);
		expect(history).toHaveLength(1);
		expect(history[0]).toEqual(query);
		expect(history[0].success).toBe(false);
		expect(history[0].error).toBe('Table not found');
	});

	it('should add queries in reverse chronological order (newest first)', () => {
		const query1: QueryHistoryItem = {
			sql: 'SELECT 1',
			timestamp: new Date('2025-01-01T10:00:00'),
			success: true
		};

		const query2: QueryHistoryItem = {
			sql: 'SELECT 2',
			timestamp: new Date('2025-01-01T10:00:01'),
			success: true
		};

		queryHistory.addQuery(query1);
		queryHistory.addQuery(query2);

		const history = get(queryHistory);
		expect(history).toHaveLength(2);
		expect(history[0]).toEqual(query2); // Newest first
		expect(history[1]).toEqual(query1);
	});

	it('should limit history to 50 queries', () => {
		// Add 60 queries
		for (let i = 0; i < 60; i++) {
			queryHistory.addQuery({
				sql: `SELECT ${i}`,
				timestamp: new Date(),
				success: true
			});
		}

		const history = get(queryHistory);
		expect(history).toHaveLength(50); // Should be capped at 50
		expect(history[0].sql).toBe('SELECT 59'); // Most recent
		expect(history[49].sql).toBe('SELECT 10'); // Oldest kept
	});

	it('should clear history', () => {
		// Add some queries
		queryHistory.addQuery({
			sql: 'SELECT 1',
			timestamp: new Date(),
			success: true
		});

		queryHistory.addQuery({
			sql: 'SELECT 2',
			timestamp: new Date(),
			success: true
		});

		let history = get(queryHistory);
		expect(history).toHaveLength(2);

		// Clear history
		queryHistory.clear();

		history = get(queryHistory);
		expect(history).toHaveLength(0);
	});

	it('should handle queries with all optional fields', () => {
		const query: QueryHistoryItem = {
			sql: 'SELECT * FROM users',
			timestamp: new Date(),
			success: true,
			executionTime: 42,
			rowCount: 100
		};

		queryHistory.addQuery(query);

		const history = get(queryHistory);
		expect(history[0].executionTime).toBe(42);
		expect(history[0].rowCount).toBe(100);
	});

	it('should handle queries without optional fields', () => {
		const query: QueryHistoryItem = {
			sql: 'SELECT * FROM users',
			timestamp: new Date(),
			success: true
		};

		queryHistory.addQuery(query);

		const history = get(queryHistory);
		expect(history[0].executionTime).toBeUndefined();
		expect(history[0].rowCount).toBeUndefined();
		expect(history[0].error).toBeUndefined();
	});
});
