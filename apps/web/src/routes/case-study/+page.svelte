<script lang="ts">
	const goals = [
		{
			title: 'Make the app understandable in seconds',
			body: 'The landing surface has to explain the request-response workflow before a visitor ever signs in or reads docs.'
		},
		{
			title: 'Keep guest mode real but constrained',
			body: 'Guests should experience the product surface, while outbound execution stays limited to allowlisted demo flows.'
		},
		{
			title: 'Prevent open-proxy behavior by design',
			body: 'The backend has to validate destination safety, request size, and response handling before anything leaves the system.'
		}
	];

	const frontendBlocks = [
		{
			title: 'Marketing shell and app shell share one visual language',
			body: 'The public pages and the `/app` workspace both use the same warm shell, rounded cards, and green action system so the product feels cohesive.',
			meta: 'SvelteKit routes, shared layout tokens, reusable UI primitives'
		},
		{
			title: 'Route structure makes the product model explicit',
			body: 'The landing page, docs, case study, and workspace are separate pages but read as one product instead of disconnected demos.',
			meta: 'Public pages on `/`, `/docs`, `/features`, `/case-study`, `/app`'
		},
		{
			title: 'Static first, backend-aware second',
			body: 'This page is static, but the content references the real product boundaries so it can survive the move from scaffold to implementation.',
			meta: 'No runtime dependency on server APIs'
		}
	];

	const backendBlocks = [
		{
			title: 'Go HTTP API owns execution and safety',
			body: 'The backend is responsible for auth, sessions, request execution, abuse logging, and the validation pipeline that prevents unsafe outbound calls.',
			meta: 'Server-side request runner, auth, abuse hooks'
		},
		{
			title: 'PostgreSQL stores product state',
			body: 'Collections, history, usage, and blocked target records belong in persistent storage so safety and auditability are first-class.',
			meta: 'Schema-driven persistence, migration-friendly'
		},
		{
			title: 'Guest limits are enforced server-side',
			body: 'The `/app` experience can show locked controls in the UI, but the actual restriction has to live in the backend so it cannot be bypassed.',
			meta: 'Guest-safe endpoints, strict execution rules'
		}
	];

	const pipeline = [
		'The user prepares a request in the shared workspace.',
		'The frontend sends a sanitized payload to the Go API.',
		'The backend checks identity, quotas, destination safety, and body limits.',
		'The request runner resolves the destination, validates redirects, and applies protocol restrictions.',
		'The response is capped, previewed, and normalized into structured metadata.',
		'The UI renders headers, pretty output, raw output, and response timing.'
	];

	const guestRules = [
		'Only allowlisted demo endpoints can be executed from guest mode.',
		'Guests can inspect the full workspace, but custom target editing stays locked.',
		'History, saved collections, environment variables, and unrestricted outbound URLs remain gated.',
		'The UI should show the locked controls instead of hiding them, so the product model is obvious.'
	];

	const safetyLayers = [
		'Protocol validation blocks unsupported schemes before a request is issued.',
		'Destination validation blocks localhost, private ranges, metadata IPs, and other unsafe targets.',
		'DNS resolution and redirect hops are re-checked instead of assumed safe.',
		'Timeouts, body size, response preview size, and redirect depth are capped.',
		'Suspicious attempts are logged for abuse review instead of silently ignored.'
	];

	const deploymentSteps = [
		'Local development uses the SvelteKit frontend and Go API as separate processes for clarity.',
		'Docker and Compose provide a reproducible local stack once persistence and migrations are wired in.',
		'Goose migrations define the database state so schema changes stay reviewable and replayable.',
		'CI should run lint, typecheck, tests, and build checks before merge so unsafe regressions do not sneak in.'
	];
</script>

<svelte:head>
	<title>API Testing Kit - Case study</title>
	<meta
		name="description"
		content="A technical case study covering the architecture, execution pipeline, guest gating, safety controls, and deployment strategy for API Testing Kit."
	/>
</svelte:head>

<div class="relative isolate overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(31,122,77,0.15),_transparent_30%),linear-gradient(180deg,_#f4f1ea_0%,_#ece7dd_100%)] text-[#162117]">
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-28 top-16 h-72 w-72 rounded-full bg-[#1f7a4d]/10 blur-3xl"></div>
		<div class="absolute right-[-5rem] top-1/3 h-96 w-96 rounded-full bg-[#dcefe3] blur-3xl"></div>
		<div class="absolute bottom-[-6rem] left-1/2 h-72 w-72 -translate-x-1/2 rounded-full bg-[#145336]/10 blur-3xl"></div>
	</div>

	<div class="mx-auto min-h-screen max-w-[1440px] px-4 py-4 sm:px-6 lg:px-8">
		<div class="overflow-hidden rounded-[32px] border border-[#e7e3d8] bg-[rgba(247,245,240,0.94)] shadow-[0_24px_60px_rgba(21,31,23,0.08)] backdrop-blur">
			<header class="border-b border-[#e7e3d8] bg-white/75 px-5 py-4 sm:px-6 lg:px-8">
				<div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
					<div class="flex items-center gap-3">
						<div class="grid h-11 w-11 place-items-center rounded-2xl bg-[#1f7a4d] text-sm font-semibold text-white shadow-[0_10px_24px_rgba(31,122,77,0.28)]">
							AT
						</div>
						<div>
							<p class="text-sm font-semibold tracking-tight text-[#162117]">API Testing Kit</p>
							<p class="text-xs text-[#7a847d]">Case study and engineering narrative</p>
						</div>
					</div>

					<nav class="flex flex-wrap items-center gap-2 text-sm">
						<a href="#problem" class="rounded-full px-4 py-2 text-[#445046] transition hover:bg-[#f5f3ed] hover:text-[#162117]">Problem</a>
						<a href="#architecture" class="rounded-full px-4 py-2 text-[#445046] transition hover:bg-[#f5f3ed] hover:text-[#162117]">Architecture</a>
						<a href="#safety" class="rounded-full px-4 py-2 text-[#445046] transition hover:bg-[#f5f3ed] hover:text-[#162117]">Safety</a>
						<a
							href="/app"
							class="rounded-full bg-[#1f7a4d] px-5 py-2.5 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-[#19663f]"
						>
							Open app
						</a>
					</nav>
				</div>
			</header>

			<main class="space-y-12 px-5 py-6 sm:px-6 lg:px-8 lg:py-8">
				<section class="grid gap-6 lg:grid-cols-[1.15fr_0.85fr] lg:items-start">
					<div class="space-y-6">
						<div class="inline-flex items-center gap-2 rounded-full border border-[#d9e7d8] bg-white/80 px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-[#1f7a4d]">
							<span class="h-2 w-2 rounded-full bg-[#1f7a4d]"></span>
							Engineering narrative
						</div>

						<div class="space-y-4">
							<h1 class="max-w-3xl text-4xl font-semibold tracking-tight text-[#162117] sm:text-5xl lg:text-6xl">
								How the product stays useful for guests without becoming an open proxy
							</h1>
							<p class="max-w-2xl text-sm leading-7 text-[#445046] sm:text-base">
								API Testing Kit is built around a narrow but credible idea: let visitors experience a real request-response workspace, then move signed-in users into a safer, more flexible execution model without changing the core interface.
							</p>
						</div>

						<div class="grid gap-3 sm:grid-cols-3">
							<div class="rounded-2xl border border-[#e7e3d8] bg-white/80 p-4 shadow-[0_10px_24px_rgba(21,31,23,0.04)]">
								<p class="text-xs font-medium uppercase tracking-[0.22em] text-[#7a847d]">Mode</p>
								<p class="mt-2 text-base font-semibold text-[#162117]">Guest + authenticated</p>
							</div>
							<div class="rounded-2xl border border-[#e7e3d8] bg-white/80 p-4 shadow-[0_10px_24px_rgba(21,31,23,0.04)]">
								<p class="text-xs font-medium uppercase tracking-[0.22em] text-[#7a847d]">Execution</p>
								<p class="mt-2 text-base font-semibold text-[#162117]">Server validated</p>
							</div>
							<div class="rounded-2xl border border-[#e7e3d8] bg-white/80 p-4 shadow-[0_10px_24px_rgba(21,31,23,0.04)]">
								<p class="text-xs font-medium uppercase tracking-[0.22em] text-[#7a847d]">Storage</p>
								<p class="mt-2 text-base font-semibold text-[#162117]">PostgreSQL-backed</p>
							</div>
						</div>
					</div>

					<aside class="rounded-[28px] border border-[#dfe8dd] bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.96))] p-6 shadow-[0_18px_40px_rgba(21,31,23,0.08)]">
						<p class="text-xs font-semibold uppercase tracking-[0.28em] text-[#1f7a4d]">Why this exists</p>
						<div class="mt-4 space-y-4">
							<p class="text-sm leading-6 text-[#445046]">
								The project needs to demonstrate frontend craft, backend judgment, and security awareness in a single product. A case study page is where those decisions become explicit.
							</p>
							<div class="rounded-[22px] border border-[#e7e3d8] bg-white/85 p-4">
								<p class="text-sm font-semibold text-[#162117]">Source of truth</p>
								<p class="mt-2 text-xs leading-5 text-[#7a847d]">
									The narrative follows the execution plan, UI map, and design system docs so it stays aligned with the repo instead of drifting into generic marketing copy.
								</p>
							</div>
						</div>
					</aside>
				</section>

				<section id="problem" class="scroll-mt-24 space-y-6">
					<div class="max-w-3xl">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Problem statement</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117] sm:text-3xl">
							The product has to feel demo-ready, but its execution path still has to be defensible
						</h2>
						<p class="mt-3 text-sm leading-6 text-[#445046] sm:text-base">
							That creates a tension most request tools ignore. If everything is wide open, the site becomes a proxy. If everything is locked down, the product feels fake. The design here keeps both sides visible and moves the real restrictions into backend enforcement.
						</p>
					</div>

					<div class="grid gap-4 lg:grid-cols-3">
						{#each goals as goal}
							<article class="rounded-[24px] border border-[#e7e3d8] bg-white p-6 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
								<p class="text-sm font-semibold text-[#162117]">{goal.title}</p>
								<p class="mt-3 text-sm leading-6 text-[#445046]">{goal.body}</p>
							</article>
						{/each}
					</div>
				</section>

				<section id="architecture" class="scroll-mt-24 space-y-6">
					<div class="max-w-3xl">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Architecture</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117] sm:text-3xl">
							Frontend and backend are split by responsibility, not just by framework
						</h2>
					</div>

					<div class="grid gap-4 xl:grid-cols-2">
						<article class="rounded-[28px] border border-[#e7e3d8] bg-white p-6 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
							<div class="flex items-center justify-between gap-3">
								<div>
									<p class="text-xs font-semibold uppercase tracking-[0.28em] text-[#1f7a4d]">Frontend</p>
									<h3 class="mt-2 text-xl font-semibold tracking-tight text-[#162117]">SvelteKit owns presentation and product flow</h3>
								</div>
								<span class="rounded-full border border-[#dfe8dd] bg-[#f5f9f4] px-3 py-1 text-xs font-medium text-[#1f7a4d]">Public-facing</span>
							</div>

							<div class="mt-5 grid gap-3">
								{#each frontendBlocks as block}
									<div class="rounded-[22px] border border-[#e7e3d8] bg-[#f9f8f4] p-4">
										<p class="text-sm font-semibold text-[#162117]">{block.title}</p>
										<p class="mt-2 text-sm leading-6 text-[#445046]">{block.body}</p>
										<p class="mt-3 text-xs font-mono text-[#7a847d]">{block.meta}</p>
									</div>
								{/each}
							</div>
						</article>

						<article class="rounded-[28px] border border-[#e7e3d8] bg-white p-6 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
							<div class="flex items-center justify-between gap-3">
								<div>
									<p class="text-xs font-semibold uppercase tracking-[0.28em] text-[#1f7a4d]">Backend</p>
									<h3 class="mt-2 text-xl font-semibold tracking-tight text-[#162117]">Go owns safety, sessions, and execution</h3>
								</div>
								<span class="rounded-full border border-[#dfe8dd] bg-[#f5f9f4] px-3 py-1 text-xs font-medium text-[#1f7a4d]">Control plane</span>
							</div>

							<div class="mt-5 grid gap-3">
								{#each backendBlocks as block}
									<div class="rounded-[22px] border border-[#e7e3d8] bg-[#f9f8f4] p-4">
										<p class="text-sm font-semibold text-[#162117]">{block.title}</p>
										<p class="mt-2 text-sm leading-6 text-[#445046]">{block.body}</p>
										<p class="mt-3 text-xs font-mono text-[#7a847d]">{block.meta}</p>
									</div>
								{/each}
							</div>
						</article>
					</div>
				</section>

				<section class="grid gap-4 lg:grid-cols-[1fr_0.92fr]">
					<article class="rounded-[28px] border border-[#e7e3d8] bg-white p-6 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Request pipeline</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117]">Execution is intentionally staged</h2>
						<div class="mt-5 space-y-3">
							{#each pipeline as step, index}
								<div class="flex gap-4 rounded-[22px] border border-[#e7e3d8] bg-[#f9f8f4] p-4">
									<div class="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-[#dcefe3] text-sm font-semibold text-[#1f7a4d]">
										{index + 1}
									</div>
									<p class="text-sm leading-6 text-[#445046]">{step}</p>
								</div>
							{/each}
						</div>
					</article>

					<article class="rounded-[28px] border border-[#e7e3d8] bg-[linear-gradient(180deg,rgba(255,255,255,0.96),rgba(245,241,235,0.92))] p-6 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Guest gating</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117]">Visible controls, server-enforced limits</h2>
						<div class="mt-5 space-y-3">
							{#each guestRules as rule}
								<div class="rounded-[22px] border border-[#e7e3d8] bg-white p-4">
									<p class="text-sm leading-6 text-[#445046]">{rule}</p>
								</div>
							{/each}
						</div>
					</article>
				</section>

				<section id="safety" class="scroll-mt-24 grid gap-4 xl:grid-cols-2">
					<article class="rounded-[28px] border border-[#e7e3d8] bg-white p-6 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Abuse prevention</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117]">Safety is part of the product contract</h2>
						<div class="mt-5 space-y-3">
							{#each safetyLayers as layer}
								<div class="rounded-[22px] border border-[#e7e3d8] bg-[#f9f8f4] p-4">
									<p class="text-sm leading-6 text-[#445046]">{layer}</p>
								</div>
							{/each}
						</div>
					</article>

					<article class="rounded-[28px] border border-[#e7e3d8] bg-white p-6 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Deployment</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117]">The stack is simple enough to ship locally and production-ready enough to extend</h2>
						<div class="mt-5 space-y-3">
							{#each deploymentSteps as step}
								<div class="flex items-start gap-3 rounded-[22px] border border-[#e7e3d8] bg-[#f9f8f4] p-4">
									<span class="mt-1 h-2.5 w-2.5 rounded-full bg-[#1f7a4d]"></span>
									<p class="text-sm leading-6 text-[#445046]">{step}</p>
								</div>
							{/each}
						</div>
					</article>
				</section>

				<section class="rounded-[30px] border border-[#dfe8dd] bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.96))] p-6 shadow-[0_18px_40px_rgba(21,31,23,0.08)] sm:p-8">
					<div class="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
						<div class="max-w-2xl">
							<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Outcome</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117] sm:text-3xl">
								The case study should make the product feel deliberate, not improvised
							</h2>
							<p class="mt-3 text-sm leading-6 text-[#445046] sm:text-base">
								The point of the page is to show that the product is architected around the experience and the risk model at the same time. That is the difference between a marketing mockup and something you can actually deploy.
							</p>
						</div>

						<div class="flex flex-wrap gap-3">
							<a
								href="/app"
								class="rounded-full bg-[#1f7a4d] px-6 py-3 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-[#19663f]"
							>
								Open /app
							</a>
							<a
								href="/docs"
								class="rounded-full border border-[#d4ded1] bg-white px-6 py-3 font-medium text-[#162117] transition hover:bg-[#f5f3ed]"
							>
								Read the quick start
							</a>
						</div>
					</div>
				</section>
			</main>

			<footer class="border-t border-[#e7e3d8] bg-white/60 px-5 py-5 sm:px-6 lg:px-8">
				<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
					<div>
						<p class="text-sm font-semibold text-[#162117]">API Testing Kit</p>
						<p class="mt-1 text-xs text-[#7a847d]">Case study, architecture notes, and delivery context.</p>
					</div>

					<div class="flex flex-wrap items-center gap-3 text-sm">
						<a href="#problem" class="rounded-full px-4 py-2 text-[#445046] transition hover:bg-[#f5f3ed] hover:text-[#162117]">Problem</a>
						<a href="#architecture" class="rounded-full px-4 py-2 text-[#445046] transition hover:bg-[#f5f3ed] hover:text-[#162117]">Architecture</a>
						<a href="#safety" class="rounded-full px-4 py-2 text-[#445046] transition hover:bg-[#f5f3ed] hover:text-[#162117]">Safety</a>
					</div>
				</div>
			</footer>
		</div>
	</div>
</div>
