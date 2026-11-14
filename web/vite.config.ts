import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const backendPort = process.env.BACKEND_PORT || '8081';
const webPort = parseInt(process.env.WEB_PORT || '5174');

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: webPort,
		host: true,
		allowedHosts: true,
		proxy: {
			'/api': {
				target: `http://localhost:${backendPort}`,
				changeOrigin: true
			},
			'/ws': {
				target: `ws://localhost:${backendPort}`,
				ws: true
			}
		}
	}
});
