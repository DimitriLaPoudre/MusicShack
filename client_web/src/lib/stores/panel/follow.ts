import { writable } from "svelte/store";
import type { FollowsResponse } from "$lib/types/response";

export const followList = writable<null | FollowsResponse>(null);
