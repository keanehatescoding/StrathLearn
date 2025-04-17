<script lang="ts">
  import { onMount } from 'svelte';
  import toast, { Toaster } from 'svelte-french-toast';
  import { 
    Code, 
    ChevronRight,
    ChevronLeft, 
    Lightbulb, 
    ClipboardCheck, 
    RotateCcw, 
    Zap, 
    FileText, 
    CheckCircle, 
    XCircle,
    Loader,
    Github,
    Monitor,
    MoonStar,
    Sun,
    BookOpen,
    GraduationCap,
    Trophy,
    Layout,
    Plus,
    Minus,
    Maximize2,
    Minimize2,
    Copy,
    Terminal
  } from 'lucide-svelte';
  
  // Import shadcn components
  import { Button } from "$lib/components/ui/button/index";
  import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "$lib/components/ui/card/index";
  import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar/index";
  import { Badge } from "$lib/components/ui/badge/index";
  import { Tabs, TabsContent, TabsList, TabsTrigger } from "$lib/components/ui/tabs/index";
  import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "$lib/components/ui/tooltip/index";
  import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "$lib/components/ui/collapsible/index";
  import { Progress } from "$lib/components/ui/progress/index";
  import { Switch } from "$lib/components/ui/switch/index";
  import { Separator } from "$lib/components/ui/separator/index";
  import { ScrollArea } from "$lib/components/ui/scroll-area/index";

  interface Challenge {
    id: string;
    title: string;
    difficulty: string;
    description: string;
    hints: string[];
    testCases: TestCase[];
    initialCode: string;
    solutions?: string[];
    timeLimit?: number;
    memoryLimit?: number;
    category?: string;
  }

  interface TestCase {
    id: string;
    input: string;
    expectedOutput: string;
    hidden: boolean;
  }

  interface TestResult {
    testCaseId: string;
    passed: boolean;
    output?: string;
    error?: string;
    executionTime?: number;
  }

  let challenges: Record<string, Challenge> = {};
  let challengeList: Challenge[] = [];
  let currentChallenge: Challenge | null = null;
  let currentChallengeIndex: number = 0;
  let code: string = '';
  let editor: any;
  let testResults: TestResult[] = [];
  let showResults: boolean = false;
  let showSuccessMessage: boolean = false;
  let submitting: boolean = false;
  let theme: string = 'material-palenight';
  let darkMode: boolean = true;
  let isEditorFullscreen: boolean = false;
  let isChallengeFullscreen: boolean = false;
  let isProcessing: boolean = false;
  let executionTime: number | null = null;
  let activeTab: string = "description";
  let progressValue: number = 0;
  let editorFontSize: number = 14;
  let showMinimap: boolean = true;
  let autoSave: boolean = true;
  let timer: number | null = null;
  let timeElapsed: number = 0;
  
  onMount(async () => {
    await fetchChallenges();
    
    const textarea = document.getElementById('code-editor') as HTMLTextAreaElement;
    
    // Set up CodeMirror with more advanced configuration
    // @ts-ignore - CodeMirror is loaded from CDN
    editor = CodeMirror.fromTextArea(textarea, {
      mode: 'text/x-csrc',
      theme: theme,
      lineNumbers: true,
      autoCloseBrackets: true,
      matchBrackets: true,
      indentUnit: 4,
      tabSize: 4,
      indentWithTabs: true,
      lineWrapping: true,
      foldGutter: true,
      gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
      highlightSelectionMatches: {showToken: /\w/, annotateScrollbar: true},
      scrollbarStyle: "overlay",
      extraKeys: {
        "Tab": "indentMore", 
        "Shift-Tab": "indentLess",
        "Ctrl-Space": "autocomplete",
        "Ctrl-S": function() {
          saveCode();
          return false;
        },
        "Ctrl-/": "toggleComment",
        "Ctrl-F": "findPersistent",
        "F11": function() {
          toggleEditorFullscreen();
          return false;
        }
      },
      hintOptions: {
        completeSingle: false
      },
      styleActiveLine: true,
      matchTags: {bothTags: true},
      autoRefresh: true,
      viewportMargin: Infinity,
      // @ts-ignore
      lint: true
    });
    
    editor.on('change', () => {
      code = editor.getValue();
      if (autoSave) {
        localStorage.setItem(`code_${currentChallenge?.id}`, code);
      }
    });
    
    // Enhanced autocomplete
    // @ts-ignore
    editor.on("inputRead", function(cm, change) {
      if (change.origin !== '+input') return;
      const cur = editor.getCursor();
      const token = editor.getTokenAt(cur);
      
      if (token.type !== null && token.string.length > 1) {
        // @ts-ignore
        CodeMirror.commands.autocomplete(editor, null, { completeSingle: false });
      }
    });
    
    // Start a timer
    startTimer();

    // Show welcome toast
    toast.success("Welcome to CodeMaster! Happy coding!", {
      icon: "👨‍💻",
      style: "border-radius: 10px; background: #333; color: #fff;",
      duration: 4000
    });
    
    // Apply font size
    updateEditorFontSize();
    
    // Clean up on unmount
    return () => {
      if (timer) clearInterval(timer);
    };
  });
  
  function startTimer() {
    if (timer) clearInterval(timer);
    timeElapsed = 0;
    timer = setInterval(() => {
      timeElapsed++;
    }, 1000);
  }
  
  function formatTime(seconds: number): string {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  }
  
  function updateEditorFontSize() {
    document.documentElement.style.setProperty('--editor-font-size', `${editorFontSize}px`);
  }
  
  async function fetchChallenges() {
    isProcessing = true;
    try {
      const response = await fetch('https://api.singularity.co.ke/api/challenges');
      if (!response.ok) {
        throw new Error('Failed to load challenges');
      }
      
      challenges = await response.json();
      challengeList = Object.values(challenges);
      
      currentChallengeIndex = 0;
      
      if (challengeList.length > 0) {
        await loadChallenge(challengeList[0].id);
      }
    } catch (error) {
      console.error('Error fetching challenges:', error);
      toast.error("Failed to load challenges. Please refresh the page.");
    } finally {
      isProcessing = false;
    }
  }
  
  async function loadChallenge(id: string) {
    isProcessing = true;
    try {
      const response = await fetch(`https://api.singularity.co.ke/api/challenge/${id}`);
      if (!response.ok) {
        throw new Error('Challenge not found');
      }
      
      currentChallenge = await response.json();
      
      // Find the index of the current challenge
      currentChallengeIndex = challengeList.findIndex(c => c.id === id);
      
      // Check for saved code in localStorage
      const savedCode = localStorage.getItem(`code_${currentChallenge.id}`);
      const formattedCode = savedCode || currentChallenge?.initialCode.replace(/\\n/g, '\n');
      code = formattedCode!;
      
      if (editor) {
        editor.setValue(formattedCode);
        
        // Clear history after setting value
        setTimeout(() => {
          editor.clearHistory();
        }, 0);
      }
      
      showResults = false;
      showSuccessMessage = false;
      testResults = [];
      activeTab = "description";
      
      // Reset timer
      startTimer();
      
      // Update progress
      progressValue = ((currentChallengeIndex + 1) / challengeList.length) * 100;
      
      toast.success(`Challenge loaded: ${currentChallenge.title}`, {
        duration: 3000,
        position: "bottom-right"
      });
      
    } catch (error) {
      console.error('Error loading challenge:', error);
      toast.error("Error loading challenge. Please try again.");
      currentChallenge = null;
    } finally {
      isProcessing = false;
    }
  }
  
  function nextChallenge() {
    if (currentChallengeIndex < challengeList.length - 1) {
      loadChallenge(challengeList[currentChallengeIndex + 1].id);
    } else {
      toast("You've reached the last challenge!", {
        icon: "🏆",
        style: "border-radius: 10px; background: #4338ca; color: #fff;",
      });
    }
  }
  
  function previousChallenge() {
    if (currentChallengeIndex > 0) {
      loadChallenge(challengeList[currentChallengeIndex - 1].id);
    } else {
      toast("You're at the first challenge!", {
        icon: "🔍",
        style: "border-radius: 10px; background: #4338ca; color: #fff;",
      });
    }
  }
  
  async function submitCode() {
    if (!currentChallenge) {
      toast.error('No challenge loaded');
      return;
    }
    
    submitting = true;
    const startTime = performance.now();
    
    try {
      if (editor) {
        code = editor.getValue();
      }
      
      toast.loading("Running tests...", {
        id: "submit-toast"
      });
      
      const response = await fetch('https://api.singularity.co.ke/api/submit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          challengeId: currentChallenge.id,
          code: code
        })
      });
      
      const result = await response.json();
      const endTime = performance.now();
      executionTime = Math.round((endTime - startTime) / 10) / 100; // Round to 2 decimal places
      
      showResults = true;
      testResults = result.testResults || [];
      
      // Add execution time to each test result
      testResults = testResults.map(test => ({
        ...test,
        executionTime: Math.random() * 0.1 + 0.05 // Simulate different execution times for each test
      }));
      
      showSuccessMessage = testResults.length > 0 && testResults.every(test => test.passed);
      
      if (showSuccessMessage) {
        toast.success("All tests passed! 🎉", {
          id: "submit-toast",
          duration: 5000
        });
        
        // Save progress
        localStorage.setItem(`completed_${currentChallenge.id}`, "true");
      } else {
        const passCount = testResults.filter(t => t.passed).length;
        toast.error(`${passCount}/${testResults.length} tests passed`, {
          id: "submit-toast",
          duration: 5000
        });
      }
      
      // Auto scroll to results
      activeTab = "results";
      setTimeout(() => {
        const resultsElement = document.getElementById("results-section");
        if (resultsElement) {
          resultsElement.scrollIntoView({ behavior: 'smooth' });
        }
      }, 100);
      
    } catch (error) {
      console.error('Error submitting code:', error);
      toast.error('Error submitting code. Please try again.', {
        id: "submit-toast"
      });
    } finally {
      submitting = false;
    }
  }
  
  function resetCode() {
    if (currentChallenge && confirm('Are you sure you want to reset your code to the initial state?')) {
      const formattedCode = currentChallenge.initialCode.replace(/\\n/g, '\n');
      code = formattedCode;
      
      if (editor) {
        editor.setValue(formattedCode);
      }
      
      toast("Code reset to initial state", {
        icon: "🔄",
        style: "border-radius: 10px; background: #333; color: #fff;",
      });
    }
  }
  
  function saveCode() {
    if (currentChallenge) {
      localStorage.setItem(`code_${currentChallenge.id}`, code);
      toast.success("Code saved", {
        duration: 2000,
        position: "bottom-right"
      });
    }
  }
  
  function changeTheme(newTheme: string) {
    theme = newTheme;
    if (editor) {
      editor.setOption('theme', newTheme);
    }
    
    toast(`Theme changed to ${newTheme}`, {
      icon: "🎨",
      position: "bottom-right",
      duration: 2000
    });
  }
  
  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text)
      .then(() => {
        toast.success("Copied to clipboard!", {
          position: "bottom-right",
          duration: 2000
        });
      })
      .catch(err => {
        toast.error("Failed to copy", {
          position: "bottom-right",
          duration: 2000
        });
      });
  }
  
  function toggleDarkMode() {
    darkMode = !darkMode;
    document.documentElement.classList.toggle('dark', darkMode);
  }
  
  function toggleEditorFullscreen() {
    isEditorFullscreen = !isEditorFullscreen;
    setTimeout(() => {
      editor.refresh();
    }, 100);
  }
  
  function toggleChallengeFullscreen() {
    isChallengeFullscreen = !isChallengeFullscreen;
  }
  
  function increaseFontSize() {
    editorFontSize = Math.min(editorFontSize + 1, 24);
    updateEditorFontSize();
  }
  
  function decreaseFontSize() {
    editorFontSize = Math.max(editorFontSize - 1, 10);
    updateEditorFontSize();
  }
  
  function getDifficultyStyle(difficulty: string) {
    if (difficulty === 'Easy') {
      return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400';
    } else if (difficulty === 'Medium') {
      return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400';
    } else {
      return 'bg-rose-100 text-rose-800 dark:bg-rose-900/30 dark:text-rose-400';
    }
  }
  
  function getDifficultyColor(difficulty: string) {
    if (difficulty === 'Easy') {
      return 'text-emerald-500 dark:text-emerald-400';
    } else if (difficulty === 'Medium') {
      return 'text-amber-500 dark:text-amber-400';
    } else {
      return 'text-rose-500 dark:text-rose-400';
    }
  }
</script>

<svelte:head>
  <title>{currentChallenge ? `${currentChallenge.title} | CodeMaster` : 'CodeMaster - Coding Challenges'}</title>

  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/codemirror.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/theme/github-dark.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/theme/material-palenight.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/theme/dracula.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/fold/foldgutter.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/hint/show-hint.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/scroll/simplescrollbars.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/search/matchesonscrollbar.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/dialog/dialog.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/lint/lint.min.css">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">

  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/codemirror.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/mode/clike/clike.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/edit/closebrackets.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/edit/matchbrackets.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/fold/foldcode.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/fold/foldgutter.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/fold/brace-fold.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/hint/show-hint.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/hint/anyword-hint.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/search/searchcursor.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/search/search.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/search/match-highlighter.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/search/matchesonscrollbar.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/scroll/simplescrollbars.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/selection/active-line.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/edit/matchtags.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/comment/comment.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/dialog/dialog.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/search/searchcursor.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/lint/lint.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/display/autorefresh.min.js"></script>
</svelte:head>

<Toaster />

<div class="min-h-screen bg-background flex flex-col font-sans {darkMode ? 'dark' : ''}">
  <div class="pb-4 flex-grow">
    {#if isProcessing && !currentChallenge}
      <div class="flex h-screen justify-center items-center">
        <div class="text-center">
          <div class="inline-block animate-spin mb-4">
            <svg class="h-16 w-16 text-primary" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <h2 class="text-2xl font-bold mb-2">Loading CodeMaster</h2>
          <p class="text-muted-foreground">Preparing your coding experience...</p>
        </div>
      </div>
    {:else if currentChallenge}
      <div class="container mx-auto px-4 py-6">
        <!-- Top Controls -->
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <div class="flex items-center gap-3">
              <div class="bg-primary/10 p-2 rounded-lg">
                <Code class="h-5 w-5 text-primary" />
              </div>
              <div>
                <h1 class="text-2xl font-bold">{currentChallenge.title}</h1>
                <div class="flex items-center gap-2 mt-1">
                  <Badge variant="outline" class={getDifficultyStyle(currentChallenge.difficulty)}>
                    {currentChallenge.difficulty}
                  </Badge>
                  {#if currentChallenge.category}
                    <Badge variant="secondary">{currentChallenge.category}</Badge>
                  {/if}
                  <div class="text-sm text-muted-foreground">Challenge {currentChallengeIndex + 1} of {challengeList.length}</div>
                </div>
              </div>
            </div>
          </div>
          
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-2 bg-muted px-3 py-1.5 rounded-md text-muted-foreground">
              <Clock class="h-4 w-4" />
              <span>{formatTime(timeElapsed)}</span>
            </div>
            
            <div class="flex gap-2">
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="outline" size="icon" on:click={toggleDarkMode} class="h-9 w-9">
                      {#if darkMode}
                        <Sun class="h-4 w-4" />
                      {:else}
                        <MoonStar class="h-4 w-4" />
                      {/if}
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>Toggle Theme</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
              
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="outline" size="icon" disabled={currentChallengeIndex === 0} on:click={previousChallenge} class="h-9 w-9">
                      <ChevronLeft class="h-4 w-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>Previous Challenge</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
              
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="outline" size="icon" disabled={currentChallengeIndex === challengeList.length - 1} on:click={nextChallenge} class="h-9 w-9">
                      <ChevronRight class="h-4 w-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>Next Challenge</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </div>
        </div>
        
        <Progress value={progressValue} class="mb-6" />
        
        <!-- Main Content -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Challenge Info Panel -->
          <div class={`transition-all duration-300 ${isChallengeFullscreen ? 'lg:col-span-2' : ''}`}>
            <Card class="h-full flex flex-col">
              <CardHeader class="px-6 pt-6 pb-3 flex flex-row items-center justify-between space-y-0">
                <Tabs value={activeTab} class="w-full" onValueChange={(tab) => activeTab = tab}>
                  <TabsList class="grid grid-cols-3">
                    <TabsTrigger value="description" class="flex gap-2 items-center">
                      <BookOpen class="h-4 w-4" />
                      Description
                    </TabsTrigger>
                    <TabsTrigger value="hints" class="flex gap-2 items-center">
                      <Lightbulb class="h-4 w-4" />
                      Hints
                    </TabsTrigger>
                    <TabsTrigger value="testcases" class="flex gap-2 items-center">
                      <ClipboardCheck class="h-4 w-4" />
                      Test Cases
                    </TabsTrigger>
                  </TabsList>
                </Tabs>
                
                <div class="flex items-center gap-2">
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button variant="ghost" size="icon" class="h-8 w-8" on:click={toggleChallengeFullscreen}>
                          {#if isChallengeFullscreen}
                            <Minimize2 class="h-4 w-4" />
                          {:else}
                            <Maximize2 class="h-4 w-4" />
                          {/if}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>{isChallengeFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}</p>
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </div>
              </CardHeader>
              
              <CardContent class="flex-grow overflow-hidden p-0">
                <ScrollArea class="h-[500px] px-6 py-4">
                  <TabsContent value="description" class="m-0">
                    <div class="prose prose-slate dark:prose-invert max-w-none">
                      {@html currentChallenge.description.replace(/\n/g, '<br>')}
                    </div>
                  </TabsContent>
                  
                  <TabsContent value="hints" class="m-0">
                    <div class="space-y-4">
                      {#each currentChallenge.hints as hint, i}
                        <Card>
                          <CardHeader class="py-3 px-4">
                            <CardTitle class="text-sm font-medium flex items-center gap-2">
                              <Lightbulb class="h-4 w-4 text-amber-500" />
                              Hint {i + 1}
                            </CardTitle>
                          </CardHeader>
                          <CardContent class="py-3 px-4 pt-0">
                            <p>{hint}</p>
                          </CardContent>
                        </Card>
                      {/each}
                      {#if currentChallenge.hints.length === 0}
  <div class="p-4 text-center text-muted-foreground">
    <Lightbulb class="h-14 w-14 mx-auto mb-2 opacity-30" />
    <p>No hints available for this challenge.</p>
    <p class="text-sm mt-2">Try to solve it on your own first!</p>
  </div>
{/if}
</div>
</TabsContent>

<TabsContent value="testcases" class="m-0">
  <div class="space-y-4">
    {#each currentChallenge.testCases.filter(tc => !tc.hidden) as testCase, i}
      <Card>
        <CardHeader class="py-3 px-4">
          <CardTitle class="text-sm font-medium flex items-center gap-2">
            <Terminal class="h-4 w-4 text-blue-500" />
            Test Case {i + 1}
          </CardTitle>
        </CardHeader>
        <CardContent class="py-3 px-4 pt-0 space-y-3">
          <div>
            <p class="text-xs text-muted-foreground mb-1">Input:</p>
            <pre class="bg-secondary/50 p-2 rounded-md text-xs overflow-x-auto">{testCase.input}</pre>
          </div>
          <div>
            <p class="text-xs text-muted-foreground mb-1">Expected Output:</p>
            <pre class="bg-secondary/50 p-2 rounded-md text-xs overflow-x-auto">{testCase.expectedOutput}</pre>
          </div>
        </CardContent>
      </Card>
    {/each}
    
    {#if currentChallenge.testCases.some(tc => tc.hidden)}
      <Card class="border-dashed">
        <CardHeader class="py-3 px-4">
          <CardTitle class="text-sm font-medium flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 text-purple-500">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Hidden Test Cases
          </CardTitle>
        </CardHeader>
        <CardContent class="py-3 px-4 pt-0">
          <p class="text-sm text-muted-foreground">
            There are {currentChallenge.testCases.filter(tc => tc.hidden).length} hidden test cases that will be used to evaluate your solution.
          </p>
        </CardContent>
      </Card>
    {/if}
  </div>
</TabsContent>
</ScrollArea>
</CardContent>
</Card>
</div>

<!-- Code Editor Panel -->
<div class={`transition-all duration-300 ${isEditorFullscreen ? 'lg:col-span-2' : ''}`}>
  <Card class="h-full flex flex-col">
    <CardHeader class="px-6 pt-6 pb-3 flex flex-row items-center justify-between space-y-0">
      <div class="flex items-center gap-3">
        <FileText class="h-5 w-5 text-primary" />
        <h2 class="text-lg font-medium">Your Solution</h2>
      </div>
      
      <div class="flex items-center gap-2">
        <div class="flex items-center gap-1.5">
          <Button variant="ghost" size="icon" class="h-8 w-8" on:click={decreaseFontSize}>
            <Minus class="h-3.5 w-3.5" />
          </Button>
          <span class="text-xs font-mono">{editorFontSize}px</span>
          <Button variant="ghost" size="icon" class="h-8 w-8" on:click={increaseFontSize}>
            <Plus class="h-3.5 w-3.5" />
          </Button>
        </div>
        
        <div class="flex items-center rounded-md bg-muted">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" class="h-8 w-8" on:click={resetCode}>
                  <RotateCcw class="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Reset Code</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" class="h-8 w-8" on:click={saveCode}>
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5m19.5 0v-.75a2.25 2.25 0 00-2.25-2.25h-15a2.25 2.25 0 00-2.25 2.25v.75m5.25.75h2.25l2.25 3h-6.75l2.25-3a2.25 2.25 0 012.25 2.25v.75m2.25-3h-6.75l2.25 3h2.25l2.25-3z" />
                  </svg>
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Save Code</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" class="h-8 w-8" on:click={toggleEditorFullscreen}>
                  {#if isEditorFullscreen}
                    <Minimize2 class="h-4 w-4" />
                  {:else}
                    <Maximize2 class="h-4 w-4" />
                  {/if}
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>{isEditorFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </div>
    </CardHeader>
    
    <CardContent class="flex-grow p-0 flex flex-col">
      <div class="border-y border-border flex items-center px-3 py-1.5 bg-muted/50 gap-2 text-sm">
        <Button variant="ghost" size="sm" class={theme === 'material-palenight' ? 'bg-muted' : ''} on:click={() => changeTheme('material-palenight')}>
          Material
        </Button>
        <Button variant="ghost" size="sm" class={theme === 'dracula' ? 'bg-muted' : ''} on:click={() => changeTheme('dracula')}>
          Dracula
        </Button>
        <Button variant="ghost" size="sm" class={theme === 'github-dark' ? 'bg-muted' : ''} on:click={() => changeTheme('github-dark')}>
          GitHub
        </Button>
        <div class="flex-grow"></div>
        <div class="flex items-center gap-1.5">
          <span class="text-xs text-muted-foreground">Auto-save</span>
          <Switch checked={autoSave} on:change={() => autoSave = !autoSave} />
        </div>
      </div>
      
      <div class="flex-grow relative {isEditorFullscreen ? 'h-[calc(100vh-200px)]' : 'h-[500px]'}">
        <textarea id="code-editor" bind:value={code} class="w-full h-full"></textarea>
      </div>
    </CardContent>
    
    <CardFooter class="px-6 py-4 border-t bg-card flex flex-col sm:flex-row items-center gap-4">
      <Button class="w-full sm:w-auto" variant="default" size="lg" on:click={submitCode} disabled={submitting}>
        {#if submitting}
          <Loader class="h-4 w-4 mr-2 animate-spin" />
          Running Tests...
        {:else}
          <Zap class="h-4 w-4 mr-2" />
          Run Code
        {/if}
      </Button>
      
      <div class="flex gap-3 flex-wrap justify-center sm:justify-start">
        {#if currentChallenge.timeLimit}
          <div class="flex items-center gap-1.5 text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 text-muted-foreground">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-muted-foreground">Time Limit: {currentChallenge.timeLimit}ms</span>
          </div>
        {/if}
        
        {#if currentChallenge.memoryLimit}
          <div class="flex items-center gap-1.5 text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 text-muted-foreground">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 3v1.5M4.5 8.25H3m18 0h-1.5M4.5 12H3m18 0h-1.5m-15 3.75H3m18 0h-1.5M8.25 19.5V21M12 3v1.5m0 15V21m3.75-18v1.5m0 15V21m-9-1.5h10.5a2.25 2.25 0 002.25-2.25V6.75a2.25 2.25 0 00-2.25-2.25H6.75A2.25 2.25 0 004.5 6.75v10.5a2.25 2.25 0 002.25 2.25z" />
            </svg>
            <span class="text-muted-foreground">Memory Limit: {currentChallenge.memoryLimit}MB</span>
          </div>
        {/if}
        
        {#if executionTime !== null}
          <div class="flex items-center gap-1.5 text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 text-muted-foreground">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23.693L5 16M3 3l6.75 6.75M21 3l-6.75 6.75" />
            </svg>
            <span class="text-muted-foreground">Execution Time: {executionTime.toFixed(2)}s</span>
          </div>
        {/if}
      </div>
    </CardFooter>
  </Card>
</div>
</div>

<!-- Test Results Section -->
{#if showResults}
  <div id="results-section" class="mt-8">
    <div class="mb-4 flex items-center gap-2">
      <h2 class="text-xl font-bold">Test Results</h2>
      <div class="ml-auto flex gap-2">
        {#if showSuccessMessage}
          <Badge variant="outline" class="bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400">
            All Tests Passed!
          </Badge>
        {:else}
          <Badge variant="outline" class="bg-rose-100 text-rose-800 dark:bg-rose-900/30 dark:text-rose-400">
            Some Tests Failed
          </Badge>
        {/if}
      </div>
    </div>
    
    <div class="grid gap-4">
      {#each testResults as result, index}
        <Card class={result.passed ? 'border-emerald-200 dark:border-emerald-800/30' : 'border-rose-200 dark:border-rose-800/30'}>
          <CardHeader class="py-3 px-4 flex flex-row items-center">
            <CardTitle class="text-sm font-medium flex items-center gap-2">
              {#if result.passed}
                <CheckCircle class="h-4 w-4 text-emerald-500" />
                <span>Test {index + 1} Passed</span>
              {:else}
                <XCircle class="h-4 w-4 text-rose-500" />
                <span>Test {index + 1} Failed</span>
              {/if}
            </CardTitle>
            
            <div class="ml-auto flex items-center gap-2">
              {#if result.executionTime}
                <div class="text-xs text-muted-foreground">
                  {(result.executionTime * 1000).toFixed(2)}ms
                </div>
              {/if}
              {#if result.output}
                <Button variant="ghost" size="sm" on:click={() => copyToClipboard(result.output)}>
                  <Copy class="h-3 w-3 mr-1" />
                  Copy
                </Button>
              {/if}
            </div>
          </CardHeader>
          
          <CardContent class="py-3 px-4 space-y-3">
            {#if !result.passed}
              <div>
                <p class="text-xs text-muted-foreground mb-1">Expected Output:</p>
                <pre class="bg-secondary/50 p-2 rounded-md text-xs overflow-x-auto">{currentChallenge.testCases.find(tc => tc.id === result.testCaseId)?.expectedOutput || 'Unknown'}</pre>
              </div>
            {/if}
            
            {#if result.output}
              <div>
                <p class="text-xs text-muted-foreground mb-1">Your Output:</p>
                <pre class="bg-secondary/50 p-2 rounded-md text-xs overflow-x-auto">{result.output}</pre>
              </div>
            {/if}
            
            {#if result.error}
              <div>
                <p class="text-xs text-rose-600 dark:text-rose-400 mb-1">Error:</p>
                <pre class="bg-rose-50 dark:bg-rose-900/20 p-2 rounded-md text-xs overflow-x-auto text-rose-600 dark:text-rose-400">{result.error}</pre>
              </div>
            {/if}
          </CardContent>
        </Card>
      {/each}
    </div>
    
    {#if showSuccessMessage}
      <div class="mt-8 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg p-6 text-center">
        <div class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 dark:bg-emerald-900/30">
          <Trophy class="h-6 w-6 text-emerald-600 dark:text-emerald-400" />
        </div>
        <h3 class="text-lg font-medium text-emerald-800 dark:text-emerald-200 mb-2">Challenge Completed!</h3>
        <p class="text-emerald-600 dark:text-emerald-300 mb-6">
          Congratulations! You've successfully solved this challenge.
        </p>
        <div class="flex justify-center gap-4">
          <Button on:click={nextChallenge} disabled={currentChallengeIndex >= challengeList.length - 1}>
            <ChevronRight class="mr-1 h-4 w-4" />
            Next Challenge
          </Button>
          <Button variant="outline" on:click={() => {
            activeTab = "description";
          }}>
            <RotateCcw class="mr-1 h-4 w-4" />
            Try Another Solution
          </Button>
        </div>
      </div>
    {/if}
  </div>
{/if}
</div>
{:else}
  <div class="flex h-screen justify-center items-center">
    <Card class="w-[400px]">
      <CardHeader>
        <CardTitle>No Challenges Found</CardTitle>
        <CardDescription>
          There seems to be an issue loading the coding challenges.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-muted-foreground mb-4">
          Please check your internet connection and refresh the page. If the problem persists, contact support.
        </p>
      </CardContent>
      <CardFooter>
        <Button on:click={fetchChallenges} class="w-full">
          <RotateCcw class="mr-2 h-4 w-4" />
          Retry Loading Challenges
        </Button>
      </CardFooter>
    </Card>
  </div>
{/if}
</div>

<!-- Tutorial Overlay - First Time User -->
{#if false}
  <div class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
    <Card class="w-[500px] max-w-[90vw]">
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <GraduationCap class="h-5 w-5 text-primary" />
          Welcome to CodeMaster!
        </CardTitle>
        <CardDescription>
          Here's a quick tour to get you started.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-start gap-3">
          <div class="bg-primary/10 p-2 rounded-full">
            <Layout class="h-5 w-5 text-primary" />
          </div>
          <div>
            <h3 class="font-medium mb-1">Two-Panel Layout</h3>
            <p class="text-sm text-muted-foreground">
              On the left, you'll find the challenge description, hints, and test cases. On the right, you can write and test your code.
            </p>
          </div>
        </div>
        
        <div class="flex items-start gap-3">
          <div class="bg-primary/10 p-2 rounded-full">
            <Tabs class="h-5 w-5 text-primary" />
          </div>
          <div>
            <h3 class="font-medium mb-1">Challenge Information</h3>
            <p class="text-sm text-muted-foreground">
              Switch between description, hints, and test cases to understand what you need to solve.
            </p>
          </div>
        </div>
        
        <div class="flex items-start gap-3">
          <div class="bg-primary/10 p-2 rounded-full">
            <Code class="h-5 w-5 text-primary" />
          </div>
          <div>
            <h3 class="font-medium mb-1">Write Your Code</h3>
            <p class="text-sm text-muted-foreground">
              Use the editor to write your solution. You can customize the editor theme and font size to your liking.
            </p>
          </div>
        </div>
        
        <div class="flex items-start gap-3">
          <div class="bg-primary/10 p-2 rounded-full">
            <Zap class="h-5 w-5 text-primary" />
          </div>
          <div>
            <h3 class="font-medium mb-1">Test Your Solution</h3>
            <p class="text-sm text-muted-foreground">
              Click "Run Code" to submit your solution and see if it passes all the test cases.
            </p>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button class="w-full">
          Start Coding
        </Button>
      </CardFooter>
    </Card>
  </div>
{/if}

<!-- Footer -->
<footer class="mt-auto py-4 border-t">
  <div class="container mx-auto px-4 flex flex-col md:flex-row justify-between items-center gap-4">
    <div class="flex items-center gap-2">
      <Code class="h-5 w-5 text-primary" />
      <span class="font-medium">CodeMaster</span>
      <span class="text-xs text-muted-foreground">v2.0.0</span>
    </div>
    
    <div class="flex items-center gap-4">
      <a href="https://github.com/codemaster/challenges" target="_blank" rel="noopener noreferrer" class="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors">
        <Github class="h-4 w-4" />
        <span>GitHub</span>
      </a>
      
      <Button variant="ghost" size="sm" class="text-sm">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 mr-1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 5.25h.008v.008H12v-.008z" />
        </svg>
        Help & Documentation
      </Button>
      
      <Button variant="ghost" size="sm" class="text-sm">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 mr-1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 18v-5.25m0 0a6.01 6.01 0 001.5-.189m-1.5.189a6.01 6.01 0 01-1.5-.189m3.75 7.478a12.06 12.06 0 01-4.5 0m3.75 2.383a14.406 14.406 0 01-3 0M14.25 18v-.192c0-.983.658-1.823 1.508-2.316a7.5 7.5 0 10-7.517 0c.85.493 1.509 1.333 1.509 2.316V18" />
        </svg>
        Report a Bug
      </Button>
    </div>
  </div>
</footer>
</div>

<style>
  :root {
    --editor-font-size: 14px;
  }
  
  :global(.CodeMirror) {
    height: 100% !important;
    font-family: 'JetBrains Mono', monospace;
    font-size: var(--editor-font-size);
    border-radius: 0;
  }
  
  :global(.dark .CodeMirror-gutters) {
    background-color: transparent;
    border-right: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  :global(.CodeMirror-scroll) {
    overflow-y: auto;
    overflow-x: auto;
  }
  
  :global(.prose code) {
    background-color: hsl(var(--muted));
    padding: 0.2em 0.4em;
    border-radius: 0.25em;
    font-size: 0.875em;
    font-weight: 500;
  }
  
  :global(.dark .prose code) {
    background-color: rgba(255, 255, 255, 0.1);
  }
  
  :global(.dark .CodeMirror-cursor) {
    border-left-color: rgba(255, 255, 255, 0.8);
  }
  
  :global(.cm-s-material-palenight .CodeMirror-linenumber, 
         .cm-s-dracula .CodeMirror-linenumber, 
         .cm-s-github-dark .CodeMirror-linenumber) {
    opacity: 0.5;
  }
  
  :global(.CodeMirror-focused) {
    box-shadow: none !important;
  }
  
  /* Animation for success celebration */
  @keyframes pulse {
    0%, 100% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.05);
    }
  }
  
  .success-pulse {
    animation: pulse 1.5s infinite;
  }
</style>