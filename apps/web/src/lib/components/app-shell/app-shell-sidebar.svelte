<script lang="ts">
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		SidebarContent,
		SidebarFooter,
		SidebarGroup,
		SidebarGroupContent,
		SidebarGroupLabel,
		SidebarHeader,
		SidebarMenu,
		SidebarMenuItem,
		SidebarMenuButton,
	} from "$lib/components/ui/sidebar/index.js";
	import Separator from "$lib/components/ui/separator/separator.svelte";
	import LayoutDashboardIcon from "@lucide/svelte/icons/layout-dashboard";
	import FolderKanbanIcon from "@lucide/svelte/icons/folder-kanban";
	import HistoryIcon from "@lucide/svelte/icons/history";
	import BookOpenTextIcon from "@lucide/svelte/icons/book-open-text";
	import Settings2Icon from "@lucide/svelte/icons/settings-2";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";
	import TerminalSquareIcon from "@lucide/svelte/icons/terminal-square";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";

	const primaryNav = [
		{ label: "Workspace", href: "/app", icon: LayoutDashboardIcon, active: true, badge: "Live" },
		{ label: "Templates", href: "/templates", icon: BookOpenTextIcon, badge: "12" },
		{ label: "Collections", href: "/app/collections/demo", icon: FolderKanbanIcon, badge: "3" },
		{ label: "History", href: "/app/history", icon: HistoryIcon },
		{ label: "Settings", href: "/app/settings", icon: Settings2Icon },
	];

	const utilityLinks = [
		{ label: "Safety rules", href: "/case-study", icon: ShieldCheckIcon },
		{ label: "Quick snippets", href: "/docs", icon: TerminalSquareIcon },
	];
</script>

<SidebarHeader class="gap-4 px-4 py-5">
	<div class="flex items-center gap-3">
		<div class="bg-primary/10 text-primary flex size-10 items-center justify-center rounded-2xl border border-emerald-800/10">
			<SparklesIcon class="size-5" />
		</div>
		<div class="min-w-0">
			<p class="text-sm font-semibold tracking-tight text-foreground">API Testing Kit</p>
			<p class="text-xs text-muted-foreground">Guest and authenticated workspace</p>
		</div>
	</div>

	<div class="bg-panel-soft/90 rounded-2xl border border-border/70 p-3">
		<div class="flex items-center justify-between gap-3">
			<div>
				<p class="text-xs uppercase tracking-[0.24em] text-muted-foreground">Session</p>
				<p class="text-sm font-medium text-foreground">Guest preview</p>
			</div>
			<Badge variant="secondary">Locked</Badge>
		</div>
		<p class="mt-2 text-xs leading-5 text-muted-foreground">
			Guests can explore the live shell. Custom URLs and saved data stay behind sign-in.
		</p>
	</div>
</SidebarHeader>

<SidebarContent class="px-3 pb-3">
	<SidebarGroup>
		<SidebarGroupLabel>Primary</SidebarGroupLabel>
		<SidebarGroupContent>
			<SidebarMenu class="gap-1">
				{#each primaryNav as item}
					{@const Icon = item.icon}
					<SidebarMenuItem>
						<SidebarMenuButton isActive={item.active} class="justify-between">
							{#snippet child({ props })}
								<a href={item.href} {...props}>
									<span class="flex items-center gap-2">
										<Icon class="size-4" />
										<span>{item.label}</span>
									</span>
									{#if item.badge}
										<Badge variant={item.active ? "default" : "outline"} class="ml-auto">
											{item.badge}
										</Badge>
									{/if}
								</a>
							{/snippet}
						</SidebarMenuButton>
					</SidebarMenuItem>
				{/each}
			</SidebarMenu>
		</SidebarGroupContent>
	</SidebarGroup>

	<SidebarGroup>
		<SidebarGroupLabel>Utilities</SidebarGroupLabel>
		<SidebarGroupContent>
			<SidebarMenu class="gap-1">
				{#each utilityLinks as item}
					{@const Icon = item.icon}
					<SidebarMenuItem>
						<SidebarMenuButton class="justify-start gap-2">
							{#snippet child({ props })}
								<a href={item.href} {...props}>
									<Icon class="size-4" />
									<span>{item.label}</span>
								</a>
							{/snippet}
						</SidebarMenuButton>
					</SidebarMenuItem>
				{/each}
			</SidebarMenu>
		</SidebarGroupContent>
	</SidebarGroup>

	<div class="px-2">
		<div class="rounded-2xl border border-border/70 bg-gradient-to-br from-emerald-50 to-white p-4 shadow-sm">
			<p class="text-xs uppercase tracking-[0.24em] text-muted-foreground">Guest limits</p>
			<p class="mt-2 text-sm font-medium text-foreground">Allowlisted endpoints only</p>
			<p class="mt-1 text-xs leading-5 text-muted-foreground">
				The shell is real, but request execution stays constrained until authentication lands.
			</p>
			<Button variant="outline" size="sm" class="mt-3 w-full">Sign in to unlock</Button>
		</div>
	</div>
</SidebarContent>

<SidebarFooter class="gap-3 px-4 pb-4">
	<Separator />
	<div class="flex items-center justify-between text-xs text-muted-foreground">
		<span>v0 shell</span>
		<span>Safe by default</span>
	</div>
</SidebarFooter>
