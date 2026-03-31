<script lang="ts">
	import {
		ArrowRight,
		BadgeCheck,
		Code2,
		History,
		Lock,
		PanelTop,
		ShieldCheck,
		SquareStack,
		Workflow
	} from '@lucide/svelte';

	const heroStats = [
		{ label: 'Workspace mode', value: 'Guest + signed in' },
		{ label: 'Safety posture', value: 'Outbound controlled' },
		{ label: 'Primary focus', value: 'Request + response flow' }
	];

	const capabilityCards = [
		{
			icon: Workflow,
			title: 'Request builder',
			description:
				'Compose methods, URLs, params, headers, auth, and body modes inside one visible editing surface.'
		},
		{
			icon: PanelTop,
			title: 'Response inspection',
			description:
				'Inspect pretty, raw, and header views alongside status, timing, payload size, and content type.'
		},
		{
			icon: SquareStack,
			title: 'Collections and history',
			description:
				'Group reusable requests, revisit recent runs, and keep a durable trail of useful API interactions.'
		},
		{
			icon: ShieldCheck,
			title: 'Safety controls',
			description:
				'Validate destinations, block private ranges, and keep the product from behaving like an open proxy.'
		}
	];

	const runtimeCards = [
		{
			label: 'Request',
			title: 'Request state',
			text: 'Method, URL, headers, auth, and body remain visible on one surface.'
		},
		{
			label: 'Execution',
			title: 'Validated runner',
			text: 'The backend enforces destination checks, size limits, and session rules before any outbound call.'
		},
		{
			label: 'Response',
			title: 'Response state',
			text: 'Status, timing, size, headers, and body views stay attached to the request context.'
		}
	];

	const modeRows = [
		{
			label: 'Request targets',
			guest: 'Curated examples only',
			authenticated: 'Custom validated URLs'
		},
		{
			label: 'Persistence',
			guest: 'No saved history',
			authenticated: 'Collections and request history'
		},
		{
			label: 'Secrets',
			guest: 'No environment variables',
			authenticated: 'Account-level variables and auth'
		},
		{
			label: 'Execution limits',
			guest: 'Tighter size and timeout caps',
			authenticated: 'Higher quotas with abuse protection'
		}
	];

	const safetyItems = [
		'Block localhost, RFC1918 private ranges, link-local targets, and metadata IPs.',
		'Re-check DNS and redirect hops before the request is allowed to proceed.',
		'Limit body size, response preview size, redirect depth, and request duration.',
		'Log suspicious traffic patterns so abuse review has an audit trail.'
	];

	const responsePreview = String.raw`{
  "status": "ready",
  "mode": "pretty",
  "trust": "validated"
}`;
</script>

<svelte:head>
	<title>Features - API Testing Kit</title>
	<meta
		name="description"
		content="Reference view for the request builder, response viewer, collections, history, guest limits, and safety controls in API Testing Kit."
	/>
</svelte:head>

<div class="relative isolate overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(31,122,77,0.14),transparent_30%),linear-gradient(180deg,#f2f0ea_0%,#ece8df_100%)] text-[#162117]">
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-32 top-16 h-80 w-80 rounded-full bg-[#1f7a4d]/10 blur-3xl"></div>
		<div class="absolute right-[-5rem] top-1/4 h-96 w-96 rounded-full bg-[#dcefe3] blur-3xl"></div>
		<div class="absolute bottom-[-7rem] left-1/2 h-72 w-72 -translate-x-1/2 rounded-full bg-[#145336]/10 blur-3xl"></div>
	</div>

	<div class="mx-auto min-h-screen max-w-[1440px] px-4 py-4 sm:px-6 lg:px-8">
		<div class="overflow-hidden rounded-[32px] border border-[#e7e3d8] bg-[rgba(247,245,240,0.92)] shadow-[0_24px_60px_rgba(21,31,23,0.08)] backdrop-blur">
			<header class="border-b border-[#e7e3d8] bg-white/75 px-5 py-4 sm:px-6 lg:px-8">
				<div class="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
					<div class="max-w-3xl">
						<div class="inline-flex items-center gap-2 rounded-full border border-[#d9e7d8] bg-white/80 px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-[#1f7a4d]">
							<span class="h-2 w-2 rounded-full bg-[#1f7a4d]"></span>
							Product features
						</div>
						<h1 class="mt-4 text-4xl font-semibold tracking-tight text-[#162117] sm:text-5xl lg:text-6xl">
							Feature reference for the shared API workspace
						</h1>
						<p class="mt-4 max-w-2xl text-sm leading-7 text-[#445046] sm:text-base">
							Request editing, response inspection, persistence, and safety controls share one product surface.
						</p>
					</div>

					<div class="flex flex-wrap gap-3 lg:justify-end">
						<a
							href="/app"
							class="inline-flex items-center gap-2 rounded-full bg-[#1f7a4d] px-6 py-3 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-[#19663f]"
						>
							Open app
							<ArrowRight size={16} />
						</a>
						<a
							href="/docs"
							class="rounded-full border border-[#d4ded1] bg-white px-6 py-3 font-medium text-[#162117] transition hover:bg-[#f5f3ed]"
						>
							Docs
						</a>
					</div>
				</div>

				<div class="mt-6 grid gap-3 sm:grid-cols-3">
					{#each heroStats as stat}
						<div class="rounded-2xl border border-[#e7e3d8] bg-white/80 p-4 shadow-[0_10px_24px_rgba(21,31,23,0.04)]">
							<p class="text-xs font-medium uppercase tracking-[0.22em] text-[#7a847d]">{stat.label}</p>
							<p class="mt-2 text-base font-semibold text-[#162117]">{stat.value}</p>
						</div>
					{/each}
				</div>
			</header>

			<main class="space-y-12 px-5 py-6 sm:px-6 lg:px-8 lg:py-8">
				<section class="space-y-6">
					<div class="max-w-3xl">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Core capabilities</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117] sm:text-3xl">
							Core workspace capabilities
						</h2>
						<p class="mt-3 text-sm leading-6 text-[#445046] sm:text-base">
							The workspace stays focused on request state, response state, persistence, and safety.
						</p>
					</div>

					<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
						{#each capabilityCards as card}
							{@const Icon = card.icon}
							<article class="rounded-[26px] border border-[#e7e3d8] bg-white/88 p-5 shadow-[0_12px_30px_rgba(21,31,23,0.05)] transition hover:-translate-y-0.5 hover:shadow-[0_16px_34px_rgba(21,31,23,0.08)]">
								<div class="flex items-start justify-between gap-4">
									<div class="grid h-11 w-11 place-items-center rounded-2xl bg-[#dcefe3] text-[#145336]">
										<Icon size={18} strokeWidth={2.1} />
									</div>
									<span class="rounded-full bg-[#f5f3ed] px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.24em] text-[#7a847d]">
										Feature
									</span>
								</div>
								<h3 class="mt-5 text-base font-semibold tracking-tight text-[#162117]">{card.title}</h3>
								<p class="mt-2 text-sm leading-6 text-[#445046]">{card.description}</p>
							</article>
						{/each}
					</div>
				</section>

				<section class="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
					<div class="rounded-[30px] border border-[#e7e3d8] bg-white/88 p-6 shadow-[0_12px_30px_rgba(21,31,23,0.05)] sm:p-7">
						<div class="max-w-2xl">
							<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Runtime surfaces</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117] sm:text-3xl">
								Request, execution, and response stay visible together
							</h2>
						</div>

						<div class="mt-6 grid gap-4 md:grid-cols-3">
							{#each runtimeCards as item}
								<div class="rounded-[24px] border border-[#e7e3d8] bg-[#f5f3ed] p-5">
									<p class="text-xs font-semibold uppercase tracking-[0.24em] text-[#1f7a4d]">{item.label}</p>
									<h3 class="mt-3 text-base font-semibold tracking-tight text-[#162117]">{item.title}</h3>
									<p class="mt-2 text-sm leading-6 text-[#445046]">{item.text}</p>
								</div>
							{/each}
						</div>
					</div>

					<div class="rounded-[30px] border border-[#e7e3d8] bg-[linear-gradient(135deg,rgba(31,122,77,0.12),rgba(255,255,255,0.96))] p-6 shadow-[0_12px_30px_rgba(21,31,23,0.05)] sm:p-7">
						<div class="flex items-center gap-3">
							<div class="grid h-11 w-11 place-items-center rounded-2xl bg-[#dcefe3] text-[#145336]">
								<Code2 size={18} strokeWidth={2.1} />
							</div>
							<div>
								<p class="text-sm font-semibold text-[#162117]">Response inspection</p>
								<p class="text-xs text-[#7a847d]">Pretty, raw, and headers views</p>
							</div>
						</div>

						<div class="mt-5 space-y-3 rounded-[24px] border border-[#e7e3d8] bg-white/92 p-4">
							<div class="flex items-center justify-between gap-3">
								<p class="text-xs font-semibold uppercase tracking-[0.22em] text-[#7a847d]">Metadata</p>
								<span class="rounded-full bg-[#dcefe3] px-3 py-1 text-xs font-semibold text-[#145336]">200 OK</span>
							</div>
							<div class="grid grid-cols-3 gap-2 text-xs text-[#445046]">
								<div class="rounded-2xl border border-[#e7e3d8] bg-[#f5f3ed] px-3 py-2">
									<p class="text-[11px] uppercase tracking-[0.2em] text-[#7a847d]">Time</p>
									<p class="mt-1 font-semibold text-[#162117]">186 ms</p>
								</div>
								<div class="rounded-2xl border border-[#e7e3d8] bg-[#f5f3ed] px-3 py-2">
									<p class="text-[11px] uppercase tracking-[0.2em] text-[#7a847d]">Size</p>
									<p class="mt-1 font-semibold text-[#162117]">1.2 KB</p>
								</div>
								<div class="rounded-2xl border border-[#e7e3d8] bg-[#f5f3ed] px-3 py-2">
									<p class="text-[11px] uppercase tracking-[0.2em] text-[#7a847d]">Type</p>
									<p class="mt-1 font-semibold text-[#162117]">json</p>
								</div>
							</div>
							<pre class="overflow-x-auto rounded-[22px] border border-[#e7e3d8] bg-[#162117] px-4 py-4 text-xs leading-6 text-[#d8ebdf]"><code>{responsePreview}</code></pre>
						</div>
					</div>
				</section>

				<section class="space-y-6">
					<div class="max-w-3xl">
						<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Capability split</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117] sm:text-3xl">
							Access changes by mode, not by route
						</h2>
					</div>

					<div class="overflow-hidden rounded-[30px] border border-[#e7e3d8] bg-white/90 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
						<div class="grid grid-cols-[1.1fr_1fr_1fr] border-b border-[#e7e3d8] bg-[#f5f3ed] px-5 py-4 text-xs font-semibold uppercase tracking-[0.22em] text-[#7a847d]">
							<span>Capability</span>
							<span>Guest</span>
							<span>Authenticated</span>
						</div>
						{#each modeRows as row}
							<div class="grid grid-cols-[1.1fr_1fr_1fr] gap-4 border-b border-[#eef0ea] px-5 py-4 last:border-b-0">
								<div class="font-medium text-[#162117]">{row.label}</div>
								<div class="text-sm text-[#445046]">{row.guest}</div>
								<div class="text-sm text-[#445046]">{row.authenticated}</div>
							</div>
						{/each}
					</div>
				</section>

				<section class="grid gap-4 lg:grid-cols-[1fr_0.95fr]">
					<div class="rounded-[30px] border border-[#e7e3d8] bg-white/88 p-6 shadow-[0_12px_30px_rgba(21,31,23,0.05)] sm:p-7">
						<div class="flex items-center gap-3">
							<div class="grid h-11 w-11 place-items-center rounded-2xl bg-[#dcefe3] text-[#145336]">
								<History size={18} strokeWidth={2.1} />
							</div>
							<div>
								<p class="text-sm font-semibold text-[#162117]">Collections and history</p>
								<p class="text-xs text-[#7a847d]">Saved work for signed-in users</p>
							</div>
						</div>
						<p class="mt-4 text-sm leading-6 text-[#445046]">
							Collections group reusable requests, while history preserves recent runs and response snapshots.
						</p>
						<div class="mt-5 grid gap-3 sm:grid-cols-2">
							<div class="rounded-[22px] border border-[#e7e3d8] bg-[#f5f3ed] p-4">
								<p class="text-xs font-semibold uppercase tracking-[0.22em] text-[#7a847d]">Collections</p>
								<p class="mt-2 text-sm leading-6 text-[#162117]">Reusable request groups with clear ownership and easy re-entry.</p>
							</div>
							<div class="rounded-[22px] border border-[#e7e3d8] bg-[#f5f3ed] p-4">
								<p class="text-xs font-semibold uppercase tracking-[0.22em] text-[#7a847d]">History</p>
								<p class="mt-2 text-sm leading-6 text-[#162117]">Durable execution records for review, reuse, and debugging.</p>
							</div>
						</div>
					</div>

					<div class="rounded-[30px] border border-[#e7e3d8] bg-[linear-gradient(135deg,rgba(31,122,77,0.08),rgba(255,255,255,0.96))] p-6 shadow-[0_12px_30px_rgba(21,31,23,0.05)] sm:p-7">
						<div class="flex items-center gap-3">
							<div class="grid h-11 w-11 place-items-center rounded-2xl bg-[#dcefe3] text-[#145336]">
								<Lock size={18} strokeWidth={2.1} />
							</div>
							<div>
								<p class="text-sm font-semibold text-[#162117]">Safety and abuse prevention</p>
								<p class="text-xs text-[#7a847d]">Security is part of the product surface</p>
							</div>
						</div>

						<ul class="mt-5 space-y-3">
							{#each safetyItems as item}
								<li class="flex gap-3 rounded-[22px] border border-[#e7e3d8] bg-white/92 p-4 text-sm leading-6 text-[#445046]">
									<BadgeCheck size={18} class="mt-0.5 shrink-0 text-[#1f7a4d]" strokeWidth={2.1} />
									<span>{item}</span>
								</li>
							{/each}
						</ul>
					</div>
				</section>

				<section class="rounded-[30px] border border-[#dfe8dd] bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.96))] p-6 shadow-[0_18px_40px_rgba(21,31,23,0.08)] sm:p-8">
					<div class="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
						<div class="max-w-2xl">
							<p class="text-xs font-semibold uppercase tracking-[0.3em] text-[#1f7a4d]">Workspace</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight text-[#162117] sm:text-3xl">
								Open the shared workspace
							</h2>
							<p class="mt-3 text-sm leading-6 text-[#445046] sm:text-base">
								The same route contains the request builder, response viewer, templates, and lock states.
							</p>
						</div>

						<div class="flex flex-wrap gap-3">
							<a
								href="/app"
								class="inline-flex items-center gap-2 rounded-full bg-[#1f7a4d] px-6 py-3 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-[#19663f]"
							>
								Open /app
								<ArrowRight size={16} />
							</a>
							<a
								href="/"
								class="rounded-full border border-[#d4ded1] bg-white px-6 py-3 font-medium text-[#162117] transition hover:bg-[#f5f3ed]"
							>
								Back to landing
							</a>
						</div>
					</div>
				</section>
			</main>
		</div>
	</div>
</div>
