export interface StatusResponse {
	status: "ok";
}

export interface ErrorResponse {
	error: string;
}

export interface User {
	username: string;
	hiRes: boolean;
}

export type UserResponse = User;

export type UsersResponse = User[];

export interface InstanceItem {
	id: number;
	api: string;
	provider: string;
	url: string;
	ping: number;
}
export type InstancesResponse = InstanceItem[];

export interface FollowItem {
	id: number;
	provider: string;
	artistId: string;
	artistName: string;
	artistPictureUrl: string;
}
export type Follow = FollowItem;

export type FollowsResponse = FollowItem[];

export interface DownloadData {
	id: number;
	data: SongData;
	provider: string;
	status: "pending" | "running" | "done" | "failed" | "cancel";
}

export type DownloadListResponse = DownloadData[];

export interface Quality {
	name: string;
	color: string;
}

export interface SongData {
	provider: string;
	api: string;
	id: string;
	title: string;
	duration: number;
	replayGain: number;
	peak: number;
	releaseDate: string;
	trackNumber: number;
	volumeNumber: number;
	audioQuality: Quality;
	popularity: number;
	isrc: string;
	explicit: boolean;
	artists: SongDataArtist[];
	album: SongDataAlbum;
}
export interface SongDataArtist {
	id: string;
	name: string;
}
export interface SongDataAlbum {
	id: string;
	title: string;
	coverUrl: string;
}

export interface AlbumData {
	provider: string;
	api: string;
	id: string;
	title: string;
	duration: number;
	releaseDate: string;
	numberTracks: number;
	numberVolumes: number;
	coverUrl: string;
	audioQuality: Quality;
	explicit: boolean;
	artists: AlbumDataArtist[];
	songs: AlbumDataSong[];
}
export interface AlbumDataArtist {
	id: string;
	name: string;
}
export interface AlbumDataSong {
	id: string;
	title: string;
	duration: number;
	trackNumber: number;
	volumeNumber: number;
	audioQuality: Quality;
	explicit: boolean;
	artists: SongDataArtist[];
}

export interface ArtistData {
	provider: string;
	api: string;
	followed: number;
	id: string;
	name: string;
	pictureUrl: string;
	albums: ArtistDataAlbum[];
	ep: ArtistDataAlbum[];
	singles: ArtistDataAlbum[];
}
export interface ArtistDataAlbum {
	id: string;
	title: string;
	duration: number;
	releaseDate: string;
	coverUrl: string;
	audioQuality: Quality;
	explicit: boolean;
	artists: AlbumDataArtist[];
}

export interface SearchData {
	songs: SearchDataSong[];
	albums: SearchDataAlbum[];
	artists: SearchDataArtist[];
}
export interface SearchDataSong {
	id: string;
	title: string;
	duration: number;
	audioQuality: Quality;
	popularity: number;
	explicit: boolean;
	artists: SongDataArtist[];
	album: SongDataAlbum;
}
export interface SearchDataAlbum {
	id: string;
	title: string;
	duration: number;
	coverUrl: string;
	audioQuality: Quality;
	explicit: boolean;
	popularity: number;
	artists: AlbumDataArtist[];
}
export interface SearchDataArtist {
	followed: number;
	id: string;
	name: string;
	pictureUrl: string;
	popularity: number;
}

export interface SearchResult {
	[key: string]: SearchData;
}

export interface UrlItem {
	provider: string;
	type: "artist" | "album" | "song";
	id: string;
}

export type SearchResponse =
	| {
		result: SearchResult;
	}
	| {
		url: UrlItem;
	};

export interface ResponseSong {
	id: number;
	title: string;
	album: string;
	albumArtists: string[];
	artists: string[];
	releaseDate: string;
	trackNumber: number;
	volumeNumber: number;
	explicit: boolean;
	isrc: string;
	albumGain: number;
	albumPeak: number;
	trackGain: number;
	trackPeak: number;
}

export interface ResponseLibrary {
	total: number;
	count: number;
	limit: number;
	offset: number;
	items: ResponseSong[];
}
