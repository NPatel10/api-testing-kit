import type { RequestHandler } from "./$types";

import { proxyBackendJson } from "$lib/server/backend";

export const POST: RequestHandler = async (event) => {
	return proxyBackendJson(event, "/api/v1/runs", true);
};
