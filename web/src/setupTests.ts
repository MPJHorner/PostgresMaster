/**
 * Test setup file for vitest
 */

import '@testing-library/jest-dom';

// Set up DOM environment for Svelte 5
if (typeof window !== 'undefined') {
	// Ensure we're in browser mode
	globalThis.window = window;
	globalThis.document = document;
}
