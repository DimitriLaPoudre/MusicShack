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

			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch instances");
			}
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

			await res.json();
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

<div style="display: flex; flex-direction: column; gap: 8px; padding: 8px;">
	<div style="display: flex; flex-direction: column; gap: 8px; padding: 8px">
		<p>Instances</p>
		<form
			on:submit|preventDefault={addInstance}
			style="display: flex; flex-direction: row; justify-content: space-between; padding: 4px; gap: 8px; border: 2px solid var(--color-background-dark); background-color: var(--color-background-light);"
		>
			<input
				placeholder="API"
				bind:value={newInstanceAPI}
				style="flex: 1; text-align: left; margin: 0;"
			/>
			<input
				placeholder="URL"
				bind:value={newInstanceURL}
				style="flex: 1; text-align: left; margin: 0;"
			/>
			<button style="margin-left: auto;"><Plus /></button>
		</form>
		{#if newInstanceError}
			<p
				style="border: 4px solid var(--color-error-dark); background-color: var(--color-error);"
			>
				{newInstanceError}
			</p>
		{/if}
		{#if instanceError}
			<p
				style="padding: 8px; border: 2px solid var(--color-error-dark); background-color: var(--color-error);"
			>
				{instanceError}
			</p>
		{:else if !instances}
			<p
				style="text-align: center; padding: 4px; gap: 8px; border: 2px solid var(--color-background-dark); background-color: var(--color-background-light);"
			>
				Loading...
			</p>
		{:else}
			{#each instances as instance}
				<div
					style="display: flex; flex-direction: row; justify-content: space-between; padding: 4px; gap: 8px; border: 2px solid var(--color-background-dark); background-color: var(--color-background-light);"
				>
					<p style="flex: 1; text-align: left; margin: 0;">
						{instance.Api}
					</p>
					<p style="flex: 1; text-align: left; margin: 0;">
						{instance.Url}
					</p>
					<button
						on:click={() => deleteInstance(instance.ID)}
						style="margin-left: auto;"
					>
						<Trash /></button
					>
				</div>
			{/each}
		{/if}
	</div>
	<button
		on:click={Logout}
		style="padding: 4px; border: 2px solid var(--color-error-dark); background-color: var(--color-error-light);"
		>Logout</button
	>
</div>
