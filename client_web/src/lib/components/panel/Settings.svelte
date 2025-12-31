<script lang="ts">
	import { goto } from "$app/navigation";
	import { apiFetch } from "$lib/functions/fetch";
	import { onMount } from "svelte";
	import { Pencil, Plus, Trash } from "lucide-svelte";
	import type { RequestInstance, RequestUser } from "$lib/types/request";
	import type {
		InstancesResponse,
		StatusResponse,
		UserResponse,
	} from "$lib/types/response";

	let errorUser = $state<null | string>(null);
	let inputUser = $state<RequestUser>({
		username: "",
		password: "",
		bestQuality: true,
	});
	let username = $state<null | string>(null);

	let errorInstances = $state<null | string>(null);
	let inputInstance = $state<RequestInstance>({ url: "" });
	let instances = $state<null | InstancesResponse>(null);

	onMount(() => {
		loadInstance();
		getUser();
	});

	async function getUser() {
		try {
			const data = await apiFetch<UserResponse>("/me");
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch me");
			}
			username = data.username;
			inputUser.bestQuality = data.bestQuality;
			errorUser = null;
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to get user info";
		}
	}

	async function changeUser(event: SubmitEvent) {
		event.preventDefault();
		try {
			const data = await apiFetch<UserResponse>("/me", "PUT", inputUser);
			if ("error" in data) {
				throw new Error(data.error || "Failed to update me");
			}
			username = data.username;
			errorUser = null;
			if (inputUser.username !== "" || inputUser.password !== "") {
				inputUser = {
					username: "",
					password: "",
					bestQuality: data.bestQuality,
				};
				await logout();
			} else {
				inputUser = {
					username: "",
					password: "",
					bestQuality: data.bestQuality,
				};
			}
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to update user info";
		}
	}

	async function loadInstance() {
		try {
			const data = await apiFetch<InstancesResponse>(`/instances`);
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch instances");
			}
			instances = data;
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
			inputInstance.url = inputInstance.url.trim();
			if (inputInstance.url.endsWith("/")) {
				inputInstance.url = inputInstance.url.substring(
					0,
					inputInstance.url.length - 1,
				);
			}

			if (!inputInstance.url) {
				errorInstances = "fill url with valid value";
				return;
			}

			const data = await apiFetch<StatusResponse>(
				`/instances`,
				"POST",
				inputInstance,
			);
			if ("error" in data) {
				throw new Error(
					data.error || "error while trying to delete Instance",
				);
			}

			inputInstance = { url: "" };
			loadInstance();
		} catch (e) {
			errorInstances =
				e instanceof Error ? e.message : "Failed to add instance";
			return;
		}
	}

	async function deleteInstance(id: number) {
		try {
			const data = await apiFetch<StatusResponse>(
				`/instances/${id}`,
				"DELETE",
			);
			if ("error" in data) {
				throw new Error(
					data.error || "error while trying to delete Instance",
				);
			}
			loadInstance();
		} catch (e) {
			errorInstances =
				e instanceof Error ? e.message : "Failed to delete instance";
		}
	}

	async function logout() {
		try {
			const data = await apiFetch<StatusResponse>(`/logout`, "POST");
			if ("error" in data) {
				throw new Error(data.error || "error while trying to logout");
			}
			goto("/login");
		} catch (e) {
			errorUser = e instanceof Error ? e.message : "Failed to logout";
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
			<div class="wrap-form">
				<div class="inputs">
					<input
						placeholder={username}
						bind:value={inputUser.username}
					/>
					<input
						placeholder="Password"
						bind:value={inputUser.password}
					/>
				</div>
				<div class="qualities">
					<button
						class="item"
						class:active={inputUser.bestQuality !== true}
						type="button"
						onclick={() => {
							inputUser.bestQuality = !inputUser.bestQuality
								? true
								: false;
						}}
					>
						Compressed Quality
					</button>
					<button
						class="item"
						class:active={inputUser.bestQuality === true}
						type="button"
						onclick={() => {
							inputUser.bestQuality = inputUser.bestQuality
								? false
								: true;
						}}
					>
						Best possible Quality
					</button>
				</div>
			</div>
			<button class="edit">
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
							<p>{instance.api}</p>
							<p>{instance.url}</p>
							{#if instance.ping === 0}
								<p>failed</p>
							{:else}
								<p>{instance.ping}ms</p>
							{/if}
						</div>
						<button onclick={() => deleteInstance(instance.id)}>
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
		}
		.form {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			.wrap-form {
				display: flex;
				flex-direction: column;
				gap: 8px;

				.inputs {
					display: grid;
					grid-template-columns: 1fr 1fr;
					gap: 8px;
				}
				.qualities {
					display: grid;
					grid-template-columns: 1fr 1fr;
					gap: 8px;

					.item {
						padding: 18.5px 0;
						outline: none;
						border: none;
						background-color: inherit;
						color: inherit;
					}
					.item:hover {
						outline: 1px solid #ffffff;
						outline-offset: -1px;
						border: none;
						background-color: inherit;
						color: inherit;
					}
					.item:active {
						background-color: #ffffff;
						color: #0e0e0e;
					}
					.item.active {
						text-decoration: underline;
					}
				}
			}
			.edit {
				aspect-ratio: 1/1;
			}
			@container (max-width: 420px) {
				.inputs {
					grid-template-columns: 1fr;
				}
				.qualities {
					grid-template-columns: 1fr;
				}
				.edit {
					aspect-ratio: auto;
				}
			}
		}
	}

	.instances {
		display: flex;
		flex-direction: column;
		padding: 8px;
		gap: 16px;
		.error {
			text-align: center;
			background-color: var(--err);
			padding: 0.5rem;
		}
		.form {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			.inputs {
				display: grid;
				grid-template-columns: 1fr;
				gap: 8px;
			}
			button {
				aspect-ratio: 1/1;
			}
			@container (max-width: 420px) {
				.inputs {
					grid-template-columns: 1fr;
				}
				button {
					aspect-ratio: auto;
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
					grid-template-columns: 1fr 1fr auto;
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
					button {
						aspect-ratio: auto;
					}
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
