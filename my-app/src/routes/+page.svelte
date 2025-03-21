<script lang="ts">
    import { onMount } from 'svelte';
    

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
    let theme: string = 'github-dark';
    

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
        const response = await fetch('/api/challenges');
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
        const response = await fetch(`/api/challenge/${id}`);
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
        
        const response = await fetch('/api/submit', {
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
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
    
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
  
  <div class="min-h-screen bg-slate-100 flex flex-col">
    <header class="bg-gradient-to-r from-slate-800 to-slate-900 text-white shadow-md">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <div class="text-2xl font-bold tracking-tight flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
          </svg>
          StrathLearn
        </div>
        <nav>
          <ul class="flex space-x-8">
            <li><a href="/" class="hover:text-blue-300 font-medium text-blue-100 transition">Challenges</a></li>
            <li><a href="#" class="hover:text-blue-300 text-white/80 transition">Leaderboard</a></li>
            <li><a href="#" class="hover:text-blue-300 text-white/80 transition">About</a></li>
          </ul>
        </nav>
      </div>
    </header>
  
    <main class="flex-grow">
      <div class="container mx-auto px-4 py-8">
        <div class="mb-6 flex items-center justify-between flex-wrap gap-4">
          <div class="flex-1 max-w-xs">
            <label for="challenge-select" class="block text-sm font-medium text-slate-700 mb-2">Select Challenge:</label>
            <select 
              id="challenge-select"
              class="block w-full px-3 py-2 border border-slate-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white"
              bind:value={selectedChallengeId}
              on:change={handleChallengeSelect}
            >
              {#each Object.entries(challenges) as [id, challenge]}
                <option value={id}>{challenge.title}</option>
              {/each}
            </select>
          </div>
          
          <div class="flex-1 max-w-xs">
            <label for="theme-select" class="block text-sm font-medium text-slate-700 mb-2">Editor Theme:</label>
            <select 
              id="theme-select"
              class="block w-full px-3 py-2 border border-slate-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white"
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
            <div class="lg:col-span-2 bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
              <div class="bg-gradient-to-r from-slate-50 to-white p-6 border-b border-slate-200">
                <h1 class="text-2xl font-bold mb-2 text-slate-800">{currentChallenge.title}</h1>
                <div class="mb-4">
                  <span class="inline-block px-2.5 py-0.5 rounded-full text-xs font-semibold 
                    {currentChallenge.difficulty === 'Easy' ? 'bg-green-100 text-green-800' : 
                    currentChallenge.difficulty === 'Medium' ? 'bg-yellow-100 text-yellow-800' : 
                    'bg-red-100 text-red-800'}">
                    {currentChallenge.difficulty}
                  </span>
                </div>
                
                <div class="prose prose-slate mb-6 max-w-none">
                  {@html currentChallenge.description.replace(/\n/g, '<br>')}
                </div>
                
       
                <div class="mb-4 border rounded-lg overflow-hidden shadow-sm">
                  <button 
                    class="w-full flex justify-between items-center p-3 bg-blue-50 hover:bg-blue-100 text-left font-medium text-blue-700 transition-colors"
                    on:click={toggleHints}
                  >
                    <div class="flex items-center">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                      </svg>
                      Hints
                    </div>
                    <span class="text-xl">{hintsOpen ? '−' : '+'}</span>
                  </button>
                  
                  {#if hintsOpen}
                    <div class="p-4 space-y-2 bg-white">
                      {#each currentChallenge.hints as hint, i}
                        <div class="p-3 bg-yellow-50 border-l-4 border-yellow-400 text-yellow-700 rounded">
                          <p class="font-medium mb-1">Hint {i + 1}:</p>
                          {hint}
                        </div>
                      {/each}
                    </div>
                  {/if}
                </div>
                

                <div class="border rounded-lg overflow-hidden shadow-sm">
                  <button 
                    class="w-full flex justify-between items-center p-3 bg-purple-50 hover:bg-purple-100 text-left font-medium text-purple-700 transition-colors"
                    on:click={toggleTestCases}
                  >
                    <div class="flex items-center">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                      </svg>
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
              <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden flex-grow flex flex-col">
                <div class="flex justify-between items-center px-6 py-4 border-b border-slate-200 bg-slate-50">
                  <h2 class="text-lg font-semibold text-slate-800 flex items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                    </svg>
                    Your Code
                  </h2>
                  <div class="flex space-x-2 items-center">
                    <span class="text-xs text-slate-500">Ctrl+Space for autocomplete</span>
                    <button 
                      on:click={resetCode}
                      class="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded text-sm font-medium text-slate-700 transition-colors flex items-center"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                      </svg>
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
                    class="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-sm disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                  >
                    {#if submitting}
                      <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Submitting...
                    {:else}
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                      </svg>
                      Run & Submit Solution
                    {/if}
                  </button>
                </div>
                
                {#if showResults}
                  <div class="p-6 border-t border-slate-200">
                    <h3 class="text-lg font-semibold mb-4 text-slate-800 flex items-center">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                      Test Results
                    </h3>
                    
                    <div class="space-y-4">
                      {#each testResults as result, i}
                        <div class={`border rounded-lg p-4 ${result.passed ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}`}>
                          <div class="flex justify-between items-center mb-2">
                            <h4 class="font-medium flex items-center">
                              {result.passed ? 
                                '<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-green-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>' : 
                                '<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-red-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg>'
                              }
                              Test Case {i + 1}: {result.passed ? 'PASSED' : 'FAILED'}
                            </h4>
                            <span class={`text-xs font-semibold px-2 py-1 rounded-full ${result.passed ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                              {result.passed ? 'Success' : 'Error'}
                            </span>
                          </div>
                          
                          {#if result.error}
                            <div class="mt-3">
                              <p class="text-sm font-medium text-red-700 mb-1">Error Message:</p>
                              <pre class="mt-2 p-3 bg-white rounded-lg border border-red-200 text-sm overflow-x-auto whitespace-pre-wrap text-red-600 font-mono">{result.error}</pre>
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
                      <div class="mt-6 p-5 bg-green-50 border border-green-200 rounded-lg shadow-sm">
                        <div class="flex items-center">
                          <div class="flex-shrink-0">
                            <svg class="h-10 w-10 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                            </svg>
                          </div>
                          <div class="ml-4">
                            <h2 class="text-xl font-bold text-green-800">All tests passed! Congratulations!</h2>
                            <p class="text-green-700 mt-1">You've successfully completed this challenge.</p>
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
          <div class="text-center py-16 bg-white rounded-xl shadow-sm border border-slate-200">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-slate-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <p class="text-xl text-slate-600">Loading challenge...</p>
          </div>
        {/if}
      </div>
    </main>
  
    <footer class="bg-slate-800 text-white/80 py-6 mt-8">
      <div class="container mx-auto px-4 flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center space-x-4 mt-4 md:mt-0">
            <a href="#" class="text-white/80 hover:text-white transition">Terms</a>
            <a href="#" class="text-white/80 hover:text-white transition">Privacy</a>
            <a href="#" class="text-white/80 hover:text-white transition">Help</a>
          </div>
        </div>
      </footer>
    </div>
    
    <style>

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
        width: 12px;
        border-radius: 6px;
        margin-right: 4px;
      }
      
      :global(.CodeMirror-simplescroll-horizontal) {
        height: 12px;
        border-radius: 6px;
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
    </style>