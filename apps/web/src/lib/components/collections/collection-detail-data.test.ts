import { describe, expect, it } from "vitest";

import { buildCollectionDetail, buildLiveCollectionDetail } from "./collection-detail-data";

describe("buildCollectionDetail", () => {
	it("builds an authenticated collection view with grouped requests", () => {
		const detail = buildCollectionDetail("saved-workflows", "authenticated");

		expect(detail.status).toBe("ready");
		expect(detail.requestGroups.length).toBeGreaterThan(0);
		expect(detail.requestGroups[0].requests.length).toBeGreaterThan(0);
		expect(detail.actionLinks[0]?.href).toContain("/app?collection=saved-workflows");
	});

	it("returns a locked guest fallback for authenticated-only collections", () => {
		const detail = buildCollectionDetail("saved-workflows", "guest");

		expect(detail.status).toBe("locked");
		expect(detail.accessCopy).toContain("Sign in");
		expect(detail.previewRequests.length).toBeGreaterThan(0);
	});

	it("returns a missing-state fallback when the collection is not seeded", () => {
		const detail = buildCollectionDetail("does-not-exist", "guest");

		expect(detail.status).toBe("missing");
		expect(detail.actionLinks[0]?.href).toBe("/app");
	});

	it("builds a live authenticated collection from persisted saved requests", () => {
		const detail = buildLiveCollectionDetail(
			{
				id: "collection-1",
				name: "Billing APIs",
				description: "Saved billing flows",
				visibility: "private",
				sortOrder: 2,
				createdAt: "2026-03-26T10:00:00Z",
				updatedAt: "2026-03-26T11:00:00Z",
			},
			[
				{
					id: "request-1",
					collectionId: "collection-1",
					name: "Create invoice",
					description: "POST invoice payload",
					method: "POST",
					url: "https://api.example.com/invoices",
					sortOrder: 1,
					authScheme: "bearer",
					exampleResponse: {
						responseStatus: 201,
						responseStatusText: "Created",
						responseContentType: "application/json",
						responseTimeMs: 182,
						responseSizeBytes: 2048,
					},
					createdAt: "2026-03-26T10:05:00Z",
					updatedAt: "2026-03-26T10:05:00Z",
				},
			],
			[
				{
					id: "collection-1",
					name: "Billing APIs",
					description: "Saved billing flows",
					visibility: "private",
					sortOrder: 2,
					createdAt: "2026-03-26T10:00:00Z",
					updatedAt: "2026-03-26T11:00:00Z",
				},
				{
					id: "collection-2",
					name: "Support APIs",
					description: "Saved support flows",
					visibility: "shared_readonly",
					sortOrder: 3,
					createdAt: "2026-03-26T10:00:00Z",
					updatedAt: "2026-03-26T11:10:00Z",
				},
			],
		);

		expect(detail.source).toBe("live");
		expect(detail.status).toBe("ready");
		expect(detail.requestGroups[0]?.requests[0]?.launchHref).toBe("/app?request=request-1");
		expect(detail.requestGroups[0]?.requests[0]?.responseLabel).toContain("201");
		expect(detail.relatedCollections[0]?.id).toBe("collection-2");
	});
});
