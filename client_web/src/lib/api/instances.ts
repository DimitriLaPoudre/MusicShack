import { apiFetch } from "./client";
import type { CreateInstance, Instance } from "$lib/types";

export function listInstances(refresh = false): Promise<Instance[]> {
	return apiFetch<Instance[]>(refresh ? "/me/instances?refresh" : "/me/instances");
}

export function createInstance(form: CreateInstance): Promise<Instance> {
	return apiFetch<Instance>("/me/instances", {
		method: "POST",
		body: JSON.stringify(form),
	});
}

export function deleteInstance(id: string): Promise<void> {
	return apiFetch<void>(`/me/instances/${encodeURIComponent(id)}`, {
		method: "DELETE",
	});
}
