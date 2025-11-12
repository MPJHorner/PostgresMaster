<script lang="ts">
	import type { ResultPayload } from '$lib/services/protocol';
	import {
		formatValue,
		formatColumnType,
		formatExecutionTime,
		formatRowCount
	} from '$lib/utils/format';
	import Card from './ui/card.svelte';
	import CardContent from './ui/card-content.svelte';
	import CardHeader from './ui/card-header.svelte';
	import CardTitle from './ui/card-title.svelte';
	import Table from './ui/table.svelte';
	import TableBody from './ui/table-body.svelte';
	import TableCell from './ui/table-cell.svelte';
	import TableHead from './ui/table-head.svelte';
	import TableHeader from './ui/table-header.svelte';
	import TableRow from './ui/table-row.svelte';
	import Badge from './ui/badge.svelte';

	interface ResultsProps {
		data: ResultPayload;
	}

	let { data }: ResultsProps = $props();

	/**
	 * Check if a value is NULL/undefined for special styling
	 */
	function isNull(value: unknown): boolean {
		return value === null || value === undefined;
	}
</script>

<Card>
	<CardHeader>
		<CardTitle>
			<div class="flex items-center justify-between">
				<span>Query Results</span>
				<div class="flex items-center gap-4 text-sm font-normal text-slate-500">
					<span>{formatRowCount(data.rowCount)}</span>
					<span>•</span>
					<span>{formatExecutionTime(data.executionTime)}</span>
				</div>
			</div>
		</CardTitle>
	</CardHeader>
	<CardContent>
		{#if data.rowCount === 0}
			<div class="flex items-center justify-center py-8 text-slate-500">
				<p>No rows returned</p>
			</div>
		{:else}
			<div class="max-h-[500px] overflow-auto rounded-md border">
				<Table>
					<TableHeader class="sticky top-0 bg-white dark:bg-slate-950 z-10">
						<TableRow>
							{#each data.columns as column}
								<TableHead>
									<div class="flex flex-col gap-1">
										<span class="font-semibold">{column.name}</span>
										<Badge variant="secondary" class="w-fit text-xs">
											{formatColumnType(column.dataType)}
										</Badge>
									</div>
								</TableHead>
							{/each}
						</TableRow>
					</TableHeader>
					<TableBody>
						{#each data.rows as row}
							<TableRow>
								{#each data.columns as column}
									<TableCell class={isNull(row[column.name]) ? 'text-slate-400 italic' : ''}>
										<span class="font-mono text-xs">
											{formatValue(row[column.name])}
										</span>
									</TableCell>
								{/each}
							</TableRow>
						{/each}
					</TableBody>
				</Table>
			</div>
		{/if}
	</CardContent>
</Card>
