<script lang="ts">
	/**
	 * Connection Status Component
	 * Displays the current connection state with appropriate styling
	 */

	import { ConnectionState } from '$lib/services/websocket';
	import {
		currentState,
		errorMessage,
		isConnected,
		isConnecting,
		isReconnecting,
		hasError
	} from '$lib/stores/connection';
	import Badge from '$lib/components/ui/badge.svelte';
	import Alert from '$lib/components/ui/alert.svelte';
	import AlertDescription from '$lib/components/ui/alert-description.svelte';

	/**
	 * Optional callback for retry action
	 */
	interface ConnectionStatusProps {
		onRetry?: () => void;
	}

	let { onRetry }: ConnectionStatusProps = $props();

	/**
	 * Get badge variant based on connection state
	 */
	function getBadgeVariant(
		state: ConnectionState
	): 'default' | 'secondary' | 'destructive' | 'outline' {
		switch (state) {
			case ConnectionState.CONNECTED:
				return 'default'; // Will be styled green
			case ConnectionState.CONNECTING:
			case ConnectionState.RECONNECTING:
				return 'secondary'; // Will be styled blue/yellow
			case ConnectionState.ERROR:
				return 'destructive'; // Red
			case ConnectionState.DISCONNECTED:
			default:
				return 'outline'; // Gray
		}
	}

	/**
	 * Get status text based on connection state
	 */
	function getStatusText(state: ConnectionState): string {
		switch (state) {
			case ConnectionState.CONNECTED:
				return '● Connected';
			case ConnectionState.CONNECTING:
				return '● Connecting...';
			case ConnectionState.RECONNECTING:
				return '● Reconnecting...';
			case ConnectionState.ERROR:
				return '● Connection Error';
			case ConnectionState.DISCONNECTED:
			default:
				return '● Disconnected';
		}
	}

	/**
	 * Get custom badge class for connection state colors
	 */
	function getCustomBadgeClass(state: ConnectionState): string {
		switch (state) {
			case ConnectionState.CONNECTED:
				return 'bg-green-500 hover:bg-green-600 text-white border-green-600';
			case ConnectionState.CONNECTING:
				return 'bg-blue-500 hover:bg-blue-600 text-white border-blue-600';
			case ConnectionState.RECONNECTING:
				return 'bg-yellow-500 hover:bg-yellow-600 text-white border-yellow-600';
			case ConnectionState.ERROR:
				return ''; // Use default destructive styling
			case ConnectionState.DISCONNECTED:
			default:
				return 'bg-gray-400 text-gray-700 border-gray-500';
		}
	}
</script>

<div class="connection-status" data-testid="connection-status">
	<!-- Status Badge -->
	<div class="flex items-center gap-2">
		<Badge
			variant={getBadgeVariant($currentState)}
			class={getCustomBadgeClass($currentState)}
		>
			{getStatusText($currentState)}
		</Badge>

		<!-- Retry Button for Error/Disconnected States -->
		{#if ($currentState === ConnectionState.ERROR || $currentState === ConnectionState.DISCONNECTED) && onRetry}
			<button
				onclick={onRetry}
				class="text-sm text-blue-600 hover:text-blue-800 underline"
				data-testid="retry-button"
			>
				Retry
			</button>
		{/if}
	</div>

	<!-- Error Message Display -->
	{#if $hasError && $errorMessage}
		<div class="mt-2" data-testid="error-alert">
			<Alert variant="destructive">
				<AlertDescription>
					<strong>Error:</strong>
					{$errorMessage}
				</AlertDescription>
			</Alert>
		</div>
	{/if}
</div>

<style>
	.connection-status {
		@apply w-full;
	}
</style>
