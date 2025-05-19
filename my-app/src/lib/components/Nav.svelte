<script lang="ts">
  import { page } from "$app/stores";
  import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar";
  import { Button } from "$lib/components/ui/button";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import ModeToggle from "$lib/components/mode-toggle.svelte";
  import { goto } from "$app/navigation";
  import { browser } from "$app/environment";
  import { onMount, tick } from "svelte";
  import { Menu, Home, Code, BookOpen, Users, User, Calendar, CreditCard, LogOut, X, ChevronDown } from "lucide-svelte";
  import { siDiscord } from "simple-icons";
  import { useSession, authClient } from "$lib/auth-client";
  import { slide, fade } from 'svelte/transition';
  
  // Reactive declarations
  let scrolled = false;
  let isMenuOpen = false;
  let isProfileOpen = false;
  let navVisible = true;
  let lastScrollY = 0;
  const session = useSession();
  
  // Handle scroll effects
  onMount(() => {
    if (browser) {
      const handleScroll = () => {
        const currentScrollY = window.scrollY;
        
        // Determine if scrolled for styling
        scrolled = currentScrollY > 10;
        
        // Handle nav visibility (hide on scroll down, show on scroll up)
        if (currentScrollY > lastScrollY + 50) {
          navVisible = false;
        } else if (currentScrollY < lastScrollY - 10) {
          navVisible = true;
        }
        
        lastScrollY = currentScrollY;
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

<header class="fixed top-0 left-0 right-0 z-50 w-full transition-all duration-300 transform {navVisible ? 'translate-y-0' : '-translate-y-full'} 
  {scrolled ? 'bg-background/95 shadow-md backdrop-blur-sm' : 'bg-background/80 backdrop-blur-sm'}"
  style="height: 68px;">
  <div class="container mx-auto h-full flex items-center justify-between px-4">
    <!-- Logo and mobile menu -->
    <div class="flex items-center gap-2">
      <button 
        class="lg:hidden flex items-center justify-center w-10 h-10 rounded-full bg-accent/20 hover:bg-accent/30 transition-colors menu-trigger" 
        onclick={toggleMenu}
        aria-label="Toggle menu"
      >
        {#if isMenuOpen}
          <X class="h-5 w-5 text-primary" />
        {:else}
          <Menu class="h-5 w-5 text-primary" />
        {/if}
      </button>
      
      <a href="/" class="flex items-center space-x-1 group" onclick={() => closeMenu()}>
        <img src="/logo.png" alt="Logo" class="h-10 w-auto" />
        <span class="text-xl font-bold bg-gradient-to-r from-primary to-primary/80 bg-clip-text text-transparent group-hover:from-primary/90 group-hover:to-primary/70 transition-all">
          Codepass
        </span>
      </a>
    </div>
    
    <!-- Desktop navigation -->
    <nav class="hidden lg:flex items-center">
      <div class="flex items-center bg-accent/20 rounded-full p-1 mx-2">
        {#each [
          { path: '/', label: 'Home', icon: Home },
          { path: '/projects', label: 'Projects', icon: Code },
          { path: '/learn', label: 'Learn', icon: BookOpen },
          { path: '/community', label: 'Community', icon: Users }
        ] as route}
          <a 
            href={route.path}
            class="relative font-medium transition-all px-4 py-2 rounded-full flex items-center gap-1.5 {
              isActive(route.path) 
                ? 'bg-primary text-primary-foreground shadow-sm' 
                : 'text-foreground hover:bg-accent/50'
            }"
          >
            <svelte:component this={route.icon} class="h-4 w-4" />
            {route.label}
          </a>
        {/each}
      </div>
      
      <a 
        href="https://discord.gg" 
        target="_blank" 
        rel="noopener noreferrer"
        class="transition-all ml-2 hover:scale-105"
      >
        <div class="flex items-center justify-center bg-indigo-600 hover:bg-indigo-700 text-white rounded-full px-4 py-2 shadow-sm">
          <svg
            role="img"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4 mr-1.5"
            fill="currentColor"
          >
            <path d={siDiscord.path} />
          </svg>
          <span class="text-sm font-medium">Discord</span>
        </div>
      </a>
    </nav>
    
    <!-- Right side controls -->
    <div class="flex items-center gap-3">
      <ModeToggle />
      
      {#if $session.data?.user}
        <DropdownMenu.Root>
          <DropdownMenu.Trigger>
            <Button variant="ghost" class="rounded-full h-10 w-10 p-0 hover:bg-accent/30">
              <Avatar class="h-10 w-10 border-2 border-primary/30 transition-all hover:border-primary">
                <AvatarImage src={$session.data.user.image || ''} alt={$session.data.user.name || "User"} />
                <AvatarFallback class="bg-primary text-primary-foreground">
                  {$session.data.user.name ? $session.data.user.name[0].toUpperCase() : "U"}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end" class="w-64 p-2">
            <div class="flex items-center space-x-3 p-2 mb-2 bg-accent/20 rounded-lg">
              <Avatar class="h-12 w-12 border border-border">
                <AvatarImage src={$session.data.user.image || ''} alt={$session.data.user.name || "User"} />
                <AvatarFallback class="bg-primary text-primary-foreground">
                  {$session.data.user.name ? $session.data.user.name[0].toUpperCase() : "U"}
                </AvatarFallback>
              </Avatar>
              <div class="flex flex-col">
                <p class="text-base font-semibold">{$session.data.user.name || "User"}</p>
                <p class="text-xs text-muted-foreground truncate">{$session.data.user.email || ""}</p>
              </div>
            </div>
            
            <div class="grid grid-cols-1 gap-1 mb-2">
              <DropdownMenu.Item class="py-2 cursor-pointer" on:click={() => navigateTo('/profile')}>
                <User class="mr-2 h-4 w-4" />
                <span>Profile</span>
              </DropdownMenu.Item>
              <DropdownMenu.Item class="py-2 cursor-pointer" onclick={() => navigateTo('/profile')}>
                <Calendar class="mr-2 h-4 w-4" />
                <span>Profile</span>
              </DropdownMenu.Item>
              <DropdownMenu.Item class="py-2 cursor-pointer" onclick={() => navigateTo('/api/auth/portal')}>
                <CreditCard class="mr-2 h-4 w-4" />
                <span>Payments </span>
              </DropdownMenu.Item>
            </div>
            
            <div class="pt-2 border-t border-border">
              <Button 
                variant="destructive" 
                class="w-full justify-center mt-1" 
                onclick={handleSignOut}
              >
                <LogOut class="mr-2 h-4 w-4" />
                <span>Sign Out</span>
              </Button>
            </div>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      {:else}
        <!-- Auth buttons -->
        <div class="flex items-center gap-3">
          <Button 
            variant="outline" 
            class="hover:bg-accent/30 border-primary/30 hover:border-primary font-medium px-4"
            onclick={() => navigateTo('/signin')}
          >
            Log in
          </Button>
          <Button 
            variant="default" 
            class="bg-primary hover:bg-primary/90 font-medium px-4"
            onclick={() => navigateTo('/signup')}
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
    class="fixed inset-0 z-40 bg-black/40 lg:hidden"
    onclick={closeMenu}
    transition:fade={{ duration: 200 }}
  ></div>
  
  <div 
    class="fixed top-[68px] left-0 bottom-0 z-40 w-[85%] max-w-[320px] bg-background border-r border-border shadow-xl lg:hidden mobile-menu overflow-y-auto"
    transition:slide={{ duration: 300, axis: 'x' }}
  >
    <div class="flex flex-col h-full">
      <div class="p-4">
        <div class="space-y-1">
          {#each [
            { path: '/', label: 'Home', icon: Home },
            { path: '/projects', label: 'Projects', icon: Code },
            { path: '/learn', label: 'Learn', icon: BookOpen },
            { path: '/community', label: 'Community', icon: Users }
          ] as route}
            <a 
              href={route.path}
              onclick={() => closeMenu()}
              class="flex items-center gap-3 p-3 rounded-lg transition-all {
                isActive(route.path) 
                  ? 'bg-primary text-primary-foreground font-medium' 
                  : 'hover:bg-accent/50 text-foreground'
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
            class="flex items-center gap-3 p-3 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 transition-colors mt-3"
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
        </div>
      </div>
      
      <div class="mt-auto p-4 border-t border-border bg-accent/10">
        {#if $session.data?.user}
          <div class="mb-4">
            <div class="flex items-center p-3 bg-accent/30 rounded-lg">
              <Avatar class="h-10 w-10 mr-3 border border-primary/30">
                <AvatarImage src={$session.data.user.image || ''} alt={$session.data.user.name || "User"} />
                <AvatarFallback class="bg-primary text-primary-foreground">
                  {$session.data.user.name ? $session.data.user.name[0].toUpperCase() : "U"}
                </AvatarFallback>
              </Avatar>
              <div class="overflow-hidden">
                <p class="font-medium truncate">{$session.data.user.name || "User"}</p>
                <p class="text-xs text-muted-foreground truncate">{$session.data.user.email || ""}</p>
              </div>
            </div>
          </div>
          
          <div class="grid grid-cols-1 gap-2">
            <Button variant="outline" class="justify-start" onclick={() => navigateTo('/profile')}>
              <User class="mr-2 h-4 w-4" />
              <span>Profile</span>
            </Button>
            <Button variant="outline" class="justify-start" onclick={() => navigateTo('/dashboard')}>
              <Calendar class="mr-2 h-4 w-4" />
              <span>Dashboard</span>
            </Button>
            <Button variant="outline" class="justify-start" onclick={() => navigateTo('/settings')}>
              <CreditCard class="mr-2 h-4 w-4" />
              <span>Settings</span>
            </Button>
            <Button variant="destructive" class="justify-start" onclick={handleSignOut}>
              <LogOut class="mr-2 h-4 w-4" />
              <span>Sign out</span>
            </Button>
          </div>
        {:else}
          <div class="space-y-3">
            <Button 
              variant="outline" 
              class="w-full justify-center font-medium" 
              onclick={() => navigateTo('/signin')}
            >
              Log in
            </Button>
            <Button 
              variant="default" 
              class="w-full justify-center font-medium" 
              onclick={() => navigateTo('/signup')}
            >
              Sign up
            </Button>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  /* Add smooth transitions */
  @keyframes glow {
    0% { box-shadow: 0 0 5px rgba(var(--primary-rgb), 0.3); }
    50% { box-shadow: 0 0 15px rgba(var(--primary-rgb), 0.5); }
    100% { box-shadow: 0 0 5px rgba(var(--primary-rgb), 0.3); }
  }
  
  :global(.menu-trigger) {
    transition: transform 0.2s ease;
  }
  
  :global(.menu-trigger:hover) {
    transform: rotate(15deg);
  }
  
  :global(a.active .menu-item-text) {
    font-weight: 600;
  }
</style>