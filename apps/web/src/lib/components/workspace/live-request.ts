import type { RequestBuilderDraft } from "./request-builder";
import type { ResponseHeader, ResponseViewerError } from "./response-viewer";

export interface LiveRunPayload {
	method: string;
	url: string;
	queryParams: Array<{ name: string; value: string; enabled: boolean }>;
	headers: Array<{ name: string; value: string; enabled: boolean }>;
	auth: {
		scheme: string;
		username: string;
		password: string;
		token: string;
	};
	body: {
		mode: string;
		raw: string;
		formFields: Array<{ name: string; value: string; enabled: boolean }>;
	};
}

export interface LiveRunApiResponse {
	status?: string;
	method?: string;
	url?: string;
	finalUrl?: string;
	responseStatus?: number;
	responseHeaders?: Record<string, string[]>;
	responseBodyPreview?: string;
	responseJson?: unknown;
	responseSizeBytes?: number;
	responseTimeMs?: number;
	responseContentType?: string;
	redirectCount?: number;
	blockedReason?: string;
	errorCode?: string;
	errorMessage?: string;
	truncated?: boolean;
	error?: {
		code?: string;
		message?: string;
	};
}

export interface LiveRunViewerState {
	status?: number;
	statusText?: string;
	duration?: number;
	size?: number | string;
	contentType?: string;
	headers: ResponseHeader[];
	prettyBody: string;
	rawBody: string;
	error: ResponseViewerError | null;
}

const HTTP_STATUS_TEXT: Record<number, string> = {
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",
	301: "Moved Permanently",
	302: "Found",
	304: "Not Modified",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	409: "Conflict",
	413: "Payload Too Large",
	429: "Too Many Requests",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
};

export function buildLiveRunPayload(draft: RequestBuilderDraft): LiveRunPayload {
	return {
		method: draft.method,
		url: draft.url,
		queryParams: draft.queryParams.map((row) => ({ name: row.key, value: row.value, enabled: row.enabled })),
		headers: draft.headers.map((row) => ({ name: row.key, value: row.value, enabled: row.enabled })),
		auth: {
			scheme: draft.auth.scheme,
			username: draft.auth.username,
			password: draft.auth.password,
			token: draft.auth.token,
		},
		body: {
			mode: draft.body.mode === "form" ? "form_urlencoded" : draft.body.mode,
			raw: draft.body.value,
			formFields: draft.body.formRows.map((row) => ({ name: row.key, value: row.value, enabled: row.enabled })),
		},
	};
}

export function createEmptyViewerState(): LiveRunViewerState {
	return {
		headers: [],
		prettyBody: "",
		rawBody: "",
		error: null,
	};
}

export function createPreviewViewerState(snapshot: {
	responseStatus: number;
	responseStatusText: string;
	responseTimeMs: number;
	responseSizeLabel: string;
	responseContentType: string;
	responseBody: string;
}): LiveRunViewerState {
	return {
		status: snapshot.responseStatus,
		statusText: snapshot.responseStatusText,
		duration: snapshot.responseTimeMs,
		size: snapshot.responseSizeLabel,
		contentType: snapshot.responseContentType,
		headers: [
			{ key: "content-type", value: snapshot.responseContentType },
			{ key: "x-preview-state", value: "demo" },
			{ key: "x-preview-size", value: snapshot.responseSizeLabel },
		],
		prettyBody: snapshot.responseBody,
		rawBody: snapshot.responseBody,
		error: null,
	};
}

export async function readLiveRunViewerState(response: Response): Promise<LiveRunViewerState> {
	const text = await response.text();
	const payload = safeParseJson(text);

	if (!response.ok) {
		return buildErrorState(response.status, payload ?? text);
	}

	if (!isLiveRunApiResponse(payload)) {
		return buildErrorState(response.status, payload ?? text);
	}

	if (payload.status && payload.status !== "succeeded") {
		return buildErrorState(response.status, payload);
	}

	const status = payload.responseStatus ?? response.status;
	const statusText = HTTP_STATUS_TEXT[status] ?? response.statusText ?? "";
	const prettyBody = formatResponseBody(payload.responseJson, payload.responseBodyPreview);
	const rawBody = payload.responseBodyPreview ?? prettyBody;

	return {
		status,
		statusText,
		duration: payload.responseTimeMs,
		size: payload.responseSizeBytes,
		contentType: payload.responseContentType,
		headers: normalizeHeaders(payload.responseHeaders),
		prettyBody,
		rawBody,
		error: null,
	};
}

export function createRequestFailureState(message: string, code = "request_failed"): LiveRunViewerState {
	return {
		...createEmptyViewerState(),
		error: {
			title: "Request failed",
			message,
			code,
		},
	};
}

function buildErrorState(status: number, payload: unknown): LiveRunViewerState {
	const { code, message } = extractError(payload);
	const title = inferErrorTitle(code, status);
	const details = status ? `HTTP ${status}` : undefined;

	return {
		...createEmptyViewerState(),
		error: {
			title,
			message: message || `The backend returned HTTP ${status}.`,
			code: code || `http_${status}`,
			details,
		},
	};
}

function formatResponseBody(responseJson: unknown, fallbackBody?: string): string {
	if (typeof responseJson === "string") {
		return responseJson;
	}

	if (responseJson !== undefined && responseJson !== null) {
		try {
			return JSON.stringify(responseJson, null, 2);
		} catch {
			return String(responseJson);
		}
	}

	return fallbackBody ?? "";
}

function inferErrorTitle(code: string | undefined, status: number): string {
	switch (code) {
		case "blocked_target":
			return "Blocked target";
		case "blocked_host":
			return "Blocked host";
		case "blocked_ip":
			return "Blocked IP";
		case "invalid_request":
		case "invalid_run":
			return "Invalid request";
		case "guest_rate_limited":
			return "Rate limited";
		case "unauthorized":
			return "Authentication required";
		case "upstream_request_failed":
		case "guest_request_failed":
			return "Upstream failure";
		default:
			if (status >= 500) {
				return "Server error";
			}
			if (status >= 400) {
				return "Request rejected";
			}
			return "Request failed";
	}
}

function extractError(payload: unknown): { code?: string; message?: string } {
	if (!payload || typeof payload !== "object") {
		return {};
	}

	const record = payload as Record<string, unknown>;
	const error = record.error;
	if (error && typeof error === "object") {
		const errorRecord = error as Record<string, unknown>;
		return {
			code: typeof errorRecord.code === "string" ? errorRecord.code : undefined,
			message: typeof errorRecord.message === "string" ? errorRecord.message : undefined,
		};
	}

	return {
		code: typeof record.errorCode === "string" ? record.errorCode : undefined,
		message: typeof record.errorMessage === "string" ? record.errorMessage : undefined,
	};
}

function normalizeHeaders(headers?: Record<string, string[]>): ResponseHeader[] {
	if (!headers) {
		return [];
	}

	return Object.entries(headers).flatMap(([key, values]) =>
		values.map((value) => ({ key, value })),
	);
}

function safeParseJson(value: string): unknown | null {
	if (!value.trim()) {
		return null;
	}

	try {
		return JSON.parse(value) as unknown;
	} catch {
		return null;
	}
}

function isLiveRunApiResponse(value: unknown): value is LiveRunApiResponse {
	if (!value || typeof value !== "object") {
		return false;
	}

	const record = value as Record<string, unknown>;
	return typeof record.responseStatus === "number" || typeof record.status === "string";
}
