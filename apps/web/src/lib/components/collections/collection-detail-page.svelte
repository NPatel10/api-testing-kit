<script lang="ts">
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle,
	} from "$lib/components/ui/card/index.js";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import ExternalLinkIcon from "@lucide/svelte/icons/external-link";
	import LockKeyholeIcon from "@lucide/svelte/icons/lock-keyhole";
	import Layers3Icon from "@lucide/svelte/icons/layers-3";
	import RouteIcon from "@lucide/svelte/icons/route";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";
	import type { EffectiveEntitlements } from "$lib/entitlements/access";
	import { buildEntitlementRows, getEntitlementSummary } from "$lib/entitlements/access";
	import type { WorkspaceMode } from "$lib/mocks/workspace-state";
	import type { CollectionDetailView } from "./collection-detail-data";

	type Props = {
		detail: CollectionDetailView;
		mode: WorkspaceMode;
		entitlements?: EffectiveEntitlements;
		sessionLabel?: string;
	};

	let { detail, mode, entitlements, sessionLabel = "Guest preview" }: Props = $props();

	const badgeTone = {
		ready: "default",
		locked: "secondary",
		missing: "outline",
		unavailable: "outline",
	} as const;

	const requestTone = {
		GET: "bg-primary-green-soft text-primary-green-deep",
		POST: "bg-emerald-100 text-emerald-900",
		PUT: "bg-amber-100 text-amber-900",
		PATCH: "bg-lime-100 text-lime-900",
		DELETE: "bg-rose-100 text-rose-900",
	} as const;

	const shouldLock = $derived(mode !== "authenticated");
</script>

<svelte:head>
	<title>{detail.title} | API Testing Kit</title>
	<meta
		name="description"
		content={`Collection detail view for ${detail.title}. ${detail.description}`}
	/>
</svelte:head>

<section class="space-y-4">
	<Card class="panel-card overflow-hidden">
		<CardHeader class="gap-4 border-b border-border/70 bg-gradient-to-br from-white via-white to-primary-green-soft/35">
			<div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
				<div class="max-w-3xl space-y-3">
					<div class="flex flex-wrap items-center gap-2 text-xs font-medium text-text-muted">
						<a class="hover:text-text-strong" href="/app">App</a>
						<span>/</span>
						<a class="hover:text-text-strong" href="/app">Collections</a>
						<span>/</span>
						<span class="text-text-strong">{detail.title}</span>
					</div>

					<div class="flex flex-wrap items-center gap-2">
						<Badge variant="outline" class="bg-white/80">
							<Layers3Icon class="size-3.5" />
							{detail.badge}
						</Badge>
						<Badge variant="secondary">{detail.accessLabel}</Badge>
						<Badge variant={badgeTone[detail.status]}>
							{detail.status === "missing"
								? "Unavailable"
								: detail.status === "unavailable"
									? "Unavailable"
									: detail.status === "ready"
										? "Ready"
										: "Preview"}
						</Badge>
					</div>

					<div class="space-y-2">
						<CardTitle class="text-2xl tracking-tight sm:text-3xl">{detail.title}</CardTitle>
						<CardDescription class="max-w-3xl text-sm leading-6 text-text-body">
							{detail.description}
						</CardDescription>
						<p class="max-w-3xl text-sm leading-6 text-text-body">{detail.heroCopy}</p>
					</div>
				</div>

				<div class="rounded-[20px] border border-border/70 bg-white/85 p-4 shadow-sm lg:min-w-[18rem]">
					<div class="flex items-center gap-2">
						<SparklesIcon class="size-4 text-primary-green" />
						<p class="text-sm font-semibold text-text-strong">{sessionLabel}</p>
					</div>
					<p class="mt-2 text-sm leading-6 text-text-body">{detail.accessCopy}</p>
					<div class="mt-4 flex flex-wrap gap-2">
						{#each detail.actionLinks as action}
							<Button href={action.href} variant={action.variant} size="sm">
								{action.label}
								<ArrowRightIcon class="size-4" />
							</Button>
						{/each}
					</div>
				</div>
			</div>
		</CardHeader>
	</Card>

	{#if entitlements && mode === "authenticated"}
		<Card class="panel-card">
			<CardHeader class="gap-3">
				<div class="flex items-center justify-between gap-3">
					<div>
						<CardTitle>Signed-in capabilities</CardTitle>
						<CardDescription>{getEntitlementSummary(entitlements, mode)}</CardDescription>
					</div>
					<Badge variant="outline">{entitlements.plan.name}</Badge>
				</div>
			</CardHeader>
			<CardContent class="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
				{#each buildEntitlementRows(entitlements) as row}
					<div class="rounded-[20px] border border-border/70 bg-panel-soft p-4">
						<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">{row.label}</p>
						<p class="mt-2 text-base font-semibold tracking-tight text-text-strong">{row.statusLabel}</p>
						<p class="mt-1 text-sm leading-6 text-text-body">{row.description}</p>
						{#if row.limitLabel}
							<p class="mt-2 text-xs font-medium text-text-muted">{row.limitLabel}</p>
						{/if}
					</div>
				{/each}
			</CardContent>
		</Card>
	{/if}

	{#if detail.status === "missing" || detail.status === "unavailable"}
		<Card class="panel-card">
			<CardHeader class="gap-3">
				<CardTitle>{detail.status === "unavailable" ? "Collection unavailable" : "Collection not found"}</CardTitle>
				<CardDescription>{detail.description}</CardDescription>
			</CardHeader>
			<CardContent class="flex flex-wrap gap-2">
				<Button href="/app" variant="outline">Back to app</Button>
				<Button href="/templates" class="gap-2">
					<RouteIcon class="size-4" />
					Browse templates
				</Button>
			</CardContent>
		</Card>
	{:else}
		<div class="grid gap-4 xl:grid-cols-[minmax(0,1.32fr)_minmax(320px,0.68fr)]">
			<div class="space-y-4">
				<Card class="panel-card">
					<CardHeader class="gap-3">
						<div class="flex items-center justify-between gap-3">
							<div>
								<CardTitle>Request groups</CardTitle>
								<CardDescription>
									{#if shouldLock}
										Guest mode keeps the collection shape visible, but execution and persistence remain locked.
									{:else}
										Authenticated users get live saved requests with reusable launch paths back into `/app`.
									{/if}
								</CardDescription>
							</div>
							<Badge variant="outline">{detail.requestGroups.length} groups</Badge>
						</div>
					</CardHeader>
					<CardContent class="space-y-4">
						{#each detail.requestGroups as group}
							<div class="rounded-[20px] border border-border/70 bg-panel-soft p-4">
								<div class="flex items-start justify-between gap-3">
									<div>
										<p class="text-sm font-semibold text-text-strong">{group.label}</p>
										<p class="mt-1 text-sm leading-6 text-text-body">{group.summary}</p>
									</div>
									<Badge variant="secondary">{group.requestCount} requests</Badge>
								</div>

								<div class="mt-4 space-y-3">
									{#each group.requests as request}
										<div class="rounded-[18px] border border-border/70 bg-white px-4 py-4">
											<div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
												<div class="space-y-2">
													<div class="flex flex-wrap items-center gap-2">
														<Badge variant="outline" class={requestTone[request.method as keyof typeof requestTone]}>
															{request.method}
														</Badge>
														<Badge variant="secondary">{request.category}</Badge>
													</div>
													<div>
														<p class="text-sm font-semibold text-text-strong">{request.title}</p>
														<p class="mt-1 text-sm leading-6 text-text-body">{request.summary}</p>
													</div>
													<p class="font-mono text-xs text-text-muted">{request.endpoint}</p>
												</div>

												{#if shouldLock}
													<div class="rounded-[16px] border border-dashed border-border/70 bg-panel-soft px-3 py-3 text-sm leading-6 text-text-body lg:max-w-xs">
														<div class="flex items-center gap-2 text-text-strong">
															<LockKeyholeIcon class="size-4 text-warning" />
															Gated in guest mode
														</div>
														<p class="mt-2 text-xs leading-5 text-text-muted">
															Open actions move through the workspace shell, but persisted collection execution is reserved for signed-in users.
														</p>
													</div>
												{:else}
													<div class="flex flex-wrap gap-2">
														<Button href={request.launchHref} size="sm" class="gap-2">
															Open in /app
															<ArrowRightIcon class="size-4" />
														</Button>
														{#if request.secondaryHref}
															<Button href={request.secondaryHref} variant="outline" size="sm" class="gap-2">
																{request.secondaryLabel ?? "Preview"}
																<ExternalLinkIcon class="size-4" />
															</Button>
														{/if}
													</div>
												{/if}
											</div>

													<div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
												<div class="rounded-[16px] border border-border/70 bg-panel-soft px-3 py-3">
													<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">Response</p>
													<p class="mt-2 text-sm font-semibold text-text-strong">{request.responseLabel}</p>
													<p class="mt-1 text-xs text-text-muted">{request.responseDetail}</p>
												</div>
												<div class="rounded-[16px] border border-border/70 bg-panel-soft px-3 py-3">
													<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">Duration</p>
													<p class="mt-2 text-sm font-semibold text-text-strong">{request.durationLabel}</p>
													<p class="mt-1 text-xs text-text-muted">{request.durationDetail}</p>
												</div>
												<div class="rounded-[16px] border border-border/70 bg-panel-soft px-3 py-3">
													<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">Size</p>
													<p class="mt-2 text-sm font-semibold text-text-strong">{request.sizeLabel}</p>
													<p class="mt-1 text-xs text-text-muted">{request.sizeDetail}</p>
												</div>
												<div class="rounded-[16px] border border-border/70 bg-panel-soft px-3 py-3">
													<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">Overrides</p>
													<p class="mt-2 text-sm font-semibold text-text-strong">
														{request.safeOverrides.length > 0 ? request.safeOverrides.join(", ") : "None"}
													</p>
													<p class="mt-1 text-xs text-text-muted">
														{request.safeOverrides.length > 0
															? "Template-defined safe fields."
															: "No safe override metadata was stored for this request."}
													</p>
												</div>
											</div>
										</div>
									{/each}
								</div>
							</div>
						{/each}
					</CardContent>
				</Card>
			</div>

			<div class="space-y-4">
				<Card class="panel-card">
					<CardHeader class="gap-3">
						<div class="flex items-center justify-between gap-3">
							<div>
								<CardTitle>Collection metadata</CardTitle>
								<CardDescription>
									{detail.source === "live"
										? "Live collection facts returned by the authenticated API."
										: "Route-level collection facts for the current preview state."}
								</CardDescription>
							</div>
							<Badge variant="outline">{detail.scope}</Badge>
						</div>
					</CardHeader>
					<CardContent class="space-y-3">
						{#each detail.metadata as item}
							<div class="rounded-[18px] border border-border/70 bg-panel-soft px-4 py-3">
								<div class="flex items-center justify-between gap-3">
									<p class="text-sm font-semibold text-text-strong">{item.label}</p>
									<Badge variant="outline">{item.value}</Badge>
								</div>
								<p class="mt-2 text-xs leading-5 text-text-muted">{item.detail}</p>
							</div>
						{/each}
					</CardContent>
				</Card>

				<Card class="panel-card">
					<CardHeader class="gap-3">
						<CardTitle>Related collections</CardTitle>
						<CardDescription>
							{detail.source === "live"
								? "Other persisted collections visible to this signed-in account."
								: "Nearby collections visible in the current preview mode."}
						</CardDescription>
					</CardHeader>
					<CardContent class="space-y-3">
						{#if detail.relatedCollections.length > 0}
							{#each detail.relatedCollections as collection}
								<div class="rounded-[18px] border border-border/70 bg-panel-soft px-4 py-3">
									<div class="flex items-start justify-between gap-3">
										<div>
											<p class="text-sm font-semibold text-text-strong">{collection.title}</p>
											<p class="mt-1 text-xs leading-5 text-text-muted">{collection.description}</p>
										</div>
										<Badge variant="secondary">{collection.badge}</Badge>
									</div>
									<div class="mt-3 flex flex-wrap gap-2">
										<Button href={`/app/collections/${collection.id}`} size="sm" variant="outline">
											Open
											<ArrowRightIcon class="size-4" />
										</Button>
									</div>
								</div>
							{/each}
						{:else}
							<div class="rounded-[18px] border border-border/70 bg-panel-soft px-4 py-3 text-sm leading-6 text-text-body">
								No related collections are visible in the current workspace mode.
							</div>
						{/if}
					</CardContent>
				</Card>

				<Card class="panel-card border-dashed">
					<CardHeader class="gap-3">
						<div class="flex items-center gap-2">
							{#if shouldLock}
								<LockKeyholeIcon class="size-4 text-warning" />
							{:else}
								<ShieldCheckIcon class="size-4 text-success" />
							{/if}
							<CardTitle>{shouldLock ? "Guest fallback" : "Execution ready"}</CardTitle>
						</div>
						<CardDescription>
							{#if shouldLock}
								The collection detail page stays honest: visitors can inspect the shape, but authenticated execution is still required.
							{:else}
								The collection is open in the signed-in shell and can be launched from the same workspace model.
							{/if}
						</CardDescription>
					</CardHeader>
					<CardContent class="space-y-3">
						<div class="rounded-[18px] border border-border/70 bg-panel-soft px-4 py-3">
							<p class="text-sm font-semibold text-text-strong">{detail.accessLabel}</p>
							<p class="mt-1 text-sm leading-6 text-text-body">{detail.accessCopy}</p>
						</div>
						<div class="flex flex-wrap gap-2">
							<Button href={detail.collectionHref} size="sm" class="gap-2">
								Open in /app
								<ArrowRightIcon class="size-4" />
							</Button>
							{#if detail.previewHref}
								<Button href={detail.previewHref} size="sm" variant="outline" class="gap-2">
									Preview in /app
									<ExternalLinkIcon class="size-4" />
								</Button>
							{/if}
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	{/if}
</section>
