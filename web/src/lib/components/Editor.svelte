<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import type * as Monaco from 'monaco-editor';
	import { setupAutocomplete, type SchemaInfo } from '$lib/utils/autocomplete';

	// Props
	export let value = '';
	export let onChange: ((newValue: string) => void) | undefined = undefined;
	export let onExecute: (() => void) | undefined = undefined;
	export let height = '400px';
	export let schema: SchemaInfo | undefined = undefined;

	let editorContainer: HTMLDivElement;
	let editor: Monaco.editor.IStandaloneCodeEditor | null = null;
	let monaco: typeof Monaco | null = null;
	let autocompleteDisposable: Monaco.IDisposable | null = null;

	onMount(async () => {
		// Only run in browser
		if (!browser) {
			return;
		}

		// Dynamically import Monaco Editor
		monaco = await import('monaco-editor');
		if (!editorContainer) {
			return;
		}

		// Initialize Monaco editor
		editor = monaco.editor.create(editorContainer, {
			value: value,
			language: 'sql',
			theme: 'vs-dark',
			automaticLayout: true,
			minimap: {
				enabled: false
			},
			wordWrap: 'on',
			scrollBeyondLastLine: false,
			fontSize: 14,
			lineNumbers: 'on',
			roundedSelection: true,
			padding: {
				top: 10,
				bottom: 10
			},
			suggest: {
				showKeywords: true,
				showSnippets: true
			}
		});

		// Listen for content changes
		editor.onDidChangeModelContent(() => {
			if (editor) {
				const newValue = editor.getValue();
				value = newValue;
				if (onChange) {
					onChange(newValue);
				}
			}
		});

		// Register Ctrl+Enter (or Cmd+Enter on Mac) command for execution
		if (monaco) {
			editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
				if (onExecute) {
					onExecute();
				}
			});

			// Setup SQL autocomplete with keywords, functions, and schema-aware completions
			autocompleteDisposable = setupAutocomplete(monaco, schema);
		}

		// Set initial value if provided
		if (value && editor) {
			editor.setValue(value);
		}
	});

	onDestroy(() => {
		// Dispose autocomplete provider
		if (autocompleteDisposable) {
			autocompleteDisposable.dispose();
			autocompleteDisposable = null;
		}

		// Dispose editor to free resources
		if (editor) {
			editor.dispose();
			editor = null;
		}
	});

	// Watch for external value changes
	$: if (editor && value !== editor.getValue()) {
		const position = editor.getPosition();
		editor.setValue(value);
		if (position) {
			editor.setPosition(position);
		}
	}

	// Watch for schema changes and update autocomplete
	$: if (monaco && schema) {
		// Dispose existing autocomplete provider
		if (autocompleteDisposable) {
			autocompleteDisposable.dispose();
		}
		// Re-register with new schema
		autocompleteDisposable = setupAutocomplete(monaco, schema);
	}

	// Public method to get current value
	export function getValue(): string {
		return editor?.getValue() ?? value;
	}

	// Public method to set value programmatically
	export function setValue(newValue: string): void {
		value = newValue;
		if (editor) {
			editor.setValue(newValue);
		}
	}

	// Public method to focus the editor
	export function focus(): void {
		editor?.focus();
	}
</script>

<div
	class="editor-wrapper rounded-md border border-gray-700 overflow-hidden"
	style="height: {height};"
>
	<div bind:this={editorContainer} class="editor-container w-full h-full"></div>
</div>

<style>
	.editor-wrapper {
		background: #1e1e1e;
	}

	.editor-container {
		width: 100%;
		height: 100%;
	}
</style>
