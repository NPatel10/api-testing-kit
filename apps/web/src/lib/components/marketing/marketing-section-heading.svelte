<script lang="ts">
	import { cn } from "$lib/utils.js";

	let {
		eyebrow,
		title,
		description,
		align = "left",
		class: className
	} = $props<{
		eyebrow?: string;
		title: string;
		description?: string;
		align?: "left" | "center";
		class?: string;
	}>();
</script>

<div class={cn(align === "center" ? "mx-auto max-w-3xl text-center" : "max-w-3xl", className)}>
	{#if eyebrow}
		<p class="text-xs font-semibold uppercase tracking-[0.3em] text-primary-green">{eyebrow}</p>
	{/if}
	<h2 class="mt-2 text-2xl font-semibold tracking-tight text-text-strong sm:text-3xl">
		{title}
	</h2>
	{#if description}
		<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
			{description}
		</p>
	{/if}
</div>
