import { apiFetch } from "./fetch";

export async function loadFollows() {
	let list = null
	let error = null
	try {
		const res = await apiFetch("/follows");
		const body = await res.json();
		if (!res.ok) {
			throw new Error(body.error || "Failed to fetch follows");
		}
		list = body;
	} catch (e) {
		error =
			e instanceof Error
				? e.message : "Failed to reload follows list";
	}
	return { list, error }
}

export async function addFollow(api: string, id: string) {
	let error = null
	try {
		const res = await apiFetch("/follows", "POST", {
			api, id
		});
		const body = await res.json();
		if (!res.ok) {
			throw new Error(body.error || "Failed to add follow");
		}
	} catch (e) {
		error =
			e instanceof Error
				? e.message
				: "Failed to add new follow";
	}
	return error
}

export async function removeFollow(id: string) {
	let error = null
	try {
		const res = await apiFetch(`/follows/${id}`, "DELETE");
		const body = await res.json();
		if (!res.ok) {
			throw new Error(body.error || "Failed to add follow");
		}
	} catch (e) {
		error =
			e instanceof Error
				? e.message
				: "Failed to add new follow";
	}
	return error
}
