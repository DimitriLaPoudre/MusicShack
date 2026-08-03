import { apiFetch, query } from "./client";
import type {
	AlbumInfo,
	ArtistInfo,
	ArtistPaginatedAlbums,
	PaginationQuery,
	PaginatedSongs,
	PlaylistInfo,
	SearchResult,
	SongInfo,
} from "$lib/types";

export function getSongInfo(provider: string, id: string): Promise<SongInfo> {
	return apiFetch<SongInfo>(`/song/${encodeURIComponent(provider)}/${encodeURIComponent(id)}`);
}

export function getAlbumInfo(provider: string, id: string): Promise<AlbumInfo> {
	return apiFetch<AlbumInfo>(`/album/${encodeURIComponent(provider)}/${encodeURIComponent(id)}`);
}

export function getAlbumSongs(
	provider: string,
	id: string,
	params: PaginationQuery = {},
): Promise<PaginatedSongs> {
	return apiFetch<PaginatedSongs>(
		`/album/${encodeURIComponent(provider)}/${encodeURIComponent(id)}/songs${query(params)}`,
	);
}

export function getArtistInfo(provider: string, id: string): Promise<ArtistInfo> {
	return apiFetch<ArtistInfo>(`/artist/${encodeURIComponent(provider)}/${encodeURIComponent(id)}`);
}

export function getArtistAlbums(
	provider: string,
	id: string,
	params: PaginationQuery = {},
): Promise<ArtistPaginatedAlbums> {
	return apiFetch<ArtistPaginatedAlbums>(
		`/artist/${encodeURIComponent(provider)}/${encodeURIComponent(id)}/albums${query(params)}`,
	);
}

export function getPlaylistInfo(provider: string, id: string): Promise<PlaylistInfo> {
	return apiFetch<PlaylistInfo>(
		`/playlist/${encodeURIComponent(provider)}/${encodeURIComponent(id)}`,
	);
}

export function getPlaylistSongs(
	provider: string,
	id: string,
	params: PaginationQuery = {},
): Promise<PaginatedSongs> {
	return apiFetch<PaginatedSongs>(
		`/playlist/${encodeURIComponent(provider)}/${encodeURIComponent(id)}/songs${query(params)}`,
	);
}

export function search(q: string, params: PaginationQuery = {}): Promise<SearchResult> {
	return apiFetch<SearchResult>(`/search${query({ q, ...params })}`);
}
