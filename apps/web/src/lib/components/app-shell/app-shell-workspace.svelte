<script lang="ts">
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
	import ResponseViewer from "$lib/components/workspace/response-viewer.svelte";
	import type { ResponseHeader } from "$lib/components/workspace/response-viewer";
	import TemplateBrowser from "$lib/components/workspace/template-browser.svelte";
	import { buildEntitlementRows, getEntitlementSummary, type EffectiveEntitlements } from "$lib/entitlements/access";
	import { authenticatedWorkspaceState, guestWorkspaceState } from "$lib/mocks/workspace-state";
	import type { WorkspaceMode } from "$lib/mocks/workspace-state";

	type Props = {
		mode?: WorkspaceMode;
		entitlements?: EffectiveEntitlements;
	};

	let { mode = "guest", entitlements }: Props = $props();
	const workspaceState = $derived(mode === "authenticated" ? authenticatedWorkspaceState : guestWorkspaceState);
	const primaryTemplate = $derived(workspaceState.templates[0]);

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

	function toRequestBodyMode(mode: string): RequestBodyMode {
		if (mode === "raw") {
			return "raw";
		}

		if (mode === "form-urlencoded") {
			return "form";
		}

		return "json";
	}

	function createTemplateRequestDraft(currentMode: WorkspaceMode): RequestBuilderDraft {
		const draft = createDefaultRequestDraft(currentMode);

		if (!primaryTemplate) {
			return draft;
		}

		draft.method = primaryTemplate.request.method;
		draft.url = primaryTemplate.request.url;
		draft.queryParams = primaryTemplate.request.query.map((item) => ({
			key: item.key,
			value: item.value,
			enabled: true,
		}));
		draft.headers = primaryTemplate.request.headers.map((item) => ({
			key: item.key,
			value: item.value,
			enabled: true,
		}));

		const bodyMode = toRequestBodyMode(primaryTemplate.request.bodyMode);
		draft.body = {
			...draft.body,
			mode: bodyMode,
			value:
				bodyMode === "json"
					? '{\n  "template": "' + primaryTemplate.slug + '",\n  "preview": true\n}'
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

	const requestDraft = $derived(createTemplateRequestDraft(mode));

	const responseHeaders: ResponseHeader[] = $derived(
		primaryTemplate
			? [
					{ key: "content-type", value: primaryTemplate.request.responseContentType },
					{ key: "x-guest-mode", value: mode === "authenticated" ? "authenticated-template" : "allowlisted-template" },
					{ key: "x-preview-size", value: primaryTemplate.request.responseSizeLabel },
				]
			: []
	);

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
		/>

		<ResponseViewer
			title="Response viewer"
			description={mode === "authenticated" ? "Previewed authenticated responses stay structured, readable, and safe." : "Previewed guest responses stay structured, readable, and visibly constrained."}
			status={primaryTemplate?.request.responseStatus}
			statusText={primaryTemplate?.request.responseStatusText}
			duration={primaryTemplate?.request.responseTimeMs}
			size={primaryTemplate?.request.responseSizeLabel}
			contentType={primaryTemplate?.request.responseContentType}
			headers={responseHeaders}
			prettyBody={primaryTemplate?.request.responseBody}
			rawBody={primaryTemplate?.request.responseBody}
		/>
	</div>

	<div class="grid gap-4 xl:grid-cols-[1.16fr_0.84fr]">
		<TemplateBrowser
			title={mode === "authenticated" ? "Templates, collections, and unlocked workspace state" : "Guest-safe templates and collections"}
			subtitle={workspaceState.subtitle}
			templates={templateBrowserTemplates}
			collections={templateBrowserCollections}
		/>

		<div class="space-y-4">
			<Card class="panel-card">
				<CardHeader class="gap-3">
					<div class="flex items-center justify-between gap-3">
						<div>
							<CardTitle>Session prompts</CardTitle>
							<CardDescription>
								{mode === "authenticated"
									? "Shared copy for the authenticated workspace and upgrade hints."
									: "Shared copy for the guest lock and sign-in surfaces."}
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
							<p class="mt-2 text-sm leading-6 text-text-body">{prompt.body}</p>
						</div>
					{/each}
				</CardContent>
			</Card>

			<Card class="panel-card">
				<CardHeader class="gap-3">
					<CardTitle>{mode === "authenticated" ? "Recent authenticated runs" : "Recent guest runs"}</CardTitle>
					<CardDescription>
						{mode === "authenticated"
							? "Persistence is part of the signed-in contract, so the history rail can grow with the same shell."
							: "Demo history stays visible, while persistence remains locked."}
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
							? "Authenticated limits are surfaced directly in the shell."
							: "The guest limits from the docs are surfaced directly in the shell."}
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
					<CardTitle>Authenticated contract</CardTitle>
					<CardDescription>The same shell now allows custom targets, saved history, and account-level controls.</CardDescription>
				</CardHeader>
				<CardContent class="space-y-2 text-sm leading-6 text-text-body">
					<p>Custom URLs are editable, saved requests persist, and the runner keeps validating outbound requests.</p>
					<p>Guest lock surfaces are intentionally hidden because the capability is already unlocked.</p>
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
					<CardTitle>Upgrade placeholder</CardTitle>
					<CardDescription>{workspaceState.lockedActions[0]?.description ?? "Higher-tier entitlements can layer on later."}</CardDescription>
				</CardHeader>
				<CardContent class="text-sm leading-6 text-text-body">
					<p>The authenticated contract leaves room for plan-based limits without changing the shared `/app` surface.</p>
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
							? "Shared metrics from the authenticated state model keep the shell grounded in the product rules."
							: "Shared metrics from the guest state model keep the shell grounded in the product rules."}
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
