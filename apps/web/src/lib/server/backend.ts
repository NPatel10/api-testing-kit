import { env as privateEnv } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import type { RequestEvent } from "@sveltejs/kit";

function normalizeBaseUrl(value: string | undefined) {
	return (value ?? "http://localhost:8080").replace(/\/+$/, "");
}

export function getBackendBaseUrl() {
	return normalizeBaseUrl(
		privateEnv.INTERNAL_API_BASE_URL || privateEnv.API_BASE_URL || publicEnv.PUBLIC_API_BASE_URL,
	);
}

export async function proxyBackendJson(event: RequestEvent, path: string, forwardCookie = false) {
	const headers = new Headers();
	const contentType = event.request.headers.get("content-type");
	if (contentType) {
		headers.set("content-type", contentType);
	}
	headers.set("accept", "application/json");

	if (forwardCookie) {
		const cookie = event.request.headers.get("cookie");
		if (cookie) {
			headers.set("cookie", cookie);
		}
	}

	const body = event.request.method === "GET" || event.request.method === "HEAD" ? undefined : await event.request.text();
	const response = await event.fetch(`${getBackendBaseUrl()}${path}`, {
		method: event.request.method,
		headers,
		body,
	});

	return new Response(await response.text(), {
		status: response.status,
		statusText: response.statusText,
		headers: {
			"content-type": response.headers.get("content-type") ?? "application/json",
			"cache-control": "no-store",
		},
	});
}
