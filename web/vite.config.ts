import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';
// @ts-expect-error - vite-plugin-monaco-editor types are not properly exported
import monacoEditorPluginModule from 'vite-plugin-monaco-editor';

// Handle both CommonJS and ES module exports
// @ts-expect-error - accessing default property
const monacoEditorPlugin = monacoEditorPluginModule.default || monacoEditorPluginModule;

export default defineConfig({
	plugins: [
		sveltekit(),
		monacoEditorPlugin({
			languageWorkers: ['editorWorkerService']
		})
	],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./src/setupTests.ts']
	}
});
