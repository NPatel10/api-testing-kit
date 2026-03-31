<script lang="ts">
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle,
	} from "$lib/components/ui/card/index.js";
	import Separator from "$lib/components/ui/separator/separator.svelte";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CheckIcon from "@lucide/svelte/icons/check";
	import Code2Icon from "@lucide/svelte/icons/code-2";
	import LockIcon from "@lucide/svelte/icons/lock";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";

	const overviewCards = [
		{
			label: "Route",
			value: "/app",
			detail: "Guest and signed-in sessions share the same workspace.",
		},
		{
			label: "Execution",
			value: "Server validated",
			detail: "Outbound requests pass through destination, size, and timeout checks.",
		},
		{
			label: "Templates",
			value: "Allowlisted",
			detail: "Preset request shapes launch into the same workspace.",
		},
		{
			label: "Persistence",
			value: "Mode aware",
			detail: "Collections and history unlock in signed-in sessions.",
		},
	];

	const modeCards = [
		{
			title: "Guest mode",
			badge: "Allowlisted",
			body:
				"Templates, previews, and visible lock states stay available. Custom URLs, saved history, and variables stay locked.",
		},
		{
			title: "Authenticated mode",
			badge: "Validated",
			body:
				"Custom request targets, persistence, and broader limits unlock without changing the workspace layout.",
		},
		{
			title: "Safety model",
			badge: "Fail closed",
			body:
				"Blocked destinations, capped payloads, and timeout controls still apply after sign-in.",
		},
	];

	const requestModel = [
		{
			title: "Method",
			text: "GET, POST, PUT, PATCH, and DELETE are available from the same request surface.",
		},
		{
			title: "URL",
			text: "Guest sessions stay on allowlisted targets. Signed-in sessions can submit validated custom URLs.",
		},
		{
			title: "Params and headers",
			text: "Query values and headers are edited inline with the rest of the request state.",
		},
		{
			title: "Auth and body",
			text: "No auth, basic auth, bearer tokens, JSON, raw text, and form data share one editor.",
		},
	];

	const responseFacts = [
		"Status and failure state",
		"Response time and payload size",
		"Headers and content type",
		"Pretty, raw, and header views",
	];

	const templateFacts = [
		"Templates keep fixed allowlisted targets.",
		"Guest-safe overrides remain limited to exposed fields.",
		"Collections group request sets on the same route model.",
	];

	const safetyFacts = [
		"Localhost, private ranges, and metadata IPs are blocked.",
		"Redirect hops are revalidated before they are followed.",
		"Request and response sizes stay capped.",
		"Timeouts and redirect depth stay bounded.",
	];
</script>

<svelte:head>
	<title>Docs - API Testing Kit</title>
	<meta
		name="description"
		content="Reference page for the shared workspace, access modes, request model, response surface, templates, and safety rules in API Testing Kit."
	/>
</svelte:head>

<div class="relative isolate overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(31,122,77,0.15),_transparent_30%),radial-gradient(circle_at_top_right,_rgba(111,142,163,0.12),_transparent_28%),linear-gradient(180deg,_#f2f0ea_0%,_#ece8df_100%)] text-text-strong">
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-32 top-20 h-80 w-80 rounded-full bg-[#1f7a4d]/10 blur-3xl"></div>
		<div class="absolute right-[-5rem] top-1/3 h-96 w-96 rounded-full bg-[#dcefe3] blur-3xl"></div>
		<div class="absolute bottom-[-8rem] left-1/2 h-72 w-72 -translate-x-1/2 rounded-full bg-[#145336]/10 blur-3xl"></div>
	</div>

	<div class="mx-auto min-h-screen max-w-[1440px] px-4 py-4 sm:px-6 lg:px-8">
		<div class="overflow-hidden rounded-[32px] border border-[#e7e3d8] bg-[rgba(247,245,240,0.92)] shadow-[0_24px_60px_rgba(21,31,23,0.08)] backdrop-blur">
			<div class="border-b border-[#e7e3d8] bg-white/75 px-5 py-4 sm:px-6 lg:px-8">
				<div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
					<div class="flex items-center gap-3">
						<div class="grid h-11 w-11 place-items-center rounded-2xl bg-primary-green text-sm font-semibold text-white shadow-[0_10px_24px_rgba(31,122,77,0.28)]">
							AT
						</div>
						<div>
							<p class="text-sm font-semibold tracking-tight text-text-strong">API Testing Kit Docs</p>
							<p class="text-xs text-text-muted">Workspace reference</p>
						</div>
					</div>

					<div class="flex flex-wrap items-center gap-2 text-sm">
						<a href="#overview" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Overview</a>
						<a href="#modes" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Modes</a>
						<a href="#request-model" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Requests</a>
						<a href="#safety" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Safety</a>
						<a href="/app" class="rounded-full bg-primary-green px-5 py-2.5 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-primary-green-hover">Open app</a>
					</div>
				</div>
			</div>

			<div class="grid gap-8 px-5 py-6 lg:grid-cols-[minmax(0,1fr)_320px] lg:px-8 lg:py-8">
				<main class="space-y-10">
					<section id="overview" class="grid gap-6 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,0.9fr)]">
						<div class="space-y-5">
							<Badge variant="secondary" class="w-fit bg-primary-green-soft text-primary-green-deep">Reference</Badge>
							<div class="space-y-4">
								<h1 class="max-w-3xl text-4xl font-semibold tracking-tight sm:text-5xl">
									Reference for the shared request workspace
								</h1>
								<p class="max-w-2xl text-sm leading-7 text-text-body sm:text-base">
									Core route structure, access modes, request model, response model, templates, and safety constraints.
								</p>
							</div>

							<div class="flex flex-wrap gap-3">
								<Button href="/app" size="lg" class="rounded-full bg-primary-green px-6 text-white hover:bg-primary-green-hover">
									Open `/app`
								</Button>
								<Button href="/templates" variant="outline" size="lg" class="rounded-full border-border bg-white px-6">
									<SparklesIcon class="size-4" />
									Templates
								</Button>
							</div>

							<div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
								{#each overviewCards as item}
									<div class="metric-card p-4">
										<p class="text-xs uppercase tracking-[0.24em] text-text-muted">{item.label}</p>
										<p class="mt-2 text-sm font-semibold">{item.value}</p>
										<p class="mt-1 text-xs leading-5 text-text-muted">{item.detail}</p>
									</div>
								{/each}
							</div>
						</div>

						<Card class="panel-card overflow-hidden">
							<CardHeader class="gap-3 border-b border-[#ece6db] bg-surface-soft/70">
								<div class="flex items-center justify-between gap-3">
									<div>
										<CardTitle>Workspace facts</CardTitle>
										<CardDescription>Core product boundaries and route behavior.</CardDescription>
									</div>
									<Badge variant="outline" class="border-primary-green-soft bg-white text-primary-green-deep">Shared shell</Badge>
								</div>
							</CardHeader>
							<CardContent class="space-y-4 p-5">
								{#each overviewCards as item}
									<div class="rounded-[20px] border border-border/70 bg-white p-4">
										<p class="text-sm font-semibold">{item.label}</p>
										<p class="mt-1 text-sm leading-6 text-text-body">{item.detail}</p>
									</div>
								{/each}
							</CardContent>
						</Card>
					</section>

					<section id="modes" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">Modes</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">Access is mode-based, not route-based</h2>
							<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
								The workspace route stays constant while capabilities change by session state.
							</p>
						</div>

						<div class="grid gap-4 md:grid-cols-3">
							{#each modeCards as card}
								<Card class="panel-card">
									<CardHeader class="gap-2">
										<div class="flex items-center justify-between gap-3">
											<CardTitle class="text-base">{card.title}</CardTitle>
											<Badge variant="outline">{card.badge}</Badge>
										</div>
									</CardHeader>
									<CardContent>
										<p class="rounded-[18px] border border-border/70 bg-white px-4 py-4 text-sm leading-6 text-text-body">
											{card.body}
										</p>
									</CardContent>
								</Card>
							{/each}
						</div>
					</section>

					<section id="request-model" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">Request model</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">The request surface keeps all major inputs together</h2>
						</div>

						<div class="grid gap-4 md:grid-cols-2">
							{#each requestModel as item}
								<Card class="panel-card">
									<CardHeader class="gap-2">
										<CardTitle class="text-base">{item.title}</CardTitle>
									</CardHeader>
									<CardContent>
										<div class="rounded-[18px] border border-border/70 bg-white px-4 py-4 text-sm leading-6 text-text-body">
											{item.text}
										</div>
									</CardContent>
								</Card>
							{/each}
						</div>

						<div class="code-surface">
							<p class="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-text-muted">Request example</p>
							<pre class="overflow-x-auto text-xs leading-6"><code>GET https://jsonplaceholder.typicode.com/posts/1
Accept: application/json
X-Demo-Mode: guest</code></pre>
						</div>
					</section>

					<section id="responses" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">Response model</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">Responses remain attached to request context</h2>
						</div>

						<div class="grid gap-4 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]">
							<Card class="panel-card">
								<CardHeader class="gap-2">
									<CardTitle class="text-base">Response fields</CardTitle>
									<CardDescription>Status, metadata, and payload views.</CardDescription>
								</CardHeader>
								<CardContent class="space-y-3">
									{#each responseFacts as item}
										<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
											<CheckIcon class="mt-0.5 size-4 text-success" />
											<p class="text-sm leading-6 text-text-body">{item}</p>
										</div>
									{/each}
								</CardContent>
							</Card>

							<div class="grid gap-4">
								<div class="metric-card p-5">
									<div class="flex items-center justify-between gap-3">
										<div>
											<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Status</p>
											<p class="mt-2 text-2xl font-semibold">200 OK</p>
										</div>
										<Badge class="bg-primary-green-soft text-primary-green-deep">Success</Badge>
									</div>
									<Separator class="my-4 bg-border/80" />
									<div class="grid gap-3 sm:grid-cols-3">
										<div class="rounded-[16px] border border-border/70 bg-surface-soft px-4 py-3">
											<p class="text-xs uppercase tracking-[0.18em] text-text-muted">Time</p>
											<p class="mt-1 text-sm font-semibold">186 ms</p>
										</div>
										<div class="rounded-[16px] border border-border/70 bg-surface-soft px-4 py-3">
											<p class="text-xs uppercase tracking-[0.18em] text-text-muted">Size</p>
											<p class="mt-1 text-sm font-semibold">1.2 KB</p>
										</div>
										<div class="rounded-[16px] border border-border/70 bg-surface-soft px-4 py-3">
											<p class="text-xs uppercase tracking-[0.18em] text-text-muted">Type</p>
											<p class="mt-1 text-sm font-semibold">JSON</p>
										</div>
									</div>
								</div>

								<div class="code-surface">
									<p class="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-text-muted">Pretty view</p>
									<pre class="overflow-x-auto text-xs leading-6"><code>&#123;
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
&#125;</code></pre>
								</div>
							</div>
						</div>
					</section>

					<section id="templates" class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
						<Card class="panel-card">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">Templates</CardTitle>
								<CardDescription>Preset request and response shapes.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								{#each templateFacts as item}
									<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
										<Code2Icon class="mt-0.5 size-4 text-primary-green" />
										<p class="text-sm leading-6 text-text-body">{item}</p>
									</div>
								{/each}
							</CardContent>
						</Card>

						<Card class="panel-card">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">Guest scope</CardTitle>
								<CardDescription>Visible shell, limited execution.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
									<ShieldCheckIcon class="mt-0.5 size-4 text-success" />
									<p class="text-sm leading-6 text-text-body">Templates and collections remain visible on the shared route.</p>
								</div>
								<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
									<ShieldCheckIcon class="mt-0.5 size-4 text-success" />
									<p class="text-sm leading-6 text-text-body">Custom targets, persistence, and variables stay locked.</p>
								</div>
								<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
									<ShieldCheckIcon class="mt-0.5 size-4 text-success" />
									<p class="text-sm leading-6 text-text-body">Signed-in access unlocks the same surfaces without changing routes.</p>
								</div>
							</CardContent>
						</Card>
					</section>

					<section id="safety" class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
						<Card class="panel-card">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">Safety rules</CardTitle>
								<CardDescription>Validation happens before execution.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								{#each safetyFacts as item}
									<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
										<LockIcon class="mt-0.5 size-4 text-primary-green" />
										<p class="text-sm leading-6 text-text-body">{item}</p>
									</div>
								{/each}
							</CardContent>
						</Card>

						<Card class="panel-card bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.98))]">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">Locked states</CardTitle>
								<CardDescription>Visible controls, restricted actions.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								<div class="rounded-[18px] border border-border/70 bg-white p-4">
									<div class="flex items-start gap-3">
										<div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-primary-green-soft text-primary-green-deep">
											<LockIcon class="size-4" />
										</div>
										<div>
											<p class="text-sm font-semibold">Lock surfaces remain in place.</p>
											<p class="mt-1 text-sm leading-6 text-text-body">
												Guests see the control layout, while execution and persistence stay mode-aware.
											</p>
										</div>
									</div>
								</div>

								<div class="flex flex-wrap gap-3">
									<Button href="/app" size="sm" class="rounded-full bg-primary-green px-4 text-white hover:bg-primary-green-hover">
										<ArrowRightIcon class="size-4" />
										Open `/app`
									</Button>
									<Button href="/features" variant="outline" size="sm" class="rounded-full border-border bg-white px-4">
										Reference features
									</Button>
								</div>
							</CardContent>
						</Card>
					</section>
				</main>

				<aside class="space-y-4 lg:sticky lg:top-6 lg:self-start">
					<Card class="panel-card">
						<CardHeader class="gap-2">
							<CardTitle class="text-base">On this page</CardTitle>
							<CardDescription>Section links for the reference view.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-2">
							<a href="#overview" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Overview</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#modes" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Modes</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#request-model" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Request model</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#responses" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Response model</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#templates" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Templates</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#safety" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Safety</span>
								<ArrowRightIcon class="size-4" />
							</a>
						</CardContent>
					</Card>

					<Card class="panel-card">
						<CardHeader class="gap-2">
							<CardTitle class="text-base">Route map</CardTitle>
							<CardDescription>Public and signed-in entry points.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-3 text-sm leading-6 text-text-body">
							<p>`/` for the public product surface.</p>
							<p>`/templates` for presets and collections.</p>
							<p>`/app` for request editing and response inspection.</p>
							<p>`/app/history` for persisted or preview request history.</p>
						</CardContent>
					</Card>

					<Card class="panel-card bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.98))]">
						<CardHeader class="gap-2">
							<CardTitle class="text-base">Live workspace</CardTitle>
							<CardDescription>The request builder and response viewer live on `/app`.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-3">
							<Button href="/app" class="w-full justify-between rounded-full bg-primary-green px-5 text-white hover:bg-primary-green-hover">
								<span>Open `/app`</span>
								<ArrowRightIcon class="size-4" />
							</Button>
						</CardContent>
					</Card>
				</aside>
			</div>
		</div>
	</div>
</div>
