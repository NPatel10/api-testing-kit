import { getWorkspaceState, type WorkspaceHistoryEntry, type WorkspaceMode } from "$lib/mocks/workspace-state";

export type HistoryStatusFilter = "all" | "success" | "blocked" | "error";
export type HistoryMethodFilter = "all" | "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export const historyStatusFilters: readonly HistoryStatusFilter[] = ["all", "success", "blocked", "error"] as const;
export const historyMethodFilters: readonly HistoryMethodFilter[] = ["all", "GET", "POST", "PUT", "PATCH", "DELETE"] as const;

export const historyStatusLabels: Record<HistoryStatusFilter, string> = {
	all: "All",
	success: "Success",
	blocked: "Blocked",
	error: "Error",
};

export const historyMethodLabels: Record<HistoryMethodFilter, string> = {
	all: "Any method",
	GET: "GET",
	POST: "POST",
	PUT: "PUT",
	PATCH: "PATCH",
	DELETE: "DELETE",
};

export interface HistoryPageEntry extends WorkspaceHistoryEntry {
	domainLabel: string;
	timeLabel: string;
}

export interface HistoryPageSection {
	title: string;
	description: string;
	entries: HistoryPageEntry[];
}

export interface HistoryPageMetric {
	label: string;
	value: string;
	note: string;
}

export interface HistoryPageData {
	mode: WorkspaceMode;
	status: HistoryStatusFilter;
	method: HistoryMethodFilter;
	selectedDomain: string;
	domainOptions: string[];
	entries: HistoryPageEntry[];
	sections: HistoryPageSection[];
	metrics: HistoryPageMetric[];
	previewLocked: boolean;
}

function normalizeStatus(value: string | null | undefined): HistoryStatusFilter {
	return historyStatusFilters.includes(value as HistoryStatusFilter)
		? (value as HistoryStatusFilter)
		: "all";
}

function normalizeMethod(value: string | null | undefined): HistoryMethodFilter {
	return historyMethodFilters.includes(value as HistoryMethodFilter)
		? (value as HistoryMethodFilter)
		: "all";
}

function extractDomain(target: string) {
	const withProtocol = target.includes("://") ? target : `https://${target}`;
	try {
		return new URL(withProtocol).hostname;
	} catch {
		return target.split("/")[0] || target;
	}
}

function toTimeLabel(timestampLabel: string, source: WorkspaceHistoryEntry["source"]) {
	return source === "persistent" ? `${timestampLabel} - persisted` : `${timestampLabel} - demo preview`;
}

function enrichEntry(entry: WorkspaceHistoryEntry): HistoryPageEntry {
	return {
		...entry,
		domainLabel: extractDomain(entry.target),
		timeLabel: toTimeLabel(entry.timestampLabel, entry.source),
	};
}

function matchesStatus(entry: HistoryPageEntry, status: HistoryStatusFilter) {
	if (status === "all") return true;
	return entry.outcome === status || (status === "error" && entry.statusCode >= 500);
}

function matchesMethod(entry: HistoryPageEntry, method: HistoryMethodFilter) {
	if (method === "all") return true;
	return entry.method === method;
}

function matchesDomain(entry: HistoryPageEntry, domain: string) {
	if (!domain || domain === "all") return true;
	return entry.domainLabel === domain;
}

function buildSections(entries: HistoryPageEntry[]): HistoryPageSection[] {
	const succeeded = entries.filter((entry) => entry.outcome === "success");
	const blocked = entries.filter((entry) => entry.outcome === "blocked");
	const failed = entries.filter((entry) => entry.outcome === "error");

	return [
		{
			title: "Successful runs",
			description: "Validated responses and preserved previews.",
			entries: succeeded,
		},
		{
			title: "Blocked runs",
			description: "Requests that were rejected by guest or safety controls.",
			entries: blocked,
		},
		{
			title: "Errors",
			description: "Transport failures or upstream issues that need another look.",
			entries: failed,
		},
	].filter((section) => section.entries.length > 0);
}

function buildMetrics(entries: HistoryPageEntry[]): HistoryPageMetric[] {
	const total = entries.length;
	const successful = entries.filter((entry) => entry.outcome === "success").length;
	const blocked = entries.filter((entry) => entry.outcome === "blocked").length;
	const avgDuration = total
		? Math.round(entries.reduce((sum, entry) => sum + entry.durationMs, 0) / total)
		: 0;

	return [
		{
			label: "Visible runs",
			value: String(total),
			note: "Filtered to the current view",
		},
		{
			label: "Successful",
			value: String(successful),
			note: "Completed without safety rejection",
		},
		{
			label: "Blocked",
			value: String(blocked),
			note: "Guest or validation failures",
		},
		{
			label: "Average duration",
			value: `${avgDuration} ms`,
			note: "Across the current filtered set",
		},
	];
}

export function buildHistoryPageData(
	mode: WorkspaceMode,
	statusParam: string | null | undefined,
	methodParam: string | null | undefined,
	domainParam: string | null | undefined,
): HistoryPageData {
	const state = getWorkspaceState(mode);
	const status = normalizeStatus(statusParam);
	const method = normalizeMethod(methodParam);
	const selectedDomain = domainParam?.trim() || "all";

	const entries = state.history
		.map(enrichEntry)
		.filter((entry) => matchesStatus(entry, status))
		.filter((entry) => matchesMethod(entry, method))
		.filter((entry) => matchesDomain(entry, selectedDomain));

	return {
		mode,
		status,
		method,
		selectedDomain,
		domainOptions: Array.from(new Set(state.history.map((entry) => extractDomain(entry.target)))).sort(),
		entries,
		sections: buildSections(entries),
		metrics: buildMetrics(entries),
		previewLocked: mode === "guest",
	};
}
