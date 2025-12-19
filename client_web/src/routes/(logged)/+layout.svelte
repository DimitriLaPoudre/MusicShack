<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import "../../app.css";
	import {
		DownloadIcon,
		Heart,
		Pencil,
		Plus,
		Search as SearchIcon,
		SettingsIcon,
		Trash,
	} from "lucide-svelte";
	import Follow from "$lib/components/panel/Follow.svelte";
	import { apiFetch } from "$lib/functions/apiFetch";
	import Downloads from "$lib/components/panel/Downloads.svelte";
	import Search from "$lib/components/panel/Search.svelte";

	let { children } = $props();
	let barState = $state<null | string>(null);

	// panel-settingsvariable
	let settingsUserError = $state<null | string>(null);
	let settingsUsername = $state<null | string>(null);
	let settingsUsernameInput = $state<null | string>(null);
	let settingsPasswordInput = $state<null | string>(null);
	let settingsApiInput = $state<null | string>(null);
	let settingsURLInput = $state<null | string>(null);
	let settingsInstanceError = $state<null | string>(null);
	let settingsInstanceList = $state<null | any>(null);

	afterNavigate(() => {
		barState = null;
	});

	async function getUser() {
		try {
			const res = await apiFetch("/me");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch me");
			}
			settingsUsername = body.user.Username;
			settingsUserError = null;
		} catch (e) {
			settingsUserError =
				e instanceof Error ? e.message : "Failed to get user info";
		}
	}

	async function changeUser(event: SubmitEvent) {
		event.preventDefault();
		try {
			const res = await apiFetch("/me", "PUT", {
				username: settingsUsernameInput,
				password: settingsPasswordInput,
			});
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to update me");
			}
			settingsUsername = body.user.Username;
			settingsUserError = null;
			settingsUsernameInput = null;
		} catch (e) {
			settingsUserError =
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
			settingsInstanceList = body.instances;
			settingsInstanceError = null;
		} catch (e) {
			settingsInstanceError =
				e instanceof Error
					? e.message
					: "Failed to reload instances queue";
		}
	}

	async function addInstance(event: SubmitEvent) {
		event.preventDefault();
		try {
			if (!settingsApiInput || !settingsURLInput) {
				settingsInstanceError = "fill all fields";
				return;
			}
			settingsApiInput = settingsApiInput.trim();
			settingsURLInput = settingsURLInput.trim();
			if (settingsURLInput.endsWith("/")) {
				settingsURLInput = settingsURLInput.substring(
					0,
					settingsURLInput.length - 1,
				);
			}

			const res = await apiFetch(`/instances`, "POST", {
				api: settingsApiInput,
				url: settingsURLInput,
			});
			const data = await res.json();
			if (!res.ok) {
				throw new Error(
					data.error || "error while trying to delete Instance",
				);
			}

			settingsApiInput = "";
			settingsURLInput = "";
		} catch (e) {
			settingsInstanceError =
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
			settingsInstanceError =
				e instanceof Error ? e.message : "Failed to delete instance";
			return;
		}
		loadInstance();
	}

	async function Logout() {
		await apiFetch(`/logout`, "POST");
		goto("/login");
	}
</script>

<header>
	<button
		class="logo"
		onclick={() => {
			barState = null;
			goto("/dashboard");
		}}
	>
		MusicShack
	</button>
	<div class="bar">
		<button
			onclick={() => {
				barState = barState === "search" ? null : "search";
			}}
			class:active={barState === "search"}
		>
			<SearchIcon />
		</button>
		<button
			onclick={() => {
				barState = barState === "follow" ? null : "follow";
			}}
			class:active={barState === "follow"}
		>
			<Heart />
		</button>
		<button
			onclick={() => {
				barState = barState === "download" ? null : "download";
			}}
			class:active={barState === "download"}
		>
			<DownloadIcon />
		</button>
		<button
			onclick={() => {
				if (barState === "settings") {
					barState = null;
				} else {
					barState = "settings";
					loadInstance();
					getUser();
				}
			}}
			class:active={barState === "settings"}
		>
			<SettingsIcon />
		</button>
	</div>

	{#if barState === "search"}
		<div class="panel-search">
			<Search />
		</div>
	{:else if barState === "follow"}
		<div class="panel-default">
			<Follow />
		</div>
	{:else if barState === "download"}
		<div class="panel-default">
			<Downloads />
		</div>
	{:else if barState === "settings"}
		<div class="panel-default">
			<h1>Settings</h1>
			<h2>User</h2>
			<div class="panel-settings-user">
				{#if settingsUserError}
					<p class="panel-settings-user-error">
						{settingsUserError}
					</p>
				{/if}
				<form class="panel-settings-user-form" onsubmit={changeUser}>
					<div class="panel-settings-user-section-inputs">
						<input
							placeholder={settingsUsername}
							bind:value={settingsUsernameInput}
						/>
						<input
							placeholder="Password"
							bind:value={settingsPasswordInput}
						/>
					</div>
					<button>
						<Pencil />
					</button>
				</form>
			</div>
			<h2>Instances</h2>
			<div class="panel-settings-instances">
				{#if settingsInstanceError}
					<p class="panel-settings-instances-error">
						{settingsInstanceError}
					</p>
				{/if}
				<form
					class="panel-settings-instances-form"
					onsubmit={addInstance}
				>
					<div class="inputs">
						<input
							placeholder="API"
							bind:value={settingsApiInput}
						/>
						<input
							placeholder="URL"
							bind:value={settingsURLInput}
						/>
					</div>
					<button><Plus /></button>
				</form>
				{#if !settingsInstanceList}
					<p class="panel-settings-instances-loading">Loading...</p>
				{:else}
					<div class="panel-settings-instances-items">
						{#each settingsInstanceList as instance}
							<div class="panel-settings-instances-item">
								<div class="panel-settings-instances-item-data">
									<p>{instance.Api}</p>
									<p>{instance.Url}</p>
								</div>
								<button
									onclick={() => deleteInstance(instance.ID)}
								>
									<Trash />
								</button>
							</div>
						{/each}
					</div>
				{/if}
			</div>
			<button class="panel-settings-logout" onclick={Logout}>
				Logout
			</button>
		</div>
	{/if}
</header>

<main>
	{@render children?.()}
</main>

<style>
	header {
		flex-wrap: wrap;
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		background-color: #0e0e0e;
		color: #fff;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
		z-index: 1000;
	}

	main {
		padding-top: 130px;
		width: 70vw;
		margin: 0 auto;
	}

	.logo {
		text-transform: uppercase;
		font-size: 3rem;
		border: none;
		padding: 10px;
	}

	.logo:hover {
		background-color: inherit;
		color: inherit;
	}

	.bar {
		display: flex;
		flex-direction: row;
		gap: 0px 10px;
		button {
			aspect-ratio: 1/1;
		}
	}

	.panel-search {
		width: 70vw;
		max-height: calc(95vh - 135px);
		overflow-y: auto;
		outline: 1px solid #ffffff;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.panel-default {
		width: 70vw;
		max-height: calc(95vh - 135px);
		overflow-y: auto;
		outline: 1px solid #ffffff;
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 8px;
	}

	.panel-settings-user {
		display: flex;
		flex-direction: column;
		padding: 8px;
		gap: 8px;
		.panel-settings-user-error {
			text-align: center;
			background-color: var(--err);
			padding: 0.5rem;
			margin: 0;
		}
		.panel-settings-user-form {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			.panel-settings-user-section-inputs {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 8px;
			}
			button {
				aspect-ratio: 1/1;
			}

			@container (max-width: 420px) {
				.panel-settings-user-section-inputs {
					grid-template-columns: 1fr;
				}
			}
		}
	}

	.panel-settings-instances {
		display: flex;
		flex-direction: column;
		padding: 8px;
		gap: 16px;
	}
	.panel-settings-instances-error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
		margin: 0;
	}
	.panel-settings-instances-form {
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
	.panel-settings-instances-loading {
		text-align: center;
	}
	.panel-settings-instances-items {
		display: flex;
		flex-direction: column;
		gap: 4px;
		.panel-settings-instances-item {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			.panel-settings-instances-item-data {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 8px;
				align-items: center;
				padding: 1rem;
			}
			.panel-settings-instances-item-data:hover {
				outline: 1px solid #ffffff;
				outline-offset: -1px;
			}
			button {
				aspect-ratio: 1/1;
			}
			@container (max-width: 420px) {
				.panel-settings-instances-item-data {
					grid-template-columns: 1fr;
				}
			}
		}
	}

	.panel-settings-logout {
		padding: 4px;
		border-color: var(--err);
	}
	.panel-settings-logout:hover {
		background-color: var(--err);
	}
</style>
