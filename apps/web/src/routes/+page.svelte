<script lang="ts">
	import MarketingShell from '$lib/components/marketing/marketing-shell.svelte';
	import MarketingHeader from '$lib/components/marketing/marketing-header.svelte';
	import MarketingCTABand from '$lib/components/marketing/marketing-cta-band.svelte';
	import MarketingSectionHeading from '$lib/components/marketing/marketing-section-heading.svelte';
	import HeroSection from '$lib/components/marketing/hero-section.svelte';
	import ProductPreview from '$lib/components/marketing/product-preview.svelte';
	import FeatureGrid from '$lib/components/marketing/feature-grid.svelte';
	import UseCases from '$lib/components/marketing/use-cases.svelte';
	import CredibilitySection from '$lib/components/marketing/credibility-section.svelte';
	import SiteFooter from '$lib/components/marketing/site-footer.svelte';

	const features = [
		{
			title: 'Guest-safe request exploration',
			description:
				'Visitors can try curated examples immediately while the backend keeps arbitrary outbound targets locked down.',
			icon: 'shield'
		},
		{
			title: 'Split request and response workspace',
			description:
				'The core canvas keeps method, URL, headers, body, and response metadata visible in one working surface.',
			icon: 'layout'
		},
		{
			title: 'Response inspection that feels technical',
			description:
				'Pretty, raw, and headers views stay readable with monospace surfaces, status badges, and concise metadata.',
			icon: 'code'
		},
		{
			title: 'Safety-first outbound execution',
			description:
				'Private ranges, metadata IPs, and unsupported protocols stay blocked so the product never behaves like an open proxy.',
			icon: 'lock'
		},
		{
			title: 'Reusable collections and history',
			description:
				'Authenticated users get save flows, re-run paths, and durable request history without changing the app model.',
			icon: 'stack'
		},
		{
			title: 'Built for polished demos',
			description:
				'The landing page and workspace share the same warm, premium visual language so the product feels coherent.',
			icon: 'spark'
		}
	];

	const useCases = [
		{
			title: 'Portfolio reviews',
			text: 'A public-facing workspace that shows real request and response surfaces without opening outbound execution.'
		},
		{
			title: 'API debugging',
			text: 'One surface for request editing, execution, and response inspection.'
		},
		{
			title: 'Team onboarding',
			text: 'Shared templates, collections, and history inside the same product model.'
		}
	];

	const topLinks = [
		{ href: '#features', label: 'Features' },
		{ href: '#examples', label: 'Examples' },
		{ href: '#trust', label: 'Safety' }
	];
</script>

<svelte:head>
	<title>API Testing Kit - Safe API workspace for demos and real requests</title>
	<meta
		name="description"
		content="A warm, premium API testing workspace with guest-safe examples, strong outbound controls, and a polished request-response experience."
	/>
</svelte:head>

<MarketingShell>
	<div class="overflow-hidden rounded-[32px] border border-border/80 bg-[rgba(247,245,240,0.92)] shadow-[0_24px_60px_rgba(21,31,23,0.08)] backdrop-blur">
		<MarketingHeader links={topLinks} />

		<main class="space-y-12 px-5 py-6 sm:px-6 lg:px-8 lg:py-8">
			<HeroSection />
			<ProductPreview />

			<section id="features" class="scroll-mt-24">
				<FeatureGrid {features} />
			</section>

			<section id="examples" class="scroll-mt-24">
				<UseCases {useCases} />
			</section>

			<section id="trust" class="scroll-mt-24 space-y-6">
				<MarketingSectionHeading
					align="center"
					eyebrow="Credibility"
					title="A product that is designed for public demos and honest engineering narratives"
					description="The public site stays aligned with the app surface, the runtime safety model, and the local-development story in the docs."
				/>
				<CredibilitySection />
			</section>

			<MarketingCTABand
				eyebrow="Workspace"
				title="Open the shared request workspace"
				description="Request builder, response viewer, templates, and guest lock states live on the same route."
				primaryHref="/app"
				primaryLabel="Open /app"
				secondaryHref="#features"
				secondaryLabel="View features"
			/>
		</main>

		<SiteFooter />
	</div>
</MarketingShell>
