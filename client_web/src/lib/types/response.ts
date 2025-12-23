// Types de réponses reçues de l'API

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
}

export interface UserResponse {
  user: User;
}

export type UsersResponse = User[];

export interface ApiInstanceItem {
  id: number;
  api: string;
  url: string;
}
export type ApiInstancesResponse = ApiInstanceItem[];

export interface FollowItem {
  id: number;
  api: string;
  artist: ArtistData;
}
export type FollowsResponse = FollowItem[];

export interface DownloadData {
  id: number;
  data: SongData;
  api: string;
  status: "pending" | "running" | "done" | "failed" | "cancel";
}
export interface DownloadListResponse {
  tasks: DownloadData[];
}

export interface SongData {
  id: string;
  title: string;
  duration: number;
  releaseDate: string;
  trackNumber: number;
  volumeNumber: number;
  maximalAudioQuality: string;
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
  id: string;
  title: string;
  duration: number;
  releaseDate: string;
  numberTracks: number;
  numberVolumes: number;
  coverUrl: string;
  maximalAudioQuality: string;
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
  maximalAudioQuality: string;
  artists: SongDataArtist[];
}

export interface ArtistData {
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
  maximalAudioQuality: string;
  popularity: number;
  artists: SongDataArtist[];
  album: SongDataAlbum;
}
export interface SearchDataAlbum {
  id: string;
  title: string;
  duration: number;
  coverUrl: string;
  maximalAudioQuality: string;
  popularity: number;
  artists: AlbumDataArtist[];
}
export interface SearchDataArtist {
  id: string;
  name: string;
  pictureUrl: string;
  popularity: number;
}
