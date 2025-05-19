<script lang="ts">
   import { authClient, useSession } from '$lib/auth-client.js';
   import { goto } from '$app/navigation';
   import { onMount } from 'svelte';
   import { toast } from 'svelte-sonner';
   import { Mail, Key, ArrowLeft, Code } from 'lucide-svelte';
   import { cn } from "$lib/utils";
   
   import { Button } from "$lib/components/ui/button";
   import { Input } from "$lib/components/ui/input";
   import { Label } from "$lib/components/ui/label";
   import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "$lib/components/ui/card";
   import { Checkbox } from "$lib/components/ui/checkbox";
   import { siGithub, siGoogle, siC, siRust, siGo, siTypescript } from 'simple-icons';
   
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
               description: 'Welcome back to Codepass!'
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

   async function handleGoogleSignIn() {
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
             
               <h1 class="text-4xl font-bold mb-4">Welcome back</h1>
               <p class="text-muted-foreground">Enter your credentials to access your account</p>
           </div>

           <form class="space-y-6" on:submit|preventDefault={handleSignIn}>
               <div class="space-y-2">
                   <Label for="email" class="text-foreground">Email address</Label>
                   <Input 
                       id="email" 
                       type="email" 
                       placeholder="Email address" 
                       class="bg-muted border-muted-foreground/20 text-foreground placeholder:text-muted-foreground focus:border-primary focus:ring-primary"
                       bind:value={email}
                       required
                   />
               </div>
               
               <div class="space-y-2">
                   <Label for="password" class="text-foreground">Password</Label>
                   <Input 
                       id="password" 
                       type="password" 
                       placeholder="Password" 
                       class="bg-muted border-muted-foreground/20 text-foreground placeholder:text-muted-foreground focus:border-primary focus:ring-primary"
                       bind:value={password}
                       required
                   />
               </div>
               
               <div class="flex items-center justify-between">
                   <div class="flex items-center space-x-2">
                       <Checkbox id="remember" bind:checked={rememberMe} class="border-muted-foreground data-[state=checked]:bg-primary data-[state=checked]:border-primary" />
                       <Label for="remember" class="text-sm text-foreground">Remember me</Label>
                   </div>
                   <a href="/forgot-password" class="text-sm text-primary hover:text-primary/80">
                       Forgot password?
                   </a>
               </div>
               
               <Button 
                   type="submit" 
                   class="w-full bg-primary text-primary-foreground hover:bg-primary/90 py-3 font-medium" 
                   disabled={loading}
               >
                   {#if loading}
                       <div class="animate-spin mr-2 h-4 w-4 border-2 border-current border-t-transparent rounded-full"></div>
                   {/if}
                   Sign In
                   <ArrowLeft class="ml-2 h-4 w-4 rotate-180" />
               </Button>
           </form>

           {#if error}
               <div class="bg-destructive/10 border border-destructive/30 text-destructive p-3 rounded-md text-sm" role="alert">
                   {error}
               </div>
           {/if}

           <div class="relative">
               <div class="absolute inset-0 flex items-center">
                   <span class="w-full border-t border-muted-foreground/20"></span>
               </div>
               <div class="relative flex justify-center text-xs uppercase">
                   <span class="bg-background px-2 text-muted-foreground">
                       Or continue with
                   </span>
               </div>
           </div>

           <div class="grid grid-cols-2 gap-3">
               <Button 
                   variant="outline" 
                   class="bg-muted border-muted-foreground/20 text-foreground hover:bg-muted/80" 
                   onclick={handleGoogleSignIn}
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
                   onclick={handleGitHubSignIn}
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

           <p class="text-center text-sm text-muted-foreground">
               Don't have an account?{" "}
               <a href="/signup" class="text-primary hover:text-primary/80 underline">
                   Create one
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