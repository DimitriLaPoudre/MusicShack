import type {
	CreateUserRequest,
	UpdateUserRequest,
} from "$lib/types/request/user";
import type { UserResponse } from "$lib/types/response/user";
import { apiFetch } from "./fetch";
import type { StatusResponse } from "$lib/types/response/basic";

export async function GetUsers(): Promise<UserResponse[]> {
	const users = await apiFetch<UserResponse[]>("/users", {
		method: "GET",
		redirect: "/admin/login",
	});
	return users;
}

export async function CreateUser(
	req: CreateUserRequest,
): Promise<UserResponse> {
	const user = await apiFetch<UserResponse>("/users", {
		method: "POST",
		body: req,
		redirect: "/admin/login",
	});
	return user;
}

export async function UpdateUser(
	id: string,
	req: UpdateUserRequest,
): Promise<UserResponse> {
	const user = await apiFetch<UserResponse>(`/users/${id}`, {
		method: "PUT",
		body: req,
		redirect: "/admin/login",
	});
	return user;
}

export async function DeleteUser(id: string): Promise<void> {
	await apiFetch<StatusResponse>(`/users/${id}`, {
		method: "DELETE",
		redirect: "/admin/login",
	});
}
