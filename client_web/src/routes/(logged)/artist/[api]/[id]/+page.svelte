<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let artist = $state<any | null>(null);

	afterNavigate(async () => {
		try {
			const res = await fetch(
				`http://localhost:8080/api/artist/${page.params.api}/${page.params.id}`,
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			artist = await res.json();
			if (!res.ok) {
				throw new Error(artist.error || "Failed to fetch artist");
			}
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load artist";
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title>{artist?.Name || "Artist"} - MusicShack</title>
</svelte:head>

{#if isLoading}
	<p>Loading...</p>
{:else if error}
	<h2>Error Loading Artist</h2>
	<p>{error}</p>
	<a href="/">Go to Home</a>
{:else}
	<!-- page top -->
	<div style="display: flex; flex-direction: row; gap: 10px;">
		<img
			src={artist.PictureUrl}
			alt={artist.Name}
			style="width:200px; height:auto;"
		/>
		<div style="display: flex; flex-direction: column; gap: 10px">
			<h1>{artist.Name}</h1>
			<div style="display:flex; flex-direction: row; gap: 10px">
				<button onclick={() => {}}>Favorite</button>
				<button onclick={() => {}}>Download discography</button>
			</div>
		</div>
	</div>
	<!-- page body -->
	<div>
		<div>
			<h2>Albums</h2>
			<div style="display: flex; flex-wrap: wrap; gap: 10px;">
				{#each artist.Albums as album}
					<button
						onclick={() => {
							goto(`/album/${page.params.api}/${album.Id}`);
						}}
						style="display: flex; flex-direction: column; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; width:160px; height: 200px;"
					>
						<img
							src={album.CoverUrl}
							alt={album.Title}
							style="max-width:100%; max-height:100%;"
						/>
						<p>{album.Title}</p>
						{#each album.Artists as artist}
							<a href="/artist/{page.params.api}/{artist.Id}">
								{artist.Name}
							</a>
						{/each}
					</button>
				{/each}
			</div>
		</div>
		<div>
			<h2>EPs</h2>
			<div style="display: flex; flex-wrap: wrap; gap: 10px;">
				{#each artist.Ep as ep}
					<button
						onclick={() => {
							goto(`/album/${page.params.api}/${ep.Id}`);
						}}
						style="display: flex; flex-direction: column; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; width:160px; height: 200px;"
					>
						<img
							src={ep.CoverUrl}
							alt={ep.Title}
							style="max-width:100%; max-height:100%;"
						/>
						<p>{ep.Title}</p>
						{#each ep.Artists as artist}
							<a href="/artist/{page.params.api}/{artist.Id}">
								{artist.Name}
							</a>
						{/each}
					</button>
				{/each}
			</div>
		</div>
		<div>
			<h2>Singles</h2>
			<div style="display: flex; flex-wrap: wrap; gap: 10px;">
				{#each artist.Singles as single}
					<button
						onclick={() => {
							goto(`/album/${page.params.api}/${single.Id}`);
						}}
						style="display: flex; flex-direction: column; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; width:160px; height: 200px;"
					>
						<img
							src={single.CoverUrl}
							alt={single.Title}
							style="max-width:100%; max-height:100%;"
						/>
						<p>{single.Title}</p>
						{#each single.Artists as artist}
							<a href="/artist/{page.params.api}/{artist.Id}">
								{artist.Name}
							</a>
						{/each}
					</button>
				{/each}
			</div>
		</div>
	</div>
{/if}

<style>
	h1 {
		margin: 0;
	}
	p {
		margin: 0;
	}
</style>
