<script lang="ts">
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import { Tabs, TabsContent, TabsList, TabsTrigger } from "$lib/components/ui/tabs/index.js";
	import { cn } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";
	import {
		formatBytes,
		formatDuration,
		formatStatusLabel,
		getResponseTone,
		getResponseToneLabel,
		hasBodyContent,
		normalizeHeaders,
		type ResponseHeader,
		type ResponseStatusTone,
		type ResponseViewerError,
		type ResponseViewerTab,
	} from "./response-viewer.js";

	type Props = HTMLAttributes<HTMLElement> & {
		title?: string;
		description?: string;
		activeTab?: ResponseViewerTab;
		status?: number;
		statusText?: string;
		duration?: number | string;
		size?: number | string;
		contentType?: string;
		headers?: ResponseHeader[];
		prettyBody?: string;
		rawBody?: string;
		error?: ResponseViewerError | null;
		emptyTitle?: string;
		emptyDescription?: string;
	};

	let {
		class: className,
		title = "Response viewer",
		description = "Formatted output, headers, and request metadata from the latest run.",
		activeTab = $bindable("pretty"),
		status,
		statusText,
		duration,
		size,
		contentType = "--",
		headers = [],
		prettyBody = "",
		rawBody = "",
		error = null,
		emptyTitle = "No response body yet",
		emptyDescription = "Send a request to populate the pretty, raw, and header views.",
		...restProps
	}: Props = $props();

	const normalizedHeaders = $derived(normalizeHeaders(headers));
	const responseTone = $derived(getResponseTone(status, error));
	const responseToneLabel = $derived(getResponseToneLabel(status, error));
	const statusLabel = $derived(formatStatusLabel(status, statusText, error));
	const durationLabel = $derived(formatDuration(duration));
	const sizeLabel = $derived(formatBytes(size));
	const prettyPreview = $derived(prettyBody.trim() || rawBody.trim());
	const rawPreview = $derived(rawBody.trim() || prettyBody.trim());
	const hasPrettyBody = $derived(hasBodyContent(prettyPreview));
	const hasRawBody = $derived(hasBodyContent(rawPreview));
	const hasHeaders = $derived(normalizedHeaders.length > 0);

	const toneClasses: Record<ResponseStatusTone, string> = {
		success: "border-success/20 bg-success/10 text-success",
		warning: "border-warning/20 bg-warning/10 text-warning",
		danger: "border-danger/20 bg-danger/10 text-danger",
		neutral: "border-border/70 bg-panel-soft text-text-body",
	};

	const metaCards = $derived([
		{ label: "Status", value: statusLabel, tone: responseTone },
		{ label: "Duration", value: durationLabel, tone: "neutral" as const },
		{ label: "Size", value: sizeLabel, tone: "neutral" as const },
		{ label: "Content type", value: contentType || "--", tone: "neutral" as const },
	]);

	const errorMessage = $derived(error?.message?.trim() || "");
	const errorDetails = $derived(error?.details?.trim() || "");
	const errorCode = $derived(error?.code?.trim() || "");
</script>

<Card class={cn("rounded-[24px] border-border/80 bg-card shadow-card", className)} {...restProps}>
	<CardHeader class="gap-4 border-b border-border/70 px-5 py-5 md:px-6">
		<div class="flex flex-wrap items-start justify-between gap-4">
			<div class="space-y-1">
				<CardTitle class="text-lg font-semibold tracking-tight text-text-strong">{title}</CardTitle>
				<CardDescription class="max-w-2xl text-sm leading-6 text-text-body">
					{description}
				</CardDescription>
			</div>

			<div class="flex items-center gap-2">
				<Badge variant={responseTone === "danger" ? "destructive" : responseTone === "success" ? "default" : "outline"} class={cn("border px-3", toneClasses[responseTone])}>
					{responseToneLabel}
				</Badge>
			</div>
		</div>

		<div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
			{#each metaCards as meta}
				<div class={cn("rounded-[18px] border border-border/70 px-4 py-3", meta.label === "Status" && toneClasses[meta.tone])}>
					<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">{meta.label}</p>
					<p class="mt-2 break-words text-sm font-medium text-text-strong">{meta.value}</p>
				</div>
			{/each}
		</div>
	</CardHeader>

	<CardContent class="space-y-4 px-5 py-5 md:px-6">
		{#if error}
			<div class="rounded-[20px] border border-danger/20 bg-danger/10 p-4">
				<div class="flex flex-wrap items-center justify-between gap-3">
					<div>
						<p class="text-sm font-semibold text-danger">{error.title}</p>
						<p class="mt-1 text-sm leading-6 text-text-body">{errorMessage}</p>
					</div>
					{#if errorCode}
						<Badge variant="destructive" class="bg-danger/10 text-danger">{errorCode}</Badge>
					{/if}
				</div>
				{#if errorDetails}
					<p class="mt-3 text-xs leading-5 text-text-muted">{errorDetails}</p>
				{/if}
			</div>
		{/if}

		<Tabs bind:value={activeTab} class="gap-4">
			<TabsList>
				<TabsTrigger value="pretty">Pretty</TabsTrigger>
				<TabsTrigger value="raw">Raw</TabsTrigger>
				<TabsTrigger value="headers">Headers</TabsTrigger>
			</TabsList>

			<TabsContent value="pretty" class="space-y-3">
				{#if hasPrettyBody}
					<pre class="code-surface max-h-[34rem] overflow-auto whitespace-pre-wrap break-words"><code>{prettyPreview}</code></pre>
				{:else}
					<div class="rounded-[20px] border border-border/70 bg-panel-soft p-5">
						<p class="text-sm font-semibold text-text-strong">{emptyTitle}</p>
						<p class="mt-1 text-sm leading-6 text-text-body">{emptyDescription}</p>
					</div>
				{/if}
			</TabsContent>

			<TabsContent value="raw" class="space-y-3">
				{#if hasRawBody}
					<pre class="code-surface max-h-[34rem] overflow-auto whitespace-pre"><code>{rawPreview}</code></pre>
				{:else}
					<div class="rounded-[20px] border border-border/70 bg-panel-soft p-5">
						<p class="text-sm font-semibold text-text-strong">{emptyTitle}</p>
						<p class="mt-1 text-sm leading-6 text-text-body">The raw payload will appear here when a response body is present.</p>
					</div>
				{/if}
			</TabsContent>

			<TabsContent value="headers" class="space-y-3">
				{#if hasHeaders}
					<div class="space-y-2">
						{#each normalizedHeaders as header}
							<div class="rounded-[18px] border border-border/70 bg-surface-soft px-4 py-3">
								<div class="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
									<div class="space-y-1">
										<p class="text-[11px] font-semibold uppercase tracking-[0.2em] text-text-muted">Header</p>
										<p class="font-mono text-sm font-medium text-text-strong">{header.key}</p>
									</div>
									<p class="max-w-full break-all font-mono text-sm text-text-body md:max-w-[55%] md:text-right">{header.value}</p>
								</div>
							</div>
						{/each}
					</div>
				{:else}
					<div class="rounded-[20px] border border-border/70 bg-panel-soft p-5">
						<p class="text-sm font-semibold text-text-strong">No headers captured</p>
						<p class="mt-1 text-sm leading-6 text-text-body">
							Header metadata will be shown here alongside the body tabs for quick inspection.
						</p>
					</div>
				{/if}
			</TabsContent>
		</Tabs>
	</CardContent>
</Card>
