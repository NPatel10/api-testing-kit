export type RequestBuilderMode = "guest" | "authenticated";

export type RequestMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type RequestBodyMode = "json" | "raw" | "form";

export type RequestAuthScheme = "none" | "bearer" | "basic";

export type ValidationSeverity = "info" | "warning" | "danger";

export interface RequestRow {
	key: string;
	value: string;
	enabled: boolean;
}

export interface RequestAuth {
	scheme: RequestAuthScheme;
	token: string;
	username: string;
	password: string;
}

export interface RequestBody {
	mode: RequestBodyMode;
	value: string;
	formRows: RequestRow[];
	contentType: string;
}

export interface RequestBuilderDraft {
	method: RequestMethod;
	url: string;
	queryParams: RequestRow[];
	headers: RequestRow[];
	auth: RequestAuth;
	body: RequestBody;
}

export interface RequestValidationIssue {
	severity: ValidationSeverity;
	title: string;
	description: string;
	action?: string;
}

export const requestMethods: RequestMethod[] = ["GET", "POST", "PUT", "PATCH", "DELETE"];

export const requestBodyModes: Array<{
	value: RequestBodyMode;
	label: string;
	description: string;
}> = [
	{ value: "json", label: "JSON", description: "Structured body" },
	{ value: "raw", label: "Raw", description: "Plain text payload" },
	{ value: "form", label: "Form URL encoded", description: "Key/value pairs" },
];

export const requestAuthSchemes: Array<{
	value: RequestAuthScheme;
	label: string;
	description: string;
}> = [
	{ value: "none", label: "None", description: "No authorization header" },
	{ value: "bearer", label: "Bearer", description: "Token-based auth" },
	{ value: "basic", label: "Basic", description: "Username and password" },
];

export const guestAllowlistedHosts = [
	"jsonplaceholder.typicode.com",
	"api.github.com",
	"api.open-meteo.com",
	"open-meteo.com",
	"demo.api-testing-kit.local",
];

export function createRequestRow(key = "", value = "", enabled = true): RequestRow {
	return { key, value, enabled };
}

export function createDefaultRequestDraft(mode: RequestBuilderMode = "guest"): RequestBuilderDraft {
	return {
		method: "GET",
		url:
			mode === "guest"
				? "https://jsonplaceholder.typicode.com/posts/1"
				: "https://api.github.com/users/octocat",
		queryParams: [
			createRequestRow("include", "summary", true),
			createRequestRow("limit", "10", true),
		],
		headers: [
			createRequestRow("accept", "application/json", true),
			createRequestRow("x-demo-mode", mode, true),
		],
		auth: {
			scheme: "none",
			token: "demo-token",
			username: "demo",
			password: "demo",
		},
		body: {
			mode: "json",
			value: '{\n  "name": "API Testing Kit",\n  "active": true\n}',
			formRows: [
				createRequestRow("name", "API Testing Kit", true),
				createRequestRow("active", "true", true),
			],
			contentType: "application/json",
		},
	};
}

export function cloneRequestDraft(draft: RequestBuilderDraft): RequestBuilderDraft {
	return {
		method: draft.method,
		url: draft.url,
		queryParams: draft.queryParams.map((row) => ({ ...row })),
		headers: draft.headers.map((row) => ({ ...row })),
		auth: { ...draft.auth },
		body: {
			mode: draft.body.mode,
			value: draft.body.value,
			formRows: draft.body.formRows.map((row) => ({ ...row })),
			contentType: draft.body.contentType,
		},
	};
}

export function isGuestAllowedTarget(url: string): boolean {
	if (!url.trim()) {
		return false;
	}

	try {
		const parsed = new URL(url);
		return guestAllowlistedHosts.includes(parsed.hostname);
	} catch {
		return false;
	}
}

export function getRequestValidation(mode: RequestBuilderMode, draft: RequestBuilderDraft): RequestValidationIssue | null {
	if (!draft.url.trim()) {
		return {
			severity: "warning",
			title: "Add a target URL",
			description: "The runner needs a URL before it can build a request.",
		};
	}

	try {
		new URL(draft.url);
	} catch {
		return {
			severity: "danger",
			title: "Enter a valid URL",
			description: "Use a complete absolute URL, including protocol and host.",
		};
	}

	if (mode === "guest" && !isGuestAllowedTarget(draft.url)) {
		return {
			severity: "danger",
			title: "Guest mode is locked to allowlisted demos",
			description:
				"Guests can inspect the builder, but custom targets and arbitrary domains stay behind sign-in.",
			action: "Open a curated template or sign in to unlock custom execution.",
		};
	}

	if (draft.method !== "GET" && draft.body.mode === "json") {
		const trimmed = draft.body.value.trim();
		if (trimmed.length > 0) {
			try {
				JSON.parse(trimmed);
			} catch {
				return {
					severity: "warning",
					title: "JSON body needs a fix",
					description: "The payload is not valid JSON yet. Clean it up before sending.",
				};
			}
		}
	}

	return null;
}

export function formatMethodTone(method: RequestMethod): string {
	switch (method) {
		case "GET":
			return "border-emerald-200 bg-emerald-50 text-emerald-900";
		case "POST":
			return "border-teal-200 bg-teal-50 text-teal-900";
		case "PUT":
			return "border-amber-200 bg-amber-50 text-amber-900";
		case "PATCH":
			return "border-lime-200 bg-lime-50 text-lime-900";
		case "DELETE":
			return "border-rose-200 bg-rose-50 text-rose-900";
	}
}

export function formatValidationTone(severity: ValidationSeverity): string {
	switch (severity) {
		case "info":
			return "border-border/70 bg-panel-soft text-text-body";
		case "warning":
			return "border-amber-200 bg-amber-50 text-amber-900";
		case "danger":
			return "border-rose-200 bg-rose-50 text-rose-900";
	}
}

