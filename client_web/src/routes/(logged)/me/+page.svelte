<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { Trash, Plus } from "lucide-svelte";

	let newInstanceAPI = $state("");
	let newInstanceURL = $state("");
	let newInstanceError = $state<string | null>(null);
	let instances = $state<any | null>(null);
	let instanceError = $state<string | null>(null);

	afterNavigate(async () => {
		const res = await fetch("http://localhost:8080/api/me", {
			credentials: "include",
		});

		if (!res.ok) {
			goto("/login");
			return;
		}
		loadInstance();
	});

	async function loadInstance() {
		try {
			const res = await fetch("http://localhost:8080/api/instances", {
				credentials: "include",
			});

			if (res.status === 401) {
				goto("/login");
				return;
			}
			if (!res.ok) {
				throw new Error("Failed to fetch instances");
			}

			const body = await res.json();
			instances = body.instances;
			instanceError = null;
		} catch (e) {
			instanceError = "network failed";
		}
	}

	async function addInstance() {
		try {
			if (newInstanceURL.endsWith("/")) {
				newInstanceURL = newInstanceURL.substring(
					0,
					newInstanceURL.length - 1,
				);
			}
			const res = await fetch("http://localhost:8080/api/instances", {
				method: "POST",
				credentials: "include",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					api: newInstanceAPI,
					url: newInstanceURL,
				}),
			});

			const data = await res.json();
			if (!res.ok) {
				if (res.status === 403) {
					goto("/");
					return;
				}
				newInstanceError = data.error || "error while addInstance";
				return;
			}

			newInstanceAPI = "";
			newInstanceURL = "";
			loadInstance();
		} catch (e) {
			newInstanceError = "network failed";
		}
	}

	async function deleteInstance(id: number) {
		try {
			const res = await fetch(
				`http://localhost:8080/api/instances/${id}`,
				{
					method: "DELETE",
					credentials: "include",
				},
			);

			const data = await res.json();
			if (!res.ok) {
				if (res.status === 403) {
					goto("/");
					return;
				}
				return;
			}

			loadInstance();
		} catch (e) {
			newInstanceError = "network failed";
		}
	}

	async function Logout() {
		await fetch("http://localhost:8080/api/logout", {
			method: "POST",
			credentials: "include",
		});
		goto("/login");
	}
</script>

<h1>User</h1>
<div>
	<div style="display: flex; flex-direction: column; gap: 8px;">
		<p>Instances</p>
		<form
			on:submit|preventDefault={addInstance}
			style="display: flex; flex-direction: row; justify-content: space-between; padding: 4px; border: 2px solid var(--color-background-dark); background-color: var(--color-background-light);"
		>
			<input placeholder="API" bind:value={newInstanceAPI} />
			<input placeholder="URL" bind:value={newInstanceURL} />
			<button><Plus /></button>
			{#if newInstanceError}
				<p
					style="border: 4px solid var(--color-error-dark); background-color: var(--color-error);"
				>
					{newInstanceError}
				</p>
			{/if}
		</form>
		{#if instanceError}
			<p
				style="padding: 8px; border: 4px solid var(--color-error-dark); background-color: var(--color-error);"
			>
				{instanceError}
			</p>
		{:else if !instances}
			<p>Loading...</p>
		{:else}
			{#each instances as instance}
				<div
					style="display: flex; flex-direction: row; justify-content: space-between; padding: 4px;"
				>
					<p>{instance.Api}</p>
					<p>{instance.Url}</p>
					<button on:click={() => deleteInstance(instance.ID)}>
						<Trash /></button
					>
				</div>
			{/each}
		{/if}
	</div>
	<button on:click={Logout} style="color: red">Logout</button>
</div>
