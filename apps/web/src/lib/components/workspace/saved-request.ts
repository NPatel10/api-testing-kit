import type { BackendSavedRequest } from "$lib/server/backend";

import {
	createDefaultRequestDraft,
	type RequestAuthScheme,
	type RequestBody,
	type RequestBodyMode,
	type RequestBuilderDraft,
	type RequestBuilderMode,
	type RequestMethod,
	type RequestRow,
} from "./request-builder";

function normalizeMethod(value: string): RequestMethod {
	switch (value.toUpperCase()) {
		case "POST":
		case "PUT":
		case "PATCH":
		case "DELETE":
			return value.toUpperCase() as RequestMethod;
		default:
			return "GET";
	}
}

function normalizeAuthScheme(value: string | undefined): RequestAuthScheme {
	switch ((value ?? "").toLowerCase()) {
		case "bearer":
			return "bearer";
		case "basic":
			return "basic";
		default:
			return "none";
	}
}

function normalizeBodyMode(value: string | undefined): RequestBodyMode {
	switch ((value ?? "").toLowerCase()) {
		case "raw":
			return "raw";
		case "form":
		case "form-urlencoded":
		case "form_urlencoded":
			return "form";
		default:
			return "json";
	}
}

function toText(value: unknown) {
	if (typeof value === "string") {
		return value;
	}

	if (typeof value === "number" || typeof value === "boolean") {
		return String(value);
	}

	if (value === null || value === undefined) {
		return "";
	}

	try {
		return JSON.stringify(value);
	} catch {
		return "";
	}
}

function formatJsonLike(value: unknown) {
	if (typeof value === "string") {
		const trimmed = value.trim();
		if (!trimmed) {
			return "";
		}

		try {
			return JSON.stringify(JSON.parse(trimmed), null, 2);
		} catch {
			return value;
		}
	}

	if (value === null || value === undefined) {
		return "";
	}

	try {
		return JSON.stringify(value, null, 2);
	} catch {
		return toText(value);
	}
}

function normalizeRows(value: unknown): RequestRow[] {
	if (Array.isArray(value)) {
		return value.flatMap((item) => {
			if (!item || typeof item !== "object") {
				return [];
			}

			const record = item as Record<string, unknown>;
			const key = toText(record.key ?? record.name);
			const rowValue = toText(record.value);
			const enabled = typeof record.enabled === "boolean" ? record.enabled : true;

			if (!key && !rowValue) {
				return [];
			}

			return [
				{
					key,
					value: rowValue,
					enabled,
				},
			];
		});
	}

	if (value && typeof value === "object") {
		return Object.entries(value as Record<string, unknown>).map(([key, rowValue]) => ({
			key,
			value: toText(rowValue),
			enabled: true,
		}));
	}

	return [];
}

function extractHeaderContentType(headers: RequestRow[], fallback: string) {
	const contentType = headers.find((header) => header.key.toLowerCase() === "content-type")?.value.trim();
	return contentType || fallback;
}

function readAuthConfig(value: unknown) {
	if (!value || typeof value !== "object") {
		return {
			token: "",
			username: "",
			password: "",
		};
	}

	const record = value as Record<string, unknown>;
	return {
		token: toText(record.token),
		username: toText(record.username),
		password: toText(record.password),
	};
}

function buildBody(mode: RequestBodyMode, headers: RequestRow[], value: unknown): RequestBody {
	const fallbackContentType =
		mode === "form"
			? "application/x-www-form-urlencoded"
			: mode === "raw"
				? "text/plain"
				: "application/json";
	const contentType = extractHeaderContentType(headers, fallbackContentType);

	if (mode === "form") {
		const source =
			value && typeof value === "object"
				? (value as Record<string, unknown>).formFields ?? (value as Record<string, unknown>).formRows ?? value
				: value;

		return {
			mode,
			value: "",
			formRows: normalizeRows(source),
			contentType,
		};
	}

	if (value && typeof value === "object") {
		const record = value as Record<string, unknown>;
		const raw = record.raw ?? record.value ?? record.body ?? record.payload;

		return {
			mode,
			value: mode === "json" ? formatJsonLike(raw ?? value) : toText(raw ?? ""),
			formRows: [],
			contentType,
		};
	}

	return {
		mode,
		value: mode === "json" ? formatJsonLike(value) : toText(value),
		formRows: [],
		contentType,
	};
}

export function createSavedRequestDraft(
	mode: RequestBuilderMode,
	savedRequest: BackendSavedRequest,
): RequestBuilderDraft {
	const draft = createDefaultRequestDraft(mode);
	const headers = normalizeRows(savedRequest.headers);
	const auth = readAuthConfig(savedRequest.authConfig);
	const bodyMode = normalizeBodyMode(savedRequest.bodyMode);

	return {
		method: normalizeMethod(savedRequest.method),
		url: savedRequest.url,
		queryParams: normalizeRows(savedRequest.queryParams),
		headers,
		auth: {
			scheme: normalizeAuthScheme(savedRequest.authScheme),
			token: auth.token,
			username: auth.username,
			password: auth.password,
		},
		body: buildBody(bodyMode, headers, savedRequest.bodyConfig),
	};
}
