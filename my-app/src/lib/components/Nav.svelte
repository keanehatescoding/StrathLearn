<script lang="ts">
  import { page } from "$app/stores";
  import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar";
  import { Button } from "$lib/components/ui/button";
  import { Sheet, SheetContent, SheetTrigger } from "$lib/components/ui/sheet";

  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import ModeToggle from "$lib/components/mode-toggle.svelte";
  import { goto } from "$app/navigation";
  import { browser } from "$app/environment";
  import { onMount } from "svelte";
  import { Menu, Home, Code, BookOpen, Users, User, Calendar, CreditCard, LogOut, X } from "lucide-svelte";
  import { siDiscord } from "simple-icons";
  import { useSession, authClient } from "$lib/auth-client";
  
  // Reactive declarations
  let scrolled = false;
  let isMenuOpen = false;
  let isProfileOpen = false;
  const session = useSession();
  
  // Handle scroll effects
  onMount(() => {
    if (browser) {
      const handleScroll = () => {
        scrolled = window.scrollY > 10;
      };
      
      window.addEventListener("scroll", handleScroll);
      handleScroll(); // Check initial scroll position
      
      // Cleanup event listener on unmount
      return () => window.removeEventListener("scroll", handleScroll);
    }
  });
  
  // Auth functions
  async function handleSignOut() {
    try {
      await authClient.signOut();
      goto('/signin');
    } catch (error) {
      console.error("Sign out failed:", error);
    }
  }
  
  // Navigation helpers
  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
    if (isMenuOpen) isProfileOpen = false;
  }
  
  function closeMenu() {
    if (isMenuOpen) isMenuOpen = false;
  }
  
  function navigateTo(path: string) {
    closeMenu();
    goto(path);
  }

  // Active route helper
  function isActive(path: string): boolean {
    return $page.url.pathname === path;
  }

  // Handle click outside
  function handleClickOutside(event: MouseEvent) {
    if (isMenuOpen && !(event.target as Element).closest('.mobile-menu, .menu-trigger')) {
      isMenuOpen = false;
    }
  }

  onMount(() => {
    if (browser) {
      document.addEventListener('click', handleClickOutside);
      return () => document.removeEventListener('click', handleClickOutside);
    }
  });
</script>

<header class="fixed top-0 left-0 right-0 z-50 flex h-16 w-full items-center px-4 md:px-6 transition-all duration-300 
  {scrolled ? 'bg-background/95 shadow-md backdrop-blur-sm' : 'bg-background/80 backdrop-blur-sm'}">
  <div class="container mx-auto flex items-center justify-between">
    <!-- Logo and mobile menu -->
    <div class="flex items-center">
      <Button 
        variant="ghost" 
        size="icon" 
        class="lg:hidden mr-2 hover:bg-accent {isMenuOpen ? 'bg-accent/20' : ''}" 
        on:click={toggleMenu}
      >
        {#if isMenuOpen}
          <X class="h-5 w-5" />
        {:else}
          <Menu class="h-5 w-5" />
        {/if}
        <span class="sr-only">Toggle navigation menu</span>
      </Button>
      
      <a href="/" class="flex items-center group" on:click={() => closeMenu()}>
        <div class="bg-primary h-8 w-8 rounded flex items-center justify-center mr-3 transition-transform group-hover:scale-110">
          <Code class="h-5 w-5 text-primary-foreground" />
        </div>
        <span class="text-xl font-bold text-primary hidden md:inline-block group-hover:text-primary/90 transition-colors">
          Codex
        </span>
      </a>
    </div>
    
    <!-- Desktop navigation -->
    <nav class="hidden lg:flex items-center gap-8">
      {#each [
        { path: '/', label: 'Home', icon: Home },
        { path: '/projects', label: 'Projects', icon: Code },
        { path: '/learn', label: 'Learn', icon: BookOpen },
        { path: '/community', label: 'Community', icon: Users }
      ] as route}
        <a 
          href={route.path}
          class="relative font-medium transition-colors py-1 px-1 {
            isActive(route.path) 
              ? 'text-primary font-semibold' 
              : 'text-foreground hover:text-primary'
          }"
        >
          {route.label}
          {#if isActive(route.path)}
            <span class="absolute -bottom-1 left-0 h-0.5 w-full bg-primary rounded-full" />
          {/if}
        </a>
      {/each}
      
      <a 
        href="https://discord.gg" 
        target="_blank" 
        rel="noopener noreferrer"
        class="transition-transform hover:scale-110"
      >
        <Button variant="outline" size="icon" class="rounded-full hover:bg-primary/10">
          <svg
            role="img"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="currentColor"
          >
            <path d={siDiscord.path} />
          </svg>
          <span class="sr-only">Discord</span>
        </Button>
      </a>
    </nav>
    
    <!-- Right side controls -->
    <div class="flex items-center gap-3">
      <ModeToggle />
      
      {#if $session.data?.user}

        <DropdownMenu.Root>
          <DropdownMenu.Trigger >
            <Button variant="ghost" class="rounded-full h-9 w-9 p-0 hover:bg-accent/20">
              <Avatar class="h-9 w-9 border border-border transition-all hover:border-primary/50">
                <AvatarImage src={$session.data.user.image || ''} alt={$session.data.user.name || "User"} />
                <AvatarFallback class="bg-primary text-primary-foreground">
                  {$session.data.user.name ? $session.data.user.name[0].toUpperCase() : "U"}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end" class="w-56">
            <DropdownMenu.Label class="font-normal">
              <div class="flex flex-col space-y-1">
                <p class="text-sm font-medium">{$session.data.user.name || "User"}</p>
                <p class="text-xs text-muted-foreground truncate">{$session.data.user.email || ""}</p>
              </div>
            </DropdownMenu.Label>
            <DropdownMenu.Separator />
            <DropdownMenu.Item on:click={() => navigateTo('/profile')}>
              <User class="mr-2 h-4 w-4" />
              <span>Profile</span>
            </DropdownMenu.Item>
            <DropdownMenu.Item on:click={() => navigateTo('/dashboard')}>
              <Calendar class="mr-2 h-4 w-4" />
              <span>Dashboard</span>
            </DropdownMenu.Item>
            <DropdownMenu.Item on:click={() => navigateTo('/settings')}>
              <CreditCard class="mr-2 h-4 w-4" />
              <span>Settings</span>
            </DropdownMenu.Item>
            <DropdownMenu.Separator />
            <DropdownMenu.Item on:click={handleSignOut}>
              <LogOut class="mr-2 h-4 w-4" />
              <span>Log out</span>
            </DropdownMenu.Item>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      {:else}
        <!-- Auth buttons -->
        <div class="flex gap-2">
          <Button 
            variant="ghost" 
            on:click={() => navigateTo('/login')}
            class="hover:bg-accent/20"
          >
            Log in
          </Button>
          <Button 
            variant="default" 
            on:click={() => navigateTo('/signup')}
            class="hover:bg-primary/90"
          >
            Sign up
          </Button>
        </div>
      {/if}
    </div>
  </div>
</header>

<!-- Mobile menu overlay -->
{#if isMenuOpen}
  <div 
    class="fixed top-16 left-0 right-0 z-40 bg-background/95 backdrop-blur-sm border-b border-border shadow-lg lg:hidden pt-2 pb-4 px-4 mobile-menu animate-in fade-in slide-in-from-top-2"
    transition:slide={{ duration: 200 }}
  >
    <div class="flex flex-col space-y-2">
      {#each [
        { path: '/', label: 'Home', icon: Home },
        { path: '/projects', label: 'Projects', icon: Code },
        { path: '/learn', label: 'Learn', icon: BookOpen },
        { path: '/community', label: 'Community', icon: Users }
      ] as route}
        <a 
          href={route.path}
          on:click={() => closeMenu()}
          class="flex items-center gap-3 py-3 px-3 rounded-md transition-all {
            isActive(route.path) 
              ? 'bg-primary/10 text-primary font-medium' 
              : 'hover:bg-accent/30 text-foreground'
          }"
        >
          <svelte:component this={route.icon} class="h-5 w-5" />
          <span class="font-medium">{route.label}</span>
        </a>
      {/each}
      
      <a 
        href="https://discord.gg" 
        target="_blank" 
        rel="noopener noreferrer"
        class="flex items-center gap-3 py-3 px-3 rounded-md hover:bg-accent/30 text-foreground transition-colors mt-1"
      >
        <svg
          role="img"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          fill="currentColor"
        >
          <path d={siDiscord.path} />
        </svg>
        <span class="font-medium">Join Discord</span>
      </a>
      
      {#if $session.data?.user}
        <div class="pt-4 border-t border-border mt-2">
          <div class="bg-secondary/40 border border-border rounded-lg py-2 px-3 mb-4 flex items-center">
            <div class="w-8 h-8 rounded-full overflow-hidden mr-3">
              <Avatar class="h-8 w-8">
                <AvatarImage src={$session.data.user.image || ''} alt={$session.data.user.name || "User"} />
                <AvatarFallback class="bg-primary text-primary-foreground">
                  {$session.data.user.name ? $session.data.user.name[0].toUpperCase() : "U"}
                </AvatarFallback>
              </Avatar>
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-medium">{$session.data.user.name || "User"}</span>
              <span class="text-xs text-muted-foreground truncate max-w-[200px]">{$session.data.user.email || ""}</span>
            </div>
          </div>
          
          <div class="grid grid-cols-2 gap-3">
            <Button variant="outline" on:click={() => navigateTo('/profile')} class="w-full justify-start hover:bg-accent/20">
              <User class="mr-2 h-4 w-4" />
              <span>Profile</span>
            </Button>
            <Button variant="outline" on:click={() => navigateTo('/dashboard')} class="w-full justify-start hover:bg-accent/20">
              <Calendar class="mr-2 h-4 w-4" />
              <span>Dashboard</span>
            </Button>
            <Button variant="outline" on:click={() => navigateTo('/settings')} class="w-full justify-start hover:bg-accent/20">
              <CreditCard class="mr-2 h-4 w-4" />
              <span>Settings</span>
            </Button>
            <Button variant="destructive" on:click={handleSignOut} class="w-full justify-start hover:bg-red-600">
              <LogOut class="mr-2 h-4 w-4" />
              <span>Log out</span>
            </Button>
          </div>
        </div>
      {:else}
        <div class="pt-4 border-t border-border mt-2">
          <div class="flex flex-col sm:flex-row gap-3">
            <Button variant="outline" onclick={() => navigateTo('/login')} class="w-full hover:bg-accent/20">
              Log in
            </Button>
            <Button variant="default" onclick={() => navigateTo('/signup')} class="w-full hover:bg-primary/90">
              Sign up
            </Button>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}