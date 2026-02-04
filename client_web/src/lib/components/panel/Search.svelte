<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { searchInput } from "$lib/stores/panel/search";

	let inputFocus: HTMLInputElement;

	onMount(() => {
		inputFocus.focus();
	});

	async function searchFunction(e: SubmitEvent) {
		e.preventDefault();
		const encodedSearchData = encodeURIComponent($searchInput);
		$searchInput = "";
		await goto(`/search?q=${encodedSearchData}`);
	}
</script>

<div>
	<form onsubmit={searchFunction}>
		<input
			class="w-full m-0 border-none"
			type="text"
			bind:value={$searchInput}
			bind:this={inputFocus}
			placeholder="Search"
		/>
	</form>
</div>
