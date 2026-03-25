<section class="grid gap-4 xl:grid-cols-[0.95fr_1.05fr]">
	<div class="rounded-[28px] border border-border/80 bg-white/90 p-6 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
		<p class="text-xs font-semibold uppercase tracking-[0.3em] text-primary-green">Architecture</p>
		<h2 class="mt-2 text-2xl font-semibold tracking-tight text-text-strong sm:text-3xl">
			A stack that is straightforward to ship, explain, and extend
		</h2>
		<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
			The repo docs already point to SvelteKit on the front end, Go for the HTTP API, PostgreSQL for persistence, and Goose for migrations. The page keeps that structure visible instead of hiding it behind hype.
		</p>

		<div class="mt-5 grid gap-3 sm:grid-cols-2">
			<div class="rounded-2xl border border-border/70 bg-surface-soft px-4 py-3">
				<p class="text-xs font-semibold uppercase tracking-[0.22em] text-text-muted">Frontend</p>
				<p class="mt-1 text-sm font-medium text-text-strong">SvelteKit, Tailwind, reusable sections</p>
			</div>
			<div class="rounded-2xl border border-border/70 bg-surface-soft px-4 py-3">
				<p class="text-xs font-semibold uppercase tracking-[0.22em] text-text-muted">Backend</p>
				<p class="mt-1 text-sm font-medium text-text-strong">Go HTTP API with explicit routing</p>
			</div>
			<div class="rounded-2xl border border-border/70 bg-surface-soft px-4 py-3">
				<p class="text-xs font-semibold uppercase tracking-[0.22em] text-text-muted">Data</p>
				<p class="mt-1 text-sm font-medium text-text-strong">PostgreSQL, sessions, history, and templates</p>
			</div>
			<div class="rounded-2xl border border-border/70 bg-surface-soft px-4 py-3">
				<p class="text-xs font-semibold uppercase tracking-[0.22em] text-text-muted">Local dev</p>
				<p class="mt-1 text-sm font-medium text-text-strong">Docker-friendly setup and test scripts</p>
			</div>
		</div>
	</div>

	<div class="rounded-[28px] border border-border/80 bg-[linear-gradient(180deg,rgba(31,122,77,0.12),rgba(255,255,255,0.94))] p-6 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
		<p class="text-xs font-semibold uppercase tracking-[0.3em] text-primary-green">Safety model</p>
		<h2 class="mt-2 text-2xl font-semibold tracking-tight text-text-strong sm:text-3xl">
			The app is designed to feel open without becoming an open relay
		</h2>
		<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
			The product direction in the docs is explicit: guests can only use allowlisted examples, and authenticated requests still pass destination validation, redirect checks, and conservative limits.
		</p>

		<div class="mt-5 space-y-3">
			<div class="rounded-2xl border border-border/70 bg-white/80 px-4 py-3 text-sm text-text-body">Block localhost, private ranges, and metadata targets.</div>
			<div class="rounded-2xl border border-border/70 bg-white/80 px-4 py-3 text-sm text-text-body">Validate redirects and keep response previews bounded.</div>
			<div class="rounded-2xl border border-border/70 bg-white/80 px-4 py-3 text-sm text-text-body">Log suspicious attempts for later admin review.</div>
		</div>
	</div>
</section>
