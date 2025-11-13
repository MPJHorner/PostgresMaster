/**
 * Formats a value for display in the results table.
 * Handles various data types and converts them to a string representation.
 *
 * @param value - The value to format (can be any type)
 * @returns A string representation of the value suitable for display
 *
 * @example
 * ```typescript
 * formatValue(null) // returns 'NULL'
 * formatValue(true) // returns 'true'
 * formatValue(123) // returns '123'
 * formatValue([1, 2, 3]) // returns '[1,2,3]'
 * ```
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
 * Simplifies type names for better readability by removing verbose qualifiers.
 *
 * @param type - The PostgreSQL column type (e.g., "character varying(255)")
 * @returns A simplified type string (e.g., "varchar")
 *
 * @example
 * ```typescript
 * formatColumnType('character varying(255)') // returns 'varchar'
 * formatColumnType('timestamp without time zone') // returns 'timestamp'
 * ```
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
 *
 * @param ms - Execution time in milliseconds
 * @returns A formatted string (e.g., "15ms" or "1.52s")
 *
 * @example
 * ```typescript
 * formatExecutionTime(0.5) // returns '< 1ms'
 * formatExecutionTime(150) // returns '150ms'
 * formatExecutionTime(1520) // returns '1.52s'
 * ```
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
 *
 * @param count - The number of rows
 * @returns A formatted string with proper pluralization (e.g., "No rows", "1 row", "100 rows")
 *
 * @example
 * ```typescript
 * formatRowCount(0) // returns 'No rows'
 * formatRowCount(1) // returns '1 row'
 * formatRowCount(1000) // returns '1,000 rows'
 * ```
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
