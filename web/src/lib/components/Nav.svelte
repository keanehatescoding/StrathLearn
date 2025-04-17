<script lang="ts">
    import { page } from "$app/stores";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar/index.js";
    import { 
      DropdownMenu,
      DropdownMenuTrigger,
      DropdownMenuContent,
      DropdownMenuItem,
      DropdownMenuLabel,
      DropdownMenuSeparator
    } from "$lib/components/ui/dropdown-menu/index.js";
    import { Sheet, SheetTrigger, SheetContent } from "$lib/components/ui/sheet/index.js";
    import { ModeToggle } from "$lib/components/mode-toggle";
    import { goto } from "$app/navigation";
    
    // Import Lucide icons
    import { 
      Menu, 
      Code, 
      Home, 
      Compass, 
      Codepen, 
      BookOpen, 
      User, 
      Folder, 
      Settings, 
      LogOut 
    } from "lucide-svelte";
  
    // Mock authentication - replace with your auth provider
    // let user = null; // Set to null for logged out state
    let user = { name: "Jane Doe", email: "jane@example.com", image: "" }; // Uncomment for logged in state
    
    let scrolled = false;
    
    // Track scroll position for header effects
    function handleScroll() {
      scrolled = window.scrollY > 10;
    }
  

    $: currentPath = $page.url.pathname;
  </script>
  
  <svelte:window on:scroll={handleScroll} />
  
  <header class="fixed top-0 left-0 right-0 z-50 flex h-16 w-full items-center px-4 md:px-6 transition-all duration-300 
    bg-background/80 backdrop-blur-sm {scrolled ? 'shadow-sm' : ''}">
    <div class="container mx-auto flex items-center justify-between">
      <div class="flex items-center">
        <Sheet>
          <SheetTrigger asChild>
            <Button 
              variant="ghost" 
              size="icon" 
              class="lg:hidden mr-2 hover:bg-accent/10"
            >
              <Menu class="h-5 w-5" />
              <span class="sr-only">Toggle navigation menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" class="w-[280px] p-6">
            <div class="flex items-center mb-8">
              <Code class="h-8 w-8 text-primary mr-3" />
              <span class="text-xl font-bold">Codex</span>
            </div>
            <div class="grid gap-4">
              <a 
                href="/" 
                class="flex items-center gap-3 py-3 px-4 rounded-lg hover:bg-accent text-foreground transition-colors"
              >
                <Home class="h-5 w-5" />
                <span class="font-medium">Home</span>
              </a>
              <a 
                href="/explore" 
                class="flex items-center gap-3 py-3 px-4 rounded-lg hover:bg-accent text-foreground transition-colors"
              >
                <Compass class="h-5 w-5" />
                <span class="font-medium">Explore</span>
              </a>
              <a 
                href="/playground" 
                class="flex items-center gap-3 py-3 px-4 rounded-lg hover:bg-accent text-foreground transition-colors"
              >
                <Codepen class="h-5 w-5" />
                <span class="font-medium">Playground</span>
              </a>
              <a 
                href="/docs" 
                class="flex items-center gap-3 py-3 px-4 rounded-lg hover:bg-accent text-foreground transition-colors"
              >
                <BookOpen class="h-5 w-5" />
                <span class="font-medium">Documentation</span>
              </a>
            </div>
          </SheetContent>
        </Sheet>
        
        <a href="/" class="flex items-center">
          <Code class="h-8 w-8 text-primary mr-2" />
          <span class="text-xl font-bold text-primary">Codex</span>
        </a>
      </div>
      
      <nav class="hidden lg:flex items-center gap-6">
        <a 
          href="/" 
          class="relative font-medium transition-colors {currentPath === '/' 
            ? 'text-primary font-semibold' 
            : 'text-foreground hover:text-foreground/80'}"
        >
          Home
          {#if currentPath === '/'}
            <span class="absolute -bottom-1 left-0 h-0.5 w-full bg-primary rounded-full"></span>
          {/if}
        </a>
        
        <a 
          href="/explore" 
          class="relative font-medium transition-colors {currentPath === '/explore' 
            ? 'text-primary font-semibold' 
            : 'text-foreground hover:text-foreground/80'}"
        >
          Explore
          {#if currentPath === '/explore'}
            <span class="absolute -bottom-1 left-0 h-0.5 w-full bg-primary rounded-full"></span>
          {/if}
        </a>
        
        <a 
          href="/playground" 
          class="relative font-medium transition-colors {currentPath === '/playground' 
            ? 'text-primary font-semibold' 
            : 'text-foreground hover:text-foreground/80'}"
        >
          Playground
          {#if currentPath === '/playground'}
            <span class="absolute -bottom-1 left-0 h-0.5 w-full bg-primary rounded-full"></span>
          {/if}
        </a>
        
        <a 
          href="/docs" 
          class="relative font-medium transition-colors {currentPath === '/docs' 
            ? 'text-primary font-semibold' 
            : 'text-foreground hover:text-foreground/80'}"
        >
          Documentation
          {#if currentPath === '/docs'}
            <span class="absolute -bottom-1 left-0 h-0.5 w-full bg-primary rounded-full"></span>
          {/if}
        </a>
      </nav>
          
      <div class="flex items-center gap-4">
        <ModeToggle />
        
        {#if user}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" class="rounded-full h-9 w-9 p-0">
                <Avatar class="h-9 w-9 border-2 border-primary/10">
                  <AvatarImage 
                    src={user.image || undefined} 
                    alt={user.name || "User"}
                  />
                  <AvatarFallback class="bg-primary text-primary-foreground">
                    {user.name?.charAt(0) || "U"}
                  </AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" class="w-56">
              <DropdownMenuLabel class="font-normal">
                <div class="flex flex-col space-y-1">
                  <p class="text-sm font-medium">{user.name}</p>
                  <p class="text-xs text-muted-foreground">{user.email}</p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <User class="mr-2 h-4 w-4" />
                <span>Profile</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Folder class="mr-2 h-4 w-4" />
                <span>My Projects</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Settings class="mr-2 h-4 w-4" />
                <span>Settings</span>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem on:click={() => { /* Add logout logic */ }}>
                <LogOut class="mr-2 h-4 w-4" />
                <span>Log out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        {:else}
          <div class="flex gap-3">
            <Button variant="outline" class="hover:bg-accent/10" on:click={() => goto("/signin")}>
              Log in
            </Button>
            <Button class="bg-primary text-primary-foreground hover:bg-primary/90" on:click={() => goto("/signup")}>
              Sign up
            </Button>
          </div>
        {/if}
      </div>
    </div>
  </header>
  
  <!-- Add some spacing below the navbar for content -->
  <div class="pt-16"></div>
  
  <style>
    /* Any custom styles can go here */
    :global(.dark) .logo-text {
      filter: brightness(1.2);
    }
  </style>