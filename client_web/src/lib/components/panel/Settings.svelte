<script lang="ts">
	import { goto } from "$app/navigation";
	import { apiFetch } from "$lib/functions/fetch";
	import { onMount } from "svelte";
	import { Pencil, Plus, Trash } from "lucide-svelte";
	import type { RequestInstance, RequestUser } from "$lib/types/request";

	let errorUser = $state<null | string>(null);
	let inputUser = $state<RequestUser>({ username: "", password: "" });
	let username = $state<null | string>(null);

	let errorInstances = $state<null | string>(null);
	let inputInstance = $state<RequestInstance>({ api: "", url: "" });
	let instances = $state<null | any>(null);

	onMount(() => {
		loadInstance();
		getUser();
	});

	async function getUser() {
		try {
			const res = await apiFetch("/me");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch me");
			}
			username = body.user.Username;
			errorUser = null;
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to get user info";
		}
	}

	async function changeUser(event: SubmitEvent) {
		event.preventDefault();
		try {
			const res = await apiFetch("/me", "PUT", inputUser);
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to update me");
			}
			username = body.user.Username;
			errorUser = null;
			inputUser = { username: "", password: "" };
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to update user info";
		}
	}

	async function loadInstance() {
		try {
			const res = await apiFetch(`/instances`);
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch instances");
			}
			instances = body;
			errorInstances = null;
		} catch (e) {
			errorInstances =
				e instanceof Error
					? e.message
					: "Failed to reload instances queue";
		}
	}

	async function addInstance(event: SubmitEvent) {
		event.preventDefault();
		try {
			if (!inputInstance.api || !inputInstance.url) {
				errorInstances = "fill all fields";
				return;
			}

			inputInstance.api = inputInstance.api.trim();

			inputInstance.url = inputInstance.url.trim();
			if (inputInstance.url.endsWith("/")) {
				inputInstance.url = inputInstance.url.substring(
					0,
					inputInstance.url.length - 1,
				);
			}

			if (!inputInstance.api || !inputInstance.url) {
				errorInstances = "fill fields with valid value";
				return;
			}

			const res = await apiFetch(`/instances`, "POST", inputInstance);
			const data = await res.json();
			if (!res.ok) {
				throw new Error(
					data.error || "error while trying to delete Instance",
				);
			}

			inputInstance = { api: "", url: "" };
		} catch (e) {
			errorInstances =
				e instanceof Error ? e.message : "Failed to add instance";
			return;
		}
		loadInstance();
	}

	async function deleteInstance(id: number) {
		try {
			const res = await apiFetch(`/instances/${id}`, "DELETE");
			const data = await res.json();
			if (!res.ok) {
				throw new Error(
					data.error || "error while trying to delete Instance",
				);
			}
		} catch (e) {
			errorInstances =
				e instanceof Error ? e.message : "Failed to delete instance";
			return;
		}
		loadInstance();
	}

	async function logout() {
		try {
			const res = await apiFetch(`/logout`, "POST");
			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || "error while trying to logout");
			}
			goto("/login");
		} catch (e) {
			errorUser = e instanceof Error ? e.message : "Failed to logout";
			return;
		}
	}
</script>

<div class="body">
	<h1>Settings</h1>
	<h2>User</h2>
	<div class="user">
		{#if errorUser}
			<p class="error">
				{errorUser}
			</p>
		{/if}
		<form class="form" onsubmit={changeUser}>
			<div class="inputs">
				<input placeholder={username} bind:value={inputUser.username} />
				<input placeholder="Password" bind:value={inputUser.password} />
			</div>
			<button>
				<Pencil />
			</button>
		</form>
	</div>
	<h2>Instances</h2>
	<div class="instances">
		{#if errorInstances}
			<p class="error">
				{errorInstances}
			</p>
		{/if}
		<form class="form" onsubmit={addInstance}>
			<div class="inputs">
				<input placeholder="API" bind:value={inputInstance.api} />
				<input placeholder="URL" bind:value={inputInstance.url} />
			</div>
			<button><Plus /></button>
		</form>
		{#if !instances}
			<p class="loading">Loading...</p>
		{:else}
			<div class="items">
				{#each instances as instance}
					<div class="item">
						<div class="data">
							<p>{instance.Api}</p>
							<p>{instance.Url}</p>
						</div>
						<button onclick={() => deleteInstance(instance.Id)}>
							<Trash />
						</button>
					</div>
				{/each}
			</div>
		{/if}
	</div>
	<button class="logout" onclick={logout}> Logout </button>
</div>

<style>
	.user {
		display: flex;
		flex-direction: column;
		padding: 8px;
		gap: 8px;
		.error {
			text-align: center;
			background-color: var(--err);
			padding: 0.5rem;
			margin: 0;
		}
		.form {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			.inputs {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 8px;
			}
			button {
				aspect-ratio: 1/1;
			}

			@container (max-width: 420px) {
				.inputs {
					grid-template-columns: 1fr;
				}
			}
		}
	}

	.instances {
		display: flex;
		flex-direction: column;
		padding: 8px;
		gap: 16px;
	}
	.error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
		margin: 0;
	}
	.form {
		display: grid;
		grid-template-columns: 1fr auto;
		gap: 8px;
		align-items: stretch;
		container-type: inline-size;

		.inputs {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 8px;
		}
		button {
			aspect-ratio: 1/1;
		}

		@container (max-width: 420px) {
			.inputs {
				grid-template-columns: 1fr;
			}
		}
	}
	.loading {
		text-align: center;
	}
	.items {
		display: flex;
		flex-direction: column;
		gap: 4px;
		.item {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			.data {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 8px;
				align-items: center;
				padding: 1rem;
			}
			.data:hover {
				outline: 1px solid #ffffff;
				outline-offset: -1px;
			}
			button {
				aspect-ratio: 1/1;
			}
			@container (max-width: 420px) {
				.data {
					grid-template-columns: 1fr;
				}
			}
		}
	}

	.logout {
		width: 100%;
		padding: 8px;
		border-color: var(--err);
	}
	.logout:hover {
		background-color: var(--err);
	}
</style>
