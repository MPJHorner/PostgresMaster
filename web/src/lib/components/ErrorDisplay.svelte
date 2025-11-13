<script lang="ts">
	import { Alert, AlertDescription, AlertTitle } from '$lib/components/ui';
	import { AlertCircle } from 'lucide-svelte';
	import { parseError, getErrorCodeDescription } from '$lib/utils/errorParser';

	/**
	 * Error message to display
	 */
	export let error: string;

	$: parsedError = parseError(error);
	$: errorCodeDescription = getErrorCodeDescription(parsedError.code);
</script>

<Alert variant="destructive" class="error-display">
	<AlertCircle class="h-4 w-4" />
	<AlertTitle class="flex items-center gap-2">
		<span>Query Error</span>
		{#if parsedError.code}
			<span class="text-xs font-mono bg-destructive/20 px-2 py-0.5 rounded">
				{parsedError.code}
			</span>
		{/if}
	</AlertTitle>
	<AlertDescription class="mt-2 space-y-2">
		<!-- Main error message -->
		<div class="font-mono text-sm whitespace-pre-wrap">
			{parsedError.message}
		</div>

		<!-- Error code description -->
		{#if errorCodeDescription}
			<div class="text-xs opacity-90">
				<span class="font-semibold">Type:</span>
				{errorCodeDescription}
			</div>
		{/if}

		<!-- Position indicator -->
		{#if parsedError.position}
			<div class="text-xs opacity-90">
				<span class="font-semibold">Position:</span>
				Character {parsedError.position}
			</div>
		{/if}

		<!-- Detail section -->
		{#if parsedError.detail}
			<div class="text-xs border-l-2 border-destructive/50 pl-2 py-1">
				<div class="font-semibold opacity-90">Detail:</div>
				<div class="font-mono opacity-90">{parsedError.detail}</div>
			</div>
		{/if}

		<!-- Hint section -->
		{#if parsedError.hint}
			<div class="text-xs border-l-2 border-blue-500/50 pl-2 py-1">
				<div class="font-semibold opacity-90">Hint:</div>
				<div class="font-mono opacity-90">{parsedError.hint}</div>
			</div>
		{/if}
	</AlertDescription>
</Alert>

<style>
	.error-display :global(.font-mono) {
		font-family: 'Fira Code', 'Consolas', 'Monaco', monospace;
	}
</style>
