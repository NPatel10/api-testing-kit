import type { PageServerLoad } from "./$types";

import { buildCollectionDetail } from "$lib/components/collections/collection-detail-data";

export const load = (async ({ params, parent }) => {
	const { mode } = await parent();
	return {
		detail: buildCollectionDetail(params.id, mode),
	};
}) satisfies PageServerLoad;
