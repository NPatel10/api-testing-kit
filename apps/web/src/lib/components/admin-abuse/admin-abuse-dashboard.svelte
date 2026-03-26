<script lang="ts">
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import { Separator } from "$lib/components/ui/separator/index.js";
	import {
		Table,
		TableBody,
		TableCell,
		TableHead,
		TableHeader,
		TableRow,
	} from "$lib/components/ui/table/index.js";
	import {
		formatAdminTimestamp,
		type AdminAbuseDashboardData,
	} from "./admin-abuse-data";

	let { dashboard }: { dashboard: AdminAbuseDashboardData } = $props();

	const overviewCards = [
		{
			label: "Summary rows",
			getValue: (data: AdminAbuseDashboardData) => String(data.summary.length),
			detail: "Grouped by severity and category",
		},
		{
			label: "Recent events",
			getValue: (data: AdminAbuseDashboardData) => String(data.recent.length),
			detail: "Latest blocked and suspicious attempts",
		},
		{
			label: "Blocked targets",
			getValue: (data: AdminAbuseDashboardData) => String(data.blockedTargets.length),
			detail: "Active and historical denylist entries",
		},
	] as const;

	function severityVariant(severity: string) {
		switch (severity) {
			case "critical":
				return "destructive";
			case "high":
				return "secondary";
			case "medium":
				return "outline";
			default:
				return "ghost";
		}
	}

	function actionVariant(actionTaken: string) {
		return actionTaken === "blocked" ? "destructive" : "outline";
	}

	function activeVariant(isActive: boolean) {
		return isActive ? "default" : "outline";
	}
</script>

<section class="relative isolate overflow-hidden rounded-[32px] border border-[#e1d8cf] bg-[radial-gradient(circle_at_top_left,rgba(160,52,52,0.16),transparent_26%),linear-gradient(180deg,#f4efe9_0%,#ece6dd_100%)] text-[#1c1d1f] shadow-[0_24px_60px_rgba(21,31,23,0.08)]">
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-28 top-20 h-72 w-72 rounded-full bg-[#9d3b3b]/10 blur-3xl"></div>
		<div class="absolute right-[-4rem] top-1/3 h-96 w-96 rounded-full bg-[#d9c5b8]/40 blur-3xl"></div>
		<div class="absolute bottom-[-5rem] left-1/2 h-72 w-72 -translate-x-1/2 rounded-full bg-[#7f2d2d]/10 blur-3xl"></div>
	</div>

	<div class="relative mx-auto max-w-[1440px] px-4 py-4 sm:px-6 lg:px-8">
		<div class="overflow-hidden rounded-[28px] border border-[#e7ddd3] bg-[rgba(248,244,239,0.92)] backdrop-blur">
			<header class="border-b border-[#e7ddd3] bg-white/80 px-5 py-5 sm:px-6 lg:px-8">
				<div class="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
					<div class="max-w-3xl">
						<div class="inline-flex items-center gap-2 rounded-full border border-[#ead7d0] bg-[#fff8f4] px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-[#a04444]">
							<span class="h-2 w-2 rounded-full bg-[#a04444]"></span>
							Admin abuse monitor
						</div>
						<h1 class="mt-4 text-4xl font-semibold tracking-tight text-[#1c1d1f] sm:text-5xl lg:text-6xl">
							Review suspicious traffic, blocked targets, and the current denylist surface
						</h1>
						<p class="mt-4 max-w-2xl text-sm leading-7 text-[#4f534f] sm:text-base">
							This page is the internal abuse surface described in the docs. When the backend is available and the session is admin or owner, it shows live records; otherwise it falls back to a clearly labeled structural preview.
						</p>
					</div>

					<div class="flex flex-wrap gap-3 lg:justify-end">
						<a
							href="/app"
							class="inline-flex items-center gap-2 rounded-full bg-[#a04444] px-6 py-3 font-medium text-white shadow-[0_12px_28px_rgba(160,68,68,0.28)] transition hover:bg-[#8e3b3b]"
						>
							Open app
						</a>
						<a
							href="/case-study"
							class="rounded-full border border-[#d9cfc5] bg-white px-6 py-3 font-medium text-[#1c1d1f] transition hover:bg-[#f5f1eb]"
						>
							Read safety notes
						</a>
					</div>
				</div>

				<div class="mt-6 grid gap-3 sm:grid-cols-3">
					{#each overviewCards as card}
						<div class="rounded-2xl border border-[#e7ddd3] bg-white/90 p-4 shadow-[0_10px_24px_rgba(21,31,23,0.04)]">
							<p class="text-xs font-medium uppercase tracking-[0.22em] text-[#7c746b]">{card.label}</p>
							<p class="mt-2 text-3xl font-semibold tracking-tight text-[#1c1d1f]">{card.getValue(dashboard)}</p>
							<p class="mt-2 text-xs leading-5 text-[#5f645f]">{card.detail}</p>
						</div>
					{/each}
				</div>
			</header>

			<main class="space-y-10 px-5 py-6 sm:px-6 lg:px-8 lg:py-8">
				<Card class="border-[#e7ddd3] bg-[linear-gradient(135deg,rgba(160,68,68,0.08),rgba(255,255,255,0.94))]">
					<CardHeader class="gap-3">
						<div class="flex flex-wrap items-center gap-2">
							<Badge variant={dashboard.mode === "live" ? "default" : "secondary"}>{dashboard.sourceLabel}</Badge>
							<Badge variant="outline">Updated {formatAdminTimestamp(dashboard.generatedAt)}</Badge>
						</div>
						<CardTitle class="text-2xl text-[#1c1d1f]">{dashboard.mode === "live" ? "Live monitoring data" : "Preview monitoring surface"}</CardTitle>
						<CardDescription class="max-w-3xl text-sm leading-6 text-[#4f534f]">
							{dashboard.message}
						</CardDescription>
					</CardHeader>
					<CardContent class="grid gap-3 md:grid-cols-3">
						<div class="rounded-[22px] border border-[#e7ddd3] bg-white/90 p-4">
							<p class="text-xs font-semibold uppercase tracking-[0.22em] text-[#7c746b]">Access model</p>
							<p class="mt-2 text-sm font-medium text-[#1c1d1f]">{dashboard.mode === "live" ? "Admin / owner session" : "Preview only"}</p>
							<p class="mt-2 text-xs leading-5 text-[#5f645f]">
								{dashboard.mode === "live"
									? "Role checks happen before any data is returned."
									: "The surface stays honest about what the live backend will show."}
							</p>
						</div>
						<div class="rounded-[22px] border border-[#e7ddd3] bg-white/90 p-4">
							<p class="text-xs font-semibold uppercase tracking-[0.22em] text-[#7c746b]">Monitoring scope</p>
							<p class="mt-2 text-sm font-medium text-[#1c1d1f]">Abuse events and blocked targets</p>
							<p class="mt-2 text-xs leading-5 text-[#5f645f]">
								Grouped summaries, recent events, and denylist entries are surfaced together so review can move fast.
							</p>
						</div>
						<div class="rounded-[22px] border border-[#e7ddd3] bg-white/90 p-4">
							<p class="text-xs font-semibold uppercase tracking-[0.22em] text-[#7c746b]">Source of truth</p>
							<p class="mt-2 text-sm font-medium text-[#1c1d1f]">Go API + PostgreSQL</p>
							<p class="mt-2 text-xs leading-5 text-[#5f645f]">
								The page is server-fetched when possible and falls back only when the API is unavailable or the session is not privileged.
							</p>
						</div>
					</CardContent>
				</Card>

				<section class="grid gap-4 xl:grid-cols-[1fr_1.1fr]">
					<Card class="border-[#e7ddd3] bg-white/92">
						<CardHeader class="gap-2">
							<CardTitle class="text-xl text-[#1c1d1f]">Abuse summary</CardTitle>
							<CardDescription class="text-sm leading-6 text-[#5f645f]">
								Grouped by severity, category, and action taken.
							</CardDescription>
						</CardHeader>
						<CardContent class="pt-0">
							<Table>
								<TableHeader>
									<TableRow>
										<TableHead>Severity</TableHead>
										<TableHead>Category</TableHead>
										<TableHead>Action</TableHead>
										<TableHead class="text-right">Count</TableHead>
										<TableHead>Last observed</TableHead>
									</TableRow>
								</TableHeader>
								<TableBody>
									{#if dashboard.summary.length === 0}
										<TableRow>
											<TableCell colspan={5} class="py-8 text-center text-sm text-[#5f645f]">
												No abuse summary rows available.
											</TableCell>
										</TableRow>
									{:else}
										{#each dashboard.summary as row}
											<TableRow>
												<TableCell><Badge variant={severityVariant(row.severity)}>{row.severity}</Badge></TableCell>
												<TableCell class="font-medium text-[#1c1d1f]">{row.category}</TableCell>
												<TableCell><Badge variant={actionVariant(row.actionTaken)}>{row.actionTaken}</Badge></TableCell>
												<TableCell class="text-right font-semibold text-[#1c1d1f]">{row.count}</TableCell>
												<TableCell class="text-[#5f645f]">{formatAdminTimestamp(row.lastCreatedAt)}</TableCell>
											</TableRow>
										{/each}
									{/if}
								</TableBody>
							</Table>
						</CardContent>
					</Card>

					<Card class="border-[#e7ddd3] bg-white/92">
						<CardHeader class="gap-2">
							<CardTitle class="text-xl text-[#1c1d1f]">Recent events</CardTitle>
							<CardDescription class="text-sm leading-6 text-[#5f645f]">
								The latest blocked or suspicious attempts recorded by the backend.
							</CardDescription>
						</CardHeader>
						<CardContent class="pt-0">
							<div class="space-y-4">
								{#if dashboard.recent.length === 0}
									<div class="rounded-[22px] border border-dashed border-[#d9cfc5] bg-[#faf7f3] p-6 text-sm text-[#5f645f]">
										No recent events available.
									</div>
								{:else}
									{#each dashboard.recent as event, index}
										<div class="rounded-[22px] border border-[#e7ddd3] bg-[#faf7f3] p-4 shadow-[0_10px_24px_rgba(21,31,23,0.03)]">
											<div class="flex flex-wrap items-start justify-between gap-3">
												<div class="space-y-2">
													<div class="flex flex-wrap items-center gap-2">
														<Badge variant={severityVariant(event.severity)}>{event.severity}</Badge>
														<Badge variant={actionVariant(event.actionTaken)}>{event.actionTaken}</Badge>
														<span class="text-xs uppercase tracking-[0.22em] text-[#7c746b]"># {index + 1}</span>
													</div>
													<p class="text-sm font-semibold text-[#1c1d1f]">{event.category}</p>
													<p class="text-sm leading-6 text-[#4f534f]">{event.message}</p>
												</div>
												<div class="text-right text-xs text-[#5f645f]">
													<p>{formatAdminTimestamp(event.createdAt)}</p>
													<p class="mt-1 font-mono text-[11px] uppercase tracking-[0.22em] text-[#7c746b]">{event.ruleKey}</p>
												</div>
											</div>

											<Separator class="my-4 bg-[#e7ddd3]" />

											<div class="grid gap-3 sm:grid-cols-2">
												<div>
													<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-[#7c746b]">Target</p>
													<p class="mt-1 text-sm text-[#1c1d1f]">{event.target ?? "Not recorded"}</p>
												</div>
												<div>
													<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-[#7c746b]">Source IP</p>
													<p class="mt-1 text-sm text-[#1c1d1f]">{event.sourceIp ?? "Not recorded"}</p>
												</div>
											</div>
										</div>
									{/each}
								{/if}
							</div>
						</CardContent>
					</Card>
				</section>

				<Card class="border-[#e7ddd3] bg-white/92">
					<CardHeader class="gap-2">
						<CardTitle class="text-xl text-[#1c1d1f]">Blocked targets</CardTitle>
						<CardDescription class="text-sm leading-6 text-[#5f645f]">
							Active and historical denylist entries surfaced by the backend.
						</CardDescription>
					</CardHeader>
					<CardContent class="pt-0">
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead>Target</TableHead>
									<TableHead>Value</TableHead>
									<TableHead>Reason</TableHead>
									<TableHead>Source</TableHead>
									<TableHead>Status</TableHead>
									<TableHead>Updated</TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								{#if dashboard.blockedTargets.length === 0}
									<TableRow>
										<TableCell colspan={6} class="py-8 text-center text-sm text-[#5f645f]">
											No blocked targets are currently configured.
										</TableCell>
									</TableRow>
								{:else}
									{#each dashboard.blockedTargets as target}
										<TableRow>
											<TableCell class="font-medium text-[#1c1d1f]">{target.targetType}</TableCell>
											<TableCell class="font-mono text-sm text-[#1c1d1f]">{target.targetValue}</TableCell>
											<TableCell class="text-[#4f534f]">{target.reason}</TableCell>
											<TableCell class="text-[#4f534f]">{target.source}</TableCell>
											<TableCell><Badge variant={activeVariant(target.isActive)}>{target.isActive ? "active" : "inactive"}</Badge></TableCell>
											<TableCell class="text-[#5f645f]">{formatAdminTimestamp(target.updatedAt)}</TableCell>
										</TableRow>
									{/each}
								{/if}
							</TableBody>
						</Table>
					</CardContent>
				</Card>

				<div class="rounded-[28px] border border-[#ead7d0] bg-[linear-gradient(135deg,rgba(160,68,68,0.12),rgba(255,255,255,0.96))] p-6 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
					<div class="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
						<div class="max-w-2xl">
							<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#a04444]">Operational note</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#1c1d1f]">
								This surface is intentionally internal, not a public abuse console
							</h2>
							<p class="mt-3 text-sm leading-6 text-[#4f534f] sm:text-base">
								The page exists so the project can show a real monitoring flow without exposing it to guests. If a session is not privileged, the preview still preserves the structure but not the live records.
							</p>
						</div>

						<div class="flex flex-wrap gap-3">
							<a
								href="/docs"
								class="rounded-full border border-[#d9cfc5] bg-white px-6 py-3 font-medium text-[#1c1d1f] transition hover:bg-[#f5f1eb]"
							>
								Read quick start
							</a>
							<a
								href="/case-study"
								class="rounded-full bg-[#a04444] px-6 py-3 font-medium text-white shadow-[0_12px_28px_rgba(160,68,68,0.28)] transition hover:bg-[#8e3b3b]"
							>
								Review the architecture
							</a>
						</div>
					</div>
				</div>
			</main>
		</div>
	</div>
</section>
