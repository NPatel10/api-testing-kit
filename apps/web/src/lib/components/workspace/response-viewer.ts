export type ResponseViewerTab = "pretty" | "raw" | "headers";

export type ResponseHeader = {
	key: string;
	value: string;
};

export type ResponseViewerError = {
	title: string;
	message: string;
	code?: string;
	details?: string;
};

export type ResponseStatusTone = "success" | "warning" | "danger" | "neutral";

export function formatDuration(duration?: number | string): string {
	if (duration === undefined || duration === null || duration === "") {
		return "--";
	}

	if (typeof duration === "string") {
		return duration;
	}

	if (!Number.isFinite(duration)) {
		return "--";
	}

	if (duration < 1000) {
		return `${Math.round(duration)} ms`;
	}

	const seconds = duration / 1000;
	return `${seconds >= 10 ? seconds.toFixed(0) : seconds.toFixed(1)} s`;
}

export function formatBytes(size?: number | string): string {
	if (size === undefined || size === null || size === "") {
		return "--";
	}

	if (typeof size === "string") {
		return size;
	}

	if (!Number.isFinite(size)) {
		return "--";
	}

	if (size < 1024) {
		return `${Math.round(size)} B`;
	}

	const units = ["KB", "MB", "GB", "TB"];
	let value = size / 1024;
	let unit = units[0];

	for (let index = 0; index < units.length; index += 1) {
		unit = units[index];
		if (value < 1024 || index === units.length - 1) {
			break;
		}
		value /= 1024;
	}

	return `${value >= 10 ? value.toFixed(0) : value.toFixed(1)} ${unit}`;
}

export function formatStatusLabel(
	status?: number,
	statusText?: string,
	error?: ResponseViewerError | null
): string {
	if (error) {
		return error.code ? `${error.title} (${error.code})` : error.title;
	}

	if (status === undefined || status === null || Number.isNaN(status)) {
		return "Idle";
	}

	if (statusText?.trim()) {
		return `${status} ${statusText}`;
	}

	return String(status);
}

export function getResponseTone(
	status?: number,
	error?: ResponseViewerError | null
): ResponseStatusTone {
	if (error) {
		return "danger";
	}

	if (status === undefined || status === null || Number.isNaN(status)) {
		return "neutral";
	}

	if (status >= 200 && status < 300) {
		return "success";
	}

	if (status >= 300 && status < 400) {
		return "warning";
	}

	if (status >= 400) {
		return "danger";
	}

	return "neutral";
}

export function getResponseToneLabel(
	status?: number,
	error?: ResponseViewerError | null
): string {
	if (error) {
		return "Error";
	}

	if (status === undefined || status === null || Number.isNaN(status)) {
		return "Waiting";
	}

	if (status >= 200 && status < 300) {
		return "Success";
	}

	if (status >= 300 && status < 400) {
		return "Redirect";
	}

	if (status >= 400 && status < 500) {
		return "Client error";
	}

	if (status >= 500) {
		return "Server error";
	}

	return "Ready";
}

export function normalizeHeaders(headers: ResponseHeader[]): ResponseHeader[] {
	return headers
		.map((header) => ({
			key: header.key.trim(),
			value: header.value.trim(),
		}))
		.filter((header) => header.key.length > 0 || header.value.length > 0);
}

export function hasBodyContent(body?: string): boolean {
	return Boolean(body?.trim());
}
