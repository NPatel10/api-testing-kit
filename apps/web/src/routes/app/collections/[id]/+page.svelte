<script lang="ts">
	import CollectionDetailPage from "$lib/components/collections/collection-detail-page.svelte";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();
</script>

<CollectionDetailPage detail={data.detail} mode={data.mode} entitlements={data.entitlements} sessionLabel={data.sessionLabel} />
