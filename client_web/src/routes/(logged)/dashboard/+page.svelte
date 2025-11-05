<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";

	afterNavigate(async () => {
		const res = await fetch("http://localhost:8080/api/me", {
			credentials: "include",
		});

		if (!res.ok) {
			goto("/login");
			return;
		}
	});
</script>

<h1>Dashboard</h1>
