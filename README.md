# MusicShack

MusicShack is a self-hosted music management platform that lets you browse, download and organize music on your own server. It pairs a modern Svelte frontend with a lightweight Go backend and supports extensible plugins for additional music sources.

---

## Purpose

Provide a private, easy-to-run solution to manage your personal music library. MusicShack focuses on privacy, extensibility (plugins) and a simple deployment flow via Docker.

---

## Key features

- Browse and search music catalogs through plugins (example: hifi)
- Download tracks and albums directly to the server
- Follow artists and manage your library
- User authentication and simple user management
- Plugin architecture to add new data sources
- Deployable with Docker / Docker Compose

---

## Recommended deployment: Docker Compose

Use Docker Compose to run MusicShack quickly and reliably. The repository includes a `docker-compose.yml` sample. For production, store secrets in a `.env` file (not in version control).

**Use the Docker image from Github Container Repository (ghcr.io/dimitrilapoudre/musicshack:latest) for better and easier update handling.**

1. Clone the repository:

```bash
git clone https://github.com/DimitriLaPoudre/MusicShack.git
cd MusicShack
```

2. Create a `.env` file at the repository root (recommended):

```env
JWT_SECRET=change_me_super_secret
POSTGRES_USER=musicshack
POSTGRES_PASSWORD=strong_password
POSTGRES_DB=musicshack
POSTGRES_HOST=database
```

3. Start the stack:

```bash
docker compose up --build -d
```

The main service listens on port `8080` by default. Downloaded files are persisted under the `./downloads` volume.

> Tip: keep your `.env` file out of source control and never commit secrets.

---

### Example `docker-compose.yml` (excerpt)

The repository already contains a `docker-compose.yml`. The example below demonstrates the primary services and environment variables used by MusicShack:

```yaml
services:
  musicshack:
    image: ghcr.io/dimitrilapoudre/musicshack:latest
    depends_on:
      database:
        condition: service_healthy
    restart: unless-stopped
    user: "1000:1000"
    ports:
      - "8080:8080"
    volumes:
      - ./downloads:/downloads
    environment:
      URL: http://localhost
      PORT: 8080
      DOWNLOAD: /downloads
      JWT_SECRET: ${JWT_SECRET}
      POSTGRES_HOST: ${POSTGRES_HOST}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}

  database:
    image: postgres:16
    ports:
      - "5432:5432"
    volumes:
      - ./db:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 5

networks:
  default:
    name: musicshack
    driver: bridge
```

---

## Local development

If you prefer to run the components locally for development:

### Frontend

```bash
cd client_web
npm install
npm run dev
```

### Backend

```bash
cd server
go mod download
go run main.go
```

The backend uses environment variables (see `server/internal/config/config.go`): `URL`, `JWT_SECRET`, `DOWNLOAD_FOLDER`, and Postgres connection variables (see `server/internal/db/database.go`) like the Docker Image.

---

## Architecture & plugins

The backend loads plugins from `server/internal/plugins`. Each plugin implements the `Plugin` interface (see `server/internal/models/plugin.go`) and lets you add new sources (e.g. `hifi`, `hifiv2`) with plugins.Register().

---

## Contributing

- Open an issue to discuss features or bugs

---

## License

MIT — see the `LICENSE` file.
