import {
	getWorkspaceCollection,
	workspaceStates,
	workspaceTemplates,
	type WorkspaceMode,
	type WorkspaceScope,
	type WorkspaceRequestMethod,
	type WorkspaceTemplateCategory,
} from "$lib/mocks/workspace-state";

export type CollectionDetailStatus = "missing" | "locked" | "ready";

export interface CollectionDetailActionLink {
	label: string;
	href: string;
	variant: "default" | "outline";
}

export interface CollectionDetailMetadataItem {
	label: string;
	value: string;
	detail: string;
}

export interface CollectionDetailRequestItem {
	slug: string;
	title: string;
	category: WorkspaceTemplateCategory;
	method: WorkspaceRequestMethod;
	endpoint: string;
	summary: string;
	description: string;
	safeOverrides: readonly string[];
	responseStatus: number;
	responseStatusText: string;
	responseTimeMs: number;
	responseSizeLabel: string;
	responseContentType: string;
	launchHref: string;
	previewHref: string;
}

export interface CollectionDetailGroup {
	id: string;
	label: string;
	summary: string;
	requestCount: number;
	requests: CollectionDetailRequestItem[];
}

export interface CollectionDetailCollection {
	id: string;
	title: string;
	description: string;
	scope: WorkspaceScope;
	requestCount: number;
	templateSlugs: readonly string[];
	badge: string;
	featured: boolean;
}

export interface CollectionDetailView {
	id: string;
	title: string;
	description: string;
	scope: WorkspaceScope;
	badge: string;
	status: CollectionDetailStatus;
	accessLabel: string;
	accessCopy: string;
	heroCopy: string;
	collectionHref: string;
	previewHref: string;
	actionLinks: readonly CollectionDetailActionLink[];
	metadata: readonly CollectionDetailMetadataItem[];
	requestGroups: readonly CollectionDetailGroup[];
	previewRequests: readonly CollectionDetailRequestItem[];
	relatedCollections: readonly CollectionDetailCollection[];
}

const allCollections = Array.from(
	new Map(
		workspaceStates.guest.collections
			.concat(workspaceStates.authenticated.collections)
			.map((collection) => [collection.id, collection] as const)
	).values()
);

function categorySummary(category: WorkspaceTemplateCategory) {
	switch (category) {
		case "REST basics":
			return "Core request/response examples with readable defaults.";
		case "Authentication flows":
			return "Login-style payloads that keep the request shell grounded.";
		case "CRUD examples":
			return "Write-oriented requests that show save and reuse paths.";
		case "Pagination examples":
			return "Paging-oriented lists that make metadata and repeated runs obvious.";
		case "Webhooks":
			return "Event delivery shapes that stay inside the curated app surface.";
		case "Error handling":
			return "Blocked and failure states that keep the viewer honest.";
	}
}

function collectionScopeLabel(scope: WorkspaceScope) {
	switch (scope) {
		case "guest":
			return "Guest-safe";
		case "authenticated":
			return "Authenticated";
		case "both":
			return "Shared";
	}
}

function buildRequestHref(collectionId: string, templateSlug: string) {
	return `/app?collection=${encodeURIComponent(collectionId)}&template=${encodeURIComponent(templateSlug)}`;
}

function buildPreviewHref(collectionId: string) {
	return `/app?collection=${encodeURIComponent(collectionId)}&mode=preview`;
}

export function buildCollectionDetail(collectionId: string, mode: WorkspaceMode): CollectionDetailView {
	const collection = getWorkspaceCollection(collectionId);

	if (!collection) {
		return {
			id: collectionId,
			title: "Collection unavailable",
			description: "The requested collection does not exist in the current workspace state.",
			scope: "guest",
			badge: "Missing",
			status: "missing",
			accessLabel: "Unavailable",
			accessCopy: "The current mock data does not contain this collection id.",
			heroCopy: "Return to the templates browser or open the workspace overview to continue.",
			collectionHref: "/app",
			previewHref: "/templates",
			actionLinks: [
				{ label: "Open workspace", href: "/app", variant: "default" },
				{ label: "Browse templates", href: "/templates", variant: "outline" },
			],
			metadata: [
				{ label: "Collection id", value: collectionId, detail: "No matching preview was found." },
				{ label: "Visibility", value: "Missing", detail: "This record is not seeded in the current mock state." },
			],
			requestGroups: [],
			previewRequests: [],
			relatedCollections: [],
		};
	}

	const requests: CollectionDetailRequestItem[] = [];

	for (const slug of collection.templateSlugs) {
		const template = workspaceTemplates.find((item) => item.slug === slug);
		if (!template) {
			continue;
		}

		requests.push({
			slug: template.slug,
			title: template.title,
			category: template.category,
			method: template.request.method as WorkspaceRequestMethod,
			endpoint: template.request.url,
			summary: template.summary,
			description: template.description,
			safeOverrides: template.safeOverrides,
			responseStatus: template.request.responseStatus,
			responseStatusText: template.request.responseStatusText,
			responseTimeMs: template.request.responseTimeMs,
			responseSizeLabel: template.request.responseSizeLabel,
			responseContentType: template.request.responseContentType,
			launchHref: buildRequestHref(collection.id, template.slug),
			previewHref: buildPreviewHref(collection.id),
		});
	}

	const groupMap = new Map<string, CollectionDetailGroup>();

	for (const request of requests) {
		const group = groupMap.get(request.category) ?? {
			id: request.category.toLowerCase().replace(/\s+/g, "-"),
			label: request.category,
			summary: categorySummary(request.category),
			requestCount: 0,
			requests: [],
		};

		group.requests.push(request);
		group.requestCount = group.requests.length;
		groupMap.set(request.category, group);
	}

	const requestGroups = Array.from(groupMap.values()).sort((a, b) => a.label.localeCompare(b.label));
	const previewRequests = requests.slice(0, 3);
	const visibleCollections = mode === "authenticated" ? allCollections : workspaceStates.guest.collections;
	const relatedCollections = visibleCollections
		.filter((item) => item.id !== collection.id)
		.map((item) => ({
			...item,
		}))
		.sort((a, b) => {
			const scoreA = a.templateSlugs.filter((slug) => collection.templateSlugs.includes(slug)).length;
			const scoreB = b.templateSlugs.filter((slug) => collection.templateSlugs.includes(slug)).length;
			if (scoreA !== scoreB) return scoreB - scoreA;
			return a.title.localeCompare(b.title);
		})
		.slice(0, 3);

	const requestCount = requests.length;
	const authenticated = mode === "authenticated";
	const locked = !authenticated && collection.scope === "authenticated";
	const status: CollectionDetailStatus = collection.scope === "guest" && !authenticated ? "locked" : authenticated ? "ready" : "locked";
	const accessLabel = authenticated ? "Authenticated access" : collectionScopeLabel(collection.scope);
	const accessCopy = authenticated
		? "This collection can be opened, saved, and reused from the same workspace shell."
		: collection.scope === "authenticated"
			? "Sign in to unlock the full collection detail, including the request list and launch actions."
			: "Guest users can preview the collection shape, but execution and persistence stay constrained.";

	const metadata: CollectionDetailMetadataItem[] = [
		{ label: "Scope", value: collectionScopeLabel(collection.scope), detail: "Aligned to the route contract in the docs." },
		{ label: "Templates", value: `${requestCount}`, detail: "Derived from seeded workspace templates." },
		{ label: "Groups", value: `${requestGroups.length}`, detail: "Grouped by template category." },
		{ label: "Visibility", value: collection.badge, detail: collection.featured ? "Featured in the workspace state." : "Shown as a reusable collection preview." },
	];

	if (authenticated) {
		metadata.push(
			{ label: "Mode", value: "Signed in", detail: "Authenticated users can keep using the same shell." },
			{ label: "Execution", value: "Ready", detail: "Open the collection in /app to keep working." },
		);
	} else {
		metadata.push(
			{ label: "Mode", value: "Guest", detail: "The route remains visible but constrained." },
			{ label: "Execution", value: locked ? "Locked" : "Preview only", detail: "Signing in is required for durable execution." },
		);
	}

	return {
		id: collection.id,
		title: collection.title,
		description: collection.description,
		scope: collection.scope,
		badge: collection.badge,
		status,
		accessLabel,
		accessCopy,
		heroCopy: authenticated
			? "Open the same collection shell, review the grouped requests, and keep moving into the workspace."
			: "The page shows the seeded collection structure, but guest mode keeps the execution surface intentionally narrow.",
		collectionHref: `/app?collection=${encodeURIComponent(collection.id)}`,
		previewHref: buildPreviewHref(collection.id),
		actionLinks: authenticated
			? [
					{ label: "Open in workspace", href: `/app?collection=${encodeURIComponent(collection.id)}`, variant: "default" },
					{ label: "Preview workspace", href: buildPreviewHref(collection.id), variant: "outline" },
				]
			: [
					{ label: "Open workspace", href: `/app?collection=${encodeURIComponent(collection.id)}`, variant: "default" },
					{ label: "Browse templates", href: "/templates", variant: "outline" },
				],
		metadata,
		requestGroups,
		previewRequests,
		relatedCollections,
	};
}
