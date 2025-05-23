<script lang="ts">
	import '../app.css';
	import Nav from '$lib/components/Nav.svelte';
	import { ModeWatcher } from "mode-watcher";
	import { Toaster } from "$lib/components/ui/sonner/index.js";
</script>

<div class="min-h-screen flex flex-col">
	<ModeWatcher />
	<Nav />
	<main class="flex-1">
		<slot />
	</main>
	<Toaster />
</div>