import { writable } from "svelte/store";
import type { InstancesResponse } from "$lib/types/response";

export const instanceList = writable<null | InstancesResponse>(null);
export const username = writable<string>("");
