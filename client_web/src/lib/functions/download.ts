import type { RequestDownload } from "$lib/types/request";
import type {
	DownloadData,
	DownloadListResponse,
	StatusResponse,
} from "$lib/types/response";
import { apiFetch } from "./fetch";
import { downloadList } from "$lib/stores/panel/download";

export async function download(req: RequestDownload) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(`/downloads`, "POST", req);
		if ("error" in data) {
			throw new Error(data.error || "Failed to download " + req.type);
		}
	} catch (e) {
		error =
			e instanceof Error ? e.message : "Failed to load download " + req.type;
	}
	return error;
}

export async function loadDownloads() {
	let error = null;
	const statusOrder: DownloadData["status"][] = [
		"running",
		"pending",
		"done",
		"failed",
		"cancel",
	];
	try {
		const data = await apiFetch<DownloadListResponse>("/downloads");
		if ("error" in data) {
			throw new Error(data.error || "Failed to fetch downloads");
		}
		downloadList.set(
			data.sort((a, b) => {
				const statusDiff =
					statusOrder.indexOf(a.status) - statusOrder.indexOf(b.status);
				if (statusDiff !== 0) return statusDiff;

				return b.id - a.id;
			}),
		);
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to reload download queue";
	}
	return error;
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
		const data = await apiFetch<StatusResponse>(`/downloads/retry`, "POST");
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
		error = e instanceof Error ? e.message : "Failed to cancel download";
	}
	return error;
}

export async function doneDownload() {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(`/downloads/done`, "POST");
		if ("error" in data) {
			throw new Error(data.error || "Failed to cancel download");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to cancel download";
	}
	return error;
}

export async function deleteDownload(id: number) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(`/downloads/${id}`, "DELETE");
		if ("error" in data) {
			throw new Error(data.error || "Failed to delete download");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to delete download";
	}
	return error;
}
