<script lang="ts">
	import { cn } from "$lib/utils.js";
	import Loader2Icon from '@lucide/svelte/icons/loader-2';
	import type { ComponentProps } from "svelte";
	import type { SVGAttributes } from "svelte/elements";

	let {
		class: className,
		role = "status",
		// we add color and stroke for compatibility with different icon libraries props
		color,
		stroke,
		"aria-label": ariaLabel = "Loading",
		...restProps
	}: SVGAttributes<SVGSVGElement> = $props();

	const iconProps = $derived(
		{
			...Object.fromEntries(
				Object.entries(restProps).filter(([, value]) => value !== null && value !== undefined)
			),
			role,
			color: color === null ? undefined : color,
			stroke: stroke === null ? undefined : stroke,
			"aria-label": ariaLabel,
			class: cn("size-4 animate-spin", className),
		} as ComponentProps<typeof Loader2Icon>
	);
</script>

<Loader2Icon {...iconProps} />
