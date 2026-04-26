import type { StatusResponse } from "$lib/types/response/basic";
import type { LoginAdminRequest } from "$lib/types/request/admin";
import { apiFetch } from "./fetch";
import { goto } from "$app/navigation";

export async function LoginAdmin(password: string): Promise<void> {
	await apiFetch<StatusResponse>("/admin/login", {
		method: "POST",
		body: { password: password } as LoginAdminRequest,
		redirect: "/admin/login",
	});
	goto("/admin/dashboard");
}
