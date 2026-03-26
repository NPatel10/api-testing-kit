<script lang="ts">
	import AppShellUtilityDrawer from "$lib/components/app-shell/app-shell-utility-drawer.svelte";
	import AppShellWorkspace from "$lib/components/app-shell/app-shell-workspace.svelte";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();
</script>

<section class="space-y-4">
	<AppShellWorkspace mode={data.mode} entitlements={data.entitlements} />
	<AppShellUtilityDrawer mode={data.mode} entitlements={data.entitlements} />
</section>
