<script lang="ts">
	import { goto } from "$app/navigation";
	import { tick } from "svelte";
	import "../../app.css";
	import { Search, User } from "lucide-svelte";

	let { children } = $props();

	let searchData = "";

	async function search() {
		const encodedSearchData = encodeURI(searchData);
		window.location.assign(`/search?q=${encodedSearchData}`);
	}
</script>

<header
	style="display: flex; justify-content: space-between; align-items: center; position: relative; padding: 20px 20px; gap: 10px; flex-wrap: wrap; background-color: var(--color-background-light);"
>
	<a href="/dashboard">MusicShack</a>
	<form
		on:submit={search}
		style="position: absolute; left: 50%; transform: translateX(-50%); width: 300px; justify-content: space-between; padding: 8px; background-color: var(--color-secondary); display: flex; flex-direction: row; gap: 10px"
	>
		<input type="search" placeholder="Search" bind:value={searchData} />
		<button><Search size="28" color="var(--color-primary)" /></button>
	</form>
	<a href="/me"> <User /> </a>
</header>

{@render children?.()}
