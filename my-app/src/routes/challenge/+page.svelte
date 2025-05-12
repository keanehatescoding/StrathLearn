<script lang="ts">
  import { onMount } from 'svelte';
  import { toast } from 'svelte-sonner';
  import { 
    Code, ChevronRight, ChevronLeft, Terminal, 
    RotateCcw, Zap, Info, ArrowRight,  
    Maximize2, Minimize2, Loader, Keyboard,
    CheckCircle, XCircle, Clock, Award
  } from 'lucide-svelte';
  
  // Import UI components
  import { Button } from "$lib/components/ui/button/index.js";
  import { Avatar, AvatarFallback } from "$lib/components/ui/avatar";
  import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "$lib/components/ui/card";
  import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "$lib/components/ui/select";
  import { Badge } from "$lib/components/ui/badge";
  import { Alert, AlertDescription } from "$lib/components/ui/alert";
  import { Separator } from "$lib/components/ui/separator";
  import { Progress } from "$lib/components/ui/progress";
  import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "$lib/components/ui/accordion";
  import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "$lib/components/ui/collapsible";
  import ModeToggle from "$lib/components/mode-toggle.svelte";
  import * as Tooltip from "$lib/components/ui/tooltip";

  // Type definitions
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
    points?: number;
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


// Then use it


  // State variables
  let challenges: Record<string, Challenge> = {};
  let challengeIds: string[] = [];
  let currentChallenge: Challenge | null = null;
  let currentChallengeIndex: number = 0;
  let code: string = '';
  let editor: any;
  let testResults: TestResult[] = [];
  let showResults: boolean = false;
  let submitting: boolean = false;
  let editorTheme: string = 'github-dark';
  let isFullScreen: boolean = false;
  let progress: number = 0;
  let editorHeight: string = '100%';
  let showHints: boolean = false;
  
  // Editor themes
  const editorThemes = [
    { value: 'vs-dark', label: 'VS Dark' },
    { value: 'github-dark', label: 'GitHub Dark' },
    { value: 'dracula', label: 'Dracula' },
    { value: 'nord', label: 'Nord' },
    { value: 'ayu-dark', label: 'Ayu Dark' },
    { value: 'tomorrow-night', label: 'Tomorrow Night' }
  ];

  // Initialization on component mount
  onMount(async () => {
    await fetchChallenges();
    
    const textarea = document.getElementById('code-editor');
    
    if (!textarea) {
      console.error("Editor element not found");
      toast.error("Failed to initialize editor");
      return;
    }
    

    
    // Wait for Monaco to load
    if (typeof monaco === 'undefined') {
      window.require(['vs/editor/editor.main'], initEditor);

    } else {
      initEditor();

     
    }

    // Add keyboard listeners
    window.addEventListener('keydown', handleKeyboardShortcuts);
    
    // Initial progress calculation
    updateProgress();
    
    // Clean up on component destruction
    return () => {
      window.removeEventListener('keydown', handleKeyboardShortcuts);
      if (editor) editor.dispose();
    };
  });
  
  // Initialize Monaco editor
  function initEditor() {
    try {
      monaco.editor.defineTheme('github-dark', {
      base: 'vs-dark',
      inherit: true,
      rules: [
        { token: 'comment', foreground: '6a737d' },
        { token: 'keyword', foreground: 'ff7b72' },
        { token: 'string', foreground: 'a5d6ff' },
        { token: 'number', foreground: '79c0ff' },
        { token: 'regexp', foreground: 'a5d6ff' },
        { token: 'operator', foreground: 'ff7b72' },
        { token: 'namespace', foreground: 'ff7b72' },
        { token: 'type.identifier', foreground: '79c0ff' },
        { token: 'identifier', foreground: 'c9d1d9' },
        { token: 'variable', foreground: 'ffa657' },
        { token: 'variable.predefined', foreground: '79c0ff' },
        { token: 'function', foreground: 'd2a8ff' },
      ],
      colors: {
        'editor.background': '#0d1117',
        'editor.foreground': '#c9d1d9',
        'editorCursor.foreground': '#c9d1d9',
        'editor.lineHighlightBackground': '#161b22',
        'editorLineNumber.foreground': '#6e7681',
        'editor.selectionBackground': '#3b5070',
        'editor.inactiveSelectionBackground': '#282e33'
      }
    });
      editor = monaco.editor.create(document.getElementById('code-editor'), {
        value: code,
        language: 'c',
        theme: editorTheme,
        automaticLayout: true,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        fontSize: 14,
        fontFamily: '"JetBrains Mono", Menlo, Monaco, "Courier New", monospace',
        lineNumbers: 'on',
        renderLineHighlight: 'all',
        wordWrap: 'on',
        tabSize: 4,
        glyphMargin: true,
        bracketPairColorization: { enabled: true },
        "semanticHighlighting.enabled": true,
        formatOnPaste: true,
        formatOnType: true
      });
      
      
      editor.onDidChangeModelContent(() => {
        code = editor.getValue();
      });
      
      // Ensure editor layouts properly when resized
      const resizeObserver = new ResizeObserver(() => {
        if (editor) editor.layout();
      });
      
      resizeObserver.observe(document.getElementById('code-editor'));
      
    } catch (error) {
      console.error("Error initializing editor:", error);
      toast.error("Failed to initialize code editor. Please refresh the page.");
    }
  }

  
  // Handle keyboard shortcuts
  function handleKeyboardShortcuts(e: KeyboardEvent) {
    // Ctrl+Enter or Cmd+Enter to submit
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      e.preventDefault();
      submitCode();
    }
    
    // Ctrl+S or Cmd+S to save
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
      e.preventDefault();
      saveCodeToLocalStorage();
      toast.success('Code saved');
    }
    
    // F11 for fullscreen
    if (e.key === 'F11') {
      e.preventDefault();
      toggleFullScreen();
    }
  }
  
  // Update progress indicator
  function updateProgress() {
    if (!challengeIds.length) return;
    progress = ((currentChallengeIndex + 1) / challengeIds.length) * 100;
  }
  
  // Fetch all challenges from API
  async function fetchChallenges() {
    try {
      // Show loading toast
      const toastId = toast.loading('Loading challenges...');
      
      const response = await fetch('https://api.singularity.co.ke/api/challenges');
      if (!response.ok) {
        throw new Error('Failed to load challenges');
      }
      
      challenges = await response.json();
      challengeIds = Object.keys(challenges);
      
      if (challengeIds.length > 0) {
        currentChallengeIndex = 0;
        await loadChallenge(challengeIds[currentChallengeIndex]);
        toast.dismiss(toastId);
        toast.success('Challenges loaded successfully');
      } else {
        toast.dismiss(toastId);
        toast.error('No challenges found');
      }
    } catch (error) {
      console.error('Error fetching challenges:', error);
      toast.error('Failed to load challenges');
    }
  }
  


  async function loadChallenge(id: string) {
  try {
    const toastId = toast.loading('Loading challenge...');
    

    const tokenResponse = await fetch('https://codex.singularity.co.ke/api/auth/token');
    if (!tokenResponse.ok) {
      throw new Error('Failed to retrieve authentication token');
    }
    
    const { token } = await tokenResponse.json();
    
    const response = await fetch(`https://api.singularity.co.ke/api/challenge/${id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Challenge not found');
    }
    
    currentChallenge = await response.json();
    
    // Format the code and set it in the editor
    const formattedCode = currentChallenge?.initialCode.replace(/\\n/g, '\n');
    code = formattedCode!;
    
    if (editor) {
      // Set language to C
      monaco.editor.setModelLanguage(editor.getModel(), 'c');
      editor.setValue(formattedCode);
    }
    
    
    const savedCode = localStorage.getItem(`c_master_challenge_${id}`);
    if (savedCode) {
      code = savedCode;
      if (editor) {
        editor.setValue(savedCode);
      }
    }
    
    
    showResults = false;
    testResults = [];
    showHints = false;
    
    updateProgress();
    
    toast.dismiss(toastId);
  } catch (error) {
    console.error('Error loading challenge:', error);
    toast.error('Error loading challenge');
    currentChallenge = null;
  }
}


async function submitCode() {
  if (!currentChallenge) {
    toast.error('No challenge loaded');
    return;
  }

  submitting = true;

  try {
    if (editor) {
      code = editor.getValue();
    }

    const tokenResponse = await fetch('https://codex.singularity.co.ke/api/auth/token');
    if (!tokenResponse.ok) {
      throw new Error('Failed to retrieve authentication token');
    }
    
    const { token } = await tokenResponse.json();

    const response = await fetch('https://api.singularity.co.ke/api/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        challengeId: currentChallenge.id,
        code: code
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    showResults = true;
    testResults = result.testResults || [];

    // Scroll to results
    setTimeout(() => {
      const resultsSection = document.getElementById('results-section');
      if (resultsSection) {
        resultsSection.scrollIntoView({ behavior: 'smooth' });
      }
    }, 100);

  } catch (error) {
    console.error('Error submitting code:', error);
    toast.error('Error submitting code. Please try again.');
  } finally {
    submitting = false;
  }
}
  
  // Reset code to initial state
  function resetCode() {
    if (currentChallenge && confirm('Are you sure you want to reset your code to the initial state?')) {
      const formattedCode = currentChallenge.initialCode.replace(/\\n/g, '\n');
      code = formattedCode;
      
      if (editor) {
        editor.setValue(formattedCode);
      }
      
      toast.info('Code has been reset to initial state');
    }
  }
  
  // Save code to local storage
  function saveCodeToLocalStorage() {
    if (currentChallenge) {
      localStorage.setItem(`c_master_challenge_${currentChallenge.id}`, code);
    }
  }
  
  // Change editor theme
  function changeEditorTheme(newTheme: string) {
    editorTheme = newTheme;
    if (editor) {
      monaco.editor.setTheme(newTheme);
      toast.success(`Theme changed to ${newTheme}`);
    }
  }
  
  // Navigate between challenges
  function navigateToChallenge(direction: 'prev' | 'next') {
    saveCodeToLocalStorage();
    
    if (direction === 'next' && currentChallengeIndex < challengeIds.length - 1) {
      currentChallengeIndex++;
      loadChallenge(challengeIds[currentChallengeIndex]);
    } else if (direction === 'prev' && currentChallengeIndex > 0) {
      currentChallengeIndex--;
      loadChallenge(challengeIds[currentChallengeIndex]);
    }
  }
  
  // Toggle fullscreen mode
  function toggleFullScreen() {
    isFullScreen = !isFullScreen;
    
    const appContainer = document.getElementById('c-master-app');
    
    if (isFullScreen) {
      // Enter fullscreen
      if (appContainer?.requestFullscreen) {
        appContainer.requestFullscreen();
      }
    } else {
      // Exit fullscreen
      if (document.exitFullscreen) {
        document.exitFullscreen();
      }
    }
    
    // Ensure editor resizes properly
    setTimeout(() => {
      if (editor) editor.layout();
    }, 100);
  }
  
  // Get appropriate color for difficulty badge
  function getDifficultyColor(difficulty: string) {
    switch(difficulty.toLowerCase()) {
      case 'easy': return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900 dark:text-emerald-300';
      case 'medium': return 'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-300';
      case 'hard': return 'bg-rose-100 text-rose-800 dark:bg-rose-900 dark:text-rose-300';
      default: return 'bg-slate-100 text-slate-800 dark:bg-slate-800 dark:text-slate-300';
    }
  }
</script>

<svelte:head>
  <title>Codex</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
  
  <!-- Monaco Editor CDN -->
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.44.0/min/vs/editor/editor.main.min.css">
  <script src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.44.0/min/vs/loader.min.js"></script>
  <script>
    require.config({ paths: { 'vs': 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.44.0/min/vs' } });
    window.MonacoEnvironment = {
      getWorkerUrl: function(workerId, label) {
        return `data:text/javascript;charset=utf-8,${encodeURIComponent(`
          self.MonacoEnvironment = {
            baseUrl: 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.44.0/min/'
          };
          importScripts('https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.44.0/min/vs/base/worker/workerMain.js');`
        )}`;
      }
    };
    require(['vs/editor/editor.main']);
  </script>
</svelte:head>

<div id="c-master-app" class="h-screen w-full flex flex-col bg-background">

  <!-- Main content with LeetCode-style layout -->
  <div class="flex-1 flex overflow-hidden">
    {#if currentChallenge}
      <!-- Problem description panel (left) -->
      <div class="w-[450px] flex flex-col border-r overflow-hidden">
        <!-- Challenge navigation -->
        <div class="border-b px-4 py-3 flex items-center justify-between bg-muted/30">
          <div class="flex items-center gap-2">
            <Button variant="outline" size="sm" disabled={currentChallengeIndex === 0} onclick={() => navigateToChallenge('prev')}>
              <ChevronLeft class="h-4 w-4 mr-1" />
              Previous
            </Button>
            
            <Button variant="outline" size="sm" disabled={currentChallengeIndex === challengeIds.length - 1} onclick={() => navigateToChallenge('next')}>
              Next
              <ChevronRight class="h-4 w-4 ml-1" />
            </Button>
          </div>
          
          <div class="flex items-center">
            <span class="text-sm font-medium">
              {currentChallengeIndex + 1}/{challengeIds.length}
            </span>
            <Progress value={progress} class="w-20 h-1.5 ml-2" />
          </div>
        </div>
        
        <!-- Challenge details -->
        <div class="flex-1 overflow-y-auto p-4">
          <div class="space-y-6">
            <!-- Challenge title and metadata -->
            <div>
              <h2 class="text-2xl font-bold mb-2">{currentChallenge.title}</h2>
              
              <div class="flex flex-wrap gap-2 mb-4">
                <Badge class={getDifficultyColor(currentChallenge.difficulty)}>
                  {currentChallenge.difficulty}
                </Badge>
                
                {#if currentChallenge.category}
                  <Badge variant="outline">{currentChallenge.category}</Badge>
                {/if}
                
                {#if currentChallenge.points}
                  <Badge variant="secondary">
                    <Award class="h-3.5 w-3.5 mr-1" />
                    {currentChallenge.points} pts
                  </Badge>
                {/if}
                
                {#if currentChallenge.timeLimit}
                  <Badge variant="outline">
                    <Clock class="h-3.5 w-3.5 mr-1" />
                    {currentChallenge.timeLimit}ms
                  </Badge>
                {/if}
              </div>
            </div>
            
            <!-- Challenge description -->
            <div class="prose dark:prose-invert max-w-none">
              {@html currentChallenge.description.replace(/\n/g, '<br>')}
            </div>
            
            <Separator />
            
            <!-- Test cases -->
            <div>
              <h3 class="text-lg font-semibold mb-3 flex items-center">
                <Terminal class="h-4 w-4 mr-2" />
                Test Cases
              </h3>
              
              <div class="space-y-4">
                {#each currentChallenge.testCases.filter(tc => !tc.hidden) as testCase, i}
                  <Card>
                    <CardHeader class="py-2 px-4">
                      <CardTitle class="text-sm">Example {i + 1}</CardTitle>
                    </CardHeader>
                    <CardContent class="py-3 px-4">
                      <div class="space-y-3">
                        <div>
                          <p class="text-sm font-medium mb-1">Input:</p>
                          <div class="font-mono text-xs bg-muted/50 p-2 rounded overflow-x-auto">
                            {testCase.input ? testCase.input : '(empty)'}
                          </div>
                        </div>
                        <div>
                          <p class="text-sm font-medium mb-1">Expected Output:</p>
                          <div class="font-mono text-xs bg-muted/50 p-2 rounded overflow-x-auto">
                            {testCase.expectedOutput}
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                {/each}
              </div>
            </div>
            
            <!-- Hints (collapsible) -->
            {#if currentChallenge.hints && currentChallenge.hints.length > 0}
              <Separator />
              
              <Collapsible open={showHints} onOpenChange={(open) => showHints = open}>
                <div class="flex items-center justify-between">
                  <h3 class="text-lg font-semibold flex items-center">
                    <Info class="h-4 w-4 mr-2" />
                    Hints
                  </h3>
                  
                  <CollapsibleTrigger asChild>
                    <Button variant="ghost" size="sm">
                      {#if showHints}
                        <ChevronRight class="h-4 w-4 rotate-90 transition-transform" />
                      {:else}
                        <ChevronRight class="h-4 w-4 transition-transform" />
                      {/if}
                    </Button>
                  </CollapsibleTrigger>
                </div>
                
                <CollapsibleContent>
                  <div class="space-y-2 mt-3">
                    <Alert variant="warning" class="bg-amber-50 dark:bg-amber-950/20">
                      <AlertDescription>
                        Hints can help you solve this challenge, but using them might reduce points earned.
                      </AlertDescription>
                    </Alert>
                    
                    {#each currentChallenge.hints as hint, i}
                      <Accordion type="single" collapsible class="w-full">
                        <AccordionItem value={`hint-${i}`}>
                          <AccordionTrigger class="text-sm">Hint {i + 1}</AccordionTrigger>
                          <AccordionContent>
                            <p class="text-sm">{hint}</p>
                          </AccordionContent>
                        </AccordionItem>
                      </Accordion>
                    {/each}
                  </div>
                </CollapsibleContent>
              </Collapsible>
            {/if}
          </div>
        </div>
      </div>
      
      <!-- Code editor and results (right) -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <!-- Editor toolbar -->
        <div class="border-b bg-muted/20 px-4 py-2 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <span class="font-semibold">Solution.c</span>
            <Badge>C</Badge>
          </div>
          
          <div class="flex items-center gap-2">
            <Tooltip.Root>
              <Tooltip.Trigger asChild>
               
                <button class="flex align-middle items-center border px-2 py-1 border-gray-200 rounded-md" on:click={()=>resetCode()}>
                  <RotateCcw class="h-4 w-4 mr-1" />
                  Reset
                </button>
              </Tooltip.Trigger>
              <Tooltip.Content>Reset to starter code</Tooltip.Content>
            </Tooltip.Root>
            
            <Tooltip.Root>
              <Tooltip.Trigger asChild>
                <button class="flex align-middle items-center border px-2 py-1 border-gray-200 rounded-md" on:click={()=>toggleFullScreen()}>
                  {#if isFullScreen}
                    <Minimize2 class="h-4 w-4 mr-1" />
                    Exit Fullscreen
                  {:else}
                    <Maximize2 class="h-4 w-4 mr-1" />
                    Fullscreen
                  {/if}
                </button>
              </Tooltip.Trigger>
              <Tooltip.Content>Toggle fullscreen (F11)</Tooltip.Content>
            </Tooltip.Root>
          </div>
        </div>
        

        <div class="flex-1 overflow-hidden relative">
          <div id="code-editor" class="absolute inset-0"></div>
        </div>
        
 
        <div class="border-t bg-muted/20 p-2 flex items-center justify-between">
          <div class="text-xs text-muted-foreground flex items-center">
            <Keyboard class="h-3 w-3 mr-1" />
            <span>Ctrl+Enter to run | Ctrl+S to save</span>
          </div>
          
   
   <button 
   on:click={() => submitCode()}
   disabled={submitting} 
   class="inline-flex items-center justify-center rounded-md text-sm font-medium bg-green-600 hover:bg-green-700 text-white px-4 py-2"
 >
   {#if submitting}
     <Loader class="h-4 w-4 mr-2 animate-spin" />
     Compiling...
   {:else}
     <Zap class="h-4 w-4 mr-2" />
     Run Code
   {/if}
 </button>
        </div>
        
      
        <div id="results-section" class="border-t overflow-y-auto" style="max-height: 40%;">
          {#if showResults && testResults.length > 0}
            <div class="p-4">
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-lg font-semibold flex items-center">
                  <Terminal class="h-5 w-5 mr-2" />
                  Test Results
                </h3>
                
                <div class="flex items-center">
                  <span class="text-sm font-medium mr-2">
                    {testResults.filter(r => r.passed).length}/{testResults.length} tests passed
                  </span>
                  <Progress 
                    value={(testResults.filter(r => r.passed).length / testResults.length) * 100} 
                    class="w-20 h-2" 
                    indicatorClass={testResults.every(t => t.passed) ? "bg-green-600" : "bg-amber-500"}
                  />
                </div>
              </div>
              
              <div class="space-y-3">
                {#each testResults as result, i}
                  <Card class={result.passed ? 
                    "border-green-200 dark:border-green-800" : 
                    "border-red-200 dark:border-red-800"}>
                    <CardHeader class="py-2 px-4 flex flex-row items-center justify-between">
                      <CardTitle class="text-sm flex items-center">
                        {#if result.passed}
                          <CheckCircle class="h-4 w-4 mr-2 text-green-600 dark:text-green-500" />
                          <span class="text-green-700 dark:text-green-400">Test {i + 1} Passed</span>
                        {:else}
                          <XCircle class="h-4 w-4 mr-2 text-red-600 dark:text-red-500" />
                          <span class="text-red-700 dark:text-red-400">Test {i + 1} Failed</span>
                        {/if}
                      </CardTitle>
                      
                      {#if result.executionTime !== undefined}
                        <Badge variant="outline" class="ml-auto">
                          <Clock class="h-3 w-3 mr-1" />
                          {result.executionTime}ms
                        </Badge>
                      {/if}
                    </CardHeader>
                    
                    <Collapsible>
                      <CollapsibleTrigger class="flex items-center justify-center w-full py-1 hover:bg-muted/50 text-xs text-muted-foreground">
                        <ChevronRight class="h-4 w-4 rotate-90 transition-transform" />
                        Details
                      </CollapsibleTrigger>
                      
                      <CollapsibleContent>
                        <CardContent class="py-3 px-4 border-t">
                          <div class="space-y-3">
                            {#if result.error}
                              <div>
                                <p class="text-sm font-medium mb-1">Error:</p>
                                <pre class="text-xs p-3 bg-red-50 dark:bg-red-950/20 text-red-800 dark:text-red-300 rounded-md overflow-x-auto">{result.error}</pre>
                              </div>
                            {/if}
                            
                            {#if result.output !== undefined}
                              <div>
                                <p class="text-sm font-medium mb-1">Your Output:</p>
                                <pre class="text-xs p-3 bg-muted rounded-md overflow-x-auto">{result.output}</pre>
                              </div>
                            {/if}
                          </div>
                        </CardContent>
                      </CollapsibleContent>
                    </Collapsible>
                  </Card>
                {/each}
              </div>
              
              {#if testResults.every(test => test.passed)}
                <div class="mt-4">
                  <Card class="border-green-200 dark:border-green-800 bg-green-50/50 dark:bg-green-950/20">
                    <CardContent class="py-4 flex items-center justify-between">
                      <div>
                        <h3 class="font-semibold text-green-700 dark:text-green-400">
                          Challenge Completed Successfully!
                        </h3>
                        <p class="text-sm text-green-600 dark:text-green-500">
                          All tests passed. Great job!
                        </p>
                      </div>
                      
                      <Button 
                        on:click={() => navigateToChallenge('next')} 
                        disabled={currentChallengeIndex === challengeIds.length - 1}
                        class="bg-green-600 hover:bg-green-700"
                      >
                        Next Challenge
                        <ArrowRight class="h-4 w-4 ml-2" />
                      </Button>
                    </CardContent>
                  </Card>
                </div>
              {/if}
            </div>
          {:else}
          <div class="flex flex-col items-center justify-center py-8 text-muted-foreground">
            <Terminal class="h-8 w-8 mb-3" />
            <h3 class="text-base font-medium">Run your code to see test results</h3>
            <p class="text-sm">Press Ctrl+Enter or click the Run Code button above</p>
          </div>
        {/if}
      </div>
    </div>
  {:else}
  
    <div class="flex-1 flex flex-col items-center justify-center">
      <Loader class="h-10 w-10 animate-spin text-muted-foreground mb-4" />
      <p class="text-lg font-medium">Loading challenges...</p>
      <p class="text-sm text-muted-foreground mt-2">Please wait while we load your programming challenges</p>
    </div>
  {/if}
</div>
</div>

<style>
:global(.monaco-editor .minimap) {
  display: none !important;
}

:global(.monaco-editor .margin) {
  background-color: transparent !important;
}

:global(.monaco-editor .line-numbers) {
  color: var(--text-muted) !important;
}

:global(body) {
  font-family: 'Inter', sans-serif;
}

.prose pre {
  padding: 0.75rem;
  border-radius: 0.375rem;
  background-color: rgb(var(--muted) / 0.5);
  overflow-x: auto;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.875rem;
}

.prose code {
  font-family: 'JetBrains Mono', monospace;
  background-color: rgb(var(--muted) / 0.5);
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
}
</style>