<script lang="ts">
	import { goto } from "$app/navigation";
	import "../../app.css";
	import {
		CircleAlert,
		CircleCheck,
		CircleDashed,
		CircleX,
		Disc,
		DownloadIcon,
		Heart,
		LoaderCircleIcon,
		Pencil,
		Plus,
		RotateCcw,
		Search,
		SettingsIcon,
		Trash,
	} from "lucide-svelte";
	import Follow from "$lib/components/panel/Follow.svelte";
	import { apiFetch } from "$lib/functions/apiFetch";

	let { children } = $props();
	let barState = $state<null | string>(null);

	// panel-search variable
	let searchInput: string = "";

	// panel-download variable
	let downloadList = $state<null | any>(null);
	let downloadError = $state<null | string>(null);
	let downloadManageHover = $state<null | number>(null);

	// panel-settingsvariable
	let settingsUserError = $state<null | string>(null);
	let settingsUsername = $state<null | string>(null);
	let settingsUsernameInput = $state<null | string>(null);
	let settingsPasswordInput = $state<null | string>(null);
	let settingsApiInput = $state<null | string>(null);
	let settingsURLInput = $state<null | string>(null);
	let settingsInstanceError = $state<null | string>(null);
	let settingsInstanceList = $state<null | any>(null);

	async function searchFunction() {
		const encodedSearchData = encodeURI(searchInput);
		window.location.assign(`/search?q=${encodedSearchData}`);
	}

	$effect(() => {
		if (barState === "download") {
			const interval = setInterval(() => {
				loadDownloads();
			}, 500);
			return () => clearInterval(interval);
		}
	});

	async function loadDownloads() {
		try {
			const res = await apiFetch(`/users/downloads/`);
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch downloads");
			}
			downloadList = body.tasks
				.slice()
				.sort((a: any, b: any) => Number(b.Id) - Number(a.Id));
			downloadError = null;
		} catch (e) {
			downloadError =
				e instanceof Error
					? e.message
					: "Failed to reload download queue";
		}
	}

	async function retryDownload(id: string) {
		try {
			const res = await apiFetch(`/users/downloads/retry/${id}`, "POST");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to retry download");
			}
		} catch (e) {
			downloadError =
				e instanceof Error ? e.message : "Failed to retry download";
		}
		loadDownloads();
	}

	async function cancelDownload(id: string) {
		try {
			const res = await apiFetch(`/users/downloads/cancel/${id}`, "POST");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to cancel download");
			}
		} catch (e) {
			downloadError =
				e instanceof Error ? e.message : "Failed to cancel download";
		}
		loadDownloads();
	}

	async function deleteDownload(id: string) {
		try {
			const res = await apiFetch(`/users/downloads/${id}`, "DELETE");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to delete download");
			}
		} catch (e) {
			downloadError =
				e instanceof Error ? e.message : "Failed to delete download";
		}
		loadDownloads();
	}

	async function getUser() {
		try {
			const res = await apiFetch("/me/");
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
			const res = await apiFetch(`/instances/`);
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

			const res = await apiFetch(`/instances/`, "POST", {
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
			<Search />
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
				if (barState === "download") {
					barState = null;
				} else {
					barState = "download";
					loadDownloads();
				}
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
			<form onsubmit={searchFunction}>
				<input
					type="text"
					bind:value={searchInput}
					placeholder="Search"
				/>
			</form>
		</div>
	{:else if barState === "follow"}
		<div class="panel-default">
			<Follow />
		</div>
	{:else if barState === "download"}
		<div class="panel-default">
			<h1>Download Queue</h1>
			{#if downloadError}
				<p class="panel-download-error">{downloadError}</p>
			{/if}
			{#if !downloadList}
				<p class="panel-download-loading">Loading...</p>
			{:else}
				<div class="panel-download-items">
					{#each downloadList as download, index}
						<div class="panel-download-item">
							{#if download.Data.Album.CoverUrl !== ""}
								<img
									src={download.Data.Album.CoverUrl}
									alt={download.Data.Album.CoverUrl}
								/>
							{:else}
								<Disc />
							{/if}
							<div class="panel-download-item-data">
								<p>{download.Data.Title}</p>
								<p>{download.Data.Artist.Name}</p>
							</div>
							<div class="panel-download-item-btn">
								{#if download.Status === "done"}
									<button>
										<CircleCheck />
									</button>
								{:else if download.Status === "pending"}
									<button
										onmouseenter={() =>
											(downloadManageHover = index)}
										onmouseleave={() =>
											(downloadManageHover = null)}
										onclick={() => {
											cancelDownload(download.Id);
										}}
									>
										{#if downloadManageHover != null && downloadManageHover === index}
											<CircleX />
										{:else}
											<CircleDashed />
										{/if}
									</button>
								{:else if download.Status === "running"}
									<button
										onmouseenter={() =>
											(downloadManageHover = index)}
										onmouseleave={() =>
											(downloadManageHover = null)}
										onclick={() => {
											cancelDownload(download.Id);
										}}
									>
										{#if downloadManageHover != null && downloadManageHover === index}
											<CircleX />
										{:else}
											<LoaderCircleIcon />
										{/if}
									</button>
								{:else if download.Status === "failed" || download.Status === "cancel"}
									<button
										onmouseenter={() =>
											(downloadManageHover = index)}
										onmouseleave={() =>
											(downloadManageHover = null)}
										onclick={() => {
											retryDownload(download.Id);
										}}
									>
										{#if downloadManageHover != null && downloadManageHover === index}
											<RotateCcw />
										{:else}
											<CircleAlert />
										{/if}
									</button>
								{/if}
								<button
									onclick={() => {
										deleteDownload(download.Id);
									}}
								>
									<Trash />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}
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
		outline: 1px solid #ffffff;
	}

	.panel-search input {
		width: 100%;
		border: none;
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

	.panel-download-loading {
		text-align: center;
	}
	.panel-download-error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
		margin: 0;
	}
	.panel-download-items {
		display: flex;
		flex-direction: column;
		gap: 4px;
		.panel-download-item {
			display: grid;
			grid-template-columns: auto 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			img {
				margin: auto;
				width: 58px;
				height: 58px;
				aspect-ratio: 1/1;
			}

			.panel-download-item-data {
				display: grid;
				grid-template-columns: 1fr 1fr;
				align-items: center;
				p {
					padding-left: 0.5rem;
					margin: 0;
				}
			}
			.panel-download-item-data:hover {
				outline: 1px solid #ffffff;
				outline-offset: -1px;
			}

			.panel-download-item-btn {
				display: grid;
				grid-template-columns: 1fr 1fr;
				button {
					aspect-ratio: 1/1;
				}
			}

			@container (max-width: 420px) {
				.panel-download-item-data {
					grid-template-columns: 1fr;
				}

				.panel-download-item-btn {
					grid-template-columns: 1fr;
				}
			}
		}
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
				p {
					margin: 0;
				}
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
