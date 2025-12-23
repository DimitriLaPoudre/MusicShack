<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import type { RequestUser } from "$lib/types/request";
	import type { ErrorResponse } from "$lib/types/response";

	let credentials = $state<RequestUser>({ username: "", password: "" });
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

<div class="body">
	<h1>Login</h1>
	<form onsubmit={handleLogin}>
		<div class="form">
			{#if error}
				<p>{error}</p>
			{/if}
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
			<button>Login</button>
		</div>
	</form>
</div>

<style>
	.body {
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-top: 10vh;
		height: 100vh;

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
	}
</style>
