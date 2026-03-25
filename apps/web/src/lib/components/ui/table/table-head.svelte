<script lang="ts">
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLThAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLThAttributes> = $props();
</script>

<th
	bind:this={ref}
	data-slot="table-head"
	class={cn("h-11 whitespace-nowrap px-3 text-left align-middle font-semibold tracking-tight text-text-strong [&:has([role=checkbox])]:pr-0", className)}
	{...restProps}
>
	{@render children?.()}
</th>
