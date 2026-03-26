import type { PageServerLoad } from "./$types";

import { readBackendJson, type BackendSavedRequest } from "$lib/server/backend";

export const load = (async ({ parent, url, fetch, request }) => {
	const { mode } = await parent();
	const requestId = url.searchParams.get("request")?.trim();

	if (mode !== "authenticated" || !requestId) {
		return {
			savedRequest: null,
		};
	}

	try {
		const result = await readBackendJson<BackendSavedRequest>(
			fetch,
			`/api/v1/requests/${encodeURIComponent(requestId)}`,
			{ cookie: request.headers.get("cookie") },
		);

		return {
			savedRequest: result.response.ok ? result.data : null,
		};
	} catch {
		return {
			savedRequest: null,
		};
	}
}) satisfies PageServerLoad;
