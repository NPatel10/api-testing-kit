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
	import PlayIcon from "@lucide/svelte/icons/play";
	import SearchIcon from "@lucide/svelte/icons/search";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";

	const quickStart = [
		{
			step: "1",
			title: "Open `/app`",
			description:
				"Guests land on the same workspace as signed-in users. The difference is capability: guest mode stays on allowlisted templates, while authenticated mode unlocks custom execution.",
		},
		{
			step: "2",
			title: "Pick a template or build a request",
			description:
				"Start from a curated guest-safe example, or shape your own request once you are signed in. The editor keeps method, URL, params, headers, auth, and body visible together.",
		},
		{
			step: "3",
			title: "Send through backend checks",
			description:
				"Authenticated requests still get destination validation, redirect checks, payload limits, and timeout limits. The backend is not a generic proxy.",
		},
		{
			step: "4",
			title: "Read the response in one pass",
			description:
				"Use status, duration, size, headers, and the pretty/raw body tabs to decide whether the request is useful, blocked, or broken.",
		},
	];

	const modeCards = [
		{
			title: "Guest mode",
			badge: "Allowlisted only",
			body:
				"Guests can explore the live app, browse templates, inspect responses, and edit only the safe fields the template exposes. They cannot swap in an arbitrary target URL, save durable history, or use secret environment variables.",
		},
		{
			title: "Authenticated mode",
			badge: "Full execution",
			body:
				"Signed-in users keep the same shell and unlock custom URLs, saved requests, history, collections, and richer request editors. The server still revalidates every target before sending anything outbound.",
		},
		{
			title: "Safety model",
			badge: "SSRf-aware",
			body:
				"The docs and backend agree on the same rules: block localhost, private ranges, metadata IPs, unsupported protocols, large payloads, and long-running requests.",
		},
	];

	const requestGuide = [
		{
			title: "Method",
			text: "Use GET for inspection and the write methods when you need to verify a mutation path. The app supports GET, POST, PUT, PATCH, and DELETE.",
		},
		{
			title: "URL",
			text: "Guest mode keeps the target locked to curated templates. Authenticated mode accepts a custom URL, but the backend validates host, IP, protocol, port, and redirects.",
		},
		{
			title: "Params and headers",
			text: "Keep overrides explicit. That keeps the request easy to reason about and makes retries, snippets, and debugging much simpler.",
		},
		{
			title: "Auth and body",
			text: "Use no auth, basic auth, or bearer tokens as needed. Bodies support JSON, raw text, and form-urlencoded modes.",
		},
	];

	const responseGuide = [
		"Status and failure state first",
		"Response time and payload size next",
		"Headers and content type for context",
		"Pretty, raw, and headers views for the payload",
	];

	const templateGuide = [
		"Templates are curated examples, not free-form proxies.",
		"Each template keeps a fixed allowlisted target.",
		"Guest-safe overrides are limited to the fields the template exposes.",
		"Templates are meant to get a user productive fast, not hide the app behind demo-only fake data.",
	];

	const safetyGuide = [
		"Guests cannot target arbitrary external domains.",
		"Authenticated requests still fail closed on blocked destinations.",
		"Private networks, metadata IPs, and unsupported protocols are blocked.",
		"Request and response sizes are capped so the UI stays usable.",
		"Redirects are limited and revalidated before following them.",
	];
</script>

<svelte:head>
	<title>Docs - API Testing Kit quick start</title>
	<meta
		name="description"
		content="Quick-start docs for API Testing Kit covering guest mode, authenticated execution, templates, response inspection, and outbound safety rules."
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
							<p class="text-xs text-text-muted">Quick start for the shared guest and authenticated workspace</p>
						</div>
					</div>

					<div class="flex flex-wrap items-center gap-2 text-sm">
						<a href="#quick-start" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Quick start</a>
						<a href="#modes" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Modes</a>
						<a href="#request-builder" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Requests</a>
						<a href="#safety" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Safety</a>
						<a href="/app" class="rounded-full bg-primary-green px-5 py-2.5 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-primary-green-hover">Open app</a>
					</div>
				</div>
			</div>

			<div class="grid gap-8 px-5 py-6 lg:grid-cols-[minmax(0,1fr)_320px] lg:px-8 lg:py-8">
				<main class="space-y-10">
					<section id="quick-start" class="grid gap-6 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,0.9fr)]">
						<div class="space-y-5">
							<Badge variant="secondary" class="w-fit bg-primary-green-soft text-primary-green-deep">Quick start guide</Badge>
							<div class="space-y-4">
								<h1 class="max-w-3xl text-4xl font-semibold tracking-tight sm:text-5xl">
									Use the docs to get productive fast, not to hide the product behind fluff.
								</h1>
								<p class="max-w-2xl text-sm leading-7 text-text-body sm:text-base">
									API Testing Kit is a shared `/app` workspace with two capability levels. Guests explore curated templates and locked surfaces. Signed-in users unlock custom request execution, while the backend still enforces safety checks on every outbound call.
								</p>
							</div>

							<div class="flex flex-wrap gap-3">
								<Button href="/app" size="lg" class="rounded-full bg-primary-green px-6 text-white hover:bg-primary-green-hover">
									<PlayIcon class="size-4" />
									Open `/app`
								</Button>
								<Button href="/templates" variant="outline" size="lg" class="rounded-full border-border bg-white px-6">
									<SparklesIcon class="size-4" />
									Browse templates
								</Button>
							</div>

							<div class="grid gap-3 sm:grid-cols-3">
								<div class="metric-card p-4">
									<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Workspace</p>
									<p class="mt-2 text-sm font-semibold">One route, two modes</p>
									<p class="mt-1 text-xs leading-5 text-text-muted">Guest and authenticated users share the same shell.</p>
								</div>
								<div class="metric-card p-4">
									<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Execution</p>
									<p class="mt-2 text-sm font-semibold">Validated server-side</p>
									<p class="mt-1 text-xs leading-5 text-text-muted">The server blocks unsafe targets and oversized payloads.</p>
								</div>
								<div class="metric-card p-4">
									<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Output</p>
									<p class="mt-2 text-sm font-semibold">Readable response state</p>
									<p class="mt-1 text-xs leading-5 text-text-muted">Inspect headers, timing, size, and body together.</p>
								</div>
							</div>
						</div>

						<Card class="panel-card overflow-hidden">
							<CardHeader class="gap-3 border-b border-[#ece6db] bg-surface-soft/70">
								<div class="flex items-center justify-between gap-3">
									<div>
										<CardTitle>Fast path</CardTitle>
										<CardDescription>Use this sequence the first time you open the app.</CardDescription>
									</div>
									<Badge variant="outline" class="border-primary-green-soft bg-white text-primary-green-deep">4 steps</Badge>
								</div>
							</CardHeader>
							<CardContent class="space-y-4 p-5">
								{#each quickStart as item}
									<div class="flex gap-4 rounded-[20px] border border-border/70 bg-white p-4">
										<div class="grid h-10 w-10 shrink-0 place-items-center rounded-full bg-primary-green-soft text-sm font-semibold text-primary-green-deep">
											{item.step}
										</div>
										<div class="space-y-1">
											<p class="text-sm font-semibold">{item.title}</p>
											<p class="text-sm leading-6 text-text-body">{item.description}</p>
										</div>
									</div>
								{/each}
							</CardContent>
						</Card>
					</section>

					<section id="modes" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">App modes</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">Guest mode stays real, authenticated mode unlocks execution</h2>
							<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
								The workspace is intentionally split by capability, not by route. Guests see the actual interface with locked actions. Signed-in users get custom request execution and persistence, but not an unsafe proxy.
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

					<section id="request-builder" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">Request builder</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">Keep the full request visible while you edit it</h2>
							<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
								The editor is built around the same request shape documented in the plan: method, URL, params, headers, auth, and body. That keeps the workflow fast enough for debugging and precise enough for demos.
							</p>
						</div>

						<div class="grid gap-4 md:grid-cols-2">
							{#each requestGuide as item}
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
							<p class="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-text-muted">Example request shape</p>
							<pre class="overflow-x-auto text-xs leading-6"><code>GET https://jsonplaceholder.typicode.com/posts/1
Accept: application/json
X-Demo-Mode: guest</code></pre>
						</div>
					</section>

					<section id="responses" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">Response inspection</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">Read the result before you decide what to do next</h2>
							<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
								The response viewer is meant to answer the practical questions quickly: did the call succeed, how long did it take, how much data came back, and what exactly did the server return.
							</p>
						</div>

						<div class="grid gap-4 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]">
							<Card class="panel-card">
								<CardHeader class="gap-2">
									<CardTitle class="text-base">What to read first</CardTitle>
									<CardDescription>Use metadata before diving into the body.</CardDescription>
								</CardHeader>
								<CardContent class="space-y-3">
									{#each responseGuide as item}
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
								<CardDescription>Templates are the fastest safe entry point into `/app`.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								{#each templateGuide as item}
									<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
										<Code2Icon class="mt-0.5 size-4 text-primary-green" />
										<p class="text-sm leading-6 text-text-body">{item}</p>
									</div>
								{/each}
							</CardContent>
						</Card>

						<Card class="panel-card">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">What guests can do</CardTitle>
								<CardDescription>Guests should understand the product, not bounce off a fake demo.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
									<ShieldCheckIcon class="mt-0.5 size-4 text-success" />
									<p class="text-sm leading-6 text-text-body">Browse templates and example collections.</p>
								</div>
								<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
									<ShieldCheckIcon class="mt-0.5 size-4 text-success" />
									<p class="text-sm leading-6 text-text-body">Edit only the template-defined safe overrides.</p>
								</div>
								<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
									<ShieldCheckIcon class="mt-0.5 size-4 text-success" />
									<p class="text-sm leading-6 text-text-body">Run only allowlisted endpoints with visible lock states for everything else.</p>
								</div>
							</CardContent>
						</Card>
					</section>

					<section id="safety" class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
						<Card class="panel-card">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">Safety rules</CardTitle>
								<CardDescription>The backend is designed to fail closed.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								{#each safetyGuide as item}
									<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
										<LockIcon class="mt-0.5 size-4 text-primary-green" />
										<p class="text-sm leading-6 text-text-body">{item}</p>
									</div>
								{/each}
							</CardContent>
						</Card>

						<Card class="panel-card bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.98))]">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">If something is locked</CardTitle>
								<CardDescription>Visible, deliberate, and explained.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								<div class="rounded-[18px] border border-border/70 bg-white p-4">
									<div class="flex items-start gap-3">
										<div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-primary-green-soft text-primary-green-deep">
											<LockIcon class="size-4" />
										</div>
										<div>
											<p class="text-sm font-semibold">The control stays on screen.</p>
											<p class="mt-1 text-sm leading-6 text-text-body">
												The UI should show what exists, why it is blocked, and what signing in changes. Hidden controls make the product feel unfinished.
											</p>
										</div>
									</div>
								</div>

								<div class="flex flex-wrap gap-3">
									<Button href="/app" size="sm" class="rounded-full bg-primary-green px-4 text-white hover:bg-primary-green-hover">
										<ArrowRightIcon class="size-4" />
										Continue in `/app`
									</Button>
									<Button href="/features" variant="outline" size="sm" class="rounded-full border-border bg-white px-4">
										<SearchIcon class="size-4" />
										Review features
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
							<CardDescription>Jump to the section you need.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-2">
							<a href="#quick-start" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Quick start</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#modes" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>App modes</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#request-builder" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Request builder</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#responses" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Responses</span>
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
							<CardTitle class="text-base">Reading order</CardTitle>
							<CardDescription>Use this as a concise tour through the product.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-3 text-sm leading-6 text-text-body">
							<p>1. Open `/app` and pick a template.</p>
							<p>2. Edit the request with the fields the mode allows.</p>
							<p>3. Inspect the response metadata and body.</p>
							<p>4. Sign in when you need persistence or custom execution.</p>
						</CardContent>
					</Card>

					<Card class="panel-card bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.98))]">
						<CardHeader class="gap-2">
							<CardTitle class="text-base">Need the live app?</CardTitle>
							<CardDescription>Docs explain the surface, but the product lives in `/app`.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-3">
							<p class="text-sm leading-6 text-text-body">
								Use the docs to orient yourself, then move into the shared workspace. That is where guest mode, authenticated mode, templates, and response inspection all come together.
							</p>
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
