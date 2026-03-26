<script lang="ts">
	import AdminAbuseDashboard from "$lib/components/admin-abuse/admin-abuse-dashboard.svelte";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Admin Abuse Monitor - API Testing Kit</title>
	<meta
		name="description"
		content="Internal abuse monitoring and blocked target review for API Testing Kit."
	/>
</svelte:head>

<AdminAbuseDashboard dashboard={data.dashboard} />
