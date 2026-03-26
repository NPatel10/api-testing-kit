import type { PageServerLoad } from "./$types";

import { buildHistoryPageData } from "$lib/components/history/history-page-data";

export const load = (async ({ parent, url }) => {
	const data = await parent();

	return buildHistoryPageData(
		data.mode,
		url.searchParams.get("status"),
		url.searchParams.get("method"),
		url.searchParams.get("domain"),
	);
}) satisfies PageServerLoad;

