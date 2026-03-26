import { env as privateEnv } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import type { RequestEvent } from "@sveltejs/kit";

export interface BackendCollection {
	id: string;
	ownerUserId?: string | null;
	name: string;
	slug?: string | null;
	description?: string;
	visibility: "private" | "shared_readonly" | "internal";
	color?: string;
	sortOrder: number;
	sharedToken?: string | null;
	metadata?: unknown;
	createdAt: string;
	updatedAt: string;
	deletedAt?: string | null;
}

export interface BackendSavedRequest {
	id: string;
	collectionId?: string | null;
	ownerUserId?: string | null;
	name: string;
	description?: string;
	method: string;
	url: string;
	sortOrder: number;
	queryParams?: unknown;
	headers?: unknown;
	authScheme?: string;
	authConfig?: unknown;
	bodyMode?: string;
	bodyConfig?: unknown;
	exampleResponse?: unknown;
	metadata?: unknown;
	createdAt: string;
	updatedAt: string;
	deletedAt?: string | null;
}

export interface BackendRunRecord {
	id: string;
	userId?: string | null;
	collectionId?: string | null;
	savedRequestId?: string | null;
	source: string;
	status: string;
	method: string;
	url: string;
	finalUrl?: string | null;
	targetHost?: string;
	requestHeaders?: unknown;
	requestQueryParams?: unknown;
	requestAuth?: unknown;
	requestBody?: unknown;
	responseStatus?: number | null;
	responseHeaders?: unknown;
	responseBodyPreview?: string;
	responseSizeBytes?: number | null;
	responseTimeMs?: number | null;
	responseContentType?: string;
	redirectCount: number;
	blockedReason?: string;
	errorCode?: string;
	errorMessage?: string;
	startedAt?: string | null;
	completedAt?: string | null;
	createdAt: string;
	metadata?: unknown;
}

export interface BackendHistoryListResponse {
	history: BackendRunRecord[];
	pagination: {
		page: number;
		limit: number;
		hasMore: boolean;
	};
}

export interface BackendCollectionsListResponse {
	collections: BackendCollection[];
}

export interface BackendCollectionDetailResponse {
	collection: BackendCollection;
	savedRequests: BackendSavedRequest[];
}

function normalizeBaseUrl(value: string | undefined) {
	return (value ?? "http://localhost:8080").replace(/\/+$/, "");
}

export function getBackendBaseUrl() {
	return normalizeBaseUrl(
		privateEnv.INTERNAL_API_BASE_URL || privateEnv.API_BASE_URL || publicEnv.PUBLIC_API_BASE_URL,
	);
}

interface ReadBackendJsonOptions {
	cookie?: string | null;
	headers?: HeadersInit;
}

export async function readBackendJson<T>(
	fetchFn: typeof fetch,
	path: string,
	options: ReadBackendJsonOptions = {},
) {
	const headers = new Headers(options.headers);
	headers.set("accept", "application/json");

	if (options.cookie) {
		headers.set("cookie", options.cookie);
	}

	const response = await fetchFn(`${getBackendBaseUrl()}${path}`, {
		headers,
		cache: "no-store",
	});

	return {
		response,
		data: await safeReadJson<T>(response),
	};
}

export async function proxyBackendJson(event: RequestEvent, path: string, forwardCookie = false) {
	const headers = new Headers();
	const contentType = event.request.headers.get("content-type");
	if (contentType) {
		headers.set("content-type", contentType);
	}
	headers.set("accept", "application/json");

	if (forwardCookie) {
		const cookie = event.request.headers.get("cookie");
		if (cookie) {
			headers.set("cookie", cookie);
		}
	}

	const body = event.request.method === "GET" || event.request.method === "HEAD" ? undefined : await event.request.text();
	const response = await event.fetch(`${getBackendBaseUrl()}${path}`, {
		method: event.request.method,
		headers,
		body,
	});

	return new Response(await response.text(), {
		status: response.status,
		statusText: response.statusText,
		headers: {
			"content-type": response.headers.get("content-type") ?? "application/json",
			"cache-control": "no-store",
		},
	});
}

async function safeReadJson<T>(response: Response): Promise<T | null> {
	const text = await response.text();
	if (!text.trim()) {
		return null;
	}

	try {
		return JSON.parse(text) as T;
	} catch {
		return null;
	}
}
