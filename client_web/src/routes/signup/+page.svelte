<script lang="ts">
	import { goto } from "$app/navigation";
	let username: string = "";
	let password: string = "";
	let confirmPassword: string = "";
	let error: string = "";

	async function handleSignup() {
		error = "";
		if (password !== confirmPassword) {
			error = "password diff";
			return;
		}
		try {
			const res = await fetch("http://localhost:8080/api/signup", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ username, password }),
			});

			if (res.ok) {
				const data = await res.json();
				alert(`account create for ${data.username}`);
				goto("/login");
			} else {
				const errData = await res.json();
				error = errData.error || "error while signup";
			}
		} catch (e) {
			error = "network failed";
		}
	}
</script>

<h1>Signup</h1>
{#if error}
	<p style="color:red">{error}</p>
{/if}
<form on:submit|preventDefault={handleSignup}>
	<input placeholder="Username" bind:value={username} />
	<input type="password" placeholder="Password" bind:value={password} />
	<input
		type="password"
		placeholder="Confirm Password"
		bind:value={confirmPassword}
	/>
	<button type="submit" on:click={handleSignup}>Signup</button>
	<a href="/login">connect to an existing account</a>
</form>
