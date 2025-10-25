<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";

	let title = $state("Unreleased");
	let artistName = $state("Unknown");
	let artistId = $state(0);
	let albumTitle = $state("Total");
	let albumId = $state(0);
	let duration = $state("0:00");
	let quality = $state("None");

	afterNavigate(async () => {
		try {
			const res = await fetch(
				`http://localhost:8080/api/song/${page.params.api}/${page.params.id}`,
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			const data = await res.json();
			console.log(data);
			title = data.Title;
			artistName = data.Artist.Name;
			artistId = data.Artist.Id;
			albumTitle = data.Album.Title;
			albumId = data.Album.Id;
			duration = `${Math.floor(data.Duration / 60)}:${(data.Duration % 60).toString().padStart(2, "0")}`;
			quality = data.AudioQuality;
		} catch (e) {
			console.log(e);
		}
	});
</script>

<p>{title}</p>
<a href="/album/{page.params.api}/{albumId}">{albumTitle}</a>
<a href="/artist/{page.params.api}/{artistId}">{artistName}</a>
<p>{duration}</p>
<p>{quality}</p>
