<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle,
	} from "$lib/components/ui/card/index.js";
	import Input from "$lib/components/ui/input/input.svelte";
	import Separator from "$lib/components/ui/separator/separator.svelte";
	import {
		Tabs,
		TabsContent,
		TabsList,
		TabsTrigger,
	} from "$lib/components/ui/tabs/index.js";
	import SearchIcon from "@lucide/svelte/icons/search";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CompassIcon from "@lucide/svelte/icons/compass";
	import FolderKanbanIcon from "@lucide/svelte/icons/folder-kanban";
	import ExternalLinkIcon from "@lucide/svelte/icons/external-link";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";

	type TemplateCategoryKey =
		| "all"
		| "rest-basics"
		| "auth-flows"
		| "crud"
		| "pagination"
		| "webhooks"
		| "error-handling";

	type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

	interface WorkspaceTemplate {
		id: string;
		name: string;
		slug: string;
		category: Exclude<TemplateCategoryKey, "all">;
		method: HttpMethod;
		endpoint: string;
		summary: string;
		notes: string;
		tags: string[];
		launchHref?: string;
		previewHref?: string;
		featured?: boolean;
	}

	interface WorkspaceCollection {
		id: string;
		name: string;
		slug: string;
		category: Exclude<TemplateCategoryKey, "all">;
		description: string;
		templateIds: string[];
		launchHref?: string;
		previewHref?: string;
	}

	interface Props {
		templates?: WorkspaceTemplate[];
		collections?: WorkspaceCollection[];
		title?: string;
		subtitle?: string;
		launchBaseHref?: string;
		defaultCategory?: TemplateCategoryKey;
		class?: string;
	}

	const categoryLabels: Record<TemplateCategoryKey, string> = {
		all: "All templates",
		"rest-basics": "REST basics",
		"auth-flows": "Auth flows",
		crud: "CRUD",
		pagination: "Pagination",
		webhooks: "Webhooks",
		"error-handling": "Error handling",
	};

	const categoryDescriptions: Record<Exclude<TemplateCategoryKey, "all">, string> = {
		"rest-basics": "Public request presets and baseline response shapes.",
		"auth-flows": "Authentication-oriented request presets.",
		crud: "Write-oriented request presets.",
		pagination: "Paged and list-oriented responses.",
		webhooks: "Webhook-shaped request presets.",
		"error-handling": "Blocked and failure response states.",
	};

	const categoryOrder: TemplateCategoryKey[] = [
		"all",
		"rest-basics",
		"auth-flows",
		"crud",
		"pagination",
		"webhooks",
		"error-handling",
	];

	const allCategoryDescription = "All templates and collections in one view.";

	const defaultTemplates: WorkspaceTemplate[] = [
		{
			id: "jsonplaceholder-starter-pack",
			name: "JSONPlaceholder starter pack",
			slug: "jsonplaceholder",
			category: "rest-basics",
			method: "GET",
			endpoint: "https://jsonplaceholder.typicode.com/posts/1",
			summary: "A clean first request that demonstrates a fast JSON response with readable fields.",
			notes: "Good for showing the request builder, headers, and the pretty response view.",
			tags: ["guest-safe", "read-only", "instant response"],
			launchHref: "/app?template=jsonplaceholder",
			previewHref: "/app?template=jsonplaceholder&mode=preview",
			featured: true,
		},
		{
			id: "github-public-api-digest",
			name: "GitHub public API digest",
			slug: "github-public-api",
			category: "pagination",
			method: "GET",
			endpoint: "https://api.github.com/repos/sveltejs/svelte/issues?per_page=5",
			summary: "A familiar public API example that hints at paging, rate limits, and headers.",
			notes: "Useful for demonstrating list-style payloads and a response header summary.",
			tags: ["pagination", "public api", "headers"],
			launchHref: "/app?template=github-public-api",
			previewHref: "/app?template=github-public-api&mode=preview",
		},
		{
			id: "weather-demo-glance",
			name: "Weather demo glance",
			slug: "weather-demo",
			category: "rest-basics",
			method: "GET",
			endpoint: "https://api.open-meteo.com/v1/forecast?latitude=22.57&longitude=88.36",
			summary: "A demo endpoint that gives the workspace a realistic service shape and a friendly payload.",
			notes: "Pairs well with the guest mode story because the target stays allowlisted and predictable.",
			tags: ["demo", "allowlisted", "response preview"],
			launchHref: "/app?template=weather-demo",
			previewHref: "/app?template=weather-demo&mode=preview",
		},
		{
			id: "auth-flow-mock",
			name: "Auth flow mock",
			slug: "auth-flow-mock",
			category: "auth-flows",
			method: "POST",
			endpoint: "https://mock.api-testing-kit.dev/auth/login",
			summary: "Authentication-shaped request and response state.",
			notes: "Signed-in request model preview.",
			tags: ["session", "bearer", "sign-in"],
			launchHref: "/app?template=auth-flow-mock",
			previewHref: "/app?template=auth-flow-mock&mode=preview",
		},
		{
			id: "jsonplaceholder-create-post",
			name: "JSONPlaceholder create post",
			slug: "jsonplaceholder-create-post",
			category: "crud",
			method: "POST",
			endpoint: "https://jsonplaceholder.typicode.com/posts",
			summary: "A write-oriented template for creating a post and showing how JSON bodies are assembled.",
			notes: "Pairs with the request body tab and makes the CRUD story concrete.",
			tags: ["create", "json body", "mutation"],
			launchHref: "/app?template=jsonplaceholder-create-post",
			previewHref: "/app?template=jsonplaceholder-create-post&mode=preview",
		},
		{
			id: "github-pagination-page",
			name: "GitHub pagination page",
			slug: "github-pagination-page",
			category: "pagination",
			method: "GET",
			endpoint: "https://api.github.com/search/issues?q=repo:sveltejs/svelte&page=2",
			summary: "A paging-focused template that makes cursor or page-number navigation easy to explain.",
			notes: "Good for calls that need repeatable fetches and response-size awareness.",
			tags: ["page 2", "list data", "rate limit"],
			launchHref: "/app?template=github-pagination-page",
			previewHref: "/app?template=github-pagination-page&mode=preview",
		},
		{
			id: "webhook-delivery-probe",
			name: "Webhook delivery probe",
			slug: "webhook-delivery-probe",
			category: "webhooks",
			method: "POST",
			endpoint: "https://hooks.example.dev/events/customer.created",
			summary: "A webhook-shaped POST request that shows headers, retries, and payload structure.",
			notes: "Helps visitors understand outbound event delivery without exposing arbitrary targets.",
			tags: ["event", "signature", "retry"],
			launchHref: "/app?template=webhook-delivery-probe",
			previewHref: "/app?template=webhook-delivery-probe&mode=preview",
		},
		{
			id: "error-envelope-sample",
			name: "Error envelope sample",
			slug: "error-envelope-sample",
			category: "error-handling",
			method: "GET",
			endpoint: "https://api-testing-kit.dev/errors/429",
			summary: "A deliberate error case that keeps the response viewer honest about blocked and failed states.",
			notes: "Useful for the headers tab, status chips, and the retry/blocked messaging path.",
			tags: ["429", "blocked", "retryable"],
			launchHref: "/app?template=error-envelope-sample",
			previewHref: "/app?template=error-envelope-sample&mode=preview",
		},
	];

	const defaultCollections: WorkspaceCollection[] = [
		{
			id: "starter-set",
			name: "Starter set",
			slug: "starter-set",
			category: "rest-basics",
			description: "The first five minutes of the product: a clean JSON request, a weather check, and a safe launch path.",
			templateIds: ["jsonplaceholder-starter-pack", "weather-demo-glance", "error-envelope-sample"],
			launchHref: "/app?collection=starter-set",
			previewHref: "/app?collection=starter-set&mode=preview",
		},
		{
			id: "auth-walkthrough",
			name: "Auth set",
			slug: "auth-walkthrough",
			category: "auth-flows",
			description: "Authentication-focused templates grouped as one collection.",
			templateIds: ["auth-flow-mock", "error-envelope-sample"],
			launchHref: "/app?collection=auth-walkthrough",
			previewHref: "/app?collection=auth-walkthrough&mode=preview",
		},
		{
			id: "reliability-lab",
			name: "Reliability lab",
			slug: "reliability-lab",
			category: "error-handling",
			description: "Paging, headers, and error states in one compact set that feels safe to browse.",
			templateIds: ["github-public-api-digest", "github-pagination-page", "webhook-delivery-probe"],
			launchHref: "/app?collection=reliability-lab",
			previewHref: "/app?collection=reliability-lab&mode=preview",
		},
	];

	let {
		templates = defaultTemplates,
		collections = defaultCollections,
		title = "Templates and collections",
		subtitle = "Allowlisted presets for REST basics, auth flows, CRUD, pagination, webhooks, and error handling.",
		launchBaseHref = "/app",
		defaultCategory = "all",
		class: className,
	}: Props = $props();

	let activeCategory = $state<TemplateCategoryKey>("all");
	let searchQuery = $state("");

	$effect(() => {
		activeCategory = defaultCategory;
	});

	function normalize(value: string) {
		return value.trim().toLowerCase();
	}

	function methodTone(method: HttpMethod) {
		switch (method) {
			case "GET":
				return "bg-primary-green-soft text-primary-green-deep";
			case "POST":
				return "bg-emerald-100 text-emerald-900";
			case "PUT":
				return "bg-amber-100 text-amber-900";
			case "PATCH":
				return "bg-lime-100 text-lime-900";
			case "DELETE":
				return "bg-rose-100 text-rose-900";
		}
	}

	function matchesTemplate(template: WorkspaceTemplate, query: string) {
		if (!query) {
			return true;
		}

		const haystack = [
			template.name,
			template.summary,
			template.notes,
			template.endpoint,
			template.category,
			...template.tags,
		]
			.join(" ")
			.toLowerCase();

		return haystack.includes(query);
	}

	function filteredTemplates(category: TemplateCategoryKey) {
		const query = normalize(searchQuery);

		return templates.filter((template) => {
			const categoryMatch = category === "all" || template.category === category;
			return categoryMatch && matchesTemplate(template, query);
		});
	}

	function filteredCollections(category: TemplateCategoryKey) {
		const query = normalize(searchQuery);

		return collections.filter((collection) => {
			const categoryMatch = category === "all" || collection.category === category;
			if (!categoryMatch) {
				return false;
			}

			if (!query) {
				return true;
			}

			const templateNames = collection.templateIds
				.map((id) => templates.find((template) => template.id === id)?.name ?? "")
				.join(" ")
				.toLowerCase();

			return [
				collection.name,
				collection.description,
				collection.category,
				templateNames,
			]
				.join(" ")
				.toLowerCase()
				.includes(query);
		});
	}

	function countTemplates(category: TemplateCategoryKey) {
		return filteredTemplates(category).length;
	}

	function countCollections(category: TemplateCategoryKey) {
		return filteredCollections(category).length;
	}

	function resolveTemplateHref(template: WorkspaceTemplate) {
		return template.launchHref ?? `${launchBaseHref}?template=${template.slug}`;
	}

	function resolvePreviewHref(template: WorkspaceTemplate) {
		return template.previewHref ?? `${launchBaseHref}?template=${template.slug}&mode=preview`;
	}

	function resolveCollectionHref(collection: WorkspaceCollection) {
		return collection.launchHref ?? `${launchBaseHref}?collection=${collection.slug}`;
	}

	function resolveCollectionPreviewHref(collection: WorkspaceCollection) {
		return collection.previewHref ?? `${launchBaseHref}?collection=${collection.slug}&mode=preview`;
	}
</script>

<section class={className}>
	<Card class="border-border/70 bg-gradient-to-br from-white via-white to-primary-green-soft/35 shadow-[0_18px_40px_rgba(21,31,23,0.06)]">
		<CardHeader class="gap-5">
			<div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
				<div class="max-w-2xl space-y-2">
					<div class="flex flex-wrap items-center gap-2">
						<Badge variant="outline" class="bg-white/80">
							<ShieldCheckIcon class="size-3.5" />
							Shared workspace
						</Badge>
						<Badge variant="secondary">Templates</Badge>
						<Badge variant="outline">Collections</Badge>
					</div>
					<div>
						<CardTitle class="text-2xl tracking-tight sm:text-3xl">{title}</CardTitle>
						<CardDescription class="mt-2 max-w-2xl text-sm leading-6 text-text-body">
							{subtitle}
						</CardDescription>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-3 sm:grid-cols-4 lg:w-auto">
					<div class="metric-card p-4">
						<p class="section-title">Templates</p>
						<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{templates.length}</p>
						<p class="text-xs text-text-muted">curated examples</p>
					</div>
					<div class="metric-card p-4">
						<p class="section-title">Categories</p>
						<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{categoryOrder.length - 1}</p>
						<p class="text-xs text-text-muted">grouped views</p>
					</div>
					<div class="metric-card p-4">
						<p class="section-title">Collections</p>
						<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{collections.length}</p>
						<p class="text-xs text-text-muted">launch sets</p>
					</div>
					<div class="metric-card p-4">
						<p class="section-title">Launch</p>
						<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">/app</p>
						<p class="text-xs text-text-muted">safe handoff</p>
					</div>
				</div>
			</div>

			<div class="grid gap-3 md:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
				<label class="space-y-2">
					<span class="section-title">Search</span>
					<div class="relative">
						<SearchIcon class="pointer-events-none absolute left-4 top-1/2 size-4 -translate-y-1/2 text-text-muted" />
						<Input
							bind:value={searchQuery}
							placeholder="Search templates, collections, headers, and notes"
							class="h-12 rounded-[16px] border-border/80 bg-white pl-11 text-sm shadow-sm"
						/>
					</div>
				</label>

				<div class="rounded-[18px] border border-border/70 bg-panel-soft p-4">
					<div class="flex items-center gap-2">
						<CompassIcon class="size-4 text-primary-green" />
						<p class="text-sm font-semibold text-text-strong">Launch</p>
					</div>
					<p class="mt-2 text-sm leading-6 text-text-body">
						Each card includes open and preview actions for `/app`.
					</p>
				</div>
			</div>
		</CardHeader>
	</Card>

	<div class="mt-4 grid gap-4 xl:grid-cols-[minmax(0,1.45fr)_minmax(320px,0.95fr)]">
		<Card class="border-border/70 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
			<CardHeader class="gap-4">
				<div class="flex flex-col gap-2 lg:flex-row lg:items-end lg:justify-between">
					<div>
						<CardTitle>Categories</CardTitle>
						<CardDescription>Templates grouped by request shape.</CardDescription>
					</div>
					<div class="flex flex-wrap gap-2 text-xs text-text-muted">
						<span class="rounded-full border border-border/70 bg-panel-soft px-3 py-1">Search is live</span>
						<span class="rounded-full border border-border/70 bg-panel-soft px-3 py-1">Launch to /app</span>
					</div>
				</div>

				<Tabs bind:value={activeCategory} class="gap-4">
					<TabsList class="w-full flex-wrap justify-start bg-panel-soft p-1">
						{#each categoryOrder as category}
							<TabsTrigger value={category} class="gap-2">
								<span>{categoryLabels[category]}</span>
								<span class="rounded-full bg-white/80 px-2 py-0.5 text-[11px] font-medium text-text-muted">
									{countTemplates(category)}
								</span>
							</TabsTrigger>
						{/each}
					</TabsList>

					{#each categoryOrder as category}
						<TabsContent value={category} class="space-y-4">
							<div class="grid gap-4 md:grid-cols-2">
								{#if filteredTemplates(category).length > 0}
									{#each filteredTemplates(category) as template}
										<Card class="border-border/70 bg-white/95 shadow-none">
											<CardHeader class="gap-3">
												<div class="flex items-start justify-between gap-3">
													<div class="space-y-2">
														<div class="flex flex-wrap items-center gap-2">
															<Badge variant="outline" class={methodTone(template.method)}>
																{template.method}
															</Badge>
															<Badge variant="secondary">{categoryLabels[template.category]}</Badge>
															{#if template.featured}
																<Badge variant="outline" class="bg-white">
																	<FolderKanbanIcon class="size-3.5" />
																	Featured
																</Badge>
															{/if}
														</div>
														<CardTitle class="text-base">{template.name}</CardTitle>
													</div>
													<Badge variant="outline">Guest safe</Badge>
												</div>

												<div class="rounded-[16px] border border-border/70 bg-panel-soft px-4 py-3">
													<p class="font-mono text-sm text-text-strong">{template.endpoint}</p>
												</div>
											</CardHeader>

											<CardContent class="space-y-4">
												<p class="text-sm leading-6 text-text-body">{template.summary}</p>
												<p class="text-xs leading-5 text-text-muted">{template.notes}</p>

												<div class="flex flex-wrap gap-2">
													{#each template.tags as tag}
														<span class="rounded-full border border-border/70 bg-panel-soft px-3 py-1 text-xs text-text-muted">
															{tag}
														</span>
													{/each}
												</div>

												<div class="flex flex-wrap gap-2">
													<Button href={resolveTemplateHref(template)} size="sm" class="pill-button">
														Open in /app
														<ArrowRightIcon class="size-4" />
													</Button>
													<Button href={resolvePreviewHref(template)} variant="outline" size="sm">
														Preview payload
														<ExternalLinkIcon class="size-4" />
													</Button>
												</div>
											</CardContent>
										</Card>
									{/each}
								{:else}
									<Card class="border-dashed border-border/70 bg-panel-soft md:col-span-2">
										<CardContent class="flex flex-col items-start gap-3 p-6">
											<p class="text-sm font-semibold text-text-strong">No templates match this category yet.</p>
											<p class="text-sm leading-6 text-text-body">
												No templates match this category.
											</p>
										</CardContent>
									</Card>
								{/if}
							</div>

							<Separator />
							<div class="grid gap-3 sm:grid-cols-3">
								<div class="soft-card p-4">
									<p class="section-title">Category intent</p>
									<p class="mt-2 text-sm font-semibold text-text-strong">{categoryLabels[category]}</p>
									<p class="mt-1 text-sm leading-6 text-text-body">
										{category === "all" ? allCategoryDescription : categoryDescriptions[category]}
									</p>
								</div>
								<div class="soft-card p-4">
									<p class="section-title">Matches</p>
									<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{countTemplates(category)}</p>
									<p class="mt-1 text-sm text-text-body">templates in this view</p>
								</div>
								<div class="soft-card p-4">
									<p class="section-title">Collections</p>
									<p class="mt-2 text-2xl font-semibold tracking-tight text-text-strong">{countCollections(category)}</p>
									<p class="mt-1 text-sm text-text-body">launch sets in this view</p>
								</div>
							</div>
						</TabsContent>
					{/each}
				</Tabs>
			</CardHeader>
		</Card>

		<div class="space-y-4">
			<Card class="border-border/70 shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
				<CardHeader class="gap-3">
					<div class="flex items-center justify-between gap-3">
						<div>
							<CardTitle>Collections</CardTitle>
							<CardDescription>Grouped request sets for the same workspace route.</CardDescription>
						</div>
						<Badge variant="secondary">{filteredCollections(activeCategory).length} visible</Badge>
					</div>
				</CardHeader>

				<CardContent class="space-y-3">
					{#each filteredCollections(activeCategory) as collection}
						<div class="rounded-[18px] border border-border/70 bg-panel-soft p-4">
							<div class="flex items-start justify-between gap-3">
								<div class="space-y-2">
									<div class="flex flex-wrap items-center gap-2">
										<Badge variant="outline">{categoryLabels[collection.category]}</Badge>
										<Badge variant="secondary">{collection.templateIds.length} templates</Badge>
									</div>
									<p class="text-sm font-semibold text-text-strong">{collection.name}</p>
								</div>
								<Button href={resolveCollectionHref(collection)} size="icon-sm" variant="outline" aria-label={`Open ${collection.name}`}>
									<ArrowRightIcon class="size-4" />
								</Button>
							</div>

							<p class="mt-3 text-sm leading-6 text-text-body">{collection.description}</p>

							<div class="mt-3 flex flex-wrap gap-2">
								{#each collection.templateIds as templateId}
									<span class="rounded-full bg-white px-3 py-1 text-xs text-text-muted">
										{templates.find((template) => template.id === templateId)?.name ?? templateId}
									</span>
								{/each}
							</div>

							<div class="mt-4 flex flex-wrap gap-2">
								<Button href={resolveCollectionHref(collection)} size="sm" class="pill-button">
									Launch collection
									<ArrowRightIcon class="size-4" />
								</Button>
								<Button href={resolveCollectionPreviewHref(collection)} size="sm" variant="outline">
									Preview set
									<ExternalLinkIcon class="size-4" />
								</Button>
							</div>
						</div>
					{/each}
				</CardContent>
			</Card>

			<Card class="border-border/70 bg-gradient-to-br from-primary-green-soft/80 to-white shadow-[0_12px_30px_rgba(21,31,23,0.05)]">
				<CardHeader class="gap-3">
						<div class="flex items-center gap-2">
							<CompassIcon class="size-4 text-primary-green" />
							<CardTitle class="text-base">Access</CardTitle>
						</div>
						<CardDescription>Visible templates and mode-aware launch behavior.</CardDescription>
					</CardHeader>

					<CardContent class="space-y-3">
						<div class="rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm leading-6 text-text-body">
							Guest sessions stay on allowlisted templates and collections.
						</div>
						<div class="rounded-[16px] border border-border/70 bg-white px-4 py-3 text-sm leading-6 text-text-body">
							Signed-in sessions unlock the same route with validated custom execution.
						</div>
					</CardContent>
			</Card>
		</div>
	</div>
</section>
