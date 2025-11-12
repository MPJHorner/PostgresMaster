import { describe, it, expect } from 'vitest';
import {
	SQL_KEYWORDS,
	DATA_TYPES,
	AGGREGATE_FUNCTIONS,
	STRING_FUNCTIONS,
	DATE_TIME_FUNCTIONS,
	JSON_FUNCTIONS,
	MATH_FUNCTIONS,
	WINDOW_FUNCTIONS,
	CONDITIONAL_FUNCTIONS,
	SYSTEM_FUNCTIONS,
	ARRAY_FUNCTIONS,
	ALL_FUNCTIONS,
	ALL_KEYWORDS,
	ALL_COMPLETIONS
} from './autocomplete';

describe('SQL Keywords', () => {
	it('should include core SQL commands', () => {
		expect(SQL_KEYWORDS).toContain('SELECT');
		expect(SQL_KEYWORDS).toContain('FROM');
		expect(SQL_KEYWORDS).toContain('WHERE');
		expect(SQL_KEYWORDS).toContain('INSERT');
		expect(SQL_KEYWORDS).toContain('UPDATE');
		expect(SQL_KEYWORDS).toContain('DELETE');
		expect(SQL_KEYWORDS).toContain('CREATE');
		expect(SQL_KEYWORDS).toContain('DROP');
	});

	it('should include JOIN types', () => {
		expect(SQL_KEYWORDS).toContain('JOIN');
		expect(SQL_KEYWORDS).toContain('INNER');
		expect(SQL_KEYWORDS).toContain('LEFT');
		expect(SQL_KEYWORDS).toContain('RIGHT');
		expect(SQL_KEYWORDS).toContain('FULL');
		expect(SQL_KEYWORDS).toContain('OUTER');
		expect(SQL_KEYWORDS).toContain('CROSS');
	});

	it('should include PostgreSQL-specific keywords', () => {
		expect(SQL_KEYWORDS).toContain('RETURNING');
		expect(SQL_KEYWORDS).toContain('ON CONFLICT');
		expect(SQL_KEYWORDS).toContain('DO NOTHING');
		expect(SQL_KEYWORDS).toContain('DO UPDATE');
		expect(SQL_KEYWORDS).toContain('LATERAL');
		expect(SQL_KEYWORDS).toContain('ILIKE');
	});

	it('should include constraint keywords', () => {
		expect(SQL_KEYWORDS).toContain('PRIMARY KEY');
		expect(SQL_KEYWORDS).toContain('FOREIGN KEY');
		expect(SQL_KEYWORDS).toContain('UNIQUE');
		expect(SQL_KEYWORDS).toContain('CHECK');
		expect(SQL_KEYWORDS).toContain('DEFAULT');
		expect(SQL_KEYWORDS).toContain('NOT NULL');
	});

	it('should include window function keywords', () => {
		expect(SQL_KEYWORDS).toContain('PARTITION BY');
		expect(SQL_KEYWORDS).toContain('OVER');
		expect(SQL_KEYWORDS).toContain('ROWS');
		expect(SQL_KEYWORDS).toContain('RANGE');
	});

	it('should have a reasonable number of keywords', () => {
		expect(SQL_KEYWORDS.length).toBeGreaterThan(50);
	});

	it('should not have duplicate keywords', () => {
		const uniqueKeywords = new Set(SQL_KEYWORDS);
		expect(uniqueKeywords.size).toBe(SQL_KEYWORDS.length);
	});
});

describe('Data Types', () => {
	it('should include numeric types', () => {
		expect(DATA_TYPES).toContain('INTEGER');
		expect(DATA_TYPES).toContain('BIGINT');
		expect(DATA_TYPES).toContain('DECIMAL');
		expect(DATA_TYPES).toContain('NUMERIC');
		expect(DATA_TYPES).toContain('REAL');
		expect(DATA_TYPES).toContain('SERIAL');
	});

	it('should include character types', () => {
		expect(DATA_TYPES).toContain('VARCHAR');
		expect(DATA_TYPES).toContain('CHAR');
		expect(DATA_TYPES).toContain('TEXT');
	});

	it('should include date/time types', () => {
		expect(DATA_TYPES).toContain('DATE');
		expect(DATA_TYPES).toContain('TIME');
		expect(DATA_TYPES).toContain('TIMESTAMP');
		expect(DATA_TYPES).toContain('TIMESTAMPTZ');
		expect(DATA_TYPES).toContain('INTERVAL');
	});

	it('should include PostgreSQL-specific types', () => {
		expect(DATA_TYPES).toContain('UUID');
		expect(DATA_TYPES).toContain('JSON');
		expect(DATA_TYPES).toContain('JSONB');
		expect(DATA_TYPES).toContain('ARRAY');
		expect(DATA_TYPES).toContain('BYTEA');
	});

	it('should include network types', () => {
		expect(DATA_TYPES).toContain('INET');
		expect(DATA_TYPES).toContain('CIDR');
		expect(DATA_TYPES).toContain('MACADDR');
	});

	it('should include range types', () => {
		expect(DATA_TYPES).toContain('INT4RANGE');
		expect(DATA_TYPES).toContain('INT8RANGE');
		expect(DATA_TYPES).toContain('TSRANGE');
		expect(DATA_TYPES).toContain('DATERANGE');
	});

	it('should not have duplicate types', () => {
		const uniqueTypes = new Set(DATA_TYPES);
		expect(uniqueTypes.size).toBe(DATA_TYPES.length);
	});
});

describe('Aggregate Functions', () => {
	it('should include standard aggregate functions', () => {
		expect(AGGREGATE_FUNCTIONS).toContain('COUNT');
		expect(AGGREGATE_FUNCTIONS).toContain('SUM');
		expect(AGGREGATE_FUNCTIONS).toContain('AVG');
		expect(AGGREGATE_FUNCTIONS).toContain('MIN');
		expect(AGGREGATE_FUNCTIONS).toContain('MAX');
	});

	it('should include PostgreSQL-specific aggregate functions', () => {
		expect(AGGREGATE_FUNCTIONS).toContain('ARRAY_AGG');
		expect(AGGREGATE_FUNCTIONS).toContain('JSON_AGG');
		expect(AGGREGATE_FUNCTIONS).toContain('JSONB_AGG');
		expect(AGGREGATE_FUNCTIONS).toContain('STRING_AGG');
	});

	it('should include boolean aggregates', () => {
		expect(AGGREGATE_FUNCTIONS).toContain('BOOL_AND');
		expect(AGGREGATE_FUNCTIONS).toContain('BOOL_OR');
		expect(AGGREGATE_FUNCTIONS).toContain('EVERY');
	});
});

describe('String Functions', () => {
	it('should include basic string functions', () => {
		expect(STRING_FUNCTIONS).toContain('CONCAT');
		expect(STRING_FUNCTIONS).toContain('LENGTH');
		expect(STRING_FUNCTIONS).toContain('LOWER');
		expect(STRING_FUNCTIONS).toContain('UPPER');
		expect(STRING_FUNCTIONS).toContain('TRIM');
		expect(STRING_FUNCTIONS).toContain('SUBSTRING');
		expect(STRING_FUNCTIONS).toContain('REPLACE');
	});

	it('should include regex functions', () => {
		expect(STRING_FUNCTIONS).toContain('REGEXP_REPLACE');
		expect(STRING_FUNCTIONS).toContain('REGEXP_MATCH');
		expect(STRING_FUNCTIONS).toContain('REGEXP_MATCHES');
		expect(STRING_FUNCTIONS).toContain('REGEXP_SPLIT_TO_ARRAY');
	});

	it('should include formatting functions', () => {
		expect(STRING_FUNCTIONS).toContain('FORMAT');
		expect(STRING_FUNCTIONS).toContain('QUOTE_IDENT');
		expect(STRING_FUNCTIONS).toContain('QUOTE_LITERAL');
	});
});

describe('Date/Time Functions', () => {
	it('should include current date/time functions', () => {
		expect(DATE_TIME_FUNCTIONS).toContain('NOW');
		expect(DATE_TIME_FUNCTIONS).toContain('CURRENT_DATE');
		expect(DATE_TIME_FUNCTIONS).toContain('CURRENT_TIME');
		expect(DATE_TIME_FUNCTIONS).toContain('CURRENT_TIMESTAMP');
	});

	it('should include date manipulation functions', () => {
		expect(DATE_TIME_FUNCTIONS).toContain('DATE_TRUNC');
		expect(DATE_TIME_FUNCTIONS).toContain('DATE_PART');
		expect(DATE_TIME_FUNCTIONS).toContain('EXTRACT');
		expect(DATE_TIME_FUNCTIONS).toContain('AGE');
	});

	it('should include formatting functions', () => {
		expect(DATE_TIME_FUNCTIONS).toContain('TO_CHAR');
		expect(DATE_TIME_FUNCTIONS).toContain('TO_DATE');
		expect(DATE_TIME_FUNCTIONS).toContain('TO_TIMESTAMP');
	});

	it('should include constructor functions', () => {
		expect(DATE_TIME_FUNCTIONS).toContain('MAKE_DATE');
		expect(DATE_TIME_FUNCTIONS).toContain('MAKE_TIME');
		expect(DATE_TIME_FUNCTIONS).toContain('MAKE_TIMESTAMP');
		expect(DATE_TIME_FUNCTIONS).toContain('MAKE_INTERVAL');
	});
});

describe('JSON Functions', () => {
	it('should include JSON builders', () => {
		expect(JSON_FUNCTIONS).toContain('JSON_BUILD_ARRAY');
		expect(JSON_FUNCTIONS).toContain('JSON_BUILD_OBJECT');
		expect(JSON_FUNCTIONS).toContain('JSONB_BUILD_ARRAY');
		expect(JSON_FUNCTIONS).toContain('JSONB_BUILD_OBJECT');
	});

	it('should include JSON extractors', () => {
		expect(JSON_FUNCTIONS).toContain('JSON_EXTRACT_PATH');
		expect(JSON_FUNCTIONS).toContain('JSON_EXTRACT_PATH_TEXT');
		expect(JSON_FUNCTIONS).toContain('JSONB_EXTRACT_PATH');
	});

	it('should include JSON conversion functions', () => {
		expect(JSON_FUNCTIONS).toContain('JSON_TO_RECORD');
		expect(JSON_FUNCTIONS).toContain('JSONB_TO_RECORD');
		expect(JSON_FUNCTIONS).toContain('JSON_TO_RECORDSET');
		expect(JSON_FUNCTIONS).toContain('JSONB_TO_RECORDSET');
	});

	it('should include JSONB path functions', () => {
		expect(JSON_FUNCTIONS).toContain('JSONB_PATH_EXISTS');
		expect(JSON_FUNCTIONS).toContain('JSONB_PATH_QUERY');
		expect(JSON_FUNCTIONS).toContain('JSONB_PATH_QUERY_ARRAY');
	});

	it('should include JSONB set functions', () => {
		expect(JSON_FUNCTIONS).toContain('JSONB_SET');
		expect(JSON_FUNCTIONS).toContain('JSONB_INSERT');
	});
});

describe('Math Functions', () => {
	it('should include basic math functions', () => {
		expect(MATH_FUNCTIONS).toContain('ABS');
		expect(MATH_FUNCTIONS).toContain('CEIL');
		expect(MATH_FUNCTIONS).toContain('FLOOR');
		expect(MATH_FUNCTIONS).toContain('ROUND');
		expect(MATH_FUNCTIONS).toContain('MOD');
	});

	it('should include exponential functions', () => {
		expect(MATH_FUNCTIONS).toContain('EXP');
		expect(MATH_FUNCTIONS).toContain('LN');
		expect(MATH_FUNCTIONS).toContain('LOG');
		expect(MATH_FUNCTIONS).toContain('POWER');
		expect(MATH_FUNCTIONS).toContain('SQRT');
	});

	it('should include trigonometric functions', () => {
		expect(MATH_FUNCTIONS).toContain('SIN');
		expect(MATH_FUNCTIONS).toContain('COS');
		expect(MATH_FUNCTIONS).toContain('TAN');
		expect(MATH_FUNCTIONS).toContain('ASIN');
		expect(MATH_FUNCTIONS).toContain('ACOS');
		expect(MATH_FUNCTIONS).toContain('ATAN');
	});

	it('should include random functions', () => {
		expect(MATH_FUNCTIONS).toContain('RANDOM');
		expect(MATH_FUNCTIONS).toContain('SETSEED');
	});
});

describe('Window Functions', () => {
	it('should include ranking functions', () => {
		expect(WINDOW_FUNCTIONS).toContain('ROW_NUMBER');
		expect(WINDOW_FUNCTIONS).toContain('RANK');
		expect(WINDOW_FUNCTIONS).toContain('DENSE_RANK');
		expect(WINDOW_FUNCTIONS).toContain('PERCENT_RANK');
		expect(WINDOW_FUNCTIONS).toContain('CUME_DIST');
		expect(WINDOW_FUNCTIONS).toContain('NTILE');
	});

	it('should include value functions', () => {
		expect(WINDOW_FUNCTIONS).toContain('LAG');
		expect(WINDOW_FUNCTIONS).toContain('LEAD');
		expect(WINDOW_FUNCTIONS).toContain('FIRST_VALUE');
		expect(WINDOW_FUNCTIONS).toContain('LAST_VALUE');
		expect(WINDOW_FUNCTIONS).toContain('NTH_VALUE');
	});
});

describe('Conditional Functions', () => {
	it('should include common conditional functions', () => {
		expect(CONDITIONAL_FUNCTIONS).toContain('COALESCE');
		expect(CONDITIONAL_FUNCTIONS).toContain('NULLIF');
		expect(CONDITIONAL_FUNCTIONS).toContain('GREATEST');
		expect(CONDITIONAL_FUNCTIONS).toContain('LEAST');
		expect(CONDITIONAL_FUNCTIONS).toContain('CAST');
	});
});

describe('System Functions', () => {
	it('should include version and user info functions', () => {
		expect(SYSTEM_FUNCTIONS).toContain('VERSION');
		expect(SYSTEM_FUNCTIONS).toContain('CURRENT_DATABASE');
		expect(SYSTEM_FUNCTIONS).toContain('CURRENT_USER');
		expect(SYSTEM_FUNCTIONS).toContain('SESSION_USER');
	});

	it('should include network info functions', () => {
		expect(SYSTEM_FUNCTIONS).toContain('INET_CLIENT_ADDR');
		expect(SYSTEM_FUNCTIONS).toContain('INET_SERVER_ADDR');
	});

	it('should include backend info functions', () => {
		expect(SYSTEM_FUNCTIONS).toContain('PG_BACKEND_PID');
		expect(SYSTEM_FUNCTIONS).toContain('PG_POSTMASTER_START_TIME');
	});
});

describe('Array Functions', () => {
	it('should include array manipulation functions', () => {
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_APPEND');
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_PREPEND');
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_CAT');
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_REMOVE');
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_REPLACE');
	});

	it('should include array query functions', () => {
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_LENGTH');
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_POSITION');
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_POSITIONS');
		expect(ARRAY_FUNCTIONS).toContain('CARDINALITY');
	});

	it('should include array conversion functions', () => {
		expect(ARRAY_FUNCTIONS).toContain('ARRAY_TO_STRING');
		expect(ARRAY_FUNCTIONS).toContain('STRING_TO_ARRAY');
		expect(ARRAY_FUNCTIONS).toContain('UNNEST');
	});
});

describe('Combined Collections', () => {
	it('should combine all functions correctly', () => {
		const expectedFunctionCount =
			AGGREGATE_FUNCTIONS.length +
			STRING_FUNCTIONS.length +
			DATE_TIME_FUNCTIONS.length +
			JSON_FUNCTIONS.length +
			MATH_FUNCTIONS.length +
			WINDOW_FUNCTIONS.length +
			CONDITIONAL_FUNCTIONS.length +
			SYSTEM_FUNCTIONS.length +
			ARRAY_FUNCTIONS.length;

		expect(ALL_FUNCTIONS.length).toBe(expectedFunctionCount);
	});

	it('should combine all keywords correctly', () => {
		const expectedKeywordCount = SQL_KEYWORDS.length + DATA_TYPES.length;
		expect(ALL_KEYWORDS.length).toBe(expectedKeywordCount);
	});

	it('should combine all completions correctly', () => {
		const expectedCompletionCount = ALL_KEYWORDS.length + ALL_FUNCTIONS.length;
		expect(ALL_COMPLETIONS.length).toBe(expectedCompletionCount);
	});

	it('should have a comprehensive set of completions', () => {
		// Should have at least 300 total completions (keywords + functions)
		expect(ALL_COMPLETIONS.length).toBeGreaterThan(300);
	});

	it('should not have duplicates in ALL_COMPLETIONS', () => {
		const uniqueCompletions = new Set(ALL_COMPLETIONS);
		expect(uniqueCompletions.size).toBe(ALL_COMPLETIONS.length);
	});
});

describe('Specific PostgreSQL Features', () => {
	it('should support UPSERT syntax', () => {
		expect(SQL_KEYWORDS).toContain('ON CONFLICT');
		expect(SQL_KEYWORDS).toContain('DO NOTHING');
		expect(SQL_KEYWORDS).toContain('DO UPDATE');
		expect(SQL_KEYWORDS).toContain('EXCLUDED');
	});

	it('should support CTEs', () => {
		expect(SQL_KEYWORDS).toContain('WITH');
		expect(SQL_KEYWORDS).toContain('RECURSIVE');
	});

	it('should support full-text search', () => {
		expect(DATA_TYPES).toContain('TSVECTOR');
		expect(DATA_TYPES).toContain('TSQUERY');
	});

	it('should support JSON operations', () => {
		expect(DATA_TYPES).toContain('JSON');
		expect(DATA_TYPES).toContain('JSONB');
		expect(JSON_FUNCTIONS.length).toBeGreaterThan(20);
	});

	it('should support array operations', () => {
		expect(DATA_TYPES).toContain('ARRAY');
		expect(ARRAY_FUNCTIONS.length).toBeGreaterThan(10);
	});

	it('should support window functions', () => {
		expect(SQL_KEYWORDS).toContain('PARTITION BY');
		expect(SQL_KEYWORDS).toContain('OVER');
		expect(WINDOW_FUNCTIONS.length).toBeGreaterThan(5);
	});
});

describe('SQL Standards Compliance', () => {
	it('should support standard SQL clauses', () => {
		const standardClauses = ['SELECT', 'FROM', 'WHERE', 'GROUP BY', 'HAVING', 'ORDER BY', 'LIMIT'];

		standardClauses.forEach((clause) => {
			expect(SQL_KEYWORDS).toContain(clause);
		});
	});

	it('should support standard SQL joins', () => {
		const standardJoins = ['INNER', 'LEFT', 'RIGHT', 'FULL', 'CROSS'];

		standardJoins.forEach((join) => {
			expect(SQL_KEYWORDS).toContain(join);
		});
	});

	it('should support standard SQL data types', () => {
		const standardTypes = [
			'INTEGER',
			'VARCHAR',
			'CHAR',
			'DATE',
			'TIME',
			'TIMESTAMP',
			'BOOLEAN',
			'DECIMAL'
		];

		standardTypes.forEach((type) => {
			expect(DATA_TYPES).toContain(type);
		});
	});

	it('should support standard aggregate functions', () => {
		const standardAggregates = ['COUNT', 'SUM', 'AVG', 'MIN', 'MAX'];

		standardAggregates.forEach((func) => {
			expect(AGGREGATE_FUNCTIONS).toContain(func);
		});
	});
});
