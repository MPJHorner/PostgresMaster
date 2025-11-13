<script lang="ts">
	import { queryHistory } from '$lib/stores/history';
	import { Button, Card, Badge } from '$lib/components/ui';
	import { Clock, CheckCircle, XCircle, RotateCcw } from 'lucide-svelte';

	/**
	 * Callback when a query is clicked to load it into the editor
	 */
	export let onQuerySelect: (sql: string) => void = () => {};

	/**
	 * Maximum number of queries to display
	 */
	const MAX_DISPLAY = 20;

	$: recentQueries = $queryHistory.slice(0, MAX_DISPLAY);

	/**
	 * Format timestamp to relative time
	 */
	function formatTimestamp(date: Date): string {
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		const diffHours = Math.floor(diffMs / 3600000);
		const diffDays = Math.floor(diffMs / 86400000);

		if (diffMins < 1) return 'just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		if (diffHours < 24) return `${diffHours}h ago`;
		return `${diffDays}d ago`;
	}

	/**
	 * Truncate SQL for display
	 */
	function truncateSQL(sql: string, maxLength: number = 60): string {
		const singleLine = sql.replace(/\s+/g, ' ').trim();
		if (singleLine.length <= maxLength) return singleLine;
		return singleLine.substring(0, maxLength) + '...';
	}

	/**
	 * Clear all history
	 */
	function handleClearHistory() {
		if (confirm('Clear all query history?')) {
			queryHistory.clear();
		}
	}
</script>

<Card class="p-4">
	<div class="mb-4 flex items-center justify-between">
		<h3 class="text-lg font-semibold">Query History</h3>
		<div class="flex gap-2">
			{#if recentQueries.length > 0}
				<Button variant="ghost" size="sm" onclick={handleClearHistory}>
					<RotateCcw class="mr-2 h-4 w-4" />
					Clear
				</Button>
			{/if}
		</div>
	</div>

	{#if recentQueries.length === 0}
		<div class="py-8 text-center text-sm text-muted-foreground">
			No queries executed yet. Run a query to see it here.
		</div>
	{:else}
		<div class="space-y-2">
			{#each recentQueries as item (item.timestamp.getTime())}
				<button
					class="w-full rounded-lg border p-3 text-left transition-colors hover:bg-accent"
					on:click={() => onQuerySelect(item.sql)}
				>
					<div class="mb-2 flex items-start justify-between gap-2">
						<code class="flex-1 text-sm">{truncateSQL(item.sql)}</code>
						{#if item.success}
							<Badge variant="default" class="flex items-center gap-1 bg-green-600">
								<CheckCircle class="h-3 w-3" />
								Success
							</Badge>
						{:else}
							<Badge variant="destructive" class="flex items-center gap-1">
								<XCircle class="h-3 w-3" />
								Error
							</Badge>
						{/if}
					</div>

					<div class="flex items-center gap-4 text-xs text-muted-foreground">
						<span class="flex items-center gap-1">
							<Clock class="h-3 w-3" />
							{formatTimestamp(item.timestamp)}
						</span>
						{#if item.executionTime !== undefined}
							<span>{item.executionTime}ms</span>
						{/if}
						{#if item.rowCount !== undefined}
							<span>{item.rowCount} {item.rowCount === 1 ? 'row' : 'rows'}</span>
						{/if}
					</div>

					{#if !item.success && item.error}
						<div class="mt-2 text-xs text-destructive">
							{truncateSQL(item.error, 80)}
						</div>
					{/if}
				</button>
			{/each}
		</div>
	{/if}
</Card>
