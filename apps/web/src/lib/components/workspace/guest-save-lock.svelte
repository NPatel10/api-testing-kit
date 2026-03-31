<script lang="ts">
	import GuestLockSurface from "./guest-lock-surface.svelte";

const bullets = [
	"Save actions remain visible.",
	"Permanent storage unlocks after sign-in.",
	"Guest mode does not write data.",
];

	const highlights = [
		{ label: "Guest save", value: "Preview only" },
		{ label: "Authenticated mode", value: "Persistent" },
		{ label: "Surface type", value: "Card + overlay" },
	];
</script>

<GuestLockSurface
	eyebrow="Save requests"
	title="Saving requests is locked"
	description="Guest mode does not persist requests."
	badge="Save control"
	summary="Persistence and collection management unlock after sign-in."
	{bullets}
	{highlights}
	primaryActionLabel="Sign in to save"
	primaryActionHref="/app"
	secondaryActionLabel="Docs"
	secondaryActionHref="/docs"
	tone="emerald"
/>
