# MusicShack

MusicShack is a self-hosted music management platform that lets you browse and download music on your own server. It pairs a modern Svelte frontend with a lightweight Go backend and supports extensible plugins for additional music sources.

---

## Purpose

Provides a private, easy-to-run solution to add new content to your personal music library. MusicShack focuses on privacy, extensibility (plugins) and a simple deployment flow via Docker.

---

## Key features

- Browse and search music catalogs through plugins (example: [hifi-api](https://github.com/uimaxbai/hifi-api))
- Download tracks and albums directly to the server with unified formatting across all plugin sources
- Follow artists — every Friday at 1AM, MusicShack will download any new songs released by artists you follow
- User authentication and simple user management
- Plugin architecture to add new data sources in the future, like Qobuz-DL or even your friends' servers
- Deployable with Docker / Docker Compose

---

## Important Info

This project is my first "big" solo project.

Feel free to ask for any features you think are relevant.

Song tags are based on the Navidrome standard (I don't know about other music servers).

The next big step will be to add a new private source to fetch data from projects like DAB Music Player or other Navidrome servers.

---

## Deployment (Docker Compose)

MusicShack provides two example files for quick deployment:
- `docker-compose.yml.example`
- `.env.example`

### Deployment steps

1. Copy the example files:
   ```bash
   cp docker-compose.yml.example docker-compose.yml
   cp .env.example .env
   ```
2. Edit the `.env` file to match your environment:
   - Set the `URL` (e.g. http://localhost or https://mywebsite.com)
   - Set the `PORT` (e.g. 8080)
   - Adjust Postgres credentials for better security
3. Create the `downloads` folder at the project root (if it doesn't exist):
   ```bash
   mkdir downloads
   ```
4. Start the stack:
   ```bash
   docker compose up -d
   ```
5. Access the admin interface:
   - Go to `http://URL:PORT/admin` (default: http://localhost:8080/admin)
   - The default admin password is: `changemenow`
   - Change it immediately after your first login!
6. Create a user via the admin interface
7. Access the main interface:
   - Go to `http://URL:PORT/`
   - Log in with the user you created

> ⚠️ Never commit your `.env` file with plain secrets.

> Downloaded files will be stored in the `downloads` folder.

Enjoy 🎶

---

## Contributing

- Open an issue to discuss features or bugs

---

## License

MIT — see the `LICENSE` file.
