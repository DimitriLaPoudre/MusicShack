CREATE TABLE instances (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	provider TEXT NOT NULL,
	plugin TEXT NOT NULL,
	url TEXT NOT NULL,
	ping INTERVAL,
	UNIQUE (user_id, url)
);
