## Quick Start

You can run MusicShack in two ways:

---

### 🚀 Option 1: Run with Docker (production, no build required)

**No need to clone the repo or install anything except Docker.**

Pull and run the official images from GitHub Container Registry (GHCR):

```sh
# Backend (Go API)


# Frontend (SvelteKit)
**Prerequisites:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
```

Configure your environment variables as needed (see documentation for details).

---

### 🛠️ Option 2: Run from the repository (development mode)

**1. Clone the repository**
```sh
git clone https://github.com/DimitriLaPoudre/MusicShack.git
cd MusicShack
```

**2. Configure environment variables**
Copy the example file below to .env at the repository root and adapt the values. Docker Compose and the Go/SvelteKit services will load these variables. Do not commit secrets.

```env
# .env.example - common variables
# Go backend

# Database (Postgres)
POSTGRES_USER=musicshack
POSTGRES_PASSWORD=changeme
POSTGRES_DB=musicshack_db

# Auth / sessions
JWT_SECRET=changeme_jwt_secret
SESSION_SECRET=changeme_session_secret
SESSION_MAX_AGE=86400

# Music file storage from the go binary position
DOWNLOAD_FOLDER=/users for docker by example and ../users for repository run

# Docker Compose helper (optional)
```

Explanations (where and why)
- PORT: backend Go listen port. Docker Compose maps this to the host.
- NODE_ENV / LOG_LEVEL: control backend mode and logging verbosity.
- DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME / DATABASE_URL: Postgres connection used by the backend. docker-compose.yml typically reads POSTGRES_* to initialize the Postgres container; the app uses DB_* or DATABASE_URL to connect.
- JWT_SECRET / SESSION_SECRET / SESSION_MAX_AGE: secrets for signing JWTs and session cookies. Use strong values in production.
- STORAGE_TYPE / STORAGE_PATH: choose local or s3 for storing music files. STORAGE_PATH should match the volume mount configured in docker-compose.
- S3_*: credentials and bucket info when using S3. Never publish these keys.
- HIFI_API_URL / HIFI_API_KEY: configuration for third‑party music API instances referenced by the UI.
- VITE_API_BASE: public API base the frontend uses to call the backend. Browser-exposed variables must start with VITE_.
- POSTGRES_*: used by the Postgres service in docker-compose to create the DB and user; keep these consistent with your DB_* values if you want the container initialized accordingly.

Quick tips
- Run: cp .env.example .env && edit .env.
- To run with Docker Compose using the file: docker-compose --env-file .env up --build
- For frontend dev: Vite reads VITE_* vars; restart the dev server after changes.
- Never commit secrets to Git. Use a secret manager for production values.
- Verify volume paths in docker-compose match STORAGE_PATH for local storage.
- If you provide DATABASE_URL, prefer it over separate DB_* vars to avoid mismatch.

**3. Start the frontend**
Download npm if not already install.
```sh
cd client_web
npm install
npm run dev
```

**4. Start the backend**
Download go if not already install.
```sh
cd server
go mod download
go run main.go
```

```sh
docker-compose up --build
```

The frontend will be available on the configured port by the .env variable PORT.

---

## Usage

1. Open the web interface in your browser.
2. Sign up or log in.
3. Add a music API instance (e.g. hifi) in the settings.
4. Browse, search, and download music directly to your server.
5. Manage your library and follow artists.

## Contributing

Pull requests are welcome! Please open an issue first to discuss major changes.

## License

MIT License — see [LICENSE](LICENSE)
