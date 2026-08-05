import { apiFetch, query } from "./client";
import type {
	AlbumInfo,
	ArtistInfo,
	ArtistPaginatedAlbums,
	PaginationQuery,
	PaginatedAlbums,
	PaginatedArtists,
	PaginatedPlaylists,
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

export function searchSetup(q: string, params: PaginationQuery = {}): Promise<SearchResult> {
	return apiFetch<SearchResult>(`/search/setup${query({ q, ...params })}`);
}

export function searchSong(
	q: string,
	provider: string,
	params: PaginationQuery = {},
): Promise<PaginatedSongs> {
	return apiFetch<PaginatedSongs>(
		`/search/song${query({ q, provider, ...params })}`,
	);
}

export function searchAlbum(
	q: string,
	provider: string,
	params: PaginationQuery = {},
): Promise<PaginatedAlbums> {
	return apiFetch<PaginatedAlbums>(
		`/search/album${query({ q, provider, ...params })}`,
	);
}

export function searchArtist(
	q: string,
	provider: string,
	params: PaginationQuery = {},
): Promise<PaginatedArtists> {
	return apiFetch<PaginatedArtists>(
		`/search/artist${query({ q, provider, ...params })}`,
	);
}

export function searchPlaylist(
	q: string,
	provider: string,
	params: PaginationQuery = {},
): Promise<PaginatedPlaylists> {
	return apiFetch<PaginatedPlaylists>(
		`/search/playlist${query({ q, provider, ...params })}`,
	);
}
