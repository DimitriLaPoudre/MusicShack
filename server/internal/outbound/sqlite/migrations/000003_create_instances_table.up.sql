CREATE TABLE instances (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	provider TEXT NOT NULL,
	plugin TEXT NOT NULL,
	url TEXT NOT NULL,
	ping INTEGER,
	UNIQUE (user_id, url)
);
