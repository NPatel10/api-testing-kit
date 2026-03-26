import {
	getWorkspaceCollection,
	workspaceStates,
	workspaceTemplates,
	type WorkspaceMode,
	type WorkspaceRequestMethod,
	type WorkspaceScope,
} from "$lib/mocks/workspace-state";
import type { BackendCollection, BackendSavedRequest } from "$lib/server/backend";

export type CollectionDetailStatus = "missing" | "locked" | "ready" | "unavailable";

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
	id: string;
	title: string;
	category: string;
	method: string;
	endpoint: string;
	summary: string;
	description: string;
	safeOverrides: readonly string[];
	responseLabel: string;
	responseDetail: string;
	durationLabel: string;
	durationDetail: string;
	sizeLabel: string;
	sizeDetail: string;
	launchHref: string;
	secondaryHref?: string;
	secondaryLabel?: string;
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
	scope: string;
	requestCount: number;
	badge: string;
}

export interface CollectionDetailView {
	id: string;
	title: string;
	description: string;
	scope: string;
	badge: string;
	status: CollectionDetailStatus;
	source: "preview" | "live";
	accessLabel: string;
	accessCopy: string;
	heroCopy: string;
	collectionHref: string;
	previewHref?: string;
	actionLinks: readonly CollectionDetailActionLink[];
	metadata: readonly CollectionDetailMetadataItem[];
	requestGroups: readonly CollectionDetailGroup[];
	previewRequests: readonly CollectionDetailRequestItem[];
	relatedCollections: readonly CollectionDetailCollection[];
}

const allPreviewCollections = Array.from(
	new Map(
		workspaceStates.guest.collections
			.concat(workspaceStates.authenticated.collections)
			.map((collection) => [collection.id, collection] as const),
	).values(),
);

function buildRequestHref(collectionId: string, templateSlug: string) {
	return `/app?collection=${encodeURIComponent(collectionId)}&template=${encodeURIComponent(templateSlug)}`;
}

function buildPreviewHref(collectionId: string) {
	return `/app?collection=${encodeURIComponent(collectionId)}&mode=preview`;
}

function toGroupId(label: string) {
	return label.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
}

function formatBytes(value: number | null | undefined) {
	if (value === null || value === undefined || Number.isNaN(value)) {
		return "n/a";
	}

	if (value < 1024) {
		return `${value} B`;
	}

	if (value < 1024 * 1024) {
		return `${(value / 1024).toFixed(1)} KB`;
	}

	return `${(value / (1024 * 1024)).toFixed(1)} MB`;
}

function formatDateLabel(value: string) {
	const timestamp = new Date(value);
	if (Number.isNaN(timestamp.getTime())) {
		return "Recently updated";
	}

	return new Intl.DateTimeFormat("en-US", {
		month: "short",
		day: "numeric",
		hour: "numeric",
		minute: "2-digit",
	}).format(timestamp);
}

function collectionScopeLabel(scope: WorkspaceScope | BackendCollection["visibility"]) {
	switch (scope) {
		case "guest":
			return "Guest-safe";
		case "authenticated":
			return "Authenticated";
		case "both":
			return "Shared";
		case "private":
			return "Private";
		case "shared_readonly":
			return "Shared read-only";
		case "internal":
			return "Internal";
	}
}

function buildMissingCollectionDetail(
	collectionId: string,
	source: "preview" | "live",
	message: string,
): CollectionDetailView {
	return {
		id: collectionId,
		title: "Collection unavailable",
		description: message,
		scope: source === "preview" ? "Guest-safe" : "Authenticated",
		badge: "Unavailable",
		status: "missing",
		source,
		accessLabel: "Unavailable",
		accessCopy: message,
		heroCopy: "Return to the shared workspace or browse templates to keep moving.",
		collectionHref: "/app",
		actionLinks: [
			{ label: "Open workspace", href: "/app", variant: "default" },
			{ label: "Browse templates", href: "/templates", variant: "outline" },
		],
		metadata: [
			{ label: "Collection id", value: collectionId, detail: "No collection matched the requested id." },
			{ label: "Source", value: source === "preview" ? "Preview" : "Live", detail: "The route could not resolve this collection." },
		],
		requestGroups: [],
		previewRequests: [],
		relatedCollections: [],
	};
}

function buildUnavailableCollectionDetail(collectionId: string): CollectionDetailView {
	return {
		id: collectionId,
		title: "Collection temporarily unavailable",
		description: "The collection detail API did not return a usable payload for this signed-in view.",
		scope: "Authenticated",
		badge: "Unavailable",
		status: "unavailable",
		source: "live",
		accessLabel: "Unavailable",
		accessCopy: "The live collection detail could not be loaded right now.",
		heroCopy: "Try the workspace again after the backend recovers, or open another collection.",
		collectionHref: "/app",
		actionLinks: [
			{ label: "Open workspace", href: "/app", variant: "default" },
			{ label: "Go to history", href: "/app/history", variant: "outline" },
		],
		metadata: [
			{ label: "Collection id", value: collectionId, detail: "The signed-in route attempted a live lookup." },
			{ label: "Mode", value: "Signed in", detail: "No preview fallback is shown for authenticated persistence." },
		],
		requestGroups: [],
		previewRequests: [],
		relatedCollections: [],
	};
}

function inferLiveCategory(request: BackendSavedRequest) {
	const method = request.method.toUpperCase();
	const url = request.url.toLowerCase();

	if ((request.authScheme ?? "").toLowerCase() !== "none" && (request.authScheme ?? "").trim() !== "") {
		return "Authentication";
	}

	if (url.includes("webhook")) {
		return "Webhooks";
	}

	if (url.includes("page=") || url.includes("limit=") || url.includes("cursor=")) {
		return "Pagination";
	}

	if (["POST", "PUT", "PATCH", "DELETE"].includes(method)) {
		return "Write flows";
	}

	return "Saved requests";
}

function parseExampleResponse(value: unknown) {
	if (!value || typeof value !== "object") {
		return null;
	}

	const record = value as Record<string, unknown>;
	const responseStatus =
		typeof record.responseStatus === "number"
			? record.responseStatus
			: typeof record.status === "number"
				? record.status
				: undefined;
	const responseText =
		typeof record.responseStatusText === "string"
			? record.responseStatusText
			: typeof record.statusText === "string"
				? record.statusText
				: undefined;
	const responseSize =
		typeof record.responseSizeBytes === "number"
			? record.responseSizeBytes
			: typeof record.sizeBytes === "number"
				? record.sizeBytes
				: undefined;
	const responseSizeLabel =
		typeof record.responseSizeLabel === "string"
			? record.responseSizeLabel
			: typeof record.sizeLabel === "string"
				? record.sizeLabel
				: undefined;
	const responseTime =
		typeof record.responseTimeMs === "number"
			? record.responseTimeMs
			: typeof record.durationMs === "number"
				? record.durationMs
				: undefined;
	const responseContentType =
		typeof record.responseContentType === "string"
			? record.responseContentType
			: typeof record.contentType === "string"
				? record.contentType
				: undefined;

	return {
		responseLabel:
			typeof responseStatus === "number"
				? `${responseStatus}${responseText ? ` ${responseText}` : ""}`
				: null,
		responseDetail: responseContentType?.trim() || null,
		durationLabel: typeof responseTime === "number" ? `${responseTime} ms` : null,
		sizeLabel: responseSizeLabel || formatBytes(responseSize),
	};
}

function buildPreviewRequestItem(collectionId: string, templateSlug: string): CollectionDetailRequestItem | null {
	const template = workspaceTemplates.find((item) => item.slug === templateSlug);
	if (!template) {
		return null;
	}

	return {
		id: template.slug,
		title: template.title,
		category: template.category,
		method: template.request.method,
		endpoint: template.request.url,
		summary: template.summary,
		description: template.description,
		safeOverrides: template.safeOverrides,
		responseLabel: `${template.request.responseStatus} ${template.request.responseStatusText}`,
		responseDetail: template.request.responseContentType,
		durationLabel: `${template.request.responseTimeMs} ms`,
		durationDetail: "Measured from the current workspace preview.",
		sizeLabel: template.request.responseSizeLabel,
		sizeDetail: "Preview payload size.",
		launchHref: buildRequestHref(collectionId, template.slug),
		secondaryHref: buildPreviewHref(collectionId),
		secondaryLabel: "Preview",
	};
}

function buildLiveRequestItem(request: BackendSavedRequest): CollectionDetailRequestItem {
	const example = parseExampleResponse(request.exampleResponse);

	return {
		id: request.id,
		title: request.name,
		category: inferLiveCategory(request),
		method: request.method.toUpperCase(),
		endpoint: request.url,
		summary:
			request.description?.trim() ||
			"Persisted request definition ready to reopen in the shared workspace.",
		description:
			request.description?.trim() ||
			"Saved request metadata comes from the authenticated collections API rather than seeded template copy.",
		safeOverrides: [],
		responseLabel: example?.responseLabel || "No stored response preview",
		responseDetail: example?.responseDetail || "Send the request in /app to capture a fresh live response.",
		durationLabel: example?.durationLabel || "Captured on send",
		durationDetail: example?.durationLabel
			? "Stored with the saved request preview."
			: "Duration is recorded after the next live execution.",
		sizeLabel: example?.sizeLabel || "Live on rerun",
		sizeDetail: example?.sizeLabel
			? "Stored response preview size."
			: "Payload size appears after the request runs.",
		launchHref: `/app?request=${encodeURIComponent(request.id)}`,
	};
}

function buildGroupedRequests(items: CollectionDetailRequestItem[]) {
	const groupMap = new Map<string, CollectionDetailGroup>();

	for (const item of items) {
		const label = item.category || "Saved requests";
		const group = groupMap.get(label) ?? {
			id: toGroupId(label),
			label,
			summary:
				label === "Saved requests"
					? "Persisted request definitions in collection order."
					: `Saved requests grouped under ${label.toLowerCase()}.`,
			requestCount: 0,
			requests: [],
		};

		group.requests.push(item);
		group.requestCount = group.requests.length;
		groupMap.set(label, group);
	}

	return Array.from(groupMap.values()).sort((left, right) => left.label.localeCompare(right.label));
}

export function buildCollectionDetail(collectionId: string, mode: WorkspaceMode): CollectionDetailView {
	const collection = getWorkspaceCollection(collectionId);

	if (!collection) {
		return buildMissingCollectionDetail(collectionId, "preview", "The requested collection does not exist in the current preview set.");
	}

	const requests = collection.templateSlugs
		.map((slug) => buildPreviewRequestItem(collection.id, slug))
		.filter((item): item is CollectionDetailRequestItem => Boolean(item));
	const requestGroups = buildGroupedRequests(requests);
	const visibleCollections = mode === "authenticated" ? allPreviewCollections : workspaceStates.guest.collections;
	const relatedCollections = visibleCollections
		.filter((item) => item.id !== collection.id)
		.sort((left, right) => left.title.localeCompare(right.title))
		.slice(0, 3)
		.map((item) => ({
			id: item.id,
			title: item.title,
			description: item.description,
			scope: collectionScopeLabel(item.scope),
			requestCount: item.requestCount,
			badge: item.badge,
		}));
	const authenticated = mode === "authenticated";

	return {
		id: collection.id,
		title: collection.title,
		description: collection.description,
		scope: collectionScopeLabel(collection.scope),
		badge: collection.badge,
		status: authenticated ? "ready" : "locked",
		source: "preview",
		accessLabel: authenticated ? "Authenticated access" : collectionScopeLabel(collection.scope),
		accessCopy: authenticated
			? "This preview collection can be opened from the same signed-in workspace shell."
			: "Sign in to unlock persistence and saved execution while keeping the same collection route visible.",
		heroCopy: authenticated
			? "This route keeps the same layout while the live collection surface is loading from persistence."
			: "The route stays visible for guests, but the preview remains intentionally constrained.",
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
		metadata: [
			{ label: "Scope", value: collectionScopeLabel(collection.scope), detail: "Aligned to the shared /app contract." },
			{ label: "Requests", value: `${requests.length}`, detail: "Visible request previews in this collection." },
			{ label: "Groups", value: `${requestGroups.length}`, detail: "Grouped by preview request category." },
			{ label: "Visibility", value: collection.badge, detail: collection.featured ? "Featured in the workspace." : "Visible in the current preview set." },
			{
				label: "Mode",
				value: authenticated ? "Signed in" : "Guest",
				detail: authenticated ? "The signed-in shell keeps the same route structure." : "Guests keep the route, but not durable persistence.",
			},
		],
		requestGroups,
		previewRequests: requests.slice(0, 3),
		relatedCollections,
	};
}

export function buildLiveCollectionDetail(
	collection: BackendCollection,
	savedRequests: BackendSavedRequest[],
	allCollections: BackendCollection[],
): CollectionDetailView {
	const requests = [...savedRequests]
		.sort((left, right) => {
			if (left.sortOrder !== right.sortOrder) {
				return left.sortOrder - right.sortOrder;
			}

			return left.createdAt.localeCompare(right.createdAt);
		})
		.map(buildLiveRequestItem);
	const requestGroups = buildGroupedRequests(requests);
	const firstRequest = requests[0];
	const relatedCollections = [...allCollections]
		.filter((item) => item.id !== collection.id)
		.sort((left, right) => {
			if (left.sortOrder !== right.sortOrder) {
				return left.sortOrder - right.sortOrder;
			}

			return left.name.localeCompare(right.name);
		})
		.slice(0, 3)
		.map((item) => ({
			id: item.id,
			title: item.name,
			description: item.description?.trim() || "Persisted collection available to this signed-in account.",
			scope: collectionScopeLabel(item.visibility),
			requestCount: 0,
			badge: collectionScopeLabel(item.visibility),
		}));

	return {
		id: collection.id,
		title: collection.name,
		description: collection.description?.trim() || "Persisted collection loaded from the authenticated collections API.",
		scope: collectionScopeLabel(collection.visibility),
		badge: collectionScopeLabel(collection.visibility),
		status: "ready",
		source: "live",
		accessLabel: "Authenticated access",
		accessCopy: "This collection is loaded from persisted data, and its saved requests reopen inside the shared workspace.",
		heroCopy: "Review the saved requests in collection order, then jump back into /app on the exact persisted request you want to run.",
		collectionHref: firstRequest?.launchHref || "/app",
		actionLinks: firstRequest
			? [
					{ label: "Open first request", href: firstRequest.launchHref, variant: "default" },
					{ label: "Open history", href: "/app/history", variant: "outline" },
				]
			: [
					{ label: "Open workspace", href: "/app", variant: "default" },
					{ label: "Open history", href: "/app/history", variant: "outline" },
				],
		metadata: [
			{ label: "Visibility", value: collectionScopeLabel(collection.visibility), detail: "Returned by the live collections API." },
			{ label: "Saved requests", value: `${requests.length}`, detail: "Ordered by persisted sort order." },
			{ label: "Groups", value: `${requestGroups.length}`, detail: "Grouped from the saved request payloads." },
			{ label: "Sort order", value: `${collection.sortOrder}`, detail: "Collection order from persistence." },
			{ label: "Updated", value: formatDateLabel(collection.updatedAt), detail: "Last collection update returned by the API." },
			{ label: "Mode", value: "Signed in", detail: "Authenticated routes use live collection persistence instead of preview data." },
		],
		requestGroups,
		previewRequests: requests.slice(0, 3),
		relatedCollections,
	};
}

export function buildLiveMissingCollectionDetail(collectionId: string) {
	return buildMissingCollectionDetail(
		collectionId,
		"live",
		"The live collections API did not return a collection for this id.",
	);
}

export { buildUnavailableCollectionDetail };
