<script lang="ts">
	import { ArrowRight, FileCode2, PanelTop, Users } from '@lucide/svelte';
	import MarketingSectionHeading from '$lib/components/marketing/marketing-section-heading.svelte';

	let {
		useCases = []
	}: {
		useCases: {
			title: string;
			text: string;
			steps: string[];
		}[];
	} = $props();

	const icons = [PanelTop, FileCode2, Users];
</script>

<div class="space-y-6">
	<MarketingSectionHeading
		eyebrow="Examples"
		title="Use cases that show the product, not just the copy"
		description="The page uses a few concrete workflows so visitors can see how the shell behaves across exploration, debugging, and team onboarding."
	/>

	<div class="grid gap-4 xl:grid-cols-3">
		{#each useCases as useCase, index}
			{@const Icon = icons[index]}
			<article class="rounded-[28px] border border-border/80 bg-[linear-gradient(180deg,rgba(255,255,255,0.92),rgba(245,243,237,0.96))] p-5 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
				<div class="flex items-start justify-between gap-3">
					<div class="grid h-11 w-11 place-items-center rounded-2xl bg-primary-green text-white shadow-[0_10px_24px_rgba(31,122,77,0.22)]">
						<Icon size={18} strokeWidth={2.1} />
					</div>
					<span class="rounded-full border border-primary-green-soft bg-white px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.22em] text-primary-green">
						0{index + 1}
					</span>
				</div>

				<h3 class="mt-5 text-base font-semibold tracking-tight text-text-strong">{useCase.title}</h3>
				<p class="mt-2 text-sm leading-6 text-text-body">{useCase.text}</p>

				<ul class="mt-4 space-y-2">
					{#each useCase.steps as step}
						<li class="flex items-start gap-2 text-sm leading-6 text-text-body">
							<ArrowRight size={15} class="mt-1 shrink-0 text-primary-green" />
							<span>{step}</span>
						</li>
					{/each}
				</ul>
			</article>
		{/each}
	</div>
</div>
