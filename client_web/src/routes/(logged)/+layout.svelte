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
		Plus,
		RotateCcw,
		Search,
		SettingsIcon,
		Trash,
	} from "lucide-svelte";

	let { children } = $props();
	let barState = $state<null | string>(null);

	// panel-search variable
	let searchInput: string = "";

	// panel-download variable
	let downloadList = $state<null | any>(null);
	let downloadError = $state<null | string>(null);
	let downloadManageHover = $state<null | number>(null);

	// panel-settingsvariable
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
			const res = await fetch(
				"http://localhost:8080/api/users/downloads/",
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}

			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch downloads");
			}
			downloadList = body.tasks
				.slice()
				.sort((a: any, b: any) => Number(b.Id) - Number(a.Id));
			downloadError = null;
		} catch (e) {
			downloadError = "network failed";
		}
	}

	async function retryDownload(id: string) {
		try {
			const res = await fetch(
				`http://localhost:8080/api/users/downloads/retry/${id}`,
				{
					method: "POST",
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			if (res.status === 403) {
				return;
			}

			loadDownloads();
		} catch (e) {
			downloadError = "network failed";
		}
	}

	async function cancelDownload(id: string) {
		try {
			const res = await fetch(
				`http://localhost:8080/api/users/downloads/cancel/${id}`,
				{
					method: "POST",
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			if (res.status === 403) {
				return;
			}

			loadDownloads();
		} catch (e) {
			downloadError = "network failed";
		}
	}

	async function deleteDownload(id: string) {
		try {
			const res = await fetch(
				`http://localhost:8080/api/users/downloads/${id}`,
				{
					method: "DELETE",
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			if (res.status === 403) {
				return;
			}

			loadDownloads();
		} catch (e) {
			downloadError = "network failed";
		}
	}

	async function loadInstance() {
		try {
			const res = await fetch("http://localhost:8080/api/instances/", {
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
			settingsInstanceList = body.instances;
			settingsInstanceError = null;
		} catch (e) {
			settingsInstanceError = "network failed";
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
			const res = await fetch("http://localhost:8080/api/instances/", {
				method: "POST",
				credentials: "include",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					api: settingsApiInput,
					url: settingsURLInput,
				}),
			});

			const data = await res.json();
			if (!res.ok) {
				settingsInstanceError =
					data.error || "error while trying to add Instance";
				return;
			}

			settingsApiInput = "";
			settingsURLInput = "";
			loadInstance();
		} catch (e) {
			settingsInstanceError = "network failed";
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
				settingsInstanceError =
					data.error || "error while trying to delete Instance";
				return;
			}

			loadInstance();
		} catch (e) {
			settingsInstanceError = "network failed";
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
			<h1>Followed Artist</h1>
			<!-- {#if downloadError} -->
			<!-- 	<p class="panel-download-error">{downloadError}</p> -->
			<!-- {/if} -->
			<!-- {#if !downloadList} -->
			<!-- 	<p class="panel-download-loading">Loading...</p> -->
			<!-- {:else} -->
			<!-- 	<div class="panel-download-items"> -->
			<!-- 		{#each downloadList as download} -->
			<!-- 			<div class="panel-download-item"> -->
			<!-- 				<p>{download.Data.Title}</p> -->
			<!-- 				<p>{download.Data.Artist.Name}</p> -->
			<!-- 				<p>{download.Status}</p> -->
			<!-- 				<button -->
			<!-- 					onclick={() => { -->
			<!-- 						deleteDownload(download.Id); -->
			<!-- 					}} -->
			<!-- 				> -->
			<!-- 					<Trash /> -->
			<!-- 				</button> -->
			<!-- 			</div> -->
			<!-- 		{/each} -->
			<!-- 	</div> -->
			<!-- {/if} -->
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
									style="width: 64px; height: 64px;"
								/>
							{:else}
								<Disc
									style="width: auto; height: auto; aspect-ratio: 1;"
								/>
							{/if}
							<p>{download.Data.Title}</p>
							<p>{download.Data.Artist.Name}</p>
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
					{/each}
				</div>
			{/if}
		</div>
	{:else if barState === "settings"}
		<div class="panel-default">
			<h1>Settings</h1>
			<h2>User</h2>
			<div class="panel-settings-user"></div>
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
					<input placeholder="API" bind:value={settingsApiInput} />
					<input placeholder="URL" bind:value={settingsURLInput} />
					<button><Plus /></button>
				</form>
				{#if !settingsInstanceList}
					<p class="panel-settings-instances-loading">Loading...</p>
				{:else}
					<div class="panel-settings-instances-items">
						{#each settingsInstanceList as instance}
							<div class="panel-settings-instances-item">
								<p>{instance.Api}</p>
								<p>{instance.Url}</p>
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
		max-height: 85vh;
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
	}
	.panel-download-item {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		p {
			flex: 1;
			text-align: left;
			padding-left: 8px;
			margin: auto;
		}
		button {
			margin: auto;
			aspect-ratio: 1/1;
		}
	}

	.panel-settings-instances {
		display: flex;
		flex-direction: column;
		padding: 8px;
		gap: 16px;
		h2 {
			margin: 0;
		}
	}
	.panel-settings-instances-error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
		margin: 0;
	}
	.panel-settings-instances-form {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		gap: 8px;

		input {
			flex: 1;
			text-align: left;
			margin: 0;
		}
		button {
			margin: auto;
			aspect-ratio: 1/1;
		}
	}
	.panel-settings-instances-loading {
		text-align: center;
	}
	.panel-settings-instances-items {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}
	.panel-settings-instances-item {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		p {
			flex: 1;
			text-align: left;
			padding-left: 8px;
			margin: auto;
		}
		button {
			margin: auto;
			aspect-ratio: 1/1;
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
