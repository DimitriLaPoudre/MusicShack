<script lang="ts">
	import { goto, onNavigate } from "$app/navigation";

	let username: string = "";
	let password: string = "";
	let error: string = "";

	onNavigate(async () => {
		const res = await fetch("http://localhost:8080/api/me", {
			credentials: "include",
		});

		if (res.ok) {
			goto("/");
			return;
		}
	});

	async function handleLogin() {
		try {
			const res = await fetch("http://localhost:8080/api/login", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ username, password }),
				credentials: "include",
			});

			if (res.ok) {
				await res.json();
				goto("/");
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

<div
	style="display: flex; flex-direction: column; align-items: center; gap: 2rem;"
>
	<h1>Login</h1>
	<form on:submit|preventDefault={handleLogin}>
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
				required
			/>
			<input
				style="padding: 8px; border: 4px solid var(--color-secondary-dark); background-color: var(--color-secondary);"
				type="password"
				placeholder="Password"
				bind:value={password}
				required
			/>
			<button
				style="width: 60%; padding: 8px; border: 4px solid var(--color-primary-dark); background-color: var(--color-primary);"
				>Login</button
			>
			<a href="/signup">create an account</a>
		</div>
	</form>
</div>
