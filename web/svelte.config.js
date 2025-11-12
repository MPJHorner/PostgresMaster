import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// Configure static adapter for pure static site generation
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: 'index.html'
		}),
		// Handle missing favicon during prerender
		prerender: {
			handleHttpError: ({ path, referrer, message }) => {
				// Ignore 404 for favicon during prerender
				if (path === '/favicon.png') {
					return;
				}
				// Throw for other errors
				throw new Error(message);
			}
		}
	}
};

export default config;
