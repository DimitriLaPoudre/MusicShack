import type { StatusResponse } from "$lib/types/response";
import { apiFetch } from "./fetch";

export async function downloadSong(api: string, id: string) {
	let error: string | null = null
	try {
		const data = await apiFetch<StatusResponse>(
			`/users/downloads/song/${api}/${id}`,
			"POST"
		);
		if ("error" in data) {
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
		const data = await apiFetch<StatusResponse>(
			`/users/downloads/album/${api}/${id}`,
			"POST"
		);
		if ("error" in data) {
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
		const data = await apiFetch<StatusResponse>(
			`/users/downloads/artist/${api}/${id}`,
			"POST"
		);
		if ("error" in data) {

			throw new Error(data.error || "Failed to download artist");
		}
	} catch (e) {
		error =
			e instanceof Error
				? e.message
				: "Failed to load download artist";
	}
	return error
}
