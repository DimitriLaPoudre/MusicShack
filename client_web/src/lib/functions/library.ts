import type { ResponseSong, StatusResponse } from "$lib/types/response";
import { apiFetch } from "./fetch";

export async function loadLibrary() {
	let list = null;
	let error = null;
	try {
		const data = await apiFetch<ResponseSong[]>(`/library`);
		if ("error" in data) {
			throw new Error(data.error || "Failed to fetch song");
		}
		list = data;
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to load song";
	}
	return { list, error };
}

export async function syncLibrary() {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(`/library`, "PUT");
		if ("error" in data) {
			throw new Error(data.error || "Failed to sync library");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to sync library";
	}
	return error;
}

export async function deleteSong(id: number) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(`/library/${id}`, "DELETE");
		if ("error" in data) {
			throw new Error(data.error || "Failed to delete song");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to delete song";
	}
	return error;
}
