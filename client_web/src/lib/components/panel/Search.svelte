<script lang="ts">
	import { goto } from "$app/navigation";

	let input = $state<string>("");

	async function searchFunction(e: SubmitEvent) {
		e.preventDefault();
		const encodedSearchData = encodeURIComponent(input);
		await goto(`/search?q=${encodedSearchData}`);
	}
</script>

<div class="body">
	<form onsubmit={searchFunction}>
		<input type="text" bind:value={input} placeholder="Search" />
	</form>
</div>

<style>
	.body {
		input {
			width: 100%;
			margin: 0;
			border: none;
		}
	}
</style>
