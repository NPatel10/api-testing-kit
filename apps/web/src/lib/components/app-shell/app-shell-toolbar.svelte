<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { SidebarTrigger } from "$lib/components/ui/sidebar/index.js";
	import SearchIcon from "@lucide/svelte/icons/search";
	import SlidersHorizontalIcon from "@lucide/svelte/icons/sliders-horizontal";
	import CirclePlusIcon from "@lucide/svelte/icons/circle-plus";
	import BellIcon from "@lucide/svelte/icons/bell";
	import { getEntitlementSummary, type EffectiveEntitlements } from "$lib/entitlements/access";
	import type { WorkspaceMode } from "$lib/mocks/workspace-state";

	let {
		mode = "guest",
		sessionLabel = "Guest preview",
		entitlements,
	}: {
		mode?: WorkspaceMode;
		sessionLabel?: string;
		entitlements?: EffectiveEntitlements;
	} = $props();
</script>

<div class="flex flex-wrap items-center gap-3 rounded-[24px] border border-border/80 bg-white px-4 py-4 shadow-sm">
	<div class="flex min-w-0 flex-1 items-center gap-3">
		<SidebarTrigger class="shrink-0" />
		<div class="flex min-w-0 flex-1 items-center gap-3 rounded-full border border-border/80 bg-shell px-4 py-2">
			<SearchIcon class="size-4 text-muted-foreground" />
			<span class="min-w-0 truncate text-sm text-muted-foreground">Search requests, templates, and history</span>
		</div>
	</div>

	<div class="flex items-center gap-2">
		<Badge variant={mode === "authenticated" ? "default" : "secondary"} class="hidden sm:inline-flex">
			{sessionLabel}
		</Badge>
		{#if entitlements}
			<Badge variant="outline" class="hidden md:inline-flex">{entitlements.plan.name}</Badge>
		{:else}
			<Badge variant="outline" class="hidden md:inline-flex">{mode}</Badge>
		{/if}
		<Button variant="outline" size="icon-sm" aria-label="Filters">
			<SlidersHorizontalIcon class="size-4" />
		</Button>
		<Button variant="outline" size="icon-sm" aria-label="Notifications">
			<BellIcon class="size-4" />
		</Button>
		<Button variant="default" size="sm" class="hidden sm:inline-flex">
			<CirclePlusIcon class="size-4" />
			New request
		</Button>
	</div>
</div>

{#if entitlements}
	<p class="px-1 pt-2 text-xs leading-5 text-muted-foreground">{getEntitlementSummary(entitlements, mode)}</p>
{/if}
