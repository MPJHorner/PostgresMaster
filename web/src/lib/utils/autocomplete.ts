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
	'UNLISTEN',
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
	'COMPOSITE',
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
	'EVERY',
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
	'DECODE',
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
	'TIMEOFDAY',
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
	'JSONB_PRETTY',
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
	'ATANH',
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
	'NTH_VALUE',
];

/**
 * Conditional and comparison functions
 */
export const CONDITIONAL_FUNCTIONS = [
	'COALESCE',
	'NULLIF',
	'GREATEST',
	'LEAST',
	'CAST',
	'CONVERT',
];

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
	'PG_CONF_LOAD_TIME',
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
	'UNNEST',
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
	...ARRAY_FUNCTIONS,
];

/**
 * All keywords and data types combined
 */
export const ALL_KEYWORDS = [
	...SQL_KEYWORDS,
	...DATA_TYPES,
];

/**
 * Complete list of all SQL completions (keywords, types, and functions)
 */
export const ALL_COMPLETIONS = [
	...ALL_KEYWORDS,
	...ALL_FUNCTIONS,
];
