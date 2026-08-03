import { apiFetch } from "./client";
import type { CreateFollow, Follow, StatusResponse } from "$lib/types";

export function addFollow(form: CreateFollow): Promise<Follow> {
	return apiFetch<Follow>("/follows/artist", {
		method: "POST",
		body: JSON.stringify(form),
	});
}

export function listFollows(): Promise<Follow[]> {
	return apiFetch<Follow[]>("/follows");
}

export function deleteFollow(id: string): Promise<StatusResponse> {
	return apiFetch<StatusResponse>(`/follows/${encodeURIComponent(id)}`, {
		method: "DELETE",
	});
}
