<script lang="ts">
    import { authClient, useSession } from '$lib/auth-client.js';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { toast } from 'svelte-sonner';
    import { User, Mail, Key, ArrowLeft, Code } from 'lucide-svelte';
    import { cn } from "$lib/utils";
    
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { siGithub, siGoogle, siC, siRust, siGo, siTypescript } from 'simple-icons';
    
    const session = useSession();

    let email = '';
    let password = '';
    let username = '';
    let confirmPassword = '';
    let agreeToTerms = false;
    let error = '';
    let loading = false;

    onMount(() => {
        if ($session.data?.user) {
            goto('/dashboard');
        }
    });

    async function handleSignUp() {
        if (password !== confirmPassword) {
            error = 'Passwords do not match';
            return;
        }

        if (!agreeToTerms) {
            error = 'Please agree to the terms of service';
            return;
        }

        try {
            loading = true;
            error = '';
            
            const toastId = toast.loading('Creating your account...');

            const { data, error: signUpError } = await authClient.signUp.email({
                name: username,
                email,
                password,
                callbackURL: "/dashboard"
            }, {
                onRequest: () => {
                    loading = true;
                },
                onError: (ctx) => {
                    error = ctx.error.message;
                    toast.dismiss(toastId);
                    toast.error('Sign up failed', {
                        description: ctx.error.message
                    });
                }
            });

            if (signUpError) {
                throw signUpError;
            }
            
            toast.dismiss(toastId);
            toast.success('Account created successfully', {
                description: 'Welcome to Codepass!'
            });
            
            goto('/dashboard');
        } catch (err) {
            console.error(err);
        } finally {
            loading = false;
        }
    }
    
    async function handleGitHubSignUp() {
        try {
            loading = true;
            error = '';
            
            const toastId = toast.loading('Connecting to GitHub...');
            
            await authClient.signIn.social({
                provider: "github",
                callbackURL: "/dashboard",
                errorCallbackURL: "/error",
                newUserCallbackURL: "/welcome",
            });
            
            toast.dismiss(toastId);
        } catch (err) {
            console.error(err);
            toast.error('GitHub authentication failed');
            loading = false;
        }
    }

    async function handleGoogleSignUp() {
        try {
            loading = true;
            error = '';
            
            const toastId = toast.loading('Connecting to Google...');
            
            await authClient.signIn.social({
                provider: "google",
                callbackURL: "/dashboard",
                errorCallbackURL: "/error",
                newUserCallbackURL: "/welcome",
            });
            
            toast.dismiss(toastId);
        } catch (err) {
            console.error(err);
            toast.error('Google authentication failed');
            loading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center bg-background text-foreground">
    <div class="flex-1 flex items-center justify-center px-8 py-12">
        <div class="w-full max-w-md space-y-8">
            <div class="text-center">
 <div class="flex items-center justify-center">
    <img src="/logo.png" alt="Logo" class="h-32" />
    <div class="text-3xl font-bold text-primary -ml-2">Codepass</div>
</div>
                <h1 class="text-4xl font-bold mb-4">Join Codepass</h1>
                <p class="text-muted-foreground">Create your account and start your coding journey</p>
            </div>

            <div class="grid grid-cols-2 gap-3">
                <Button 
                    variant="outline" 
                    class="bg-muted border-muted-foreground/20 text-foreground hover:bg-muted/80" 
                    onclick={handleGoogleSignUp}
                    disabled={loading}
                >
                    <svg
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5 mr-2"
                        fill="currentColor"
                    >
                        <path d={siGoogle.path} />
                    </svg>
                    Google
                </Button>
                
                <Button 
                    variant="outline" 
                    class="bg-muted border-muted-foreground/20 text-foreground hover:bg-muted/80" 
                    onclick={handleGitHubSignUp}
                    disabled={loading}
                >
                    <svg
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5 mr-2"
                        fill="currentColor"
                    >
                        <path d={siGithub.path} />
                    </svg>
                    GitHub
                </Button>
            </div>

            <div class="relative">
                <div class="absolute inset-0 flex items-center">
                    <span class="w-full border-t border-muted-foreground/20"></span>
                </div>
                <div class="relative flex justify-center text-xs uppercase">
                    <span class="bg-background px-2 text-muted-foreground">
                        Or continue with email
                    </span>
                </div>
            </div>

            <form class="space-y-6" on:submit|preventDefault={handleSignUp}>
                <div class="space-y-2">
                    <Label for="username" class="text-foreground">Full name</Label>
                    <div class="relative">
                        <User class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                        <Input 
                            id="username" 
                            type="text" 
                            placeholder="John Doe" 
                            class="pl-10 bg-muted border-muted-foreground/20 text-foreground placeholder:text-muted-foreground focus:border-primary focus:ring-primary"
                            bind:value={username}
                            required
                        />
                    </div>
                </div>
                
                <div class="space-y-2">
                    <Label for="email" class="text-foreground">Email address</Label>
                    <div class="relative">
                        <Mail class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                        <Input 
                            id="email" 
                            type="email" 
                            placeholder="your@email.com" 
                            class="pl-10 bg-muted border-muted-foreground/20 text-foreground placeholder:text-muted-foreground focus:border-primary focus:ring-primary"
                            bind:value={email}
                            required
                        />
                    </div>
                </div>
                
                <div class="space-y-2">
                    <Label for="password" class="text-foreground">Password</Label>
                    <div class="relative">
                        <Key class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                        <Input 
                            id="password" 
                            type="password" 
                            placeholder="••••••••" 
                            class="pl-10 bg-muted border-muted-foreground/20 text-foreground placeholder:text-muted-foreground focus:border-primary focus:ring-primary"
                            bind:value={password}
                            required
                            minlength="8"
                        />
                    </div>
                </div>

                <div class="space-y-2">
                    <Label for="confirmPassword" class="text-foreground">Confirm password</Label>
                    <div class="relative">
                        <Key class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                        <Input 
                            id="confirmPassword" 
                            type="password" 
                            placeholder="••••••••" 
                            class="pl-10 bg-muted border-muted-foreground/20 text-foreground placeholder:text-muted-foreground focus:border-primary focus:ring-primary"
                            bind:value={confirmPassword}
                            required
                        />
                    </div>
                    <p class="text-xs text-muted-foreground">
                        Password must be at least 8 characters long
                    </p>
                </div>
                
                <div class="flex items-center space-x-2">
                    <Checkbox 
                        id="terms" 
                        bind:checked={agreeToTerms} 
                        class="border-muted-foreground data-[state=checked]:bg-primary data-[state=checked]:border-primary" 
                    />
                    <Label for="terms" class="text-sm text-foreground">
                        I agree to the{" "}
                        <a href="/terms" class="text-primary hover:text-primary/80 underline">
                            Terms of Service
                        </a>
                        {" "}and{" "}
                        <a href="/privacy" class="text-primary hover:text-primary/80 underline">
                            Privacy Policy
                        </a>
                    </Label>
                </div>
                
                <Button 
                    type="submit" 
                    class="w-full bg-primary text-primary-foreground hover:bg-primary/90 py-3 font-medium" 
                    disabled={loading || !agreeToTerms}
                >
                    {#if loading}
                        <div class="animate-spin mr-2 h-4 w-4 border-2 border-current border-t-transparent rounded-full"></div>
                    {/if}
                    Create Account
                    <ArrowLeft class="ml-2 h-4 w-4 rotate-180" />
                </Button>
            </form>

            {#if error}
                <div class="bg-destructive/10 border border-destructive/30 text-destructive p-3 rounded-md text-sm" role="alert">
                    {error}
                </div>
            {/if}

            <p class="text-center text-sm text-muted-foreground">
                Already have an account?{" "}
                <a href="/login" class="text-primary hover:text-primary/80 underline">
                    Sign in
                </a>
            </p>
        </div>
    </div>

    <div class="flex-1 bg-background flex items-center justify-center relative overflow-hidden">
        <div class="relative flex h-[500px] w-[500px] items-center justify-center overflow-hidden">
            <span class="pointer-events-none whitespace-pre-wrap bg-gradient-to-b from-primary/90 to-primary/50 bg-clip-text text-center text-6xl font-bold leading-none text-transparent">
                Codepass
            </span>

            
            <svg
                xmlns="http://www.w3.org/2000/svg"
                version="1.1"
                class="pointer-events-none absolute inset-0 h-full w-full"
            >
                <circle
                    class="stroke-primary/30 stroke-[1.5]"
                    cx="50%"
                    cy="50%"
                    r="120"
                    fill="none"
                    stroke-dasharray="8 4"
                />
                <circle
                    class="stroke-primary/30 stroke-[1.5]"
                    cx="50%"
                    cy="50%"
                    r="200"
                    fill="none"
                    stroke-dasharray="8 4"
                />
            </svg>

            <!-- First inner orbit (C) -->
            <div 
                class="absolute flex h-full w-full transform-gpu animate-orbit items-center justify-center rounded-full" 
                style="--delay: 0s; --duration: 20s; --radius: 120px;"
            >
                <div class="h-[40px] w-[40px] rounded-full bg-transparent border border-primary/30 flex items-center justify-center" style="transform: translateY(120px);">
                    <svg
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6 text-primary/80"
                        fill="currentColor"
                    >
                        <path d={siC.path} />
                    </svg>
                </div>
            </div>
            
            <!-- Second inner orbit (Rust) -->
            <div 
                class="absolute flex h-full w-full transform-gpu animate-orbit items-center justify-center rounded-full" 
                style="animation-delay: -10s; --duration: 20s; --radius: 120px;"
            >
                <div class="h-[40px] w-[40px] rounded-full bg-primary/10 backdrop-blur-sm border border-primary/30 flex items-center justify-center" style="transform: translateY(120px) rotate(180deg);">
                    <svg
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6 text-primary/80"
                        fill="currentColor"
                    >
                        <path d={siRust.path} />
                    </svg>
                </div>
            </div>

            <!-- First outer orbit (Go) -->
            <div 
                class="absolute flex h-full w-full transform-gpu animate-orbit-reverse items-center justify-center rounded-full" 
                style="--delay: 0s; --duration: 25s; --radius: 200px;"
            >
                <div class="h-[50px] w-[50px] rounded-full bg-primary/10 backdrop-blur-sm border border-primary/30 flex items-center justify-center" style="transform: translateY(200px);">
                    <svg
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-7 w-7 text-primary/80"
                        fill="currentColor"
                    >
                        <path d={siGo.path} />
                    </svg>
                </div>
            </div>
            
            <!-- Second outer orbit (TypeScript) -->
            <div 
                class="absolute flex h-full w-full transform-gpu animate-orbit-reverse items-center justify-center rounded-full" 
                style="animation-delay: -12s; --duration: 25s; --radius: 200px;"
            >
                <div class="h-[50px] w-[50px] rounded-full bg-transparent border border-primary/30 flex items-center justify-center" style="transform: translateY(200px) rotate(180deg);">
                    <svg
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-7 w-7 text-primary/80"
                        fill="currentColor"
                    >
                        <path d={siTypescript.path} />
                    </svg>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    :global(body) {
        font-family: 'Inter', sans-serif;
        margin: 0;
        padding: 0;
    }

    @keyframes orbit {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }

    @keyframes orbit-reverse {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(-360deg);
        }
    }

    .animate-orbit {
        animation: orbit var(--duration, 20s) linear infinite;
        animation-delay: var(--delay, 0s);
    }

    .animate-orbit-reverse {
        animation: orbit-reverse var(--duration, 20s) linear infinite;
        animation-delay: var(--delay, 0s);
    }
</style>