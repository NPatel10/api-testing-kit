<script lang="ts">
	import MarketingSectionHeading from "$lib/components/marketing/marketing-section-heading.svelte";

	const stats = [
		{ label: "Guest-safe endpoints", value: "Curated", detail: "allowlisted demos" },
		{ label: "Response viewer", value: "Pretty + raw", detail: "monospace previews" },
		{ label: "Backend controls", value: "Validated", detail: "destination checks" }
	];
</script>

<section class="space-y-6">
	<div class="inline-flex items-center gap-2 rounded-full border border-border/80 bg-white/80 px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-primary-green">
		<span class="h-2 w-2 rounded-full bg-primary-green"></span>
		Safe demo workspace
	</div>

	<div class="grid gap-6 lg:grid-cols-[minmax(0,1.12fr)_minmax(0,0.88fr)] lg:items-start">
		<div class="space-y-5">
			<div class="space-y-4">
				<h1 class="max-w-3xl text-4xl font-semibold tracking-tight text-text-strong sm:text-5xl lg:text-6xl">
					A premium API testing workspace that is safe to put in front of people
				</h1>
				<p class="max-w-2xl text-sm leading-7 text-text-body sm:text-base">
					API Testing Kit combines a guest-friendly landing experience with a real request/response surface, explicit execution controls, and a backend that refuses to behave like an open proxy.
				</p>
			</div>

			<div class="flex flex-wrap gap-3">
				<a
					href="/app"
					class="rounded-full bg-primary-green px-6 py-3 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-primary-green-hover"
				>
					Open the app
				</a>
				<a
					href="#examples"
					class="rounded-full border border-border/80 bg-white px-6 py-3 font-medium text-text-strong transition hover:bg-surface-soft"
				>
					See example flows
				</a>
			</div>

			<MarketingSectionHeading
				class="pt-2"
				eyebrow="Why it works"
				title="The marketing page and app share the same visual language"
				description="Warm shells, rounded cards, and restrained green accents make the public surface feel like the workspace rather than a separate brochure site."
			/>
		</div>

		<div class="grid gap-3 sm:grid-cols-3 lg:grid-cols-1 xl:grid-cols-3">
			{#each stats as stat}
				<div class="metric-card p-4">
					<p class="text-xs font-medium uppercase tracking-[0.22em] text-text-muted">{stat.label}</p>
					<p class="mt-2 text-base font-semibold text-text-strong">{stat.value}</p>
					<p class="mt-1 text-xs text-text-muted">{stat.detail}</p>
				</div>
			{/each}
		</div>
	</div>
</section>
