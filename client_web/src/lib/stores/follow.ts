import { writable } from "svelte/store";
import type { Follow } from "$lib/types";

export const followList = writable<Follow[] | null>(null);
