import { writable } from "svelte/store";
import type { DownloadTask } from "$lib/types";

export const downloadList = writable<DownloadTask[] | null>(null);
