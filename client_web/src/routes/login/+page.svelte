<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";

	let username: string = "";
	let password: string = "";
	let error: string = "";

	onMount(async () => {
		const res = await fetch("http://localhost:8080/api/me", {
			credentials: "include",
		});

		if (res.ok) {
			goto("/dashboard");
			return;
		}
	});

	async function handleLogin() {
		error = "";
		try {
			const res = await fetch("http://localhost:8080/api/login", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ username, password }),
				credentials: "include",
			});

			if (res.ok) {
				const data = await res.json();
				alert(`welcome ${data.username}`);
				goto("/dashboard");
				return;
			} else {
				if (res.status === 403) {
					goto("/login");
					return;
				}
				const errData = await res.json();
				error = errData.error || "error during login";
			}
		} catch (e) {
			error = "network failed";
		}
	}
</script>

<h1>Login</h1>
{#if error}
	<p style="color:red">{error}</p>
{/if}
<form on:submit|preventDefault={handleLogin}>
	<input placeholder="Username" bind:value={username} />
	<input type="password" placeholder="Password" bind:value={password} />
	<button on:click={handleLogin}>Login</button>
	<a href="/signup">create an account</a>
</form>
