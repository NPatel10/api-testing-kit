<script lang="ts">
	import { resolve } from "$app/paths";
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle,
	} from "$lib/components/ui/card/index.js";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import Separator from "$lib/components/ui/separator/separator.svelte";
	import HistoryIcon from "@lucide/svelte/icons/history";
	import LockIcon from "@lucide/svelte/icons/lock";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
	import SearchIcon from "@lucide/svelte/icons/search";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";
	import type { HistoryPageData } from "./history-page-data";
	import {
		historyMethodFilters,
		historyMethodLabels,
		historyStatusFilters,
		historyStatusLabels,
	} from "./history-page-data";

	let { data }: { data: HistoryPageData } = $props();

	const statusTone: Record<"all" | "success" | "blocked" | "error", "default" | "secondary" | "outline"> = {
		success: "default",
		blocked: "secondary",
		error: "outline",
		all: "outline",
	};
</script>

<svelte:head>
	<title>History - API Testing Kit</title>
	<meta
		name="description"
		content="Authenticated request history and guest-safe previews for API Testing Kit."
	/>
</svelte:head>

<section class="relative isolate overflow-hidden rounded-[32px] border border-[#dfe8dd] bg-[linear-gradient(135deg,rgba(31,122,77,0.10),rgba(255,255,255,0.98))] text-text-strong shadow-[0_24px_60px_rgba(21,31,23,0.08)]">
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-24 top-10 h-64 w-64 rounded-full bg-[#1f7a4d]/10 blur-3xl"></div>
		<div class="absolute right-[-4rem] top-28 h-80 w-80 rounded-full bg-[#dcefe3] blur-3xl"></div>
	</div>

	<div class="relative border-b border-[#dfe8dd] bg-white/80 px-5 py-5 sm:px-6 lg:px-8">
		<div class="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
			<div class="max-w-3xl space-y-4">
				<div class="inline-flex items-center gap-2 rounded-full border border-[#d9e7d8] bg-white px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-[#1f7a4d]">
					<HistoryIcon class="size-4" />
					Request history
				</div>
				<div class="space-y-3">
					<h1 class="text-4xl font-semibold tracking-tight sm:text-5xl">
						{data.mode === "authenticated"
							? "Request history"
							: "Request history preview"}
					</h1>
					<p class="max-w-2xl text-sm leading-7 text-text-body sm:text-base">
						{data.mode === "authenticated"
							? "Persisted runs grouped by status, method, domain, timing, and response size."
							: "Preview history with mode-aware persistence locks."}
					</p>
				</div>
			</div>

			<div class="flex flex-wrap gap-3">
				<Button href={resolve("/app")} size="lg" class="rounded-full bg-primary-green px-6 text-white hover:bg-primary-green-hover">
					<SparklesIcon class="size-4" />
					Back to workspace
				</Button>
				<Button href={resolve("/docs")} variant="outline" size="lg" class="rounded-full border-border bg-white px-6">
					<SearchIcon class="size-4" />
					Docs
				</Button>
			</div>
		</div>

		<div class="mt-6 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
			{#each data.metrics as metric (metric.label)}
				<div class="rounded-[22px] border border-[#e7e3d8] bg-white/85 p-4 shadow-[0_10px_24px_rgba(21,31,23,0.04)]">
					<p class="text-xs uppercase tracking-[0.22em] text-text-muted">{metric.label}</p>
					<p class="mt-2 text-2xl font-semibold">{metric.value}</p>
					<p class="mt-1 text-xs leading-5 text-text-muted">{metric.note}</p>
				</div>
			{/each}
		</div>
	</div>

	<div class="grid gap-6 px-5 py-6 lg:grid-cols-[minmax(0,1fr)_320px] lg:px-8 lg:py-8">
		<main class="space-y-6">
			{#if data.notice}
				<Card class="border border-amber-200 bg-amber-50/80 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
					<CardHeader class="gap-2">
						<CardTitle class="text-lg">Live history unavailable</CardTitle>
						<CardDescription>{data.notice}</CardDescription>
					</CardHeader>
				</Card>
			{/if}

			<Card class="border border-[#e7e3d8] bg-white/90 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
				<CardHeader class="gap-3">
					<div class="flex flex-wrap items-center gap-2">
						<Badge variant={statusTone[data.status]}>{historyStatusLabels[data.status]}</Badge>
						<Badge variant="outline">{historyMethodLabels[data.method]}</Badge>
						<Badge variant="outline">{data.selectedDomain === "all" ? "Any domain" : data.selectedDomain}</Badge>
						<Badge variant={data.mode === "authenticated" ? "default" : "secondary"}>
							{data.mode === "authenticated" ? "Persistent history" : "Guest preview"}
						</Badge>
					</div>
					<div class="flex items-start justify-between gap-4">
						<div>
							<CardTitle class="text-xl">
								{data.mode === "authenticated" ? "Recent request timeline" : "Preview timeline"}
							</CardTitle>
							<CardDescription>
								{data.mode === "authenticated"
									? "Filtered persisted requests."
									: "Filtered preview requests."}
							</CardDescription>
						</div>
						<Button href="?status=all&method=all&domain=all" variant="ghost" class="rounded-full">
							<RefreshCwIcon class="size-4" />
							Reset
						</Button>
					</div>
				</CardHeader>
				<CardContent class="space-y-5">
					<div class="flex flex-wrap gap-2">
						{#each historyStatusFilters as filter (filter)}
							<a
								href={resolve(
									`/app/history?status=${filter}&method=${data.method}&domain=${encodeURIComponent(data.selectedDomain)}`,
								)}
								class={`rounded-full border px-4 py-2 text-sm transition ${
									data.status === filter
										? "border-primary-green bg-primary-green text-white shadow-[0_10px_24px_rgba(31,122,77,0.18)]"
										: "border-border bg-white text-text-body hover:border-primary-green-soft hover:bg-primary-green-soft/30 hover:text-text-strong"
								}`}
							>
								{historyStatusLabels[filter]}
							</a>
						{/each}
					</div>

					<div class="flex flex-wrap gap-2">
						{#each historyMethodFilters as filter (filter)}
							<a
								href={resolve(
									`/app/history?status=${data.status}&method=${filter}&domain=${encodeURIComponent(data.selectedDomain)}`,
								)}
								class={`rounded-full border px-4 py-2 text-sm transition ${
									data.method === filter
										? "border-[#145336] bg-[#145336] text-white"
										: "border-border bg-white text-text-body hover:border-primary-green-soft hover:bg-primary-green-soft/30 hover:text-text-strong"
								}`}
							>
								{historyMethodLabels[filter]}
							</a>
						{/each}
					</div>

					<div class="flex flex-wrap gap-2">
						<a
							href={resolve(`/app/history?status=${data.status}&method=${data.method}&domain=all`)}
							class={`rounded-full border px-4 py-2 text-sm transition ${
								data.selectedDomain === "all"
									? "border-[#145336] bg-[#145336] text-white"
									: "border-border bg-white text-text-body hover:border-primary-green-soft hover:bg-primary-green-soft/30 hover:text-text-strong"
							}`}
						>
							Any domain
						</a>
						{#each data.domainOptions as domain (domain)}
							<a
								href={resolve(
									`/app/history?status=${data.status}&method=${data.method}&domain=${encodeURIComponent(domain)}`,
								)}
								class={`rounded-full border px-4 py-2 text-sm transition ${
									data.selectedDomain === domain
										? "border-primary-green bg-primary-green text-white shadow-[0_10px_24px_rgba(31,122,77,0.18)]"
										: "border-border bg-white text-text-body hover:border-primary-green-soft hover:bg-primary-green-soft/30 hover:text-text-strong"
								}`}
							>
								{domain}
							</a>
						{/each}
					</div>

					<Separator class="bg-[#e7e3d8]" />

					{#if data.entries.length > 0}
						<div class="space-y-4">
							{#each data.sections as section (section.title)}
								<div class="space-y-3">
									<div class="flex items-center justify-between gap-3">
										<div>
											<p class="text-sm font-semibold text-text-strong">{section.title}</p>
											<p class="text-xs text-text-muted">{section.description}</p>
										</div>
										<Badge variant="outline">{section.entries.length}</Badge>
									</div>

									<div class="grid gap-3">
										{#each section.entries as entry (entry.id)}
											<article class="rounded-[24px] border border-[#e7e3d8] bg-[#fdfcf8] p-4 shadow-[0_8px_20px_rgba(21,31,23,0.04)]">
												<div class="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
													<div class="space-y-3">
														<div class="flex flex-wrap items-center gap-2">
															<Badge variant={entry.outcome === "success" ? "default" : entry.outcome === "blocked" ? "secondary" : "outline"}>
																{entry.statusLabel}
															</Badge>
															<Badge variant="outline">{entry.method}</Badge>
															<Badge variant="outline">{entry.domainLabel}</Badge>
														</div>
														<div>
															<p class="text-base font-semibold text-text-strong">{entry.title}</p>
															<p class="mt-1 text-sm leading-6 text-text-body">{entry.target}</p>
														</div>
														{#if entry.launchHref}
															<div>
																<Button href={entry.launchHref} size="sm" variant="outline" class="rounded-full">
																	{entry.launchLabel ?? "Open in /app"}
																	<ArrowRightIcon class="size-4" />
																</Button>
															</div>
														{/if}
													</div>

													<div class="grid grid-cols-2 gap-2 sm:grid-cols-4 xl:min-w-[360px]">
														<div class="rounded-[18px] border border-border/70 bg-white px-3 py-2">
															<p class="text-[11px] uppercase tracking-[0.18em] text-text-muted">Time</p>
															<p class="mt-1 text-sm font-semibold">{entry.durationLabel}</p>
														</div>
														<div class="rounded-[18px] border border-border/70 bg-white px-3 py-2">
															<p class="text-[11px] uppercase tracking-[0.18em] text-text-muted">Size</p>
															<p class="mt-1 text-sm font-semibold">{entry.responseSizeLabel}</p>
														</div>
														<div class="rounded-[18px] border border-border/70 bg-white px-3 py-2">
															<p class="text-[11px] uppercase tracking-[0.18em] text-text-muted">Type</p>
															<p class="mt-1 text-sm font-semibold">{entry.contentTypeLabel}</p>
														</div>
														<div class="rounded-[18px] border border-border/70 bg-white px-3 py-2">
															<p class="text-[11px] uppercase tracking-[0.18em] text-text-muted">When</p>
															<p class="mt-1 text-sm font-semibold">{entry.timeLabel}</p>
														</div>
													</div>
												</div>
											</article>
										{/each}
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<div class="rounded-[24px] border border-dashed border-[#d9e7d8] bg-[#fbfaf6] px-5 py-8 text-center">
							<p class="text-base font-semibold text-text-strong">No history matches this view</p>
							<p class="mt-2 text-sm leading-6 text-text-body">
								{data.mode === "authenticated"
									? "No persisted requests match the current filters."
									: "No preview requests match the current filters."}
							</p>
						</div>
					{/if}
				</CardContent>
			</Card>

			{#if data.previewLocked}
				<Card class="border border-[#e7e3d8] bg-[linear-gradient(135deg,rgba(31,122,77,0.10),rgba(255,255,255,0.96))] shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
					<CardHeader class="gap-3">
						<div class="flex items-center gap-3">
							<div class="grid size-10 place-items-center rounded-full bg-primary-green-soft text-primary-green-deep">
								<LockIcon class="size-4" />
							</div>
							<div>
								<CardTitle class="text-lg">History persistence stays locked for guests</CardTitle>
								<CardDescription>Guest mode does not store or replay runs.</CardDescription>
							</div>
						</div>
					</CardHeader>
					<CardContent class="flex flex-wrap gap-3">
						<Button href={resolve("/app")} class="rounded-full bg-primary-green px-5 text-white hover:bg-primary-green-hover">
							Open `/app`
						</Button>
						<Button href={resolve("/docs")} variant="outline" class="rounded-full border-border bg-white px-5">
							Docs
						</Button>
					</CardContent>
				</Card>
			{/if}
		</main>

		<aside class="space-y-4 lg:sticky lg:top-6 lg:self-start">
				<Card class="border border-[#e7e3d8] bg-white/90 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
					<CardHeader class="gap-2">
						<CardTitle class="text-base">View state</CardTitle>
						<CardDescription>Current route context.</CardDescription>
					</CardHeader>
					<CardContent class="space-y-3 text-sm leading-6 text-text-body">
						<p>History surfaces method, target, response state, duration, and payload size.</p>
						<p>{data.mode === "authenticated" ? "Signed-in sessions load persisted runs on this route." : "Guest mode keeps the route visible while persistence stays locked."}</p>
					</CardContent>
				</Card>

				<Card class="border border-[#e7e3d8] bg-white/90 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
					<CardHeader class="gap-2">
						<CardTitle class="text-base">Defaults</CardTitle>
						<CardDescription>Baseline timeline fields.</CardDescription>
					</CardHeader>
					<CardContent class="space-y-2">
						<div class="rounded-[18px] border border-border/70 bg-[#fdfcf8] px-4 py-3 text-sm text-text-body">
							Always show status first
					</div>
					<div class="rounded-[18px] border border-border/70 bg-[#fdfcf8] px-4 py-3 text-sm text-text-body">
						Show method and domain next
					</div>
					<div class="rounded-[18px] border border-border/70 bg-[#fdfcf8] px-4 py-3 text-sm text-text-body">
						Keep time and size visible
					</div>
						<div class="rounded-[18px] border border-border/70 bg-[#fdfcf8] px-4 py-3 text-sm text-text-body">
							Guest mode keeps lock states visible
						</div>
					</CardContent>
				</Card>
		</aside>
	</div>
</section>
