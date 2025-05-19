import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    // Allow only specific hosts
    allowedHosts: [
      'localhost',
      '127.0.0.1',
      'ed1e-196-216-95-131.ngrok-free.app', // Your specific ngrok URL
      '.ngrok-free.app', // Allow all ngrok-free.app subdomains
    ],
    cors: {
      origin: [
        'http://localhost:3000',
        'http://localhost:5173',
        'https://api.singularity.co.ke',
        'https://ed1e-196-216-95-131.ngrok-free.app'
      ],
      credentials: true,
      methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
      allowedHeaders: ['Content-Type', 'Authorization']
    }
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}']
  }
});