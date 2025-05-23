<script lang="ts">
   import { onMount } from 'svelte';
   import { toast } from 'svelte-sonner';
   import { Home, Play, BookOpen, FileText, Zap, Trophy, Clock, ChevronRight, Star, Target, Calendar, Code2, Flame, Award, ArrowRight, BookOpenCheck, Users, TrendingUp, Brain, Coffee, Sparkles } from 'lucide-svelte';
   import { cn } from "$lib/utils";
   import { Button } from "$lib/components/ui/button";
   import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
   import { Progress } from "$lib/components/ui/progress";
   import { Badge } from "$lib/components/ui/badge";
   import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar";

   let meteorStyles: any = [];
   let currentTab = 'home';

   const user = {
       name: "Alex Johnson",
       avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face",
       streak: 12,
       totalXP: 2840,
       level: 8
   };

   const currentCourse = {
       title: "C Programming Fundamentals",
       description: "Master the foundations of C programming and build efficient system-level applications",
       progress: 42,
       totalLessons: 28,
       completedLessons: 12,
       estimatedTime: "8h 45m",
       difficulty: "Beginner",
       nextLesson: "Pointers and Memory Management",
       category: "Systems Programming",
       instructor: "Dr. Sarah Chen",
       students: 15240,
       rating: 4.9
   };

   const achievements = [
       { title: "First Steps", description: "Complete your first C program", icon: Zap, earned: true, date: "3 days ago", rarity: "Common" },
       { title: "Memory Master", description: "Master pointer manipulation", icon: Brain, earned: true, date: "1 week ago", rarity: "Rare" },
       { title: "Speed Coder", description: "Complete 10 challenges in one session", icon: Flame, earned: false, progress: 7, rarity: "Epic" },
       { title: "System Builder", description: "Build a complete C project", icon: Award, earned: false, progress: 1, rarity: "Legendary" }
   ];

   const upcomingLessons = [
       { title: "Pointers and Memory Management", duration: "35 min", type: "Interactive", difficulty: "Intermediate", icon: "🧠" },
       { title: "Dynamic Memory Allocation", duration: "28 min", type: "Hands-on", difficulty: "Intermediate", icon: "⚡" },
       { title: "File I/O Operations", duration: "22 min", type: "Project", difficulty: "Beginner", icon: "📁" },
       { title: "Data Structures in C", duration: "40 min", type: "Deep Dive", difficulty: "Advanced", icon: "🏗️" }
   ];

   const stats = [
       { label: "Total XP", value: user.totalXP.toLocaleString(), icon: Star, color: "from-yellow-400 to-orange-500", bg: "bg-gradient-to-br from-yellow-400/10 to-orange-500/10" },
       { label: "Current Streak", value: `${user.streak} days`, icon: Flame, color: "from-orange-400 to-red-500", bg: "bg-gradient-to-br from-orange-400/10 to-red-500/10" },
       { label: "Course Progress", value: "42%", icon: Trophy, color: "from-green-400 to-emerald-500", bg: "bg-gradient-to-br from-green-400/10 to-emerald-500/10" },
       { label: "Study Time", value: "24h", icon: Clock, color: "from-blue-400 to-cyan-500", bg: "bg-gradient-to-br from-blue-400/10 to-cyan-500/10" }
   ];

   const quickActions = [
       { title: "Continue Course", description: "Pick up where you left off", icon: BookOpenCheck, color: "bg-gradient-to-r from-blue-500 to-purple-600" },
       { title: "Practice Problems", description: "Sharpen your skills", icon: Target, color: "bg-gradient-to-r from-green-500 to-teal-600" },
       { title: "Join Study Group", description: "Learn with others", icon: Users, color: "bg-gradient-to-r from-purple-500 to-pink-600" }
   ];

   const changeMeteors = (num: number) => {
       meteorStyles = [];
       const styles = [...new Array(num)].map(() => ({
           top: Math.floor(Math.random() * 100) + "%",
           left: Math.floor(Math.random() * 100) + "%",
           animationDelay: Math.random() * 2 + "s",
           animationDuration: Math.floor(Math.random() * 3 + 2) + "s",
       }));
       meteorStyles = styles;
   };

   onMount(() => {
       changeMeteors(25);
   });
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 text-foreground flex relative overflow-hidden">
   <div class="absolute inset-0 bg-grid-white/[0.02] bg-[size:50px_50px]" />
   <div class="absolute inset-0 bg-gradient-to-t from-slate-950 via-transparent to-transparent" />
   
   <aside class="relative z-10 w-72 backdrop-blur-xl bg-slate-900/50 border-r border-slate-800/50 shadow-2xl">
       <div class="p-6 border-b border-slate-800/50">
           <div class="flex items-center space-x-3">
               <div class="relative">
                   <div class="p-3 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl shadow-lg">
                       <Code2 class="h-6 w-6 text-white" />
                   </div>
                   <div class="absolute -top-1 -right-1 w-3 h-3 bg-green-400 rounded-full border-2 border-slate-900 animate-pulse" />
               </div>
               <div>
                   <h1 class="text-xl font-bold bg-gradient-to-r from-white to-slate-300 bg-clip-text text-transparent">Codepass</h1>
                   <p class="text-xs text-slate-400">Master. Build. Excel.</p>
               </div>
           </div>
       </div>

       <div class="p-6 border-b border-slate-800/50">
           <div class="flex items-center space-x-4">
               <div class="relative">
                   <Avatar class="h-12 w-12 ring-2 ring-blue-500/20">
                       <AvatarImage src={user.avatar} alt={user.name} />
                       <AvatarFallback class="bg-gradient-to-br from-blue-500 to-purple-600 text-white font-semibold">
                           {user.name.split(' ').map(n => n[0]).join('')}
                       </AvatarFallback>
                   </Avatar>
                   <div class="absolute -bottom-1 -right-1 bg-gradient-to-r from-yellow-400 to-orange-500 text-black text-xs font-bold px-1.5 py-0.5 rounded-full">
                       {user.level}
                   </div>
               </div>
               <div>
                   <p class="font-semibold text-white">{user.name}</p>
                   <div class="flex items-center space-x-2">
                       <div class="flex items-center space-x-1 text-orange-400">
                           <Flame class="h-3 w-3" />
                           <span class="text-xs font-medium">{user.streak}</span>
                       </div>
                       <div class="flex items-center space-x-1 text-yellow-400">
                           <Star class="h-3 w-3" />
                           <span class="text-xs font-medium">{user.totalXP.toLocaleString()}</span>
                       </div>
                   </div>
               </div>
           </div>
       </div>

       <nav class="p-4 space-y-2">
           {#each [
               { id: 'home', label: 'Dashboard', icon: Home },
               { id: 'runners', label: 'Code Runners', icon: Play },
               { id: 'tracks', label: 'Learning Tracks', icon: BookOpen },
               { id: 'changelog', label: 'Updates', icon: FileText }
           ] as item}
               <button
                   class={cn(
                       "w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm transition-all duration-200 group",
                       currentTab === item.id 
                           ? "bg-gradient-to-r from-blue-500/20 to-purple-500/20 text-white border border-blue-500/30 shadow-lg shadow-blue-500/10" 
                           : "text-slate-400 hover:text-white hover:bg-slate-800/50 border border-transparent"
                   )}
                   on:click={() => currentTab = item.id}
               >
                   <svelte:component this={item.icon} class={cn("h-5 w-5 transition-transform group-hover:scale-110", currentTab === item.id && "text-blue-400")} />
                   <span class="font-medium">{item.label}</span>
                   {#if currentTab === item.id}
                       <div class="ml-auto w-2 h-2 bg-blue-400 rounded-full animate-pulse" />
                   {/if}
               </button>
           {/each}
       </nav>

       <div class="p-4 mt-auto">
           <div class="bg-gradient-to-br from-slate-800/50 to-slate-900/50 p-4 rounded-xl border border-slate-700/50 backdrop-blur-sm">
               <div class="space-y-3">
                   <div class="flex items-center justify-between text-sm">
                       <span class="text-slate-300 font-medium">Weekly Goal</span>
                       <div class="flex items-center space-x-1">
                           <Coffee class="h-3 w-3 text-amber-400" />
                           <span class="font-bold text-white">4/7 days</span>
                       </div>
                   </div>
                   <div class="relative">
                       <Progress value={57} class="h-3 bg-slate-800" />
                       <div class="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full opacity-20" />
                   </div>
                   <p class="text-xs text-slate-400">🔥 3 more days to maintain your streak!</p>
               </div>
           </div>
       </div>
   </aside>

   <main class="flex-1 overflow-y-auto relative z-10">
       {#if currentTab === 'home'}
           <div class="p-8 space-y-8">
               <div class="flex items-center justify-between">
                   <div>
                       <h1 class="text-4xl font-bold bg-gradient-to-r from-white via-blue-100 to-purple-200 bg-clip-text text-transparent">
                           Welcome back, {user.name.split(' ')[0]}! 
                           <span class="wave inline-block">👋</span>
                       </h1>
                       <p class="text-slate-400 mt-2 text-lg">Ready to master the art of C programming?</p>
                   </div>
                   <div class="flex items-center space-x-4">
                       <div class="flex items-center space-x-2 bg-gradient-to-r from-orange-500/20 to-red-500/20 text-orange-300 px-4 py-2 rounded-full border border-orange-500/30 backdrop-blur-sm">
                           <Flame class="h-4 w-4" />
                           <span class="font-semibold">{user.streak} day streak</span>
                           <Sparkles class="h-3 w-3 text-yellow-400" />
                       </div>
                       <Badge variant="secondary" class="bg-gradient-to-r from-blue-500/20 to-purple-500/20 text-blue-300 border-blue-500/30 px-3 py-1">
                           Level {user.level}
                       </Badge>
                   </div>
               </div>

               <div class="grid grid-cols-4 gap-6">
                   {#each stats as stat, index}
                       <Card class="relative overflow-hidden bg-slate-900/50 border-slate-700/50 backdrop-blur-sm hover:scale-105 transition-transform duration-200 group">
                           <div class="absolute inset-0 bg-gradient-to-br from-slate-800/20 to-transparent" />
                           <CardContent class="p-6 relative">
                               <div class="flex items-center justify-between">
                                   <div class={cn("p-3 rounded-xl shadow-lg", stat.bg)}>
                                       <svelte:component this={stat.icon} class={cn("h-6 w-6 bg-gradient-to-r bg-clip-text text-transparent", stat.color)} />
                                   </div>
                                   <TrendingUp class="h-4 w-4 text-green-400 opacity-60 group-hover:opacity-100 transition-opacity" />
                               </div>
                               <div class="mt-4">
                                   <p class="text-3xl font-bold text-white">{stat.value}</p>
                                   <p class="text-sm text-slate-400 mt-1">{stat.label}</p>
                               </div>
                           </CardContent>
                       </Card>
                   {/each}
               </div>

               <div class="grid grid-cols-3 gap-8">
                   <div class="col-span-2 space-y-6">
                       <Card class="relative overflow-hidden bg-gradient-to-br from-slate-900/80 via-slate-800/50 to-slate-900/80 border-slate-700/50 backdrop-blur-xl">
                           <div class="absolute inset-0 overflow-hidden">
                               {#each meteorStyles as style, idx}
                                   <span
                                       class="pointer-events-none absolute size-1 rotate-45 animate-meteor rounded-full bg-blue-400/60 shadow-[0_0_0_1px_rgba(59,130,246,0.3)]"
                                       style="top: {style.top}; left: {style.left}; animation-delay: {style.animationDelay}; animation-duration: {style.animationDuration};"
                                   >
                                       <div class="pointer-events-none absolute top-1/2 -z-10 h-px w-12 -translate-y-1/2 bg-gradient-to-r from-blue-400/60 via-purple-400/30 to-transparent" />
                                   </span>
                               {/each}
                           </div>
                           
                           <div 
                               class="pointer-events-none absolute inset-[0] rounded-[inherit] [border:2px_solid_transparent] ![mask-clip:padding-box,border-box] ![mask-composite:intersect] [mask:linear-gradient(transparent,transparent),linear-gradient(white,white)] after:absolute after:aspect-square after:w-[300px] after:animate-border-beam after:[animation-delay:0s] after:[background:linear-gradient(to_left,#3b82f6,#8b5cf6,transparent)] after:[offset-anchor:90%_50%] after:[offset-path:rect(0_auto_auto_0_round_12px)]"
                           />

                           <CardHeader class="relative z-10 pb-4">
                               <div class="flex items-center justify-between">
                                   <div>
                                       <CardTitle class="text-2xl flex items-center gap-3">
                                           <div class="p-2 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg">
                                               <BookOpen class="h-6 w-6 text-white" />
                                           </div>
                                           <span class="bg-gradient-to-r from-white to-slate-300 bg-clip-text text-transparent">
                                               Continue Your Journey
                                           </span>
                                       </CardTitle>
                                       <CardDescription class="text-slate-400 mt-2">Master the fundamentals and build amazing projects</CardDescription>
                                   </div>
                                   <div class="flex flex-col items-end space-y-2">
                                       <Badge variant="outline" class="bg-green-500/10 text-green-400 border-green-500/30">
                                           {currentCourse.difficulty}
                                       </Badge>
                                       <div class="flex items-center space-x-1 text-slate-400 text-sm">
                                           <Users class="h-3 w-3" />
                                           <span>{currentCourse.students.toLocaleString()}</span>
                                       </div>
                                   </div>
                               </div>
                           </CardHeader>
                           <CardContent class="relative z-10 space-y-6">
                               <div>
                                   <div class="flex items-center justify-between mb-3">
                                       <h3 class="text-2xl font-bold text-white">{currentCourse.title}</h3>
                                       <div class="flex items-center space-x-1">
                                           <Star class="h-4 w-4 text-yellow-400 fill-current" />
                                           <span class="text-sm font-medium text-white">{currentCourse.rating}</span>
                                       </div>
                                   </div>
                                   <p class="text-slate-300 mb-6 leading-relaxed">{currentCourse.description}</p>
                                   
                                   <div class="grid grid-cols-2 gap-6 mb-6">
                                       <div class="space-y-3">
                                           <div class="flex justify-between text-sm">
                                               <span class="text-slate-400">Overall Progress</span>
                                               <span class="font-bold text-blue-400">{currentCourse.progress}%</span>
                                           </div>
                                           <div class="relative">
                                               <Progress value={currentCourse.progress} class="h-4 bg-slate-800" />
                                               <div class="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full" 
                                                    style="width: {currentCourse.progress}%; opacity: 0.8;" />
                                           </div>
                                       </div>
                                       <div class="space-y-2">
                                           <p class="text-sm text-slate-400">Lessons Completed</p>
                                           <p class="text-2xl font-bold text-white">{currentCourse.completedLessons}<span class="text-slate-400">/{currentCourse.totalLessons}</span></p>
                                           <p class="text-xs text-slate-500">{currentCourse.estimatedTime} remaining</p>
                                       </div>
                                   </div>

                                   <div class="bg-gradient-to-r from-slate-800/50 to-slate-700/50 p-6 rounded-xl border border-slate-600/30 backdrop-blur-sm">
                                       <div class="flex items-center justify-between">
                                           <div>
                                               <p class="font-semibold text-white text-lg">Next: {currentCourse.nextLesson}</p>
                                               <div class="flex items-center space-x-4 mt-2 text-slate-400">
                                                   <span class="text-sm">By {currentCourse.instructor}</span>
                                                   <span class="text-sm">• {upcomingLessons[0].duration}</span>
                                                   <Badge variant="outline" size="sm" class="bg-blue-500/10 text-blue-400 border-blue-500/30">
                                                       {upcomingLessons[0].type}
                                                   </Badge>
                                               </div>
                                           </div>
                                           <Button class="bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white px-6 py-3 shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40 transition-all duration-200">
                                               Continue Learning
                                               <ArrowRight class="ml-2 h-5 w-5" />
                                           </Button>
                                       </div>
                                   </div>
                               </div>
                           </CardContent>
                       </Card>

                       <div class="grid grid-cols-3 gap-4">
                           {#each quickActions as action}
                               <Card class="relative overflow-hidden border-slate-700/50 backdrop-blur-sm hover:scale-105 transition-all duration-200 cursor-pointer group">
                                   <div class="absolute inset-0 bg-gradient-to-br from-slate-800/30 to-slate-900/50" />
                                   <CardContent class="p-6 relative text-center">
                                       <div class={cn("w-12 h-12 rounded-xl mx-auto mb-4 flex items-center justify-center shadow-lg", action.color)}>
                                           <svelte:component this={action.icon} class="h-6 w-6 text-white" />
                                       </div>
                                       <h4 class="font-semibold text-white mb-2">{action.title}</h4>
                                       <p class="text-sm text-slate-400">{action.description}</p>
                                       <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                                   </CardContent>
                               </Card>
                           {/each}
                       </div>

                       <Card class="bg-slate-900/50 border-slate-700/50 backdrop-blur-sm">
                           <CardHeader>
                               <CardTitle class="text-xl flex items-center gap-3">
                                   <div class="p-2 bg-gradient-to-br from-purple-500 to-pink-600 rounded-lg">
                                       <Calendar class="h-5 w-5 text-white" />
                                   </div>
                                   <span class="text-white">Upcoming Lessons</span>
                               </CardTitle>
                               <CardDescription class="text-slate-400">Your personalized learning roadmap</CardDescription>
                           </CardHeader>
                           <CardContent>
                               <div class="space-y-4">
                                   {#each upcomingLessons as lesson, index}
                                       <div class="group flex items-center justify-between p-4 rounded-xl border border-slate-700/50 hover:bg-slate-800/30 transition-all duration-200 cursor-pointer hover:border-blue-500/30">
                                           <div class="flex items-center space-x-4">
                                               <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-slate-700 to-slate-800 flex items-center justify-center text-lg font-semibold text-white group-hover:from-blue-500 group-hover:to-purple-600 transition-all duration-200">
                                                   {lesson.icon}
                                               </div>
                                               <div>
                                                   <p class="font-semibold text-white group-hover:text-blue-300 transition-colors">{lesson.title}</p>
                                                   <div class="flex items-center space-x-3 text-sm text-slate-400 mt-1">
                                                       <span>{lesson.duration}</span>
                                                       <span>•</span>
                                                       <span>{lesson.type}</span>
                                                       <Badge variant="outline" size="sm" class={cn(
                                                           "ml-2",
                                                           lesson.difficulty === 'Beginner' && "bg-green-500/10 text-green-400 border-green-500/30",
                                                           lesson.difficulty === 'Intermediate' && "bg-yellow-500/10 text-yellow-400 border-yellow-500/30",
                                                           lesson.difficulty === 'Advanced' && "bg-red-500/10 text-red-400 border-red-500/30"
                                                       )}>
                                                           {lesson.difficulty}
                                                       </Badge>
                                                   </div>
                                               </div>
                                           </div>
                                           <ChevronRight class="h-5 w-5 text-slate-500 group-hover:text-blue-400 group-hover:translate-x-1 transition-all duration-200" />
                                       </div>
                                   {/each}
                               </div>
                           </CardContent>
                       </Card>
                   </div>

                   <div class="space-y-6">
                       <Card class="bg-slate-900/50 border-slate-700/50 backdrop-blur-sm">
                           <CardHeader>
                               <CardTitle class="text-lg flex items-center gap-3">
                                   <div class="p-2 bg-gradient-to-br from-yellow-500 to-orange-600 rounded-lg">
                                       <Trophy class="h-5 w-5 text-white" />
                                   </div>
                                   <span class="text-white">Achievements</span>
                               </CardTitle>
                               <CardDescription class="text-slate-400">Your learning milestones</CardDescription>
                           </CardHeader>
                           <CardContent class="space-y-4">
                               {#each achievements as achievement}
                                   <div class={cn(
                                       "relative overflow-hidden flex items-start space-x-4 p-4 rounded-xl transition-all duration-200 border",
                                       achievement.earned 
                                           ? "bg-gradient-to-r from-green-500/10 to-emerald-500/10 border-green-500/30 shadow-lg shadow-green-500/10" 
                                           : "bg-slate-800/30 border-slate-700/50 hover:border-slate-600/50"
                                   )}>
                                       {#if achievement.earned}
                                           <div class="absolute top-2 right-2">
                                               <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
                                           </div>
                                       {/if}
                                       <div class={cn(
                                           "p-3 rounded-xl",
                                           achievement.earned 
                                               ? "bg-gradient-to-br from-green-500 to-emerald-600 text-white shadow-lg" 
                                               : "bg-slate-700 text-slate-400"
                                       )}>
                                           <svelte:component this={achievement.icon} class="h-5 w-5" />
                                       </div>
                                       <div class="flex-1 min-w-0">
                                           <div class="flex items-center justify-between mb-1">
                                               <p class="font-semibold text-white text-sm">{achievement.title}</p>
                                               <Badge 
                                                   variant="outline" 
                                                   size="sm"
                                                   class={cn(
                                                       "text-xs",
                                                       achievement.rarity === 'Common' && "bg-gray-500/10 text-gray-400 border-gray-500/30",
                                                       achievement.rarity === 'Rare' && "bg-blue-500/10 text-blue-400 border-blue-500/30",
                                                       achievement.rarity === 'Epic' && "bg-purple-500/10 text-purple-400 border-purple-500/30",
                                                       achievement.rarity === 'Legendary' && "bg-yellow-500/10 text-yellow-400 border-yellow-500/30"
                                                   )}
                                               >
                                                   {achievement.rarity}
                                               </Badge>
                                           </div>
                                           <p class="text-xs text-slate-400 mb-3">{achievement.description}</p>
                                           {#if achievement.earned}
                                               <p class="text-xs text-green-400 font-medium">✨ Earned {achievement.date}</p>
                                           {:else if achievement.progress}
                                               <div class="space-y-2">
                                                   <div class="flex justify-between text-xs">
                                                       <span class="text-slate-400">Progress</span>
                                                       <span class="font-medium text-blue-400">{achievement.progress}/{achievement.rarity === 'Epic' ? '10' : '3'}</span>
                                                   </div>
                                                   <div class="relative">
                                                       <Progress value={(achievement.progress / (achievement.rarity === 'Epic' ? 10 : 3)) * 100} class="h-2 bg-slate-800" />
                                                       <div class="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full opacity-60" 
                                                            style="width: {(achievement.progress / (achievement.rarity === 'Epic' ? 10 : 3)) * 100}%;" />
                                                   </div>
                                               </div>
                                           {/if}
                                       </div>
                                   </div>
                               {/each}
                           </CardContent>
                       </Card>

                       <Card class="bg-gradient-to-br from-slate-900/80 to-slate-800/50 border-slate-700/50 backdrop-blur-sm">
                           <CardHeader>
                               <CardTitle class="text-lg flex items-center gap-3">
                                   <div class="p-2 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg">
                                       <TrendingUp class="h-5 w-5 text-white" />
                                   </div>
                                   <span class="text-white">Learning Insights</span>
                               </CardTitle>
                           </CardHeader>
                           <CardContent class="space-y-4">
                               <div class="space-y-3">
                                   <div class="flex items-center justify-between p-3 bg-slate-800/50 rounded-lg border border-slate-700/30">
                                       <div class="flex items-center space-x-2">
                                           <div class="w-2 h-2 bg-green-400 rounded-full"></div>
                                           <span class="text-sm text-slate-300">Most Active Time</span>
                                       </div>
                                       <span class="text-sm font-semibold text-white">9:00 PM</span>
                                   </div>
                                   
                                   <div class="flex items-center justify-between p-3 bg-slate-800/50 rounded-lg border border-slate-700/30">
                                       <div class="flex items-center space-x-2">
                                           <div class="w-2 h-2 bg-blue-400 rounded-full"></div>
                                           <span class="text-sm text-slate-300">Average Session</span>
                                       </div>
                                       <span class="text-sm font-semibold text-white">45 min</span>
                                   </div>
                                   
                                   <div class="flex items-center justify-between p-3 bg-slate-800/50 rounded-lg border border-slate-700/30">
                                       <div class="flex items-center space-x-2">
                                           <div class="w-2 h-2 bg-purple-400 rounded-full"></div>
                                           <span class="text-sm text-slate-300">Favorite Topic</span>
                                       </div>
                                       <span class="text-sm font-semibold text-white">Pointers</span>
                                   </div>
                               </div>
                               
                               <div class="mt-4 p-4 bg-gradient-to-r from-blue-500/10 to-purple-500/10 rounded-xl border border-blue-500/20">
                                   <div class="flex items-center space-x-2 mb-2">
                                       <Brain class="h-4 w-4 text-blue-400" />
                                       <span class="text-sm font-semibold text-blue-300">AI Recommendation</span>
                                   </div>
                                   <p class="text-xs text-slate-300 leading-relaxed">
                                       Based on your progress, consider practicing more with dynamic memory allocation before moving to advanced data structures.
                                   </p>
                               </div>
                           </CardContent>
                       </Card>
                   </div>
               </div>
           </div>
       {:else if currentTab === 'runners'}
           <div class="p-8 space-y-6">
               <div class="text-center py-16">
                   <div class="w-16 h-16 bg-gradient-to-br from-green-500 to-teal-600 rounded-full mx-auto mb-4 flex items-center justify-center">
                       <Play class="h-8 w-8 text-white" />
                   </div>
                   <h2 class="text-2xl font-bold text-white mb-2">Code Runners</h2>
                   <p class="text-slate-400 mb-6">Interactive coding challenges and live execution environment</p>
                   <Button class="bg-gradient-to-r from-green-500 to-teal-600 hover:from-green-600 hover:to-teal-700 text-white">
                       Launch Code Runner
                   </Button>
               </div>
           </div>
       {:else if currentTab === 'tracks'}
           <div class="p-8 space-y-6">
               <div class="text-center py-16">
                   <div class="w-16 h-16 bg-gradient-to-br from-purple-500 to-pink-600 rounded-full mx-auto mb-4 flex items-center justify-center">
                       <BookOpen class="h-8 w-8 text-white" />
                   </div>
                   <h2 class="text-2xl font-bold text-white mb-2">Learning Tracks</h2>
                   <p class="text-slate-400 mb-6">Structured learning paths for different skill levels</p>
                   <Button class="bg-gradient-to-r from-purple-500 to-pink-600 hover:from-purple-600 hover:to-pink-700 text-white">
                       Explore Tracks
                   </Button>
               </div>
           </div>
       {:else if currentTab === 'changelog'}
           <div class="p-8 space-y-6">
               <div class="text-center py-16">
                   <div class="w-16 h-16 bg-gradient-to-br from-orange-500 to-red-600 rounded-full mx-auto mb-4 flex items-center justify-center">
                       <FileText class="h-8 w-8 text-white" />
                   </div>
                   <h2 class="text-2xl font-bold text-white mb-2">Platform Updates</h2>
                   <p class="text-slate-400 mb-6">Latest features, improvements, and announcements</p>
                   <Button class="bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white">
                       View Changelog
                   </Button>
               </div>
           </div>
       {/if}
   </main>
</div>

<style>
   .wave {
       animation: wave 2s ease-in-out infinite;
   }
   
   @keyframes wave {
       0%, 100% { transform: rotate(0deg); }
       25% { transform: rotate(20deg); }
       75% { transform: rotate(-10deg); }
   }
   
   @keyframes meteor {
       0% {
           transform: rotate(45deg) translateX(0);
           opacity: 1;
       }
       70% {
           opacity: 1;
       }
       100% {
           transform: rotate(45deg) translateX(-500px);
           opacity: 0;
       }
   }
   
   .animate-meteor {
       animation: meteor linear infinite;
   }
   
   @keyframes border-beam {
       100% {
           offset-distance: 100%;
       }
   }
   
   .animate-border-beam {
       animation: border-beam calc(var(--duration)*1s) infinite linear;
   }
</style>
                               