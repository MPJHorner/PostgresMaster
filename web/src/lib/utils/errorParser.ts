/**
 * Postgres error structure
 */
export interface PostgresError {
	code?: string;
	message: string;
	detail?: string;
	hint?: string;
	position?: number;
}

/**
 * Parse error string to extract Postgres error details.
 * Attempts to parse JSON error or plain text error format.
 *
 * @param errorStr - The error string from the database or proxy
 * @returns A structured PostgresError object with extracted details
 *
 * @example
 * ```typescript
 * const error = parseError('ERROR: relation "users" does not exist (SQLSTATE 42P01)');
 * // returns: { code: '42P01', message: 'relation "users" does not exist', ... }
 * ```
 */
export function parseError(errorStr: string): PostgresError {
	// Try to parse as JSON first (if proxy sends structured error)
	try {
		const parsed = JSON.parse(errorStr);
		if (parsed.code || parsed.message) {
			return {
				code: parsed.code,
				message: parsed.message || errorStr,
				detail: parsed.detail,
				hint: parsed.hint,
				position: parsed.position
			};
		}
	} catch {
		// Not JSON, continue to text parsing
	}

	// Try to parse common Postgres error format
	// Example: "ERROR: relation "users" does not exist (SQLSTATE 42P01)"
	const codeMatch = errorStr.match(/\(SQLSTATE\s+([A-Z0-9]{5})\)/i);
	const code = codeMatch ? codeMatch[1] : undefined;

	// Extract detail if present
	const detailMatch = errorStr.match(/DETAIL:\s*(.+?)(?=\n|$)/i);
	const detail = detailMatch ? detailMatch[1].trim() : undefined;

	// Extract hint if present
	const hintMatch = errorStr.match(/HINT:\s*(.+?)(?=\n|$)/i);
	const hint = hintMatch ? hintMatch[1].trim() : undefined;

	// Extract position if present
	const positionMatch = errorStr.match(/at character (\d+)/i);
	const position = positionMatch ? parseInt(positionMatch[1], 10) : undefined;

	// Clean up the main message (remove SQLSTATE, DETAIL, HINT, position indicator, etc.)
	let message = errorStr
		.replace(/\(SQLSTATE\s+[A-Z0-9]{5}\)/gi, '')
		.replace(/DETAIL:\s*.+/gi, '')
		.replace(/HINT:\s*.+/gi, '')
		.replace(/\s+at character \d+/gi, '')
		.trim();

	// If message starts with "ERROR:", remove it for cleaner display
	message = message.replace(/^ERROR:\s*/i, '');

	return {
		code,
		message,
		detail,
		hint,
		position
	};
}

/**
 * Get user-friendly error code description for PostgreSQL error codes.
 *
 * @param code - PostgreSQL error code (e.g., "42P01", "23505")
 * @returns A human-readable description of the error, or undefined if code is unknown
 *
 * @example
 * ```typescript
 * getErrorCodeDescription('42P01') // returns 'Undefined table'
 * getErrorCodeDescription('23505') // returns 'Unique violation'
 * getErrorCodeDescription('99999') // returns undefined
 * ```
 */
export function getErrorCodeDescription(code?: string): string | undefined {
	if (!code) return undefined;

	const descriptions: Record<string, string> = {
		'42P01': 'Undefined table',
		'42703': 'Undefined column',
		'42601': 'Syntax error',
		'42501': 'Insufficient privilege',
		'23505': 'Unique violation',
		'23503': 'Foreign key violation',
		'23502': 'Not null violation',
		'22P02': 'Invalid text representation',
		'08006': 'Connection failure',
		'08003': 'Connection does not exist',
		'57P03': 'Cannot connect now'
	};

	return descriptions[code];
}
