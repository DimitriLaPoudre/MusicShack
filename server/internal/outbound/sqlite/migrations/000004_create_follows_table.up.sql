CREATE TABLE follows (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	provider TEXT NOT NULL,
	artist_id TEXT NOT NULL,
	artist_name TEXT NOT NULL,
	artist_picture_url TEXT NOT NULL,
	featuring BOOLEAN NOT NULL,
	UNIQUE (user_id, provider, artist_id)
);
