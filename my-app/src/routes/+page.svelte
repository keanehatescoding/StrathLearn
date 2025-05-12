<script lang="ts">
  import { onMount } from 'svelte';
  import { fade, fly, slide } from 'svelte/transition';
  import { 
    Code, Terminal, Zap, BookOpen, 
    ChevronRight, CheckCircle, Github, 
    Sparkles, ArrowRight, Braces, 
    Laptop, Users, Award, Rocket, 
    Play, CheckSquare, XSquare, Clock,
    RefreshCw, Download, Copy, Maximize2
  } from 'lucide-svelte';
  
  // Import UI components
  import { Button } from "$lib/components/ui/button";
  import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar";
  import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "$lib/components/ui/card";
  import { Badge } from "$lib/components/ui/badge";
  import { Tabs, TabsContent, TabsList, TabsTrigger } from "$lib/components/ui/tabs";
  import { Progress } from "$lib/components/ui/progress";
  
  // Import simple-icons
  import { 
    siC, siCplusplus, siPython, siJavascript, 
    siGooglecloudcomposer, siRust, siGo, siRuby, 
    siReact, siSvelte, siTailwindcss, siTypescript, 
    siVuedotjs, siNodedotjs, siDjango, siLaravel, 
    siPostgresql, siGooglecloud, 
    siFirebase, siDocker, siVercel,

	siPhp,

	siSwift,

	siHtml5,

	siHaskell,

	siKotlin





  } from 'simple-icons';
  
  // Mock editor state
  let editorContent = `#include <stdio.h>

int main() {
  // Welcome to Codex Learning Platform
  printf("Hello, World!\\n");
  
  // Your coding journey starts here
  int progress = 100;
  
  if (progress == 100) {
    printf("You're ready to begin!\\n");
  }
  
  return 0;
}`;

  let editorTheme = 'dark';
  let activeTab = 'c';
  let isTyping = false;
  let cursorPosition = 0;
  let typingInterval: ReturnType<typeof setInterval>;
  let editorElement: HTMLElement;
  let lineNumbers: number[] = [];
  let isRunning = false;
  let showOutput = false;
  let outputContent = "";
  let executionTime = 0;
  let executionSuccess = true;
  
  // Testimonials data
  const testimonials = [
    {
      name: "Alex Johnson",
      role: "Software Engineer",
      avatar: "/placeholder.svg?height=40&width=40",
      content: "Codex has transformed how I approach learning programming languages. The interactive challenges are engaging and practical."
    },
    {
      name: "Sarah Chen",
      role: "Computer Science Student",
      avatar: "/placeholder.svg?height=40&width=40",
      content: "As a student, Codex has been invaluable. The immediate feedback and progression system keeps me motivated."
    },
    {
      name: "Michael Rodriguez",
      role: "Full Stack Developer",
      avatar: "/placeholder.svg?height=40&width=40",
      content: "The best platform I've found for mastering C programming. The editor is intuitive and the challenges are well-designed."
    }
  ];
  
  // Features data
  const features = [
    {
      icon: Terminal,
      title: "Interactive Challenges",
      description: "Practice with real-world coding challenges that test your skills and reinforce learning."
    },
    {
      icon: Zap,
      title: "Instant Feedback",
      description: "Get immediate results and detailed explanations for your code submissions."
    },
    {
      icon: BookOpen,
      title: "Comprehensive Curriculum",
      description: "Follow a structured learning path from basics to advanced programming concepts."
    },
    {
      icon: Users,
      title: "Community Support",
      description: "Connect with fellow learners and mentors to solve problems together."
    },
    {
      icon: Award,
      title: "Achievement System",
      description: "Earn badges and certificates as you progress through the curriculum."
    },
    {
      icon: Rocket,
      title: "Career Preparation",
      description: "Build a portfolio of projects that demonstrate your skills to potential employers."
    }
  ];
  
  // Languages supported with icons
  const languages = [
    { id: 'c', name: 'C', icon: siC, color: '#3949AB', available: true },
    { id: 'cpp', name: 'C++', icon: siCplusplus, color: '#00599C', coming: true },
    {id: 'tailwind', name: 'Tailwind CSS', icon: siTailwindcss, color: '#38BDF8', coming: true},
  
    {id: 'html', name: "HTML", icon: siHtml5, color: "#E34F26", coming: true},
    { id: 'typescript', name: 'TypeScript', icon: siTypescript, color: '#007ACC', coming: true },
    { id: 'php', name: 'PHP', icon: siPhp, color: '#4F5B93', coming: true },
    { id: 'swift', name: 'Swift', icon: siSwift, color: '#F05138', coming: true },
    {id:"kotlin", name:"Kotlin", icon: siKotlin, color:"#F18E33", coming: true},
    { id: 'python', name: 'Python', icon: siPython, color: '#3776AB', coming: true },
    { id: 'javascript', name: 'JavaScript', icon: siJavascript, color: '#F7DF1E', coming: true },
   
    { id: 'firebase', name: 'Firebase', icon: siFirebase, color: '#FFCA28', coming: true },
    { id: 'postgresql', name: 'PostgreSQL', icon: siPostgresql, color: '#336791', coming: true },
    { id: 'docker', name: 'Docker', icon: siDocker, color: '#2496ED', coming: true },
    { id: 'rust', name: 'Rust', icon: siRust, color: '#FF4500', coming: true },
    { id: 'go', name: 'Go', icon: siGo, color: '#00ADD8', coming: true },
   
    {id: 'haskell', name:"Haskell", icon: siHaskell, color:"#5e5086", coming: true},
  ];
  
  // Initialize on component mount
  onMount(() => {
    // Generate line numbers
    updateLineNumbers();
    
    // Start typing animation
    startTypingAnimation();
    
    // Clean up on component destruction
    return () => {
      if (typingInterval) clearInterval(typingInterval);
    };
  });
  
  // Update line numbers based on content
  function updateLineNumbers() {
    const lines = editorContent.split('\n').length;
    lineNumbers = Array.from({ length: lines }, (_, i) => i + 1);
  }
  
  // Simulate typing animation
  function startTypingAnimation() {
    const textToType = `// Let's add a new function
void greet(const char* name) {
  printf("Welcome to Codex, %s!\\n", name);
}`;
    
    let typingPosition = editorContent.length;
    const insertPosition = editorContent.indexOf('return 0;');
    
    if (insertPosition !== -1) {
      editorContent = editorContent.slice(0, insertPosition) + '\n\n' + textToType + '\n\n  ' + editorContent.slice(insertPosition);
      updateLineNumbers();
      
      // Simulate cursor movement and typing
      isTyping = true;
      let charIndex = 0;
      
      typingInterval = setInterval(() => {
        charIndex++;
        cursorPosition = insertPosition + charIndex;
        
        if (charIndex >= textToType.length + 4) {
          clearInterval(typingInterval);
          setTimeout(() => {
            isTyping = false;
          }, 1000);
        }
      }, 50);
    }
  }
  
  // Handle tab switching
  function switchTab(tab: string) {
    activeTab = tab;
  }
  
  // Mock code execution
  function runCode() {
    isRunning = true;
    showOutput = false;
    
    // Simulate processing time
    setTimeout(() => {
      isRunning = false;
      showOutput = true;
      
      // Generate mock output based on the code
      if (editorContent.includes('printf("Hello, World!')) {
        outputContent = "Hello, World!\n";
      }
      
      if (editorContent.includes('greet')) {
        outputContent += "Welcome to Codex, Programmer!\n";
      }
      
      if (editorContent.includes('You\'re ready to begin!')) {
        outputContent += "You're ready to begin!\n";
      }
      
      // Random execution time between 50-200ms
      executionTime = Math.floor(Math.random() * 150) + 50;
      
      // 95% chance of success
      executionSuccess = Math.random() > 0.05;
      
      if (!executionSuccess) {
        outputContent = "Error: Compilation failed.\nLine 7: Syntax error - missing semicolon.";
      }
    }, 800);
  }
  
  // Copy code to clipboard
  function copyCode() {
    navigator.clipboard.writeText(editorContent);
  }
  
  // Reset code to initial state
  function resetCode() {
    editorContent = `#include <stdio.h>

int main() {
  // Welcome to Codex Learning Platform
  printf("Hello, World!\\n");
  
  // Your coding journey starts here
  int progress = 100;
  
  if (progress == 100) {
    printf("You're ready to begin!\\n");
  }
  
  return 0;
}`;
    updateLineNumbers();
    showOutput = false;
  }
</script>

<svelte:head>
  <title>Codex - Master Programming with Interactive Challenges</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Inter:wght@300;400;500;600;700&family=Orbitron:wght@400;500;600;700&display=swap" rel="stylesheet">
</svelte:head>

<div class="min-h-screen bg-gradient-to-b from-background to-background/95">
  <!-- Hero Section -->
  <section class="relative pt-24 pb-16 md:pt-32 md:pb-24 overflow-hidden">
    <!-- Background Elements -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[30%] -right-[10%] w-[70%] h-[70%] rounded-full bg-primary/5 blur-3xl"></div>
      <div class="absolute -bottom-[20%] -left-[10%] w-[50%] h-[50%] rounded-full bg-primary/5 blur-3xl"></div>
      
      <!-- Code Particles -->
      <div class="absolute top-[10%] left-[5%] text-primary/10 text-4xl font-mono">{`{}`}</div>
      <div class="absolute top-[20%] right-[15%] text-primary/10 text-4xl font-mono">{`</>`}</div>
      <div class="absolute bottom-[30%] left-[20%] text-primary/10 text-4xl font-mono">{`()`}</div>
      <div class="absolute bottom-[10%] right-[10%] text-primary/10 text-4xl font-mono">{`#`}</div>
    </div>
    
    <div class="container px-4 md:px-6">
      <div class="flex flex-col items-center text-center space-y-4 mb-12">
        <div in:fade={{ delay: 200, duration: 700 }} class="inline-flex items-center justify-center p-2 bg-primary/10 backdrop-blur-sm rounded-full mb-2">
          <Badge variant="outline" class="px-4 py-1 border-primary/20 text-primary bg-primary/5 backdrop-blur-sm">
            <Sparkles class="h-3.5 w-3.5 mr-1" />
            <span>Master Programming Interactively</span>
          </Badge>
        </div>
        
        <h1 in:fade={{ delay: 300, duration: 700 }} class="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight max-w-3xl">
          Learn to Code with <span class="text-primary">Interactive</span> Challenges
        </h1>
        
        <p in:fade={{ delay: 400, duration: 700 }} class="text-lg md:text-xl text-muted-foreground max-w-2xl">
          Elevate your programming skills through hands-on practice, real-time feedback, and a structured learning path.
        </p>
        
        <div in:fade={{ delay: 500, duration: 700 }} class="mt-4">
          <Button size="lg" class="gap-2 px-8 py-6 text-lg" onclick={() => window.location.href = '/signup'}>
            <Zap class="h-5 w-5" />
            Start Learning Now
          </Button>
        </div>
      </div>
      
      <!-- Enhanced Mock Editor -->
      <div in:fade={{ delay: 600, duration: 700 }} class="max-w-5xl mx-auto">
        <div class="rounded-xl overflow-hidden border shadow-xl bg-card/30 backdrop-blur-lg">
          <!-- Editor Header -->
          <div class="flex items-center justify-between px-4 py-3 bg-muted/50 border-b">
            <div class="flex items-center gap-2">
              <div class="flex space-x-2">
                <div class="h-3 w-3 rounded-full bg-red-500"></div>
                <div class="h-3 w-3 rounded-full bg-yellow-500"></div>
                <div class="h-3 w-3 rounded-full bg-green-500"></div>
              </div>
              <div class="ml-4 flex items-center">
                <Code class="h-4 w-4 mr-2 text-primary" />
                <span class="text-sm font-medium">main.c</span>
              </div>
            </div>
            
            <div class="flex items-center gap-2">
              <Button variant="ghost" size="icon" class="h-8 w-8" on:click={resetCode} title="Reset Code">
                <RefreshCw class="h-4 w-4" />
              </Button>
              
              <Button variant="ghost" size="icon" class="h-8 w-8" on:click={copyCode} title="Copy Code">
                <Copy class="h-4 w-4" />
              </Button>
              
              <Button variant="ghost" size="icon" class="h-8 w-8" title="Fullscreen">
                <Maximize2 class="h-4 w-4" />
              </Button>
            </div>
          </div>
          
          <!-- Editor Content -->
          <div class="relative font-mono text-sm overflow-hidden" style="height: 400px;">
            <div class="absolute inset-0 overflow-auto p-4 flex">
              <!-- Line Numbers -->
              <div class="select-none text-muted-foreground pr-4 text-right">
                {#each lineNumbers as lineNum}
                  <div class="h-6">{lineNum}</div>
                {/each}
              </div>
              
              <!-- Code Content -->
              <div class="flex-1 overflow-hidden relative" bind:this={editorElement}>
                <pre class="text-foreground"><code>{editorContent}</code></pre>
                
                {#if isTyping}
                  <div class="absolute top-0 left-0 right-0 bottom-0 pointer-events-none">
                    <div 
                      class="absolute w-0.5 h-6 bg-primary animate-pulse"
                      style={`top: ${Math.floor(cursorPosition / 50) * 24}px; left: ${(cursorPosition % 50) * 8}px`}
                    ></div>
                  </div>
                {/if}
              </div>
            </div>
            
            <!-- Editor Overlay Gradient -->
            <div class="absolute bottom-0 left-0 right-0 h-16 bg-gradient-to-t from-card/80 to-transparent pointer-events-none"></div>
          </div>
          
          <!-- Editor Footer -->
          <div class="px-4 py-3 bg-muted/30 border-t flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Badge variant="outline" class="bg-primary/10 text-primary border-primary/20">
                <Terminal class="h-3 w-3 mr-1" />
                Ready
              </Badge>
              <span class="text-xs text-muted-foreground">Syntax: OK</span>
            </div>
            
            <Button 
              size="sm" 
              class="gap-1 px-4 py-2 bg-green-600 hover:bg-green-700 text-white"
              on:click={runCode}
              disabled={isRunning}
            >
              {#if isRunning}
                <div class="h-3.5 w-3.5 rounded-full border-2 border-t-transparent border-white animate-spin mr-1"></div>
                Compiling...
              {:else}
                <Play class="h-3.5 w-3.5 mr-1" />
                Run Code
              {/if}
            </Button>
          </div>
          
          <!-- Output Panel (conditionally shown) -->
          {#if showOutput}
            <div in:slide={{ duration: 300 }} class="border-t">
              <div class="px-4 py-2 bg-muted/50 border-b flex items-center justify-between">
                <div class="flex items-center">
                  <Terminal class="h-4 w-4 mr-2 text-primary" />
                  <span class="text-sm font-medium">Output</span>
                </div>
                <div class="flex items-center gap-2">
                  {#if executionSuccess}
                    <Badge variant="outline" class="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 border-green-200 dark:border-green-800">
                      <CheckSquare class="h-3 w-3 mr-1" />
                      Success
                    </Badge>
                  {:else}
                    <Badge variant="outline" class="bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400 border-red-200 dark:border-red-800">
                      <XSquare class="h-3 w-3 mr-1" />
                      Error
                    </Badge>
                  {/if}
                  <Badge variant="outline" class="bg-muted/50">
                    <Clock class="h-3 w-3 mr-1" />
                    {executionTime}ms
                  </Badge>
                </div>
              </div>
              <div class="p-4 font-mono text-sm bg-muted/20 max-h-48 overflow-auto">
                <pre class={executionSuccess ? "text-foreground" : "text-red-500"}>{outputContent}</pre>
              </div>
            </div>
          {/if}
        </div>
        
        <!-- Editor Caption -->
        <div class="mt-4 text-center text-sm text-muted-foreground">
          <p>Our intelligent editor provides syntax highlighting, auto-completion, and real-time error checking</p>
        </div>
      </div>
    </div>
  </section>
  
  <!-- Improved Languages Section -->
  <section class="py-16 md:py-24 relative">
    <div class="absolute inset-0 bg-muted/30 -z-10"></div>
    <div class="absolute inset-0 overflow-hidden pointer-events-none -z-10">
      <div class="absolute top-[10%] right-[10%] w-[40%] h-[40%] rounded-full bg-primary/5 blur-3xl"></div>
      <div class="absolute bottom-[10%] left-[10%] w-[40%] h-[40%] rounded-full bg-primary/5 blur-3xl"></div>
    </div>
    
    <div class="container px-4 md:px-6">
      <div class="text-center mb-12">
        <h2 class="text-3xl md:text-4xl font-bold mb-4 font-['Orbitron'] text-primary">Master Multiple Programming Languages</h2>
        <p class="text-muted-foreground max-w-2xl mx-auto">
          Starting with C programming, we're expanding to cover all major languages. Build transferable skills across the programming ecosystem.
        </p>
      </div>
      
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-6 md:gap-8 max-w-5xl mx-auto">
        {#each languages as lang, i}
          <div in:fly={{ y: 20, delay: 100 + i * 50, duration: 400 }} class="relative group">
            <Card class="h-full bg-card/50 backdrop-blur-sm border-muted/50 group-hover:border-primary/20 transition-all overflow-hidden">
              <div class="absolute inset-0 opacity-5 group-hover:opacity-10 transition-opacity" style="background-color: {lang.color}"></div>
              <CardContent class="p-6 flex flex-col items-center justify-center text-center">
                <div class="h-16 w-16 rounded-full flex items-center justify-center mb-4 group-hover:scale-110 transition-transform" style="background-color: {lang.color}10">
                  <!-- Language Icon using simple-icons -->
                  <svg class="h-8 w-8" viewBox="0 0 24 24" fill="currentColor" style="color: {lang.color}">
                    <path d={lang.icon.path} />
                  </svg>
                </div>
                <h3 class="font-medium mb-1">{lang.name}</h3>
                {#if lang.coming}
                  <Badge variant="outline" class="text-xs bg-muted/50">Coming Soon</Badge>
                {:else}
                  <Badge variant="default" class="text-xs">Available Now</Badge>
                {/if}
              </CardContent>
            </Card>
          </div>
        {/each}
      </div>
      
      <!-- Tech Stack Section using simple-icons -->
   
    </div>
  </section>
  
  <!-- Testimonials Section -->
  <section class="py-16 md:py-24">
    <div class="container px-4 md:px-6">
      <div class="text-center mb-12">
        <h2 class="text-3xl md:text-4xl font-bold mb-4">What Our Students Say</h2>
        <p class="text-muted-foreground max-w-2xl mx-auto">
          Join thousands of learners who have transformed their coding skills with Codex.
        </p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 md:gap-8">
        {#each testimonials as testimonial, i}
          <div in:fly={{ y: 20, delay: 100 + i * 100, duration: 500 }}>
            <Card class="bg-card/50 backdrop-blur-sm border-muted/50">
              <CardContent class="p-6">
                <div class="flex items-center gap-4 mb-4">
                  <Avatar>
                    <AvatarImage src={testimonial.avatar || "/placeholder.svg"} alt={testimonial.name} />
                    <AvatarFallback>{testimonial.name[0]}{testimonial.name.split(' ')[1][0]}</AvatarFallback>
                  </Avatar>
                  <div>
                    <p class="font-medium">{testimonial.name}</p>
                    <p class="text-sm text-muted-foreground">{testimonial.role}</p>
                  </div>
                </div>
                <p class="italic text-muted-foreground">"{testimonial.content}"</p>
              </CardContent>
            </Card>
          </div>
        {/each}
      </div>
    </div>
  </section>
  
  <!-- CTA Section -->
  <section class="py-16 md:py-24 relative">
    <div class="absolute inset-0 bg-primary/5 -z-10"></div>
    <div class="absolute inset-0 overflow-hidden pointer-events-none -z-10">
      <div class="absolute -top-[10%] -right-[10%] w-[50%] h-[50%] rounded-full bg-primary/5 blur-3xl"></div>
      <div class="absolute -bottom-[10%] -left-[10%] w-[50%] h-[50%] rounded-full bg-primary/5 blur-3xl"></div>
    </div>
    
    <div class="container px-4 md:px-6">
      <Card class="bg-card/70 backdrop-blur-lg border-muted/50 overflow-hidden">
        <div class="absolute top-0 left-0 w-40 h-40 bg-primary/10 rounded-full -translate-y-1/2 -translate-x-1/2 blur-2xl pointer-events-none"></div>
        <div class="absolute bottom-0 right-0 w-40 h-40 bg-primary/10 rounded-full translate-y-1/2 translate-x-1/2 blur-2xl pointer-events-none"></div>
        
        <CardContent class="p-8 md:p-12">
          <div class="grid gap-8 md:grid-cols-2 items-center">
            <div>
              <h2 class="text-3xl md:text-4xl font-bold mb-4">Ready to Start Your Coding Journey?</h2>
              <p class="text-muted-foreground mb-6">
                Join our community of learners and start mastering programming skills that will advance your career.
              </p>
              <div>
                <Button size="lg" class="gap-2 px-8 py-6 text-lg">
                  <Zap class="h-5 w-5" />
                  Get Started 
                </Button>
              </div>
            </div>
            
            <div class="bg-muted/30 rounded-lg p-6 border border-muted/50">
              <h3 class="text-xl font-semibold mb-4 flex items-center">
                <CheckCircle class="h-5 w-5 text-primary mr-2" />
                What You'll Get
              </h3>
              <ul class="space-y-3">
                {#each [
                  "Access to 100+ interactive coding challenges",
                  "Real-time feedback on your code",
                  "Structured learning path from basics to advanced",
                  "Community support from fellow learners",
                  "Progress tracking and achievement system"
                ] as item}
                  <li class="flex items-start">
                    <CheckCircle class="h-5 w-5 text-primary shrink-0 mr-2 mt-0.5" />
                    <span>{item}</span>
                  </li>
                {/each}
              </ul>
              <div class="mt-6 pt-4 border-t border-muted/50 flex items-center justify-between">
                <div>
                  <p class="text-sm text-muted-foreground">Join 10,000+ developers</p>
                  <div class="flex -space-x-2 mt-1">
                    {#each Array(5) as _, i}
                      <Avatar class="h-6 w-6 border-2 border-background">
                        <AvatarImage src={`/placeholder.svg?height=24&width=24&text=${i+1}`} />
                        <AvatarFallback>U{i+1}</AvatarFallback>
                      </Avatar>
                    {/each}
                  </div>
                </div>
                <Button variant="ghost" class="gap-1 text-primary">
                  <span>View Success Stories</span>
                  <ArrowRight class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </section>
</div>

<style>
  :global(body) {
    font-family: 'Inter', sans-serif;
  }
  
  pre, code {
    font-family: 'JetBrains Mono', monospace;
  }
  
  /* Syntax highlighting */
  code .keyword { color: #ff7b72; }
  code .string { color: #a5d6ff; }
  code .comment { color: #8b949e; }
  code .number { color: #79c0ff; }
  code .function { color: #d2a8ff; }
  
  /* Custom scrollbar for code editor */
  :global(.overflow-auto::-webkit-scrollbar) {
    width: 8px;
    height: 8px;
  }
  
  :global(.overflow-auto::-webkit-scrollbar-track) {
    background: transparent;
  }
  
  :global(.overflow-auto::-webkit-scrollbar-thumb) {
    background: rgba(127, 127, 127, 0.2);
    border-radius: 4px;
  }
  
  :global(.overflow-auto::-webkit-scrollbar-thumb:hover) {
    background: rgba(127, 127, 127, 0.3);
  }
</style>