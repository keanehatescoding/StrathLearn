<script lang="ts">
  import { onMount } from 'svelte';
  import { useSession } from '$lib/auth-client.js';
  import { Button } from "$lib/components/ui/button";
  import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "$lib/components/ui/card";
  import { ChevronLeft, ChevronRight, Calendar, Code, Award, Clock } from 'lucide-svelte';
  import { format, parseISO, subMonths, addMonths, startOfMonth, endOfMonth, eachDayOfInterval, getDay, isEqual, isSameMonth } from 'date-fns';
  
  const session = useSession();
  
  // State variables
  let currentDate = new Date();
  let submissions: DailySubmission[] = [];
  let loadingSubmissions = true;
  let userStats = {
    totalSubmissions: 0,
    challengesSolved: 0,
    firstSubmission: null as Date | null,
    lastSubmission: null as Date | null,
    streak: 0
  };

  // Types
  type DailySubmission = {
    date: string;
    count: number;
    challengeId?: string;
  };

  // Generate calendar dates
  function getCalendarDays(date: Date) {
    const start = startOfMonth(date);
    const end = endOfMonth(date);
    
    // Create array of dates for current month
    const days = eachDayOfInterval({ start, end });
    
    // Get the day of the week (0-6) for the first day of the month
    const firstDayOfWeek = getDay(start);
    
    // Add empty cells for days before the start of the month
    const prefixDays = Array(firstDayOfWeek).fill(null);
    
    return [...prefixDays, ...days];
  }

  // Activity level based on submission count (0-4)
  function getActivityLevel(date: string) {
    const submission = submissions.find(s => s.date === date);
    if (!submission) return 0;
    
    if (submission.count >= 10) return 4;
    if (submission.count >= 6) return 3;
    if (submission.count >= 3) return 2;
    return 1;
  }

  // Move calendar to previous/next month
  function previousMonth() {
    currentDate = subMonths(currentDate, 1);
    fetchSubmissionsForCalendar();
  }
  
  function nextMonth() {
    currentDate = addMonths(currentDate, 1);
    fetchSubmissionsForCalendar();
  }

  // Format day of month with correct suffix (1st, 2nd, etc.)
  function formatDayWithSuffix(date: Date) {
    const day = date.getDate();
    if (day > 3 && day < 21) return `${day}th`;
    switch (day % 10) {
      case 1: return `${day}st`;
      case 2: return `${day}nd`;
      case 3: return `${day}rd`;
      default: return `${day}th`;
    }
  }

  // Fetch user profile stats
  async function fetchUserProfile() {
    try {
      const response = await fetch('https://api.singularity.co.ke/api/profile');
      if (!response.ok) throw new Error('Failed to fetch profile');
      
      const data = await response.json();
      
      userStats = {
        totalSubmissions: data.stats.totalSubmissions,
        challengesSolved: data.stats.challengesSolved,
        firstSubmission: data.stats.firstSubmission ? new Date(data.stats.firstSubmission) : null,
        lastSubmission: data.stats.lastSubmission ? new Date(data.stats.lastSubmission) : null,
        streak: calculateStreak(submissions)
      };
    } catch (error) {
      console.error('Error fetching user profile:', error);
    }
  }

  // Calculate current streak from submissions
  function calculateStreak(submissions: DailySubmission[]) {
    if (!submissions.length) return 0;
    
    // Sort submissions by date descending
    const sortedDates = [...submissions]
      .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
      .map(s => s.date);
    
    // Check if latest submission is today or yesterday (to maintain streak)
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const latestSubmissionDate = new Date(sortedDates[0]);
    latestSubmissionDate.setHours(0, 0, 0, 0);
    
    const diffTime = Math.abs(today.getTime() - latestSubmissionDate.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays > 1) return 0; // Streak broken
    
    // Count consecutive days
    let streak = 1;
    for (let i = 1; i < sortedDates.length; i++) {
      const currentDate = new Date(sortedDates[i-1]);
      const prevDate = new Date(sortedDates[i]);
      
      // Check if dates are consecutive
      currentDate.setDate(currentDate.getDate() - 1);
      
      if (currentDate.getFullYear() === prevDate.getFullYear() && 
          currentDate.getMonth() === prevDate.getMonth() && 
          currentDate.getDate() === prevDate.getDate()) {
        streak++;
      } else {
        break;
      }
    }
    
    return streak;
  }

  // Fetch submissions for current calendar view
  async function fetchSubmissionsForCalendar() {
    try {
      loadingSubmissions = true;
      
      const start = startOfMonth(currentDate);
      const end = endOfMonth(currentDate);
      
      const startDate = format(start, 'yyyy-MM-dd');
      const endDate = format(end, 'yyyy-MM-dd');
      
      const response = await fetch(`https://api.singularity.co.ke/api/profile/submissions?startDate=${startDate}&endDate=${endDate}`);
      if (!response.ok) throw new Error('Failed to fetch submissions');
      
      const data = await response.json();
      submissions = data.submissions;
      
      // Update streak after getting fresh submission data
      userStats.streak = calculateStreak(submissions);
    } catch (error) {
      console.error('Error fetching submissions:', error);
    } finally {
      loadingSubmissions = false;
    }
  }

  // Initialize component
  onMount(async () => {
    await fetchUserProfile();
    await fetchSubmissionsForCalendar();
  });

  // Reactively update when month changes
  $: days = getCalendarDays(currentDate);
  $: monthYearDisplay = format(currentDate, 'MMMM yyyy');
</script>

<div class="container mx-auto px-4 py-6 max-w-4xl">
  <div class="grid gap-6 md:grid-cols-4">
    <!-- Stats cards -->
    <Card class="md:col-span-4 bg-gradient-to-br from-background to-accent/20 border-primary/30">
      <CardHeader class="pb-2">
        <CardTitle class="text-2xl font-bold bg-gradient-to-r from-primary to-primary/80 bg-clip-text text-transparent">
          Your Coding Journey
        </CardTitle>
        <CardDescription>Track your progress and consistency</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="bg-accent/20 rounded-lg p-4 border border-primary/10 flex flex-col items-center justify-center text-center">
            <div class="bg-primary/10 p-2 rounded-full mb-2">
              <Code class="h-6 w-6 text-primary" />
            </div>
            <span class="text-2xl font-bold text-foreground">{userStats.totalSubmissions}</span>
            <span class="text-xs text-muted-foreground">Total Submissions</span>
          </div>
          
          <div class="bg-accent/20 rounded-lg p-4 border border-primary/10 flex flex-col items-center justify-center text-center">
            <div class="bg-primary/10 p-2 rounded-full mb-2">
              <Award class="h-6 w-6 text-primary" />
            </div>
            <span class="text-2xl font-bold text-foreground">{userStats.challengesSolved}</span>
            <span class="text-xs text-muted-foreground">Challenges Solved</span>
          </div>
          
          <div class="bg-accent/20 rounded-lg p-4 border border-primary/10 flex flex-col items-center justify-center text-center">
            <div class="bg-primary/10 p-2 rounded-full mb-2">
              <Calendar class="h-6 w-6 text-primary" />
            </div>
            <span class="text-2xl font-bold text-foreground">{userStats.streak}</span>
            <span class="text-xs text-muted-foreground">Day Streak</span>
          </div>
          
          <div class="bg-accent/20 rounded-lg p-4 border border-primary/10 flex flex-col items-center justify-center text-center">
            <div class="bg-primary/10 p-2 rounded-full mb-2">
              <Clock class="h-6 w-6 text-primary" />
            </div>
            <span class="text-2xl font-bold text-foreground">
              {userStats.firstSubmission ? format(new Date(userStats.firstSubmission), 'MMM d') : '-'}
            </span>
            <span class="text-xs text-muted-foreground">First Submission</span>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Contribution calendar -->
    <Card class="md:col-span-4 bg-gradient-to-br from-background to-accent/10 border-primary/30">
      <CardHeader class="pb-2">
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="text-xl font-bold">Contribution Calendar</CardTitle>
            <CardDescription>Your coding activity for {monthYearDisplay}</CardDescription>
          </div>
          <div class="flex items-center space-x-2">
            <Button 
              variant="outline" 
              size="icon" 
              class="h-8 w-8 rounded-full border-primary/20 hover:bg-accent/50"
              on:click={previousMonth}
            >
              <ChevronLeft class="h-4 w-4" />
            </Button>
            <Button 
              variant="outline" 
              size="icon" 
              class="h-8 w-8 rounded-full border-primary/20 hover:bg-accent/50"
              on:click={nextMonth}
            >
              <ChevronRight class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      
      <CardContent>
        {#if loadingSubmissions}
          <div class="flex items-center justify-center h-40">
            <div class="animate-spin h-6 w-6 border-2 border-primary border-t-transparent rounded-full"></div>
          </div>
        {:else}
          <div class="space-y-4">
            <!-- Calendar grid -->
            <div class="grid grid-cols-7 gap-1">
              <!-- Week day headers -->
              {#each ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as day}
                <div class="h-8 flex items-center justify-center text-xs font-medium text-muted-foreground">
                  {day}
                </div>
              {/each}
              
              <!-- Calendar days -->
              {#each days as day, i}
                {#if day === null}
                  <div class="aspect-square rounded-md"></div>
                {:else}
                  {@const dateStr = format(day, 'yyyy-MM-dd')}
                  {@const activityLevel = getActivityLevel(dateStr)}
                  {@const today = new Date()}
                  {@const isToday = isEqual(
                    new Date(day.setHours(0, 0, 0, 0)),
                    new Date(today.setHours(0, 0, 0, 0))
                  )}
                  
                  <div 
                    class="aspect-square p-0.5 relative group"
                    aria-label="{format(day, 'EEEE, MMMM d, yyyy')}: {activityLevel ? `${submissions.find(s => s.date === dateStr)?.count || 0} submissions` : 'No submissions'}"
                  >
                    <div 
                      class={`
                        w-full h-full rounded-md flex items-center justify-center transition-all
                        ${isToday ? 'ring-2 ring-primary ring-offset-1 ring-offset-background z-10' : ''}
                        ${!isSameMonth(day, currentDate) ? 'opacity-30' : ''}
                        ${activityLevel === 0 ? 'bg-accent/30 hover:bg-accent/50' : ''}
                        ${activityLevel === 1 ? 'bg-primary/30 hover:bg-primary/40' : ''}
                        ${activityLevel === 2 ? 'bg-primary/50 hover:bg-primary/60' : ''}
                        ${activityLevel === 3 ? 'bg-primary/70 hover:bg-primary/80' : ''}
                        ${activityLevel === 4 ? 'bg-primary hover:bg-primary/90' : ''}
                      `}
                    >
                      <span class={`text-xs font-medium ${activityLevel >= 3 ? 'text-primary-foreground' : 'text-foreground'}`}>
                        {day.getDate()}
                      </span>
                    </div>
                    
                    <!-- Tooltip -->
                    <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 w-auto opacity-0 group-hover:opacity-100 transition-opacity z-50 pointer-events-none">
                      <div class="bg-background/95 backdrop-blur-sm border border-border shadow-lg rounded-md p-2 text-xs whitespace-nowrap">
                        <div class="font-medium">{format(day, 'EEEE, MMMM d')}</div>
                        {#if activityLevel > 0}
                          <div class="text-primary">{submissions.find(s => s.date === dateStr)?.count || 0} submissions</div>
                        {:else}
                          <div class="text-muted-foreground">No submissions</div>
                        {/if}
                      </div>
                      <div class="w-2 h-2 bg-background border-r border-b border-border rotate-45 absolute -bottom-1 left-1/2 transform -translate-x-1/2"></div>
                    </div>
                  </div>
                {/if}
              {/each}
            </div>
            
            <!-- Legend -->
            <div class="flex items-center justify-end space-x-2 pt-2">
              <span class="text-xs text-muted-foreground">Less</span>
              {#each [0, 1, 2, 3, 4] as level}
                <div 
                  class={`h-3 w-3 rounded-sm 
                    ${level === 0 ? 'bg-accent/30' : ''} 
                    ${level === 1 ? 'bg-primary/30' : ''} 
                    ${level === 2 ? 'bg-primary/50' : ''} 
                    ${level === 3 ? 'bg-primary/70' : ''} 
                    ${level === 4 ? 'bg-primary' : ''}`}
                ></div>
              {/each}
              <span class="text-xs text-muted-foreground">More</span>
            </div>
          </div>
        {/if}
      </CardContent>
      
      <CardFooter class="flex justify-between pt-0">
        <div class="text-xs text-muted-foreground">
          {userStats.totalSubmissions > 0 
            ? `First submission on ${userStats.firstSubmission ? format(new Date(userStats.firstSubmission), 'MMMM d, yyyy') : '-'}`
            : 'No submissions yet. Start coding to see your progress!'}
        </div>
        
        <Button 
          variant="link" 
          class="text-primary p-0 h-auto text-xs hover:text-primary/80"
          on:click={() => fetchSubmissionsForCalendar()}
        >
          Refresh
        </Button>
      </CardFooter>
    </Card>
    
    <!-- Background orbit animation -->
    <div class="fixed inset-0 z-[-1] pointer-events-none overflow-hidden">
      <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
        <div class="relative flex h-[800px] w-[800px] items-center justify-center opacity-20">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            version="1.1"
            class="pointer-events-none absolute inset-0 h-full w-full"
          >
            <circle
              class="stroke-primary/30 stroke-[1.5]"
              cx="50%"
              cy="50%"
              r="300"
              fill="none"
              stroke-dasharray="8 4"
            />
            <circle
              class="stroke-primary/20 stroke-[1.5]"
              cx="50%"
              cy="50%"
              r="400"
              fill="none"
              stroke-dasharray="8 4"
            />
          </svg>

          <div 
            class="absolute flex h-full w-full transform-gpu animate-orbit items-center justify-center rounded-full" 
            style="--duration: 30s; --radius: 300px;"
          >
            <div class="h-[50px] w-[50px] rounded-full bg-transparent border border-primary/30 flex items-center justify-center" style="transform: translateY(300px);">
              <Code class="h-6 w-6 text-primary/50" />
            </div>
          </div>
          
          <div 
            class="absolute flex h-full w-full transform-gpu animate-orbit-reverse items-center justify-center rounded-full" 
            style="animation-delay: -15s; --duration: 40s; --radius: 400px;"
          >
            <div class="h-[60px] w-[60px] rounded-full bg-primary/5 backdrop-blur-sm border border-primary/20 flex items-center justify-center" style="transform: translateY(400px);">
              <Calendar class="h-8 w-8 text-primary/40" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  @keyframes orbit {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  @keyframes orbit-reverse {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(-360deg);
    }
  }

  .animate-orbit {
    animation: orbit var(--duration, 20s) linear infinite;
    animation-delay: var(--delay, 0s);
  }

  .animate-orbit-reverse {
    animation: orbit-reverse var(--duration, 20s) linear infinite;
    animation-delay: var(--delay, 0s);
  }
</style>