import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			fallback: 'index.html',
		}),
		// override files to look in the frontend dir
		files: {
			lib: 'frontend/lib',
			params: 'frontend/params',
			routes: 'frontend/routes',
			serviceWorker: 'frontend/service-worker',
			appTemplate: 'frontend/app.html',
			errorTemplate: 'frontend/error.html',
		},
	},
};

export default config;
