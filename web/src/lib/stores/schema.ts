/**
 * Schema Store
 * Manages database schema information for autocomplete and introspection
 */

import { writable, derived, type Readable } from 'svelte/store';
import type { SchemaPayload, TableInfo, ColumnInfo, FunctionInfo } from '$lib/services/protocol';

/**
 * Column information with table context for autocomplete
 */
export interface ColumnWithTable extends ColumnInfo {
	/** Table name this column belongs to */
	tableName: string;
	/** Schema name this column belongs to */
	schemaName: string;
	/** Full qualified name: schema.table.column */
	fullName: string;
}

/**
 * Initial empty schema state
 */
const initialSchema: SchemaPayload = {
	tables: [],
	functions: []
};

/**
 * Main schema store
 * Holds the database schema information (tables, columns, functions)
 */
export const schemaStore = writable<SchemaPayload>(initialSchema);

/**
 * Derived store: provides the tables array
 */
export const tables: Readable<TableInfo[]> = derived(schemaStore, ($schema) => $schema.tables);

/**
 * Derived store: provides the functions array
 */
export const functions: Readable<FunctionInfo[]> = derived(
	schemaStore,
	($schema) => $schema.functions
);

/**
 * Derived store: provides an array of table names for autocomplete
 * Returns format: "schema.table" or just "table" for public schema
 */
export const tableNames: Readable<string[]> = derived(schemaStore, ($schema) =>
	$schema.tables.map((table) => {
		// Use short name for public schema, full name for others
		if (table.schema === 'public') {
			return table.name;
		}
		return `${table.schema}.${table.name}`;
	})
);

/**
 * Derived store: provides all columns from all tables for autocomplete
 * Each column includes table context and full qualified name
 */
export const allColumns: Readable<ColumnWithTable[]> = derived(schemaStore, ($schema) => {
	const columns: ColumnWithTable[] = [];

	for (const table of $schema.tables) {
		for (const column of table.columns) {
			columns.push({
				...column,
				tableName: table.name,
				schemaName: table.schema,
				fullName: `${table.schema}.${table.name}.${column.name}`
			});
		}
	}

	return columns;
});

/**
 * Derived store: provides a map of table names to their columns
 * Useful for context-aware autocomplete (e.g., "SELECT * FROM users u WHERE u.|")
 */
export const tableColumnsMap: Readable<Map<string, ColumnInfo[]>> = derived(
	schemaStore,
	($schema) => {
		const map = new Map<string, ColumnInfo[]>();

		for (const table of $schema.tables) {
			// Add entry for short name (table only)
			map.set(table.name, table.columns);

			// Add entry for full name (schema.table)
			map.set(`${table.schema}.${table.name}`, table.columns);
		}

		return map;
	}
);

/**
 * Derived store: provides function names for autocomplete
 * Returns format: "schema.function" or just "function" for public schema
 */
export const functionNames: Readable<string[]> = derived(schemaStore, ($schema) =>
	$schema.functions.map((func) => {
		// Use short name for public schema, full name for others
		if (func.schema === 'public') {
			return func.name;
		}
		return `${func.schema}.${func.name}`;
	})
);

/**
 * Derived store: indicates if schema has been loaded
 */
export const hasSchema: Readable<boolean> = derived(
	schemaStore,
	($schema) => $schema.tables.length > 0 || $schema.functions.length > 0
);

/**
 * Derived store: provides count of tables in schema
 */
export const tableCount: Readable<number> = derived(
	schemaStore,
	($schema) => $schema.tables.length
);

/**
 * Derived store: provides count of functions in schema
 */
export const functionCount: Readable<number> = derived(
	schemaStore,
	($schema) => $schema.functions.length
);

/**
 * Updates the schema store with new schema information
 * @param schema Schema payload from introspection
 */
export function setSchema(schema: SchemaPayload): void {
	schemaStore.set(schema);
}

/**
 * Clears the schema store (resets to empty)
 */
export function clearSchema(): void {
	schemaStore.set(initialSchema);
}

/**
 * Gets columns for a specific table
 * @param tableName Name of the table (can be "table" or "schema.table")
 * @returns Array of columns or empty array if table not found
 */
export function getTableColumns(tableName: string): ColumnInfo[] {
	let columns: ColumnInfo[] = [];

	schemaStore.subscribe(($schema) => {
		const table = $schema.tables.find(
			(t) => t.name === tableName || `${t.schema}.${t.name}` === tableName
		);
		columns = table?.columns || [];
	})();

	return columns;
}

/**
 * Gets table information by name
 * @param tableName Name of the table (can be "table" or "schema.table")
 * @returns TableInfo or null if not found
 */
export function getTable(tableName: string): TableInfo | null {
	let table: TableInfo | null = null;

	schemaStore.subscribe(($schema) => {
		table =
			$schema.tables.find(
				(t) => t.name === tableName || `${t.schema}.${t.name}` === tableName
			) || null;
	})();

	return table;
}

/**
 * Searches for tables matching a pattern (case-insensitive)
 * @param pattern Search pattern (partial match)
 * @returns Array of matching table names
 */
export function searchTables(pattern: string): string[] {
	const lowerPattern = pattern.toLowerCase();
	let matches: string[] = [];

	schemaStore.subscribe(($schema) => {
		matches = $schema.tables
			.filter(
				(t) =>
					t.name.toLowerCase().includes(lowerPattern) ||
					`${t.schema}.${t.name}`.toLowerCase().includes(lowerPattern)
			)
			.map((t) => (t.schema === 'public' ? t.name : `${t.schema}.${t.name}`));
	})();

	return matches;
}

/**
 * Searches for columns matching a pattern (case-insensitive)
 * @param pattern Search pattern (partial match)
 * @returns Array of matching columns with table context
 */
export function searchColumns(pattern: string): ColumnWithTable[] {
	const lowerPattern = pattern.toLowerCase();
	let matches: ColumnWithTable[] = [];

	allColumns.subscribe(($columns) => {
		matches = $columns.filter(
			(c) =>
				c.name.toLowerCase().includes(lowerPattern) ||
				c.fullName.toLowerCase().includes(lowerPattern)
		);
	})();

	return matches;
}

/**
 * Gets the current schema state synchronously
 * Note: Prefer using the derived stores in components
 */
export function getSchemaState(): SchemaPayload {
	let state: SchemaPayload = initialSchema;
	schemaStore.subscribe((s) => {
		state = s;
	})();
	return state;
}
