CREATE TABLE admin (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	password TEXT NOT NULL,
	token TEXT,
	expires_at TIMESTAMPTZ
);
