// ---- Requests ---- //

export interface LoginForm {
	username: string;
	password: string;
	remember: boolean;
}

export interface CreateUser {
	username: string;
	password: string;
	hi_res: boolean;
	role?: string;
}

export interface UpdateUser {
	username?: string;
	password?: string;
	hi_res?: boolean;
}

export interface CreateInstance {
	url: string;
}

export interface CreateFollow {
	provider: string;
	artist_id: string;
	featuring: boolean;
}

export interface CreateDownload {
	provider: string;
	id: string;
}

export interface PaginationQuery {
	offset?: number;
	limit?: number;
}

export interface SearchQuery extends PaginationQuery {
	q: string;
	provider: string;
}

// ---- Responses ---- //

export interface StatusResponse {
	status: string;
}

export interface ErrorResponse {
	error: string;
}

export interface User {
	id: string;
	username: string;
	hi_res: boolean;
	role: string;
}

export interface Instance {
	id: string;
	user_id: string;
	provider: string;
	plugin: string;
	url: string;
	ping: number | null;
}

export interface Follow {
	id: string;
	user_id: string;
	provider: string;
	artist_id: string;
	artist_name: string;
	artist_picture_url: string;
	featuring: boolean;
}

export interface DownloadTask {
	id: string;
	user_id: string;
	cover_url: string;
	title: string;
	artists: string[];
	album: string;
	status: string;
	status_comment: string;
}

export interface AudioQuality {
	name: string;
	color: string;
}

export interface Pagination {
	limit: number;
	offset: number;
	total_number_of_items: number;
}

export interface MiniArtist {
	id: string;
	name: string;
}

export interface SongAlbum {
	id: string;
	title: string;
	cover_url: string;
}

export interface SongInfo {
	provider: string;
	id: string;
	title: string;
	duration: number;
	replay_gain: number;
	peak: number;
	album_replay_gain: number;
	album_peak: number;
	release_date: string;
	track_number: number;
	volume_number: number;
	audio_quality: AudioQuality;
	explicit: boolean;
	popularity: number;
	isrc: string;
	artists: MiniArtist[];
	album: SongAlbum;
}

export interface PaginatedSongs extends Pagination {
	provider: string;
	items: SongInfo[];
}

export interface AlbumInfo {
	provider: string;
	id: string;
	title: string;
	duration: number;
	release_date: string;
	number_tracks: number;
	number_volumes: number;
	cover_url: string;
	audio_quality: AudioQuality;
	explicit: boolean;
	artists: MiniArtist[];
}

export interface PaginatedAlbums extends Pagination {
	provider: string;
	items: AlbumInfo[];
}

export interface ArtistInfo {
	provider: string;
	followed: string | null;
	id: string;
	name: string;
	picture_url: string;
}

export interface PaginatedArtists extends Pagination {
	provider: string;
	items: ArtistInfo[];
}

export interface ArtistPaginatedAlbums {
	provider: string;
	albums: PaginatedAlbums;
	ep: PaginatedAlbums;
	singles: PaginatedAlbums;
}

export interface PlaylistInfo {
	provider: string;
	id: string;
	title: string;
	description: string;
	duration: number;
	last_updated: string;
	number_of_tracks: number;
	cover_url: string;
}

export interface PaginatedPlaylists extends Pagination {
	provider: string;
	items: PlaylistInfo[];
}

export interface SearchItem {
	type: string;
	data: unknown;
}

export interface SearchProviderResult {
	songs: PaginatedSongs;
	albums: PaginatedAlbums;
	artists: PaginatedArtists;
	playlists: PaginatedPlaylists;
}

export interface SearchResult {
	item?: SearchItem;
	provider_result?: Record<string, SearchProviderResult>;
}
