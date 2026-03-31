<script lang="ts">
	import GuestLockSurface from "./guest-lock-surface.svelte";

const bullets = [
	"History remains visible.",
	"Runs are ephemeral in guest mode.",
	"Replay and persistence stay locked.",
];

	const highlights = [
		{ label: "Guest retention", value: "Ephemeral" },
		{ label: "Signed-in retention", value: "Stored runs" },
		{ label: "Visibility", value: "Always shown" },
	];
</script>

<GuestLockSurface
	eyebrow="Request history"
	title="History persistence is locked"
	description="Guest mode does not store durable runs."
	badge="History rail"
	summary="Stored history unlocks after sign-in."
	{bullets}
	{highlights}
	primaryActionLabel="Sign in to view history"
	primaryActionHref="/app"
	secondaryActionLabel="Docs"
	secondaryActionHref="/docs"
	tone="slate"
/>
