/**
 * Formats a value for display in the results table.
 * Handles various data types and converts them to a string representation.
 */
export function formatValue(value: unknown): string {
	// Handle null/undefined
	if (value === null || value === undefined) {
		return 'NULL';
	}

	// Handle boolean
	if (typeof value === 'boolean') {
		return value.toString();
	}

	// Handle number
	if (typeof value === 'number') {
		return value.toString();
	}

	// Handle string
	if (typeof value === 'string') {
		return value;
	}

	// Handle Date
	if (value instanceof Date) {
		return value.toISOString();
	}

	// Handle Array
	if (Array.isArray(value)) {
		return JSON.stringify(value);
	}

	// Handle Object (including plain objects and other types)
	if (typeof value === 'object') {
		return JSON.stringify(value, null, 2);
	}

	// Fallback for any other type
	return String(value);
}

/**
 * Formats a PostgreSQL column type for display.
 * Simplifies type names for better readability.
 */
export function formatColumnType(type: string): string {
	// Remove common type modifiers for display
	return type
		.replace(/character varying/i, 'varchar')
		.replace(/timestamp without time zone/i, 'timestamp')
		.replace(/timestamp with time zone/i, 'timestamptz')
		.replace(/double precision/i, 'double')
		.replace(/\([^)]+\)/g, ''); // Remove length specifiers like (255) or (10,2)
}

/**
 * Formats execution time in milliseconds to a human-readable string.
 */
export function formatExecutionTime(ms: number): string {
	if (ms < 1) {
		return '< 1ms';
	}
	if (ms < 1000) {
		return `${Math.round(ms)}ms`;
	}
	const seconds = (ms / 1000).toFixed(2);
	return `${seconds}s`;
}

/**
 * Formats row count with proper pluralization.
 */
export function formatRowCount(count: number): string {
	if (count === 0) {
		return 'No rows';
	}
	if (count === 1) {
		return '1 row';
	}
	return `${count.toLocaleString()} rows`;
}
