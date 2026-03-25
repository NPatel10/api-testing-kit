<script lang="ts">
	import { Tabs as TabsPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: TabsPrimitive.TriggerProps = $props();
</script>

	<TabsPrimitive.Trigger
		bind:ref
		data-slot="tabs-trigger"
		class={cn(
		"relative inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center whitespace-nowrap rounded-full border border-transparent px-3 py-2 text-sm font-medium text-text-muted transition-all focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 group-data-[orientation=vertical]/tabs:w-full group-data-[orientation=vertical]/tabs:justify-start",
		"group-data-[variant=default]/tabs-list:data-active:bg-surface group-data-[variant=default]/tabs-list:data-active:text-text-strong group-data-[variant=default]/tabs-list:data-active:shadow-[0_6px_16px_rgba(21,31,23,0.08)]",
		"group-data-[variant=line]/tabs-list:rounded-full group-data-[variant=line]/tabs-list:px-3 group-data-[variant=line]/tabs-list:data-active:bg-primary-green-soft group-data-[variant=line]/tabs-list:data-active:text-primary-green-deep",
		"after:absolute after:opacity-0 after:transition-opacity group-data-[orientation=horizontal]/tabs:after:inset-x-0 group-data-[orientation=horizontal]/tabs:after:bottom-[-6px] group-data-[orientation=horizontal]/tabs:after:h-0.5 group-data-[orientation=vertical]/tabs:after:inset-y-0 group-data-[orientation=vertical]/tabs:after:-right-1 group-data-[orientation=vertical]/tabs:after:w-0.5 group-data-[variant=line]/tabs-list:data-active:after:bg-primary-green group-data-[variant=line]/tabs-list:data-active:after:opacity-100",
		"data-active:text-text-strong",
		className
	)}
	{...restProps}
/>
