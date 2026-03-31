<script lang="ts">
	import GuestLockSurface from "./guest-lock-surface.svelte";

const bullets = [
	"Advanced utilities remain visible.",
	"Power-user tools are gated.",
	"Signed-in access unlocks the drawer.",
];

	const highlights = [
		{ label: "Snippets", value: "Preview only" },
		{ label: "Advanced tools", value: "Locked" },
		{ label: "Unlock path", value: "Authenticated" },
	];
</script>

<GuestLockSurface
	eyebrow="Advanced tools"
	title="Advanced tools are gated"
	description="Power-user tools unlock after sign-in."
	badge="Utility drawer"
	summary="Advanced tooling is available in signed-in mode."
	{bullets}
	{highlights}
	primaryActionLabel="Sign in for tools"
	primaryActionHref="/app"
	secondaryActionLabel="Docs"
	secondaryActionHref="/docs"
	tone="emerald"
/>
