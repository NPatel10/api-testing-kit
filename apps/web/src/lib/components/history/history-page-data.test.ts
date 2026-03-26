import { describe, expect, it } from "vitest";

import { buildLiveHistoryPageData, buildUnavailableHistoryPageData } from "./history-page-data";

describe("history page data", () => {
	it("maps persisted history records into the signed-in route model", () => {
		const data = buildLiveHistoryPageData(
			"authenticated",
			{
				history: [
					{
						id: "run-1",
						savedRequestId: "request-1",
						status: "succeeded",
						source: "authenticated",
						method: "post",
						url: "https://api.example.com/invoices",
						targetHost: "api.example.com",
						responseStatus: 201,
						responseTimeMs: 182,
						responseSizeBytes: 2048,
						responseContentType: "application/json",
						redirectCount: 0,
						createdAt: "2026-03-26T11:15:00Z",
					},
				],
				pagination: {
					page: 1,
					limit: 20,
					hasMore: false,
				},
			},
			"all",
			"POST",
			"all",
		);

		expect(data.source).toBe("live");
		expect(data.previewLocked).toBe(false);
		expect(data.entries[0]?.statusLabel).toBe("HTTP 201");
		expect(data.entries[0]?.launchHref).toBe("/app?request=request-1");
		expect(data.entries[0]?.domainLabel).toBe("api.example.com");
	});

	it("returns a signed-in unavailable state without falling back to preview rows", () => {
		const data = buildUnavailableHistoryPageData("authenticated", "all", "all", "all");

		expect(data.source).toBe("live");
		expect(data.notice).toContain("temporarily unavailable");
		expect(data.entries).toHaveLength(0);
	});
});
