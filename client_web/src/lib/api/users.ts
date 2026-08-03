import { apiFetch } from "./client";
import type { CreateUser, UpdateUser, User } from "$lib/types";

export function getMe(): Promise<User> {
	return apiFetch<User>("/me");
}

export function updateMe(form: UpdateUser): Promise<User> {
	return apiFetch<User>("/me", {
		method: "PUT",
		body: JSON.stringify(form),
	});
}

export function deleteMe(): Promise<void> {
	return apiFetch<void>("/me", {
		method: "DELETE",
	});
}

export function createUser(form: CreateUser): Promise<User> {
	return apiFetch<User>("/users", {
		method: "POST",
		body: JSON.stringify(form),
	});
}

export function listUsers(): Promise<User[]> {
	return apiFetch<User[]>("/users");
}

export function getUser(userId: string): Promise<User> {
	return apiFetch<User>(`/users/${encodeURIComponent(userId)}`);
}

export function updateUser(userId: string, form: UpdateUser): Promise<User> {
	return apiFetch<User>(`/users/${encodeURIComponent(userId)}`, {
		method: "PUT",
		body: JSON.stringify(form),
	});
}

export function deleteUser(userId: string): Promise<void> {
	return apiFetch<void>(`/users/${encodeURIComponent(userId)}`, {
		method: "DELETE",
	});
}
