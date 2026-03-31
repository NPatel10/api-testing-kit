<script lang="ts">
	import { page } from "$app/state";
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import GuestAdvancedToolsLock from "$lib/components/workspace/guest-advanced-tools-lock.svelte";
	import GuestCustomUrlLock from "$lib/components/workspace/guest-custom-url-lock.svelte";
	import GuestEnvVarsLock from "$lib/components/workspace/guest-env-vars-lock.svelte";
	import GuestHistoryLock from "$lib/components/workspace/guest-history-lock.svelte";
	import GuestSaveLock from "$lib/components/workspace/guest-save-lock.svelte";
	import RequestBuilder from "$lib/components/workspace/request-builder.svelte";
	import {
		createDefaultRequestDraft,
		type RequestBodyMode,
		type RequestBuilderDraft,
	} from "$lib/components/workspace/request-builder";
	import { createSavedRequestDraft } from "$lib/components/workspace/saved-request";
	import ResponseViewer from "$lib/components/workspace/response-viewer.svelte";
	import type { ResponseHeader } from "$lib/components/workspace/response-viewer";
	import {
		buildLiveRunPayload,
		createEmptyViewerState,
		createPreviewViewerState,
		createRequestFailureState,
		readLiveRunViewerState,
		type LiveRunViewerState,
	} from "$lib/components/workspace/live-request";
	import TemplateBrowser from "$lib/components/workspace/template-browser.svelte";
	import { buildEntitlementRows, getEntitlementSummary, type EffectiveEntitlements } from "$lib/entitlements/access";
	import {
		authenticatedWorkspaceState,
		guestWorkspaceState,
		type WorkspaceCollectionPreview,
		type WorkspaceMode,
		type WorkspaceTemplate,
	} from "$lib/mocks/workspace-state";
	import type { BackendSavedRequest } from "$lib/server/backend";

	type Props = {
		mode?: WorkspaceMode;
		entitlements?: EffectiveEntitlements;
		savedRequest?: BackendSavedRequest | null;
	};

	let { mode = "guest", entitlements, savedRequest = null }: Props = $props();
	const workspaceState = $derived(mode === "authenticated" ? authenticatedWorkspaceState : guestWorkspaceState);
	const previewMode = $derived(page.url.searchParams.get("mode") === "preview");
	const selectedTemplate = $derived(
		selectWorkspaceTemplate(
			workspaceState.templates,
			workspaceState.collections,
			page.url.searchParams.get("template"),
			page.url.searchParams.get("collection"),
		),
	);
	const previewSnapshot = $derived(selectedTemplate ? buildPreviewSnapshot(selectedTemplate) : null);
	const previewResponse = $derived(previewMode && previewSnapshot ? createPreviewViewerState(previewSnapshot) : null);
	const emptyResponse = createEmptyViewerState();
	let liveResponse = $state<LiveRunViewerState | null>(null);
	let isSending = $state(false);
	let activeSelectionKey = "";
	let requestNonce = 0;
	let activeAbortController: AbortController | null = null;

	const categoryMap = {
		"REST basics": "rest-basics",
		"Authentication flows": "auth-flows",
		"CRUD examples": "crud",
		"Pagination examples": "pagination",
		Webhooks: "webhooks",
		"Error handling": "error-handling",
	} as const;

	const metricToneClasses = {
		neutral: "border-border/70 bg-panel-soft text-text-strong",
		positive: "border-success/20 bg-success/10 text-success",
		warning: "border-warning/20 bg-warning/10 text-warning",
		danger: "border-danger/20 bg-danger/10 text-danger",
	} as const;

	const historyToneClasses = {
		success: "border-success/20 bg-success/10 text-success",
		blocked: "border-warning/20 bg-warning/10 text-warning",
		error: "border-danger/20 bg-danger/10 text-danger",
	} as const;

	const currentSelectionKey = $derived(
		`${mode}:${previewMode ? "preview" : "live"}:${savedRequest?.id ?? selectedTemplate?.slug ?? "default"}`,
	);
	const liveResponseState = $derived(liveResponse ?? previewResponse ?? emptyResponse);
	const liveResponseDescription = $derived(
		liveResponse
			? "Latest backend response."
			: savedRequest
				? "Persisted request state."
			: previewResponse
				? "Seeded preview response."
				: "No response yet.",
	);
	const liveEmptyTitle = $derived(
		isSending ? "Sending request..." : previewMode ? "Preview response pending" : "No request sent yet",
	);
	const liveEmptyDescription = $derived(
		isSending
			? "Waiting for the current run to finish."
			: previewMode
				? "Preview response is available for the current selection."
				: "Send a request to populate the viewer.",
	);

	$effect(() => {
		if (currentSelectionKey === activeSelectionKey) {
			return;
		}

		activeSelectionKey = currentSelectionKey;
		liveResponse = null;
		isSending = false;
		activeAbortController?.abort();
		activeAbortController = null;
	});

	function toRequestBodyMode(value: string): RequestBodyMode {
		if (value === "raw") {
			return "raw";
		}

		if (value === "form-urlencoded") {
			return "form";
		}

		return "json";
	}

	function selectWorkspaceTemplate(
		templates: readonly WorkspaceTemplate[],
		collections: readonly WorkspaceCollectionPreview[],
		templateSlug: string | null,
		collectionId: string | null,
	): WorkspaceTemplate | undefined {
		const explicitTemplate = templateSlug?.trim()
			? templates.find((template) => template.slug === templateSlug.trim())
			: undefined;
		if (explicitTemplate) {
			return explicitTemplate;
		}

		const collectionTemplateSlug = collectionId?.trim()
			? collections
					.find((collection) => collection.id === collectionId.trim())
					?.templateSlugs.find((slug) => templates.some((template) => template.slug === slug))
			: undefined;

		if (collectionTemplateSlug) {
			return templates.find((template) => template.slug === collectionTemplateSlug);
		}

		return templates[0];
	}

	function createTemplateRequestDraft(currentMode: WorkspaceMode, template?: WorkspaceTemplate): RequestBuilderDraft {
		const draft = createDefaultRequestDraft(currentMode);

		if (!template) {
			return draft;
		}

		draft.method = template.request.method;
		draft.url = template.request.url;
		draft.queryParams = template.request.query.map((item) => ({
			key: item.key,
			value: item.value,
			enabled: true,
		}));
		draft.headers = template.request.headers.map((item) => ({
			key: item.key,
			value: item.value,
			enabled: true,
		}));

		const bodyMode = toRequestBodyMode(template.request.bodyMode);
		draft.body = {
			...draft.body,
			mode: bodyMode,
			value:
				bodyMode === "json"
					? '{\n  "template": "' + template.slug + '",\n  "preview": true\n}'
					: bodyMode === "raw"
						? "demo-preview-body"
						: draft.body.value,
			formRows:
				bodyMode === "form"
					? [
							{ key: "email", value: "guest@example.dev", enabled: true },
							{ key: "city", value: "Kolkata", enabled: true },
						]
					: draft.body.formRows,
			contentType:
				bodyMode === "form"
					? "application/x-www-form-urlencoded"
					: bodyMode === "raw"
						? "text/plain"
						: "application/json",
		};

		return draft;
	}

	const requestDraft = $derived(
		savedRequest ? createSavedRequestDraft(mode, savedRequest) : createTemplateRequestDraft(mode, selectedTemplate),
	);

	const responseHeaders: ResponseHeader[] = $derived(liveResponseState.headers);

	const templateBrowserTemplates = $derived(
		workspaceState.templates.map((template, index) => ({
			id: template.slug,
			name: template.title,
			slug: template.slug,
			category: categoryMap[template.category],
			method: template.request.method,
			endpoint: template.request.url,
			summary: template.summary,
			notes: template.description,
			tags: [...template.tags],
			featured: index === 0,
			launchHref: `/app?template=${template.slug}`,
			previewHref: `/app?template=${template.slug}&mode=preview`,
		}))
	);

	const templateBrowserCollections = $derived(
		workspaceState.collections.map((collection) => ({
			id: collection.id,
			name: collection.title,
			slug: collection.id,
			category: categoryMap[
				workspaceState.templateGroups.find((group) =>
					group.templateSlugs.some((slug) => collection.templateSlugs.includes(slug))
				)?.label ?? "REST basics"
			],
			description: collection.description,
			templateIds: [...collection.templateSlugs],
			launchHref: `/app?collection=${collection.id}`,
			previewHref: `/app?collection=${collection.id}&mode=preview`,
		}))
	);

	async function handleSend(draft: RequestBuilderDraft) {
		const requestId = ++requestNonce;
		activeAbortController?.abort();
		const controller = new AbortController();
		activeAbortController = controller;
		isSending = true;

		try {
			const response = await fetch(getRunEndpoint(mode), {
				method: "POST",
				headers: {
					accept: "application/json",
					"content-type": "application/json",
				},
				body: JSON.stringify(buildLiveRunPayload(draft)),
				signal: controller.signal,
			});

			const nextResponse = await readLiveRunViewerState(response);
			if (requestId === requestNonce) {
				liveResponse = nextResponse;
			}
		} catch (error) {
			if (requestId !== requestNonce) {
				return;
			}

			if (error instanceof DOMException && error.name === "AbortError") {
				return;
			}

			liveResponse = createRequestFailureState(error instanceof Error ? error.message : "Unexpected request failure");
		} finally {
			if (requestId === requestNonce) {
				isSending = false;
				if (activeAbortController === controller) {
					activeAbortController = null;
				}
			}
		}
	}

	function buildPreviewSnapshot(template: WorkspaceTemplate) {
		return {
			responseStatus: template.request.responseStatus,
			responseStatusText: template.request.responseStatusText,
			responseTimeMs: template.request.responseTimeMs,
			responseSizeLabel: template.request.responseSizeLabel,
			responseContentType: template.request.responseContentType,
			responseBody: template.request.responseBody,
		};
	}

	function getRunEndpoint(currentMode: WorkspaceMode) {
		return currentMode === "guest" ? "/api/v1/guest-runs" : "/api/v1/runs";
	}
</script>

<section class="space-y-4">
	{#if entitlements}
		<Card class="panel-card">
			<CardHeader class="gap-3">
				<div class="flex items-center justify-between gap-3">
					<div>
						<CardTitle>{entitlements.plan.name}</CardTitle>
						<CardDescription>{getEntitlementSummary(entitlements, mode)}</CardDescription>
					</div>
					<Badge variant="outline">{entitlements.plan.source}</Badge>
				</div>
			</CardHeader>
			<CardContent class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
				{#each buildEntitlementRows(entitlements) as row}
					<div class={`rounded-[20px] border px-4 py-4 ${row.tone === "positive" ? "border-success/20 bg-success/10 text-success" : "border-warning/20 bg-warning/10 text-warning"}`}>
						<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">{row.label}</p>
						<p class="mt-2 text-base font-semibold tracking-tight text-current">{row.statusLabel}</p>
						<p class="mt-1 text-sm text-text-body">{row.description}</p>
						{#if row.limitLabel}
							<p class="mt-2 text-xs font-medium text-text-muted">{row.limitLabel}</p>
						{/if}
					</div>
				{/each}
			</CardContent>
		</Card>
	{/if}

	<div class="grid gap-4 xl:grid-cols-[1.18fr_0.95fr]">
		<RequestBuilder
			title="Request builder"
			mode={mode}
			description={workspaceState.accessSummary}
			request={requestDraft}
			lockedNote={workspaceState.prompts[0]?.body}
			pending={isSending}
			sendLabel={mode === "guest" ? "Send guest request" : "Send request"}
			onSend={handleSend}
		/>

		<ResponseViewer
			title="Response viewer"
			description={liveResponseDescription}
			status={liveResponseState.status}
			statusText={liveResponseState.statusText}
			duration={liveResponseState.duration}
			size={liveResponseState.size}
			contentType={liveResponseState.contentType}
			headers={responseHeaders}
			prettyBody={liveResponseState.prettyBody}
			rawBody={liveResponseState.rawBody}
			error={liveResponseState.error}
			emptyTitle={liveEmptyTitle}
			emptyDescription={liveEmptyDescription}
		/>
	</div>

	<div class="grid gap-4 xl:grid-cols-[1.16fr_0.84fr]">
		<TemplateBrowser
			title="Templates and collections"
			subtitle={workspaceState.subtitle}
			templates={templateBrowserTemplates}
			collections={templateBrowserCollections}
		/>

		<div class="space-y-4">
			<Card class="panel-card">
				<CardHeader class="gap-3">
					<div class="flex items-center justify-between gap-3">
						<div>
							<CardTitle>Session state</CardTitle>
							<CardDescription>
								{mode === "authenticated"
									? "Signed-in access and quota state."
									: "Guest access and lock-state labels."}
							</CardDescription>
						</div>
						<Badge variant={mode === "authenticated" ? "default" : "secondary"}>{workspaceState.mode}</Badge>
					</div>
				</CardHeader>
				<CardContent class="space-y-3">
					{#each workspaceState.prompts as prompt}
						<div class="rounded-[20px] border border-border/70 bg-panel-soft p-4">
							<div class="flex items-center justify-between gap-3">
								<p class="text-sm font-semibold text-text-strong">{prompt.title}</p>
								<Badge
									variant="outline"
									class={metricToneClasses[prompt.tone]}
								>
									{prompt.action.label}
								</Badge>
							</div>
						</div>
					{/each}
				</CardContent>
			</Card>

			<Card class="panel-card">
				<CardHeader class="gap-3">
					<CardTitle>{mode === "authenticated" ? "Recent authenticated runs" : "Recent guest runs"}</CardTitle>
					<CardDescription>
						{mode === "authenticated"
							? "Recent persisted runs."
							: "Recent preview runs."}
					</CardDescription>
				</CardHeader>
				<CardContent class="space-y-3">
					{#each workspaceState.history.slice(0, 4) as entry}
						<div class="rounded-[18px] border border-border/70 bg-panel-soft px-4 py-3">
							<div class="flex items-start justify-between gap-3">
								<div>
									<p class="text-sm font-semibold text-text-strong">{entry.title}</p>
									<p class="mt-1 font-mono text-xs text-text-muted">{entry.target}</p>
								</div>
								<Badge variant="outline" class={historyToneClasses[entry.outcome]}>
									{entry.statusCode} {entry.statusText}
								</Badge>
							</div>
							<div class="mt-3 flex flex-wrap gap-2 text-xs text-text-muted">
								<span class="rounded-full border border-border/70 bg-white px-3 py-1">
									{entry.durationMs} ms
								</span>
								<span class="rounded-full border border-border/70 bg-white px-3 py-1">
									{entry.responseSizeLabel}
								</span>
								<span class="rounded-full border border-border/70 bg-white px-3 py-1">
									{entry.timestampLabel}
								</span>
							</div>
						</div>
					{/each}
				</CardContent>
			</Card>

			<Card class="panel-card">
				<CardHeader class="gap-3">
					<CardTitle>Quota snapshot</CardTitle>
					<CardDescription>
						{mode === "authenticated"
							? "Current signed-in limits."
							: "Current guest limits."}
					</CardDescription>
				</CardHeader>
				<CardContent class="space-y-3">
					{#each workspaceState.quotas.slice(0, 3) as quota}
						<div class="rounded-[18px] border border-border/70 bg-panel-soft px-4 py-3">
							<div class="flex items-center justify-between gap-3">
								<p class="text-sm font-semibold text-text-strong">{quota.label}</p>
								<Badge variant="outline">{quota.remainingLabel}</Badge>
							</div>
							<p class="mt-2 text-sm text-text-body">{quota.usedLabel} of {quota.limitLabel}</p>
							<p class="mt-1 text-xs leading-5 text-text-muted">{quota.note}</p>
						</div>
					{/each}
				</CardContent>
			</Card>
		</div>
	</div>

	{#if mode === "guest"}
		<div class="grid gap-4 lg:grid-cols-2 2xl:grid-cols-3">
			<GuestCustomUrlLock />
			<GuestSaveLock />
			<GuestHistoryLock />
			<GuestEnvVarsLock />
			<GuestAdvancedToolsLock />
		</div>
	{:else}
		<div class="grid gap-4 lg:grid-cols-2 2xl:grid-cols-3">
			<Card class="panel-card">
				<CardHeader class="gap-3">
					<CardTitle>Access</CardTitle>
					<CardDescription>Signed-in capabilities on the shared route.</CardDescription>
				</CardHeader>
				<CardContent class="space-y-2 text-sm leading-6 text-text-body">
					<p>Custom URLs, saved requests, and history are enabled.</p>
					<p>Outbound validation remains active.</p>
				</CardContent>
			</Card>
			<Card class="panel-card">
				<CardHeader class="gap-3">
					<CardTitle>Unlocked surfaces</CardTitle>
					<CardDescription>These controls stay visible in the authenticated workspace.</CardDescription>
				</CardHeader>
				<CardContent class="flex flex-wrap gap-2">
					{#each ["Custom URL", "Save request", "History", "Environment variables", "Advanced tools"] as item}
						<Badge variant="outline">{item}</Badge>
					{/each}
				</CardContent>
			</Card>
			<Card class="panel-card">
				<CardHeader class="gap-3">
					<CardTitle>Quota model</CardTitle>
					<CardDescription>{workspaceState.lockedActions[0]?.description ?? "Higher-tier entitlements can layer on later."}</CardDescription>
				</CardHeader>
				<CardContent class="text-sm leading-6 text-text-body">
					<p>The shared workspace can carry broader plan limits later.</p>
				</CardContent>
			</Card>
		</div>
	{/if}

	<Card class="panel-card">
		<CardHeader class="gap-3">
			<div class="flex items-center justify-between gap-3">
				<div>
					<CardTitle>Workspace pulse</CardTitle>
					<CardDescription>
						{mode === "authenticated"
							? "Signed-in workspace metrics."
							: "Guest workspace metrics."}
					</CardDescription>
				</div>
				<Badge variant="outline">{workspaceState.mode} mode</Badge>
			</div>
		</CardHeader>
		<CardContent class="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
			{#each workspaceState.metrics as metric}
				<div class={`rounded-[20px] border px-4 py-4 ${metricToneClasses[metric.tone]}`}>
					<p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">
						{metric.label}
					</p>
					<p class="mt-2 text-2xl font-semibold tracking-tight text-current">{metric.value}</p>
					<p class="mt-1 text-sm text-text-body">{metric.detail}</p>
				</div>
			{/each}
		</CardContent>
	</Card>
</section>
