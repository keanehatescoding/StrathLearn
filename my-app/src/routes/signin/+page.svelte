<script lang="ts">
    import { authClient, useSession } from '$lib/auth-client.js';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { toast } from 'svelte-sonner';
    import { Mail, Key, ArrowLeft, Code } from 'lucide-svelte';
    
    // Import UI components
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Checkbox } from "$lib/components/ui/checkbox";
	import { siGithub } from 'simple-icons';
    
    const session = useSession();

    let email = '';
    let password = '';
    let rememberMe = false;
    let error = '';
    let loading = false;

    onMount(() => {
        if ($session.data?.user) {
            goto('/dashboard');
        }
    });

    async function handleSignIn() {
        try {
            loading = true;
            error = '';
            
            const toastId = toast.loading('Signing in...');

            const { data, error: signInError } = await authClient.signIn.email({
                email,
                password,
                rememberMe
            }, {
                onRequest: () => {
                    loading = true;
                },
                onError: (ctx) => {
                    error = ctx.error.message;
                    toast.dismiss(toastId);
                    toast.error('Sign in failed', {
                        description: ctx.error.message
                    });
                }
            });

            if (signInError) {
                throw signInError;
            }
            
            toast.dismiss(toastId);
            toast.success('Signed in successfully', {
                description: 'Welcome back to Codex!'
            });
            
            goto('/dashboard');
        } catch (err) {
            console.error(err);
        } finally {
            loading = false;
        }
    }
    
    async function handleGitHubSignIn() {
        try {
            loading = true;
            error = '';
            
            const toastId = toast.loading('Connecting to GitHub...');
            
         
            await authClient.signIn.social({
    /**
     * The social provider id
     * @example "github", "google", "apple"
     */
    provider: "github",
    /**
     * A URL to redirect after the user authenticates with the provider
     * @default "/"
     */
    callbackURL: "/dashboard", 
    /**
     * A URL to redirect if an error occurs during the sign in process
     */
    errorCallbackURL: "/error",
    /**
     * A URL to redirect if the user is newly registered
     */
    newUserCallbackURL: "/welcome",
    /**

     * @default false
     */
    // disableRedirect: true,
});
            toast.dismiss(toastId);
        } catch (err) {
            console.error(err);
            toast.error('GitHub authentication failed');
            loading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-background px-4">
    <Card class="w-full max-w-md border-border shadow-lg">
        <CardHeader class="space-y-1">
            <div class="flex items-center justify-between">
                <CardTitle class="text-2xl font-bold flex items-center">
                    <Code class="h-5 w-5 mr-2 text-primary" />
                    Sign in to Codex
                </CardTitle>
                <a href="/" class="text-muted-foreground hover:text-foreground text-sm flex items-center transition-colors">
                    <ArrowLeft class="h-4 w-4 mr-1" />
                    Back to home
                </a>
            </div>
            <CardDescription>
                Enter your credentials to access your account
            </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
            <div class="space-y-4">
                <Button 
                    variant="outline" 
                    class="w-full flex items-center justify-center gap-2 h-10 hover:bg-accent/20" 
                    onclick={handleGitHubSignIn}
                    disabled={loading}
                >
                <svg
                role="img"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="currentColor"
              >
                <path d={siGithub.path} />
              </svg>
                    <span>Sign in with GitHub</span>
                </Button>
                
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <span class="w-full border-t border-border"></span>
                    </div>
                    <div class="relative flex justify-center text-xs uppercase">
                        <span class="bg-background px-2 text-muted-foreground">
                            Or continue with
                        </span>
                    </div>
                </div>
            </div>
            
            <form class="space-y-4" on:submit|preventDefault={handleSignIn}>
                <div class="space-y-2">
                    <Label for="email">Email</Label>
                    <div class="relative">
                        <Mail class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                        <Input 
                            id="email" 
                            type="email" 
                            placeholder="your@email.com" 
                            class="pl-10"
                            bind:value={email}
                            required
                        />
                    </div>
                </div>
                <div class="space-y-2">
                    <div class="flex items-center justify-between">
                        <Label for="password">Password</Label>
                        <a href="/forgot-password" class="text-xs text-primary hover:underline">
                            Forgot password?
                        </a>
                    </div>
                    <div class="relative">
                        <Key class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                        <Input 
                            id="password" 
                            type="password" 
                            placeholder="••••••••" 
                            class="pl-10"
                            bind:value={password}
                            required
                        />
                    </div>
                </div>
                
                <div class="flex items-center space-x-2">
                    <Checkbox id="remember" bind:checked={rememberMe} />
                    <p  class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                        Remember me
                    </p>
                </div>
                
                <Button 
                    type="submit" 
                    class="w-full bg-primary hover:bg-primary/90 text-primary-foreground" 
                    disabled={loading}
                >
                    {#if loading}
                        <div class="animate-spin mr-2 h-4 w-4 border-2 border-current border-t-transparent rounded-full"></div>
                    {/if}
                    Sign in
                </Button>
            </form>

            {#if error}
                <div class="bg-destructive/10 text-destructive p-3 rounded-md text-sm" role="alert">
                    {error}
                </div>
            {/if}
        </CardContent>
        <CardFooter>
            <p class="text-sm text-center text-muted-foreground w-full">
                Don't have an account?{" "}
                <a href="/signup" class="text-primary hover:underline">
                    Create one
                </a>
            </p>
        </CardFooter>
    </Card>
</div>

<style>
    :global(body) {
        font-family: 'Inter', sans-serif;
    }
</style>