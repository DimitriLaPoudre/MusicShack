import type { RequestDownload } from "$lib/types/request";
import type { DownloadListResponse, StatusResponse } from "$lib/types/response";
import { apiFetch } from "./fetch";

export async function download(req: RequestDownload) {
	let error: string | null = null
	try {
		const data = await apiFetch<StatusResponse>(
			`/downloads`,
			"POST",
			req
		);
		if ("error" in data) {
			throw new Error(data.error || "Failed to download " + req.type);
		}
	} catch (e) {
		error =
			e instanceof Error
				? e.message
				: "Failed to load download " + req.type;
	}
	return error
}

export async function loadDownloads() {
	let list = null;
	let error = null;
	try {
		const data =
			await apiFetch<DownloadListResponse>(`/downloads`);
		if ("error" in data) {
			throw new Error(data.error || "Failed to fetch downloads");
		}
		list = data.sort((a, b) => Number(b.id) - Number(a.id));
	} catch (e) {
		error =
			e instanceof Error
				? e.message
				: "Failed to reload download queue";
	}
	return { list, error };
}

export async function retryDownload(id: number) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(
			`/downloads/${id}/retry`,
			"POST",
		);
		if ("error" in data) {
			throw new Error(data.error || "Failed to retry download");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to retry download";
	}
	return error;
}

export async function retryAllDownload() {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(
			`/downloads/retry`,
			"POST",
		);
		if ("error" in data) {
			throw new Error(data.error || "Failed to retry download");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to retry download";
	}
	return error;
}

export async function cancelDownload(id: number) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(
			`/downloads/${id}/cancel`,
			"POST",
		);
		if ("error" in data) {
			throw new Error(data.error || "Failed to cancel download");
		}
	} catch (e) {
		error =
			e instanceof Error ? e.message : "Failed to cancel download";
	}
	return error;
}

export async function doneDownload() {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(
			`/downloads/done`,
			"POST",
		);
		if ("error" in data) {
			throw new Error(data.error || "Failed to cancel download");
		}
	} catch (e) {
		error =
			e instanceof Error ? e.message : "Failed to cancel download";
	}
	return error;
}

export async function deleteDownload(id: number) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(
			`/downloads/${id}`,
			"DELETE",
		);
		if ("error" in data) {
			throw new Error(data.error || "Failed to delete download");
		}
	} catch (e) {
		error =
			e instanceof Error ? e.message : "Failed to delete download";
	}
	return error;
}
