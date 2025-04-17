<script lang="ts">
    import { authClient, useSession } from '$lib/auth-client.js';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { Sparkles, User, Mail, Key } from 'lucide-svelte';

    const session = useSession();

    let email = '';
    let password = '';
    let username = '';
    let error = '';
    let loading = false;

    onMount(() => {
        if ($session.data?.user) {
            goto('/dashboard');
        }
    });

    async function handleSignUp() {
        try {
            loading = true;
            error = '';

            const { data, error: signUpError } = await authClient.signUp.email({
                name: username,	
                email,
                password,
                callbackURL: "/"
            }, {
                onRequest: () => {
                    loading = true;
                },
                onError: (ctx) => {
                    error = ctx.error.message;
                }
            });

            if (signUpError) {
                throw signUpError;
            }
            
            goto('/');
        } catch (err) {
            console.log(err);
        } finally {
            loading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-100 to-amber-100">
    <div class="max-w-md w-full space-y-8 p-8 bg-gradient-to-br from-amber-50 to-purple-100 rounded-lg border border-amber-200 shadow-lg">
        <div>
            <Sparkles class="h-12 w-12 mx-auto mb-4 text-purple-700" />
            <h2 class="mt-2 text-center text-3xl font-serif font-bold text-purple-900">
                Join the Magical Order
            </h2>
            <p class="mt-2 text-center text-amber-800">
                Already a member of our mystical circle? 
                <a href="/signin" class="font-medium text-amber-700 hover:text-amber-800 underline">
                    Sign in with your runes
                </a>
            </p>
        </div>
        <form class="mt-8 space-y-6" on:submit|preventDefault={handleSignUp}>
            <div class="space-y-4">
                <div>
                    <label for="username" class="block text-amber-800 font-medium mb-2">Your Magical Title</label>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <User class="h-5 w-5 text-amber-700" />
                        </div>
                        <input
                            id="username"
                            type="text"
                            bind:value={username}
                            placeholder="Wizard Eldritch"
                            required
                            class="pl-10 w-full bg-amber-50 border border-amber-300 text-amber-900 rounded-lg p-3 focus:ring-2 focus:ring-purple-600 focus:border-transparent"
                        />
                    </div>
                </div>
                <div>
                    <label for="email" class="block text-amber-800 font-medium mb-2">Your Mystical Address</label>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <Mail class="h-5 w-5 text-amber-700" />
                        </div>
                        <input
                            id="email"
                            type="email"
                            bind:value={email}
                            placeholder="wizard@eldoria.realm"
                            required
                            class="pl-10 w-full bg-amber-50 border border-amber-300 text-amber-900 rounded-lg p-3 focus:ring-2 focus:ring-purple-600 focus:border-transparent"
                        />
                    </div>
                </div>
                <div>
                    <label for="password" class="block text-amber-800 font-medium mb-2">Create Your Arcane Seal</label>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <Key class="h-5 w-5 text-amber-700" />
                        </div>
                        <input
                            id="password"
                            type="password"
                            bind:value={password}
                            placeholder="At least 8 mystical characters"
                            required
                            minlength="8"
                            class="pl-10 w-full bg-amber-50 border border-amber-300 text-amber-900 rounded-lg p-3 focus:ring-2 focus:ring-purple-600 focus:border-transparent"
                        />
                    </div>
                    <p class="mt-1 text-xs text-amber-700">Your seal must contain at least 8 mystical symbols for proper enchantment</p>
                </div>
            </div>

            <div>
                <button
                    type="submit"
                    disabled={loading}
                    class="w-full bg-purple-800 hover:bg-purple-900 text-amber-100 font-medium px-6 py-3 rounded-lg transition-colors border border-purple-700 flex items-center justify-center space-x-2 disabled:opacity-50"
                >
                    <Sparkles size={18} />
                    <span>{loading ? 'Casting Spell...' : 'Create Magical Account'}</span>
                </button>
            </div>

            {#if error}
                <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded" role="alert">
                    <span class="font-bold">Arcane Error:</span> {error}
                </div>
            {/if}
        </form>
    </div>
</div>