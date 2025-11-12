<script lang="ts">
	import { onMount } from 'svelte';
	import Editor from './Editor.svelte';
	import ErrorDisplay from './ErrorDisplay.svelte';
	import Results from './Results.svelte';
	import {
		Button,
		Card,
		CardContent,
		CardHeader,
		CardTitle,
		Separator
	} from '$lib/components/ui';
	import { connectionStore } from '$lib/stores/connection';
	import { schemaStore } from '$lib/stores/schema';
	import { Play, Loader2 } from 'lucide-svelte';
	import type { ResultPayload } from '$lib/services/protocol';
	import type { SchemaInfo } from '$lib/utils/autocomplete';

	// SQL query state
	let sql = '';
	let loading = false;
	let error: string | null = null;
	let results: ResultPayload | null = null;

	// Editor reference for programmatic control
	let editor: Editor;

	// Get client from connection store
	$: client = $connectionStore.client;

	// Get schema from schema store and map to autocomplete SchemaInfo format
	$: schema = mapSchemaToAutocomplete($schemaStore);

	/**
	 * Maps protocol SchemaPayload to autocomplete SchemaInfo
	 */
	function mapSchemaToAutocomplete(schemaPayload: typeof $schemaStore): SchemaInfo {
		return {
			tables: schemaPayload.tables.map((table) => ({
				schema: table.schema,
				name: table.name,
				columns: table.columns.map((col) => ({
					name: col.name,
					type: col.dataType,
					nullable: col.nullable ?? false
				}))
			})),
			functions: schemaPayload.functions.map((func) => ({
				schema: func.schema,
				name: func.name,
				returnType: func.returnType
			}))
		};
	}

	/**
	 * Executes the SQL query
	 */
	async function executeQuery() {
		// Validate SQL is not empty
		if (!sql.trim()) {
			error = 'Please enter a SQL query';
			return;
		}

		// Validate client is available
		if (!client) {
			error = 'Not connected to proxy server';
			return;
		}

		// Set loading state
		loading = true;
		error = null;
		results = null;

		try {
			// Execute query via client
			const result = await client.executeQuery(sql);

			// Store results
			results = result;
			error = null;
		} catch (err) {
			// Handle error
			const errorMsg = err instanceof Error ? err.message : String(err);
			error = errorMsg;
			results = null;
		} finally {
			// Clear loading state
			loading = false;
		}
	}

	/**
	 * Handles SQL changes from the editor
	 */
	function handleSqlChange(newSql: string) {
		sql = newSql;
	}

	/**
	 * Handles Ctrl+Enter from the editor
	 */
	function handleExecute() {
		executeQuery();
	}

	// Set initial query example on mount
	onMount(() => {
		sql =
			'-- Write your SQL query here\n-- Press Ctrl+Enter or click Run to execute\n\nSELECT 1 as example;';
	});
</script>

<div class="query-panel flex flex-col h-full gap-4">
	<!-- Query Editor Section -->
	<Card class="flex-shrink-0">
		<CardHeader class="pb-3">
			<div class="flex items-center justify-between">
				<CardTitle class="text-lg">SQL Query Editor</CardTitle>
				<Button
					onclick={executeQuery}
					disabled={loading || !client}
					class="gap-2"
					variant="default"
				>
					{#if loading}
						<Loader2 class="h-4 w-4 animate-spin" />
						Executing...
					{:else}
						<Play class="h-4 w-4" />
						Run Query
					{/if}
					<span class="text-xs opacity-70 ml-1">(Ctrl+Enter)</span>
				</Button>
			</div>
		</CardHeader>
		<CardContent>
			<Editor
				bind:this={editor}
				bind:value={sql}
				onChange={handleSqlChange}
				onExecute={handleExecute}
				height="300px"
				{schema}
			/>
		</CardContent>
	</Card>

	<Separator />

	<!-- Results Section -->
	<div class="results-section flex-1 overflow-auto">
		{#if loading}
			<!-- Loading State -->
			<Card>
				<CardContent class="py-12">
					<div class="flex flex-col items-center justify-center gap-4 text-muted-foreground">
						<Loader2 class="h-8 w-8 animate-spin" />
						<p>Executing query...</p>
					</div>
				</CardContent>
			</Card>
		{:else if error}
			<!-- Error Display -->
			<ErrorDisplay {error} />
		{:else if results}
			<!-- Results Display -->
			<Results data={results} />
		{:else}
			<!-- Empty State -->
			<Card>
				<CardContent class="py-12">
					<div class="flex flex-col items-center justify-center gap-2 text-muted-foreground">
						<Play class="h-8 w-8 opacity-50" />
						<p>Execute a query to see results</p>
						<p class="text-xs">Type your SQL above and press Ctrl+Enter or click Run</p>
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>

<style>
	.query-panel {
		width: 100%;
		max-width: 100%;
	}

	.results-section {
		min-height: 200px;
	}
</style>
