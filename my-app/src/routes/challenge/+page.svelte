<script lang="ts">
  import { onMount } from 'svelte';
  import { toast } from 'svelte-sonner';
  import { 
    Code, ChevronRight, ChevronLeft, Terminal, 
    RotateCcw, Zap, Info, ArrowRight,  
    Maximize2, Minimize2, Loader, Keyboard,
    CheckCircle, XCircle, Clock, Award,
    Server, Database, Download, Share2
  } from 'lucide-svelte';
  
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
  import { Tabs, TabsContent, TabsList, TabsTrigger } from "$lib/components/ui/tabs";
  import ModeToggle from "$lib/components/mode-toggle.svelte";
  import * as Tooltip from "$lib/components/ui/tooltip";

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
    constraints?: string;
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
    memory?: number;
  }

  interface ExecutionStats {
    totalTime: number;
    maxMemory: number;
    submissions: number;
    passRate: number;
  }

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
  let activeTab: string = 'description';
  let runStats: ExecutionStats = { totalTime: 0, maxMemory: 0, submissions: 0, passRate: 0 };
  let statusMessage: string = '';
  let lastSaveTime: Date | null = null;
  
  const editorThemes = [
    { value: 'vs-dark', label: 'VS Dark' },
    { value: 'github-dark', label: 'GitHub Dark' },
    { value: 'dracula', label: 'Dracula' },
    { value: 'nord', label: 'Nord' },
    { value: 'ayu-dark', label: 'Ayu Dark' },
    { value: 'tomorrow-night', label: 'Tomorrow Night' }
  ];

  const languageOptions = [
    { value: 50, label: 'C (GCC 9.2.0)' },
    { value: 54, label: 'C++ (GCC 9.2.0)' }
  ];
  
  let selectedLanguage = 50;

  onMount(async () => {
    await fetchChallenges();
    
    const textarea = document.getElementById('code-editor');
    
    if (!textarea) {
      console.error("Editor element not found");
      toast.error("Failed to initialize editor");
      return;
    }
    
    if (typeof monaco === 'undefined') {
      window.require(['vs/editor/editor.main'], initEditor);
    } else {
      initEditor();
    }

    window.addEventListener('keydown', handleKeyboardShortcuts);
    
    updateProgress();
    
    return () => {
      window.removeEventListener('keydown', handleKeyboardShortcuts);
      if (editor) editor.dispose();
    };
  });
  
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
        formatOnType: true,
        suggestOnTriggerCharacters: true,
        snippetSuggestions: "inline",
        quickSuggestions: true,
        folding: true
      });
      
      editor.onDidChangeModelContent(() => {
        code = editor.getValue();
        markEditorAsDirty();
      });
      
      const resizeObserver = new ResizeObserver(() => {
        if (editor) editor.layout();
      });
      
      resizeObserver.observe(document.getElementById('code-editor'));
      
    } catch (error) {
      console.error("Error initializing editor:", error);
      toast.error("Failed to initialize code editor. Please refresh the page.");
    }
  }

  function markEditorAsDirty() {
    lastSaveTime = null;
    const editorTitleElement = document.getElementById('editor-title');
    if (editorTitleElement) {
      editorTitleElement.classList.add('text-amber-500');
    }
  }
  
  function markEditorAsSaved() {
    lastSaveTime = new Date();
    const editorTitleElement = document.getElementById('editor-title');
    if (editorTitleElement) {
      editorTitleElement.classList.remove('text-amber-500');
    }
  }

  function handleKeyboardShortcuts(e: KeyboardEvent) {
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      e.preventDefault();
      submitCode();
    }
    
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
      e.preventDefault();
      saveCodeToLocalStorage();
      markEditorAsSaved();
      toast.success('Code saved');
    }
    
    if (e.key === 'F11') {
      e.preventDefault();
      toggleFullScreen();
    }
    
    if ((e.ctrlKey || e.metaKey) && e.key === 'b') {
      e.preventDefault();
      formatCode();
    }
  }
  
  function formatCode() {
    if (editor) {
      editor.getAction('editor.action.formatDocument').run();
      toast.success('Code formatted');
    }
  }
  
  function updateProgress() {
    if (!challengeIds.length) return;
    progress = ((currentChallengeIndex + 1) / challengeIds.length) * 100;
  }
  
  async function fetchChallenges() {
    try {
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
      
      const formattedCode = currentChallenge?.initialCode.replace(/\\n/g, '\n');
      code = formattedCode!;
      
      if (editor) {
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
      activeTab = 'description';
      resetRunStats();
      
      updateProgress();
      
      toast.dismiss(toastId);
      markEditorAsSaved();
    } catch (error) {
      console.error('Error loading challenge:', error);
      toast.error('Error loading challenge');
      currentChallenge = null;
    }
  }

  function resetRunStats() {
    runStats = {
      totalTime: 0,
      maxMemory: 0,
      submissions: 0,
      passRate: 0
    };
  }

  async function submitCode() {
    if (!currentChallenge) {
      toast.error('No challenge loaded');
      return;
    }

    submitting = true;
    statusMessage = 'Compiling code...';

    try {
      if (editor) {
        code = editor.getValue();
      }

      const tokenResponse = await fetch('https://codex.singularity.co.ke/api/auth/token');
      if (!tokenResponse.ok) {
        throw new Error('Failed to retrieve authentication token');
      }
      
      const { token } = await tokenResponse.json();

      statusMessage = 'Running tests...';
      
      const response = await fetch('https://api.singularity.co.ke/api/submit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          challengeId: currentChallenge.id,
          code: code,
          languageId: selectedLanguage
        })
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      showResults = true;
      testResults = result.testResults || [];
      
      // Update execution statistics
      updateExecutionStats(testResults);

      // Save code after successful run
      saveCodeToLocalStorage();
      markEditorAsSaved();

      setTimeout(() => {
        const resultsSection = document.getElementById('results-section');
        if (resultsSection) {
          resultsSection.scrollIntoView({ behavior: 'smooth' });
        }
      }, 100);
      
      if (testResults.every(test => test.passed)) {
        toast.success('All tests passed!');
      }

    } catch (error) {
      console.error('Error submitting code:', error);
      toast.error('Error submitting code. Please try again.');
    } finally {
      submitting = false;
      statusMessage = '';
    }
  }
  
  function updateExecutionStats(results: TestResult[]) {
    runStats.submissions++;
    
    let totalExecTime = 0;
    let maxMem = 0;
    let passedTests = 0;
    
    results.forEach(result => {
      if (result.executionTime) {
        totalExecTime += parseFloat(result.executionTime.toString());
      }
      
      if (result.memory && result.memory > maxMem) {
        maxMem = result.memory;
      }
      
      if (result.passed) {
        passedTests++;
      }
    });
    
    runStats.totalTime = totalExecTime;
    runStats.maxMemory = maxMem;
    runStats.passRate = results.length > 0 ? (passedTests / results.length) * 100 : 0;
  }
  
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
  
  function saveCodeToLocalStorage() {
    if (currentChallenge) {
      localStorage.setItem(`c_master_challenge_${currentChallenge.id}`, code);
    }
  }
  
  function downloadCode() {
    if (!currentChallenge) return;
    
    const blob = new Blob([code], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${currentChallenge.title.replace(/\s+/g, '_').toLowerCase()}.c`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    toast.success('Code downloaded');
  }
  
  function shareChallenge() {
    if (!currentChallenge) return;
    
    const url = `${window.location.origin}/challenge/${currentChallenge.id}`;
    
    if (navigator.clipboard) {
      navigator.clipboard.writeText(url)
        .then(() => toast.success('Challenge link copied to clipboard'))
        .catch(err => toast.error('Failed to copy link'));
    } else {
      const input = document.createElement('input');
      input.value = url;
      document.body.appendChild(input);
      input.select();
      document.execCommand('copy');
      document.body.removeChild(input);
      toast.success('Challenge link copied to clipboard');
    }
  }
  
  function changeEditorTheme(newTheme: string) {
    editorTheme = newTheme;
    if (editor) {
      monaco.editor.setTheme(newTheme);
      toast.success(`Theme changed to ${newTheme}`);
    }
  }
  
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
  
  function toggleFullScreen() {
    isFullScreen = !isFullScreen;
    
    const appContainer = document.getElementById('c-master-app');
    
    if (isFullScreen) {
      if (appContainer?.requestFullscreen) {
        appContainer.requestFullscreen();
      }
    } else {
      if (document.exitFullscreen) {
        document.exitFullscreen();
      }
    }
    
    setTimeout(() => {
      if (editor) editor.layout();
    }, 100);
  }
  
  function getDifficultyColor(difficulty: string) {
    switch(difficulty.toLowerCase()) {
      case 'easy': return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900 dark:text-emerald-300';
      case 'medium': return 'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-300';
      case 'hard': return 'bg-rose-100 text-rose-800 dark:bg-rose-900 dark:text-rose-300';
      default: return 'bg-slate-100 text-slate-800 dark:bg-slate-800 dark:text-slate-300';
    }
  }
  
  function formatMemory(bytes: number) {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  }
</script>

<svelte:head>
  <title>Codex - {currentChallenge ? currentChallenge.title : 'Loading...'}</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
  
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
  <div class="flex-1 flex overflow-hidden">
    {#if currentChallenge}
      <div class="w-[450px] flex flex-col border-r overflow-hidden">
        <div class="border-b px-4 py-3 flex items-center justify-between bg-muted/30">
          <div class="flex items-center gap-2">
            <Button variant="outline" size="sm" disabled={currentChallengeIndex === 0} on:click={() => navigateToChallenge('prev')}>
              <ChevronLeft class="h-4 w-4 mr-1" />
              Previous
            </Button>
            
            <Button variant="outline" size="sm" disabled={currentChallengeIndex === challengeIds.length - 1} on:click={() => navigateToChallenge('next')}>
              Next
              <ChevronRight class="h-4 w-4 ml-1" />
            </Button>
          </div>
          
          <div class="flex items-center gap-3">
            <span class="text-sm font-medium">
              {currentChallengeIndex + 1}/{challengeIds.length}
            </span>
            <Progress value={progress} class="w-20 h-1.5" />
            <div class="flex items-center gap-1">
              <Tooltip.Root>
                <Tooltip.Trigger asChild>
                  <button on:click={shareChallenge} class="p-1 hover:bg-muted rounded-md">
                    <Share2 class="h-4 w-4" />
                  </button>
                </Tooltip.Trigger>
                <Tooltip.Content>Share challenge</Tooltip.Content>
              </Tooltip.Root>
            </div>
          </div>
        </div>
        
        <Tabs value={activeTab} onValueChange={(val) => activeTab = val} class="flex-1 flex flex-col overflow-hidden">
          <TabsList class="px-4 pt-2 pb-0 gap-2 border-b h-auto bg-transparent justify-start">
            <TabsTrigger value="description" class="px-3 py-1.5 data-[state=active]:bg-muted rounded-t-md data-[state=active]:border-b-0 data-[state=active]:border data-[state=active]:border-muted">
              Description
            </TabsTrigger>
            <TabsTrigger value="hints" class="px-3 py-1.5 data-[state=active]:bg-muted rounded-t-md data-[state=active]:border-b-0 data-[state=active]:border data-[state=active]:border-muted">
              Hints
            </TabsTrigger>
            <TabsTrigger value="stats" class="px-3 py-1.5 data-[state=active]:bg-muted rounded-t-md data-[state=active]:border-b-0 data-[state=active]:border data-[state=active]:border-muted">
              Stats
            </TabsTrigger>
          </TabsList>
          
          <TabsContent value="description" class="flex-1 overflow-y-auto p-4 mt-0">
            <div class="space-y-6">
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

                  {#if currentChallenge.memoryLimit}
                    <Badge variant="outline">
                      <Database class="h-3.5 w-3.5 mr-1" />
                      {currentChallenge.memoryLimit}KB
                    </Badge>
                  {/if}
                </div>
              </div>
              
              <div class="prose dark:prose-invert max-w-none">
                {@html currentChallenge.description.replace(/\n/g, '<br>')}
              </div>
              
              {#if currentChallenge.constraints}
                <div class="mt-4">
                  <h3 class="text-lg font-semibold mb-2">Constraints:</h3>
                  <div class="prose dark:prose-invert max-w-none">
                    {@html currentChallenge.constraints.replace(/\n/g, '<br>')}
                  </div>
                </div>
              {/if}
              
              <Separator />
              
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
            </div>
          </TabsContent>
          
          <TabsContent value="hints" class="flex-1 overflow-y-auto p-4 mt-0">
            {#if currentChallenge.hints && currentChallenge.hints.length > 0}
              <div class="space-y-4">
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
            {:else}
              <div class="flex flex-col items-center justify-center h-full text-muted-foreground">
                <Info class="h-12 w-12 mb-4 opacity-20" />
                <p class="text-lg font-medium">No hints available</p>
                <p class="text-sm">Try solving this challenge on your own!</p>
              </div>
            {/if}
          </TabsContent>
          
          <TabsContent value="stats" class="flex-1 overflow-y-auto p-4 mt-0">
            <div class="space-y-6">
              <div>
                <h3 class="text-lg font-semibold mb-3 flex items-center">
                  <Server class="h-4 w-4 mr-2" />
                  Judge0 Execution Statistics
                </h3>
                
                <div class="grid grid-cols-2 gap-4">
                  <Card>
                    <CardContent class="pt-6">
                      <div class="flex flex-col items-center">
                        <p class="text-sm text-muted-foreground mb-1">Execution Time</p>
                        <p class="text-2xl font-bold">{runStats.totalTime.toFixed(2)}ms</p>
                        {#if currentChallenge.timeLimit}
                          <p class="text-xs text-muted-foreground mt-1">
                            {(runStats.totalTime / currentChallenge.timeLimit * 100).toFixed(1)}% of time limit
                          </p>
                          <Progress
                            value={(runStats.totalTime / currentChallenge.timeLimit * 100)}
                            class="w-full h-1 mt-2"
                            indicatorClass={runStats.totalTime / currentChallenge.timeLimit > 0.8 ? "bg-amber-500" : "bg-green-600"}
                          />
                        {/if}
                        </div>
                    </CardContent>
                  </Card>
                  
                  <Card>
                    <CardContent class="pt-6">
                      <div class="flex flex-col items-center">
                        <p class="text-sm text-muted-foreground mb-1">Memory Usage</p>
                        <p class="text-2xl font-bold">{formatMemory(runStats.maxMemory)}</p>
                        {#if currentChallenge.memoryLimit}
                          <p class="text-xs text-muted-foreground mt-1">
                            {(runStats.maxMemory / (currentChallenge.memoryLimit * 1024) * 100).toFixed(1)}% of memory limit
                          </p>
                          <Progress
                            value={(runStats.maxMemory / (currentChallenge.memoryLimit * 1024) * 100)}
                            class="w-full h-1 mt-2"
                            indicatorClass={runStats.maxMemory / (currentChallenge.memoryLimit * 1024) > 0.8 ? "bg-amber-500" : "bg-green-600"}
                          />
                        {/if}
                      </div>
                    </CardContent>
                  </Card>
                  
                  <Card>
                    <CardContent class="pt-6">
                      <div class="flex flex-col items-center">
                        <p class="text-sm text-muted-foreground mb-1">Success Rate</p>
                        <p class="text-2xl font-bold">{runStats.passRate.toFixed(1)}%</p>
                        <p class="text-xs text-muted-foreground mt-1">
                          From {runStats.submissions} submissions
                        </p>
                        <Progress
                          value={runStats.passRate}
                          class="w-full h-1 mt-2"
                          indicatorClass={runStats.passRate < 50 ? "bg-red-500" : runStats.passRate < 80 ? "bg-amber-500" : "bg-green-600"}
                        />
                      </div>
                    </CardContent>
                  </Card>
                  
                  <Card>
                    <CardContent class="pt-6">
                      <div class="flex flex-col items-center">
                        <p class="text-sm text-muted-foreground mb-1">Judge0 Processing</p>
                        <p class="text-2xl font-bold">v1.13.0</p>
                        <p class="text-xs text-muted-foreground mt-1">
                          GCC 9.2.0 compliant
                        </p>
                        <div class="flex items-center gap-2 mt-2">
                          <Badge variant="outline" class="bg-green-50 dark:bg-green-950/20 text-green-700 dark:text-green-400 border-green-200 dark:border-green-800">
                            Standard C11
                          </Badge>
                          <Badge variant="outline" class="bg-blue-50 dark:bg-blue-950/20 text-blue-700 dark:text-blue-400 border-blue-200 dark:border-blue-800">
                            Judge0 API
                          </Badge>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </div>
              </div>
              
              {#if runStats.submissions > 0}
                <div>
                  <h3 class="text-lg font-semibold mb-3">Execution History</h3>
                  <Card>
                    <div class="overflow-x-auto">
                      <table class="w-full">
                        <thead>
                          <tr class="border-b">
                            <th class="text-left p-3 text-sm font-medium text-muted-foreground">Run</th>
                            <th class="text-left p-3 text-sm font-medium text-muted-foreground">Status</th>
                            <th class="text-left p-3 text-sm font-medium text-muted-foreground">Time</th>
                            <th class="text-left p-3 text-sm font-medium text-muted-foreground">Memory</th>
                          </tr>
                        </thead>
                        <tbody>
                          {#if testResults.length > 0}
                            {#each testResults as result, i}
                              <tr class="border-b last:border-0 hover:bg-muted/30">
                                <td class="p-3 text-sm">Test Case {i + 1}</td>
                                <td class="p-3 text-sm">
                                  {#if result.passed}
                                    <span class="inline-flex items-center text-green-600 dark:text-green-400">
                                      <CheckCircle class="h-3.5 w-3.5 mr-1" />
                                      Passed
                                    </span>
                                  {:else}
                                    <span class="inline-flex items-center text-red-600 dark:text-red-400">
                                      <XCircle class="h-3.5 w-3.5 mr-1" />
                                      Failed
                                    </span>
                                  {/if}
                                </td>
                                <td class="p-3 text-sm">
                                  {result.executionTime ? `${result.executionTime}ms` : 'N/A'}
                                </td>
                                <td class="p-3 text-sm">
                                  {result.memory ? formatMemory(result.memory) : 'N/A'}
                                </td>
                              </tr>
                            {/each}
                          {:else}
                            <tr>
                              <td colspan="4" class="p-3 text-sm text-center text-muted-foreground">
                                No execution data available
                              </td>
                            </tr>
                          {/if}
                        </tbody>
                      </table>
                    </div>
                  </Card>
                </div>
              {/if}
            </div>
          </TabsContent>
        </Tabs>
      </div>
      
      <div class="flex-1 flex flex-col overflow-hidden">
        <div class="border-b bg-muted/20 px-4 py-2 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <span id="editor-title" class="font-semibold">Solution.c</span>
            <div class="flex gap-2">
              <Badge variant="secondary">GCC 9.2.0</Badge>
              <Select value={selectedLanguage.toString()} onValueChange={(val) => selectedLanguage = parseInt(val)}>
                <SelectTrigger class="h-7 w-[120px] text-xs">
                  <SelectValue placeholder="Language" />
                </SelectTrigger>
                <SelectContent>
                  {#each languageOptions as option}
                    <SelectItem value={option.value.toString()} class="text-xs">{option.label}</SelectItem>
                  {/each}
                </SelectContent>
              </Select>
            </div>
          </div>
          
          <div class="flex items-center gap-2">
            <Select value={editorTheme} onValueChange={changeEditorTheme}>
              <SelectTrigger class="h-7 w-[140px] text-xs">
                <SelectValue placeholder="Editor Theme" />
              </SelectTrigger>
              <SelectContent>
                {#each editorThemes as theme}
                  <SelectItem value={theme.value} class="text-xs">{theme.label}</SelectItem>
                {/each}
              </SelectContent>
            </Select>
            
            <div class="flex gap-1">
              <Tooltip.Root>
                <Tooltip.Trigger asChild>
                  <Button variant="outline" size="icon" class="h-7 w-7" on:click={resetCode}>
                    <RotateCcw class="h-3.5 w-3.5" />
                  </Button>
                </Tooltip.Trigger>
                <Tooltip.Content>Reset code</Tooltip.Content>
              </Tooltip.Root>
              
              <Tooltip.Root>
                <Tooltip.Trigger asChild>
                  <Button variant="outline" size="icon" class="h-7 w-7" on:click={formatCode}>
                    <Code class="h-3.5 w-3.5" />
                  </Button>
                </Tooltip.Trigger>
                <Tooltip.Content>Format code (Ctrl+B)</Tooltip.Content>
              </Tooltip.Root>
              
              <Tooltip.Root>
                <Tooltip.Trigger asChild>
                  <Button variant="outline" size="icon" class="h-7 w-7" on:click={downloadCode}>
                    <Download class="h-3.5 w-3.5" />
                  </Button>
                </Tooltip.Trigger>
                <Tooltip.Content>Download code</Tooltip.Content>
              </Tooltip.Root>
              
              <Tooltip.Root>
                <Tooltip.Trigger asChild>
                  <Button variant="outline" size="icon" class="h-7 w-7" on:click={toggleFullScreen}>
                    {#if isFullScreen}
                      <Minimize2 class="h-3.5 w-3.5" />
                    {:else}
                      <Maximize2 class="h-3.5 w-3.5" />
                    {/if}
                  </Button>
                </Tooltip.Trigger>
                <Tooltip.Content>Toggle fullscreen (F11)</Tooltip.Content>
              </Tooltip.Root>
            </div>
          </div>
        </div>
        
        <div class="flex-1 overflow-hidden relative">
          <div id="code-editor" class="absolute inset-0"></div>
        </div>
        
        <div class="border-t bg-muted/20 px-3 py-2 flex items-center justify-between">
          <div class="flex items-center text-xs text-muted-foreground">
            <Keyboard class="h-3 w-3 mr-1" />
            <span>Ctrl+Enter to run | Ctrl+S to save | Ctrl+B to format</span>
            {#if lastSaveTime}
              <span class="ml-3 text-green-600 dark:text-green-400">
                Saved {lastSaveTime.toLocaleTimeString()}
              </span>
            {/if}
          </div>
          
          <button 
            on:click={() => submitCode()}
            disabled={submitting} 
            class="inline-flex items-center justify-center rounded-md text-sm font-medium bg-green-600 hover:bg-green-700 text-white px-4 py-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {#if submitting}
              <Loader class="h-4 w-4 mr-2 animate-spin" />
              {statusMessage || 'Processing...'}
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
                      
                      <div class="flex items-center gap-2">
                        {#if result.executionTime !== undefined}
                          <Badge variant="outline" class="ml-auto">
                            <Clock class="h-3 w-3 mr-1" />
                            {result.executionTime}ms
                          </Badge>
                        {/if}
                        
                        {#if result.memory !== undefined}
                          <Badge variant="outline" class="ml-auto">
                            <Database class="h-3 w-3 mr-1" />
                            {formatMemory(result.memory)}
                          </Badge>
                        {/if}
                      </div>
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
                                <pre class="text-xs p-3 bg-red-50 dark:bg-red-950/20 text-red-800 dark:text-red-300 rounded-md overflow-x-auto whitespace-pre-wrap">{result.error}</pre>
                              </div>
                            {/if}
                            
                            {#if result.output !== undefined}
                              <div>
                                <p class="text-sm font-medium mb-1">Your Output:</p>
                                <pre class="text-xs p-3 bg-muted rounded-md overflow-x-auto whitespace-pre-wrap">{result.output}</pre>
                              </div>
                            {/if}
                            
                            {#if currentChallenge.testCases[i] && currentChallenge.testCases[i].expectedOutput}
                              <div>
                                <p class="text-sm font-medium mb-1">Expected Output:</p>
                                <pre class="text-xs p-3 bg-muted/50 rounded-md overflow-x-auto whitespace-pre-wrap">{currentChallenge.testCases[i].expectedOutput}</pre>
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