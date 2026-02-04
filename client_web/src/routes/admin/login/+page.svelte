<script lang="ts">
	import { goto } from "$app/navigation";
	import type { RequestAdmin } from "$lib/types/request";
	import type { ErrorResponse, StatusResponse } from "$lib/types/response";
	import { onMount } from "svelte";

	let credentials = $state<RequestAdmin>({ password: "" });
	let error = $state<string>("");

	onMount(async () => {
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

<div class="flex flex-col items-center mx-auto mt-[10vh] h-screen w-[clamp(320px,70vw+20px,1200px)]">
	<h1 class="font-extrabold">Admin Login</h1>
	{#if error}
		<p class="p-3 m-4 text-err">{error}</p>
	{/if}
	<form onsubmit={handleLogin}>
		<div class="flex flex-col items-center gap-4 px-8 py-4">
			<input
				class="p-3"
				type="password"
				placeholder="Password"
				bind:value={credentials.password}
				required
			/>
			<button class="hover-full w-[200px] p-3">Login</button>
		</div>
	</form>
</div>
