import type { PageServerLoad } from "./$types";

import {
	buildCollectionDetail,
	buildLiveCollectionDetail,
	buildLiveMissingCollectionDetail,
	buildUnavailableCollectionDetail,
} from "$lib/components/collections/collection-detail-data";
import {
	readBackendJson,
	type BackendCollectionDetailResponse,
	type BackendCollectionsListResponse,
} from "$lib/server/backend";

export const load = (async ({ params, parent, fetch, request }) => {
	const { mode } = await parent();

	if (mode !== "authenticated") {
		return {
			detail: buildCollectionDetail(params.id, mode),
		};
	}

	const cookie = request.headers.get("cookie");

	try {
		const [detailResult, listResult] = await Promise.all([
			readBackendJson<BackendCollectionDetailResponse>(fetch, `/api/v1/collections/${encodeURIComponent(params.id)}`, {
				cookie,
			}),
			readBackendJson<BackendCollectionsListResponse>(fetch, "/api/v1/collections", {
				cookie,
			}),
		]);

		if (detailResult.response.status === 404) {
			return {
				detail: buildLiveMissingCollectionDetail(params.id),
			};
		}

		if (!detailResult.response.ok || !detailResult.data?.collection) {
			return {
				detail: buildUnavailableCollectionDetail(params.id),
			};
		}

		return {
			detail: buildLiveCollectionDetail(
				detailResult.data.collection,
				detailResult.data.savedRequests ?? [],
				listResult.data?.collections ?? [],
			),
		};
	} catch {
		return {
			detail: buildUnavailableCollectionDetail(params.id),
		};
	}
}) satisfies PageServerLoad;
