import { createAuthClient } from 'better-auth/svelte';
import { adminClient } from 'better-auth/client/plugins';

export const authClient = createAuthClient({
  baseURL: 'https://codex.singularity.co.ke',
  credentials: 'include',
  plugin:[
    adminClient()
    
  ]

});

export const { signIn, signUp, useSession } = authClient;
