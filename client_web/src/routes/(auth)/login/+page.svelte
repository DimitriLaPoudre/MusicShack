<script lang="ts">
	import { goto } from "$app/navigation";
	import type { RequestUserLogin } from "$lib/types/request";
	import type { ErrorResponse } from "$lib/types/response";
	import { onMount } from "svelte";

	let credentials = $state<RequestUserLogin>({ username: "", password: "" });
	let error = $state<string>("");

	onMount(async () => {
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

<div class="flex flex-col items-center mx-auto mt-[10vh] w-[clamp(320px,70vw+20px,1200px)] h-screen">
	<h1 class="font-extrabold">Login</h1>
	<form onsubmit={handleLogin}>
		{#if error}
			<p class="p-3 m-4 text-err">{error}</p>
		{/if}
		<div class="flex flex-col items-center gap-4 px-8 py-4">
			<div class="flex flex-col gap-2">
				<input
					class="p-3"
					placeholder="Username"
					bind:value={credentials.username}
					required
				/>
				<input
					class="p-3"
					type="password"
					placeholder="Password"
					bind:value={credentials.password}
					required
				/>
			</div>
			<button class="hover-full w-[200px] p-3">Login</button>
		</div>
	</form>
</div>
