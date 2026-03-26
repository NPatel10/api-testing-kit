import { describe, expect, it } from "vitest";

import { buildLiveRunPayload, createPreviewViewerState, readLiveRunViewerState } from "./live-request";
import type { RequestBuilderDraft } from "./request-builder";

describe("live-request helpers", () => {
	it("builds the backend payload from the request draft", () => {
		const draft: RequestBuilderDraft = {
			method: "POST",
			url: "https://api.github.com/users/octocat",
			queryParams: [
				{ key: "page", value: "1", enabled: true },
				{ key: "ignored", value: "nope", enabled: false },
			],
			headers: [
				{ key: "accept", value: "application/json", enabled: true },
				{ key: "x-unused", value: "skip", enabled: false },
			],
			auth: {
				scheme: "bearer",
				token: "secret",
				username: "",
				password: "",
			},
			body: {
				mode: "form",
				value: "",
				formRows: [
					{ key: "firstName", value: "Ada", enabled: true },
					{ key: "lastName", value: "Lovelace", enabled: false },
				],
				contentType: "application/x-www-form-urlencoded",
			},
		};

		expect(buildLiveRunPayload(draft)).toEqual({
			method: "POST",
			url: "https://api.github.com/users/octocat",
			queryParams: [
				{ name: "page", value: "1", enabled: true },
				{ name: "ignored", value: "nope", enabled: false },
			],
			headers: [
				{ name: "accept", value: "application/json", enabled: true },
				{ name: "x-unused", value: "skip", enabled: false },
			],
			auth: {
				scheme: "bearer",
				token: "secret",
				username: "",
				password: "",
			},
			body: {
				mode: "form_urlencoded",
				raw: "",
				formFields: [
					{ name: "firstName", value: "Ada", enabled: true },
					{ name: "lastName", value: "Lovelace", enabled: false },
				],
			},
		});
	});

	it("normalizes a successful response into viewer state", async () => {
		const response = new Response(
			JSON.stringify({
				status: "succeeded",
				responseStatus: 200,
				responseHeaders: {
					"content-type": ["application/json"],
					"x-trace-id": ["trace-123"],
				},
				responseBodyPreview: "{\"ok\":true}",
				responseJson: { ok: true },
				responseSizeBytes: 12,
				responseTimeMs: 87,
				responseContentType: "application/json",
			}),
			{
				status: 200,
				headers: {
					"content-type": "application/json",
				},
			},
		);

		const state = await readLiveRunViewerState(response);

		expect(state.status).toBe(200);
		expect(state.statusText).toBe("OK");
		expect(state.prettyBody).toBe("{\n  \"ok\": true\n}");
		expect(state.rawBody).toBe("{\"ok\":true}");
		expect(state.headers).toEqual([
			{ key: "content-type", value: "application/json" },
			{ key: "x-trace-id", value: "trace-123" },
		]);
		expect(state.size).toBe(12);
		expect(state.duration).toBe(87);
		expect(state.error).toBeNull();
	});

	it("turns an error response into a readable viewer error", async () => {
		const response = new Response(
			JSON.stringify({
				error: {
					code: "blocked_target",
					message: "Requests to private network targets are not allowed.",
				},
			}),
			{
				status: 403,
				headers: {
					"content-type": "application/json",
				},
			},
		);

		const state = await readLiveRunViewerState(response);

		expect(state.error?.title).toBe("Blocked target");
		expect(state.error?.code).toBe("blocked_target");
		expect(state.error?.message).toContain("private network");
	});

	it("creates a preview state for the seeded demo mode", () => {
		const state = createPreviewViewerState({
			responseStatus: 200,
			responseStatusText: "OK",
			responseTimeMs: 150,
			responseSizeLabel: "1.2 KB",
			responseContentType: "application/json",
			responseBody: "{\"ok\":true}",
		});

		expect(state.statusText).toBe("OK");
		expect(state.headers[0]).toEqual({ key: "content-type", value: "application/json" });
		expect(state.prettyBody).toContain("\"ok\":true");
	});
});
