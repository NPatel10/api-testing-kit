<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import Input from "$lib/components/ui/input/input.svelte";
	import Separator from "$lib/components/ui/separator/separator.svelte";
	import { Tabs, TabsList, TabsTrigger } from "$lib/components/ui/tabs/index.js";
	import {
		createPreviewHref,
		groupTemplatesByCategory,
		templateMatchesFilter,
		templateMatchesSearch,
		templatesPageCategories,
		templatesPageCategoryLabels,
		templatesPageFilterLabels,
		type TemplatesPageData,
		type TemplatesPageCategory,
		type TemplatesPageFilter,
		type TemplatesPageTemplate,
		formatPreviewBody,
		selectInitialTemplate,
	} from "$lib/components/templates-page/templates-page-data";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import SearchIcon from "@lucide/svelte/icons/search";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";
	import Layers3Icon from "@lucide/svelte/icons/layers-3";
	import CircleCheckBigIcon from "@lucide/svelte/icons/circle-check-big";
	import FilterIcon from "@lucide/svelte/icons/filter";
	import Code2Icon from "@lucide/svelte/icons/code-2";
	import EyeIcon from "@lucide/svelte/icons/eye";

	let { data }: { data: TemplatesPageData } = $props();

	let searchQuery = $state("");
	let activeCategory = $state<TemplatesPageCategory>("all");
	let activeFilter = $state<TemplatesPageFilter>("all");
	let selectedSlug = $state("");
	let hasInitializedSelection = $state(false);

	const categoryCount = (category: TemplatesPageCategory) =>
		data.templates.filter((template) => category === "all" || template.category === category).length;

	const visibleTemplates = $derived(
		data.templates.filter(
			(template) =>
				templateMatchesSearch(template, searchQuery) &&
				(activeCategory === "all" || template.category === activeCategory) &&
				templateMatchesFilter(template, activeFilter),
		),
	);

	const selectedTemplate = $derived(selectInitialTemplate(visibleTemplates, selectedSlug));
	const categoryGroups = $derived(groupTemplatesByCategory(data.templates));
	const filteredCount = $derived(visibleTemplates.length);
	const totalCount = $derived(data.templates.length);
	const bodyCount = $derived(data.templates.filter((template) => template.request.bodyMode !== "none").length);
	const errorCount = $derived(data.templates.filter((template) => template.responsePreview.status >= 400).length);

	$effect(() => {
		if (!hasInitializedSelection) {
			selectedSlug = data.selectedSlug ?? data.templates[0]?.slug ?? "";
			hasInitializedSelection = true;
		}

		if (selectedTemplate && selectedTemplate.slug !== selectedSlug) {
			selectedSlug = selectedTemplate.slug;
		}
	});

	function resetFilters() {
		searchQuery = "";
		activeCategory = "all";
		activeFilter = "all";
	}

	function filterButtonClass(value: TemplatesPageFilter) {
		return value === activeFilter
			? "border-primary-green bg-primary-green-soft text-primary-green-deep"
			: "border-border/80 bg-white text-text-body hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong";
	}

	function categoryButtonClass(value: TemplatesPageCategory) {
		return value === activeCategory
			? "border-primary-green bg-primary-green text-white shadow-[0_10px_20px_rgba(31,122,77,0.18)]"
			: "border-border/80 bg-white text-text-body hover:border-primary-green-soft hover:bg-primary-green-soft/40 hover:text-text-strong";
	}

	function selectedCardClass(slug: string) {
		return slug === selectedTemplate?.slug
			? "border-primary-green bg-primary-green-soft/35 shadow-[0_16px_34px_rgba(31,122,77,0.12)]"
			: "border-border/80 bg-white shadow-[0_10px_24px_rgba(21,31,23,0.05)] hover:-translate-y-0.5 hover:border-primary-green-soft";
	}

	function openTemplate(slug: string) {
		selectedSlug = slug;
	}

	function launchHref(template: TemplatesPageTemplate) {
		return template.launchHref;
	}

	function previewHref(template: TemplatesPageTemplate) {
		return template.previewHref ?? createPreviewHref(template.slug);
	}

	function previewSnippet(template: TemplatesPageTemplate) {
		return formatPreviewBody(template.responsePreview.body).slice(0, 360);
	}
</script>

<svelte:head>
	<title>Templates - API Testing Kit</title>
	<meta
		name="description"
		content="Browse safe API templates with filters, search, and live previews that launch directly into the shared /app workspace."
	/>
</svelte:head>

<section class="relative isolate overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(31,122,77,0.14),transparent_28%),radial-gradient(circle_at_top_right,rgba(111,142,163,0.12),transparent_26%),linear-gradient(180deg,#f2f0ea_0%,#ece8df_100%)] text-text-strong">
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-28 top-16 h-72 w-72 rounded-full bg-[#1f7a4d]/10 blur-3xl"></div>
		<div class="absolute right-[-5rem] top-1/4 h-96 w-96 rounded-full bg-[#dcefe3] blur-3xl"></div>
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
							<p class="text-sm font-semibold tracking-tight text-text-strong">API Testing Kit</p>
							<p class="text-xs text-text-muted">Templates page with live backend data and safe launches</p>
						</div>
					</div>

					<div class="flex flex-wrap items-center gap-2 text-sm">
						<a href="/" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Home</a>
						<a href="/docs" class="rounded-full px-4 py-2 text-text-body transition hover:bg-surface-soft hover:text-text-strong">Docs</a>
						<a href="/app" class="rounded-full bg-primary-green px-5 py-2.5 font-medium text-white shadow-[0_12px_28px_rgba(31,122,77,0.28)] transition hover:bg-primary-green-hover">Open /app</a>
					</div>
				</div>
			</header>

			<main class="space-y-10 px-5 py-6 sm:px-6 lg:px-8 lg:py-8">
				<section class="grid gap-6 xl:grid-cols-[minmax(0,1.3fr)_minmax(320px,0.7fr)] xl:items-start">
					<div class="space-y-6">
						<div class="inline-flex items-center gap-2 rounded-full border border-[#d9e7d8] bg-white/80 px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-[#1f7a4d]">
							<span class="h-2 w-2 rounded-full bg-[#1f7a4d]"></span>
							Guided browsing
						</div>

						<div class="space-y-4">
							<h1 class="max-w-3xl text-4xl font-semibold tracking-tight text-text-strong sm:text-5xl lg:text-6xl">
								Safe API templates with category filters, search, and a real launch path
							</h1>
							<p class="max-w-2xl text-sm leading-7 text-text-body sm:text-base">
								Browse curated examples by category, inspect the request and response preview, then jump straight into the shared `/app` workspace with the selected template.
							</p>
						</div>

						<div class="grid gap-3 sm:grid-cols-3">
							<div class="metric-card p-4">
								<p class="section-title">Templates</p>
								<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{totalCount}</p>
								<p class="text-xs text-text-muted">live or fallback data</p>
							</div>
							<div class="metric-card p-4">
								<p class="section-title">Categories</p>
								<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{templatesPageCategories.length - 1}</p>
								<p class="text-xs text-text-muted">grouped views</p>
							</div>
							<div class="metric-card p-4">
								<p class="section-title">Source</p>
								<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{data.sourceLabel}</p>
								<p class="text-xs text-text-muted">GET /api/v1/templates</p>
							</div>
						</div>
					</div>

					<Card class="border-border/80 bg-gradient-to-br from-white via-white to-primary-green-soft/35 shadow-[0_18px_40px_rgba(21,31,23,0.06)]">
						<CardHeader class="gap-3">
							<div class="flex items-center gap-2">
								<SparklesIcon class="size-4 text-primary-green" />
								<CardTitle class="text-base">Launch snapshot</CardTitle>
							</div>
							<CardDescription>
								Selected template preview, launch affordance, and a clear path back into `/app`.
							</CardDescription>
						</CardHeader>
						<CardContent class="space-y-4">
							<div class="rounded-[18px] border border-border/70 bg-panel-soft p-4">
								<div class="flex items-center justify-between gap-3">
									<p class="text-xs font-semibold uppercase tracking-[0.18em] text-text-muted">Selected template</p>
									<Badge variant="secondary">{selectedTemplate?.category ?? "All templates"}</Badge>
								</div>
								<p class="mt-2 text-lg font-semibold tracking-tight text-text-strong">
									{selectedTemplate?.title ?? "No template selected"}
								</p>
								<p class="mt-2 text-sm leading-6 text-text-body">
									{selectedTemplate?.summary ?? "Use the search bar or category tabs to pick a template."}
								</p>
							</div>

							<div class="grid gap-3 sm:grid-cols-2">
								<div class="soft-card p-4">
									<p class="section-title">Response status</p>
									<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">
										{selectedTemplate?.responsePreview.status ?? 200}
									</p>
									<p class="text-xs text-text-muted">{selectedTemplate?.responsePreview.contentType ?? "application/json"}</p>
								</div>
								<div class="soft-card p-4">
									<p class="section-title">Request type</p>
									<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">
										{selectedTemplate?.request.method ?? "GET"}
									</p>
									<p class="text-xs text-text-muted">{selectedTemplate?.request.bodyMode ?? "none"}</p>
								</div>
							</div>

							<div class="flex flex-wrap gap-3">
								{#if selectedTemplate}
									<Button href={launchHref(selectedTemplate)} class="pill-button">
										Open in /app
										<ArrowRightIcon class="size-4" />
									</Button>
									<Button href={previewHref(selectedTemplate)} variant="outline" class="pill-button">
										<EyeIcon class="size-4" />
										Preview in /app
									</Button>
								{/if}
							</div>
						</CardContent>
					</Card>
				</section>

				<Card class="border-border/80 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
					<CardHeader class="gap-4">
						<div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
							<div>
								<CardTitle class="text-2xl tracking-tight sm:text-3xl">Filter templates</CardTitle>
								<CardDescription class="mt-2 max-w-2xl text-sm leading-6 text-text-body">
									Search by target, summary, tag, or override. Narrow the surface with category tabs and quick shape filters.
								</CardDescription>
							</div>
							<div class="flex flex-wrap gap-2 text-xs text-text-muted">
								<span class="rounded-full border border-border/70 bg-panel-soft px-3 py-1">{filteredCount} visible</span>
								<span class="rounded-full border border-border/70 bg-panel-soft px-3 py-1">{bodyCount} with bodies</span>
								<span class="rounded-full border border-border/70 bg-panel-soft px-3 py-1">{errorCount} error states</span>
							</div>
						</div>

						<div class="grid gap-3 xl:grid-cols-[minmax(0,1.2fr)_minmax(0,0.8fr)]">
							<label class="space-y-2">
								<span class="section-title">Search</span>
								<div class="relative">
									<SearchIcon class="pointer-events-none absolute left-4 top-1/2 size-4 -translate-y-1/2 text-text-muted" />
									<Input
										bind:value={searchQuery}
										placeholder="Search templates, tags, targets, and overrides"
										class="h-12 rounded-[16px] border-border/80 bg-white pl-11 text-sm shadow-sm"
									/>
								</div>
							</label>

							<div class="rounded-[18px] border border-border/70 bg-panel-soft p-4">
								<div class="flex items-center gap-2">
									<FilterIcon class="size-4 text-primary-green" />
									<p class="text-sm font-semibold text-text-strong">Quick filters</p>
								</div>
								<div class="mt-3 flex flex-wrap gap-2">
									{#each (["all", "body", "query", "headers", "write", "error"] as TemplatesPageFilter[]) as filter}
										<button
											type="button"
											class={`rounded-full border px-3 py-1.5 text-xs font-medium transition ${filterButtonClass(filter)}`}
											onclick={() => (activeFilter = filter)}
										>
											{templatesPageFilterLabels[filter]}
										</button>
									{/each}
								</div>
							</div>
						</div>

						<Tabs bind:value={activeCategory} class="gap-4">
							<TabsList class="w-full flex-wrap justify-start bg-panel-soft p-1">
								{#each templatesPageCategories as category}
									<TabsTrigger value={category} class={categoryButtonClass(category)}>
										<span>{templatesPageCategoryLabels[category]}</span>
										<span class="rounded-full bg-white/80 px-2 py-0.5 text-[11px] font-medium text-text-muted">
											{categoryCount(category)}
										</span>
									</TabsTrigger>
								{/each}
							</TabsList>
						</Tabs>
					</CardHeader>
				</Card>

				<div class="grid gap-4 xl:grid-cols-[minmax(0,1.45fr)_minmax(320px,0.8fr)]">
					<Card class="border-border/80 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
						<CardHeader class="gap-3">
							<div class="flex items-center justify-between gap-3">
								<div>
									<CardTitle>Template cards</CardTitle>
									<CardDescription>
										Click a card to inspect the request and response preview in the detail panel.
									</CardDescription>
								</div>
								<Badge variant="secondary">{visibleTemplates.length} match filters</Badge>
							</div>
						</CardHeader>

						<CardContent>
							{#if visibleTemplates.length > 0}
								<div class="grid gap-4 md:grid-cols-2">
									{#each visibleTemplates as template}
										<article class={`rounded-[24px] border p-5 transition ${selectedCardClass(template.slug)}`}>
											<div class="flex items-start justify-between gap-3">
												<div class="space-y-3">
													<div class="flex flex-wrap items-center gap-2">
														<Badge variant="outline" class="bg-white">{template.request.method}</Badge>
														<Badge variant="secondary">{template.category}</Badge>
														{#if template.guestSafe}
															<Badge variant="outline" class="bg-white">
																<CircleCheckBigIcon class="size-3.5" />
																Guest safe
															</Badge>
														{/if}
													</div>
													<div>
														<h3 class="text-base font-semibold tracking-tight text-text-strong">{template.title}</h3>
														<p class="mt-2 text-sm leading-6 text-text-body">{template.summary}</p>
													</div>
												</div>

												<button
													type="button"
													class="rounded-full border border-border/80 bg-white p-2 text-text-muted transition hover:border-primary-green-soft hover:text-text-strong"
													onclick={() => openTemplate(template.slug)}
													aria-label={`Select ${template.title}`}
												>
													<Layers3Icon class="size-4" />
												</button>
											</div>

											<div class="mt-4 rounded-[18px] border border-border/70 bg-panel-soft px-4 py-3">
												<p class="font-mono text-sm text-text-strong">{template.allowlistedTarget}</p>
											</div>

											<div class="mt-4 flex flex-wrap gap-2">
												{#each template.allowedOverrides.slice(0, 4) as override}
													<span class="rounded-full border border-border/70 bg-white px-3 py-1 text-xs text-text-muted">
														{override}
													</span>
												{/each}
											</div>

											<div class="mt-4 flex flex-wrap items-center justify-between gap-3">
												<div class="flex flex-wrap gap-2 text-xs text-text-muted">
													<span class="rounded-full bg-white px-3 py-1">{template.responsePreview.status}</span>
													<span class="rounded-full bg-white px-3 py-1">{template.responsePreview.durationMs} ms</span>
													<span class="rounded-full bg-white px-3 py-1">{template.responsePreview.size}</span>
												</div>

												<div class="flex flex-wrap gap-2">
													<Button href={launchHref(template)} size="sm" class="pill-button">
														Open
														<ArrowRightIcon class="size-4" />
													</Button>
													<Button href={previewHref(template)} size="sm" variant="outline">
														Preview
														<EyeIcon class="size-4" />
													</Button>
												</div>
											</div>
										</article>
									{/each}
								</div>
							{:else}
								<div class="rounded-[24px] border border-dashed border-border/70 bg-panel-soft p-8">
									<p class="text-lg font-semibold tracking-tight text-text-strong">No templates match the current filters.</p>
									<p class="mt-2 max-w-2xl text-sm leading-6 text-text-body">
										Adjust the search, switch categories, or clear the quick filters to continue browsing.
									</p>
									<div class="mt-4">
										<Button variant="outline" class="pill-button" onclick={resetFilters}>
											Reset filters
										</Button>
									</div>
								</div>
							{/if}
						</CardContent>
					</Card>

					<Card class="border-border/80 shadow-[0_12px_30px_rgba(21,31,23,0.05)] xl:sticky xl:top-6 xl:self-start">
						<CardHeader class="gap-3">
							<div class="flex items-center justify-between gap-3">
								<div>
									<CardTitle>Template preview</CardTitle>
									<CardDescription>Request, response, and launch details for the current selection.</CardDescription>
								</div>
								<Badge variant="secondary">{selectedTemplate?.source === "live" ? "Live backend" : "Fallback"}</Badge>
							</div>
						</CardHeader>

						<CardContent class="space-y-4">
							{#if selectedTemplate}
								<div class="rounded-[18px] border border-border/70 bg-panel-soft p-4">
									<p class="section-title">Selected</p>
									<p class="mt-2 text-xl font-semibold tracking-tight text-text-strong">{selectedTemplate.title}</p>
									<p class="mt-2 text-sm leading-6 text-text-body">{selectedTemplate.description}</p>
								</div>

								<div class="space-y-4 rounded-[20px] border border-border/70 bg-white p-4">
									<div class="flex flex-wrap gap-2">
										{#each selectedTemplate.tags as tag}
											<span class="rounded-full bg-primary-green-soft px-3 py-1 text-xs font-medium text-primary-green-deep">{tag}</span>
										{/each}
									</div>

									<Separator />

									<div class="grid gap-3 sm:grid-cols-2">
										<div class="soft-card p-3">
											<p class="section-title">Method</p>
											<p class="mt-1 font-mono text-sm text-text-strong">{selectedTemplate.request.method}</p>
										</div>
										<div class="soft-card p-3">
											<p class="section-title">Body mode</p>
											<p class="mt-1 font-mono text-sm text-text-strong">{selectedTemplate.request.bodyMode}</p>
										</div>
										<div class="soft-card p-3">
											<p class="section-title">Status</p>
											<p class="mt-1 font-mono text-sm text-text-strong">{selectedTemplate.responsePreview.status}</p>
										</div>
										<div class="soft-card p-3">
											<p class="section-title">Preview size</p>
											<p class="mt-1 font-mono text-sm text-text-strong">{selectedTemplate.responsePreview.size}</p>
										</div>
									</div>

									<div class="rounded-[18px] border border-border/70 bg-surface-soft p-4">
										<p class="section-title">Allowlisted target</p>
										<p class="mt-2 font-mono text-sm leading-6 text-text-strong">{selectedTemplate.allowlistedTarget}</p>
									</div>
								</div>

								<div class="rounded-[20px] border border-border/70 bg-white p-4">
									<div class="flex items-center gap-2">
										<Code2Icon class="size-4 text-primary-green" />
										<p class="text-sm font-semibold text-text-strong">Preview body</p>
									</div>
									<pre class="mt-3 overflow-x-auto rounded-[18px] border border-border/70 bg-[#162117] px-4 py-4 text-xs leading-6 text-[#d8ebdf]"><code>{previewSnippet(selectedTemplate)}</code></pre>
								</div>

								<div class="flex flex-wrap gap-3">
									<Button href={launchHref(selectedTemplate)} class="pill-button">
										Open in /app
										<ArrowRightIcon class="size-4" />
									</Button>
									<Button href={previewHref(selectedTemplate)} variant="outline" class="pill-button">
										<EyeIcon class="size-4" />
										Preview in /app
									</Button>
								</div>
							{:else}
								<div class="rounded-[24px] border border-dashed border-border/70 bg-panel-soft p-6">
									<p class="text-sm font-semibold text-text-strong">No selection available.</p>
									<p class="mt-2 text-sm leading-6 text-text-body">
										Reset the filters to reveal templates, then open one of the cards.
									</p>
									<Button variant="outline" class="mt-4 pill-button" onclick={resetFilters}>
										Reset filters
									</Button>
								</div>
							{/if}
						</CardContent>
					</Card>
				</div>

				<section class="space-y-5">
					<div class="max-w-3xl">
						<p class="section-title">Category collections</p>
						<h2 class="mt-2 text-2xl font-semibold tracking-tight text-text-strong sm:text-3xl">
							The same templates, regrouped as launchable collections
						</h2>
						<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
							Each category cluster can feed users back into the main workspace with a known starting point and a consistent safe target.
						</p>
					</div>

					<div class="grid gap-4 lg:grid-cols-2 xl:grid-cols-3">
						{#each categoryGroups as group}
							<Card class="border-border/80 bg-white/90 shadow-[0_10px_24px_rgba(21,31,23,0.05)]">
								<CardHeader class="gap-3">
									<div class="flex items-start justify-between gap-3">
										<div>
											<CardTitle class="text-base">{group.label}</CardTitle>
											<CardDescription class="mt-2 text-sm leading-6 text-text-body">{group.description}</CardDescription>
										</div>
										<Badge variant="secondary">{group.templates.length}</Badge>
									</div>
								</CardHeader>
								<CardContent class="space-y-4">
									<div class="space-y-2">
										{#each group.templates.slice(0, 3) as template}
											<button
												type="button"
												class="flex w-full items-center justify-between rounded-[16px] border border-border/70 bg-panel-soft px-4 py-3 text-left transition hover:border-primary-green-soft hover:bg-primary-green-soft/35"
												onclick={() => openTemplate(template.slug)}
											>
												<div class="min-w-0">
													<p class="truncate text-sm font-medium text-text-strong">{template.title}</p>
													<p class="mt-1 truncate text-xs text-text-muted">{template.allowlistedTarget}</p>
												</div>
												<ArrowRightIcon class="size-4 shrink-0 text-text-muted" />
											</button>
										{/each}
									</div>

									<div class="flex flex-wrap gap-2">
										<Button
											variant="outline"
											size="sm"
											class="pill-button"
											onclick={() => {
												activeCategory = group.category;
												searchQuery = "";
												activeFilter = "all";
											}}
										>
											View category
										</Button>
										{#if group.templates[0]}
											<Button href={launchHref(group.templates[0])} size="sm" class="pill-button">
												Launch first
												<ArrowRightIcon class="size-4" />
											</Button>
										{/if}
									</div>
								</CardContent>
							</Card>
						{/each}
					</div>
				</section>

				<section class="rounded-[30px] border border-[#dfe8dd] bg-[linear-gradient(135deg,rgba(31,122,77,0.16),rgba(255,255,255,0.96))] p-6 shadow-[0_18px_40px_rgba(21,31,23,0.08)] sm:p-8">
					<div class="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
						<div class="max-w-2xl">
							<p class="section-title">Ready to explore</p>
							<h2 class="mt-2 text-2xl font-semibold tracking-tight text-text-strong sm:text-3xl">
								Open `/app` with a selected template and continue in the shared workspace
							</h2>
							<p class="mt-3 text-sm leading-6 text-text-body sm:text-base">
								The templates page is a stepping stone, not a dead end. Each launch path lands in the real app shell where guest constraints and authenticated execution diverge.
							</p>
						</div>

						<div class="flex flex-wrap gap-3">
							<Button href="/app" class="pill-button">
								Open /app
								<ArrowRightIcon class="size-4" />
							</Button>
							<Button href="/" variant="outline" class="pill-button">
								Back to home
							</Button>
						</div>
					</div>
				</section>
			</main>
		</div>
	</div>
</section>
