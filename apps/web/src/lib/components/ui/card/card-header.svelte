<script lang="ts">
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div
	bind:this={ref}
	data-slot="card-header"
	class={cn(
		"group/card-header @container/card-header grid auto-rows-min items-start gap-1 rounded-t-[inherit] px-6 pt-6 pb-0 has-data-[slot=card-action]:grid-cols-[1fr_auto] has-data-[slot=card-description]:grid-rows-[auto_auto] group-data-[size=sm]/card:px-4 group-data-[size=sm]/card:pt-4 [.border-b]:pb-5 group-data-[size=sm]/card:[.border-b]:pb-4",
		className
	)}
	{...restProps}
>
	{@render children?.()}
</div>
