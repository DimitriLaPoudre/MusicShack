import { libraryPage } from "$lib/stores/panel/library";
import type { ResponseLibrary, StatusResponse } from "$lib/types/response";
import { apiFetch } from "./fetch";

export async function loadLibrary(
	search: string,
	limit: number,
	offset: number,
) {
	let error = null;
	try {
		const data = await apiFetch<ResponseLibrary>(
			`/library?q=${search}&limit=${limit}&offset=${offset}`,
		);
		libraryPage.set(data);
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to load song";
	}
	return error;
}

export async function syncLibrary() {
	let error = null;
	try {
		await apiFetch<StatusResponse>(`/library`, "PUT");
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to sync library";
	}
	return error;
}

export async function deleteSong(id: number) {
	let error = null;
	try {
		await apiFetch<StatusResponse>(`/library/${id}`, "DELETE");
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to delete song";
	}
	return error;
}
