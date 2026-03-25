<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="theme-color" content="#efede7" />
	<meta name="color-scheme" content="light" />
</svelte:head>

<div class="min-h-dvh bg-canvas text-foreground antialiased">
	{@render children()}
</div>
