<script lang="ts">
	import GuestLockSurface from "./guest-lock-surface.svelte";

const bullets = [
	"Secrets stay locked.",
	"Scoped variables are signed-in features.",
	"Request state remains visible.",
];

	const highlights = [
		{ label: "Secret storage", value: "Locked" },
		{ label: "Variable editor", value: "Preview only" },
		{ label: "Scope", value: "Per account" },
	];
</script>

<GuestLockSurface
	eyebrow="Environment variables"
	title="Env vars are locked"
	description="Scoped secrets and reusable values unlock after sign-in."
	badge="Variables"
	summary="Available in signed-in sessions."
	{bullets}
	{highlights}
	primaryActionLabel="Sign in for variables"
	primaryActionHref="/app"
	secondaryActionLabel="Docs"
	secondaryActionHref="/docs"
	tone="sand"
/>
