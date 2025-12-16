import { goto } from "$app/navigation";
import { PUBLIC_API_URL } from "$env/static/public";

export async function apiFetch(
	path: string,
	method: string = "GET",
	body?: any,
): Promise<Response> {
	const res = await fetch(`${PUBLIC_API_URL}/api` + path, {
		method: method,
		credentials: "include",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(body),
	});
	if (res.status === 401) {
		goto("/login");
		return res;
	}
	return res;
}
