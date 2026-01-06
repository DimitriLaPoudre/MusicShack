<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import type { RequestUserLogin } from "$lib/types/request";
	import type { ErrorResponse } from "$lib/types/response";

	let credentials = $state<RequestUserLogin>({ username: "", password: "" });
	let error = $state<string>("");

	afterNavigate(async () => {
		const res = await fetch("/api/me", {
			credentials: "include",
		});

		if (res.ok) {
			goto("/dashboard");
			return;
		}
	});

	async function handleLogin(e: SubmitEvent) {
		e.preventDefault();
		try {
			const res = await fetch("/api/login", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(credentials),
				credentials: "include",
			});

			if (res.ok) {
				goto("/");
				return;
			} else {
				if (res.status === 403) {
					goto("/dashboard");
					return;
				}
				const data = (await res.json()) as ErrorResponse;
				throw new Error(data.error || "error during login");
			}
		} catch (e) {
			error = e instanceof Error ? e.message : "network failed";
		}
	}
</script>

<svelte:head>
	<title>Login - MusicShack</title>
</svelte:head>

<div class="body">
	<h1 class="title">Login</h1>
	<form onsubmit={handleLogin}>
		{#if error}
			<p class="error">{error}</p>
		{/if}
		<div class="form">
			<div class="inputs">
				<input
					placeholder="Username"
					bind:value={credentials.username}
					required
				/>
				<input
					type="password"
					placeholder="Password"
					bind:value={credentials.password}
					required
				/>
			</div>
			<button class="hover-full">Login</button>
		</div>
	</form>
</div>

<style>
	.body {
		display: flex;
		flex-direction: column;
		align-items: center;
		margin: 10vh auto;
		width: clamp(320px, 70vw + 20px, 1200px);
		height: 100vh;

		.title {
			font-weight: bolder;
		}
		.error {
			padding: 0.75rem;
			margin: 1rem;
			color: var(--err);
		}
		.form {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 1rem;
			padding: 1rem 2rem;

			.inputs {
				display: flex;
				flex-direction: column;
				gap: 0.5rem;

				input {
					padding: 0.75rem;
				}
			}
			button {
				width: 200px;
				padding: 0.75rem;
			}
		}
	}
</style>
