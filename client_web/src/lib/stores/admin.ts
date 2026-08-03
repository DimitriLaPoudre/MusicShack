import { writable } from "svelte/store";
import type { User } from "$lib/types";

export const isAdmin = writable(false);
export const usersList = writable<User[] | null>(null);
