import {
	workspaceTemplateCategories,
	workspaceTemplates,
} from "$lib/mocks/workspace-state";

export const templatesPageCategories = [
	"all",
	...workspaceTemplateCategories,
] as const;

export type TemplatesPageCategory = (typeof templatesPageCategories)[number];

export const templatesPageCategoryLabels: Record<TemplatesPageCategory, string> = {
	all: "All templates",
	"REST basics": "REST basics",
	"Authentication flows": "Authentication flows",
	"CRUD examples": "CRUD examples",
	"Pagination examples": "Pagination examples",
	Webhooks: "Webhooks",
	"Error handling": "Error handling",
};

export const templatesPageCategoryDescriptions: Record<
	Exclude<TemplatesPageCategory, "all">,
	string
> = {
	"REST basics": "Public API examples that show the core request and response workflow.",
	"Authentication flows": "Login and token patterns that explain how sign-in changes the surface.",
	"CRUD examples": "Request shapes that cover create, update, and delete-style interactions.",
	"Pagination examples": "Repeated fetches and list navigation with query-driven controls.",
	Webhooks: "Outbound event payloads and delivery shapes with a fixed safe target.",
	"Error handling": "Responses that surface non-200 states and blocked or failure paths.",
};

export type TemplatesPageFilter = "all" | "body" | "query" | "headers" | "write" | "error";

export const templatesPageFilterLabels: Record<TemplatesPageFilter, string> = {
	all: "All",
	body: "Has body",
	query: "Has query params",
	headers: "Has headers",
	write: "Write methods",
	error: "Error responses",
};

export interface TemplatesPageRequestEntry {
	name: string;
	value: string;
	enabled?: boolean;
	overridable?: boolean;
}

export interface TemplatesPageResponsePreview {
	status: number;
	contentType: string;
	durationMs: number;
	size: string;
	body: string;
}

export interface TemplatesPageRequest {
	method: string;
	url: string;
	queryParams: TemplatesPageRequestEntry[];
	headers: TemplatesPageRequestEntry[];
	bodyMode: string;
}

export interface TemplatesPageTemplate {
	slug: string;
	title: string;
	category: TemplatesPageCategory;
	summary: string;
	description: string;
	tags: string[];
	guestSafe: boolean;
	allowlistedTarget: string;
	allowedOverrides: string[];
	request: TemplatesPageRequest;
	responsePreview: TemplatesPageResponsePreview;
	launchHref: string;
	previewHref: string;
	source: "live" | "fallback";
}

export interface TemplatesPageGroup {
	category: Exclude<TemplatesPageCategory, "all">;
	label: string;
	description: string;
	templates: TemplatesPageTemplate[];
}

export interface TemplatesPageData {
	templates: TemplatesPageTemplate[];
	sourceLabel: string;
	selectedSlug: string | null;
}

export interface BackendTemplateParam {
	name: string;
	value: string;
	enabled?: boolean;
	overridable?: boolean;
}

export interface BackendTemplateResponse {
	slug: string;
	title: string;
	description: string;
	category: Exclude<TemplatesPageCategory, "all">;
	tags: string[];
	guestSafe: boolean;
	allowedOverrides: string[];
	request: {
		method: string;
		url: string;
		queryParams: BackendTemplateParam[];
		headers: BackendTemplateParam[];
		auth?: {
			scheme: string;
			label?: string;
			description?: string;
			exampleValue?: string;
		};
		body: {
			mode: string;
			contentType?: string;
			example?: unknown;
			raw?: string;
			formFields?: BackendTemplateParam[];
		};
	};
	responsePreview: {
		status: number;
		contentType: string;
		durationMs: number;
		size: string;
		body: unknown;
	};
}

interface WorkspaceFallbackTemplate {
	slug: string;
	title: string;
	category: Exclude<TemplatesPageCategory, "all">;
	summary: string;
	description: string;
	availability: "guest" | "authenticated" | "both";
	allowlistedTarget: string;
	safeOverrides: readonly string[];
	tags: readonly string[];
	request: {
		method: string;
		url: string;
		bodyMode: string;
		query: readonly { key: string; value: string }[];
		headers: readonly { key: string; value: string }[];
		responseStatus: number;
		responseStatusText: string;
		responseTimeMs: number;
		responseSizeLabel: string;
		responseContentType: string;
		responseBody: string;
	};
}

function normalizeCategory(category: string): TemplatesPageCategory {
	return templatesPageCategories.includes(category as TemplatesPageCategory)
		? (category as TemplatesPageCategory)
		: "all";
}

function normalizeText(value: unknown) {
	return typeof value === "string" ? value : "";
}

export function formatPreviewBody(value: unknown): string {
	if (typeof value === "string") {
		return value;
	}

	try {
		return JSON.stringify(value, null, 2);
	} catch {
		return String(value);
	}
}

export function createLaunchHref(slug: string) {
	return `/app?template=${encodeURIComponent(slug)}`;
}

export function createPreviewHref(slug: string) {
	return `/app?template=${encodeURIComponent(slug)}&mode=preview`;
}

export function normalizeBackendTemplate(template: BackendTemplateResponse): TemplatesPageTemplate {
	return {
		slug: template.slug,
		title: template.title,
		category: normalizeCategory(template.category),
		summary: template.description,
		description: template.description,
		tags: template.tags,
		guestSafe: template.guestSafe,
		allowlistedTarget: template.request.url,
		allowedOverrides: template.allowedOverrides,
		request: {
			method: template.request.method,
			url: template.request.url,
			queryParams: template.request.queryParams.map((param) => ({
				name: param.name,
				value: param.value,
				enabled: param.enabled,
				overridable: param.overridable,
			})),
			headers: template.request.headers.map((param) => ({
				name: param.name,
				value: param.value,
				enabled: param.enabled,
				overridable: param.overridable,
			})),
			bodyMode: normalizeText(template.request.body?.mode),
		},
		responsePreview: {
			status: template.responsePreview.status,
			contentType: template.responsePreview.contentType,
			durationMs: template.responsePreview.durationMs,
			size: template.responsePreview.size,
			body: formatPreviewBody(template.responsePreview.body),
		},
		launchHref: createLaunchHref(template.slug),
		previewHref: createPreviewHref(template.slug),
		source: "live",
	};
}

export function normalizeFallbackTemplate(template: WorkspaceFallbackTemplate): TemplatesPageTemplate {
	return {
		slug: template.slug,
		title: template.title,
		category: template.category,
		summary: template.summary,
		description: template.description,
		tags: [...template.tags],
		guestSafe: template.availability !== "authenticated",
		allowlistedTarget: template.allowlistedTarget,
		allowedOverrides: [...template.safeOverrides],
		request: {
			method: template.request.method,
			url: template.request.url,
			queryParams: template.request.query.map((param) => ({
				name: param.key,
				value: param.value,
			})),
			headers: template.request.headers.map((param) => ({
				name: param.key,
				value: param.value,
			})),
			bodyMode: template.request.bodyMode,
		},
		responsePreview: {
			status: template.request.responseStatus,
			contentType: template.request.responseContentType,
			durationMs: template.request.responseTimeMs,
			size: template.request.responseSizeLabel,
			body: formatPreviewBody(template.request.responseBody),
		},
		launchHref: createLaunchHref(template.slug),
		previewHref: createPreviewHref(template.slug),
		source: "fallback",
	};
}

export function normalizeFallbackTemplates(): TemplatesPageTemplate[] {
	return workspaceTemplates.map((template) => normalizeFallbackTemplate(template as WorkspaceFallbackTemplate));
}

export function groupTemplatesByCategory(templates: TemplatesPageTemplate[]): TemplatesPageGroup[] {
	return workspaceTemplateCategories.map((category) => ({
		category,
		label: templatesPageCategoryLabels[category],
		description: templatesPageCategoryDescriptions[category],
		templates: templates.filter((template) => template.category === category),
	}));
}

export function templateMatchesSearch(template: TemplatesPageTemplate, query: string) {
	if (!query) {
		return true;
	}

	const haystack = [
		template.slug,
		template.title,
		template.summary,
		template.description,
		template.category,
		template.allowlistedTarget,
		template.request.method,
		template.request.url,
		template.request.bodyMode,
		...template.tags,
		...template.allowedOverrides,
		...template.request.queryParams.flatMap((param) => [param.name, param.value]),
		...template.request.headers.flatMap((param) => [param.name, param.value]),
	]
		.join(" ")
		.toLowerCase();

	return haystack.includes(query.trim().toLowerCase());
}

export function templateMatchesFilter(template: TemplatesPageTemplate, filter: TemplatesPageFilter) {
	switch (filter) {
		case "all":
			return true;
		case "body":
			return template.request.bodyMode !== "none";
		case "query":
			return template.request.queryParams.length > 0;
		case "headers":
			return template.request.headers.length > 0;
		case "write":
			return template.request.method !== "GET";
		case "error":
			return template.responsePreview.status >= 400;
	}
}

export function selectInitialTemplate(
	templates: TemplatesPageTemplate[],
	selectedSlug: string | null | undefined,
) {
	if (selectedSlug) {
		const selectedTemplate = templates.find((template) => template.slug === selectedSlug);
		if (selectedTemplate) {
			return selectedTemplate;
		}
	}

	return templates[0] ?? null;
}
