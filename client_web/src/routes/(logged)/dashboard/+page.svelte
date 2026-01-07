<script lang="ts">
	import { apiFetch } from "$lib/functions/fetch";
	import type { StatusResponse } from "$lib/types/response";
	import { onMount } from "svelte";

	let error = $state<null | string>(null);

	onMount(async () => {
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

<svelte:head>
	<title>Dashboard - MusicShack</title>
</svelte:head>

{#if error}
	<div class="error">
		<h2>Error loading Dashboard</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else}
	<div class="body">
		<h1>Dashboard</h1>
	</div>
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

	.body {
		margin-top: 30px;

		h1 {
			text-align: center;
			font-weight: bolder;
		}
	}
</style>
