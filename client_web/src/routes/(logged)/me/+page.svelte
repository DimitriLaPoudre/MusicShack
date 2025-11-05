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

	async function Logout() {
		await fetch("http://localhost:8080/api/logout", {
			method: "POST",
			credentials: "include",
		});
		goto("/login");
	}
</script>

<h1>User</h1>
<div>
	<button on:click={Logout} style="color: red">Logout</button>
</div>
