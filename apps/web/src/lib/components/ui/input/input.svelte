<script lang="ts">
	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from "svelte/elements";
	import { cn, type WithElementRef } from "$lib/utils.js";

	type InputType = Exclude<HTMLInputTypeAttribute, "file">;

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, "type"> &
			({ type: "file"; files?: FileList } | { type?: InputType; files?: undefined })
	>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		type,
		files = $bindable(),
		class: className,
		"data-slot": dataSlot = "input",
		...restProps
	}: Props = $props();
</script>

{#if type === "file"}
	<input
		bind:this={ref}
		data-slot={dataSlot}
		class={cn(
			"h-10 w-full min-w-0 rounded-[16px] border border-input bg-surface px-3 py-2 text-sm text-text-strong shadow-[0_4px_14px_rgba(21,31,23,0.03)] transition-[color,box-shadow,border-color] placeholder:text-text-muted file:inline-flex file:h-8 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-text-strong focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-3 aria-invalid:border-danger aria-invalid:ring-danger/20 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
			className
		)}
		type="file"
		bind:files
		bind:value
		{...restProps}
	/>
{:else}
	<input
		bind:this={ref}
		data-slot={dataSlot}
		class={cn(
			"h-10 w-full min-w-0 rounded-[16px] border border-input bg-surface px-3 py-2 text-sm text-text-strong shadow-[0_4px_14px_rgba(21,31,23,0.03)] transition-[color,box-shadow,border-color] placeholder:text-text-muted file:inline-flex file:h-8 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-text-strong focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-3 aria-invalid:border-danger aria-invalid:ring-danger/20 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
			className
		)}
		{type}
		bind:value
		{...restProps}
	/>
{/if}
