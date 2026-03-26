import { env as privateEnv } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import type { LayoutServerLoad } from "./$types";

import type { WorkspaceMode } from "$lib/mocks/workspace-state";

interface MeUser {
	id: string;
	email: string;
	displayName?: string;
}

interface MeResponse {
	user: MeUser;
}

function normalizeBaseUrl(value: string | undefined) {
	return (value ?? "http://localhost:8080").replace(/\/+$/, "");
}

function buildGuestState(): { mode: WorkspaceMode; sessionLabel: string } {
	return {
		mode: "guest",
		sessionLabel: "Guest preview",
	};
}

function buildAuthenticatedState(user: MeUser): { mode: WorkspaceMode; sessionLabel: string } {
	const sessionLabel = user.displayName?.trim() || user.email.trim() || "Signed-in user";
	return {
		mode: "authenticated",
		sessionLabel: `Signed in as ${sessionLabel}`,
	};
}

export const load = (async ({ fetch, request }) => {
	const baseUrl = normalizeBaseUrl(
		privateEnv.INTERNAL_API_BASE_URL || privateEnv.API_BASE_URL || publicEnv.PUBLIC_API_BASE_URL,
	);
	const cookie = request.headers.get("cookie");

	if (!cookie) {
		return buildGuestState();
	}

	try {
		const response = await fetch(`${baseUrl}/api/v1/me`, {
			headers: {
				accept: "application/json",
				cookie,
			},
			cache: "no-store",
		});

		if (!response.ok) {
			return buildGuestState();
		}

		const payload = (await response.json()) as Partial<MeResponse>;
		if (!payload.user?.id) {
			return buildGuestState();
		}

		return buildAuthenticatedState(payload.user);
	} catch {
		return buildGuestState();
	}
}) satisfies LayoutServerLoad;
