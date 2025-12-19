<script lang="ts">
	let input = $state<string>("");

	async function searchFunction() {
		const encodedSearchData = encodeURI(input);
		window.location.assign(`/search?q=${encodedSearchData}`);
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
