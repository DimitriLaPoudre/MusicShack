<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import { apiFetch } from "$lib/functions/fetch";
	import type { StatusResponse } from "$lib/types/response";

	let error = $state<null | string>(null);

	afterNavigate(async () => {
		try {
			const data = await apiFetch<StatusResponse>("/me");
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch me");
			}
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load dashboard";
		}
	});
</script>

{#if error}
	<div class="error">
		<h2>Error loading Dashboard</h2>
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
	}
</style>
