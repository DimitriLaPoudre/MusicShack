<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { env } from "$env/dynamic/public";

	let username: string = "";
	let password: string = "";
	let confirmPassword: string = "";
	let error: string = "";

	afterNavigate(async () => {
		const res = await fetch(`${env.PUBLIC_API_URL}/api/me`, {
			credentials: "include",
		});

		if (res.ok) {
			goto("/");
			return;
		}
	});

	async function handleSignup() {
		if (password !== confirmPassword) {
			error = "password diff";
			return;
		}
		try {
			const res = await fetch(`${env.PUBLIC_API_URL}/api/signup`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ username, password }),
				credentials: "include",
			});

			if (res.ok) {
				await res.json();
				goto("/login");
				return;
			} else {
				if (res.status === 403) {
					goto("/");
					return;
				}
				const errData = await res.json();
				error = errData.error || "error while signup";
			}
		} catch (e) {
			error = e instanceof Error ? e.message : "network failed";
		}
	}
</script>

<div class="body">
	<h1>Signup</h1>
	<form on:submit|preventDefault={handleSignup}>
		<div class="form">
			{#if error}
				<p>{error}</p>
			{/if}
			<input placeholder="Username" bind:value={username} required />
			<input
				type="password"
				placeholder="Password"
				bind:value={password}
				required
			/>
			<input
				type="password"
				placeholder="Confirm Password"
				bind:value={confirmPassword}
				required
			/>
			<button>Signup</button>
		</div>
	</form>
	<a href="/login">connect to an existing account</a>
</div>

<style>
	.body {
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-top: 10vh;
		height: 100vh;

		h1 {
			margin: 0;
		}
	}
	.form {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		padding: 1rem 2rem;

		p {
			padding: 8px;
			margin: 10px;
			color: var(--err);
		}

		input {
			padding: 8px;
		}

		button {
			width: 60%;
			padding: 8px;
		}
	}
</style>
