<script lang="ts">
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLTableRowElement>> = $props();
</script>

<tr
	bind:this={ref}
	data-slot="table-row"
	class={cn("border-b border-border/80 transition-colors hover:bg-surface-soft data-[state=selected]:bg-primary-green-soft/40", className)}
	{...restProps}
>
	{@render children?.()}
</tr>
