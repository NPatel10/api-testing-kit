import { getWorkspaceState, type WorkspaceMode } from "$lib/mocks/workspace-state";
import type { BackendHistoryListResponse, BackendRunRecord } from "$lib/server/backend";

export type HistoryStatusFilter = "all" | "success" | "blocked" | "error";
export type HistoryMethodFilter = "all" | "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
export type HistoryEntryOutcome = "success" | "blocked" | "error";

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

export interface HistoryPageEntry {
	id: string;
	title: string;
	method: string;
	target: string;
	domainLabel: string;
	timeLabel: string;
	statusLabel: string;
	durationLabel: string;
	durationMs: number;
	responseSizeLabel: string;
	contentTypeLabel: string;
	outcome: HistoryEntryOutcome;
	launchHref?: string;
	launchLabel?: string;
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
	source: "preview" | "live";
	status: HistoryStatusFilter;
	method: HistoryMethodFilter;
	selectedDomain: string;
	domainOptions: string[];
	entries: HistoryPageEntry[];
	sections: HistoryPageSection[];
	metrics: HistoryPageMetric[];
	previewLocked: boolean;
	notice?: string;
	pagination?: {
		page: number;
		limit: number;
		hasMore: boolean;
	};
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

function buildSections(entries: HistoryPageEntry[]): HistoryPageSection[] {
	const succeeded = entries.filter((entry) => entry.outcome === "success");
	const blocked = entries.filter((entry) => entry.outcome === "blocked");
	const failed = entries.filter((entry) => entry.outcome === "error");

	return [
		{
			title: "Successful runs",
			description: "Requests that completed and returned a live response snapshot.",
			entries: succeeded,
		},
		{
			title: "Blocked runs",
			description: "Requests rejected by validation, target rules, or quota controls.",
			entries: blocked,
		},
		{
			title: "Errors",
			description: "Runs that failed before a clean response was returned.",
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
			note: "Completed without a safety rejection",
		},
		{
			label: "Blocked",
			value: String(blocked),
			note: "Rejected by runner or rate-limit controls",
		},
		{
			label: "Average duration",
			value: `${avgDuration} ms`,
			note: "Across the current filtered set",
		},
	];
}

function outcomeFromRecord(record: BackendRunRecord): HistoryEntryOutcome {
	switch ((record.status ?? "").toLowerCase()) {
		case "succeeded":
			return "success";
		case "blocked":
			return "blocked";
		default:
			return "error";
	}
}

function toPreviewTimeLabel(timestampLabel: string, source: "demo" | "persistent") {
	return source === "persistent" ? `${timestampLabel} - persisted` : `${timestampLabel} - demo preview`;
}

function formatLiveTimestamp(value: string) {
	const timestamp = new Date(value);
	if (Number.isNaN(timestamp.getTime())) {
		return "Saved run";
	}

	return `${new Intl.DateTimeFormat("en-US", {
		month: "short",
		day: "numeric",
		hour: "numeric",
		minute: "2-digit",
	}).format(timestamp)} - persisted`;
}

function formatBytes(value: number | null | undefined) {
	if (value === null || value === undefined || Number.isNaN(value)) {
		return "n/a";
	}

	if (value < 1024) {
		return `${value} B`;
	}

	if (value < 1024 * 1024) {
		return `${(value / 1024).toFixed(1)} KB`;
	}

	return `${(value / (1024 * 1024)).toFixed(1)} MB`;
}

function humanizeSlug(value: string) {
	return value
		.replace(/[_-]+/g, " ")
		.replace(/\s+/g, " ")
		.trim()
		.replace(/\b\w/g, (match) => match.toUpperCase());
}

function formatLiveStatusLabel(record: BackendRunRecord) {
	if (typeof record.responseStatus === "number") {
		return `HTTP ${record.responseStatus}`;
	}

	switch ((record.status ?? "").toLowerCase()) {
		case "blocked":
			return record.blockedReason ? `Blocked: ${humanizeSlug(record.blockedReason)}` : "Blocked";
		case "timed_out":
			return "Timed out";
		case "canceled":
			return "Canceled";
		case "failed":
			return record.errorCode ? `Failed: ${humanizeSlug(record.errorCode)}` : "Failed";
		default:
			return humanizeSlug(record.status || "Unknown");
	}
}

function buildLiveTitle(record: BackendRunRecord) {
	const target = record.finalUrl || record.url;

	try {
		const parsed = new URL(target);
		if (parsed.pathname && parsed.pathname !== "/") {
			return `${record.method.toUpperCase()} ${parsed.pathname}`;
		}
		return `${record.method.toUpperCase()} ${parsed.hostname}`;
	} catch {
		return `${record.method.toUpperCase()} request`;
	}
}

function buildLaunchTarget(record: BackendRunRecord) {
	if (record.savedRequestId) {
		return {
			launchHref: `/app?request=${encodeURIComponent(record.savedRequestId)}`,
			launchLabel: "Open saved request",
		};
	}

	if (record.collectionId) {
		return {
			launchHref: `/app?collection=${encodeURIComponent(record.collectionId)}`,
			launchLabel: "Open collection",
		};
	}

	return {};
}

function buildPreviewEntry(entry: {
	id: string;
	title: string;
	method: string;
	target: string;
	statusCode: number;
	statusText: string;
	durationMs: number;
	responseSizeLabel: string;
	contentType: string;
	timestampLabel: string;
	source: "demo" | "persistent";
	outcome: HistoryEntryOutcome;
}): HistoryPageEntry {
	return {
		id: entry.id,
		title: entry.title,
		method: entry.method,
		target: entry.target,
		domainLabel: extractDomain(entry.target),
		timeLabel: toPreviewTimeLabel(entry.timestampLabel, entry.source),
		statusLabel: `${entry.statusCode} ${entry.statusText}`,
		durationLabel: `${entry.durationMs} ms`,
		durationMs: entry.durationMs,
		responseSizeLabel: entry.responseSizeLabel,
		contentTypeLabel: entry.contentType,
		outcome: entry.outcome,
	};
}

function buildLiveEntry(record: BackendRunRecord): HistoryPageEntry {
	const target = record.finalUrl || record.url;
	const { launchHref, launchLabel } = buildLaunchTarget(record);

	return {
		id: record.id,
		title: buildLiveTitle(record),
		method: record.method.toUpperCase(),
		target,
		domainLabel: record.targetHost || extractDomain(target),
		timeLabel: formatLiveTimestamp(record.createdAt),
		statusLabel: formatLiveStatusLabel(record),
		durationLabel:
			typeof record.responseTimeMs === "number" && record.responseTimeMs >= 0
				? `${record.responseTimeMs} ms`
				: "n/a",
		durationMs: typeof record.responseTimeMs === "number" && record.responseTimeMs >= 0 ? record.responseTimeMs : 0,
		responseSizeLabel: formatBytes(record.responseSizeBytes),
		contentTypeLabel: record.responseContentType?.trim() || "Unknown",
		outcome: outcomeFromRecord(record),
		launchHref,
		launchLabel,
	};
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
		.map((entry) =>
			buildPreviewEntry({
				...entry,
				outcome: entry.outcome,
			}),
		)
		.filter((entry) => matchesStatus(entry, status))
		.filter((entry) => matchesMethod(entry, method))
		.filter((entry) => matchesDomain(entry, selectedDomain));

	return {
		mode,
		source: "preview",
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

export function buildLiveHistoryPageData(
	mode: WorkspaceMode,
	payload: BackendHistoryListResponse,
	statusParam: string | null | undefined,
	methodParam: string | null | undefined,
	domainParam: string | null | undefined,
): HistoryPageData {
	const status = normalizeStatus(statusParam);
	const method = normalizeMethod(methodParam);
	const selectedDomain = domainParam?.trim() || "all";
	const entries = (payload.history ?? []).map(buildLiveEntry);
	const domainOptions = Array.from(new Set(entries.map((entry) => entry.domainLabel).filter(Boolean))).sort();

	if (selectedDomain !== "all" && !domainOptions.includes(selectedDomain)) {
		domainOptions.push(selectedDomain);
		domainOptions.sort();
	}

	return {
		mode,
		source: "live",
		status,
		method,
		selectedDomain,
		domainOptions,
		entries,
		sections: buildSections(entries),
		metrics: buildMetrics(entries),
		previewLocked: false,
		pagination: payload.pagination,
	};
}

export function buildUnavailableHistoryPageData(
	mode: WorkspaceMode,
	statusParam: string | null | undefined,
	methodParam: string | null | undefined,
	domainParam: string | null | undefined,
): HistoryPageData {
	const status = normalizeStatus(statusParam);
	const method = normalizeMethod(methodParam);
	const selectedDomain = domainParam?.trim() || "all";

	return {
		mode,
		source: "live",
		status,
		method,
		selectedDomain,
		domainOptions: selectedDomain === "all" ? [] : [selectedDomain],
		entries: [],
		sections: [],
		metrics: buildMetrics([]),
		previewLocked: false,
		notice: "History is temporarily unavailable, so the signed-in timeline could not be loaded.",
	};
}

function matchesStatus(entry: HistoryPageEntry, status: HistoryStatusFilter) {
	if (status === "all") {
		return true;
	}

	return entry.outcome === status;
}

function matchesMethod(entry: HistoryPageEntry, method: HistoryMethodFilter) {
	if (method === "all") {
		return true;
	}

	return entry.method === method;
}

function matchesDomain(entry: HistoryPageEntry, domain: string) {
	if (!domain || domain === "all") {
		return true;
	}

	return entry.domainLabel === domain;
}
