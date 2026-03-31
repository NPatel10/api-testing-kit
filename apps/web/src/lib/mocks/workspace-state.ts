export const workspaceModes = ["guest", "authenticated"] as const;

export type WorkspaceMode = (typeof workspaceModes)[number];

export const workspaceTemplateCategories = [
	"REST basics",
	"Authentication flows",
	"CRUD examples",
	"Pagination examples",
	"Webhooks",
	"Error handling",
] as const;

export type WorkspaceTemplateCategory = (typeof workspaceTemplateCategories)[number];

export const workspaceSurfaces = [
	"toolbar",
	"sidebar",
	"request-builder",
	"response-viewer",
	"utility-drawer",
] as const;

export type WorkspaceSurface = (typeof workspaceSurfaces)[number];

export type WorkspacePromptKind = "sign-in" | "upgrade" | "learn-more";

export type WorkspaceRequestMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type WorkspaceBodyMode = "none" | "json" | "raw" | "form-urlencoded";

export type WorkspaceTone = "neutral" | "positive" | "warning" | "danger";

export type WorkspaceScope = WorkspaceMode | "both";

export interface WorkspacePromptAction {
	kind: WorkspacePromptKind;
	label: string;
}

export interface WorkspacePrompt {
	id: string;
	kind: WorkspacePromptKind;
	title: string;
	body: string;
	tone: WorkspaceTone;
	action: WorkspacePromptAction;
}

export type WorkspaceControlId =
	| "template-browser"
	| "method-selector"
	| "url-bar"
	| "query-params"
	| "headers-editor"
	| "auth-selector"
	| "body-editor"
	| "send-request"
	| "save-request"
	| "history-rail"
	| "environment-variables"
	| "custom-url"
	| "snippet-drawer"
	| "advanced-tools";

export interface WorkspaceControl {
	id: WorkspaceControlId;
	label: string;
	surface: WorkspaceSurface;
	description: string;
	visible: boolean;
	editable: boolean;
	locked: boolean;
	availability: "available" | "preview" | "locked";
	constraint?: string;
	promptId?: string;
}

export interface WorkspaceLockedAction {
	id: string;
	label: string;
	description: string;
	surface: WorkspaceSurface;
	promptId: string;
	previewCopy: string;
}

export interface WorkspaceRequestPreview {
	method: WorkspaceRequestMethod;
	url: string;
	bodyMode: WorkspaceBodyMode;
	query: readonly { key: string; value: string }[];
	headers: readonly { key: string; value: string }[];
	responseStatus: number;
	responseStatusText: string;
	responseTimeMs: number;
	responseSizeLabel: string;
	responseContentType: string;
	responseBody: string;
}

export interface WorkspaceTemplate {
	slug: string;
	category: WorkspaceTemplateCategory;
	title: string;
	summary: string;
	description: string;
	availability: WorkspaceScope;
	allowlistedTarget: string;
	safeOverrides: readonly string[];
	tags: readonly string[];
	request: WorkspaceRequestPreview;
}

export interface WorkspaceTemplateGroup {
	id: string;
	label: WorkspaceTemplateCategory;
	summary: string;
	templateSlugs: readonly string[];
}

export interface WorkspaceCollectionPreview {
	id: string;
	title: string;
	description: string;
	scope: WorkspaceScope;
	requestCount: number;
	templateSlugs: readonly string[];
	badge: string;
	featured: boolean;
}

export interface WorkspaceHistoryEntry {
	id: string;
	title: string;
	scope: WorkspaceScope;
	method: WorkspaceRequestMethod;
	target: string;
	statusCode: number;
	statusText: string;
	durationMs: number;
	responseSizeLabel: string;
	contentType: string;
	timestampLabel: string;
	source: "demo" | "persistent";
	outcome: "success" | "blocked" | "error";
}

export interface WorkspaceMetric {
	id: string;
	label: string;
	value: string;
	detail: string;
	scope: WorkspaceScope;
	tone: WorkspaceTone;
}

export interface WorkspaceQuota {
	id: string;
	label: string;
	scope: WorkspaceScope;
	limit: number;
	limitLabel: string;
	used: number;
	usedLabel: string;
	remaining: number;
	remainingLabel: string;
	windowLabel: string;
	note: string;
}

export interface WorkspaceState {
	mode: WorkspaceMode;
	title: string;
	subtitle: string;
	accessSummary: string;
	prompts: readonly WorkspacePrompt[];
	controls: readonly WorkspaceControl[];
	lockedActions: readonly WorkspaceLockedAction[];
	templateGroups: readonly WorkspaceTemplateGroup[];
	templates: readonly WorkspaceTemplate[];
	collections: readonly WorkspaceCollectionPreview[];
	history: readonly WorkspaceHistoryEntry[];
	metrics: readonly WorkspaceMetric[];
	quotas: readonly WorkspaceQuota[];
}

const promptSignIn: WorkspacePrompt = {
	id: "sign-in-required",
	kind: "sign-in",
	title: "Sign in to unlock custom execution",
	body: "Custom URLs, saved requests, history, and variables unlock after sign-in.",
	tone: "warning",
	action: {
		kind: "sign-in",
		label: "Sign in",
	},
};

const promptUpgrade: WorkspacePrompt = {
	id: "future-upgrade",
	kind: "upgrade",
	title: "Higher quotas are reserved for later plan tiers",
	body: "The quota model leaves room for higher tiers without changing the workspace surface.",
	tone: "neutral",
	action: {
		kind: "upgrade",
		label: "View upgrade",
	},
};

const promptLearnMore: WorkspacePrompt = {
	id: "guest-safety",
	kind: "learn-more",
	title: "Guest access stays allowlisted",
	body: "Visible lock states and allowlisted endpoints keep guest execution constrained.",
	tone: "positive",
	action: {
		kind: "learn-more",
		label: "Read safety notes",
	},
};

const baseControls: readonly WorkspaceControl[] = [
	{
		id: "template-browser",
		label: "Templates browser",
		surface: "sidebar",
		description: "Allowlisted templates and collections.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
	},
	{
		id: "method-selector",
		label: "Method selector",
		surface: "request-builder",
		description: "HTTP method control.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
	},
	{
		id: "url-bar",
		label: "URL bar",
		surface: "request-builder",
		description: "Show the target endpoint for the template or signed-in request.",
		visible: true,
		editable: false,
		locked: true,
		availability: "preview",
		constraint: "Guests can inspect the target but cannot replace the allowlisted URL.",
		promptId: promptSignIn.id,
	},
	{
		id: "query-params",
		label: "Query params",
		surface: "request-builder",
		description: "Override safe query values that the template explicitly exposes.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
		constraint: "Only allowlisted overrides are editable in guest mode.",
	},
	{
		id: "headers-editor",
		label: "Headers editor",
		surface: "request-builder",
		description: "Edit safe request headers for the selected template.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
		constraint: "Template-defined headers remain editable while the target stays locked.",
	},
	{
		id: "auth-selector",
		label: "Auth selector",
		surface: "request-builder",
		description: "Switch between no auth and template-defined auth examples.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
		constraint: "Custom bearer tokens and account credentials unlock after sign-in.",
	},
	{
		id: "body-editor",
		label: "Body editor",
		surface: "request-builder",
		description: "Edit JSON, raw, or form-urlencoded bodies when the template allows it.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
		constraint: "Guest templates only expose body modes that are safe to override.",
	},
	{
		id: "send-request",
		label: "Send request",
		surface: "request-builder",
		description: "Run the selected template or validated authenticated request.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
		constraint: "Guest sends stay confined to allowlisted endpoints.",
	},
	{
		id: "save-request",
		label: "Save request",
		surface: "sidebar",
		description: "Persist a request into a reusable collection.",
		visible: true,
		editable: false,
		locked: true,
		availability: "locked",
		promptId: promptSignIn.id,
	},
	{
		id: "history-rail",
		label: "History rail",
		surface: "sidebar",
		description: "Show recent runs and saved history snapshots.",
		visible: true,
		editable: false,
		locked: true,
		availability: "locked",
		promptId: promptSignIn.id,
	},
	{
		id: "environment-variables",
		label: "Environment variables",
		surface: "utility-drawer",
		description: "Manage variables that are applied to requests at runtime.",
		visible: true,
		editable: false,
		locked: true,
		availability: "locked",
		promptId: promptSignIn.id,
	},
	{
		id: "custom-url",
		label: "Custom URL",
		surface: "request-builder",
		description: "Replace the allowlisted target with a custom endpoint.",
		visible: true,
		editable: false,
		locked: true,
		availability: "locked",
		promptId: promptSignIn.id,
	},
	{
		id: "snippet-drawer",
		label: "Code snippets",
		surface: "utility-drawer",
		description: "Request snippets for the current state.",
		visible: true,
		editable: true,
		locked: false,
		availability: "available",
	},
	{
		id: "advanced-tools",
		label: "Advanced tools",
		surface: "utility-drawer",
		description: "Additional workspace tooling.",
		visible: true,
		editable: false,
		locked: true,
		availability: "locked",
		promptId: promptSignIn.id,
	},
] as const;

const guestControls: readonly WorkspaceControl[] = [...baseControls];

const authenticatedControls: readonly WorkspaceControl[] = [
	...baseControls.map((control) => ({
		...control,
		editable: true,
		locked: false,
		availability: "available" as const,
		promptId: undefined,
		constraint:
			control.id === "url-bar"
				? "Validated custom URLs are allowed once the request is signed in."
				: control.constraint,
	})),
];

const guestLockedActions: readonly WorkspaceLockedAction[] = [
	{
		id: "save-request",
		label: "Save request",
		description: "Persistence remains locked in guest mode.",
		surface: "sidebar",
		promptId: promptSignIn.id,
		previewCopy: "Save request to a collection",
	},
	{
		id: "history-persistence",
		label: "History persistence",
		description: "Guest history remains preview-only.",
		surface: "sidebar",
		promptId: promptSignIn.id,
		previewCopy: "View saved request history",
	},
	{
		id: "environment-variables",
		label: "Environment variables",
		description: "Variables are locked in guest mode.",
		surface: "utility-drawer",
		promptId: promptSignIn.id,
		previewCopy: "Configure request variables",
	},
	{
		id: "custom-url",
		label: "Custom URL",
		description: "Guest execution stays on allowlisted endpoints.",
		surface: "request-builder",
		promptId: promptSignIn.id,
		previewCopy: "Replace the target URL",
	},
	{
		id: "advanced-tools",
		label: "Advanced tools",
		description: "Advanced tooling remains locked.",
		surface: "utility-drawer",
		promptId: promptSignIn.id,
		previewCopy: "Open advanced tooling",
	},
];

const authenticatedLockedActions: readonly WorkspaceLockedAction[] = [
	{
		id: "higher-tier-quotas",
		label: "Higher tier quotas",
		description: "The quota model can expand without changing the workspace surface.",
		surface: "toolbar",
		promptId: promptUpgrade.id,
		previewCopy: "Increase request limits",
	},
];

export const workspaceTemplates: readonly WorkspaceTemplate[] = [
	{
		slug: "jsonplaceholder-posts",
		category: "REST basics",
		title: "JSONPlaceholder posts",
		summary: "GET a public post and inspect the response structure.",
		description:
			"A guest-safe starter template for the canonical JSONPlaceholder demo request. It keeps the target fixed while allowing safe query and header overrides.",
		availability: "both",
		allowlistedTarget: "https://jsonplaceholder.typicode.com/posts/1",
		safeOverrides: ["query", "headers"],
		tags: ["guest-safe", "public-api", "json"],
		request: {
			method: "GET",
			url: "https://jsonplaceholder.typicode.com/posts/1",
			bodyMode: "none",
			query: [],
			headers: [{ key: "accept", value: "application/json" }],
			responseStatus: 200,
			responseStatusText: "OK",
			responseTimeMs: 186,
			responseSizeLabel: "1.2 KB",
			responseContentType: "application/json; charset=utf-8",
			responseBody: `{
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}`,
		},
	},
	{
		slug: "github-public-user",
		category: "REST basics",
		title: "GitHub public user",
		summary: "Inspect a public profile with a response that feels closer to real product usage.",
		description:
			"This template demonstrates a public API with rate-limit headers, a richer response payload, and a strong preview for the response viewer.",
		availability: "both",
		allowlistedTarget: "https://api.github.com/users/octocat",
		safeOverrides: ["headers"],
		tags: ["guest-safe", "headers", "response-viewer"],
		request: {
			method: "GET",
			url: "https://api.github.com/users/octocat",
			bodyMode: "none",
			query: [],
			headers: [{ key: "accept", value: "application/vnd.github+json" }],
			responseStatus: 200,
			responseStatusText: "OK",
			responseTimeMs: 241,
			responseSizeLabel: "3.8 KB",
			responseContentType: "application/json; charset=utf-8",
			responseBody: `{
  "login": "octocat",
  "id": 1,
  "type": "User"
}`,
		},
	},
	{
		slug: "weather-demo",
		category: "REST basics",
		title: "Weather demo",
		summary: "Try safe parameter overrides on a weather-style endpoint.",
		description:
			"The weather demo shows how a guest template can expose controlled query editing while keeping the target locked and predictable.",
		availability: "both",
		allowlistedTarget: "https://api-testing-kit.example/demo/weather",
		safeOverrides: ["query"],
		tags: ["guest-safe", "query", "demo"],
		request: {
			method: "GET",
			url: "https://api-testing-kit.example/demo/weather?city=Kolkata",
			bodyMode: "none",
			query: [{ key: "city", value: "Kolkata" }],
			headers: [{ key: "accept", value: "application/json" }],
			responseStatus: 200,
			responseStatusText: "OK",
			responseTimeMs: 164,
			responseSizeLabel: "2.1 KB",
			responseContentType: "application/json",
			responseBody: `{
  "city": "Kolkata",
  "forecast": "partly cloudy",
  "temperature": 31
}`,
		},
	},
	{
		slug: "auth-flow-login",
		category: "Authentication flows",
		title: "Auth flow mock",
		summary: "Preview a login request with a templated body and safe auth guidance.",
		description:
			"This example keeps the shell honest by showing an authentication flow without exposing real credentials or arbitrary destinations.",
		availability: "both",
		allowlistedTarget: "https://api-testing-kit.example/mock/login",
		safeOverrides: ["headers", "body"],
		tags: ["auth", "json-body", "guest-safe"],
		request: {
			method: "POST",
			url: "https://api-testing-kit.example/mock/login",
			bodyMode: "json",
			query: [],
			headers: [{ key: "content-type", value: "application/json" }],
			responseStatus: 401,
			responseStatusText: "Unauthorized",
			responseTimeMs: 92,
			responseSizeLabel: "612 B",
			responseContentType: "application/json",
			responseBody: `{
  "error": "missing_token",
  "message": "This is a mock authentication example."
}`,
		},
	},
	{
		slug: "paged-crud-list",
		category: "Pagination examples",
		title: "Paged list",
		summary: "Show list navigation, paging, and count metadata in one template.",
		description:
			"The paged list template exercises query params, response metadata, and the preview of a collection-like API without requiring an editable target.",
		availability: "both",
		allowlistedTarget: "https://api-testing-kit.example/demo/items",
		safeOverrides: ["query", "headers"],
		tags: ["pagination", "list", "metadata"],
		request: {
			method: "GET",
			url: "https://api-testing-kit.example/demo/items?page=2&limit=10",
			bodyMode: "none",
			query: [
				{ key: "page", value: "2" },
				{ key: "limit", value: "10" },
			],
			headers: [{ key: "accept", value: "application/json" }],
			responseStatus: 200,
			responseStatusText: "OK",
			responseTimeMs: 154,
			responseSizeLabel: "4.4 KB",
			responseContentType: "application/json",
			responseBody: `{
  "page": 2,
  "limit": 10,
  "total": 48
}`,
		},
	},
	{
		slug: "webhook-echo",
		category: "Webhooks",
		title: "Webhook echo",
		summary: "Preview a webhook-shaped payload with a locked target.",
		description:
			"The webhook example is useful for showing how the product can model event delivery without turning guest mode into a general-purpose relay.",
		availability: "both",
		allowlistedTarget: "https://api-testing-kit.example/mock/webhook-echo",
		safeOverrides: ["headers", "body"],
		tags: ["webhook", "event", "json-body"],
		request: {
			method: "POST",
			url: "https://api-testing-kit.example/mock/webhook-echo",
			bodyMode: "json",
			query: [],
			headers: [
				{ key: "content-type", value: "application/json" },
				{ key: "x-demo-signature", value: "preview" },
			],
			responseStatus: 202,
			responseStatusText: "Accepted",
			responseTimeMs: 138,
			responseSizeLabel: "808 B",
			responseContentType: "application/json",
			responseBody: `{
  "status": "queued",
  "delivery": "echo"
}`,
		},
	},
	{
		slug: "error-envelope",
		category: "Error handling",
		title: "Error envelope",
		summary: "Show a structured error payload for blocked or failed runs.",
		description:
			"This template gives the response viewer a stable blocked/error shape to render, which is useful for both guest lock states and safety rejections.",
		availability: "both",
		allowlistedTarget: "https://api-testing-kit.example/mock/error",
		safeOverrides: ["headers"],
		tags: ["error", "blocked", "safety"],
		request: {
			method: "GET",
			url: "https://api-testing-kit.example/mock/error",
			bodyMode: "none",
			query: [],
			headers: [{ key: "accept", value: "application/json" }],
			responseStatus: 403,
			responseStatusText: "Forbidden",
			responseTimeMs: 61,
			responseSizeLabel: "540 B",
			responseContentType: "application/json",
			responseBody: `{
  "error": {
    "code": "blocked_target",
    "message": "Requests to private network targets are not allowed."
  }
}`,
		},
	},
] as const;

export const workspaceTemplateGroups: readonly WorkspaceTemplateGroup[] = [
	{
		id: "rest-basics",
		label: "REST basics",
		summary: "Guest-friendly templates for public APIs and the core request workflow.",
		templateSlugs: ["jsonplaceholder-posts", "github-public-user", "weather-demo"],
	},
	{
		id: "auth-flows",
		label: "Authentication flows",
		summary: "Mocked sign-in patterns and safe auth preview requests.",
		templateSlugs: ["auth-flow-login"],
	},
	{
		id: "crud-examples",
		label: "CRUD examples",
		summary: "Collection-style requests that can later grow into authenticated workflows.",
		templateSlugs: ["paged-crud-list"],
	},
	{
		id: "pagination",
		label: "Pagination examples",
		summary: "List responses and query-driven paging metadata.",
		templateSlugs: ["paged-crud-list"],
	},
	{
		id: "webhooks",
		label: "Webhooks",
		summary: "Event delivery shapes that stay within allowlisted demo endpoints.",
		templateSlugs: ["webhook-echo"],
	},
	{
		id: "error-handling",
		label: "Error handling",
		summary: "Blocked and failure envelopes for the response viewer and safety states.",
		templateSlugs: ["error-envelope"],
	},
] as const;

const guestCollections: readonly WorkspaceCollectionPreview[] = [
	{
		id: "guest-demo-pack",
		title: "Guest demo pack",
		description: "Four curated requests that show the shell, response viewer, and safety model.",
		scope: "guest",
		requestCount: 4,
		templateSlugs: ["jsonplaceholder-posts", "github-public-user", "weather-demo", "auth-flow-login"],
		badge: "Featured",
		featured: true,
	},
	{
		id: "rest-foundations",
		title: "REST foundations",
		description: "A compact starter pack for public GET and POST examples.",
		scope: "both",
		requestCount: 3,
		templateSlugs: ["jsonplaceholder-posts", "github-public-user", "paged-crud-list"],
		badge: "Starter",
		featured: false,
	},
	{
		id: "safety-tour",
		title: "Safety states",
		description: "Blocked-target and constrained execution examples.",
		scope: "guest",
		requestCount: 2,
		templateSlugs: ["error-envelope", "webhook-echo"],
		badge: "Locked",
		featured: false,
	},
];

const authenticatedCollections: readonly WorkspaceCollectionPreview[] = [
	...guestCollections,
	{
		id: "saved-workflows",
		title: "Saved workflows",
		description: "Authenticated examples for requests that should persist across sessions.",
		scope: "authenticated",
		requestCount: 6,
		templateSlugs: ["paged-crud-list", "webhook-echo", "error-envelope"],
		badge: "Saved",
		featured: true,
	},
	{
		id: "release-checks",
		title: "Release checks",
		description: "A lightweight set of requests that can later expand into a runbook.",
		scope: "authenticated",
		requestCount: 4,
		templateSlugs: ["github-public-user", "error-envelope"],
		badge: "Pinned",
		featured: false,
	},
];

const guestHistory: readonly WorkspaceHistoryEntry[] = [
	{
		id: "hist-jsonplaceholder",
		title: "JSONPlaceholder post",
		scope: "guest",
		method: "GET",
		target: "jsonplaceholder.typicode.com/posts/1",
		statusCode: 200,
		statusText: "OK",
		durationMs: 186,
		responseSizeLabel: "1.2 KB",
		contentType: "application/json",
		timestampLabel: "2 min ago",
		source: "demo",
		outcome: "success",
	},
	{
		id: "hist-github",
		title: "GitHub public user",
		scope: "guest",
		method: "GET",
		target: "api.github.com/users/octocat",
		statusCode: 200,
		statusText: "OK",
		durationMs: 241,
		responseSizeLabel: "3.8 KB",
		contentType: "application/json",
		timestampLabel: "9 min ago",
		source: "demo",
		outcome: "success",
	},
	{
		id: "hist-weather",
		title: "Weather demo",
		scope: "guest",
		method: "GET",
		target: "api-testing-kit.example/demo/weather?city=Kolkata",
		statusCode: 200,
		statusText: "OK",
		durationMs: 164,
		responseSizeLabel: "2.1 KB",
		contentType: "application/json",
		timestampLabel: "18 min ago",
		source: "demo",
		outcome: "success",
	},
	{
		id: "hist-auth-mock",
		title: "Auth mock login",
		scope: "guest",
		method: "POST",
		target: "api-testing-kit.example/mock/login",
		statusCode: 401,
		statusText: "Unauthorized",
		durationMs: 92,
		responseSizeLabel: "612 B",
		contentType: "application/json",
		timestampLabel: "26 min ago",
		source: "demo",
		outcome: "blocked",
	},
];

const authenticatedHistory: readonly WorkspaceHistoryEntry[] = [
	...guestHistory,
	{
		id: "hist-authenticated-crud",
		title: "Paginated list",
		scope: "authenticated",
		method: "GET",
		target: "api-testing-kit.example/demo/items?page=2&limit=10",
		statusCode: 200,
		statusText: "OK",
		durationMs: 154,
		responseSizeLabel: "4.4 KB",
		contentType: "application/json",
		timestampLabel: "Today",
		source: "persistent",
		outcome: "success",
	},
	{
		id: "hist-webhook",
		title: "Webhook echo",
		scope: "authenticated",
		method: "POST",
		target: "api-testing-kit.example/mock/webhook-echo",
		statusCode: 202,
		statusText: "Accepted",
		durationMs: 138,
		responseSizeLabel: "808 B",
		contentType: "application/json",
		timestampLabel: "Today",
		source: "persistent",
		outcome: "success",
	},
];

const guestMetrics: readonly WorkspaceMetric[] = [
	{
		id: "requests-today",
		label: "Requests today",
		value: "12",
		detail: "guest-safe sends",
		scope: "guest",
		tone: "neutral",
	},
	{
		id: "success-rate",
		label: "Success rate",
		value: "98%",
		detail: "preview responses",
		scope: "both",
		tone: "positive",
	},
	{
		id: "average-response-time",
		label: "Average response time",
		value: "186 ms",
		detail: "demo endpoints",
		scope: "guest",
		tone: "neutral",
	},
	{
		id: "collections-saved",
		label: "Collections saved",
		value: "Locked",
		detail: "available after sign-in",
		scope: "guest",
		tone: "warning",
	},
	{
		id: "remaining-quota",
		label: "Remaining quota",
		value: "8 / 10",
		detail: "per 10 minutes",
		scope: "guest",
		tone: "warning",
	},
];

const authenticatedMetrics: readonly WorkspaceMetric[] = [
	{
		id: "requests-today",
		label: "Requests today",
		value: "38",
		detail: "signed-in sends",
		scope: "authenticated",
		tone: "neutral",
	},
	{
		id: "success-rate",
		label: "Success rate",
		value: "99%",
		detail: "validated requests",
		scope: "both",
		tone: "positive",
	},
	{
		id: "average-response-time",
		label: "Average response time",
		value: "142 ms",
		detail: "recent runs",
		scope: "authenticated",
		tone: "neutral",
	},
	{
		id: "collections-saved",
		label: "Collections saved",
		value: "12",
		detail: "persisted requests",
		scope: "authenticated",
		tone: "positive",
	},
	{
		id: "remaining-quota",
		label: "Remaining quota",
		value: "162 / 200",
		detail: "per day",
		scope: "authenticated",
		tone: "warning",
	},
];

const guestQuotas: readonly WorkspaceQuota[] = [
	{
		id: "guest-10-min",
		label: "Requests per 10 minutes",
		scope: "guest",
		limit: 10,
		limitLabel: "10 requests",
		used: 2,
		usedLabel: "2 sent",
		remaining: 8,
		remainingLabel: "8 left",
		windowLabel: "10 minutes",
		note: "Guest sends stay allowlisted and rate limited per IP.",
	},
	{
		id: "guest-day",
		label: "Requests per day",
		scope: "guest",
		limit: 30,
		limitLabel: "30 requests",
		used: 12,
		usedLabel: "12 sent",
		remaining: 18,
		remainingLabel: "18 left",
		windowLabel: "24 hours",
		note: "The guest experience is intentionally small but real.",
	},
	{
		id: "guest-concurrency",
		label: "Concurrent requests",
		scope: "guest",
		limit: 1,
		limitLabel: "1 active request",
		used: 1,
		usedLabel: "1 active",
		remaining: 0,
		remainingLabel: "No spare capacity",
		windowLabel: "at a time",
		note: "Only one guest request can be active per IP.",
	},
	{
		id: "guest-body-size",
		label: "Request body size",
		scope: "guest",
		limit: 64 * 1024,
		limitLabel: "64 KB",
		used: 24 * 1024,
		usedLabel: "24 KB used",
		remaining: 40 * 1024,
		remainingLabel: "40 KB left",
		windowLabel: "per request",
		note: "Large payloads are blocked to reduce abuse risk.",
	},
	{
		id: "guest-response-size",
		label: "Response preview size",
		scope: "guest",
		limit: 512 * 1024,
		limitLabel: "512 KB",
		used: 192 * 1024,
		usedLabel: "192 KB shown",
		remaining: 320 * 1024,
		remainingLabel: "320 KB left",
		windowLabel: "previewed",
		note: "The viewer trims oversized responses for safety.",
	},
	{
		id: "guest-timeout",
		label: "Timeout",
		scope: "guest",
		limit: 10,
		limitLabel: "10 seconds",
		used: 6,
		usedLabel: "6 seconds",
		remaining: 4,
		remainingLabel: "4 seconds left",
		windowLabel: "per request",
		note: "Guest requests should fail fast and stay demo-friendly.",
	},
];

const authenticatedQuotas: readonly WorkspaceQuota[] = [
	{
		id: "auth-hour",
		label: "Requests per hour",
		scope: "authenticated",
		limit: 60,
		limitLabel: "60 requests",
		used: 18,
		usedLabel: "18 sent",
		remaining: 42,
		remainingLabel: "42 left",
		windowLabel: "1 hour",
		note: "Authenticated users get a larger but still bounded budget.",
	},
	{
		id: "auth-day",
		label: "Requests per day",
		scope: "authenticated",
		limit: 200,
		limitLabel: "200 requests",
		used: 38,
		usedLabel: "38 sent",
		remaining: 162,
		remainingLabel: "162 left",
		windowLabel: "24 hours",
		note: "Daily usage stays conservative enough for public hosting.",
	},
	{
		id: "auth-concurrency",
		label: "Concurrent requests",
		scope: "authenticated",
		limit: 5,
		limitLabel: "5 active requests",
		used: 2,
		usedLabel: "2 active",
		remaining: 3,
		remainingLabel: "3 spare",
		windowLabel: "at a time",
		note: "Concurrency stays capped to prevent bursty abuse.",
	},
	{
		id: "auth-body-size",
		label: "Request body size",
		scope: "authenticated",
		limit: 256 * 1024,
		limitLabel: "256 KB",
		used: 72 * 1024,
		usedLabel: "72 KB used",
		remaining: 184 * 1024,
		remainingLabel: "184 KB left",
		windowLabel: "per request",
		note: "Bodies can grow for real workflows, but they remain bounded.",
	},
	{
		id: "auth-response-limit",
		label: "Response processing limit",
		scope: "authenticated",
		limit: 1024 * 1024,
		limitLabel: "1 MB",
		used: 340 * 1024,
		usedLabel: "340 KB processed",
		remaining: 684 * 1024,
		remainingLabel: "684 KB left",
		windowLabel: "per response",
		note: "The runner trims large responses before rendering them in the UI.",
	},
	{
		id: "auth-timeout",
		label: "Timeout",
		scope: "authenticated",
		limit: 15,
		limitLabel: "15 seconds",
		used: 8,
		usedLabel: "8 seconds",
		remaining: 7,
		remainingLabel: "7 seconds left",
		windowLabel: "per request",
		note: "Authenticated requests get a slightly larger safety window.",
	},
	{
		id: "auth-redirects",
		label: "Redirect count",
		scope: "authenticated",
		limit: 3,
		limitLabel: "3 redirects",
		used: 1,
		usedLabel: "1 redirect",
		remaining: 2,
		remainingLabel: "2 left",
		windowLabel: "per request",
		note: "Redirects stay low to keep destination validation strict.",
	},
];

export const workspaceStates: Record<WorkspaceMode, WorkspaceState> = {
	guest: {
		mode: "guest",
		title: "Guest preview",
		subtitle: "Allowlisted templates, visible lock states, and preview history.",
		accessSummary: "Guest mode supports allowlisted execution and preview-only persistence surfaces.",
		prompts: [promptSignIn, promptLearnMore, promptUpgrade],
		controls: guestControls,
		lockedActions: guestLockedActions,
		templateGroups: workspaceTemplateGroups,
		templates: workspaceTemplates,
		collections: guestCollections,
		history: guestHistory,
		metrics: guestMetrics,
		quotas: guestQuotas,
	},
	authenticated: {
		mode: "authenticated",
		title: "Signed-in workspace",
		subtitle: "Validated custom URLs, saved requests, history, and broader request limits.",
		accessSummary: "Signed-in mode unlocks custom execution and persistence on the same workspace route.",
		prompts: [promptUpgrade, promptLearnMore],
		controls: authenticatedControls,
		lockedActions: authenticatedLockedActions,
		templateGroups: workspaceTemplateGroups,
		templates: workspaceTemplates,
		collections: authenticatedCollections,
		history: authenticatedHistory,
		metrics: authenticatedMetrics,
		quotas: authenticatedQuotas,
	},
} as const;

export const guestWorkspaceState = workspaceStates.guest;

export const authenticatedWorkspaceState = workspaceStates.authenticated;

export function getWorkspaceState(mode: WorkspaceMode): WorkspaceState {
	return workspaceStates[mode];
}

export function getWorkspaceControl(mode: WorkspaceMode, controlId: WorkspaceControlId): WorkspaceControl | undefined {
	return workspaceStates[mode].controls.find((control) => control.id === controlId);
}

export function isWorkspaceControlLocked(mode: WorkspaceMode, controlId: WorkspaceControlId): boolean {
	return getWorkspaceControl(mode, controlId)?.locked ?? false;
}

export function getWorkspacePrompt(mode: WorkspaceMode, promptId: string): WorkspacePrompt | undefined {
	return workspaceStates[mode].prompts.find((prompt) => prompt.id === promptId);
}

export function getWorkspaceTemplate(slug: string): WorkspaceTemplate | undefined {
	return workspaceTemplates.find((template) => template.slug === slug);
}

export function getWorkspaceCollection(id: string): WorkspaceCollectionPreview | undefined {
	return workspaceStates.guest.collections
		.concat(workspaceStates.authenticated.collections)
		.find((collection) => collection.id === id);
}

export function getWorkspaceTemplateGroup(id: string): WorkspaceTemplateGroup | undefined {
	return workspaceTemplateGroups.find((group) => group.id === id);
}
