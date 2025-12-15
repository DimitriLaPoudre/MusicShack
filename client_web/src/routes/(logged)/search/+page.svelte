<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Disc, DiscAlbum, Download, User } from "lucide-svelte";
	import { PUBLIC_API_URL } from "$env/static/public";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let result = $state<any | null>(null);
	let searchData = $state<string | null>(null);
	let api = $state<string>("");
	let type = $state<string>("songs");

	afterNavigate(async () => {
		try {
			searchData = page.url.searchParams.get("q");
			if (!searchData) {
				throw new Error("No Search");
			}
			const res = await fetch(
				`${PUBLIC_API_URL}/api/search?q=${searchData}`,
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			result = await res.json();
			if (!res.ok) {
				throw new Error(result.error || "Failed to fetch search");
			}
			api = Object.keys(result)[0];
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load song";
			isLoading = false;
		}
	});

	async function downloadSong(api: string, id: string) {
		try {
			const res = await fetch(
				`${PUBLIC_API_URL}/api/users/downloads/song/${api}/${id}`,
				{
					method: "POST",
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || "Failed to download song");
			}
		} catch (e) {
			error =
				e instanceof Error ? e.message : "Failed to load download song";
		}
	}

	async function downloadAlbum(api: string, id: string) {
		try {
			const res = await fetch(
				`${PUBLIC_API_URL}/api/users/downloads/album/${api}/${id}`,
				{
					method: "POST",
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || "Failed to download album");
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Failed to load download album";
		}
	}

	async function downloadArtist(api: string, id: string) {
		try {
			const res = await fetch(
				`${PUBLIC_API_URL}/api/users/downloads/artist/${api}/${id}`,
				{
					method: "POST",
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || "Failed to download artist");
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Failed to load download artist";
		}
	}
</script>

<svelte:head>
	<title>{"Search"} | {searchData} - MusicShack</title>
</svelte:head>

{#if isLoading}
	<p class="loading">Loading...</p>
{:else if error}
	<div class="error">
		<h2>Error Loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else}
	<div class="top">
		<div class="section">
			{#each Object.entries(result as Record<string, any>) as [key, _]}
				<button onclick={() => (api = key)} class:active={api === key}>
					{key}</button
				>
			{/each}
		</div>
		<div class="section">
			<button
				onclick={() => (type = "songs")}
				class:active={type === "songs"}
			>
				Songs</button
			>
			<button
				onclick={() => (type = "albums")}
				class:active={type === "albums"}>Albums</button
			>
			<button
				onclick={() => (type = "artists")}
				class:active={type === "artists"}>Artists</button
			>
		</div>
	</div>
	<div class="items">
		{#if type === "songs"}
			{#each result[api].Songs as song}
				<div class="wrap-item">
					<button
						class="item"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/song/${api}/${song.Id}`);
						}}
					>
						<div class="cover">
							{#if song.CoverUrl !== ""}
								<img src={song.CoverUrl} alt={song.Title} />
							{:else}
								<Disc size={140} />
							{/if}
						</div>
						<p>{song.Title}</p>
						<nav>
							{#each song.Artists as artist}
								<a href="/artist/{api}/{artist.Id}">
									{artist.Name}
								</a>
							{/each}
						</nav>
					</button>
					<button
						class="download"
						onclick={() => downloadSong(api, song.Id)}
					>
						<Download />
					</button>
				</div>
			{/each}
		{:else if type === "albums"}
			{#each result[api].Albums as album}
				<div class="wrap-item">
					<button
						class="item"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/song/${api}/${album.Id}`);
						}}
					>
						<div class="cover">
							{#if album.CoverUrl !== ""}
								<img src={album.CoverUrl} alt={album.Title} />
							{:else}
								<DiscAlbum size={140} />
							{/if}
						</div>
						<p>{album.Title}</p>
						<nav>
							{#each album.Artists as artist}
								<a href="/artist/{api}/{artist.Id}">
									{artist.Name}
								</a>
							{/each}
						</nav>
					</button>
					<button
						class="download"
						onclick={() => downloadAlbum(api, album.Id)}
					>
						<Download />
					</button>
				</div>
			{/each}
		{:else}
			{#each result[api].Artists as artist}
				<button
					class="artist"
					onclick={() => goto(`/artist/${api}/${artist.Id}`)}
				>
					<div class="picture">
						{#if artist.PictureUrl !== ""}
							<img src={artist.PictureUrl} alt={artist.Name} />
						{:else}
							<User size={140} />
						{/if}
					</div>
					<p>{artist.Name}</p>
				</button>
			{/each}
		{/if}
	</div>
{/if}

<style>
	.loading {
		margin-top: 30px;
		text-align: center;
	}
	.error {
		margin-top: 30px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 10px;

		* {
			margin: 0;
		}
	}

	.top {
		display: flex;
		flex-direction: column;
		gap: 10px;
		padding: 10px 0;
		.section {
			display: flex;
			flex-direction: row;
			gap: 10px;
			button {
				padding: 10px;
			}
		}
	}

	.items {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;

		.wrap-item {
			width: 200px;
			height: auto;
			.item {
				display: flex;
				flex-direction: column;
				align-items: center;
				width: 200px;
				height: auto;
				overflow: hidden;
				border-bottom: none;

				.cover {
					width: 160px;
					height: 160px;
				}

				nav {
					display: flex;
					flex-direction: column;
					gap: 0.2rem 1rem;
				}
			}
			.download {
				width: 100%;
				border-top: none;
			}
		}
	}

	.artist {
		width: 200px;
		height: auto;
		.picture {
			width: 160px;
			height: 160px;
			border-radius: 50%;
		}
	}
</style>
