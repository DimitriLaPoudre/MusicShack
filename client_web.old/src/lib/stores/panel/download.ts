import { writable } from "svelte/store";
import type { DownloadListResponse } from "$lib/types/response";

export const downloadList = writable<null | DownloadListResponse>(null);
