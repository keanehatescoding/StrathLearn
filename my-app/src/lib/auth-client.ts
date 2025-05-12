import { createAuthClient } from 'better-auth/svelte';

export const authClient = createAuthClient({
  baseURL: 'https://codex.singularity.co.ke/',
  credentials: 'include'
});

export const { signIn, signUp, useSession } = authClient;
