<script lang="ts">
	import { onMount } from "svelte";
	import { afterNavigate, goto } from "$app/navigation";

	let username: string = "";
	let password: string = "";
	let confirmPassword: string = "";
	let error: string = "";

	afterNavigate(async () => {
		const res = await fetch("http://localhost:8080/api/me", {
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
			const res = await fetch("http://localhost:8080/api/signup", {
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
			error = "network failed";
		}
	}
</script>

<div
	style="display: flex; flex-direction: column; align-items: center; gap: 2rem;"
>
	<h1>Signup</h1>
	<form on:submit|preventDefault={handleSignup}>
		<div
			style="display: flex; flex-direction: column; align-items: center; gap: 1rem; background-color: var(--color-background-light); padding: 2rem 3rem; border: 4px solid var(--color-background-dark);"
		>
			{#if error}
				<p
					style="padding: 8px; border: 4px solid var(--color-error-dark); background-color: var(--color-error);"
				>
					{error}
				</p>
			{/if}
			<input
				style="padding: 8px; border: 4px solid var(--color-secondary-dark); background-color: var(--color-secondary);"
				placeholder="Username"
				bind:value={username}
			/>
			<input
				style="padding: 8px; border: 4px solid var(--color-secondary-dark); background-color: var(--color-secondary);"
				type="password"
				placeholder="Password"
				bind:value={password}
			/>
			<input
				style="padding: 8px; border: 4px solid var(--color-secondary-dark); background-color: var(--color-secondary);"
				type="password"
				placeholder="Confirm Password"
				bind:value={confirmPassword}
			/>

			<button
				style="width: 60%; padding: 8px; border: 4px solid var(--color-primary-dark); background-color: var(--color-primary);"
				>Signup</button
			>
			<a href="/login">connect to an existing account</a>
		</div>
	</form>
</div>
