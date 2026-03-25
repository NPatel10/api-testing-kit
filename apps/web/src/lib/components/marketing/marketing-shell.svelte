<script lang="ts">
	import { cn } from "$lib/utils.js";

	let { class: className, children } = $props<{
		class?: string;
		children?: () => unknown;
	}>();
</script>

<div
	class={cn(
		"relative isolate overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(31,122,77,0.16),transparent_32%),radial-gradient(circle_at_top_right,rgba(111,142,163,0.12),transparent_28%),linear-gradient(180deg,#f2f0ea_0%,#ece8df_100%)] text-text-strong",
		className
	)}
>
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-32 top-24 h-80 w-80 rounded-full bg-primary-green/10 blur-3xl"></div>
		<div class="absolute right-[-4rem] top-1/3 h-96 w-96 rounded-full bg-primary-green-soft blur-3xl"></div>
		<div class="absolute bottom-[-8rem] left-1/2 h-72 w-72 -translate-x-1/2 rounded-full bg-primary-green-deep/10 blur-3xl"></div>
	</div>

	<div class="mx-auto min-h-screen max-w-[1440px] px-4 py-4 sm:px-6 lg:px-8">
		{@render children?.()}
	</div>
</div>
