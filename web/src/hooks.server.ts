import { auth } from "$lib/auth.js"; 
import { svelteKitHandler } from "better-auth/svelte-kit";
 
export async function handle({ event, resolve }) {
    return svelteKitHandler({ event, resolve, auth });
}