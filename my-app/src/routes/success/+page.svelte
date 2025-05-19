<script lang="ts">
  import { onMount } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';

  // Import icons
  let CheckCircle, Zap, ArrowRight, Sparkles;
  
  onMount(async () => {
    // Dynamically import icons to prevent SSR issues
    const module = await import('lucide-svelte');
    CheckCircle = module.CheckCircle;
    Zap = module.Zap;
    ArrowRight = module.ArrowRight;
    Sparkles = module.Sparkles;
  });

  // Animation timing
  let visible = false;
  onMount(() => {
    setTimeout(() => {
      visible = true;
    }, 100);
  });

  // Utility function for class names
  const cn = (...classes) => {
    return classes.filter(Boolean).join(' ');
  };
</script>

<div class="min-h-screen w-full bg-gradient-to-b from-background to-background/95 flex items-center justify-center p-4">
  <!-- Aurora Background Effect -->
  <div class="absolute inset-0 overflow-hidden pointer-events-none">
    <div class="pointer-events-none absolute -top-1/2 right-0 h-[40vw] w-[40vw] animate-aurora-1 bg-purple-500/20 mix-blend-overlay blur-[6rem]"></div>
    <div class="pointer-events-none absolute left-0 top-0 h-[40vw] w-[40vw] animate-aurora-2 bg-blue-500/20 mix-blend-overlay blur-[6rem]"></div>
    <div class="pointer-events-none absolute bottom-0 left-0 h-[40vw] w-[40vw] animate-aurora-3 bg-green-500/20 mix-blend-overlay blur-[6rem]"></div>
    <div class="pointer-events-none absolute -bottom-1/2 right-0 h-[40vw] w-[40vw] animate-aurora-4 bg-teal-500/20 mix-blend-overlay blur-[6rem]"></div>
  </div>
  
  <!-- Subtle Grid Background -->
  <div class="absolute inset-0 bg-grid-pattern opacity-[0.03] pointer-events-none"></div>
  
  <!-- Content -->
  <div class="relative z-10 max-w-3xl w-full px-4">
    {#if visible}
      <div in:fly={{ y: 50, duration: 800, easing: quintOut }} class="relative overflow-hidden">
        <!-- Success Card -->
        <div class="bg-card/80 backdrop-blur-md border border-primary/20 rounded-2xl overflow-hidden shadow-2xl">
          <!-- Glass Highlight Effect -->
          <div class="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-primary/30 to-transparent"></div>
          <div class="absolute inset-y-0 right-0 w-px bg-gradient-to-b from-transparent via-primary/30 to-transparent"></div>
          <div class="absolute inset-x-0 bottom-0 h-px bg-gradient-to-r from-transparent via-primary/30 to-transparent"></div>
          <div class="absolute inset-y-0 left-0 w-px bg-gradient-to-b from-transparent via-primary/30 to-transparent"></div>
          
          <div class="p-10 md:p-12">
            <!-- Success Icon -->
            <div in:fade={{ delay: 300, duration: 700 }} class="flex justify-center mb-8">
              <div class="relative">
                <div class="absolute inset-0 animate-pulse-slow bg-primary/20 rounded-full blur-xl"></div>
                <div class="relative h-24 w-24 rounded-full bg-gradient-to-br from-green-100/80 to-green-400/30 dark:from-green-900/50 dark:to-green-600/20 flex items-center justify-center border border-green-400/30 shadow-lg">
                  {#if CheckCircle}
                    <svelte:component this={CheckCircle} class="h-12 w-12 text-green-600 dark:text-green-400" />
                  {:else}
                    <div class="h-12 w-12 bg-green-600 dark:bg-green-400 rounded-full"></div>
                  {/if}
                </div>
              </div>
            </div>
            
            <!-- Success Message -->
            <div class="text-center mb-10">
              <div in:fade={{ delay: 400, duration: 700 }} class="flex items-center justify-center gap-3 mb-4">
                {#if Sparkles}
                  <svelte:component this={Sparkles} class="h-6 w-6 text-primary" />
                {/if}
                <h1 class="text-4xl md:text-5xl font-bold tracking-tight font-display bg-clip-text text-transparent bg-gradient-to-r from-primary to-green-500 dark:from-primary dark:to-green-400">
                  Payment Successful!
                </h1>
                {#if Sparkles}
                  <svelte:component this={Sparkles} class="h-6 w-6 text-primary" />
                {/if}
              </div>
              <p in:fade={{ delay: 500, duration: 700 }} class="text-lg text-foreground/80 max-w-lg mx-auto leading-relaxed">
                Welcome to <span class="font-semibold text-primary">Codex</span>! Your subscription is now active and you have full access to all our interactive challenges and learning materials.
              </p>
            </div>
            
            <!-- What's Next -->
            <div in:fade={{ delay: 600, duration: 700 }} class="relative bg-muted/30 rounded-xl p-7 border border-primary/10 mb-10 backdrop-blur-sm">
              <!-- Subtle glow behind the card -->
              <div class="absolute -inset-px bg-gradient-to-tr from-primary/5 to-green-500/5 blur-sm rounded-xl"></div>
              
              <div class="relative">
                <h3 class="text-xl font-semibold mb-6 flex items-center font-display">
                  {#if Zap}
                    <svelte:component this={Zap} class="h-5 w-5 text-primary mr-2" />
                  {/if}
                  What's Next
                </h3>
                
                <ul class="space-y-4">
                  {#each [
                    "Access 100+ interactive coding challenges",
                    "Start with C programming fundamentals",
                    "Track your progress with our achievement system",
                    "Join our community of learners",
                    "Build your first project within 30 days"
                  ] as item, i}
                    <li in:fade={{ delay: 700 + (i * 100), duration: 300 }} class="flex items-start group">
                      <div class="h-6 w-6 rounded-full bg-primary/10 flex items-center justify-center shrink-0 mr-3 mt-0.5 group-hover:bg-primary/20 transition-colors">
                        {#if CheckCircle}
                          <svelte:component this={CheckCircle} class="h-4 w-4 text-primary" />
                        {/if}
                      </div>
                      <span class="text-foreground/90">{item}</span>
                    </li>
                  {/each}
                </ul>
              </div>
            </div>
            
            <!-- CTA Buttons -->
            <div in:fade={{ delay: 1200, duration: 700 }} class="flex flex-col sm:flex-row gap-4 justify-center">
              <a href="/challenge" class="relative inline-flex items-center justify-center gap-2 px-8 py-3.5 font-medium rounded-lg overflow-hidden bg-gradient-to-br from-primary to-green-500 text-primary-foreground hover:shadow-lg hover:shadow-primary/20 transition-all duration-300 group">
                <span class="absolute inset-0 bg-white/10 opacity-0 group-hover:opacity-20 transition-opacity"></span>
                {#if Zap}
                  <svelte:component this={Zap} class="h-5 w-5" />
                {/if}
                <span>Start Learning Now</span>
              </a>
              <a href="/dashboard" class="relative inline-flex items-center justify-center gap-2 px-8 py-3.5 font-medium rounded-lg bg-card border border-muted hover:border-muted/80 hover:bg-card/80 transition-all duration-300">
                <span>View Dashboard</span>
                {#if ArrowRight}
                  <svelte:component this={ArrowRight} class="h-5 w-5 transition-transform group-hover:translate-x-0.5" />
                {/if}
              </a>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  /* Aurora animations */
  @keyframes aurora-1 {
    0%, 100% { transform: translateY(-30%) translateX(10%) scale(1); opacity: 0.8; }
    50% { transform: translateY(-35%) translateX(15%) scale(1.2); opacity: 0.6; }
  }
  
  @keyframes aurora-2 {
    0%, 100% { transform: translateY(0%) translateX(0%) scale(1); opacity: 0.6; }
    50% { transform: translateY(5%) translateX(-5%) scale(1.1); opacity: 0.8; }
  }
  
  @keyframes aurora-3 {
    0%, 100% { transform: translateY(0%) translateX(0%) scale(1); opacity: 0.7; }
    50% { transform: translateY(-5%) translateX(5%) scale(1.15); opacity: 0.5; }
  }
  
  @keyframes aurora-4 {
    0%, 100% { transform: translateY(30%) translateX(-10%) scale(1); opacity: 0.7; }
    50% { transform: translateY(35%) translateX(-15%) scale(1.1); opacity: 0.5; }
  }
  
  .animate-aurora-1 {
    animation: aurora-1 15s ease-in-out infinite;
  }
  
  .animate-aurora-2 {
    animation: aurora-2 18s ease-in-out infinite;
  }
  
  .animate-aurora-3 {
    animation: aurora-3 21s ease-in-out infinite;
  }
  
  .animate-aurora-4 {
    animation: aurora-4 24s ease-in-out infinite;
  }
  
  .animate-pulse-slow {
    animation: pulse 3s ease-in-out infinite;
  }
  
  /* Grid background pattern */
  .bg-grid-pattern {
    background-image: 
      linear-gradient(to right, rgba(127, 127, 127, 0.1) 1px, transparent 1px),
      linear-gradient(to bottom, rgba(127, 127, 127, 0.1) 1px, transparent 1px);
    background-size: 40px 40px;
  }
  
  /* Font styles */
  .font-display {
    font-family: 'Orbitron', sans-serif;
  }
  
  /* Custom scrollbar */
  :global(.overflow-auto::-webkit-scrollbar) {
    width: 6px;
    height: 6px;
  }
  
  :global(.overflow-auto::-webkit-scrollbar-track) {
    background: transparent;
  }
  
  :global(.overflow-auto::-webkit-scrollbar-thumb) {
    background: rgba(127, 127, 127, 0.2);
    border-radius: 3px;
  }
  
  :global(.overflow-auto::-webkit-scrollbar-thumb:hover) {
    background: rgba(127, 127, 127, 0.3);
  }
</style>