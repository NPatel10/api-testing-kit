<script lang="ts">
	import type { HTMLAttributes } from "svelte/elements";
	import { cn, type WithElementRef } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		children,
		size = "default",
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & { size?: "default" | "sm" } = $props();
</script>

<div
	bind:this={ref}
	data-slot="card"
	data-size={size}
	class={cn(
		"group/card flex flex-col gap-6 overflow-hidden rounded-[24px] border border-border/80 bg-card text-card-foreground text-sm shadow-[0_10px_30px_rgba(21,31,23,0.05)] has-[>img:first-child]:pt-0 data-[size=sm]:gap-4 data-[size=sm]:py-4 *:[img:first-child]:rounded-t-[24px] *:[img:last-child]:rounded-b-[24px]",
		className
	)}
	{...restProps}
>
	{@render children?.()}
</div>
