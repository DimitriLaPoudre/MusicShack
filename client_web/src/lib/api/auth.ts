import { apiFetch } from "./client";
import type { LoginForm } from "$lib/types";

export function login(form: LoginForm): Promise<void> {
	return apiFetch<void>("/auth/login", {
		method: "POST",
		body: JSON.stringify(form),
	});
}

export function logout(): Promise<void> {
	return apiFetch<void>("/auth/logout", {
		method: "DELETE",
	});
}
