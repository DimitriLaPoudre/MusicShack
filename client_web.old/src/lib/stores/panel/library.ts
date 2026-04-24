import type { ResponseLibrary } from "$lib/types/response";
import { writable } from "svelte/store";

export const libraryPage = writable<null | ResponseLibrary>(null);
