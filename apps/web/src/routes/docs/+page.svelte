<script lang="ts">
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
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
			title: "Open /app",
			description:
				'Guests can start immediately on the shared workspace. Open a curated template or a sample collection and you are already looking at the real product.'
		},
		{
			step: "2",
			title: "Choose a request shape",
			description:
				'Pick a method, review the target URL, and inspect the headers, auth, and body tabs. The workspace keeps the request editor visible instead of hiding it behind dialogs.'
		},
		{
			step: "3",
			title: "Send a safe request",
			description:
				'Guests can only use allowlisted demo targets. Signed-in users unlock validated custom URLs, but private ranges, metadata IPs, and unsupported protocols remain blocked.'
		},
		{
			step: "4",
			title: "Read the response",
			description:
				'Check status, timing, size, headers, and the pretty/raw tabs together so the result is immediately legible.'
		}
	];

	const requestGuide = [
		{
			title: "Method",
			text: "Start with GET for inspection, then move to POST, PUT, PATCH, or DELETE when you need to validate write paths.",
			tone: "bg-primary-green-soft text-primary-green-deep"
		},
		{
			title: "URL",
			text: "Guests stay on curated endpoints. Authenticated users can type custom URLs, but the backend validates every target before sending.",
			tone: "bg-white text-text-strong"
		},
		{
			title: "Params and headers",
			text: "Keep query parameters and headers explicit. That makes retries, debugging, and code snippet generation much easier.",
			tone: "bg-white text-text-strong"
		},
		{
			title: "Body and auth",
			text: "Use JSON, raw text, or form-urlencoded bodies. Bearer tokens and other custom auth flows unlock after sign-in.",
			tone: "bg-white text-text-strong"
		}
	];

	const responseGuide = [
		"Status badge and success or failure signal",
		"Response time and payload size",
		"Content type and header inspection",
		"Pretty, raw, and headers views"
	];

	const collectionFlow = [
		"Save a request once it is useful.",
		"Group related requests into a collection.",
		"Re-run from history or collection detail later.",
		"Keep reusable examples in the same workspace language."
	];

	const guestCapabilities = [
		"Browse templates and sample collections",
		"Edit only the safe fields exposed by the template",
		"Run only allowlisted demo endpoints",
		"See the full shell with locked actions instead of fake placeholders"
	];

	const signedInUnlocks = [
		"Custom target URLs with backend validation",
		"Saved requests and durable history",
		"Collections with re-run and organize flows",
		"Auth headers, body editors, and richer execution controls"
	];
</script>

<svelte:head>
	<title>Docs - API Testing Kit quick start</title>
	<meta
		name="description"
		content="Quick-start documentation for API Testing Kit with request building, response inspection, collections, and the guest versus signed-in model."
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
							<p class="text-xs text-text-muted">Quick start guidance for guests and signed-in users</p>
						</div>
					</div>

					<div class="flex flex-wrap items-center gap-2 text-sm">
						<a href="#quick-start" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Quick start</a>
						<a href="#request-builder" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Requests</a>
						<a href="#response-reading" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Responses</a>
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
									Learn the product in a few minutes, then move straight into the workspace.
								</h1>
								<p class="max-w-2xl text-sm leading-7 text-text-body sm:text-base">
									This page is the shortest path from first visit to a useful request. It explains how to open the app, shape a request, read the response, and understand what stays locked in guest mode.
								</p>
							</div>

							<div class="flex flex-wrap gap-3">
								<Button href="/app" size="lg" class="rounded-full bg-primary-green px-6 text-white hover:bg-primary-green-hover">
									<PlayIcon class="size-4" />
									Open /app
								</Button>
								<Button href="/features" variant="outline" size="lg" class="rounded-full border-border bg-white px-6">
									<SparklesIcon class="size-4" />
									Review features
								</Button>
							</div>

							<div class="grid gap-3 sm:grid-cols-3">
								<div class="metric-card p-4">
									<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Focus</p>
									<p class="mt-2 text-sm font-semibold">Request first</p>
									<p class="mt-1 text-xs leading-5 text-text-muted">Build and inspect before you save anything.</p>
								</div>
								<div class="metric-card p-4">
									<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Mode</p>
									<p class="mt-2 text-sm font-semibold">Guest or signed in</p>
									<p class="mt-1 text-xs leading-5 text-text-muted">Same workspace, different capability levels.</p>
								</div>
								<div class="metric-card p-4">
									<p class="text-xs uppercase tracking-[0.24em] text-text-muted">Safety</p>
									<p class="mt-2 text-sm font-semibold">Validated execution</p>
									<p class="mt-1 text-xs leading-5 text-text-muted">The backend blocks dangerous outbound targets.</p>
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
									<Badge variant="outline" class="border-primary-green-soft bg-white text-primary-green-deep">3 min</Badge>
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

					<section id="request-builder" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">Request builder</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">Shape the request before you hit send</h2>
							<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
								The editor keeps the important pieces visible together: method, URL, params, headers, auth, and body. That makes the workflow easy to explain and easy to repeat.
							</p>
						</div>

						<div class="grid gap-4 md:grid-cols-2">
							{#each requestGuide as item}
								<Card class="panel-card">
									<CardHeader class="gap-2">
										<CardTitle class="text-base">{item.title}</CardTitle>
									</CardHeader>
									<CardContent>
										<div class={`rounded-[18px] border border-border/70 px-4 py-4 text-sm leading-6 ${item.tone}`}>
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

					<section id="response-reading" class="scroll-mt-24 space-y-4">
						<div class="max-w-3xl">
							<p class="section-title">Response inspection</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight sm:text-3xl">Read status, size, and body together</h2>
							<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
								The response area should answer the only questions that matter at a glance: did it work, how long did it take, how much data came back, and what was actually returned.
							</p>
						</div>

						<div class="grid gap-4 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]">
							<Card class="panel-card">
								<CardHeader class="gap-2">
									<CardTitle class="text-base">What to look for</CardTitle>
									<CardDescription>Use the metadata first, then inspect the body.</CardDescription>
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

					<section id="collections" class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
						<Card class="panel-card">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">Collections</CardTitle>
								<CardDescription>Use collections to keep related requests grouped and re-runnable.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								{#each collectionFlow as item}
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
								<CardDescription>The guest experience is real, but intentionally constrained.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								{#each guestCapabilities as item}
									<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
										<ShieldCheckIcon class="mt-0.5 size-4 text-success" />
										<p class="text-sm leading-6 text-text-body">{item}</p>
									</div>
								{/each}
							</CardContent>
						</Card>
					</section>

					<section class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
						<Card class="panel-card">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">What unlocks after sign-in</CardTitle>
								<CardDescription>Authenticated mode keeps the same shell and removes the artificial limits.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								{#each signedInUnlocks as item}
									<div class="flex items-start gap-3 rounded-[18px] border border-border/70 bg-white px-4 py-3">
										<LockIcon class="mt-0.5 size-4 text-primary-green" />
										<p class="text-sm leading-6 text-text-body">{item}</p>
									</div>
								{/each}
							</CardContent>
						</Card>

						<Card class="panel-card bg-[linear-gradient(135deg,rgba(31,122,77,0.14),rgba(255,255,255,0.98))]">
							<CardHeader class="gap-2">
								<CardTitle class="text-base">If something is locked</CardTitle>
								<CardDescription>Locked controls should be visible, not hidden.</CardDescription>
							</CardHeader>
							<CardContent class="space-y-3">
								<div class="rounded-[18px] border border-border/70 bg-white p-4">
									<div class="flex items-start gap-3">
										<div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-primary-green-soft text-primary-green-deep">
											<LockIcon class="size-4" />
										</div>
										<div>
											<p class="text-sm font-semibold">The workspace stays visible.</p>
											<p class="mt-1 text-sm leading-6 text-text-body">
												Gated features should feel deliberate: the user can see what exists, why it is blocked, and what signing in changes.
											</p>
										</div>
									</div>
								</div>

								<div class="flex flex-wrap gap-3">
									<Button href="/app" size="sm" class="rounded-full bg-primary-green px-4 text-white hover:bg-primary-green-hover">
										<ArrowRightIcon class="size-4" />
										Continue in /app
									</Button>
									<Button href="/" variant="outline" size="sm" class="rounded-full border-border bg-white px-4">
										<SearchIcon class="size-4" />
										Review landing page
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
							<CardDescription>Jump straight to the section you need.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-2">
							<a href="#quick-start" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Quick start</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#request-builder" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Request builder</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#response-reading" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Response reading</span>
								<ArrowRightIcon class="size-4" />
							</a>
							<a href="#collections" class="flex items-center justify-between rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm text-text-body transition hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong">
								<span>Collections</span>
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
							<p>1. Open the app and choose a template.</p>
							<p>2. Edit only the safe fields you need.</p>
							<p>3. Send the request and inspect the response.</p>
							<p>4. Save useful work into a collection once you sign in.</p>
						</CardContent>
					</Card>

					<Card class="panel-card bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.98))]">
						<CardHeader class="gap-2">
							<CardTitle class="text-base">Need the fastest path?</CardTitle>
							<CardDescription>Go straight to the live workspace.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-3">
							<p class="text-sm leading-6 text-text-body">
								The docs are here to reduce friction, but the product experience lives in the shared `/app` route.
							</p>
							<Button href="/app" class="w-full justify-between rounded-full bg-primary-green px-5 text-white hover:bg-primary-green-hover">
								<span>Open /app</span>
								<ArrowRightIcon class="size-4" />
							</Button>
						</CardContent>
					</Card>
				</aside>
			</div>
		</div>
	</div>
</div>
