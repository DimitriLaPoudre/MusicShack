import { goto } from "$app/navigation";
import type { ErrorResponse } from "$lib/types/response";

export async function apiFetch<T>(
	path: string,
	method: string = "GET",
	body?: any,
): Promise<T | ErrorResponse> {
	const res = await fetch("/api" + path, {
		method: method,
		credentials: "include",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(body),
	});
	let data;
	if (res.ok) {
		data = await res.json() as T;
	} else {
		data = await res.json() as ErrorResponse;
	}
	if (res.status === 401) {
		await goto("/login");
	}
	return data;
}

export async function adminFetch<T>(
	path: string,
	method: string = "GET",
	body?: any,
): Promise<T | ErrorResponse> {
	const res = await fetch("/api" + path, {
		method: method,
		credentials: "include",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(body),
	});
	let data;
	if (res.ok) {
		data = await res.json() as T;
	} else {
		data = await res.json() as ErrorResponse;
	}
	if (res.status === 401) {
		await goto("/admin/login");
	}
	return data;
}
