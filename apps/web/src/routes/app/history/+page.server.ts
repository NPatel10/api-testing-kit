import type { PageServerLoad } from "./$types";

import {
	buildHistoryPageData,
	buildLiveHistoryPageData,
	buildUnavailableHistoryPageData,
} from "$lib/components/history/history-page-data";
import { readBackendJson, type BackendHistoryListResponse } from "$lib/server/backend";

export const load = (async ({ parent, url, fetch, request }) => {
	const data = await parent();
	const status = url.searchParams.get("status");
	const method = url.searchParams.get("method");
	const domain = url.searchParams.get("domain");

	if (data.mode !== "authenticated") {
		return buildHistoryPageData(data.mode, status, method, domain);
	}

	const query = new URLSearchParams();
	for (const key of ["status", "method", "domain", "date", "page", "limit"] as const) {
		const value = url.searchParams.get(key);
		if (value) {
			query.set(key, value);
		}
	}

	try {
		const result = await readBackendJson<BackendHistoryListResponse>(
			fetch,
			`/api/v1/history${query.size > 0 ? `?${query.toString()}` : ""}`,
			{ cookie: request.headers.get("cookie") },
		);

		if (!result.response.ok || !result.data) {
			return buildUnavailableHistoryPageData(data.mode, status, method, domain);
		}

		return buildLiveHistoryPageData(data.mode, result.data, status, method, domain);
	} catch {
		return buildUnavailableHistoryPageData(data.mode, status, method, domain);
	}

}) satisfies PageServerLoad;

