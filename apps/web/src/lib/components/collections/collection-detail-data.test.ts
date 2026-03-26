import { describe, expect, it } from "vitest";

import { buildCollectionDetail } from "./collection-detail-data";

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
});
