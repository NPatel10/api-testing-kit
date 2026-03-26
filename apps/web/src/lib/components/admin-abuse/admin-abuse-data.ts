export interface AdminAbuseSummaryRow {
	severity: string;
	category: string;
	actionTaken: string;
	count: number;
	lastCreatedAt: string;
}

export interface AdminAbuseEvent {
	id: string;
	ruleKey: string;
	category: string;
	severity: string;
	actionTaken: string;
	message: string;
	target?: string | null;
	sourceIp?: string | null;
	requestId?: string | null;
	createdAt: string;
}

export interface AdminBlockedTarget {
	id: string;
	targetType: string;
	targetValue: string;
	reason: string;
	source: string;
	isActive: boolean;
	expiresAt?: string | null;
	createdByUserId?: string | null;
	metadata?: Record<string, unknown> | null;
	createdAt: string;
	updatedAt: string;
}

export interface AdminAbuseApiResponse {
	generatedAt: string;
	summary: AdminAbuseSummaryRow[];
	recent: AdminAbuseEvent[];
	blockedTargets: AdminBlockedTarget[];
}

export interface AdminAbuseDashboardData {
	mode: "live" | "preview";
	sourceLabel: string;
	message: string;
	generatedAt: string;
	summary: AdminAbuseSummaryRow[];
	recent: AdminAbuseEvent[];
	blockedTargets: AdminBlockedTarget[];
}

const previewSummary: AdminAbuseSummaryRow[] = [
	{
		severity: "high",
		category: "blocked_host",
		actionTaken: "blocked",
		count: 4,
		lastCreatedAt: "2026-03-26T08:00:00Z",
	},
];

const previewRecent: AdminAbuseEvent[] = [
	{
		id: "preview-abuse-1",
		ruleKey: "blocked-target",
		category: "blocked_host",
		severity: "high",
		actionTaken: "blocked",
		message: "Repeated attempts against blocked targets were observed.",
		target: "169.254.169.254",
		sourceIp: "203.0.113.42",
		requestId: "preview-request-1",
		createdAt: "2026-03-26T08:01:00Z",
	},
];

const previewBlockedTargets: AdminBlockedTarget[] = [
	{
		id: "preview-blocked-1",
		targetType: "ip",
		targetValue: "169.254.169.254",
		reason: "metadata IP",
		source: "manual",
		isActive: true,
		metadata: {
			note: "Preview data only. Live backend data appears for admin sessions.",
		},
		createdAt: "2026-03-26T07:00:00Z",
		updatedAt: "2026-03-26T07:00:00Z",
	},
];

export function normalizeAdminAbuseDashboard(
	response: Partial<AdminAbuseApiResponse>,
): AdminAbuseDashboardData {
	return {
		mode: "live",
		sourceLabel: "Live backend",
		message: "Live abuse monitoring data loaded from the Go API.",
		generatedAt: response.generatedAt ?? new Date().toISOString(),
		summary: response.summary ?? [],
		recent: response.recent ?? [],
		blockedTargets: response.blockedTargets ?? [],
	};
}

export function buildPreviewAdminAbuseDashboard(message: string): AdminAbuseDashboardData {
	return {
		mode: "preview",
		sourceLabel: "Preview data",
		message,
		generatedAt: new Date().toISOString(),
		summary: previewSummary,
		recent: previewRecent,
		blockedTargets: previewBlockedTargets,
	};
}

export function formatAdminTimestamp(value: string | undefined | null): string {
	if (!value) {
		return "Unknown";
	}

	const date = new Date(value);
	if (Number.isNaN(date.getTime())) {
		return value;
	}

	return new Intl.DateTimeFormat("en-US", {
		dateStyle: "medium",
		timeStyle: "short",
		timeZone: "UTC",
	}).format(date);
}
