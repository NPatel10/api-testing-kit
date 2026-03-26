<script lang="ts">
	import HistoryPage from "$lib/components/history/history-page.svelte";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();
</script>

<HistoryPage {data} />

