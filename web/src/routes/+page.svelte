<script lang="ts">
	/**
	 * Main Page Component
	 * Routes between landing page and connected editor based on secret parameter
	 */

	import { onMount } from 'svelte';
	import LandingPage from '$lib/components/LandingPage.svelte';
	import ConnectionStatus from '$lib/components/ConnectionStatus.svelte';
	import QueryPanel from '$lib/components/QueryPanel.svelte';
	import { connect, disconnect, clearError } from '$lib/stores/connection';
	import { connectionStore } from '$lib/stores/connection';
	import { setSchema, clearSchema } from '$lib/stores/schema';
	import { ConnectionState } from '$lib/services/websocket';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui';

	// Connection state
	let secret = $state<string | null>(null);
	let isInitializing = $state(true);
	let schemaLoaded = $state(false);

	/**
	 * Initialize connection when component mounts
	 */
	onMount(() => {
		// Parse secret from URL query parameter
		const urlParams = new URLSearchParams(window.location.search);
		const secretParam = urlParams.get('secret');

		if (secretParam) {
			secret = secretParam;
			// Start connection process
			connectToProxy(secretParam);
		} else {
			// No secret, show landing page
			isInitializing = false;
		}

		// Cleanup on unmount
		return () => {
			disconnect();
		};
	});

	/**
	 * Connects to the proxy server and introspects schema
	 */
	async function connectToProxy(secretValue: string) {
		try {
			clearError();
			schemaLoaded = false;

			// Connect to proxy
			await connect(secretValue);

			// Connection successful, introspect schema
			const client = $connectionStore.client;
			if (client) {
				try {
					const schema = await client.introspectSchema();
					setSchema(schema);
					schemaLoaded = true;
				} catch (schemaError) {
					console.error('Failed to introspect schema:', schemaError);
					// Continue even if schema introspection fails
					// User can still use the editor without autocomplete
					schemaLoaded = true;
				}
			}
		} catch (error) {
			console.error('Connection failed:', error);
			// Error is already set by the connection store
		} finally {
			isInitializing = false;
		}
	}

	/**
	 * Retry connection on error
	 */
	function handleRetry() {
		if (secret) {
			isInitializing = true;
			clearSchema();
			connectToProxy(secret);
		}
	}

	// Reactive check for connected state
	$effect(() => {
		const state = $connectionStore.state;
		const isConnected = state === ConnectionState.CONNECTED;

		// If we lost connection, reset schema
		if (!isConnected && schemaLoaded) {
			schemaLoaded = false;
			clearSchema();
		}
	});
</script>

<!-- Main Page Layout -->
<div class="min-h-screen bg-background">
	{#if !secret}
		<!-- No secret in URL: Show landing page -->
		<LandingPage />
	{:else if isInitializing}
		<!-- Connecting to proxy -->
		<div class="container mx-auto px-4 py-16">
			<Card class="max-w-2xl mx-auto">
				<CardHeader>
					<CardTitle>Connecting to Proxy...</CardTitle>
					<CardDescription>Establishing connection to local proxy server</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="flex items-center justify-center py-8">
						<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
					</div>
					<p class="text-center text-muted-foreground text-sm">This may take a few seconds...</p>
				</CardContent>
			</Card>
		</div>
	{:else if $connectionStore.state === ConnectionState.ERROR}
		<!-- Connection error -->
		<div class="container mx-auto px-4 py-16">
			<Card class="max-w-2xl mx-auto">
				<CardHeader>
					<CardTitle class="text-destructive">Connection Failed</CardTitle>
					<CardDescription>Unable to connect to the local proxy server</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="space-y-4">
						<ConnectionStatus onRetry={handleRetry} />

						<div class="bg-muted p-4 rounded-lg text-sm">
							<p class="font-semibold mb-2">Troubleshooting:</p>
							<ul class="list-disc list-inside space-y-1 text-muted-foreground">
								<li>Make sure the proxy server is running</li>
								<li>Check that the secret in the URL matches the one from the proxy</li>
								<li>Verify the proxy is running on localhost:8080</li>
								<li>Check your browser console for detailed errors</li>
							</ul>
						</div>
					</div>
				</CardContent>
			</Card>
		</div>
	{:else if $connectionStore.state === ConnectionState.CONNECTED && schemaLoaded}
		<!-- Successfully connected -->
		<div class="container mx-auto px-4 py-8">
			<!-- Header with connection status -->
			<div class="mb-6">
				<div class="flex items-center justify-between">
					<h1 class="text-3xl font-bold">PostgreSQL Client</h1>
					<ConnectionStatus onRetry={handleRetry} />
				</div>
			</div>

			<!-- Query Panel with Editor and Results -->
			<QueryPanel />
		</div>
	{/if}
</div>
