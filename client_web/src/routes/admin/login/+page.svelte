<script lang="ts">
	import { goto, afterNavigate } from "$app/navigation";
	import type { RequestAdmin } from "$lib/types/request";
	import type { ErrorResponse, StatusResponse } from "$lib/types/response";

	let credentials = $state<RequestAdmin>({ password: "" });
	let error = $state<string>("");

	afterNavigate(async () => {
		const res = await fetch("/api/admin", {
			credentials: "include",
		});

		if (res.ok) {
			goto("/admin/dashboard");
			return;
		}
	});

	async function handleLogin(e: SubmitEvent) {
		e.preventDefault();
		try {
			const res = await fetch("/api/admin/login", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(credentials),
				credentials: "include",
			});

			let data;
			if (res.ok) {
				data = (await res.json()) as StatusResponse;
				goto("/admin/dashboard");
				return;
			} else {
				data = (await res.json()) as ErrorResponse;
				if (res.status === 403) {
					goto("/login");
					return;
				}
				error = data.error || "error during login";
			}
		} catch (e) {
			error = e instanceof Error ? e.message : "network failed";
		}
		credentials.password = "";
	}
</script>

<svelte:head>
	<title>Admin | Login - MusicShack</title>
</svelte:head>

<div class="body">
	<h1 class="title">Admin Login</h1>
	{#if error}
		<p class="error">{error}</p>
	{/if}
	<form onsubmit={handleLogin}>
		<div class="form">
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
		margin: 10vh auto;
		height: 100vh;
		width: clamp(320px, 70vw + 20px, 1200px);

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

			input {
				padding: 0.75rem;
			}
			button {
				width: 200px;
				padding: 0.75rem;
			}
		}
	}
</style>
