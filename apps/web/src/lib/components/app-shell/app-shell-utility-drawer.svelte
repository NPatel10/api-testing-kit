<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import {
		Drawer,
		DrawerContent,
		DrawerDescription,
		DrawerHeader,
		DrawerTitle,
	} from "$lib/components/ui/drawer/index.js";
	import { Tabs, TabsContent, TabsList, TabsTrigger } from "$lib/components/ui/tabs/index.js";
	import TerminalIcon from "@lucide/svelte/icons/terminal";
	import CodeXmlIcon from "@lucide/svelte/icons/code-xml";
	import HistoryIcon from "@lucide/svelte/icons/history";
	import { buildEntitlementRows, type EffectiveEntitlements } from "$lib/entitlements/access";
	import type { WorkspaceMode } from "$lib/mocks/workspace-state";

	let {
		mobileLabel = "Utilities",
		mode = "guest",
		entitlements,
	}: {
		mobileLabel?: string;
		mode?: WorkspaceMode;
		entitlements?: EffectiveEntitlements;
	} = $props();
	let open = $state(false);

	const snippets = [
		"curl -X GET https://api.example.com/status",
		"fetch('/api/v1/health').then((res) => res.json())",
		"requests.get('https://api.example.com/status')",
	];
</script>

<div class="hidden lg:block">
	<div class="rounded-[24px] border border-border/80 bg-white shadow-sm">
		<div class="flex items-center justify-between border-b border-border/70 px-5 py-4">
			<div>
				<p class="text-xs uppercase tracking-[0.24em] text-muted-foreground">Utility drawer</p>
				<p class="text-sm font-medium text-foreground">Snippets, history, and locked tools</p>
			</div>
			<Badge variant="outline">Docked</Badge>
		</div>

		<div class="p-4">
			<Tabs value="snippets" class="gap-4">
				<TabsList>
					<TabsTrigger value="snippets">
						<TerminalIcon class="size-4" />
						Snippets
					</TabsTrigger>
					<TabsTrigger value="history">
						<HistoryIcon class="size-4" />
						History
					</TabsTrigger>
					<TabsTrigger value="advanced">
						<CodeXmlIcon class="size-4" />
						Advanced
					</TabsTrigger>
				</TabsList>

				<TabsContent value="snippets" class="space-y-3">
					{#each snippets as snippet}
						<pre class="overflow-x-auto rounded-2xl border border-border/70 bg-slate-950 px-4 py-3 text-xs leading-6 text-slate-100"><code>{snippet}</code></pre>
					{/each}
				</TabsContent>

				<TabsContent value="history" class="space-y-3">
					<div class="rounded-2xl border border-border/70 bg-panel-soft p-4">
						<p class="text-sm font-medium text-foreground">
							{mode === "authenticated" ? "Authenticated history ready" : "No saved history yet"}
						</p>
						<p class="mt-1 text-xs leading-5 text-muted-foreground">
							{mode === "authenticated"
								? "Saved requests and replayable runs will populate here once persistence is wired."
								: "Authenticated requests will populate here once persistence is wired."}
						</p>
					</div>
				</TabsContent>

				<TabsContent value="advanced" class="space-y-3">
					<div class="rounded-2xl border border-dashed border-border/80 bg-gradient-to-br from-emerald-50 to-white p-4">
						<p class="text-sm font-medium text-foreground">
							{mode === "authenticated" ? "Advanced tools available" : "Locked for guests"}
						</p>
						<p class="mt-1 text-xs leading-5 text-muted-foreground">
							{mode === "authenticated"
								? "Environment variables, custom targets, and advanced tooling are part of the authenticated contract."
								: "Environment variables, custom targets, and advanced tooling will unlock after sign-in."}
						</p>
					</div>
					{#if entitlements}
						<div class="grid gap-2">
							{#each buildEntitlementRows(entitlements) as row}
								<div class="rounded-2xl border border-border/70 bg-panel-soft p-3">
									<div class="flex items-center justify-between gap-3">
										<p class="text-sm font-medium text-foreground">{row.label}</p>
										<Badge variant={row.tone === "positive" ? "default" : "secondary"}>{row.statusLabel}</Badge>
									</div>
									<p class="mt-1 text-xs leading-5 text-muted-foreground">{row.description}</p>
									{#if row.limitLabel}
										<p class="mt-2 text-xs font-medium text-foreground">{row.limitLabel}</p>
									{/if}
								</div>
							{/each}
						</div>
					{/if}
				</TabsContent>
			</Tabs>
		</div>
	</div>
</div>

<div class="lg:hidden">
	<Button variant="outline" class="w-full justify-between" onclick={() => (open = true)}>
		<span class="flex items-center gap-2">
			<TerminalIcon class="size-4" />
			{mobileLabel}
		</span>
		<Badge variant="secondary">Drawer</Badge>
	</Button>

	<Drawer bind:open direction="bottom">
		<DrawerContent class="border-border/70 bg-shell">
			<DrawerHeader class="border-b border-border/70">
				<DrawerTitle>Utility drawer</DrawerTitle>
				<DrawerDescription>Quick access to snippets and locked tools.</DrawerDescription>
			</DrawerHeader>
			<div class="p-4">
				<Tabs value="snippets" class="gap-4">
					<TabsList>
						<TabsTrigger value="snippets">Snippets</TabsTrigger>
						<TabsTrigger value="history">History</TabsTrigger>
						<TabsTrigger value="advanced">Advanced</TabsTrigger>
					</TabsList>

					<TabsContent value="snippets" class="space-y-3">
						{#each snippets as snippet}
							<pre class="overflow-x-auto rounded-2xl border border-border/70 bg-slate-950 px-4 py-3 text-xs leading-6 text-slate-100"><code>{snippet}</code></pre>
						{/each}
					</TabsContent>

					<TabsContent value="history" class="space-y-3">
						<div class="rounded-2xl border border-border/70 bg-panel-soft p-4">
							<p class="text-sm font-medium text-foreground">
								{mode === "authenticated" ? "Authenticated history ready" : "No saved history yet"}
							</p>
							<p class="mt-1 text-xs leading-5 text-muted-foreground">
								{mode === "authenticated"
									? "Saved requests and replayable runs will populate here once persistence is wired."
									: "Authenticated requests will populate here once persistence is wired."}
							</p>
						</div>
					</TabsContent>

					<TabsContent value="advanced" class="space-y-3">
						<div class="rounded-2xl border border-dashed border-border/80 bg-gradient-to-br from-emerald-50 to-white p-4">
							<p class="text-sm font-medium text-foreground">
								{mode === "authenticated" ? "Advanced tools available" : "Locked for guests"}
							</p>
							<p class="mt-1 text-xs leading-5 text-muted-foreground">
								{mode === "authenticated"
									? "Environment variables, custom targets, and advanced tooling are part of the authenticated contract."
									: "Environment variables, custom targets, and advanced tooling will unlock after sign-in."}
							</p>
						</div>
						{#if entitlements}
							<div class="grid gap-2">
								{#each buildEntitlementRows(entitlements) as row}
									<div class="rounded-2xl border border-border/70 bg-panel-soft p-3">
										<div class="flex items-center justify-between gap-3">
											<p class="text-sm font-medium text-foreground">{row.label}</p>
											<Badge variant={row.tone === "positive" ? "default" : "secondary"}>{row.statusLabel}</Badge>
										</div>
										<p class="mt-1 text-xs leading-5 text-muted-foreground">{row.description}</p>
										{#if row.limitLabel}
											<p class="mt-2 text-xs font-medium text-foreground">{row.limitLabel}</p>
										{/if}
									</div>
								{/each}
							</div>
						{/if}
					</TabsContent>
				</Tabs>
			</div>
		</DrawerContent>
	</Drawer>
</div>
