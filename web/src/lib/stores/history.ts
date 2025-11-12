import { writable } from 'svelte/store';

/**
 * Represents a single query history item
 */
export interface QueryHistoryItem {
	sql: string;
	timestamp: Date;
	success: boolean;
	executionTime?: number;
	rowCount?: number;
	error?: string;
}

const MAX_HISTORY_SIZE = 50;

function createQueryHistory() {
	const { subscribe, update } = writable<QueryHistoryItem[]>([]);

	return {
		subscribe,
		/**
		 * Add a query to the history
		 * Keeps only the last MAX_HISTORY_SIZE queries
		 */
		addQuery: (item: QueryHistoryItem) => {
			update((history) => {
				const newHistory = [item, ...history];
				// Keep only the last MAX_HISTORY_SIZE queries
				return newHistory.slice(0, MAX_HISTORY_SIZE);
			});
		},
		/**
		 * Clear all history
		 */
		clear: () => {
			update(() => []);
		}
	};
}

/**
 * Query history store
 * Stores the last 50 queries executed (in-memory only, cleared on refresh)
 */
export const queryHistory = createQueryHistory();
