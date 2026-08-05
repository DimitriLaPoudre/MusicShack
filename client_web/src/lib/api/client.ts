import type { ErrorResponse } from "$lib/types";

export class ApiError extends Error {
	readonly status: number;
	readonly body?: ErrorResponse;

	constructor(status: number, body?: ErrorResponse) {
		super(body?.error ?? `Request failed with status ${status}`);
		this.name = "ApiError";
		this.status = status;
		this.body = body;
	}
}

export async function apiFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
	const response = await fetch(`/api${path}`, {
		credentials: "include",
		...init,
		headers: {
			"Content-Type": "application/json",
			...init.headers,
		},
	});

	if (response.status === 204) {
		return undefined as T;
	}

	if (!response.ok) {
		let body: ErrorResponse | undefined;
		try {
			body = await response.json();
		} catch {
			// body not parseable as JSON
		}
		if (response.status === 401) {
			window.location.href = "/login";
		}
		throw new ApiError(response.status, body);
	}

	return (await response.json()) as T;
}

export function query(params: object): string {
	const search = new URLSearchParams();
	for (const [key, value] of Object.entries(params)) {
		if (value !== undefined && value !== null) {
			search.set(key, String(value));
		}
	}
	const s = search.toString();
	return s ? `?${s}` : "";
}
