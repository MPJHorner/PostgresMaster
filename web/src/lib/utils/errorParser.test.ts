import { describe, it, expect } from 'vitest';
import { parseError, getErrorCodeDescription } from './errorParser';

describe('parseError', () => {
	it('parses simple error message', () => {
		const result = parseError('Something went wrong');
		expect(result.message).toBe('Something went wrong');
		expect(result.code).toBeUndefined();
		expect(result.detail).toBeUndefined();
		expect(result.hint).toBeUndefined();
		expect(result.position).toBeUndefined();
	});

	it('parses Postgres error with SQLSTATE code', () => {
		const result = parseError('ERROR: relation "users" does not exist (SQLSTATE 42P01)');
		expect(result.message).toBe('relation "users" does not exist');
		expect(result.code).toBe('42P01');
	});

	it('parses error with position', () => {
		const result = parseError('ERROR: syntax error at or near "SELCT" at character 1');
		expect(result.message).toContain('syntax error');
		expect(result.position).toBe(1);
	});

	it('parses error with DETAIL', () => {
		const result = parseError(
			'ERROR: duplicate key value violates unique constraint "users_pkey"\nDETAIL: Key (id)=(1) already exists.'
		);
		expect(result.message).toContain('duplicate key value');
		expect(result.detail).toBe('Key (id)=(1) already exists.');
	});

	it('parses error with HINT', () => {
		const result = parseError(
			'ERROR: column "nam" does not exist\nHINT: Perhaps you meant to reference the column "name".'
		);
		expect(result.message).toBe('column "nam" does not exist');
		expect(result.hint).toBe('Perhaps you meant to reference the column "name".');
	});

	it('parses error with multiple components', () => {
		const result = parseError(
			'ERROR: column "test" does not exist (SQLSTATE 42703) at character 8\nDETAIL: The column was not found\nHINT: Check the spelling'
		);
		expect(result.message).toBe('column "test" does not exist');
		expect(result.code).toBe('42703');
		expect(result.position).toBe(8);
		expect(result.detail).toBe('The column was not found');
		expect(result.hint).toBe('Check the spelling');
	});

	it('parses JSON error format', () => {
		const jsonError = JSON.stringify({
			code: '42703',
			message: 'Column does not exist',
			detail: 'The column "test" was not found',
			hint: 'Check your column names',
			position: 5
		});

		const result = parseError(jsonError);
		expect(result.message).toBe('Column does not exist');
		expect(result.code).toBe('42703');
		expect(result.detail).toBe('The column "test" was not found');
		expect(result.hint).toBe('Check your column names');
		expect(result.position).toBe(5);
	});

	it('handles JSON with only message', () => {
		const jsonError = JSON.stringify({
			message: 'Connection timeout'
		});

		const result = parseError(jsonError);
		expect(result.message).toBe('Connection timeout');
	});

	it('removes ERROR: prefix', () => {
		const result = parseError('ERROR: connection failed');
		expect(result.message).toBe('connection failed');
	});

	it('removes SQLSTATE from message', () => {
		const result = parseError('ERROR: test error (SQLSTATE 42601)');
		expect(result.message).toBe('test error');
		expect(result.message).not.toContain('SQLSTATE');
	});
});

describe('getErrorCodeDescription', () => {
	it('returns undefined for undefined code', () => {
		expect(getErrorCodeDescription(undefined)).toBeUndefined();
	});

	it('returns undefined for unknown code', () => {
		expect(getErrorCodeDescription('99999')).toBeUndefined();
	});

	it('returns description for known codes', () => {
		expect(getErrorCodeDescription('42P01')).toBe('Undefined table');
		expect(getErrorCodeDescription('42703')).toBe('Undefined column');
		expect(getErrorCodeDescription('42601')).toBe('Syntax error');
		expect(getErrorCodeDescription('42501')).toBe('Insufficient privilege');
		expect(getErrorCodeDescription('23505')).toBe('Unique violation');
		expect(getErrorCodeDescription('23503')).toBe('Foreign key violation');
		expect(getErrorCodeDescription('23502')).toBe('Not null violation');
		expect(getErrorCodeDescription('22P02')).toBe('Invalid text representation');
		expect(getErrorCodeDescription('08006')).toBe('Connection failure');
		expect(getErrorCodeDescription('08003')).toBe('Connection does not exist');
		expect(getErrorCodeDescription('57P03')).toBe('Cannot connect now');
	});
});
