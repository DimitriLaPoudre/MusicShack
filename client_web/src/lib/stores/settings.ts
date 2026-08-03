import { writable } from "svelte/store";
import type { Instance, User } from "$lib/types";

export const userData = writable<User | null>(null);
export const instanceList = writable<Instance[] | null>(null);
