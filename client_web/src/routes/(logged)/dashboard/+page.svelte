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

<h1 style="text-align: center;">Dashboard</h1>
{#if error}
	<p class="error-msg">{error}</p>
{/if}
