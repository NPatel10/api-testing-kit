<script lang="ts">
	import Badge from "$lib/components/ui/badge/badge.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import { Checkbox } from "$lib/components/ui/checkbox/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Tabs, TabsContent, TabsList, TabsTrigger } from "$lib/components/ui/tabs/index.js";
	import { Textarea } from "$lib/components/ui/textarea/index.js";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CirclePlusIcon from "@lucide/svelte/icons/circle-plus";
	import FileJsonIcon from "@lucide/svelte/icons/file-json";
	import LockIcon from "@lucide/svelte/icons/lock";
	import PlayIcon from "@lucide/svelte/icons/play";
	import ShieldIcon from "@lucide/svelte/icons/shield";

	import {
		cloneRequestDraft,
		createDefaultRequestDraft,
		createRequestRow,
		formatMethodTone,
		formatValidationTone,
		getRequestValidation,
		requestAuthSchemes,
		requestBodyModes,
		requestMethods,
		type RequestBuilderDraft,
		type RequestBuilderMode,
		type RequestAuthScheme,
		type RequestBodyMode,
		type RequestRow,
	} from "./request-builder";

	type Props = {
		title?: string;
		description?: string;
		mode?: RequestBuilderMode;
		request?: RequestBuilderDraft;
		sendLabel?: string;
		pending?: boolean;
		lockedNote?: string;
		onSend?: (draft: RequestBuilderDraft) => void;
		onOpenSignIn?: () => void;
		onUpgrade?: () => void;
	};

	let {
		title = "Request builder",
		description = "HTTP request editor for the shared workspace.",
		mode = "guest",
		request,
		sendLabel = "Send request",
		pending = false,
		lockedNote = "Custom targets remain locked in guest mode.",
		onSend,
		onOpenSignIn,
		onUpgrade,
	}: Props = $props();

	let editor = $state(createDefaultRequestDraft("guest"));
	let activeTab = $state("params");

	$effect(() => {
		editor = cloneRequestDraft(request ?? createDefaultRequestDraft(mode));
	});

	function setMethod(method: RequestBuilderDraft["method"]) {
		editor.method = method;
	}

	function addQueryParam() {
		editor.queryParams = [...editor.queryParams, createRequestRow()];
	}

	function addHeader() {
		editor.headers = [...editor.headers, createRequestRow()];
	}

	function addFormRow() {
		editor.body.formRows = [...editor.body.formRows, createRequestRow()];
	}

	function updateQueryParam(index: number, patch: Partial<RequestRow>) {
		editor.queryParams = editor.queryParams.map((row, currentIndex) =>
			currentIndex === index ? { ...row, ...patch } : row
		);
	}

	function updateHeader(index: number, patch: Partial<RequestRow>) {
		editor.headers = editor.headers.map((row, currentIndex) =>
			currentIndex === index ? { ...row, ...patch } : row
		);
	}

	function updateFormRow(index: number, patch: Partial<RequestRow>) {
		editor.body.formRows = editor.body.formRows.map((row, currentIndex) =>
			currentIndex === index ? { ...row, ...patch } : row
		);
	}

	function removeQueryParam(index: number) {
		editor.queryParams = editor.queryParams.filter((_, currentIndex) => currentIndex !== index);
	}

	function removeHeader(index: number) {
		editor.headers = editor.headers.filter((_, currentIndex) => currentIndex !== index);
	}

	function removeFormRow(index: number) {
		editor.body.formRows = editor.body.formRows.filter((_, currentIndex) => currentIndex !== index);
	}

	function setAuthScheme(scheme: RequestAuthScheme) {
		editor.auth = { ...editor.auth, scheme };
	}

	function setBodyMode(modeValue: RequestBodyMode) {
		editor.body = {
			...editor.body,
			mode: modeValue,
			contentType:
				modeValue === "json"
					? "application/json"
					: modeValue === "form"
						? "application/x-www-form-urlencoded"
						: "text/plain",
		};
	}

	function handleSend() {
		const validation = getRequestValidation(mode, editor);
		if (validation) {
			return;
		}

		onSend?.(cloneRequestDraft(editor));
	}

	function requestContentType() {
		return editor.body.mode === "form" ? "application/x-www-form-urlencoded" : editor.body.contentType;
	}
</script>

<Card class="panel-card overflow-hidden bg-gradient-to-br from-white via-white to-panel-soft/70">
	<CardHeader class="gap-4 border-b border-border/70 bg-[linear-gradient(180deg,rgba(255,255,255,0.92),rgba(245,241,235,0.84))]">
		<div class="flex flex-wrap items-start justify-between gap-3">
			<div class="space-y-1">
				<p class="section-title">Workspace module</p>
				<CardTitle class="text-xl tracking-tight">{title}</CardTitle>
				<CardDescription class="max-w-2xl text-sm leading-6 text-text-body">
					{description}
				</CardDescription>
			</div>

			<div class="flex flex-wrap items-center gap-2">
				<Badge variant={mode === "guest" ? "secondary" : "default"}>
					{mode === "guest" ? "Guest-safe" : "Authenticated"}
				</Badge>
				<Badge
					variant="outline"
					class={formatValidationTone(getRequestValidation(mode, editor)?.severity ?? "info")}
				>
					{getRequestValidation(mode, editor)?.title ?? "Ready to send"}
				</Badge>
			</div>
		</div>

		<div class="flex flex-wrap items-center gap-2">
			{#each requestMethods as method}
				<Button
					type="button"
					variant={editor.method === method ? "default" : "outline"}
					size="sm"
					class="rounded-full px-4 text-xs font-semibold tracking-[0.08em]"
					onclick={() => setMethod(method)}
				>
					{method}
				</Button>
			{/each}
		</div>
	</CardHeader>

	<CardContent class="space-y-4 p-5">
		{@const validation = getRequestValidation(mode, editor)}
		<div class="grid gap-3 xl:grid-cols-[minmax(0,1fr)_auto]">
			<div class="relative">
				<Input
					bind:value={editor.url}
					placeholder="https://jsonplaceholder.typicode.com/posts/1"
					class="h-12 rounded-full border-border/80 bg-white px-4 pr-28 font-mono text-sm shadow-xs"
					aria-invalid={validation?.severity === "danger"}
				/>
				<div class="pointer-events-none absolute inset-y-0 right-3 flex items-center gap-2">
					<Badge variant="outline" class="rounded-full border-border/70 bg-panel-soft text-[11px] uppercase tracking-[0.16em] text-text-muted">
						URL
					</Badge>
					<Badge variant="outline" class="rounded-full border-border/70 bg-panel-soft text-[11px] uppercase tracking-[0.16em] text-text-muted">
						Preview
					</Badge>
				</div>
			</div>

			<Button
				type="button"
				variant="default"
				size="lg"
				class="pill-button shadow-sm"
				disabled={Boolean(validation) || pending}
				onclick={handleSend}
			>
				<PlayIcon class="size-4" />
				{pending ? "Sending..." : sendLabel}
			</Button>
		</div>

		<div class="grid gap-3 lg:grid-cols-[1.1fr_0.9fr]">
			<div
				class={`rounded-[20px] border px-4 py-3 shadow-sm ${
					validation ? formatValidationTone(validation.severity) : "border-emerald-200 bg-emerald-50/80 text-emerald-950"
				}`}
			>
				<div class="flex items-start gap-3">
					<div class="mt-0.5 rounded-full border border-current/20 bg-white/70 p-2">
						{#if validation?.severity === "danger"}
							<LockIcon class="size-4" />
						{:else if validation?.severity === "warning"}
							<ShieldIcon class="size-4" />
						{:else}
							<FileJsonIcon class="size-4" />
						{/if}
					</div>
					<div class="min-w-0 flex-1">
						<p class="text-sm font-semibold tracking-tight">{validation?.title ?? "Request looks ready"}</p>
						<p class="mt-1 text-xs leading-5">
							{validation?.description ?? "Request is valid."}
						</p>
						{#if validation?.action}
							<p class="mt-2 text-xs font-medium">{validation.action}</p>
						{/if}
					</div>
				</div>
			</div>

			<div class="rounded-[20px] border border-border/70 bg-white/80 px-4 py-3 shadow-sm">
				<div class="flex items-center justify-between gap-3">
					<div>
						<p class="text-xs uppercase tracking-[0.22em] text-text-muted">Request summary</p>
						<p class="mt-1 text-sm font-semibold text-foreground">{editor.method} {editor.url || "No URL yet"}</p>
					</div>
					<Badge variant="outline" class={formatMethodTone(editor.method)}>{editor.method}</Badge>
				</div>
				<div class="mt-3 flex flex-wrap gap-2 text-xs">
					<Badge variant="outline" class="border-border/70 bg-panel-soft text-text-body">
						{editor.body.mode.toUpperCase()}
					</Badge>
					<Badge variant="outline" class="border-border/70 bg-panel-soft text-text-body">
						{requestContentType()}
					</Badge>
					<Badge variant="outline" class="border-border/70 bg-panel-soft text-text-body">
						{mode === "guest" ? "Allowlisted only" : "Custom targets allowed"}
					</Badge>
				</div>
			</div>
		</div>

		{#if mode === "guest" && validation?.severity === "danger"}
			<div class="locked-overlay flex items-start gap-4 px-4 py-4">
				<div class="rounded-full border border-border/70 bg-panel-soft p-2 text-text-muted">
					<LockIcon class="size-4" />
				</div>
				<div class="min-w-0 flex-1">
					<div class="flex flex-wrap items-center gap-2">
						<p class="text-sm font-semibold tracking-tight text-foreground">Guest restrictions stay in force</p>
						<Badge variant="secondary">Locked</Badge>
					</div>
					<p class="mt-1 text-sm leading-6 text-text-body">{lockedNote}</p>
					<div class="mt-3 flex flex-wrap gap-2">
						<Button
							type="button"
							variant="default"
							size="sm"
							class="rounded-full"
							onclick={() => onOpenSignIn?.()}
							disabled={!onOpenSignIn}
						>
							Sign in
							<ArrowRightIcon class="size-4" />
						</Button>
						<Button
							type="button"
							variant="outline"
							size="sm"
							class="rounded-full"
							onclick={() => onUpgrade?.()}
							disabled={!onUpgrade}
						>
							Open templates
						</Button>
					</div>
				</div>
			</div>
		{/if}

		<Tabs bind:value={activeTab} class="gap-4">
			<TabsList class="flex flex-wrap rounded-full border border-border/70 bg-panel-soft p-1">
				<TabsTrigger value="params">Params</TabsTrigger>
				<TabsTrigger value="headers">Headers</TabsTrigger>
				<TabsTrigger value="auth">Auth</TabsTrigger>
				<TabsTrigger value="body">Body</TabsTrigger>
			</TabsList>

			<TabsContent value="params" class="space-y-3">
				<div class="flex items-center justify-between gap-3">
					<div>
						<p class="text-sm font-semibold text-foreground">Query params</p>
						<p class="text-xs leading-5 text-text-muted">Mode-aware query fields.</p>
					</div>
					<Button type="button" variant="outline" size="sm" class="rounded-full" onclick={addQueryParam}>
						<CirclePlusIcon class="size-4" />
						Add row
					</Button>
				</div>

				{#if editor.queryParams.length}
					<div class="space-y-2">
						{#each editor.queryParams as row, index}
							<div class="grid gap-2 rounded-[18px] border border-border/70 bg-panel-soft p-3 md:grid-cols-[1fr_1.2fr_auto]">
								<Input
									value={row.key}
									placeholder="Key"
									class="h-10 rounded-full bg-white font-mono text-sm"
									oninput={(event) => updateQueryParam(index, { key: event.currentTarget.value })}
								/>
								<Input
									value={row.value}
									placeholder="Value"
									class="h-10 rounded-full bg-white font-mono text-sm"
									oninput={(event) => updateQueryParam(index, { value: event.currentTarget.value })}
								/>
								<div class="flex items-center justify-between gap-2">
									<label class="flex items-center gap-2 text-xs text-text-body">
										<Checkbox bind:checked={row.enabled} />
										Enabled
									</label>
									<Button type="button" variant="ghost" size="icon-sm" class="rounded-full" onclick={() => removeQueryParam(index)}>
										<span class="sr-only">Remove param</span>
										<ArrowRightIcon class="size-4 rotate-45" />
									</Button>
								</div>
							</div>
						{/each}
					</div>
				{:else}
					<div class="soft-card flex items-start gap-3 px-4 py-4">
						<CirclePlusIcon class="mt-0.5 size-4 text-primary" />
						<div>
							<p class="text-sm font-semibold text-foreground">No query params yet</p>
							<p class="mt-1 text-xs leading-5 text-text-muted">No query entries.</p>
						</div>
					</div>
				{/if}
			</TabsContent>

			<TabsContent value="headers" class="space-y-3">
				<div class="flex items-center justify-between gap-3">
					<div>
						<p class="text-sm font-semibold text-foreground">Headers</p>
						<p class="text-xs leading-5 text-text-muted">Request headers.</p>
					</div>
					<Button type="button" variant="outline" size="sm" class="rounded-full" onclick={addHeader}>
						<CirclePlusIcon class="size-4" />
						Add row
					</Button>
				</div>

				{#if editor.headers.length}
					<div class="space-y-2">
						{#each editor.headers as row, index}
							<div class="grid gap-2 rounded-[18px] border border-border/70 bg-panel-soft p-3 md:grid-cols-[1fr_1.2fr_auto]">
								<Input
									value={row.key}
									placeholder="Header"
									class="h-10 rounded-full bg-white font-mono text-sm"
									oninput={(event) => updateHeader(index, { key: event.currentTarget.value })}
								/>
								<Input
									value={row.value}
									placeholder="Value"
									class="h-10 rounded-full bg-white font-mono text-sm"
									oninput={(event) => updateHeader(index, { value: event.currentTarget.value })}
								/>
								<div class="flex items-center justify-between gap-2">
									<label class="flex items-center gap-2 text-xs text-text-body">
										<Checkbox bind:checked={row.enabled} />
										Enabled
									</label>
									<Button type="button" variant="ghost" size="icon-sm" class="rounded-full" onclick={() => removeHeader(index)}>
										<span class="sr-only">Remove header</span>
										<ArrowRightIcon class="size-4 rotate-45" />
									</Button>
								</div>
							</div>
						{/each}
					</div>
				{:else}
					<div class="soft-card flex items-start gap-3 px-4 py-4">
						<FileJsonIcon class="mt-0.5 size-4 text-primary" />
						<div>
							<p class="text-sm font-semibold text-foreground">No headers yet</p>
							<p class="mt-1 text-xs leading-5 text-text-muted">No header entries.</p>
						</div>
					</div>
				{/if}
			</TabsContent>

			<TabsContent value="auth" class="space-y-3">
				<div class="flex items-center justify-between gap-3">
					<div>
						<p class="text-sm font-semibold text-foreground">Authorization</p>
						<p class="text-xs leading-5 text-text-muted">Mode-aware auth controls.</p>
					</div>
					<Badge variant={mode === "guest" ? "secondary" : "outline"}>{mode === "guest" ? "Locked in guest mode" : "Editable"}</Badge>
				</div>

				<div class="grid gap-2 md:grid-cols-3">
					{#each requestAuthSchemes as scheme}
						<Button
							type="button"
							variant={editor.auth.scheme === scheme.value ? "default" : "outline"}
							class="h-auto flex-col items-start justify-start rounded-[20px] px-4 py-3 text-left"
							onclick={() => setAuthScheme(scheme.value)}
							disabled={mode === "guest"}
						>
							<span class="text-sm font-semibold tracking-tight">{scheme.label}</span>
							<span class="mt-1 text-xs font-normal leading-5 text-inherit/70">{scheme.description}</span>
						</Button>
					{/each}
				</div>

				<div class="grid gap-3 rounded-[20px] border border-border/70 bg-panel-soft p-4 md:grid-cols-2">
					<div class="space-y-2">
						<p class="text-xs uppercase tracking-[0.2em] text-text-muted">Bearer token</p>
						<Input
							value={editor.auth.token}
							placeholder="Paste token"
							class="h-11 rounded-full bg-white font-mono text-sm"
							disabled={mode === "guest" || editor.auth.scheme !== "bearer"}
							oninput={(event) => (editor.auth = { ...editor.auth, token: event.currentTarget.value })}
						/>
					</div>
					<div class="grid gap-2 sm:grid-cols-2">
						<div class="space-y-2">
							<p class="text-xs uppercase tracking-[0.2em] text-text-muted">Username</p>
							<Input
								value={editor.auth.username}
								placeholder="Username"
								class="h-11 rounded-full bg-white font-mono text-sm"
								disabled={mode === "guest" || editor.auth.scheme !== "basic"}
								oninput={(event) => (editor.auth = { ...editor.auth, username: event.currentTarget.value })}
							/>
						</div>
						<div class="space-y-2">
							<p class="text-xs uppercase tracking-[0.2em] text-text-muted">Password</p>
							<Input
								value={editor.auth.password}
								type="password"
								placeholder="Password"
								class="h-11 rounded-full bg-white font-mono text-sm"
								disabled={mode === "guest" || editor.auth.scheme !== "basic"}
								oninput={(event) => (editor.auth = { ...editor.auth, password: event.currentTarget.value })}
							/>
						</div>
					</div>
				</div>

				{#if mode === "guest"}
					<div class="locked-overlay flex items-start gap-3 px-4 py-4">
						<ShieldIcon class="mt-0.5 size-4 text-primary" />
						<div>
							<p class="text-sm font-semibold text-foreground">Auth controls are locked</p>
							<p class="mt-1 text-xs leading-5 text-text-muted">
								Available after sign-in.
							</p>
						</div>
					</div>
				{/if}
			</TabsContent>

			<TabsContent value="body" class="space-y-3">
				<div class="flex flex-wrap items-center justify-between gap-3">
					<div>
						<p class="text-sm font-semibold text-foreground">Body</p>
						<p class="text-xs leading-5 text-text-muted">JSON, raw text, or form data.</p>
					</div>
					<div class="flex flex-wrap gap-2">
						{#each requestBodyModes as bodyMode}
							<Button
								type="button"
								variant={editor.body.mode === bodyMode.value ? "default" : "outline"}
								size="sm"
								class="rounded-full px-4"
								onclick={() => setBodyMode(bodyMode.value)}
							>
								{bodyMode.label}
							</Button>
						{/each}
					</div>
				</div>

				<div class="rounded-[20px] border border-border/70 bg-panel-soft p-4">
					<div class="mb-3 flex items-center justify-between gap-3">
						<div>
							<p class="text-xs uppercase tracking-[0.2em] text-text-muted">Current mode</p>
							<p class="mt-1 text-sm font-semibold text-foreground">
								{requestBodyModes.find((modeOption) => modeOption.value === editor.body.mode)?.label ?? "Body"}
							</p>
						</div>
						<Badge variant="outline" class="border-border/70 bg-white text-text-body">
							{requestContentType()}
						</Badge>
					</div>

					{#if editor.body.mode === "json"}
						<Textarea
							bind:value={editor.body.value}
							class="min-h-56 rounded-[20px] border-border/80 bg-white font-mono text-sm leading-6 shadow-xs"
							placeholder={`{\n  "name": "API Testing Kit"\n}`}
						/>
					{:else if editor.body.mode === "raw"}
						<Textarea
							bind:value={editor.body.value}
							class="min-h-56 rounded-[20px] border-border/80 bg-white font-mono text-sm leading-6 shadow-xs"
							placeholder="Plain text, XML, or another raw payload"
						/>
					{:else}
						<div class="space-y-3">
							<div class="grid gap-2 md:grid-cols-[1fr_1fr_auto]">
								<div>
									<p class="text-xs uppercase tracking-[0.2em] text-text-muted">Key</p>
								</div>
								<div>
									<p class="text-xs uppercase tracking-[0.2em] text-text-muted">Value</p>
								</div>
								<div></div>
							</div>

							{#if editor.body.formRows.length}
								{#each editor.body.formRows as row, index}
									<div class="grid gap-2 rounded-[18px] border border-border/70 bg-white p-3 md:grid-cols-[1fr_1fr_auto]">
										<Input
											value={row.key}
											class="h-10 rounded-full bg-white font-mono text-sm"
											placeholder="name"
											oninput={(event) => updateFormRow(index, { key: event.currentTarget.value })}
										/>
										<Input
											value={row.value}
											class="h-10 rounded-full bg-white font-mono text-sm"
											placeholder="value"
											oninput={(event) => updateFormRow(index, { value: event.currentTarget.value })}
										/>
										<div class="flex items-center justify-between gap-2">
											<label class="flex items-center gap-2 text-xs text-text-body">
												<Checkbox bind:checked={row.enabled} />
												Enabled
											</label>
											<Button type="button" variant="ghost" size="icon-sm" class="rounded-full" onclick={() => removeFormRow(index)}>
												<span class="sr-only">Remove form field</span>
												<ArrowRightIcon class="size-4 rotate-45" />
											</Button>
										</div>
									</div>
								{/each}
							{:else}
								<div class="soft-card flex items-start gap-3 px-4 py-4">
									<FileJsonIcon class="mt-0.5 size-4 text-primary" />
									<div>
										<p class="text-sm font-semibold text-foreground">No form fields yet</p>
										<p class="mt-1 text-xs leading-5 text-text-muted">No form entries.</p>
									</div>
								</div>
							{/if}

							<Button type="button" variant="outline" size="sm" class="rounded-full" onclick={addFormRow}>
								<CirclePlusIcon class="size-4" />
								Add form field
							</Button>
						</div>
					{/if}
				</div>
			</TabsContent>
		</Tabs>
	</CardContent>
</Card>
