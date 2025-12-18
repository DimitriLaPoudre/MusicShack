import { goto } from "$app/navigation";
import { apiFetch } from "./apiFetch";

export async function downloadSong(api: string, id: string) {
	let error: string | null = null
	try {
		const res = await apiFetch(
			`/users/downloads/song/${api}/${id}`,
			"POST"
		);
		const data = await res.json();
		if (!res.ok) {
			throw new Error(data.error || "Failed to download song");
		}
	} catch (e) {
		error =
			e instanceof Error
				? e.message
				: "Failed to load download song";
	}
	return error
}

export async function downloadAlbum(api: string, id: string) {
	let error: string | null = null
	try {
		const res = await apiFetch(
			`/users/downloads/album/${api}/${id}`,
			"POST"
		);
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
	return error
}


export async function downloadArtist(api: string, id: string) {
	let error: string | null = null
	try {
		const res = await apiFetch(
			`/users/downloads/artist/${api}/${id}`,
			"POST"
		);
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
	return error
}
