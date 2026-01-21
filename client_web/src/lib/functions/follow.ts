import { followList } from "$lib/stores/panel/follow";
import type { RequestFollow } from "$lib/types/request";
import type { FollowsResponse, StatusResponse } from "$lib/types/response";
import { apiFetch } from "./fetch";

export async function loadFollows() {
	let error = null;
	try {
		const data = await apiFetch<FollowsResponse>("/follows");
		if ("error" in data) {
			throw new Error(data.error || "Failed to fetch follows");
		}
		followList.set(data.sort((a, b) => Number(b.id) - Number(a.id)));
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to reload follows list";
	}
	return error;
}

export async function addFollow(req: RequestFollow) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>("/follows", "POST", req);
		if ("error" in data) {
			throw new Error(data.error || "Failed to add follow");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to add new follow";
	}
	return error;
}

export async function removeFollow(id: number) {
	let error = null;
	try {
		const data = await apiFetch<StatusResponse>(`/follows/${id}`, "DELETE");
		if ("error" in data) {
			throw new Error(data.error || "Failed to add follow");
		}
	} catch (e) {
		error = e instanceof Error ? e.message : "Failed to add new follow";
	}
	return error;
}
