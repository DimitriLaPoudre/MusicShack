CREATE TABLE songs (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	isrc TEXT,
	path TEXT NOT NULL,
	updated_at TIMESTAMP NOT NULL
	UNIQUE (user_id, path)
);
