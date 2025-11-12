import { describe, it, expect } from 'vitest';
import {
	formatValue,
	formatColumnType,
	formatExecutionTime,
	formatRowCount
} from './format';

describe('formatValue', () => {
	describe('null and undefined', () => {
		it('should format null as "NULL"', () => {
			expect(formatValue(null)).toBe('NULL');
		});

		it('should format undefined as "NULL"', () => {
			expect(formatValue(undefined)).toBe('NULL');
		});
	});

	describe('boolean values', () => {
		it('should format true as "true"', () => {
			expect(formatValue(true)).toBe('true');
		});

		it('should format false as "false"', () => {
			expect(formatValue(false)).toBe('false');
		});
	});

	describe('number values', () => {
		it('should format integer as string', () => {
			expect(formatValue(42)).toBe('42');
		});

		it('should format negative number as string', () => {
			expect(formatValue(-123)).toBe('-123');
		});

		it('should format decimal as string', () => {
			expect(formatValue(3.14159)).toBe('3.14159');
		});

		it('should format zero as "0"', () => {
			expect(formatValue(0)).toBe('0');
		});

		it('should format large number as string', () => {
			expect(formatValue(1234567890)).toBe('1234567890');
		});

		it('should format Infinity', () => {
			expect(formatValue(Infinity)).toBe('Infinity');
		});

		it('should format -Infinity', () => {
			expect(formatValue(-Infinity)).toBe('-Infinity');
		});

		it('should format NaN', () => {
			expect(formatValue(NaN)).toBe('NaN');
		});
	});

	describe('string values', () => {
		it('should return string as-is', () => {
			expect(formatValue('hello world')).toBe('hello world');
		});

		it('should return empty string as-is', () => {
			expect(formatValue('')).toBe('');
		});

		it('should handle strings with special characters', () => {
			expect(formatValue('hello\nworld\ttab')).toBe('hello\nworld\ttab');
		});

		it('should handle strings with quotes', () => {
			expect(formatValue('hello "world"')).toBe('hello "world"');
		});

		it('should handle unicode strings', () => {
			expect(formatValue('Hello 世界 🌍')).toBe('Hello 世界 🌍');
		});
	});

	describe('Date values', () => {
		it('should format Date as ISO string', () => {
			const date = new Date('2024-01-15T10:30:00.000Z');
			expect(formatValue(date)).toBe('2024-01-15T10:30:00.000Z');
		});

		it('should format current date as ISO string', () => {
			const date = new Date();
			const expected = date.toISOString();
			expect(formatValue(date)).toBe(expected);
		});
	});

	describe('Array values', () => {
		it('should format empty array as JSON', () => {
			expect(formatValue([])).toBe('[]');
		});

		it('should format simple array as JSON', () => {
			expect(formatValue([1, 2, 3])).toBe('[1,2,3]');
		});

		it('should format array with mixed types', () => {
			expect(formatValue([1, 'hello', true, null])).toBe('[1,"hello",true,null]');
		});

		it('should format nested arrays', () => {
			expect(formatValue([[1, 2], [3, 4]])).toBe('[[1,2],[3,4]]');
		});

		it('should format array with objects', () => {
			const result = formatValue([{ id: 1 }, { id: 2 }]);
			expect(result).toBe('[{"id":1},{"id":2}]');
		});
	});

	describe('Object values', () => {
		it('should format empty object as pretty JSON', () => {
			expect(formatValue({})).toBe('{}');
		});

		it('should format simple object as pretty JSON', () => {
			const obj = { name: 'John', age: 30 };
			const result = formatValue(obj);
			expect(result).toContain('"name": "John"');
			expect(result).toContain('"age": 30');
		});

		it('should format nested object as pretty JSON', () => {
			const obj = {
				user: {
					name: 'John',
					address: {
						city: 'NYC'
					}
				}
			};
			const result = formatValue(obj);
			expect(result).toContain('"user"');
			expect(result).toContain('"name": "John"');
			expect(result).toContain('"city": "NYC"');
		});

		it('should format object with null values', () => {
			const obj = { name: 'John', middleName: null };
			const result = formatValue(obj);
			expect(result).toContain('"middleName": null');
		});

		it('should format object with array values', () => {
			const obj = { tags: ['tag1', 'tag2'] };
			const result = formatValue(obj);
			expect(result).toContain('"tags"');
			expect(result).toContain('tag1');
			expect(result).toContain('tag2');
		});
	});
});

describe('formatColumnType', () => {
	it('should simplify "character varying" to "varchar"', () => {
		expect(formatColumnType('character varying')).toBe('varchar');
	});

	it('should simplify "CHARACTER VARYING" (case insensitive)', () => {
		expect(formatColumnType('CHARACTER VARYING')).toBe('varchar');
	});

	it('should simplify "timestamp without time zone"', () => {
		expect(formatColumnType('timestamp without time zone')).toBe('timestamp');
	});

	it('should simplify "timestamp with time zone" to "timestamptz"', () => {
		expect(formatColumnType('timestamp with time zone')).toBe('timestamptz');
	});

	it('should simplify "double precision"', () => {
		expect(formatColumnType('double precision')).toBe('double');
	});

	it('should remove length specifiers like (255)', () => {
		expect(formatColumnType('varchar(255)')).toBe('varchar');
	});

	it('should remove length specifiers from numeric types', () => {
		expect(formatColumnType('numeric(10,2)')).toBe('numeric');
	});

	it('should handle multiple simplifications', () => {
		expect(formatColumnType('character varying(255)')).toBe('varchar');
	});

	it('should leave simple types unchanged', () => {
		expect(formatColumnType('integer')).toBe('integer');
		expect(formatColumnType('text')).toBe('text');
		expect(formatColumnType('boolean')).toBe('boolean');
	});
});

describe('formatExecutionTime', () => {
	it('should format time less than 1ms', () => {
		expect(formatExecutionTime(0.5)).toBe('< 1ms');
		expect(formatExecutionTime(0)).toBe('< 1ms');
	});

	it('should format time in milliseconds (< 1000ms)', () => {
		expect(formatExecutionTime(1)).toBe('1ms');
		expect(formatExecutionTime(50)).toBe('50ms');
		expect(formatExecutionTime(999)).toBe('999ms');
	});

	it('should round milliseconds to nearest integer', () => {
		expect(formatExecutionTime(45.7)).toBe('46ms');
		expect(formatExecutionTime(123.4)).toBe('123ms');
	});

	it('should format time in seconds (>= 1000ms)', () => {
		expect(formatExecutionTime(1000)).toBe('1.00s');
		expect(formatExecutionTime(1500)).toBe('1.50s');
		expect(formatExecutionTime(2345)).toBe('2.35s');
	});

	it('should format large times', () => {
		expect(formatExecutionTime(10000)).toBe('10.00s');
		expect(formatExecutionTime(60000)).toBe('60.00s');
	});

	it('should format with 2 decimal places for seconds', () => {
		expect(formatExecutionTime(1234)).toBe('1.23s');
		expect(formatExecutionTime(5678)).toBe('5.68s');
	});
});

describe('formatRowCount', () => {
	it('should format zero rows', () => {
		expect(formatRowCount(0)).toBe('No rows');
	});

	it('should format one row (singular)', () => {
		expect(formatRowCount(1)).toBe('1 row');
	});

	it('should format multiple rows (plural)', () => {
		expect(formatRowCount(2)).toBe('2 rows');
		expect(formatRowCount(10)).toBe('10 rows');
		expect(formatRowCount(100)).toBe('100 rows');
	});

	it('should format large numbers with locale formatting', () => {
		expect(formatRowCount(1000)).toBe('1,000 rows');
		expect(formatRowCount(1234567)).toBe('1,234,567 rows');
	});
});
