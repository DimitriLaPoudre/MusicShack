import { writable } from "svelte/store";

export type PanelId = "search" | "follow" | "download" | "settings" | "admin";

export const openPanel = writable<PanelId | null>(null);
