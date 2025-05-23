<script lang="ts">
   import { onMount } from 'svelte';
   import { toast } from 'svelte-sonner';
   import { Home, Play, BookOpen, FileText, Zap, Trophy, Clock, ChevronRight, Star, Target, Calendar, Code2, Flame, Award, TrendingUp, Users, BookMarked, CloudLightning } from 'lucide-svelte';
   import { cn } from "$lib/utils";
   import { Button } from "$lib/components/ui/button";
   import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
   import { Progress } from "$lib/components/ui/progress";
   import { Badge } from "$lib/components/ui/badge";
   import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar";

   let currentTab = 'home';

   const user = {
       name: "Alex Johnson",
       avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face",
       streak: 12,
       totalXP: 2840,
       level: 8,
       rank: "Advanced Coder",
       joinDate: "March 2024"
   };

   const currentCourse = {
       title: "Complete C Programming Masterclass",
       description: "Master C programming from basics to advanced concepts with hands-on projects",
       progress: 68,
       totalLessons: 32,
       completedLessons: 22,
       estimatedTime: "4h 15m remaining",
       difficulty: "Intermediate",
       nextLesson: "Pointers and Memory Management",
       category: "Systems Programming",
       instructor: "Dr. Sarah Chen"
   };

   const recentCourses = [
       {
           title: "C Programming Fundamentals",
           progress: 100,
           totalLessons: 18,
           completedLessons: 18,
           badge: "Completed",
           color: "green",
           rating: 4.9
       },
       {
           title: "Data Structures in C",
           progress: 45,
           totalLessons: 24,
           completedLessons: 11,
           badge: "In Progress",
           color: "blue",
           rating: 4.8
       },
       {
           title: "Advanced C Concepts",
           progress: 0,
           totalLessons: 20,
           completedLessons: 0,
           badge: "Locked",
           color: "gray",
           rating: 4.9
       }
   ];

   const achievements = [
       { title: "C Master", description: "Complete 3 C programming courses", icon: Code2, earned: true, date: "2 days ago", rarity: "Epic" },
       { title: "Speed Coder", description: "Write 100 lines of C code in one session", icon: CloudLightning, earned: true, date: "1 week ago", rarity: "Rare" },
       { title: "Memory Manager", description: "Master pointer concepts", icon: Target, earned: false, progress: 75, rarity: "Legendary" },
       { title: "Problem Crusher", description: "Solve 50 C coding challenges", icon: Award, earned: false, progress: 42, rarity: "Epic" }
   ];

   const upcomingLessons = [
       { title: "Pointers and Memory Management", duration: "35 min", type: "Interactive", difficulty: "Intermediate", icon: Target },
       { title: "Dynamic Memory Allocation", duration: "28 min", type: "Hands-on", difficulty: "Advanced", icon: Zap },
       { title: "File I/O Operations", duration: "22 min", type: "Project", difficulty: "Intermediate", icon: FileText },
       { title: "Data Structures Implementation", duration: "45 min", type: "Workshop", difficulty: "Advanced", icon: BookMarked }
   ];

   const stats = [
       { label: "Total XP", value: user.totalXP.toLocaleString(), icon: Star, color: "from-yellow-400 to-orange-500", bg: "bg-yellow-500/10" },
       { label: "Current Streak", value: `${user.streak} days`, icon: Flame, color: "from-orange-400 to-red-500", bg: "bg-orange-500/10" },
       { label: "Courses Completed", value: "3", icon: Trophy, color: "from-green-400 to-emerald-500", bg: "bg-green-500/10" },
       { label: "Code Written", value: "12.5k", icon: Code2, color: "from-blue-400 to-cyan-500", bg: "bg-blue-500/10" }
   ];

   // BorderBeam Component
   let BorderBeam = {
       size: 200,
       duration: 15,
       colorFrom: "#3b82f6",
       colorTo: "#8b5cf6"
   };
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 text-white">
   <!-- Animated Background -->
   <div class="fixed inset-0 overflow-hidden pointer-events-none">
       <div class="absolute -top-1/2 -left-1/2 w-full h-full bg-gradient-to-br from-blue-500/5 via-purple-500/5 to-transparent rounded-full blur-3xl animate-pulse"></div>
       <div class="absolute -bottom-1/2 -right-1/2 w-full h-full bg-gradient-to-tl from-cyan-500/5 via-blue-500/5 to-transparent rounded-full blur-3xl animate-pulse" style="animation-delay: 2s;"></div>
   </div>

   <div class="relative z-10 flex">
       <!-- Sidebar -->
       <aside class="w-72 border-r border-slate-700/50 bg-slate-800/30 backdrop-blur-xl">
           <!-- Brand Header -->
           <div class="p-6 border-b border-slate-700/50">
               <div class="flex items-center space-x-3">
                   <div class="relative">
                       <div class="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl blur opacity-75"></div>
                       <div class="relative p-3 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl">
                           <Code2 class="h-6 w-6 text-white" />
                       </div>
                   </div>
                   <div>
                       <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">Codepass</h1>
                       <p class="text-xs text-slate-400">Master. Build. Excel.</p>
                   </div>
               </div>
           </div>

           <!-- User Profile -->
           <div class="p-6 border-b border-slate-700/50">
               <div class="flex items-center space-x-4">
                   <div class="relative">
                       <Avatar class="h-12 w-12 ring-2 ring-blue-500/50">
                           <AvatarImage src={user.avatar} alt={user.name} />
                           <AvatarFallback class="bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold">
                               {user.name.split(' ').map(n => n[0]).join('')}
                           </AvatarFallback>
                       </Avatar>
                       <div class="absolute -bottom-1 -right-1 w-5 h-5 bg-green-500 rounded-full border-2 border-slate-800 flex items-center justify-center">
                           <div class="w-2 h-2 bg-white rounded-full"></div>
                       </div>
                   </div>
                   <div>
                       <p class="font-semibold text-white">{user.name}</p>
                       <p class="text-sm text-blue-400">{user.rank}</p>
                       <p class="text-xs text-slate-400">Level {user.level}</p>
                   </div>
               </div>
           </div>

           <!-- Navigation -->
           <nav class="p-6 space-y-2">
               {#each [
                   { id: 'home', label: 'Dashboard', icon: Home, count: null },
                   { id: 'courses', label: 'My Courses', icon: BookOpen, count: 3 },
                   { id: 'practice', label: 'Code Practice', icon: Code2, count: 15 },
                   { id: 'community', label: 'Community', icon: Users, count: null }
               ] as item}
                   <button
                       class={cn(
                           "w-full flex items-center justify-between px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 group",
                           currentTab === item.id 
                               ? "bg-gradient-to-r from-blue-500/20 to-purple-500/20 text-white border border-blue-500/30 shadow-lg shadow-blue-500/10" 
                               : "text-slate-300 hover:text-white hover:bg-slate-700/30 hover:translate-x-1"
                       )}
                       on:click={() => currentTab = item.id}
                   >
                       <div class="flex items-center space-x-3">
                           <svelte:component this={item.icon} class={cn("h-5 w-5 transition-colors", currentTab === item.id ? "text-blue-400" : "text-slate-400 group-hover:text-blue-400")} />
                           <span>{item.label}</span>
                       </div>
                       {#if item.count}
                           <Badge variant="secondary" class="bg-slate-700 text-slate-300 text-xs">
                               {item.count}
                           </Badge>
                       {/if}
                   </button>
               {/each}
           </nav>

           <!-- Weekly Progress -->
           <div class="p-6 mt-auto">
               <Card class="bg-slate-700/30 border-slate-600/50 backdrop-blur-sm">
                   <CardContent class="p-4 space-y-4">
                       <div class="flex items-center justify-between">
                           <span class="text-sm font-medium text-slate-300">Weekly Goal</span>
                           <span class="text-sm font-bold text-white">4/7 days</span>
                       </div>
                       <div class="space-y-2">
                           <Progress value={57} class="h-2 bg-slate-600" />
                           <div class="flex justify-between text-xs text-slate-400">
                               <span>Keep it up!</span>
                               <span>57%</span>
                           </div>
                       </div>
                   </CardContent>
               </Card>
           </div>
       </aside>

       <!-- Main Content -->
       <main class="flex-1 overflow-y-auto">
           {#if currentTab === 'home'}
               <div class="p-8 space-y-8 max-w-7xl mx-auto">
                   <!-- Header -->
                   <div class="flex items-center justify-between">
                       <div>
                           <h1 class="text-4xl font-bold text-white mb-2">
                               Welcome back, <span class="bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">{user.name.split(' ')[0]}!</span> 👋
                           </h1>
                           <p class="text-slate-400 text-lg">Ready to continue your C programming journey?</p>
                       </div>
                       <div class="flex items-center space-x-4">
                           <div class="flex items-center space-x-2 bg-gradient-to-r from-orange-500/20 to-red-500/20 backdrop-blur-sm text-orange-300 px-4 py-2 rounded-full border border-orange-500/30">
                               <Flame class="h-5 w-5" />
                               <span class="font-bold">{user.streak} day streak</span>
                           </div>
                           <Badge class="bg-gradient-to-r from-blue-500/20 to-purple-500/20 text-blue-300 border-blue-500/30 px-4 py-2 text-sm font-semibold">
                               Level {user.level}
                           </Badge>
                       </div>
                   </div>

                   <!-- Stats Grid -->
                   <div class="grid grid-cols-4 gap-6">
                       {#each stats as stat}
                           <Card class="relative overflow-hidden bg-slate-800/40 border-slate-700/50 backdrop-blur-sm hover:bg-slate-800/60 transition-all duration-300 group">
                               <!-- BorderBeam Effect -->
                               <div class="absolute inset-[0] rounded-[inherit] [border:1.5px_solid_transparent] ![mask-clip:padding-box,border-box] ![mask-composite:intersect] [mask:linear-gradient(transparent,transparent),linear-gradient(white,white)] after:absolute after:aspect-square after:w-[100px] after:animate-spin after:[animation-duration:8s] after:bg-gradient-to-r after:from-blue-500 after:to-purple-500 after:[offset-anchor:50%_50%] after:[offset-path:rect(0_auto_auto_0_round_8px)] opacity-60"></div>
                               
                               <CardContent class="relative p-6">
                                   <div class="flex items-center space-x-4">
                                       <div class={cn("p-3 rounded-xl bg-gradient-to-r", stat.color, stat.bg, "shadow-lg")}>
                                           <svelte:component this={stat.icon} class="h-6 w-6 text-white" />
                                       </div>
                                       <div>
                                           <p class="text-3xl font-bold text-white group-hover:scale-105 transition-transform">{stat.value}</p>
                                           <p class="text-sm text-slate-400 font-medium">{stat.label}</p>
                                       </div>
                                   </div>
                               </CardContent>
                           </Card>
                       {/each}
                   </div>

                   <!-- Main Content Grid -->
                   <div class="grid grid-cols-3 gap-8">
                       <!-- Left Column - Current Course -->
                       <div class="col-span-2 space-y-8">
                           <!-- Featured Course Card -->
                           <Card class="relative overflow-hidden bg-gradient-to-br from-slate-800/50 via-slate-800/30 to-slate-700/50 border-slate-600/50 backdrop-blur-sm">
                               <!-- Animated Background -->
                               <div class="absolute inset-0 bg-gradient-to-r from-blue-500/5 via-purple-500/5 to-cyan-500/5"></div>
                               <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-500 via-purple-500 to-cyan-500"></div>
                               
                               <!-- BorderBeam -->
                               <div class="absolute inset-[0] rounded-[inherit] [border:2px_solid_transparent] ![mask-clip:padding-box,border-box] ![mask-composite:intersect] [mask:linear-gradient(transparent,transparent),linear-gradient(white,white)] after:absolute after:aspect-square after:w-[200px] after:animate-spin after:[animation-duration:15s] after:bg-gradient-to-r after:from-blue-400 after:to-purple-600 after:[offset-anchor:50%_50%] after:[offset-path:rect(0_auto_auto_0_round_12px)]"></div>

                               <CardHeader class="relative z-10 pb-4">
                                   <div class="flex items-center justify-between">
                                       <div>
                                           <CardTitle class="text-2xl flex items-center gap-3 text-white">
                                               <div class="p-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg">
                                                   <BookOpen class="h-6 w-6 text-white" />
                                               </div>
                                               Continue Learning
                                           </CardTitle>
                                           <CardDescription class="text-slate-300 mt-2">Pick up where you left off in your C journey</CardDescription>
                                       </div>
                                       <Badge class="bg-gradient-to-r from-emerald-500/20 to-cyan-500/20 text-emerald-300 border-emerald-500/30 px-3 py-1">
                                           {currentCourse.difficulty}
                                       </Badge>
                                   </div>
                               </CardHeader>

                               <CardContent class="relative z-10 space-y-6">
                                   <div>
                                       <h3 class="text-3xl font-bold mb-3 text-white">{currentCourse.title}</h3>
                                       <p class="text-slate-300 mb-2">{currentCourse.description}</p>
                                       <p class="text-sm text-slate-400">Instructor: <span class="text-blue-400 font-medium">{currentCourse.instructor}</span></p>
                                   </div>
                                   
                                   <div class="grid grid-cols-2 gap-6">
                                       <div class="space-y-3">
                                           <div class="flex justify-between items-center">
                                               <span class="text-slate-300 font-medium">Progress</span>
                                               <span class="text-xl font-bold text-blue-400">{currentCourse.progress}%</span>
                                           </div>
                                           <div class="space-y-2">
                                               <Progress value={currentCourse.progress} class="h-3 bg-slate-600" />
                                               <div class="flex justify-between text-sm text-slate-400">
                                                   <span>{currentCourse.completedLessons} completed</span>
                                                   <span>{currentCourse.totalLessons} total</span>
                                               </div>
                                           </div>
                                       </div>
                                       <div class="space-y-3">
                                           <p class="text-slate-300 font-medium">Time Remaining</p>
                                           <div class="flex items-center space-x-2">
                                               <Clock class="h-5 w-5 text-blue-400" />
                                               <span class="text-xl font-bold text-white">{currentCourse.estimatedTime}</span>
                                           </div>
                                       </div>
                                   </div>

                                   <div class="bg-slate-700/50 backdrop-blur-sm p-6 rounded-xl border border-slate-600/50">
                                       <div class="flex items-center justify-between">
                                           <div>
                                               <p class="font-bold text-white text-lg">Next: {currentCourse.nextLesson}</p>
                                               <p class="text-slate-400 mt-1">Master one of C's most important concepts</p>
                                           </div>
                                           <Button class="bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white px-6 py-3 rounded-xl font-semibold shadow-lg hover:shadow-blue-500/25 transition-all duration-300 hover:scale-105">
                                               Continue Learning
                                               <ChevronRight class="ml-2 h-5 w-5" />
                                           </Button>
                                       </div>
                                   </div>
                               </CardContent>
                           </Card>

                           <!-- Upcoming Lessons -->
                           <Card class="bg-slate-800/40 border-slate-700/50 backdrop-blur-sm">
                               <CardHeader>
                                   <CardTitle class="flex items-center gap-3 text-white text-xl">
                                       <Calendar class="h-6 w-6 text-blue-400" />
                                       Upcoming Lessons
                                   </CardTitle>
                                   <CardDescription class="text-slate-300">Your personalized C programming curriculum</CardDescription>
                               </CardHeader>
                               <CardContent>
                                   <div class="space-y-4">
                                       {#each upcomingLessons as lesson, index}
                                           <div class="flex items-center justify-between p-4 rounded-xl bg-slate-700/30 border border-slate-600/50 hover:bg-slate-700/50 transition-all duration-300 cursor-pointer group hover:scale-[1.02]">
                                               <div class="flex items-center space-x-4">
                                                   <div class="relative">
                                                       <div class="w-12 h-12 rounded-xl bg-gradient-to-r from-blue-500/20 to-purple-500/20 flex items-center justify-center border border-blue-500/30">
                                                           <svelte:component this={lesson.icon} class="h-5 w-5 text-blue-400" />
                                                       </div>
                                                       <div class="absolute -top-1 -right-1 w-6 h-6 bg-slate-600 rounded-full flex items-center justify-center text-xs font-bold text-white">
                                                           {index + 1}
                                                       </div>
                                                   </div>
                                                   <div>
                                                       <p class="font-semibold text-white group-hover:text-blue-400 transition-colors">{lesson.title}</p>
                                                       <div class="flex items-center space-x-3 text-sm text-slate-400 mt-1">
                                                           <span class="flex items-center gap-1">
                                                               <Clock class="h-3 w-3" />
                                                               {lesson.duration}
                                                           </span>
                                                           <span>•</span>
                                                           <span>{lesson.type}</span>
                                                           <Badge variant="outline" class="text-xs bg-slate-600/50 text-slate-300 border-slate-500">
                                                               {lesson.difficulty}
                                                           </Badge>
                                                       </div>
                                                   </div>
                                               </div>
                                               <ChevronRight class="h-5 w-5 text-slate-400 group-hover:text-blue-400 transition-colors" />
                                           </div>
                                       {/each}
                                   </div>
                               </CardContent>
                           </Card>
                       </div>

                       <!-- Right Column -->
                       <div class="space-y-8">
                           <!-- Progress Overview -->
                           <Card class="bg-slate-800/40 border-slate-700/50 backdrop-blur-sm">
                               <CardHeader>
                                   <CardTitle class="flex items-center gap-3 text-white">
                                       <TrendingUp class="h-5 w-5 text-green-400" />
                                       Learning Progress
                                   </CardTitle>
                                   <CardDescription class="text-slate-300">Your C programming journey</CardDescription>
                               </CardHeader>
                               <CardContent class="space-y-6">
                                   {#each recentCourses as course}
                                       <div class="space-y-3">
                                           <div class="flex items-center justify-between">
                                               <div>
                                                   <p class="font-semibold text-white text-sm">{course.title}</p>
                                                   <div class="flex items-center space-x-2 mt-1">
                                                       <div class="flex items-center">
                                                           {#each Array(5) as _, i}
                                                               <Star class={cn("h-3 w-3", i < Math.floor(course.rating) ? "text-yellow-400 fill-current" : "text-slate-600")} />
                                                           {/each}
                                                       </div>
                                                       <span class="text-xs text-slate-400">{course.rating}</span>
                                                   </div>
                                               </div>
                                               <Badge 
                                                   class={cn(
                                                       "text-xs font-medium",
                                                       course.color === 'green' && "bg-green-500/20 text-green-300 border-green-500/30",
                                                       course.color === 'blue' && "bg-blue-500/20 text-blue-300 border-blue-500/30",
                                                       course.color === 'gray' && "bg-slate-500/20 text-slate-400 border-slate-500/30"
                                                   )}
                                               >
                                                   {course.badge}
                                               </Badge>
                                           </div>
                                           {#if course.progress > 0}
                                               <div class="space-y-2">
                                                   <Progress value={course.progress} class="h-2 bg-slate-600" />
                                                   <p class="text-xs text-slate-400 flex justify-between">
                                                       <span>{course.completedLessons}/{course.totalLessons} lessons</span>
                                                       <span>{course.progress}%</span>
                                                   </p>
                                               </div>
                                           {:else}
                                               <p class="text-xs text-slate-500">Complete prerequisites to unlock</p>
                                           {/if}
                                       </div>
                                   {/each}
                               </CardContent>
                           </Card>

                           <!-- Achievements -->
                           <Card class="bg-slate-800/40 border-slate-700/50 backdrop-blur-sm">
                               <CardHeader>
                                   <CardTitle class="flex items-center gap-3 text-white">
                                       <Award class="h-5 w-5 text-yellow-400" />
                                       Achievements
                                   </CardTitle>
                                   <CardDescription class="text-slate-300">Your coding milestones</CardDescription>
                               </CardHeader>
                               <CardContent class="space-y-4">
                                   {#each achievements as achievement}
                                       <div class={cn(
                                           "flex items-start space-x-4 p-4 rounded-xl border transition-all duration-300",
                                           achievement.earned 
                                               ? "bg-gradient-to-r from-green-500/10 to-emerald-500/10 border-green-500/30 shadow-lg shadow-green-500/10" 
                                               : "bg-slate-700/30 border-slate-600/50 hover:bg-slate-700/50"
                                       )}>
                                           <div class={cn(
                                               "p-3 rounded-xl flex-shrink-0",
                                               achievement.earned 
                                                   ? "bg-gradient-to-r from-green-500 to-emerald-500 text-white shadow-lg" 
                                                   : "bg-slate-600 text-slate-300"
                                           )}>
                                               <svelte:component this={achievement.icon} class="h-5 w-5" />
                                           </div>
                                           <div class="flex-1 min-w-0">
                                               <div class="flex items-center gap-2">
                                                   <p class="font-semibold text-white text-sm">{achievement.title}</p>
                                                   <Badge class={cn(
                                                       "text-xs",
                                                       achievement.rarity === 'Legendary' && "bg-purple-500/20 text-purple-300 border-purple-500/30",
                                                       achievement.rarity === 'Epic' && "bg-orange-500/20 text-orange-300 border-orange-500/30",
                                                       achievement.rarity === 'Rare' && "bg-blue-500/20 text-blue-300 border-blue-500/30"
                                                   )}>
                                                       {achievement.rarity}
                                                   </Badge>
                                               </div>
                                               <p class="text-xs text-slate-400 mb-3">{achievement.description}</p>
                                               {#if achievement.earned}
                                                   <p class="text-xs text-green-400 font-medium">Earned {achievement.date}</p>
                                               {:else if achievement.progress}
                                                   <div class="space-y-2">
                                                       <div class="flex justify-between text-xs">
                                                           <span class="text-slate-400">Progress</span>
                                                           <span class="font-semibold text-white">{achievement.progress}/{achievement.title === 'Memory Manager' ? '100' : '50'}</span>
                                                       </div>
                                                       <Progress 
                                                           value={(achievement.progress / (achievement.title === 'Memory Manager' ? 100 : 50)) * 100} 
                                                           class="h-2 bg-slate-600" 
                                                       />
                                                   </div>
                                               {/if}
                                           </div>
                                       </div>
                                   {/each}
                               </CardContent>
                           </Card>
                       </div>
                   </div>
               </div>
           {:else}
               <!-- Coming Soon Pages -->
               <div class="p-8 flex items-center justify-center min-h-full">
                   <div class="text-center space-y-6 max-w-md">
                       <div class="relative">
                           <div class="w-20 h-20 bg-gradient-to-r from-slate-700 to-slate-600 rounded-2xl flex items-center justify-center mx-auto shadow-2xl">
                               <svelte:component this={currentTab === 'courses' ? BookOpen : currentTab === 'practice' ? Code2 : Users} class="h-10 w-10 text-slate-300" />
                           </div>
                           <div class="absolute inset-0 bg-gradient-to-r from-blue-500/20 to-purple-500/20 rounded-2xl blur-xl"></div>
                       </div>
                       <div>
                           <h2 class="text-3xl font-bold text-white mb-2">Coming Soon</h2>
                           <p class="text-slate-400 text-lg">
                               {currentTab === 'courses' ? 'Advanced Course Library' : 
                                currentTab === 'practice' ? 'Interactive Code Playground' : 
                                'Learning Community Hub'} 
                               is under development
                           </p>
                       </div>
                       <div class="space-y-4">
                          <div class="flex items-center justify-center space-x-2 text-sm text-slate-400">
                              <div class="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
                              <div class="w-2 h-2 bg-purple-500 rounded-full animate-pulse" style="animation-delay: 0.5s;"></div>
                              <div class="w-2 h-2 bg-cyan-500 rounded-full animate-pulse" style="animation-delay: 1s;"></div>
                          </div>
                          <Button 
                              variant="outline" 
                              class="bg-slate-800/50 border-slate-600 text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-300 px-6 py-3 rounded-xl font-medium"
                              on:click={() => currentTab = 'home'}
                          >
                              Back to Dashboard
                          </Button>
                      </div>
                  </div>
              </div>
          {/if}
      </main>
  </div>
</div>

<style>
  @keyframes meteor {
      0% {
          transform: rotate(215deg) translateX(0);
          opacity: 1;
      }
      70% {
          opacity: 1;
      }
      100% {
          transform: rotate(215deg) translateX(-500px);
          opacity: 0;
      }
  }

  @keyframes border-beam {
      0% {
          offset-distance: 0%;
      }
      100% {
          offset-distance: 100%;
      }
  }

  @keyframes spin {
      from {
          transform: rotate(0deg);
      }
      to {
          transform: rotate(360deg);
      }
  }

  .animate-meteor {
      animation: meteor 4s linear infinite;
  }

  .animate-border-beam {
      animation: border-beam 15s linear infinite;
  }

  .animate-spin {
      animation: spin 1s linear infinite;
  }

  /* Custom scrollbar */
  :global(::-webkit-scrollbar) {
      width: 8px;
  }

  :global(::-webkit-scrollbar-track) {
      background: rgba(51, 65, 85, 0.3);
      border-radius: 4px;
  }

  :global(::-webkit-scrollbar-thumb) {
      background: rgba(100, 116, 139, 0.5);
      border-radius: 4px;
  }

  :global(::-webkit-scrollbar-thumb:hover) {
      background: rgba(100, 116, 139, 0.7);
  }

  /* Global styles */
  :global(body) {
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      margin: 0;
      padding: 0;
      background: #0f172a;
  }

  /* Enhanced button hover effects */
  :global(.hover-glow:hover) {
      box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
  }

  /* Custom progress bar styling */
  :global(.progress-glow) {
      background: linear-gradient(90deg, #3b82f6, #8b5cf6, #06b6d4);
      background-size: 200% 100%;
      animation: gradient-shift 3s ease-in-out infinite;
  }

  @keyframes gradient-shift {
      0%, 100% {
          background-position: 0% 50%;
      }
      50% {
          background-position: 100% 50%;
      }
  }

  /* Floating animation for cards */
  @keyframes float {
      0%, 100% {
          transform: translateY(0px);
      }
      50% {
          transform: translateY(-5px);
      }
  }

  .float-animation {
      animation: float 6s ease-in-out infinite;
  }

  /* Pulse glow effect */
  @keyframes pulse-glow {
      0%, 100% {
          box-shadow: 0 0 5px rgba(59, 130, 246, 0.2);
      }
      50% {
          box-shadow: 0 0 20px rgba(59, 130, 246, 0.4), 0 0 30px rgba(139, 92, 246, 0.3);
      }
  }

  .pulse-glow {
      animation: pulse-glow 2s ease-in-out infinite;
  }

  /* Shimmer effect for loading states */
  @keyframes shimmer {
      0% {
          background-position: -200px 0;
      }
      100% {
          background-position: calc(200px + 100%) 0;
      }
  }

  .shimmer {
      background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
      background-size: 200px 100%;
      animation: shimmer 2s infinite;
  }

  /* Glass morphism effects */
  .glass-effect {
      background: rgba(30, 41, 59, 0.1);
      backdrop-filter: blur(10px);
      -webkit-backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
  }

  /* Gradient text animations */
  @keyframes gradient-text {
      0%, 100% {
          background-position: 0% 50%;
      }
      50% {
          background-position: 100% 50%;
      }
  }

  .animate-gradient-text {
      background: linear-gradient(-45deg, #3b82f6, #8b5cf6, #06b6d4, #10b981);
      background-size: 400% 400%;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      animation: gradient-text 3s ease infinite;
  }

  /* Enhanced hover states */
  .card-hover {
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .card-hover:hover {
      transform: translateY(-4px) scale(1.02);
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                  0 10px 10px -5px rgba(0, 0, 0, 0.04),
                  0 0 0 1px rgba(59, 130, 246, 0.1);
  }

  /* Improved focus states for accessibility */
  :global(button:focus-visible) {
      outline: 2px solid #3b82f6;
      outline-offset: 2px;
  }

  /* Better responsive behavior */
  @media (max-width: 1024px) {
      .grid-cols-4 {
          grid-template-columns: repeat(2, minmax(0, 1fr));
      }
      
      .grid-cols-3 {
          grid-template-columns: 1fr;
      }
      
      .col-span-2 {
          grid-column: span 1 / span 1;
      }
  }

  @media (max-width: 768px) {
      .grid-cols-4 {
          grid-template-columns: 1fr;
      }
      
      aside {
          width: 100%;
          position: fixed;
          bottom: 0;
          height: auto;
          z-index: 50;
          border-top: 1px solid rgba(100, 116, 139, 0.3);
          border-right: none;
      }
      
      main {
          padding-bottom: 80px;
      }
  }

  /* Loading states */
  .loading-skeleton {
      background: linear-gradient(90deg, #1e293b 25%, #334155 50%, #1e293b 75%);
      background-size: 200% 100%;
      animation: loading 1.5s infinite;
  }

  @keyframes loading {
      0% {
          background-position: 200% 0;
      }
      100% {
          background-position: -200% 0;
      }
  }

  /* Custom selection */
  :global(::selection) {
      background: rgba(59, 130, 246, 0.3);
      color: white;
  }

  /* Enhanced focus indicators */
  :global(.focus-ring:focus) {
      box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.5);
  }

  /* Smooth page transitions */
  .page-transition {
      animation: fadeInUp 0.5s ease-out;
  }

  @keyframes fadeInUp {
      from {
          opacity: 0;
          transform: translateY(20px);
      }
      to {
          opacity: 1;
          transform: translateY(0);
      }
  }

  /* Interactive hover states for better UX */
  .interactive-hover {
      position: relative;
      overflow: hidden;
  }

  .interactive-hover::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
      transition: left 0.5s;
  }

  .interactive-hover:hover::before {
      left: 100%;
  }

  /* Achievement rarity colors */
  .rarity-common {
      border-color: #6b7280;
      background: rgba(107, 114, 128, 0.1);
  }

  .rarity-rare {
      border-color: #3b82f6;
      background: rgba(59, 130, 246, 0.1);
  }

  .rarity-epic {
      border-color: #a855f7;
      background: rgba(168, 85, 247, 0.1);
  }

  .rarity-legendary {
      border-color: #f59e0b;
      background: rgba(245, 158, 11, 0.1);
  }

  /* Enhanced animations for better visual feedback */
  .bounce-in {
      animation: bounceIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55);
  }

  @keyframes bounceIn {
      0% {
          transform: scale(0.3);
          opacity: 0;
      }
      50% {
          transform: scale(1.05);
          opacity: 0.8;
      }
      70% {
          transform: scale(0.9);
          opacity: 1;
      }
      100% {
          transform: scale(1);
          opacity: 1;
      }
  }

  /* Progress bar enhancements */
  :global(.progress-bar) {
      background: linear-gradient(90deg, #3b82f6, #8b5cf6);
      border-radius: 9999px;
      position: relative;
      overflow: hidden;
  }

  :global(.progress-bar::after) {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
      background-image: linear-gradient(
          -45deg,
          rgba(255, 255, 255, 0.2) 25%,
          transparent 25%,
          transparent 50%,
          rgba(255, 255, 255, 0.2) 50%,
          rgba(255, 255, 255, 0.2) 75%,
          transparent 75%,
          transparent
      );
      background-size: 30px 30px;
      animation: move 2s linear infinite;
  }

  @keyframes move {
      0% {
          background-position: 0 0;
      }
      100% {
          background-position: 30px 30px;
      }
  }
</style>