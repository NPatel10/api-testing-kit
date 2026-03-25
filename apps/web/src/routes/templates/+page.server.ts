import { env as privateEnv } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import type { PageServerLoad } from "./$types";

import {
	normalizeBackendTemplate,
	normalizeFallbackTemplates,
	type BackendTemplateResponse,
	type TemplatesPageTemplate,
} from "$lib/components/templates-page/templates-page-data";

interface TemplatesApiResponse {
	templates: BackendTemplateResponse[];
}

function normalizeBaseUrl(value: string | undefined) {
	return (value ?? "http://localhost:8080").replace(/\/+$/, "");
}

async function loadTemplates(fetchFn: typeof fetch): Promise<{
	templates: TemplatesPageTemplate[];
	sourceLabel: string;
}> {
	const baseUrl = normalizeBaseUrl(
		privateEnv.INTERNAL_API_BASE_URL || privateEnv.API_BASE_URL || publicEnv.PUBLIC_API_BASE_URL,
	);

	try {
		const response = await fetchFn(`${baseUrl}/api/v1/templates`, {
			headers: {
				accept: "application/json",
			},
		});

		if (!response.ok) {
			throw new Error(`template API request failed with status ${response.status}`);
		}

		const payload = (await response.json()) as TemplatesApiResponse;
		if (!Array.isArray(payload.templates) || payload.templates.length === 0) {
			throw new Error("template API response did not include templates");
		}

		return {
			templates: payload.templates.map((template) => normalizeBackendTemplate(template)),
			sourceLabel: "Live backend",
		};
	} catch {
		return {
			templates: normalizeFallbackTemplates(),
			sourceLabel: "Local fallback",
		};
	}
}

export const load = (async ({ fetch, url }) => {
	const { templates, sourceLabel } = await loadTemplates(fetch);

	return {
		templates,
		sourceLabel,
		selectedSlug: url.searchParams.get("template") ?? templates[0]?.slug ?? null,
	};
}) satisfies PageServerLoad;
