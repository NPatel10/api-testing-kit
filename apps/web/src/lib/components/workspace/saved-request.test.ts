import { describe, expect, it } from "vitest";

import { createSavedRequestDraft } from "./saved-request";

describe("createSavedRequestDraft", () => {
	it("hydrates a saved request into the workspace builder shape", () => {
		const draft = createSavedRequestDraft("authenticated", {
			id: "request-1",
			name: "Create invoice",
			method: "POST",
			url: "https://api.example.com/invoices",
			sortOrder: 1,
			queryParams: [{ name: "expand", value: "customer", enabled: true }],
			headers: { accept: "application/json", "content-type": "application/json" },
			authScheme: "bearer",
			authConfig: { token: "secret-token" },
			bodyMode: "json",
			bodyConfig: {
				raw: { invoiceId: "inv_123", amount: 42 },
			},
			createdAt: "2026-03-26T10:00:00Z",
			updatedAt: "2026-03-26T10:00:00Z",
		});

		expect(draft.method).toBe("POST");
		expect(draft.url).toBe("https://api.example.com/invoices");
		expect(draft.queryParams[0]).toMatchObject({ key: "expand", value: "customer", enabled: true });
		expect(draft.headers).toEqual(
			expect.arrayContaining([
				expect.objectContaining({ key: "accept", value: "application/json" }),
				expect.objectContaining({ key: "content-type", value: "application/json" }),
			]),
		);
		expect(draft.auth.scheme).toBe("bearer");
		expect(draft.auth.token).toBe("secret-token");
		expect(draft.body.mode).toBe("json");
		expect(draft.body.value).toContain('"invoiceId": "inv_123"');
	});
});
