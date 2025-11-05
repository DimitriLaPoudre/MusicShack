<script lang="ts">
	import "../../app.css";
	import { Search, User } from "lucide-svelte";

	let { children } = $props();
</script>

<header
	style="display: flex; justify-content: space-between; align-items: center; padding: 10px 20px; gap: 10px; flex-wrap: wrap; background-color: var(--color-background-light);"
>
	<a href="/">MusicShack</a>
	<div
		style="width: 300px; max-width: 100%; padding: 8px; background-color: var(--color-secondary); display: flex; flex-direction: row; gap: 10px"
	>
		<Search size="28" color="var(--color-primary)" />
		<input type="text" placeholder="Search" />
	</div>
	<a href="/me"> <User /> </a>
</header>

{@render children?.()}
