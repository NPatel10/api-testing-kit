import type { WorkspaceMode } from "$lib/mocks/workspace-state";

export const entitlementCapabilityKeys = [
	"custom_url_execution",
	"history_depth",
	"environment_variables",
	"shared_links",
] as const;

export type EntitlementCapabilityKey = (typeof entitlementCapabilityKeys)[number];

export interface EntitlementPlan {
	code: string;
	name: string;
	source: string;
}

export interface EntitlementCapability {
	key: EntitlementCapabilityKey;
	label: string;
	description: string;
	enabled: boolean;
	scope: "guest" | "authenticated" | "plan";
	limit?: number;
	limitLabel?: string;
	reason?: string;
}

export interface EffectiveEntitlements {
	plan: EntitlementPlan;
	capabilities: EntitlementCapability[];
}

export interface EntitlementCapabilityRow {
	key: EntitlementCapabilityKey;
	label: string;
	description: string;
	statusLabel: string;
	tone: "neutral" | "positive" | "warning";
	limitLabel?: string;
}

const guestEntitlements: EffectiveEntitlements = {
	plan: {
		code: "guest",
		name: "Guest Preview",
		source: "guest",
	},
	capabilities: [
		{
			key: "custom_url_execution",
			label: "Custom URL execution",
			description: "Guest sessions stay on allowlisted template targets.",
			enabled: false,
			scope: "guest",
			reason: "Sign in to send custom outbound requests.",
		},
		{
			key: "history_depth",
			label: "History depth",
			description: "Guest history is preview only and does not persist.",
			enabled: false,
			scope: "guest",
			reason: "Sign in to retain request history.",
		},
		{
			key: "environment_variables",
			label: "Environment variables",
			description: "Variables remain locked in guest mode.",
			enabled: false,
			scope: "guest",
			reason: "Sign in to use request variables.",
		},
		{
			key: "shared_links",
			label: "Shared links",
			description: "Guests can inspect sharing surfaces but cannot create them.",
			enabled: false,
			scope: "guest",
			reason: "Sign in to publish shared links.",
		},
	],
};

export const authenticatedEntitlements: EffectiveEntitlements = {
	plan: {
		code: "starter",
		name: "Starter",
		source: "system_default",
	},
	capabilities: [
		{
			key: "custom_url_execution",
			label: "Custom URL execution",
			description: "Validated custom URLs are available after sign-in.",
			enabled: true,
			scope: "authenticated",
		},
		{
			key: "history_depth",
			label: "History depth",
			description: "Signed-in users keep a bounded replay history.",
			enabled: true,
			scope: "authenticated",
			limit: 50,
			limitLabel: "50 runs",
		},
		{
			key: "environment_variables",
			label: "Environment variables",
			description: "Variables stay locked on the base plan.",
			enabled: false,
			scope: "plan",
			reason: "Upgrade to unlock environment variable storage.",
		},
		{
			key: "shared_links",
			label: "Shared links",
			description: "Readonly sharing remains part of a later tier.",
			enabled: false,
			scope: "plan",
			reason: "Upgrade to publish shared links.",
		},
	],
};

export function getDefaultEntitlements(mode: WorkspaceMode): EffectiveEntitlements {
	return mode === "authenticated" ? authenticatedEntitlements : guestEntitlements;
}

export function normalizeEntitlements(
	value: Partial<EffectiveEntitlements> | undefined,
	mode: WorkspaceMode,
): EffectiveEntitlements {
	const capabilities = value?.capabilities ?? [];

	if (!value?.plan || !Array.isArray(capabilities) || capabilities.length === 0) {
		return getDefaultEntitlements(mode);
	}

	return {
		plan: {
			code: value.plan.code || (mode === "authenticated" ? "starter" : "guest"),
			name: value.plan.name || (mode === "authenticated" ? "Starter" : "Guest Preview"),
			source: value.plan.source || (mode === "authenticated" ? "system_default" : "guest"),
		},
		capabilities: entitlementCapabilityKeys.map((key) => {
			const capability = capabilities.find((item) => item.key === key);
			if (capability) {
				return capability;
			}

			return getDefaultEntitlements(mode).capabilities.find((item) => item.key === key) ?? {
				key,
				label: key,
				description: "",
				enabled: false,
				scope: "plan",
			};
		}),
	};
}

export function getCapability(entitlements: EffectiveEntitlements, key: EntitlementCapabilityKey): EntitlementCapability {
	return (
		entitlements.capabilities.find((item) => item.key === key) ?? {
			key,
			label: key,
			description: "",
			enabled: false,
			scope: "plan",
		}
	);
}

export function buildEntitlementRows(entitlements: EffectiveEntitlements): readonly EntitlementCapabilityRow[] {
	return entitlements.capabilities.map((capability) => {
		const tone = capability.enabled ? "positive" : "warning";
		const statusLabel = capability.enabled ? "Unlocked" : "Locked";

		return {
			key: capability.key,
			label: capability.label,
			description: capability.description,
			statusLabel,
			tone,
			limitLabel: capability.limitLabel,
		};
	});
}

export function getEntitlementSummary(entitlements: EffectiveEntitlements, mode: WorkspaceMode): string {
	if (mode === "guest") {
		return "Guest preview mode keeps custom URLs, environment variables, and shared links locked behind sign-in.";
	}

	const lockedCount = entitlements.capabilities.filter((capability) => !capability.enabled).length;
	if (lockedCount === 0) {
		return "All plan-backed capabilities are currently unlocked for this session.";
	}

	return `${lockedCount} plan-backed ${lockedCount === 1 ? "capability" : "capabilities"} remain locked for this session.`;
}
