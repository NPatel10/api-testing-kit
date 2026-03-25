<script lang="ts">
	import { cn } from "$lib/utils.js";

	type MarketingLink = {
		href: string;
		label: string;
	};

	let {
		brand = "API Testing Kit",
		tagline = "A safer way to demo and ship HTTP workflows",
		links = [],
		ctaHref = "/app",
		ctaLabel = "Open app",
		class: className
	} = $props<{
		brand?: string;
		tagline?: string;
		links?: MarketingLink[];
		ctaHref?: string;
		ctaLabel?: string;
		class?: string;
	}>();
</script>

<header
	class={cn(
		"border-b border-border/80 bg-white/80 px-5 py-4 backdrop-blur sm:px-6 lg:px-8",
		className
	)}
>
	<div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
		<div class="flex items-center gap-3">
			<div class="grid h-11 w-11 place-items-center rounded-2xl bg-primary-green text-sm font-semibold text-white shadow-[0_10px_24px_rgba(31,122,77,0.28)]">
				AT
			</div>
			<div>
				<p class="text-sm font-semibold tracking-tight text-text-strong">{brand}</p>
				<p class="text-xs text-text-muted">{tagline}</p>
			</div>
		</div>

		<nav class="flex flex-wrap items-center gap-2 text-sm">
			{#each links as link}
				<a
					href={link.href}
					class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong"
				>
					{link.label}
				</a>
			{/each}
			<a
				href={ctaHref}
				class="rounded-full bg-primary-green px-5 py-2.5 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-primary-green-hover"
			>
				{ctaLabel}
			</a>
		</nav>
	</div>
</header>
