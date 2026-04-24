import type { StatusResponse } from "$lib/types/response/basic";
import { toast } from "svelte-sonner";
import { apiFetch } from "./fetch";

export async function LoginAdmin(password: string): Promise<void> {
	try {
		await apiFetch<StatusResponse>("/admin/login", "POST", {
			password: password,
		});
	} catch (e) {
		toast.error(e instanceof Error ? e.message : "Login unknown error");
	}
}
