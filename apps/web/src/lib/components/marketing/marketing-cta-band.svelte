<script lang="ts">
	import { ArrowRight } from "@lucide/svelte";
	import { cn } from "$lib/utils.js";

	let {
		eyebrow,
		title,
		description,
		primaryHref,
		primaryLabel,
		secondaryHref,
		secondaryLabel,
		class: className
	} = $props<{
		eyebrow?: string;
		title: string;
		description?: string;
		primaryHref: string;
		primaryLabel: string;
		secondaryHref?: string;
		secondaryLabel?: string;
		class?: string;
	}>();
</script>

<section
	class={cn(
		"rounded-[30px] border border-border/80 bg-[linear-gradient(135deg,rgba(31,122,77,0.14),rgba(255,255,255,0.96))] p-6 shadow-[0_18px_40px_rgba(21,31,23,0.08)] sm:p-8",
		className
	)}
>
	<div class="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
		<div class="max-w-2xl">
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

		<div class="flex flex-wrap gap-3">
			<a
				href={primaryHref}
				class="inline-flex items-center gap-2 rounded-full bg-primary-green px-6 py-3 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-primary-green-hover"
			>
				{primaryLabel}
				<ArrowRight size={16} />
			</a>
			{#if secondaryHref && secondaryLabel}
				<a
					href={secondaryHref}
					class="rounded-full border border-border/80 bg-white px-6 py-3 font-medium text-text-strong transition hover:bg-surface-soft"
				>
					{secondaryLabel}
				</a>
			{/if}
		</div>
	</div>
</section>
