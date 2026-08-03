import { apiFetch } from "./client";
import type { CreateDownload, DownloadTask } from "$lib/types";

type DownloadKind = "song" | "album" | "artist" | "playlist";

function download(kind: DownloadKind, form: CreateDownload): Promise<void> {
	return apiFetch<void>(`/download/${kind}`, {
		method: "POST",
		body: JSON.stringify(form),
	});
}

export function downloadSong(form: CreateDownload): Promise<void> {
	return download("song", form);
}

export function downloadAlbum(form: CreateDownload): Promise<void> {
	return download("album", form);
}

export function downloadArtist(form: CreateDownload): Promise<void> {
	return download("artist", form);
}

export function downloadPlaylist(form: CreateDownload): Promise<void> {
	return download("playlist", form);
}

export function retryDownload(id: string): Promise<void> {
	return apiFetch<void>(`/download/${encodeURIComponent(id)}/retry`, { method: "PUT" });
}

export function retryAllDownloads(): Promise<void> {
	return apiFetch<void>("/download/retry", { method: "PUT" });
}

export function cancelDownload(id: string): Promise<void> {
	return apiFetch<void>(`/download/${encodeURIComponent(id)}/cancel`, { method: "PUT" });
}

export function removeDoneDownloads(): Promise<void> {
	return apiFetch<void>("/download/remove", { method: "PUT" });
}

export function removeDownload(id: string): Promise<void> {
	return apiFetch<void>(`/download/${encodeURIComponent(id)}/remove`, { method: "PUT" });
}

export function listDownloads(): Promise<DownloadTask[]> {
	return apiFetch<DownloadTask[]>("/download");
}
