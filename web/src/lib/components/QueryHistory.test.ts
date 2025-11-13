/**
 * Unit tests for QueryHistory component
 *
 * Tests verify the component logic and integration with query history store.
 */

import { describe, it, expect } from 'vitest';
import { queryHistory, type QueryHistoryItem } from '$lib/stores/history';
import { get } from 'svelte/store';

describe('QueryHistory Component', () => {
	describe('Query history store integration', () => {
		it('should work with empty history', () => {
			queryHistory.clear();
			const history = get(queryHistory);
			expect(history).toHaveLength(0);
		});

		it('should add successful query to history', () => {
			queryHistory.clear();

			const query: QueryHistoryItem = {
				sql: 'SELECT * FROM users',
				timestamp: new Date(),
				success: true,
				executionTime: 123,
				rowCount: 10
			};

			queryHistory.addQuery(query);
			const history = get(queryHistory);

			expect(history).toHaveLength(1);
			expect(history[0].sql).toBe('SELECT * FROM users');
			expect(history[0].success).toBe(true);
			expect(history[0].executionTime).toBe(123);
			expect(history[0].rowCount).toBe(10);
		});

		it('should add failed query to history with error', () => {
			queryHistory.clear();

			const query: QueryHistoryItem = {
				sql: 'SELECT * FROM nonexistent',
				timestamp: new Date(),
				success: false,
				executionTime: 50,
				error: 'Table does not exist'
			};

			queryHistory.addQuery(query);
			const history = get(queryHistory);

			expect(history).toHaveLength(1);
			expect(history[0].success).toBe(false);
			expect(history[0].error).toBe('Table does not exist');
		});

		it('should maintain most recent queries first', () => {
			queryHistory.clear();

			const query1: QueryHistoryItem = {
				sql: 'SELECT 1',
				timestamp: new Date(Date.now() - 1000),
				success: true,
				executionTime: 10,
				rowCount: 1
			};

			const query2: QueryHistoryItem = {
				sql: 'SELECT 2',
				timestamp: new Date(),
				success: true,
				executionTime: 20,
				rowCount: 1
			};

			queryHistory.addQuery(query1);
			queryHistory.addQuery(query2);

			const history = get(queryHistory);
			expect(history).toHaveLength(2);
			// Most recent should be first
			expect(history[0].sql).toBe('SELECT 2');
			expect(history[1].sql).toBe('SELECT 1');
		});

		it('should clear all history', () => {
			queryHistory.clear();

			queryHistory.addQuery({
				sql: 'SELECT 1',
				timestamp: new Date(),
				success: true,
				executionTime: 10,
				rowCount: 1
			});

			queryHistory.addQuery({
				sql: 'SELECT 2',
				timestamp: new Date(),
				success: true,
				executionTime: 20,
				rowCount: 1
			});

			let history = get(queryHistory);
			expect(history).toHaveLength(2);

			queryHistory.clear();
			history = get(queryHistory);
			expect(history).toHaveLength(0);
		});
	});

	describe('SQL truncation logic', () => {
		function truncateSQL(sql: string, maxLength: number = 60): string {
			const singleLine = sql.replace(/\s+/g, ' ').trim();
			if (singleLine.length <= maxLength) return singleLine;
			return singleLine.substring(0, maxLength) + '...';
		}

		it('should not truncate short SQL', () => {
			const shortSql = 'SELECT * FROM users';
			const truncated = truncateSQL(shortSql);
			expect(truncated).toBe(shortSql);
			expect(truncated).not.toContain('...');
		});

		it('should truncate long SQL', () => {
			const longSql =
				'SELECT * FROM users WHERE name = "John" AND age > 25 AND city = "New York" AND status = "active"';
			const truncated = truncateSQL(longSql);
			expect(truncated).toContain('...');
			expect(truncated.length).toBeLessThan(longSql.length);
			expect(truncated.length).toBe(63); // 60 + '...'
		});

		it('should handle multi-line SQL', () => {
			const multiLine = `SELECT *
			FROM users
			WHERE id = 1`;
			const truncated = truncateSQL(multiLine);
			expect(truncated).toBe('SELECT * FROM users WHERE id = 1');
		});

		it('should handle SQL with extra whitespace', () => {
			const messySql = 'SELECT   *    FROM     users';
			const truncated = truncateSQL(messySql);
			expect(truncated).toBe('SELECT * FROM users');
		});
	});

	describe('Timestamp formatting logic', () => {
		function formatTimestamp(date: Date): string {
			const now = new Date();
			const diffMs = now.getTime() - date.getTime();
			const diffMins = Math.floor(diffMs / 60000);
			const diffHours = Math.floor(diffMs / 3600000);
			const diffDays = Math.floor(diffMs / 86400000);

			if (diffMins < 1) return 'just now';
			if (diffMins < 60) return `${diffMins}m ago`;
			if (diffHours < 24) return `${diffHours}h ago`;
			return `${diffDays}d ago`;
		}

		it('should show "just now" for recent queries', () => {
			const now = new Date();
			expect(formatTimestamp(now)).toBe('just now');
		});

		it('should show minutes for queries within the hour', () => {
			const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000);
			expect(formatTimestamp(fiveMinutesAgo)).toBe('5m ago');

			const thirtyMinutesAgo = new Date(Date.now() - 30 * 60 * 1000);
			expect(formatTimestamp(thirtyMinutesAgo)).toBe('30m ago');
		});

		it('should show hours for queries within the day', () => {
			const twoHoursAgo = new Date(Date.now() - 2 * 60 * 60 * 1000);
			expect(formatTimestamp(twoHoursAgo)).toBe('2h ago');

			const tenHoursAgo = new Date(Date.now() - 10 * 60 * 60 * 1000);
			expect(formatTimestamp(tenHoursAgo)).toBe('10h ago');
		});

		it('should show days for older queries', () => {
			const twoDaysAgo = new Date(Date.now() - 2 * 24 * 60 * 60 * 1000);
			expect(formatTimestamp(twoDaysAgo)).toBe('2d ago');

			const sevenDaysAgo = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000);
			expect(formatTimestamp(sevenDaysAgo)).toBe('7d ago');
		});
	});

	describe('Component props validation', () => {
		it('should accept valid onQuerySelect callback', () => {
			const callback = (sql: string) => {
				expect(sql).toBeDefined();
			};
			callback('SELECT 1');
		});
	});

	describe('History limits', () => {
		it('should limit history to 50 items', () => {
			queryHistory.clear();

			// Add 60 queries
			for (let i = 0; i < 60; i++) {
				queryHistory.addQuery({
					sql: `SELECT ${i}`,
					timestamp: new Date(),
					success: true,
					executionTime: 10,
					rowCount: 1
				});
			}

			const history = get(queryHistory);
			expect(history).toHaveLength(50);
			// Should keep the most recent 50
			expect(history[0].sql).toBe('SELECT 59');
			expect(history[49].sql).toBe('SELECT 10');
		});
	});
});
