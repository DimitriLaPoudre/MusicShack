<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";

	let error = $state<string | null>(null);

	afterNavigate(async () => {
		try {
			const res = await fetch(`http://localhost:8080/api/me`, {
				credentials: "include",
			});
			if (res.status === 401) {
				goto("/login");
				return;
			}
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load dashboard";
		}
	});
</script>

{#if error}
	<div class="error">
		<h2>Error Loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else}
	<h1 style="text-align: center;">Dashboard</h1>
{/if}

<style>
	.error {
		margin-top: 30px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 10px;

		* {
			margin: 0;
		}
	}
</style>
