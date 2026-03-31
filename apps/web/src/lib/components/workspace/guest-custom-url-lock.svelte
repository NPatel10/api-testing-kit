<script lang="ts">
	import GuestLockSurface from "./guest-lock-surface.svelte";

const bullets = [
	"Guest mode stays allowlisted.",
	"Custom domains are disabled.",
	"Outbound validation remains server enforced.",
];

	const highlights = [
		{ label: "Guest URLs", value: "Allowlisted" },
		{ label: "Custom domains", value: "Locked" },
		{ label: "Safety model", value: "Server enforced" },
	];
</script>

<GuestLockSurface
	eyebrow="Custom URLs"
	title="Custom outbound targets are locked"
	description="Guest sessions stay on allowlisted targets."
	badge="Outbound"
	summary="Available after sign-in and validation."
	{bullets}
	{highlights}
	primaryActionLabel="Unlock custom targets"
	primaryActionHref="/app"
	secondaryActionLabel="Docs"
	secondaryActionHref="/docs"
	tone="rose"
/>
