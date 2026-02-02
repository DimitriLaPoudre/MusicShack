<h1 align="center">

**MusicShack**
</h1>

<div align="center">

**Self-Host Web App for Music Library Import and Management**

[![Go](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![TypeScript](https://img.shields.io/badge/Language-TypeScript-blue.svg)](https://www.typescriptlang.org/)
[![Svelte](https://img.shields.io/badge/Framework-Svelte-green.svg)](https://svelte.dev/)
[![Docker](https://img.shields.io/badge/Container-Docker-blue.svg)](https://www.docker.com/)


</div>

---

## Overview

MusicShack is a self-hosted web app to import, organize, and manage a music library. It centralizes browsing, downloading, and metadata handling through pluggable sources, with automated artist tracking and scheduled updates. The server stores downloads locally, while the web UI provides search, playback, and admin features for multiple users.  

---

> [!IMPORTANT]
>
> This is my first **usefull** solo project with completely new techno to me.
> Feel free to correct me or ask for **any** features you think are relevant.
>
> Song tagging are based on the Navidrome standard (I don't know about other music servers).
>
> The next big step will be to add a new private source to fetch data from projects like DAB Music Player or other Subsonic servers.

---

## Key features

- Browse and search music catalogs through plugins
- Download tracks and albums directly from all plugin sources
- Follow artists — everyday at 1AM, MusicShack will download any new songs released by artists you follow
- User authentication and simple user management
- Deployable with Docker / Docker Compose
- Add source Url for plugins 
- Plugin architecture to add new data sources in the future, like [DAB](https://dab.yeet.su/) or even your friends' servers

---

## Deployment (Docker Compose)

MusicShack provides two example files for quick deployment:
- `example.docker-compose.yml`
- `example.env`

### Deployment steps

1. Copy the example files:
   ```bash
   cp example.docker-compose.yml docker-compose.yml
   cp example.env .env
   ```
2. Edit the `.env` file to match your environment:
   - `HTTPS` = *boolean* (**false** by default) set at true if your domain use https
   - `PORT` = *number* (**8080** by default)
   - `LIBRARY_PATH` = *string* (mandatory) path to the library
   - `ADMIN_PASSWORD` = *string* (mandatory)
   - Adjust Postgres credentials for better security
3. Create the `LIBRARY_PATH` folder at the project root (if it doesn't exist):
   ```bash
   mkdir ./download
   ```
4. Launch MusicShack:
   ```bash
   docker compose up -d
   ```
5. Access the admin interface:
   - Go to `http://URL:PORT/admin`
   - Use your password set in `ADMIN_PASSWORD`
6. Create a user via the admin interface
7. Access the main interface:
   - Go to `http://URL:PORT/`
   - Log in with the user you created
   - Add sources in Settings -> Instances
   - Follow your favorites artists Search -> Artist section -> Artist page -> Follow

> [!NOTE]
>
> Downloaded files will be stored in the `LIBRARY_PATH` folder.

Enjoy 🎶

---

## Contributing

- Open an issue to discuss features or bugs
- Feel free to add a plugin for new provider and open a PR

---

## License

MIT — see the `LICENSE` file.

> [!IMPORTANT]
>
> **Disclaimer**
>
> This project is intended for educational and personal use only. The developers do not encourage or endorse piracy.
>
> - Users are solely responsible for complying with copyright laws in their jurisdiction.
> - All music rights remain with their respective copyright holders.
> - This tool serves as a interface for personal, non-commercial use.
>
> MusicShack assumes no responsibility for any misuse or legal violations.
