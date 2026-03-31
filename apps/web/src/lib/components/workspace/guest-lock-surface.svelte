<script lang="ts">
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import { cn } from "$lib/utils.js";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import LockIcon from "@lucide/svelte/icons/lock";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";

	type GuestLockTone = "emerald" | "sand" | "slate" | "rose";

	type Highlight = {
		label: string;
		value: string;
	};

	interface Props {
		eyebrow?: string;
		title: string;
		description: string;
		badge?: string;
		summary?: string;
		bullets?: string[];
		highlights?: Highlight[];
		primaryActionLabel?: string;
		primaryActionHref?: string;
		secondaryActionLabel?: string;
		secondaryActionHref?: string;
		tone?: GuestLockTone;
	}

	const toneStyles: Record<
		GuestLockTone,
		{
			wash: string;
			accent: string;
			badge: string;
			icon: string;
		}
	> = {
		emerald: {
			wash: "bg-[radial-gradient(circle_at_top_right,rgba(31,122,77,0.12),transparent_45%)]",
			accent: "bg-primary-green-soft text-primary-green-deep",
			badge: "bg-primary-green-soft text-primary-green-deep",
			icon: "text-primary-green-deep",
		},
		sand: {
			wash: "bg-[radial-gradient(circle_at_top_right,rgba(240,180,76,0.16),transparent_45%)]",
			accent: "bg-amber-100 text-amber-900",
			badge: "bg-amber-100 text-amber-900",
			icon: "text-amber-900",
		},
		slate: {
			wash: "bg-[radial-gradient(circle_at_top_right,rgba(111,142,163,0.14),transparent_45%)]",
			accent: "bg-slate-100 text-slate-800",
			badge: "bg-slate-100 text-slate-800",
			icon: "text-slate-800",
		},
		rose: {
			wash: "bg-[radial-gradient(circle_at_top_right,rgba(227,109,93,0.14),transparent_45%)]",
			accent: "bg-rose-100 text-rose-900",
			badge: "bg-rose-100 text-rose-900",
			icon: "text-rose-900",
		},
	};

	let {
		eyebrow = "Guest lock",
		title,
		description,
		badge = "Guest mode",
		summary = "Available after sign-in.",
		bullets = [],
		highlights = [],
		primaryActionLabel = "Sign in to unlock",
		primaryActionHref = "/app",
		secondaryActionLabel = "Docs",
		secondaryActionHref = "/docs",
		tone = "emerald",
	}: Props = $props();
</script>

<section class="panel-card relative overflow-hidden">
	<div class={cn("absolute inset-0 opacity-100", toneStyles[tone].wash)}></div>
	<div class="relative grid gap-5 p-5 md:grid-cols-[minmax(0,1.2fr)_minmax(260px,0.8fr)] md:p-6">
		<div class="space-y-4">
			<div class="flex flex-wrap items-center gap-2">
				<Badge variant="secondary">{eyebrow}</Badge>
				<Badge class={toneStyles[tone].badge}>{badge}</Badge>
			</div>

			<div class="space-y-2">
				<h3 class="text-xl font-semibold tracking-tight text-foreground">{title}</h3>
				<p class="max-w-2xl text-sm leading-6 text-text-body">{description}</p>
			</div>

			<div class="rounded-[20px] border border-border/70 bg-white/85 p-4 shadow-sm">
				<div class="flex items-center gap-2">
					<LockIcon class={cn("size-4", toneStyles[tone].icon)} />
					<p class="text-sm font-medium text-foreground">State</p>
				</div>
				<p class="mt-2 text-sm leading-6 text-text-body">{summary}</p>
			</div>

			{#if bullets.length}
				<ul class="grid gap-2 sm:grid-cols-2">
					{#each bullets as bullet}
						<li class="flex items-start gap-2 rounded-[18px] border border-border/70 bg-white/80 px-3 py-2 text-sm leading-6 text-text-body shadow-sm">
							<span class={cn("mt-2 size-1.5 rounded-full", toneStyles[tone].accent)}></span>
							<span>{bullet}</span>
						</li>
					{/each}
				</ul>
			{/if}

			<div class="flex flex-wrap gap-2">
				<Button href={primaryActionHref} variant="default" size="sm">
					{primaryActionLabel}
					<ArrowRightIcon class="size-4" />
				</Button>
				<Button href={secondaryActionHref} variant="outline" size="sm">
					{secondaryActionLabel}
				</Button>
			</div>
		</div>

		<aside class="soft-card relative overflow-hidden p-4">
			<div class={cn("absolute inset-0", toneStyles[tone].wash)}></div>
			<div class="relative space-y-4">
				<div class="flex items-center justify-between gap-3">
					<div>
						<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Surface state</p>
						<p class="text-sm font-medium text-foreground">Visible, but constrained</p>
					</div>
					<SparklesIcon class={cn("size-4", toneStyles[tone].icon)} />
				</div>

				<div class="rounded-[18px] border border-border/70 bg-white/90 p-4 shadow-sm">
					<div class="flex items-center gap-3">
						<div class={cn("grid size-10 place-items-center rounded-2xl", toneStyles[tone].accent)}>
							<LockIcon class="size-5" />
						</div>
						<div>
							<p class="text-sm font-semibold text-foreground">Access</p>
							<p class="text-xs leading-5 text-text-muted">Locked in guest mode.</p>
						</div>
					</div>
				</div>

				{#if highlights.length}
					<div class="space-y-2">
						{#each highlights as highlight}
							<div class="flex items-center justify-between gap-4 rounded-[16px] border border-border/70 bg-white/85 px-3 py-2 text-sm shadow-sm">
								<span class="text-text-muted">{highlight.label}</span>
								<span class="font-medium text-foreground">{highlight.value}</span>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</aside>
	</div>
</section>
