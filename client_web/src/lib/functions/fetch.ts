import { goto } from "$app/navigation";

export async function apiFetch(
	path: string,
	method: string = "GET",
	body?: any,
): Promise<Response> {
	const res = await fetch("/api" + path, {
		method: method,
		credentials: "include",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(body),
	});
	if (res.status === 401) {
		await goto("/login");
		return res;
	}
	return res;
}

export async function adminFetch(
	path: string,
	method: string = "GET",
	body?: any,
): Promise<Response> {
	const res = await fetch("/api/admin" + path, {
		method: method,
		credentials: "include",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(body),
	});
	if (res.status === 401) {
		await goto("/admin/login");
		return res;
	}
	return res;
}
