import { writable } from "svelte/store";
import type { InstancesResponse, UserResponse } from "$lib/types/response";

export const instanceList = writable<null | InstancesResponse>(null);
export const userData = writable<null | UserResponse>(null);
