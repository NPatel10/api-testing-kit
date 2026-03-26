<script lang="ts">
	import AppShellShell from "$lib/components/app-shell/app-shell-shell.svelte";
	import type { Snippet } from "svelte";
	import type { LayoutData } from "./$types";

	let { children, data }: { children: Snippet; data: LayoutData } = $props();
</script>

<AppShellShell mode={data.mode} sessionLabel={data.sessionLabel} entitlements={data.entitlements}>
	{@render children()}
</AppShellShell>
