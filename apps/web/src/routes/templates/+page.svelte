<script lang="ts">
	import TemplatesPage from "$lib/components/templates-page/templates-page.svelte";
	import type { PageData } from "./$types";

	let { data } = $props<{ data: PageData }>();
</script>

<TemplatesPage {data} />
