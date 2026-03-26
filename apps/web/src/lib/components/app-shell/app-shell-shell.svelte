<script lang="ts">
	import { Sidebar, SidebarInset, SidebarProvider } from "$lib/components/ui/sidebar/index.js";
	import { getDefaultEntitlements, type EffectiveEntitlements } from "$lib/entitlements/access";
	import AppShellSidebar from "./app-shell-sidebar.svelte";
	import AppShellToolbar from "./app-shell-toolbar.svelte";
	import type { WorkspaceMode } from "$lib/mocks/workspace-state";
	import type { Snippet } from "svelte";

	let {
		children,
		mode = "guest",
		sessionLabel = "Guest preview",
		entitlements = getDefaultEntitlements("guest"),
	}: {
		children: Snippet;
		mode?: WorkspaceMode;
		sessionLabel?: string;
		entitlements?: EffectiveEntitlements;
	} = $props();
	let sidebarOpen = $state(true);
</script>

<SidebarProvider bind:open={sidebarOpen}>
	<div class="min-h-screen bg-[radial-gradient(circle_at_top_left,rgba(31,122,77,0.12),transparent_36%),linear-gradient(180deg,#f7f5f0_0%,#efeeea_100%)] p-3 md:p-4">
		<div class="mx-auto flex min-h-[calc(100vh-1.5rem)] max-w-[1440px] overflow-hidden rounded-[32px] border border-white/70 bg-shell shadow-[0_24px_60px_rgba(21,31,23,0.08)] md:min-h-[calc(100vh-2rem)]">
			<Sidebar class={mode === "authenticated" ? "bg-[#f3f7f1]" : undefined}>
				<AppShellSidebar mode={mode} sessionLabel={sessionLabel} entitlements={entitlements} />
			</Sidebar>

			<SidebarInset class="flex min-w-0 flex-1 flex-col bg-shell">
				<div class="flex-1 px-4 py-4 md:px-5 md:py-5">
					<div class="space-y-4">
						<AppShellToolbar mode={mode} sessionLabel={sessionLabel} entitlements={entitlements} />
						<main class="space-y-4">
							{@render children()}
						</main>
					</div>
				</div>
			</SidebarInset>
		</div>
	</div>
</SidebarProvider>
