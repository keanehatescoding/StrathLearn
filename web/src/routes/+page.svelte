<script lang="ts">
    import { onMount } from 'svelte';
    import { 
      Code, 
      ChevronRight, 
      Lightbulb, 
      ClipboardCheck, 
      RotateCcw, 
      Zap, 
      FileText, 
      CheckCircle, 
      XCircle,
      Loader
    } from 'lucide-svelte';
    

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
    }
  

    let challenges: Record<string, Challenge> = {};
    let currentChallenge: Challenge | null = null;
    let selectedChallengeId: string = '';
    let code: string = '';
    let editor: any;
    let testResults: TestResult[] = [];
    let showResults: boolean = false;
    let showSuccessMessage: boolean = false;
    let submitting: boolean = false;
    let theme: string = 'material-palenight';
    

    let hintsOpen = false;
    let testCasesOpen = false;
    
    onMount(async () => {
      await fetchChallenges();
      

      const textarea = document.getElementById('code-editor') as HTMLTextAreaElement;
      
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
          "Ctrl-Space": "autocomplete"
        },
        hintOptions: {
          completeSingle: false
        }
      });
      
      editor.on('change', () => {
        code = editor.getValue();
      });
      
      //Autocomplete
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
    });
    
    async function fetchChallenges() {
      try {
        const response = await fetch('https://strathlearn-3ba26ebd38c9.herokuapp.com/api/challenges');
        if (!response.ok) {
          throw new Error('Failed to load challenges');
        }
        
        challenges = await response.json();
        
  
       
          selectedChallengeId = Object.keys(challenges)[0];
        
        
        if (selectedChallengeId) {
          await loadChallenge(selectedChallengeId);
        }
      } catch (error) {
        console.error('Error fetching challenges:', error);
      }
    }
    
    async function loadChallenge(id: string) {
      try {
        const response = await fetch(`https://strathlearn-3ba26ebd38c9.herokuapp.com/api/challenge/${id}`);
        if (!response.ok) {
          throw new Error('Challenge not found');
        }
        
        currentChallenge = await response.json();
        const formattedCode = currentChallenge?.initialCode.replace(/\\n/g, '\n');
        code = formattedCode!;
        
        if (editor) {
          editor.setValue(formattedCode);
        }
        
   
        showResults = false;
        showSuccessMessage = false;
        testResults = [];
        
      } catch (error) {
        console.error('Error loading challenge:', error);
        currentChallenge = null;
      }
    }
    
    async function submitCode() {
      if (!currentChallenge) {
        alert('No challenge loaded');
        return;
      }
      
      submitting = true;
      
      try {
  
        if (editor) {
          code = editor.getValue();
        }
        
        const response = await fetch('https://strathlearn-3ba26ebd38c9.herokuapp.com/api/submit', {
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
        
        showResults = true;
        testResults = result.testResults || [];
        
   
        showSuccessMessage = testResults.length > 0 && testResults.every(test => test.passed);
        
      } catch (error) {
        console.error('Error submitting code:', error);
        alert('Error submitting code. Please try again.');
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
      }
    }
    
    function toggleHints() {
      hintsOpen = !hintsOpen;
    }
    
    function toggleTestCases() {
      testCasesOpen = !testCasesOpen;
    }
    
    function handleChallengeSelect() {
      if (selectedChallengeId) {
        loadChallenge(selectedChallengeId);
      }
    }
    
    function changeTheme(newTheme: string) {
      theme = newTheme;
      if (editor) {
        editor.setOption('theme', newTheme);
      }
    }
  </script>
  
  <svelte:head>
    <title>StrathLearn - Code Challenge Platform</title>

    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/codemirror.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/theme/github-dark.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/theme/material-palenight.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/theme/dracula.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/fold/foldgutter.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/hint/show-hint.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/scroll/simplescrollbars.min.css">
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
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/search/match-highlighter.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.2/addon/scroll/simplescrollbars.min.js"></script>
  </svelte:head>
  
  <div class="min-h-screen bg-slate-50 flex flex-col font-sans">
    <header class="bg-gradient-to-r from-indigo-900 to-blue-800 text-white shadow-lg">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <div class="text-2xl font-bold tracking-tight flex items-center gap-2">
          <Code class="h-6 w-6 text-cyan-300" />
          StrathLearn
        </div>
        <nav>
          <ul class="flex space-x-8">
            <li><a href="/" class="hover:text-cyan-300 font-medium text-cyan-100 transition-colors duration-200 py-2 flex items-center border-b-2 border-cyan-400">Challenges</a></li>
            <li><a href="#" class="hover:text-cyan-300 text-white/80 transition-colors duration-200 py-2 flex items-center border-b-2 border-transparent">Leaderboard</a></li>
            <li><a href="#" class="hover:text-cyan-300 text-white/80 transition-colors duration-200 py-2 flex items-center border-b-2 border-transparent">About</a></li>
          </ul>
        </nav>
      </div>
    </header>
  
    <main class="flex-grow">
      <div class="container mx-auto px-4 py-8">
        <div class="mb-6 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label for="challenge-select" class="block text-sm font-medium text-slate-700 mb-2">Select Challenge:</label>
            <select 
              id="challenge-select"
              class="block w-full px-3 py-2 border border-slate-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white transition-all duration-200"
              bind:value={selectedChallengeId}
              on:change={handleChallengeSelect}
            >
              {#each Object.entries(challenges) as [id, challenge]}
                <option value={id}>{challenge.title}</option>
              {/each}
            </select>
          </div>
          
          <div>
            <label for="theme-select" class="block text-sm font-medium text-slate-700 mb-2">Editor Theme:</label>
            <select 
              id="theme-select"
              class="block w-full px-3 py-2 border border-slate-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white transition-all duration-200"
              bind:value={theme}
              on:change={() => changeTheme(theme)}
            >
              <option value="github-dark">GitHub Dark</option>
              <option value="material-palenight">Material Palenight</option>
              <option value="dracula">Dracula</option>
            </select>
          </div>
        </div>
        
        {#if currentChallenge}
          <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
            <div class="lg:col-span-2 bg-white rounded-xl shadow-md border border-slate-200 overflow-hidden">
              <div class="bg-gradient-to-r from-slate-50 to-white p-6 border-b border-slate-200">
                <h1 class="text-2xl font-bold mb-2 text-indigo-900">{currentChallenge.title}</h1>
                <div class="mb-4">
                  <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold 
                    {currentChallenge.difficulty === 'Easy' ? 'bg-emerald-100 text-emerald-800' : 
                    currentChallenge.difficulty === 'Medium' ? 'bg-amber-100 text-amber-800' : 
                    'bg-rose-100 text-rose-800'}">
                    {#if currentChallenge.difficulty === 'Easy'}
                      <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                    {:else if currentChallenge.difficulty === 'Medium'}
                      <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span>
                    {:else}
                      <span class="w-1.5 h-1.5 rounded-full bg-rose-500"></span>
                    {/if}
                    {currentChallenge.difficulty}
                  </span>
                </div>
                
                <div class="prose prose-slate mb-6 max-w-none">
                  {@html currentChallenge.description.replace(/\n/g, '<br>')}
                </div>
                
       
                <div class="mb-4 border rounded-lg overflow-hidden shadow-sm transition-all duration-200">
                  <button 
                    class="w-full flex justify-between items-center p-3 bg-indigo-50 hover:bg-indigo-100 text-left font-medium text-indigo-700 transition-colors"
                    on:click={toggleHints}
                  >
                    <div class="flex items-center">
                      <Lightbulb class="h-5 w-5 mr-2" />
                      Hints
                    </div>
                    <span class="text-xl">{hintsOpen ? '−' : '+'}</span>
                  </button>
                  
                  {#if hintsOpen}
                    <div class="p-4 space-y-2 bg-white">
                      {#each currentChallenge.hints as hint, i}
                        <div class="p-3 bg-amber-50 border-l-4 border-amber-400 text-amber-800 rounded">
                          <p class="font-medium mb-1">Hint {i + 1}:</p>
                          {hint}
                        </div>
                      {/each}
                    </div>
                  {/if}
                </div>
                

                <div class="border rounded-lg overflow-hidden shadow-sm transition-all duration-200">
                  <button 
                    class="w-full flex justify-between items-center p-3 bg-violet-50 hover:bg-violet-100 text-left font-medium text-violet-700 transition-colors"
                    on:click={toggleTestCases}
                  >
                    <div class="flex items-center">
                      <ClipboardCheck class="h-5 w-5 mr-2" />
                      Test Cases
                    </div>
                    <span class="text-xl">{testCasesOpen ? '−' : '+'}</span>
                  </button>
                  
                  {#if testCasesOpen}
                    <div class="p-4 space-y-4 bg-white">
                      {#each currentChallenge.testCases.filter(tc => !tc.hidden) as testCase, i}
                        <div class="border rounded-lg p-4 bg-gray-50 shadow-sm">
                          <p class="text-xs uppercase tracking-wider text-gray-500 mb-2 font-semibold">Test Case {i + 1}</p>
                          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                              <p class="font-medium text-sm text-gray-700 mb-1">Input:</p>
                              <div class="font-mono text-sm bg-white p-2 rounded border">
                                {testCase.input ? `"${testCase.input}"` : '<em>(empty)</em>'}
                              </div>
                            </div>
                            <div>
                              <p class="font-medium text-sm text-gray-700 mb-1">Expected Output:</p>
                              <div class="font-mono text-sm bg-white p-2 rounded border">
                                "{testCase.expectedOutput.replace(/\\n/g, '\\\\n')}"
                              </div>
                            </div>
                          </div>
                        </div>
                      {/each}
                    </div>
                  {/if}
                </div>
              </div>
            </div>
            
            <div class="lg:col-span-3 flex flex-col">
              <div class="bg-white rounded-xl shadow-md border border-slate-200 overflow-hidden flex-grow flex flex-col">
                <div class="flex justify-between items-center px-6 py-4 border-b border-slate-200 bg-slate-50">
                  <h2 class="text-lg font-semibold text-slate-800 flex items-center">
                    <Code class="h-5 w-5 mr-2 text-blue-500" />
                    Your Code
                  </h2>
                  <div class="flex space-x-2 items-center">
                    <span class="text-xs text-slate-500 bg-slate-100 px-2 py-1 rounded">Ctrl+Space for autocomplete</span>
                    <button 
                      on:click={resetCode}
                      class="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded text-sm font-medium text-slate-700 transition-colors flex items-center"
                    >
                      <RotateCcw class="h-4 w-4 mr-1" />
                      Reset
                    </button>
                  </div>
                </div>
                
                <div class="h-[500px] border-b border-slate-200 flex-grow">
                  <textarea id="code-editor" class="hidden">{code}</textarea>
                </div>
                
                <div class="p-4 bg-slate-50 border-t border-slate-200">
                  <button 
                    on:click={submitCode} 
                    disabled={submitting}
                    class="w-full py-3 px-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-medium rounded-lg shadow-sm disabled:opacity-50 transition-all duration-200 flex items-center justify-center gap-2"
                  >
                    {#if submitting}
                      <Loader class="h-5 w-5 animate-spin" />
                      Submitting...
                    {:else}
                      <Zap class="h-5 w-5" />
                      Run & Submit Solution
                    {/if}
                  </button>
                </div>
                
                {#if showResults}
                  <div class="p-6 border-t border-slate-200">
                    <h3 class="text-lg font-semibold mb-4 text-slate-800 flex items-center">
                      <FileText class="h-5 w-5 mr-2 text-slate-600" />
                      Test Results
                    </h3>
                    
                    <div class="space-y-4">
                      {#each testResults as result, i}
                        <div class={`border rounded-lg p-4 ${result.passed ? 'bg-emerald-50 border-emerald-200' : 'bg-rose-50 border-rose-200'} transition-all duration-200`}>
                          <div class="flex justify-between items-center mb-2">
                            <h4 class="font-medium flex items-center">
                              {#if result.passed}
                                <CheckCircle class="h-5 w-5 mr-2 text-emerald-500" />
                              {:else}
                                <XCircle class="h-5 w-5 mr-2 text-rose-500" />
                              {/if}
                              Test Case {i + 1}: {result.passed ? 'PASSED' : 'FAILED'}
                            </h4>
                            <span class={`text-xs font-semibold px-2 py-1 rounded-full ${result.passed ? 'bg-emerald-100 text-emerald-800' : 'bg-rose-100 text-rose-800'}`}>
                              {result.passed ? 'Success' : 'Error'}
                            </span>
                          </div>
                          
                          {#if result.error}
                            <div class="mt-3">
                              <p class="text-sm font-medium text-rose-700 mb-1">Error Message:</p>
                              <pre class="mt-2 p-3 bg-white rounded-lg border border-rose-200 text-sm overflow-x-auto whitespace-pre-wrap text-rose-600 font-mono">{result.error}</pre>
                            </div>
                          {/if}
                          
                          {#if result.output !== undefined}
                            <div class="mt-3">
                              <p class="text-sm font-medium text-slate-700 mb-1">Your Output:</p>
                              <pre class="p-3 bg-white rounded-lg border border-slate-200 text-sm overflow-x-auto whitespace-pre-wrap font-mono">{result.output}</pre>
                            </div>
                          {/if}
                        </div>
                      {/each}
                    </div>
                    
                    {#if showSuccessMessage}
                      <div class="mt-6 p-5 bg-emerald-50 border border-emerald-200 rounded-lg shadow-sm animate-fadeIn">
                        <div class="flex items-center">
                          <div class="flex-shrink-0">
                            <CheckCircle class="h-10 w-10 text-emerald-400" />
                          </div>
                          <div class="ml-4">
                            <h2 class="text-xl font-bold text-emerald-800">All tests passed! Congratulations!</h2>
                            <p class="text-emerald-700 mt-1">You've successfully completed this challenge.</p>
                          </div>
                        </div>
                      </div>
                    {/if}
                  </div>
                {/if}
              </div>
            </div>
          </div>
        {:else}
          <div class="text-center py-16 bg-white rounded-xl shadow-md border border-slate-200">
            <Loader class="h-16 w-16 mx-auto text-slate-400 mb-4 animate-spin" />
            <p class="text-xl text-slate-600">Loading challenge...</p>
          </div>
        {/if}
      </div>
    </main>
  
    <footer class="bg-slate-800 text-white/80 py-6 mt-8">
      <div class="container mx-auto px-4 flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center gap-2 mb-4 md:mb-0">
          <Code class="h-5 w-5 text-blue-400" />
          <span class="font-semibold text-white">StrathLearn</span>
        </div>
        <div class="flex items-center space-x-6 mt-4 md:mt-0">
          <a href="#" class="text-white/80 hover:text-white transition-colors duration-200">Terms</a>
          <a href="#" class="text-white/80 hover:text-white transition-colors duration-200">Privacy</a>
          <a href="#" class="text-white/80 hover:text-white transition-colors duration-200">Help</a>
        </div>
      </div>
    </footer>
  </div>
    
  <style>
    :global(body) {
      font-family: 'Inter', sans-serif;
    }

    :global(.CodeMirror) {
      height: 100% !important;
      font-family: 'JetBrains Mono', monospace;
      font-size: 14px;
      line-height: 1.6;
    }
    
    :global(.CodeMirror-gutters) {
      border-right: 1px solid rgba(0,0,0,0.05);
      background-color: rgba(0,0,0,0.02);
    }
    
    :global(.CodeMirror-linenumber) {
      color: rgba(0,0,0,0.4);
      padding: 0 8px;
    }
    
    :global(.CodeMirror-selected) {
      background: rgba(66, 153, 225, 0.15) !important;
    }
    
    :global(.CodeMirror-activeline-background) {
      background: rgba(0,0,0,0.04);
    }
    
    :global(.CodeMirror-matchingbracket) {
      color: #38a169 !important;
      border-bottom: 1px solid #38a169;
      font-weight: bold;
    }
    
    :global(.CodeMirror-scrollbar-filler) {
      background-color: transparent;
    }
    
    :global(.CodeMirror-hint) {
      font-family: 'JetBrains Mono', monospace;
      font-size: 13px;
      padding: 4px 8px;
      border-radius: 2px;
    }
    
    :global(.CodeMirror-hints) {
      border-radius: 4px;
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
      padding: 2px 0;
    }
    
    :global(.cm-s-github-dark.CodeMirror), :global(.cm-s-material-palenight.CodeMirror), :global(.cm-s-dracula.CodeMirror) {
      background-color: #0d1117;
      color: #e1e4e8;
    }
    
    :global(.CodeMirror-simplescroll-vertical) {
      width: 8px;
      border-radius: 4px;
      margin-right: 4px;
    }
    
    :global(.CodeMirror-simplescroll-horizontal) {
      height: 8px;
      border-radius: 4px;
      margin-bottom: 4px;
    }
    
    /* Smooth transitions */
    :global(.transition-height) {
      transition: max-height 0.3s ease-in-out;
      overflow: hidden;
    }
    
    :global(.prose) {
      line-height: 1.7;
    }
    
    :global(.prose p) {
      margin-bottom: 1rem;
    }
    
    :global(.prose ul) {
      list-style-type: disc;
      margin-left: 1.25rem;
      margin-bottom: 1rem;
    }
    
    :global(.prose ol) {
      list-style-type: decimal;
      margin-left: 1.25rem;
      margin-bottom: 1rem;
    }
    
    :global(.prose code) {
      font-family: 'JetBrains Mono', monospace;
      background-color: rgba(0,0,0,0.05);
      padding: 0.1em 0.3em;
      border-radius: 3px;
      font-size: 0.9em;
    }
    
    :global(.prose pre) {
      background-color: #f6f8fa;
      border-radius: 6px;
      padding: 1rem;
      overflow-x: auto;
      margin: 1rem 0;
    }
    
    :global(.prose pre code) {
      background-color: transparent;
      padding: 0;
      font-size: 0.9em;
      color: #24292e;
    }
    
    /* Animation utilities */
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(10px); }
      to { opacity: 1; transform: translateY(0); }
    }
    
    .animate-fadeIn {
      animation: fadeIn 0.5s ease-out forwards;
    }
  </style>