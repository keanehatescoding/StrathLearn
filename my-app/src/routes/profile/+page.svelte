<script lang="ts">
  import { page } from "$app/stores";
  import { useSession, authClient } from "$lib/auth-client";
  import { onMount, tick } from "svelte";
  import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar";
  import { Button } from "$lib/components/ui/button";
  import { Badge } from "$lib/components/ui/badge";
  import { Tabs, TabsContent, TabsList, TabsTrigger } from "$lib/components/ui/tabs";
  import { Calendar, Code, Award, Clock, Activity, Github, Twitter, Linkedin, Mail, Edit3, ChevronRight } from "lucide-svelte";
  import { Motion } from "svelte-motion";
  import GridBeam from "$lib/components/GridBeam.svelte";
  import { format, parseISO, differenceInDays } from "date-fns";
  
  // Get user session
  const session = useSession();
  
  // User stats and submissions data
  let userStats = {
    challengesSolved: 0,
    firstSubmission: "",
    lastSubmission: "",
    totalSubmissions: 0
  };
  
  let submissions = [];
  let loadingData = true;
  let activeStreak = 0;
  let longestStreak = 0;
  let contributionData = {};
  let weeklyActivity = [];
  
  // Calculate streaks & prepare calendar data
  function processSubmissionData(data) {
    if (!data || data.length === 0) return {};
    
    // Sort submissions by date
    const sortedData = [...data].sort((a, b) => new Date(a.date) - new Date(b.date));
    
    // Get unique dates with their submission counts
    const uniqueDates = {};
    sortedData.forEach(item => {
      if (!uniqueDates[item.date]) {
        uniqueDates[item.date] = item.count;
      } else {
        uniqueDates[item.date] += item.count;
      }
    });
    
    // Calculate streaks
    let currentStreak = 0;
    let maxStreak = 0;
    let lastDate = null;
    
    Object.keys(uniqueDates).sort().forEach(date => {
      const currentDate = new Date(date);
      
      if (lastDate) {
        const dayDiff = differenceInDays(currentDate, lastDate);
        if (dayDiff === 1) {
          currentStreak++;
        } else if (dayDiff > 1) {
          maxStreak = Math.max(maxStreak, currentStreak);
          currentStreak = 1;
        }
      } else {
        currentStreak = 1;
      }
      
      lastDate = currentDate;
    });
    
    // Final check for max streak
    maxStreak = Math.max(maxStreak, currentStreak);
    
    // Calculate activity for last 7 days
    const last7Days = [];
    const today = new Date();
    for (let i = 6; i >= 0; i--) {
      const d = new Date(today);
      d.setDate(d.getDate() - i);
      const dateStr = d.toISOString().split('T')[0];
      last7Days.push({
        date: dateStr,
        count: uniqueDates[dateStr] || 0
      });
    }
    
    return {
      calendarData: uniqueDates,
      activeStreak: currentStreak,
      longestStreak: maxStreak,
      weeklyActivity: last7Days
    };
  }
  
  // Generate full year calendar grid
  function generateYearCalendar() {
    const today = new Date();
    const calendarCells = [];
    const monthLabels = [];
    
    // Go back 12 months
    const startDate = new Date(today);
    startDate.setMonth(startDate.getMonth() - 11);
    
    // Reset to first day of month
    startDate.setDate(1);
    
    // If it's not Sunday, go back to previous Sunday
    const day = startDate.getDay();
    if (day !== 0) {
      startDate.setDate(startDate.getDate() - day);
    }
    
    // Track months for labels
    let currentMonth = -1;
    
    // Generate 53 weeks (371 days) for full year view
    for (let week = 0; week < 53; week++) {
      const weekData = [];
      
      for (let dayOfWeek = 0; dayOfWeek < 7; dayOfWeek++) {
        const cellDate = new Date(startDate);
        cellDate.setDate(startDate.getDate() + (week * 7) + dayOfWeek);
        
        // Store month labels
        const month = cellDate.getMonth();
        if (month !== currentMonth) {
          currentMonth = month;
          monthLabels.push({
            week: week,
            label: format(cellDate, 'MMM')
          });
        }
        
        const dateStr = cellDate.toISOString().split('T')[0];
        const intensity = getIntensity(dateStr);
        
        weekData.push({
          date: dateStr,
          value: contributionData[dateStr] || 0,
          intensity
        });
      }
      
      calendarCells.push(weekData);
    }
    
    return { cells: calendarCells, months: monthLabels };
  }
  
  // Calculate intensity class based on submission count
  function getIntensity(date) {
    const count = contributionData[date] || 0;
    if (count === 0) return 'intensity-0';
    if (count <= 2) return 'intensity-1';
    if (count <= 5) return 'intensity-2';
    if (count <= 8) return 'intensity-3';
    return 'intensity-4';
  }
  
  async function fetchUserData() {
    loadingData = true;
    
    try {
      // In a real app, these would be API calls
      // Simulating API responses with the data you provided
      const submissionsResponse = {
        submissions: [
          { date: "2025-05-18", count: 5, challengeId: "challenge1" },
          { date: "2025-05-19", count: 5, challengeId: "challenge2" },
          { date: "2025-05-20", count: 5, challengeId: "challenge3" }
        ]
      };
      
      const statsResponse = {
        stats: {
          challengesSolved: 3,
          firstSubmission: "2025-05-18T10:00:00Z",
          lastSubmission: "2025-05-20T09:15:00Z",
          totalSubmissions: 15
        },
        user: {
          email: "wambua_kelvin@outlook.com",
          id: "rCJwbmAkfmdTC4FUOSPpptoPZ2NmjTDr",
          name: "Kelvin Wambua"
        }
      };
      
      submissions = submissionsResponse.submissions;
      userStats = statsResponse.stats;
      
      // Process submission data
      const processedData = processSubmissionData(submissions);
      contributionData = processedData.calendarData;
      activeStreak = processedData.activeStreak;
      longestStreak = processedData.longestStreak;
      weeklyActivity = processedData.weeklyActivity;
      
    } catch (error) {
      console.error("Failed to fetch user data:", error);
    } finally {
      loadingData = false;
    }
  }
  
  // Fetch data on mount
  onMount(async () => {
    await fetchUserData();
  });
  
  // Calendar data
  $: calendar = generateYearCalendar();
</script>

<div class="pt-24 pb-16 w-full">
  <!-- Hero Section with Grid Beam Background -->
  <div class="relative w-full h-[240px] dark:bg-grid-white/[0.05] bg-grid-black/[0.07] rounded-lg overflow-hidden mb-6">
    <GridBeam class="w-full h-full">
      <div class="absolute inset-0 bg-gradient-to-r from-background/80 via-background/50 to-background/30"></div>
      
      <!-- Profile Header Content -->
      <div class="relative z-10 container mx-auto px-4 h-full flex flex-col justify-end pb-6">
        <div class="flex items-end gap-6">
          <!-- Profile Avatar -->
          <Motion
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            let:motion
          >
            <div use:motion>
              <Avatar class="w-28 h-28 border-4 border-background shadow-xl">
                <AvatarImage 
                  src={$session?.data?.user?.image || "https://github.com/shadcn.png"} 
                  alt={$session?.data?.user?.name || "User"} 
                />
                <AvatarFallback class="text-3xl bg-primary text-primary-foreground">
                  {$session?.data?.user?.name ? $session.data.user.name[0].toUpperCase() : "K"}
                </AvatarFallback>
              </Avatar>
            </div>
          </Motion>
          
          <!-- User Info -->
          <Motion
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.3 }}
            let:motion
          >
            <div use:motion class="flex-1">
              <div class="flex items-center gap-3">
                <h1 class="text-3xl font-bold">{$session?.data?.user?.name || "Kelvin Wambua"}</h1>
                <Badge variant="outline" class="bg-primary/10 text-primary border-primary/20">Pro Member</Badge>
              </div>
              <p class="text-muted-foreground">{$session?.data?.user?.email || "wambua_kelvin@outlook.com"}</p>
            </div>
          </Motion>
          
          <!-- Actions -->
          <Motion
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
            let:motion
          >
            <div use:motion class="flex gap-3">
              <Button variant="outline" class="gap-2">
                <Edit3 class="h-4 w-4" /> Edit Profile
              </Button>
              <Button class="gap-2">
                <Code class="h-4 w-4" /> New Challenge
              </Button>
            </div>
          </Motion>
        </div>
      </div>
    </GridBeam>
  </div>
  
  <!-- Main Content -->
  <div class="container mx-auto px-4">
    <div class="grid lg:grid-cols-4 gap-6">
      <!-- Left Sidebar - Stats -->
      <div class="lg:col-span-1">
        <Motion
          initial={{ x: -20, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          transition={{ delay: 0.2 }}
          let:motion
        >
          <div use:motion class="space-y-6">
            <!-- Stats Cards -->
            <div class="bg-card rounded-xl border shadow-sm overflow-hidden">
              <div class="p-6 space-y-6">
                <h2 class="text-xl font-semibold">Coding Stats</h2>
                
                <div class="grid grid-cols-2 gap-4">
                  <div class="rounded-lg bg-accent/30 p-4 flex flex-col items-center justify-center text-center">
                    <Award class="h-6 w-6 text-primary mb-2" />
                    <span class="text-2xl font-bold">{userStats.challengesSolved}</span>
                    <span class="text-xs text-muted-foreground">Challenges</span>
                  </div>
                  
                  <div class="rounded-lg bg-accent/30 p-4 flex flex-col items-center justify-center text-center">
                    <Activity class="h-6 w-6 text-primary mb-2" />
                    <span class="text-2xl font-bold">{userStats.totalSubmissions}</span>
                    <span class="text-xs text-muted-foreground">Submissions</span>
                  </div>
                  
                  <div class="rounded-lg bg-accent/30 p-4 flex flex-col items-center justify-center text-center">
                    <Calendar class="h-6 w-6 text-primary mb-2" />
                    <span class="text-2xl font-bold">{activeStreak}</span>
                    <span class="text-xs text-muted-foreground">Day Streak</span>
                  </div>
                  
                  <div class="rounded-lg bg-accent/30 p-4 flex flex-col items-center justify-center text-center">
                    <Clock class="h-6 w-6 text-primary mb-2" />
                    <span class="text-2xl font-bold">{longestStreak}</span>
                    <span class="text-xs text-muted-foreground">Best Streak</span>
                  </div>
                </div>
                
                <!-- First/Last Activity -->
                <div class="space-y-3 border-t pt-4">
                  <div class="flex justify-between items-center text-sm">
                    <span class="text-muted-foreground">First submission</span>
                    <span class="font-medium">
                      {userStats.firstSubmission ? format(parseISO(userStats.firstSubmission), 'MMM d, yyyy') : 'N/A'}
                    </span>
                  </div>
                  
                  <div class="flex justify-between items-center text-sm">
                    <span class="text-muted-foreground">Last submission</span>
                    <span class="font-medium">
                      {userStats.lastSubmission ? format(parseISO(userStats.lastSubmission), 'MMM d, yyyy') : 'N/A'}
                    </span>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Weekly Activity Card -->
            <div class="bg-card rounded-xl border shadow-sm overflow-hidden">
              <div class="p-6">
                <h2 class="text-xl font-semibold mb-4">Weekly Activity</h2>
                
                <div class="space-y-3">
                  {#each weeklyActivity as day}
                    <div class="flex items-center">
                      <div class="w-14 text-xs text-muted-foreground">
                        {format(new Date(day.date), 'EEE')}
                      </div>
                      <div class="flex-1 h-4 bg-accent/30 rounded-full overflow-hidden">
                        <div 
                          class="h-full bg-primary rounded-full" 
                          style="width: {Math.min(day.count * 10, 100)}%"
                        ></div>
                      </div>
                      <div class="w-10 text-xs font-medium text-right">
                        {day.count}
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            </div>
            
            <!-- Socials Card -->
            <div class="bg-card rounded-xl border shadow-sm overflow-hidden">
              <div class="p-6">
                <h2 class="text-xl font-semibold mb-4">Connect</h2>
                
                <div class="space-y-3">
                  <Button variant="outline" class="w-full justify-start">
                    <Github class="h-4 w-4 mr-2" /> GitHub
                  </Button>
                  
                  <Button variant="outline" class="w-full justify-start">
                    <Twitter class="h-4 w-4 mr-2" /> Twitter
                  </Button>
                  
                  <Button variant="outline" class="w-full justify-start">
                    <Linkedin class="h-4 w-4 mr-2" /> LinkedIn
                  </Button>
                  
                  <Button variant="outline" class="w-full justify-start">
                    <Mail class="h-4 w-4 mr-2" /> Email
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </Motion>
      </div>
      
      <!-- Main Content Area -->
      <div class="lg:col-span-3">
        <Motion
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.3 }}
          let:motion
        >
          <div use:motion class="space-y-6">
            <!-- Contribution Calendar -->
            <div class="bg-card rounded-xl border shadow-sm overflow-hidden">
              <div class="p-6">
                <div class="flex items-center justify-between mb-6">
                  <h2 class="text-xl font-semibold">Contribution History</h2>
                  <div class="flex items-center gap-2 text-xs">
                    <span class="text-muted-foreground">Less</span>
                    <div class="flex gap-1">
                      <div class="w-3 h-3 rounded-sm bg-accent/50 intensity-0"></div>
                      <div class="w-3 h-3 rounded-sm intensity-1"></div>
                      <div class="w-3 h-3 rounded-sm intensity-2"></div>
                      <div class="w-3 h-3 rounded-sm intensity-3"></div>
                      <div class="w-3 h-3 rounded-sm intensity-4"></div>
                    </div>
                    <span class="text-muted-foreground">More</span>
                  </div>
                </div>
                
                <div class="relative overflow-x-auto pb-2">
                  <!-- Month labels -->
                  <div class="flex text-xs text-muted-foreground mb-2 pl-8">
                    {#each calendar.months as month}
                      <div class="absolute" style="left: {month.week * 16 + 32}px;">
                        {month.label}
                      </div>
                    {/each}
                  </div>
                  
                  <!-- Calendar grid -->
                  <div class="flex">
                    <!-- Day of week labels -->
                    <div class="flex flex-col gap-[3px] mr-2 text-xs text-muted-foreground mt-[3px]">
                      <span class="h-[15px]">Mon</span>
                      <span class="h-[15px]">Wed</span>
                      <span class="h-[15px]">Fri</span>
                    </div>
                    
                    <!-- Calendar grid -->
                    <div class="flex gap-[3px]">
                      {#each calendar.cells as week}
                        <div class="flex flex-col gap-[3px]">
                          {#each week as day}
                            <div 
                              class="w-[15px] h-[15px] rounded-sm {day.intensity} hover:ring-2 hover:ring-primary transition-all cursor-pointer" 
                              title="{format(new Date(day.date), 'PP')}: {day.value} submissions"
                            ></div>
                          {/each}
                        </div>
                      {/each}
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Tabs Section -->
            <Tabs defaultValue="challenges" class="w-full">
              <TabsList class="w-full justify-start border-b rounded-none h-12 bg-transparent pb-0 mb-6">
                <TabsTrigger value="challenges" class="data-[state=active]:border-b-2 data-[state=active]:border-primary rounded-none data-[state=active]:shadow-none">
                  Completed Challenges
                </TabsTrigger>
                <TabsTrigger value="skills" class="data-[state=active]:border-b-2 data-[state=active]:border-primary rounded-none data-[state=active]:shadow-none">
                  Skills & Badges
                </TabsTrigger>
                <TabsTrigger value="achievements" class="data-[state=active]:border-b-2 data-[state=active]:border-primary rounded-none data-[state=active]:shadow-none">
                  Achievements
                </TabsTrigger>
              </TabsList>
              
              <TabsContent value="challenges" class="m-0">
                <div class="grid md:grid-cols-2 gap-4">
                  {#each submissions as submission, i}
                    <Motion
                      initial={{ y: 20, opacity: 0 }}
                      animate={{ y: 0, opacity: 1 }}
                      transition={{ delay: 0.1 * i }}
                      let:motion
                    >
                      <div use:motion class="bg-card rounded-xl border shadow-sm overflow-hidden hover:shadow-md transition-all group">
                        <div class="p-5 flex items-center justify-between">
                          <div>
                            <div class="flex items-center gap-2">
                              <Badge variant="secondary" class="bg-primary/10 text-primary border-primary/20">
                                {new Date(submission.date).toLocaleDateString()}
                              </Badge>
                              <span class="text-xs text-muted-foreground">{submission.count} submissions</span>
                            </div>
                            <h3 class="text-lg font-medium mt-2">Challenge {submission.challengeId}</h3>
                            <p class="text-sm text-muted-foreground mt-1">Algorithmic problem solving</p>
                          </div>
                          
                          <Button variant="ghost" size="icon" class="opacity-0 group-hover:opacity-100 transition-opacity">
                            <ChevronRight class="h-5 w-5" />
                          </Button>
                        </div>
                      </div>
                    </Motion>
                  {/each}
                </div>
              </TabsContent>
              
              <TabsContent value="skills" class="m-0">
                <div class="grid md:grid-cols-3 gap-4">
                  {#each ["JavaScript", "Python", "React", "Data Structures", "Algorithms", "System Design"] as skill, i}
                    <Motion
                      initial={{ scale: 0.95, opacity: 0 }}
                      animate={{ scale: 1, opacity: 1 }}
                      transition={{ delay: 0.05 * i }}
                      let:motion
                    >
                      <div use:motion class="bg-card rounded-xl border shadow-sm overflow-hidden p-5 text-center hover:border-primary/50 transition-all">
                        <div class="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-3">
                          <Code class="h-6 w-6 text-primary" />
                        </div>
                        <h3 class="font-medium">{skill}</h3>
                        <div class="mt-2 h-1.5 bg-accent/50 rounded-full w-full">
                          <div class="h-full bg-primary rounded-full" style="width: {Math.floor(Math.random() * 50) + 50}%"></div>
                        </div>
                      </div>
                    </Motion>
                  {/each}
                </div>
              </TabsContent>
              
              <TabsContent value="achievements" class="m-0">
                <div class="grid md:grid-cols-3 gap-4">
                  {#each ["First Challenge", "3-Day Streak", "5 Submissions"] as achievement, i}
                    <Motion
                      initial={{ y: 10, opacity: 0 }}
                      animate={{ y: 0, opacity: 1 }}
                      transition={{ delay: 0.1 * i }}
                      let:motion
                    >
                      <div use:motion class="bg-card rounded-xl border shadow-sm overflow-hidden hover:shadow-md transition-all p-5 flex items-center gap-4">
                        <div class="min-w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                          <Award class="h-6 w-6 text-primary" />
                        </div>
                        <div>
                          <h3 class="font-medium">{achievement}</h3>
                          <p class="text-xs text-muted-foreground mt-1">Achieved on May 20, 2025</p>
                        </div>
                      </div>
                    </Motion>
                  {/each}
                </div>
              </TabsContent>
            </Tabs>
          </div>
        </Motion>
      </div>
    </div>
  </div>
</div>

<style>
  /* Contribution intensity colors */
  .intensity-0 {
    @apply bg-accent/50;
  }
  
  .intensity-1 {
    @apply bg-primary/30;
  }
  
  .intensity-2 {
    @apply bg-primary/50;
  }
  
  .intensity-3 {
    @apply bg-primary/70;
  }
  
  .intensity-4 {
    @apply bg-primary;
  }
</style>