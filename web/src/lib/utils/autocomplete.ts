/**
 * SQL Keywords and Functions for PostgreSQL Autocomplete
 *
 * This module provides comprehensive lists of SQL keywords, data types, and functions
 * specifically for PostgreSQL, used in Monaco Editor autocomplete functionality.
 */

/**
 * Core SQL keywords and PostgreSQL-specific keywords
 */
export const SQL_KEYWORDS = [
	// Core SQL commands
	'SELECT',
	'FROM',
	'WHERE',
	'INSERT',
	'INTO',
	'UPDATE',
	'DELETE',
	'CREATE',
	'ALTER',
	'DROP',
	'TRUNCATE',

	// Table and schema operations
	'TABLE',
	'INDEX',
	'VIEW',
	'MATERIALIZED',
	'SCHEMA',
	'DATABASE',
	'SEQUENCE',
	'TRIGGER',
	'FUNCTION',
	'PROCEDURE',

	// Query clauses
	'JOIN',
	'INNER',
	'LEFT',
	'RIGHT',
	'FULL',
	'OUTER',
	'CROSS',
	'NATURAL',
	'ON',
	'USING',
	'GROUP BY',
	'ORDER BY',
	'HAVING',
	'LIMIT',
	'OFFSET',
	'FETCH',
	'UNION',
	'INTERSECT',
	'EXCEPT',
	'ALL',
	'DISTINCT',
	'DISTINCT ON',

	// Filtering and conditions
	'AND',
	'OR',
	'NOT',
	'IN',
	'EXISTS',
	'BETWEEN',
	'LIKE',
	'ILIKE',
	'SIMILAR TO',
	'IS',
	'NULL',
	'TRUE',
	'FALSE',
	'CASE',
	'WHEN',
	'THEN',
	'ELSE',
	'END',

	// Constraints
	'PRIMARY KEY',
	'FOREIGN KEY',
	'REFERENCES',
	'UNIQUE',
	'CHECK',
	'DEFAULT',
	'NOT NULL',
	'CONSTRAINT',

	// Indexes
	'CONCURRENTLY',
	'BTREE',
	'HASH',
	'GIST',
	'GIN',
	'BRIN',

	// Transaction control
	'BEGIN',
	'COMMIT',
	'ROLLBACK',
	'SAVEPOINT',
	'RELEASE',
	'START TRANSACTION',

	// PostgreSQL-specific
	'RETURNING',
	'ON CONFLICT',
	'DO NOTHING',
	'DO UPDATE',
	'SET',
	'EXCLUDED',
	'LATERAL',
	'WITH',
	'RECURSIVE',
	'AS',
	'CTE',
	'WINDOW',
	'PARTITION BY',
	'OVER',
	'ROWS',
	'RANGE',
	'UNBOUNDED',
	'PRECEDING',
	'FOLLOWING',
	'CURRENT ROW',

	// Access control
	'GRANT',
	'REVOKE',
	'ROLE',
	'USER',
	'PRIVILEGES',
	'TO',
	'FOR',

	// Data modification
	'VALUES',
	'ADD',
	'COLUMN',
	'RENAME',
	'CASCADE',
	'RESTRICT',

	// Query modifiers
	'EXPLAIN',
	'ANALYZE',
	'VERBOSE',
	'COSTS',
	'BUFFERS',

	// Advanced features
	'COPY',
	'VACUUM',
	'CLUSTER',
	'REINDEX',
	'LISTEN',
	'NOTIFY',
	'UNLISTEN'
];

/**
 * PostgreSQL data types
 */
export const DATA_TYPES = [
	// Numeric types
	'INTEGER',
	'INT',
	'BIGINT',
	'SMALLINT',
	'DECIMAL',
	'NUMERIC',
	'REAL',
	'DOUBLE PRECISION',
	'SERIAL',
	'BIGSERIAL',
	'SMALLSERIAL',
	'MONEY',

	// Character types
	'VARCHAR',
	'CHARACTER VARYING',
	'CHAR',
	'CHARACTER',
	'TEXT',

	// Binary types
	'BYTEA',

	// Date/time types
	'DATE',
	'TIME',
	'TIMESTAMP',
	'TIMESTAMPTZ',
	'TIMESTAMP WITH TIME ZONE',
	'TIMESTAMP WITHOUT TIME ZONE',
	'TIME WITH TIME ZONE',
	'TIME WITHOUT TIME ZONE',
	'INTERVAL',

	// Boolean
	'BOOLEAN',
	'BOOL',

	// Geometric types
	'POINT',
	'LINE',
	'LSEG',
	'BOX',
	'PATH',
	'POLYGON',
	'CIRCLE',

	// Network types
	'INET',
	'CIDR',
	'MACADDR',
	'MACADDR8',

	// Bit string types
	'BIT',
	'BIT VARYING',
	'VARBIT',

	// Text search types
	'TSVECTOR',
	'TSQUERY',

	// UUID
	'UUID',

	// JSON types
	'JSON',
	'JSONB',

	// Array
	'ARRAY',

	// Range types
	'INT4RANGE',
	'INT8RANGE',
	'NUMRANGE',
	'TSRANGE',
	'TSTZRANGE',
	'DATERANGE',

	// XML
	'XML',

	// Other
	'ENUM',
	'COMPOSITE'
];

/**
 * SQL aggregate functions
 */
export const AGGREGATE_FUNCTIONS = [
	'COUNT',
	'SUM',
	'AVG',
	'MIN',
	'MAX',
	'ARRAY_AGG',
	'JSON_AGG',
	'JSONB_AGG',
	'JSON_OBJECT_AGG',
	'JSONB_OBJECT_AGG',
	'STRING_AGG',
	'XMLAGG',
	'BIT_AND',
	'BIT_OR',
	'BOOL_AND',
	'BOOL_OR',
	'EVERY'
];

/**
 * String functions
 */
export const STRING_FUNCTIONS = [
	'CONCAT',
	'CONCAT_WS',
	'LENGTH',
	'CHAR_LENGTH',
	'CHARACTER_LENGTH',
	'LOWER',
	'UPPER',
	'INITCAP',
	'TRIM',
	'LTRIM',
	'RTRIM',
	'BTRIM',
	'LPAD',
	'RPAD',
	'SUBSTRING',
	'SUBSTR',
	'REPLACE',
	'TRANSLATE',
	'REVERSE',
	'REPEAT',
	'SPLIT_PART',
	'REGEXP_REPLACE',
	'REGEXP_MATCH',
	'REGEXP_MATCHES',
	'REGEXP_SPLIT_TO_ARRAY',
	'REGEXP_SPLIT_TO_TABLE',
	'FORMAT',
	'QUOTE_IDENT',
	'QUOTE_LITERAL',
	'QUOTE_NULLABLE',
	'MD5',
	'ENCODE',
	'DECODE'
];

/**
 * Date and time functions
 */
export const DATE_TIME_FUNCTIONS = [
	'NOW',
	'CURRENT_DATE',
	'CURRENT_TIME',
	'CURRENT_TIMESTAMP',
	'LOCALTIME',
	'LOCALTIMESTAMP',
	'DATE_TRUNC',
	'DATE_PART',
	'EXTRACT',
	'AGE',
	'TO_CHAR',
	'TO_DATE',
	'TO_TIMESTAMP',
	'TO_NUMBER',
	'MAKE_DATE',
	'MAKE_TIME',
	'MAKE_TIMESTAMP',
	'MAKE_TIMESTAMPTZ',
	'MAKE_INTERVAL',
	'JUSTIFY_DAYS',
	'JUSTIFY_HOURS',
	'JUSTIFY_INTERVAL',
	'CLOCK_TIMESTAMP',
	'STATEMENT_TIMESTAMP',
	'TRANSACTION_TIMESTAMP',
	'TIMEOFDAY'
];

/**
 * JSON functions
 */
export const JSON_FUNCTIONS = [
	'JSON_BUILD_ARRAY',
	'JSON_BUILD_OBJECT',
	'JSON_OBJECT',
	'JSON_ARRAY',
	'JSONB_BUILD_ARRAY',
	'JSONB_BUILD_OBJECT',
	'JSONB_OBJECT',
	'JSON_EXTRACT_PATH',
	'JSON_EXTRACT_PATH_TEXT',
	'JSONB_EXTRACT_PATH',
	'JSONB_EXTRACT_PATH_TEXT',
	'JSON_ARRAY_LENGTH',
	'JSONB_ARRAY_LENGTH',
	'JSON_EACH',
	'JSON_EACH_TEXT',
	'JSONB_EACH',
	'JSONB_EACH_TEXT',
	'JSON_OBJECT_KEYS',
	'JSONB_OBJECT_KEYS',
	'JSON_POPULATE_RECORD',
	'JSONB_POPULATE_RECORD',
	'JSON_POPULATE_RECORDSET',
	'JSONB_POPULATE_RECORDSET',
	'JSON_TO_RECORD',
	'JSONB_TO_RECORD',
	'JSON_TO_RECORDSET',
	'JSONB_TO_RECORDSET',
	'JSON_STRIP_NULLS',
	'JSONB_STRIP_NULLS',
	'JSONB_SET',
	'JSONB_INSERT',
	'JSONB_PATH_EXISTS',
	'JSONB_PATH_MATCH',
	'JSONB_PATH_QUERY',
	'JSONB_PATH_QUERY_ARRAY',
	'JSONB_PATH_QUERY_FIRST',
	'JSONB_PRETTY'
];

/**
 * Mathematical functions
 */
export const MATH_FUNCTIONS = [
	'ABS',
	'CEIL',
	'CEILING',
	'FLOOR',
	'ROUND',
	'TRUNC',
	'SIGN',
	'MOD',
	'DIV',
	'EXP',
	'LN',
	'LOG',
	'LOG10',
	'POWER',
	'SQRT',
	'CBRT',
	'PI',
	'RANDOM',
	'SETSEED',
	'DEGREES',
	'RADIANS',
	'SIN',
	'COS',
	'TAN',
	'ASIN',
	'ACOS',
	'ATAN',
	'ATAN2',
	'SINH',
	'COSH',
	'TANH',
	'ASINH',
	'ACOSH',
	'ATANH'
];

/**
 * Window functions
 */
export const WINDOW_FUNCTIONS = [
	'ROW_NUMBER',
	'RANK',
	'DENSE_RANK',
	'PERCENT_RANK',
	'CUME_DIST',
	'NTILE',
	'LAG',
	'LEAD',
	'FIRST_VALUE',
	'LAST_VALUE',
	'NTH_VALUE'
];

/**
 * Conditional and comparison functions
 */
export const CONDITIONAL_FUNCTIONS = ['COALESCE', 'NULLIF', 'GREATEST', 'LEAST', 'CAST', 'CONVERT'];

/**
 * System information functions
 */
export const SYSTEM_FUNCTIONS = [
	'VERSION',
	'CURRENT_DATABASE',
	'CURRENT_SCHEMA',
	'CURRENT_USER',
	'SESSION_USER',
	'INET_CLIENT_ADDR',
	'INET_CLIENT_PORT',
	'INET_SERVER_ADDR',
	'INET_SERVER_PORT',
	'PG_BACKEND_PID',
	'PG_POSTMASTER_START_TIME',
	'PG_CONF_LOAD_TIME'
];

/**
 * Array functions
 */
export const ARRAY_FUNCTIONS = [
	'ARRAY_APPEND',
	'ARRAY_CAT',
	'ARRAY_DIMS',
	'ARRAY_FILL',
	'ARRAY_LENGTH',
	'ARRAY_LOWER',
	'ARRAY_NDIMS',
	'ARRAY_POSITION',
	'ARRAY_POSITIONS',
	'ARRAY_PREPEND',
	'ARRAY_REMOVE',
	'ARRAY_REPLACE',
	'ARRAY_TO_STRING',
	'ARRAY_UPPER',
	'CARDINALITY',
	'STRING_TO_ARRAY',
	'UNNEST'
];

/**
 * All functions combined
 */
export const ALL_FUNCTIONS = [
	...AGGREGATE_FUNCTIONS,
	...STRING_FUNCTIONS,
	...DATE_TIME_FUNCTIONS,
	...JSON_FUNCTIONS,
	...MATH_FUNCTIONS,
	...WINDOW_FUNCTIONS,
	...CONDITIONAL_FUNCTIONS,
	...SYSTEM_FUNCTIONS,
	...ARRAY_FUNCTIONS
];

/**
 * All keywords and data types combined
 */
export const ALL_KEYWORDS = [...SQL_KEYWORDS, ...DATA_TYPES];

/**
 * Complete list of all SQL completions (keywords, types, and functions)
 */
export const ALL_COMPLETIONS = [...ALL_KEYWORDS, ...ALL_FUNCTIONS];

/**
 * Schema information for autocomplete
 */
export interface SchemaInfo {
	tables: TableInfo[];
	functions?: FunctionInfo[];
}

export interface TableInfo {
	schema: string;
	name: string;
	columns: ColumnInfo[];
}

export interface ColumnInfo {
	name: string;
	type: string;
	nullable: boolean;
}

export interface FunctionInfo {
	schema: string;
	name: string;
	returnType: string;
}

/**
 * SQL Context information for context-aware autocomplete
 */
interface SQLContext {
	inFromClause: boolean;
	inWhereClause: boolean;
	inSelectClause: boolean;
	afterTableDot: string | null; // Table name if we're after "tablename."
	availableTables: string[]; // Tables mentioned in FROM/JOIN clauses
}

/**
 * Parse SQL context at cursor position
 *
 * Analyzes the SQL text before the cursor to determine:
 * - Which clause we're in (FROM, WHERE, SELECT)
 * - If we're after a table name with a dot (for column completion)
 * - Which tables are available in the current query
 *
 * @param model - Monaco editor model
 * @param position - Current cursor position
 * @returns Context information for smart suggestions
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function parseSQLContext(model: any, position: any): SQLContext {
	const context: SQLContext = {
		inFromClause: false,
		inWhereClause: false,
		inSelectClause: false,
		afterTableDot: null,
		availableTables: []
	};

	// Get all text before cursor
	const textBeforeCursor = model.getValueInRange({
		startLineNumber: 1,
		startColumn: 1,
		endLineNumber: position.lineNumber,
		endColumn: position.column
	});

	// Convert to uppercase for easier matching (preserve original for table names)
	const upperText = textBeforeCursor.toUpperCase();

	// Check if we're after a table name with a dot
	// Match pattern: tablename. or "tablename". (with quotes)
	const dotMatch = textBeforeCursor.match(/(?:^|\s)([a-zA-Z_][a-zA-Z0-9_]*)\.\s*$/);
	if (dotMatch) {
		context.afterTableDot = dotMatch[1];
	}

	// Find all table names mentioned in FROM and JOIN clauses
	// Match: FROM table_name, JOIN table_name
	const fromMatches = textBeforeCursor.matchAll(/\bFROM\s+([a-zA-Z_][a-zA-Z0-9_]*)/gi);
	const joinMatches = textBeforeCursor.matchAll(/\bJOIN\s+([a-zA-Z_][a-zA-Z0-9_]*)/gi);

	for (const match of fromMatches) {
		context.availableTables.push(match[1]);
	}
	for (const match of joinMatches) {
		context.availableTables.push(match[1]);
	}

	// Determine which clause we're in
	// Split by major SQL keywords to find context
	const lastSelect = upperText.lastIndexOf('SELECT');
	const lastFrom = upperText.lastIndexOf('FROM');
	const lastWhere = upperText.lastIndexOf('WHERE');
	const lastGroupBy = upperText.lastIndexOf('GROUP BY');
	const lastOrderBy = upperText.lastIndexOf('ORDER BY');
	const lastJoin = upperText.lastIndexOf('JOIN');

	// Find the most recent clause keyword
	const clausePositions = [
		{ pos: lastSelect, clause: 'SELECT' },
		{ pos: lastFrom, clause: 'FROM' },
		{ pos: lastWhere, clause: 'WHERE' },
		{ pos: lastGroupBy, clause: 'GROUP BY' },
		{ pos: lastOrderBy, clause: 'ORDER BY' },
		{ pos: lastJoin, clause: 'JOIN' }
	].filter((c) => c.pos !== -1);

	if (clausePositions.length > 0) {
		const lastClause = clausePositions.reduce((prev, curr) => (prev.pos > curr.pos ? prev : curr));

		switch (lastClause.clause) {
			case 'FROM':
			case 'JOIN':
				context.inFromClause = true;
				break;
			case 'WHERE':
				context.inWhereClause = true;
				break;
			case 'SELECT':
				context.inSelectClause = true;
				break;
		}
	}

	return context;
}

/**
 * Get columns for a specific table
 */
function getTableColumns(tableName: string, schema?: SchemaInfo): ColumnInfo[] {
	if (!schema) return [];

	const table = schema.tables.find((t) => t.name.toLowerCase() === tableName.toLowerCase());

	return table ? table.columns : [];
}

/**
 * Setup SQL autocomplete for Monaco Editor
 *
 * Registers a completion item provider that provides:
 * - SQL keywords (SELECT, FROM, WHERE, etc.)
 * - PostgreSQL data types (INTEGER, TEXT, JSONB, etc.)
 * - SQL functions (COUNT, SUM, NOW, etc.)
 * - Schema-aware completions (tables and columns when schema is provided)
 * - Context-aware suggestions based on cursor position
 *
 * @param monaco - The Monaco Editor instance
 * @param schema - Optional schema information for table/column completion
 * @returns A disposable that can be used to unregister the provider
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function setupAutocomplete(monaco: any, schema?: SchemaInfo): any {
	return monaco.languages.registerCompletionItemProvider('sql', {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		provideCompletionItems: (model: any, position: any) => {
			// Get the word at the current cursor position
			const word = model.getWordUntilPosition(position);
			const range = {
				startLineNumber: position.lineNumber,
				endLineNumber: position.lineNumber,
				startColumn: word.startColumn,
				endColumn: word.endColumn
			};

			// Parse SQL context for smart suggestions
			const context = parseSQLContext(model, position);

			const suggestions: any[] = [];

			// If we're after "tablename.", only suggest columns from that table
			if (context.afterTableDot && schema) {
				const columns = getTableColumns(context.afterTableDot, schema);
				columns.forEach((column) => {
					suggestions.push({
						label: column.name,
						kind: monaco.languages.CompletionItemKind.Field,
						insertText: column.name,
						range: range,
						sortText: '0_' + column.name, // Prioritize columns
						detail: column.type,
						documentation: `Column: ${context.afterTableDot}.${column.name}\nType: ${column.type}\nNullable: ${column.nullable ? 'Yes' : 'No'}`
					});
				});

				// Return early - only show columns for this table
				return { suggestions };
			}

			// Context-aware suggestions based on current clause

			// In FROM clause: prioritize table names
			if (context.inFromClause && schema) {
				schema.tables.forEach((table) => {
					suggestions.push({
						label: table.name,
						kind: monaco.languages.CompletionItemKind.Class,
						insertText: table.name,
						range: range,
						sortText: '0_' + table.name, // Prioritize tables in FROM
						detail: `Table (${table.schema})`,
						documentation: `Table: ${table.schema}.${table.name}\nColumns: ${table.columns.map((c) => c.name).join(', ')}`
					});
				});

				// Add relevant keywords for FROM clause
				[
					'JOIN',
					'INNER',
					'LEFT',
					'RIGHT',
					'FULL',
					'OUTER',
					'CROSS',
					'ON',
					'USING',
					'WHERE',
					'GROUP BY',
					'ORDER BY',
					'LIMIT'
				].forEach((keyword) => {
					suggestions.push({
						label: keyword,
						kind: monaco.languages.CompletionItemKind.Keyword,
						insertText: keyword,
						range: range,
						sortText: '1_' + keyword,
						documentation: `SQL keyword: ${keyword}`
					});
				});

				return { suggestions };
			}

			// In WHERE clause: prioritize columns from available tables
			if (context.inWhereClause && schema && context.availableTables.length > 0) {
				// Add columns from all available tables
				context.availableTables.forEach((tableName) => {
					const columns = getTableColumns(tableName, schema);
					columns.forEach((column) => {
						suggestions.push({
							label: `${tableName}.${column.name}`,
							kind: monaco.languages.CompletionItemKind.Field,
							insertText: `${tableName}.${column.name}`,
							range: range,
							sortText: '0_' + tableName + '_' + column.name,
							detail: column.type,
							documentation: `Column: ${tableName}.${column.name}\nType: ${column.type}\nNullable: ${column.nullable ? 'Yes' : 'No'}`
						});
					});
				});

				// Add WHERE-specific keywords and operators
				[
					'AND',
					'OR',
					'NOT',
					'IN',
					'EXISTS',
					'BETWEEN',
					'LIKE',
					'ILIKE',
					'IS',
					'NULL',
					'TRUE',
					'FALSE'
				].forEach((keyword) => {
					suggestions.push({
						label: keyword,
						kind: monaco.languages.CompletionItemKind.Keyword,
						insertText: keyword,
						range: range,
						sortText: '1_' + keyword,
						documentation: `SQL keyword: ${keyword}`
					});
				});

				// Add comparison operators
				['=', '!=', '<>', '<', '>', '<=', '>='].forEach((op) => {
					suggestions.push({
						label: op,
						kind: monaco.languages.CompletionItemKind.Operator,
						insertText: op,
						range: range,
						sortText: '1_' + op,
						documentation: `Comparison operator: ${op}`
					});
				});

				return { suggestions };
			}

			// In SELECT clause: prioritize columns and functions
			if (context.inSelectClause && schema && context.availableTables.length > 0) {
				// Add columns from available tables
				context.availableTables.forEach((tableName) => {
					const columns = getTableColumns(tableName, schema);
					columns.forEach((column) => {
						suggestions.push({
							label: `${tableName}.${column.name}`,
							kind: monaco.languages.CompletionItemKind.Field,
							insertText: `${tableName}.${column.name}`,
							range: range,
							sortText: '0_' + tableName + '_' + column.name,
							detail: column.type,
							documentation: `Column: ${tableName}.${column.name}\nType: ${column.type}`
						});
					});
				});

				// Add aggregate functions (common in SELECT)
				AGGREGATE_FUNCTIONS.forEach((func) => {
					suggestions.push({
						label: func,
						kind: monaco.languages.CompletionItemKind.Function,
						insertText: `${func}($1)`,
						insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
						range: range,
						sortText: '0_' + func,
						documentation: `Aggregate function: ${func}()`
					});
				});

				// Add other functions
				ALL_FUNCTIONS.forEach((func) => {
					if (!AGGREGATE_FUNCTIONS.includes(func)) {
						suggestions.push({
							label: func,
							kind: monaco.languages.CompletionItemKind.Function,
							insertText: `${func}($1)`,
							insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
							range: range,
							sortText: '1_' + func,
							documentation: `SQL function: ${func}()`
						});
					}
				});

				// Add SELECT-specific keywords
				['FROM', 'AS', 'DISTINCT', 'ALL', 'CASE', 'WHEN', 'THEN', 'ELSE', 'END'].forEach(
					(keyword) => {
						suggestions.push({
							label: keyword,
							kind: monaco.languages.CompletionItemKind.Keyword,
							insertText: keyword,
							range: range,
							sortText: '2_' + keyword,
							documentation: `SQL keyword: ${keyword}`
						});
					}
				);

				return { suggestions };
			}

			// Default: provide all suggestions (no specific context detected)

			// Add SQL keywords
			ALL_KEYWORDS.forEach((keyword) => {
				suggestions.push({
					label: keyword,
					kind: monaco.languages.CompletionItemKind.Keyword,
					insertText: keyword,
					range: range,
					sortText: '0_' + keyword, // Prioritize keywords
					documentation: `SQL keyword: ${keyword}`
				});
			});

			// Add SQL functions with snippet support
			ALL_FUNCTIONS.forEach((func) => {
				suggestions.push({
					label: func,
					kind: monaco.languages.CompletionItemKind.Function,
					insertText: `${func}($1)`,
					insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
					range: range,
					sortText: '1_' + func, // Functions after keywords
					documentation: `SQL function: ${func}()`
				});
			});

			// Add schema-aware completions if schema is provided
			if (schema) {
				// Add table names
				schema.tables.forEach((table) => {
					suggestions.push({
						label: table.name,
						kind: monaco.languages.CompletionItemKind.Class,
						insertText: table.name,
						range: range,
						sortText: '2_' + table.name, // Tables after functions
						detail: `Table (${table.schema})`,
						documentation: `Table: ${table.schema}.${table.name}\nColumns: ${table.columns.map((c) => c.name).join(', ')}`
					});

					// Add columns for each table (as table.column)
					table.columns.forEach((column) => {
						suggestions.push({
							label: `${table.name}.${column.name}`,
							kind: monaco.languages.CompletionItemKind.Field,
							insertText: `${table.name}.${column.name}`,
							range: range,
							sortText: '3_' + table.name + '_' + column.name, // Columns last
							detail: column.type,
							documentation: `Column: ${table.name}.${column.name}\nType: ${column.type}\nNullable: ${column.nullable ? 'Yes' : 'No'}`
						});
					});
				});

				// Add custom functions from schema
				if (schema.functions) {
					schema.functions.forEach((func) => {
						suggestions.push({
							label: func.name,
							kind: monaco.languages.CompletionItemKind.Function,
							insertText: `${func.name}($1)`,
							insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
							range: range,
							sortText: '1_' + func.name, // With other functions
							detail: `Returns ${func.returnType}`,
							documentation: `Custom function: ${func.schema}.${func.name}\nReturns: ${func.returnType}`
						});
					});
				}
			}

			return { suggestions };
		}
	});
}
