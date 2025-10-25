<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let album = $state<any | null>(null);

	afterNavigate(async () => {
		try {
			const res = await fetch(
				`http://localhost:8080/api/album/${page.params.api}/${page.params.id}`,
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			if (!res.ok) {
				throw new Error("Failed to fetch album");
			}
			album = await res.json();
			album.Duration = `${Math.floor(album.Duration / 60)}:${(album.Duration % 60).toString().padStart(2, "0")}`;
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load album";
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title
		>{album?.Title || "Album"} | {album?.Artist?.Name || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if isLoading}
	<p>Loading...</p>
{:else if error}
	<h2>Error Loading Album</h2>
	<p>{error}</p>
	<a href="/dashboard"> Go to Dashboard </a>
{:else}
	<div style="display: flex; align-items: flex-start; gap: 10px;">
		<img
			src={album.CoverUrl}
			alt={album.title}
			style="width:300px; height:auto;"
		/>
		<div>
			<p>[{album.Type}]</p>
			<h1>{album.Title}</h1>
			{#each album.Artists as artist}
				<a href="/artist/{page.params.api}/{artist.Id}">{artist.Name}</a
				>
			{/each}
			<p>{album.Duration}</p>
			<p>{album.AudioQuality}</p>
			<p>{album.ReleaseDate}</p>
			<button onclick={() => {}}>Download album</button>
		</div>
	</div>
	<div>
		{#each album.Songs as song}
			<div style="display: flex">
				<button
					style="display: flex; gap: 10px;"
					onclick={() => {
						goto(`/song/${page.params.api}/${song.Id}`);
					}}
				>
					<span>
						{song.TrackNumber}
					</span>
					<span>
						{song.Title}
					</span>
					{#each song.Artists as artist}
						<a href="/artist/{page.params.api}/{artist.Id}">
							{artist.Name}
						</a>
					{/each}
					<span>
						{`${Math.floor(song.Duration / 60)}:${(song.Duration % 60).toString().padStart(2, "0")}`}
					</span>
				</button>
				<button onclick={() => {}}> Download</button>
			</div>
		{/each}
	</div>
{/if}
