<script lang="ts" module>
	import { type VariantProps, tv } from "tailwind-variants";

	export const badgeVariants = tv({
		base: "group/badge inline-flex w-fit shrink-0 items-center justify-center overflow-hidden whitespace-nowrap rounded-full border border-transparent px-2.5 py-1 text-xs font-medium transition-all focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:border-danger aria-invalid:ring-danger/20 has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&>svg]:pointer-events-none [&>svg]:size-3!",
		variants: {
			variant: {
				default: "bg-primary text-primary-foreground shadow-[0_8px_18px_rgba(31,122,77,0.16)] [a]:hover:bg-primary/90",
				secondary: "border-border bg-secondary text-secondary-foreground [a]:hover:bg-secondary/80",
				destructive:
					"border border-danger/15 bg-danger/10 text-danger [a]:hover:bg-danger/15 focus-visible:ring-danger/20",
				outline: "border-border bg-surface text-text-strong [a]:hover:bg-surface-soft [a]:hover:text-text-strong",
				ghost: "bg-transparent text-text-body hover:bg-surface-soft hover:text-text-strong",
				link: "text-primary underline-offset-4 hover:underline",
			},
		},
		defaultVariants: {
			variant: "default",
		},
	});

	export type BadgeVariant = VariantProps<typeof badgeVariants>["variant"];
</script>

<script lang="ts">
	import type { HTMLAnchorAttributes } from "svelte/elements";
	import { cn, type WithElementRef } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		href,
		class: className,
		variant = "default",
		children,
		...restProps
	}: WithElementRef<HTMLAnchorAttributes> & {
		variant?: BadgeVariant;
	} = $props();
</script>

<svelte:element
	this={href ? "a" : "span"}
	bind:this={ref}
	data-slot="badge"
	{href}
	class={cn(badgeVariants({ variant }), className)}
	{...restProps}
>
	{@render children?.()}
</svelte:element>
