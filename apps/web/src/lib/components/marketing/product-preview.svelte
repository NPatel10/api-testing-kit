<script lang="ts">
	const response = String.raw`{
  "name": "sveltejs/kit",
  "stars": 18800,
  "safety": "validated"
}`;
</script>

<section class="space-y-5">
	<div class="max-w-3xl">
		<p class="text-xs font-semibold uppercase tracking-[0.3em] text-primary-green">Workspace preview</p>
		<h2 class="mt-2 text-2xl font-semibold tracking-tight text-text-strong sm:text-3xl">
			The request builder and response viewer stay visible together
		</h2>
		<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
			Guests can inspect curated examples, while signed-in users get the full workspace with saved requests, collections, and stricter outbound control.
		</p>
	</div>

	<div class="relative">
		<div class="absolute inset-0 -z-10 rounded-[32px] bg-[radial-gradient(circle_at_top_right,rgba(31,122,77,0.18),transparent_52%)]"></div>
		<div class="rounded-[32px] border border-border/80 bg-[linear-gradient(180deg,rgba(255,255,255,0.9),rgba(245,243,237,0.92))] p-4 shadow-[0_18px_40px_rgba(21,31,23,0.08)] sm:p-5">
			<div class="rounded-[28px] border border-border/80 bg-white p-4">
				<div class="flex items-center justify-between border-b border-border/70 pb-4">
					<div>
						<p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-green">Live preview</p>
						<p class="mt-1 text-sm text-text-muted">Curated request + response canvas</p>
					</div>
					<div class="rounded-full bg-primary-green-soft px-3 py-1 text-xs font-semibold text-primary-green-deep">Guest mode</div>
				</div>

				<div class="mt-4 grid gap-3 lg:grid-cols-[0.92fr_1.08fr]">
					<div class="soft-card p-4">
						<div class="flex items-center justify-between">
							<p class="text-sm font-semibold text-text-strong">Request</p>
							<span class="rounded-full bg-primary-green-soft px-2.5 py-1 text-xs font-semibold text-primary-green-deep">GET</span>
						</div>
						<div class="mt-3 rounded-2xl border border-border/70 bg-white px-3 py-2 text-xs font-mono text-text-body">
							https://api.github.com/repos/sveltejs/kit
						</div>
						<div class="mt-3 space-y-2">
							<div class="flex items-center justify-between rounded-2xl bg-white px-3 py-2 text-xs text-text-body">
								<span>Headers</span><span>4 items</span>
							</div>
							<div class="flex items-center justify-between rounded-2xl bg-white px-3 py-2 text-xs text-text-body">
								<span>Body</span><span>Locked for this example</span>
							</div>
							<div class="flex items-center justify-between rounded-2xl bg-white px-3 py-2 text-xs text-text-body">
								<span>Auth</span><span>None</span>
							</div>
						</div>
					</div>

					<div class="rounded-[24px] border border-border/80 bg-white p-4 shadow-[0_12px_30px_rgba(21,31,23,0.06)]">
						<div class="flex items-center justify-between border-b border-border/70 pb-4">
							<div>
								<p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-green">Response</p>
								<p class="mt-1 text-sm text-text-muted">Pretty JSON with metadata</p>
							</div>
							<div class="rounded-full bg-primary-green px-3 py-1 text-xs font-semibold text-white">200 OK</div>
						</div>

						<div class="mt-4 grid gap-3 sm:grid-cols-3">
							<div class="rounded-2xl bg-surface-soft px-3 py-3">
								<p class="text-[11px] uppercase tracking-[0.22em] text-text-muted">Time</p>
								<p class="mt-1 text-sm font-semibold text-text-strong">183 ms</p>
							</div>
							<div class="rounded-2xl bg-surface-soft px-3 py-3">
								<p class="text-[11px] uppercase tracking-[0.22em] text-text-muted">Size</p>
								<p class="mt-1 text-sm font-semibold text-text-strong">18.2 KB</p>
							</div>
							<div class="rounded-2xl bg-surface-soft px-3 py-3">
								<p class="text-[11px] uppercase tracking-[0.22em] text-text-muted">Type</p>
								<p class="mt-1 text-sm font-semibold text-text-strong">application/json</p>
							</div>
						</div>

						<div class="mt-4 overflow-hidden rounded-[22px] border border-border/80 bg-surface-soft">
							<div class="flex gap-3 border-b border-border/70 px-4 py-3 text-xs font-semibold text-text-muted">
								<span class="text-text-strong">Pretty</span>
								<span>Raw</span>
								<span>Headers</span>
							</div>
							<pre class="overflow-x-auto px-4 py-4 text-xs leading-6 text-text-strong"><code>{response}</code></pre>
						</div>
					</div>
				</div>

				<div class="mt-4 rounded-[24px] border border-border/80 bg-surface-soft px-4 py-3">
					<div class="flex flex-wrap items-center justify-between gap-3">
						<p class="text-sm font-semibold text-text-strong">Snippet drawer</p>
						<div class="flex flex-wrap gap-2 text-xs">
							<span class="rounded-full bg-white px-3 py-1 text-text-body">curl</span>
							<span class="rounded-full bg-white px-3 py-1 text-text-body">fetch</span>
							<span class="rounded-full bg-white px-3 py-1 text-text-body">python requests</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
