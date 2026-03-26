import { env as privateEnv } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import type { PageServerLoad } from "./$types";

import {
	buildPreviewAdminAbuseDashboard,
	normalizeAdminAbuseDashboard,
	type AdminAbuseApiResponse,
} from "$lib/components/admin-abuse/admin-abuse-data";

function normalizeBaseUrl(value: string | undefined) {
	return (value ?? "http://localhost:8080").replace(/\/+$/, "");
}

async function loadAdminAbuseDashboard(fetchFn: typeof fetch, cookie: string | null) {
	const baseUrl = normalizeBaseUrl(
		privateEnv.INTERNAL_API_BASE_URL || privateEnv.API_BASE_URL || publicEnv.PUBLIC_API_BASE_URL,
	);

	if (!cookie) {
		return buildPreviewAdminAbuseDashboard(
			"Admin access is required to load live abuse records. Showing the surface structure instead.",
		);
	}

	try {
		const response = await fetchFn(`${baseUrl}/api/v1/admin/abuse`, {
			headers: {
				accept: "application/json",
				cookie,
			},
			cache: "no-store",
		});

		if (response.ok) {
			const payload = (await response.json()) as Partial<AdminAbuseApiResponse>;
			if (Array.isArray(payload.summary) && Array.isArray(payload.recent) && Array.isArray(payload.blockedTargets)) {
				return normalizeAdminAbuseDashboard(payload);
			}
		}

		if (response.status === 401 || response.status === 403) {
			return buildPreviewAdminAbuseDashboard(
				"Live records are restricted to admin and owner sessions. Showing the monitored surface instead.",
			);
		}

		throw new Error(`admin abuse request failed with status ${response.status}`);
	} catch {
		return buildPreviewAdminAbuseDashboard(
			"The admin API is temporarily unavailable, so a structural preview is shown instead.",
		);
	}
}

export const load = (async ({ fetch, request }) => {
	const cookie = request.headers.get("cookie");
	const dashboard = await loadAdminAbuseDashboard(fetch, cookie);

	return {
		dashboard,
	};
}) satisfies PageServerLoad;
