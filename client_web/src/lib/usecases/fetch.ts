import { goto } from "$app/navigation";
import type { ErrorResponse } from "$lib/types/response/basic";

type ApiFetchOptions = {
	method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
	body?: any;
	redirect?: string;
};

export async function apiFetch<T>(
	path: string,
	{ method = "GET", body, redirect = "/login" }: ApiFetchOptions = {},
): Promise<T> {
	const options: RequestInit = {
		method: method,
		credentials: "include",
		headers: { "Content-Type": "application/json" },
	};

	if (body) {
		options.body = JSON.stringify(body);
	}

	const res = await fetch("/api" + path, options);
	if (res.status === 401) {
		await goto(redirect);
		throw new Error(res.statusText);
	}

	const data = await res.json().catch(() => {
		throw new Error(`Invalid JSON response: ${res.statusText}`);
	});

	if (!res.ok) {
		const error = data as ErrorResponse;
		throw new Error(error.error || `Request failed with status ${res.status}`);
	}
	return data as T;
}

export async function apiFetchFormData<T>(
	path: string,
	fd: FormData,
	method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE" = "POST",
): Promise<T> {
	const options: RequestInit = {
		method: method,
		credentials: "include",
	};

	if (fd) {
		options.body = fd;
	}

	const res = await fetch("/api" + path, options);
	if (res.status === 401) {
		await goto("/login");
		throw new Error(res.statusText);
	}

	const data = await res.json().catch(() => {
		throw new Error(`Invalid JSON response: ${res.status}`);
	});

	if (!res.ok) {
		const error = data as ErrorResponse;
		throw new Error(error.error || `Request failed with status ${res.status}`);
	}
	return data as T;
}
