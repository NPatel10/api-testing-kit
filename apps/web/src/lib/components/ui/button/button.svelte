<script lang="ts" module>
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from "svelte/elements";
	import { type VariantProps, tv } from "tailwind-variants";

	export const buttonVariants = tv({
		base: "group/button inline-flex shrink-0 items-center justify-center whitespace-nowrap border border-transparent bg-clip-padding text-sm font-medium transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-3 active:translate-y-px aria-invalid:border-danger aria-invalid:ring-danger/20 aria-invalid:ring-3 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
		variants: {
			variant: {
				default:
					"rounded-full bg-primary text-primary-foreground shadow-[0_12px_28px_rgba(31,122,77,0.18)] hover:bg-primary/90 hover:shadow-[0_14px_32px_rgba(31,122,77,0.22)]",
				outline:
					"rounded-full border-border bg-surface text-text-strong shadow-none hover:border-primary-green/35 hover:bg-primary-green-soft/40 hover:text-text-strong data-[state=open]:bg-primary-green-soft/50",
				secondary:
					"rounded-full border border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80 data-[state=open]:bg-secondary/85",
				ghost:
					"rounded-full border border-transparent text-text-body hover:bg-surface-soft hover:text-text-strong data-[state=open]:bg-surface-soft",
				destructive:
					"rounded-full border border-danger/15 bg-danger/10 text-danger hover:bg-danger/15 focus-visible:border-danger focus-visible:ring-danger/20",
				link: "rounded-none px-0 text-primary underline-offset-4 hover:underline",
			},
			size: {
				default:
					"h-10 gap-1.5 rounded-full px-4 in-data-[slot=button-group]:rounded-full has-data-[icon=inline-end]:pr-3 has-data-[icon=inline-start]:pl-3",
				xs:
					"h-7 gap-1 rounded-full px-2.5 text-xs in-data-[slot=button-group]:rounded-full has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 [&_svg:not([class*='size-'])]:size-3",
				sm:
					"h-9 gap-1.5 rounded-full px-3 in-data-[slot=button-group]:rounded-full has-data-[icon=inline-end]:pr-2.5 has-data-[icon=inline-start]:pl-2.5",
				lg: "h-11 gap-2 rounded-full px-5 has-data-[icon=inline-end]:pr-4 has-data-[icon=inline-start]:pl-4",
				icon: "size-10 rounded-full",
				"icon-xs": "size-7 rounded-full [&_svg:not([class*='size-'])]:size-3",
				"icon-sm": "size-8 rounded-full",
				"icon-lg": "size-11 rounded-full",
			},
		},
		defaultVariants: {
			variant: "default",
			size: "default",
		},
	});

	export type ButtonVariant = VariantProps<typeof buttonVariants>["variant"];
	export type ButtonSize = VariantProps<typeof buttonVariants>["size"];

	export type ButtonProps = WithElementRef<HTMLButtonAttributes> &
		WithElementRef<HTMLAnchorAttributes> & {
			variant?: ButtonVariant;
			size?: ButtonSize;
		};
</script>

<script lang="ts">
	let {
		class: className,
		variant = "default",
		size = "default",
		ref = $bindable(null),
		href = undefined,
		type = "button",
		disabled,
		children,
		...restProps
	}: ButtonProps = $props();
</script>

{#if href}
	<a
		bind:this={ref}
		data-slot="button"
		class={cn(buttonVariants({ variant, size }), className)}
		href={disabled ? undefined : href}
		aria-disabled={disabled}
		role={disabled ? "link" : undefined}
		tabindex={disabled ? -1 : undefined}
		{...restProps}
	>
		{@render children?.()}
	</a>
{:else}
	<button
		bind:this={ref}
		data-slot="button"
		class={cn(buttonVariants({ variant, size }), className)}
		{type}
		{disabled}
		{...restProps}
	>
		{@render children?.()}
	</button>
{/if}
