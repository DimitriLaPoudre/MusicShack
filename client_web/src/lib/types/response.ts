export interface StatusResponse {
	status: "ok";
}

export interface ErrorResponse {
	error: string;
}

export interface User {
	id: number;
	username: string;
	password: string;
	bestQuality: boolean;
}

export type UserResponse = User;

export type UsersResponse = User[];

export interface InstanceItem {
	id: number;
	api: string;
	url: string;
}
export type InstancesResponse = InstanceItem[];

export interface FollowItem {
	id: number;
	api: string;
	artistId: string;
	artistName: string;
	artistPictureUrl: string;
}
export type FollowsResponse = FollowItem[];

export interface DownloadData {
	id: number;
	data: SongData;
	api: string;
	status: "pending" | "running" | "done" | "failed" | "cancel";
}

export type DownloadListResponse = DownloadData[];

export interface SongData {
	api: string;
	id: string;
	title: string;
	duration: number;
	replayGain: number;
	peak: number;
	releaseDate: string;
	trackNumber: number;
	volumeNumber: number;
	audioQuality: number;
	popularity: number;
	isrc: string;
	coverUrl: string;
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
	api: string;
	id: string;
	title: string;
	duration: number;
	releaseDate: string;
	numberTracks: number;
	numberVolumes: number;
	coverUrl: string;
	audioQuality: number;
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
	audioQuality: number;
	artists: SongDataArtist[];
}

export interface ArtistData {
	api: string;
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
	audioQuality: number;
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
	audioQuality: number;
	popularity: number;
	artists: SongDataArtist[];
	album: SongDataAlbum;
}
export interface SearchDataAlbum {
	id: string;
	title: string;
	duration: number;
	coverUrl: string;
	audioQuality: number;
	popularity: number;
	artists: AlbumDataArtist[];
}
export interface SearchDataArtist {
	id: string;
	name: string;
	pictureUrl: string;
	popularity: number;
}

export type SearchResponse = { [key: string]: SearchData }
