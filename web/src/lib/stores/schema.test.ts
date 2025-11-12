/**
 * Schema Store Tests
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import {
	schemaStore,
	tables,
	functions,
	tableNames,
	allColumns,
	tableColumnsMap,
	functionNames,
	hasSchema,
	tableCount,
	functionCount,
	setSchema,
	clearSchema,
	getTableColumns,
	getTable,
	searchTables,
	searchColumns,
	getSchemaState
} from './schema';
import type { SchemaPayload, TableInfo, ColumnInfo, FunctionInfo } from '$lib/services/protocol';

// Test data
const mockColumn1: ColumnInfo = {
	name: 'id',
	dataType: 'integer',
	nullable: false
};

const mockColumn2: ColumnInfo = {
	name: 'name',
	dataType: 'text',
	nullable: true
};

const mockColumn3: ColumnInfo = {
	name: 'email',
	dataType: 'text',
	nullable: false
};

const mockColumn4: ColumnInfo = {
	name: 'user_id',
	dataType: 'integer',
	nullable: false
};

const mockColumn5: ColumnInfo = {
	name: 'content',
	dataType: 'text',
	nullable: true
};

const mockTable1: TableInfo = {
	schema: 'public',
	name: 'users',
	type: 'r',
	columns: [mockColumn1, mockColumn2, mockColumn3]
};

const mockTable2: TableInfo = {
	schema: 'public',
	name: 'posts',
	type: 'r',
	columns: [mockColumn1, mockColumn4, mockColumn5]
};

const mockTable3: TableInfo = {
	schema: 'auth',
	name: 'sessions',
	type: 'r',
	columns: [mockColumn1, mockColumn4]
};

const mockFunction1: FunctionInfo = {
	schema: 'public',
	name: 'get_user_count',
	returnType: 'integer'
};

const mockFunction2: FunctionInfo = {
	schema: 'auth',
	name: 'verify_token',
	returnType: 'boolean'
};

const mockSchema: SchemaPayload = {
	tables: [mockTable1, mockTable2, mockTable3],
	functions: [mockFunction1, mockFunction2]
};

describe('Schema Store', () => {
	beforeEach(() => {
		// Clear schema before each test
		clearSchema();
	});

	describe('Initial State', () => {
		it('should have empty tables array initially', () => {
			const state = get(schemaStore);
			expect(state.tables).toEqual([]);
		});

		it('should have empty functions array initially', () => {
			const state = get(schemaStore);
			expect(state.functions).toEqual([]);
		});

		it('should not have schema initially', () => {
			expect(get(hasSchema)).toBe(false);
		});

		it('should have zero table count initially', () => {
			expect(get(tableCount)).toBe(0);
		});

		it('should have zero function count initially', () => {
			expect(get(functionCount)).toBe(0);
		});
	});

	describe('setSchema', () => {
		it('should update the schema store', () => {
			setSchema(mockSchema);
			const state = get(schemaStore);
			expect(state.tables).toEqual(mockSchema.tables);
			expect(state.functions).toEqual(mockSchema.functions);
		});

		it('should update hasSchema to true', () => {
			setSchema(mockSchema);
			expect(get(hasSchema)).toBe(true);
		});

		it('should update table and function counts', () => {
			setSchema(mockSchema);
			expect(get(tableCount)).toBe(3);
			expect(get(functionCount)).toBe(2);
		});
	});

	describe('clearSchema', () => {
		it('should reset schema to empty', () => {
			setSchema(mockSchema);
			clearSchema();

			const state = get(schemaStore);
			expect(state.tables).toEqual([]);
			expect(state.functions).toEqual([]);
		});

		it('should set hasSchema to false', () => {
			setSchema(mockSchema);
			clearSchema();
			expect(get(hasSchema)).toBe(false);
		});

		it('should reset counts to zero', () => {
			setSchema(mockSchema);
			clearSchema();
			expect(get(tableCount)).toBe(0);
			expect(get(functionCount)).toBe(0);
		});
	});

	describe('Derived Store: tables', () => {
		it('should return empty array when no schema', () => {
			expect(get(tables)).toEqual([]);
		});

		it('should return tables array from schema', () => {
			setSchema(mockSchema);
			expect(get(tables)).toEqual(mockSchema.tables);
		});
	});

	describe('Derived Store: functions', () => {
		it('should return empty array when no schema', () => {
			expect(get(functions)).toEqual([]);
		});

		it('should return functions array from schema', () => {
			setSchema(mockSchema);
			expect(get(functions)).toEqual(mockSchema.functions);
		});
	});

	describe('Derived Store: tableNames', () => {
		it('should return empty array when no schema', () => {
			expect(get(tableNames)).toEqual([]);
		});

		it('should return short names for public schema tables', () => {
			setSchema(mockSchema);
			const names = get(tableNames);
			expect(names).toContain('users');
			expect(names).toContain('posts');
		});

		it('should return qualified names for non-public schema tables', () => {
			setSchema(mockSchema);
			const names = get(tableNames);
			expect(names).toContain('auth.sessions');
		});

		it('should return correct number of table names', () => {
			setSchema(mockSchema);
			expect(get(tableNames)).toHaveLength(3);
		});
	});

	describe('Derived Store: allColumns', () => {
		it('should return empty array when no schema', () => {
			expect(get(allColumns)).toEqual([]);
		});

		it('should return all columns from all tables', () => {
			setSchema(mockSchema);
			const columns = get(allColumns);

			// Should have 8 columns total (3 + 3 + 2)
			expect(columns).toHaveLength(8);
		});

		it('should include table context for each column', () => {
			setSchema(mockSchema);
			const columns = get(allColumns);

			const idColumn = columns.find((c) => c.name === 'id' && c.tableName === 'users');
			expect(idColumn).toBeDefined();
			expect(idColumn?.schemaName).toBe('public');
			expect(idColumn?.tableName).toBe('users');
		});

		it('should include full qualified name', () => {
			setSchema(mockSchema);
			const columns = get(allColumns);

			const nameColumn = columns.find((c) => c.name === 'name' && c.tableName === 'users');
			expect(nameColumn?.fullName).toBe('public.users.name');
		});

		it('should handle columns from different schemas', () => {
			setSchema(mockSchema);
			const columns = get(allColumns);

			const authColumn = columns.find((c) => c.tableName === 'sessions');
			expect(authColumn?.schemaName).toBe('auth');
			expect(authColumn?.fullName).toContain('auth.sessions');
		});
	});

	describe('Derived Store: tableColumnsMap', () => {
		it('should return empty map when no schema', () => {
			const map = get(tableColumnsMap);
			expect(map.size).toBe(0);
		});

		it('should contain entries for short table names', () => {
			setSchema(mockSchema);
			const map = get(tableColumnsMap);

			expect(map.has('users')).toBe(true);
			expect(map.has('posts')).toBe(true);
		});

		it('should contain entries for qualified table names', () => {
			setSchema(mockSchema);
			const map = get(tableColumnsMap);

			expect(map.has('public.users')).toBe(true);
			expect(map.has('auth.sessions')).toBe(true);
		});

		it('should return correct columns for table', () => {
			setSchema(mockSchema);
			const map = get(tableColumnsMap);

			const userColumns = map.get('users');
			expect(userColumns).toHaveLength(3);
			expect(userColumns?.map((c) => c.name)).toEqual(['id', 'name', 'email']);
		});
	});

	describe('Derived Store: functionNames', () => {
		it('should return empty array when no schema', () => {
			expect(get(functionNames)).toEqual([]);
		});

		it('should return short names for public schema functions', () => {
			setSchema(mockSchema);
			const names = get(functionNames);
			expect(names).toContain('get_user_count');
		});

		it('should return qualified names for non-public schema functions', () => {
			setSchema(mockSchema);
			const names = get(functionNames);
			expect(names).toContain('auth.verify_token');
		});
	});

	describe('Derived Store: hasSchema', () => {
		it('should be false when no schema', () => {
			expect(get(hasSchema)).toBe(false);
		});

		it('should be true when tables exist', () => {
			setSchema({ tables: [mockTable1], functions: [] });
			expect(get(hasSchema)).toBe(true);
		});

		it('should be true when functions exist', () => {
			setSchema({ tables: [], functions: [mockFunction1] });
			expect(get(hasSchema)).toBe(true);
		});

		it('should be true when both exist', () => {
			setSchema(mockSchema);
			expect(get(hasSchema)).toBe(true);
		});
	});

	describe('Derived Store: tableCount', () => {
		it('should return 0 when no schema', () => {
			expect(get(tableCount)).toBe(0);
		});

		it('should return correct count', () => {
			setSchema(mockSchema);
			expect(get(tableCount)).toBe(3);
		});
	});

	describe('Derived Store: functionCount', () => {
		it('should return 0 when no schema', () => {
			expect(get(functionCount)).toBe(0);
		});

		it('should return correct count', () => {
			setSchema(mockSchema);
			expect(get(functionCount)).toBe(2);
		});
	});

	describe('getTableColumns', () => {
		beforeEach(() => {
			setSchema(mockSchema);
		});

		it('should return columns for table by short name', () => {
			const columns = getTableColumns('users');
			expect(columns).toHaveLength(3);
			expect(columns.map((c) => c.name)).toEqual(['id', 'name', 'email']);
		});

		it('should return columns for table by qualified name', () => {
			const columns = getTableColumns('public.users');
			expect(columns).toHaveLength(3);
			expect(columns.map((c) => c.name)).toEqual(['id', 'name', 'email']);
		});

		it('should return columns for non-public schema table', () => {
			const columns = getTableColumns('auth.sessions');
			expect(columns).toHaveLength(2);
			expect(columns.map((c) => c.name)).toEqual(['id', 'user_id']);
		});

		it('should return empty array for non-existent table', () => {
			const columns = getTableColumns('nonexistent');
			expect(columns).toEqual([]);
		});
	});

	describe('getTable', () => {
		beforeEach(() => {
			setSchema(mockSchema);
		});

		it('should return table by short name', () => {
			const table = getTable('users');
			expect(table).not.toBeNull();
			expect(table?.name).toBe('users');
			expect(table?.schema).toBe('public');
		});

		it('should return table by qualified name', () => {
			const table = getTable('public.posts');
			expect(table).not.toBeNull();
			expect(table?.name).toBe('posts');
		});

		it('should return table from non-public schema', () => {
			const table = getTable('auth.sessions');
			expect(table).not.toBeNull();
			expect(table?.name).toBe('sessions');
			expect(table?.schema).toBe('auth');
		});

		it('should return null for non-existent table', () => {
			const table = getTable('nonexistent');
			expect(table).toBeNull();
		});
	});

	describe('searchTables', () => {
		beforeEach(() => {
			setSchema(mockSchema);
		});

		it('should return empty array when no match', () => {
			const results = searchTables('xyz');
			expect(results).toEqual([]);
		});

		it('should find tables by partial name match', () => {
			const results = searchTables('user');
			expect(results).toContain('users');
		});

		it('should be case-insensitive', () => {
			const results = searchTables('USER');
			expect(results).toContain('users');
		});

		it('should find tables by schema prefix', () => {
			const results = searchTables('auth');
			expect(results).toContain('auth.sessions');
		});

		it('should return short names for public schema', () => {
			const results = searchTables('post');
			expect(results).toContain('posts');
			expect(results).not.toContain('public.posts');
		});

		it('should return qualified names for non-public schema', () => {
			const results = searchTables('session');
			expect(results).toContain('auth.sessions');
		});
	});

	describe('searchColumns', () => {
		beforeEach(() => {
			setSchema(mockSchema);
		});

		it('should return empty array when no match', () => {
			const results = searchColumns('xyz');
			expect(results).toEqual([]);
		});

		it('should find columns by name', () => {
			const results = searchColumns('id');
			expect(results.length).toBeGreaterThan(0);
			expect(results.every((c) => c.name.includes('id'))).toBe(true);
		});

		it('should be case-insensitive', () => {
			const results = searchColumns('NAME');
			expect(results.length).toBeGreaterThan(0);
			expect(results.some((c) => c.name === 'name')).toBe(true);
		});

		it('should find columns by partial name', () => {
			const results = searchColumns('user');
			expect(results.some((c) => c.name === 'user_id')).toBe(true);
		});

		it('should include table context', () => {
			const results = searchColumns('id');
			const userIdColumn = results.find((c) => c.tableName === 'users');
			expect(userIdColumn).toBeDefined();
			expect(userIdColumn?.schemaName).toBe('public');
		});

		it('should search by full qualified name', () => {
			const results = searchColumns('public.users.email');
			expect(results.length).toBeGreaterThan(0);
			expect(results.some((c) => c.name === 'email')).toBe(true);
		});
	});

	describe('getSchemaState', () => {
		it('should return initial state when no schema', () => {
			const state = getSchemaState();
			expect(state.tables).toEqual([]);
			expect(state.functions).toEqual([]);
		});

		it('should return current state synchronously', () => {
			setSchema(mockSchema);
			const state = getSchemaState();
			expect(state.tables).toEqual(mockSchema.tables);
			expect(state.functions).toEqual(mockSchema.functions);
		});
	});

	describe('Store Reactivity', () => {
		it('should update derived stores when main store changes', () => {
			const tableCountValues: number[] = [];
			const unsubscribe = tableCount.subscribe((value) => {
				tableCountValues.push(value);
			});

			setSchema(mockSchema);

			expect(tableCountValues.length).toBeGreaterThan(1);
			expect(tableCountValues[tableCountValues.length - 1]).toBe(3);

			unsubscribe();
		});

		it('should update multiple derived stores simultaneously', () => {
			setSchema(mockSchema);

			expect(get(tableCount)).toBe(3);
			expect(get(functionCount)).toBe(2);
			expect(get(hasSchema)).toBe(true);
			expect(get(tableNames)).toHaveLength(3);
		});

		it('should update allColumns when schema changes', () => {
			setSchema({ tables: [mockTable1], functions: [] });
			expect(get(allColumns)).toHaveLength(3);

			setSchema(mockSchema);
			expect(get(allColumns)).toHaveLength(8);
		});
	});

	describe('Edge Cases', () => {
		it('should handle schema with no tables', () => {
			setSchema({ tables: [], functions: [mockFunction1] });
			expect(get(tableNames)).toEqual([]);
			expect(get(allColumns)).toEqual([]);
			expect(get(hasSchema)).toBe(true); // Still has functions
		});

		it('should handle schema with no functions', () => {
			setSchema({ tables: [mockTable1], functions: [] });
			expect(get(functionNames)).toEqual([]);
			expect(get(hasSchema)).toBe(true); // Still has tables
		});

		it('should handle tables with no columns', () => {
			const emptyTable: TableInfo = {
				schema: 'public',
				name: 'empty_table',
				type: 'r',
				columns: []
			};
			setSchema({ tables: [emptyTable], functions: [] });

			expect(get(allColumns)).toEqual([]);
			expect(getTableColumns('empty_table')).toEqual([]);
		});

		it('should handle multiple tables with same name in different schemas', () => {
			const table1: TableInfo = {
				schema: 'public',
				name: 'logs',
				type: 'r',
				columns: [mockColumn1]
			};
			const table2: TableInfo = {
				schema: 'audit',
				name: 'logs',
				type: 'r',
				columns: [mockColumn2]
			};

			setSchema({ tables: [table1, table2], functions: [] });

			expect(get(tableNames)).toContain('logs'); // public.logs as short name
			expect(get(tableNames)).toContain('audit.logs');
		});
	});
});
