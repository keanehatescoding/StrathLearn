document.addEventListener('DOMContentLoaded', () => {
    // Initialize CodeMirror
    const editor = CodeMirror.fromTextArea(document.getElementById('code-editor'), {
        mode: 'text/x-csrc',
        theme: 'dracula',
        lineNumbers: true,
        autoCloseBrackets: true,
        matchBrackets: true,
        indentUnit: 4,
        tabSize: 4,
        indentWithTabs: true,
        extraKeys: {"Tab": "indentMore", "Shift-Tab": "indentLess"}
    });

    // Current challenge data
    let currentChallenge = null;

    // Challenge selector
    const challengeSelect = document.getElementById('challenge-select');
    challengeSelect.addEventListener('change', () => {
        loadChallenge(challengeSelect.value);
    });

    // First, fetch all available challenges
    async function fetchChallenges() {
        try {
            const response = await fetch('/api/challenges');
            if (!response.ok) {
                throw new Error('Failed to load challenges');
            }
            
            const challenges = await response.json();
            
            // Clear the dropdown
            challengeSelect.innerHTML = '';
            
            // Add all challenges to the dropdown
            Object.values(challenges).forEach(challenge => {
                const option = document.createElement('option');
                option.value = challenge.id;
                option.textContent = challenge.title;
                challengeSelect.appendChild(option);
            });
            
            // Set temperature-converter as selected if available
            if (challenges['temperature-converter']) {
                challengeSelect.value = 'temperature-converter';
                loadChallenge('temperature-converter');
            } else {
                // Otherwise load the first challenge
                loadChallenge(challengeSelect.value);
            }
        } catch (error) {
            console.error('Error fetching challenges:', error);
        }
    }

    // Collapsible sections
    document.querySelectorAll('.collapsible-header').forEach(header => {
        header.addEventListener('click', function() {
            // Toggle this collapsible
            const content = this.nextElementSibling;
            const isOpen = content.classList.toggle('open');
            
            // Update the toggle icon
            const icon = this.querySelector('.toggle-icon');
            icon.textContent = isOpen ? '-' : '+';
        });
    });

    // Load challenge
    async function loadChallenge(id = 'hello-world') {
        try {
            const response = await fetch(`/api/challenge/${id}`);
            if (!response.ok) {
                throw new Error('Challenge not found');
            }
            
            currentChallenge = await response.json();
            
            // Update UI
            document.getElementById('challenge-title').textContent = currentChallenge.title;
            document.getElementById('challenge-difficulty').textContent = currentChallenge.difficulty;
            document.getElementById('challenge-description').innerHTML = currentChallenge.description.replace(/\n/g, '<br>');
            
            // Set initial code
            editor.setValue(currentChallenge.initialCode.replace(/\\n/g, '\n'));
            
            // Display hints
            const hintsContainer = document.getElementById('hints-container');
            hintsContainer.innerHTML = '';
            currentChallenge.hints.forEach(hint => {
                const hintElem = document.createElement('div');
                hintElem.className = 'hint';
                hintElem.textContent = hint;
                hintsContainer.appendChild(hintElem);
            });
            
            // Display test cases
            const testCasesContainer = document.getElementById('test-cases-container');
            testCasesContainer.innerHTML = '';
            currentChallenge.testCases.forEach(testCase => {
                if (!testCase.hidden) {
                    const testCaseElem = document.createElement('div');
                    testCaseElem.className = 'test-case';
                    testCaseElem.innerHTML = `
                        <p><strong>Input:</strong> ${testCase.input ? `"${testCase.input}"` : '<em>(empty)</em>'}</p>
                        <p><strong>Expected Output:</strong> "${testCase.expectedOutput.replace(/\\n/g, '\\\\n')}"</p>
                    `;
                    testCasesContainer.appendChild(testCaseElem);
                }
            });
            
            // Reset results
            document.getElementById('results-container').style.display = 'none';
            document.getElementById('success-message').style.display = 'none';
        } catch (error) {
            console.error('Error loading challenge:', error);
            document.getElementById('challenge-title').textContent = 'Error loading challenge';
            document.getElementById('challenge-description').textContent = error.message;
        }
    }

    // Submit code
    document.getElementById('submit-btn').addEventListener('click', async () => {
        if (!currentChallenge) {
            alert('No challenge loaded');
            return;
        }
        
        const submitBtn = document.getElementById('submit-btn');
        submitBtn.disabled = true;
        submitBtn.textContent = 'Submitting...';
        
        const code = editor.getValue();
        
        try {
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
            
            // Display results
            const resultsContainer = document.getElementById('results-container');
            resultsContainer.style.display = 'block';
            
            const testResultsContainer = document.getElementById('test-results');
            testResultsContainer.innerHTML = '';
            
            let allPassed = true;
            
            if (result.testResults && result.testResults.length > 0) {
                result.testResults.forEach(testResult => {
                    const resultElem = document.createElement('div');
                    resultElem.className = `test-result ${testResult.passed ? 'pass' : 'fail'}`;
                    
                    if (!testResult.passed) {
                        allPassed = false;
                    }
                    
                    let resultText = `<h4>Test ${testResult.testCaseId}: ${testResult.passed ? 'PASSED ✅' : 'FAILED ❌'}</h4>`;
                    
                    if (testResult.error) {
                        resultText += `<pre>${escapeHTML(testResult.error)}</pre>`;
                    }
                    
                    if (testResult.output !== undefined) {
                        resultText += `<p>Your Output:</p><pre>${escapeHTML(testResult.output)}</pre>`;
                    }
                    
                    resultElem.innerHTML = resultText;
                    testResultsContainer.appendChild(resultElem);
                });
                
                // Show success message if all tests passed
                document.getElementById('success-message').style.display = allPassed ? 'block' : 'none';
            } else {
                testResultsContainer.innerHTML = `<p>${result.message || 'No test results available'}</p>`;
            }
        } catch (error) {
            console.error('Error submitting code:', error);
            alert('Error submitting code. Please try again.');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Submit Solution';
        }
    });

    // Reset code to initial state
    document.getElementById('reset-btn').addEventListener('click', () => {
        if (currentChallenge) {
            if (confirm('Are you sure you want to reset your code to the initial state?')) {
                editor.setValue(currentChallenge.initialCode.replace(/\\n/g, '\n'));
            }
        }
    });

    function escapeHTML(str) {
        return str
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }
    
    // Start by fetching all challenges
    fetchChallenges();
});