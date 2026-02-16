<h1 align="center">
   <img src="https://raw.githubusercontent.com/DimitriLaPoudre/MusicShack/main/client_web/static/assets/apple-touch-icon.png" align="center">
   
   MusicShack
</h1>

<h3 align="center">
   Self-Hosted Web Application for Music Library Import and Management
</h3>

<p align="center">
   <a href="https://github.com/DimitriLaPoudre/MusicShack/commits/main"><img src="https://img.shields.io/github/last-commit/DimitriLaPoudre/MusicShack?style=for-the-badge&color=green" alt="Last Commit"></a>
   <a href="https://github.com/DimitriLaPoudre/MusicShack/blob/main/LICENSE"><img src="https://img.shields.io/github/license/DimitriLaPoudre/MusicShack?style=for-the-badge&color=5D6D7E" alt="License"></a>
   <a href="https://github.com/DimitriLaPoudre/MusicShack/stargazers"><img src="https://img.shields.io/github/stars/DimitriLaPoudre/MusicShack?logo=github&style=for-the-badge&color=E67E22" alt="GitHubStars"></a>
   <a href="https://github.com/DimitriLaPoudre/MusicShack/pkgs/container/musicshack/"><img src="https://img.shields.io/badge/docker-package-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker Package"></a>
   <br>
   <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
   <img src="https://img.shields.io/badge/SvelteKit-FF3E00?style=for-the-badge&logo=svelte&logoColor=white" alt="SvelteKit">
   <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
   <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
</p>

---

## About

MusicShack is a self-hosted web application to import, format, and manage a music library.
It centralizes browsing, downloading, and metadata handling via multiple sources, with automatic downloads of your favorite artists’ new releases.
The server stores downloads locally, while the web UI provides search, download, follow, library editing, and admin features for multiple users.

---

> [!IMPORTANT]
>
> This is my first **useful** solo project using technologies that are completely new to me.
> I refactor many things every time I notice past architecture error.  
> Feel free to correct me or ask for **any** features you think are relevant.
>
> Song tagging is based on the Navidrome standard (I don't know if it's compatible with other music servers; I hope so^^).

---

## Features

- Browse and search music catalogs through different sources
- Download tracks and albums directly from these sources
- Follow artists — every day at 1AM, MusicShack will download any new songs released by the artists you follow
- Add new source URL
- User authentication and simple user management
- Admin panel for adding new users
- Deployable with Docker / Docker Compose
- Plugin architecture to add new data sources in the future, like [DAB](https://dab.yeet.su/) or even your friends' servers

---

## Screenshots

TO-DO

---

## Deployment (Docker)

MusicShack provides an `example.env` file for quick deployment  
For more info about [environment variable](https://github.com/DimitriLaPoudre/MusicShack/README.md#environment-variable)

### Deployment steps

1. Copy the .env file:
   ```bash
   cp example.env .env
   ```
2. Launch MusicShack:
   ```bash
   docker compose up -d
   ```
3. Access the admin panel:
   - Go to `http://localhost:8080/admin`
   - Use your password set at `ADMIN_PASSWORD` in the .env file
4. Create a user via the admin panel:
   - Enter an username and a password
   - Click on the `+` button or press `Enter` key
5. Access the main interface:
   - Go to `http://localhost:8080/`
   - Log in as the new user you created

   Welcome to **MusicShack**^^

---

## Usage

- Add sources:
  - Click on the `Settings` button
  - Enter an instance URL (find some [here](https://github.com/EduardPrigoana/hifi-instances))
  - Click on the `+` button or press `Enter` key
- Follow an artist:
  - Click on the `Search` button
  - Select the artist name
  - Go to the `Artist` section
  - Click on their profile
  - Click on the `Follow` button
- Download a song:
  - Click on the `Search` button
  - Select the song name
  - Go to the `Song` section
  - Click on the `Download` button under the song
- Upload a song:
  - Click on the `Library` button
  - Click on the `Upload` button
  - Select a file
  - Choose a cover, title, album name, etc.
  - Click on the `Save` button
- Edit a song:
  - Click on the `Library` button
  - Click on the `Edit` button of the song you want to edit
  - Choose a new cover, title, album name, etc.
  - Click on the `Save` button

---

## Environment Variable

- `HTTPS` = _boolean_ (**false** by default) set at true if your domain use https
- `PORT` = _number_ (**8080** by default) port where the app will be accessible
- `LIBRARY_PATH` = _string_ (mandatory) path to the library (downloads/uploads will go into that directory)
- `ADMIN_PASSWORD` = _string_ (mandatory) default password for admin panel
- `POSTGRES_HOST` = _string_ (mandatory) localhost or name of the service that contains PostgreSQL
- `POSTGRES_USER` = _string_ (mandatory)
- `POSTGRES_PASSWORD` = _string_ (mandatory)
- `POSTGRES_DB` = _string_ (mandatory)

---

## Roadmap

The project roadmap is managed via [GitHub Projects](https://github.com/users/DimitriLaPoudre/projects/5)

---

## Contributing

- Open an issue to discuss features or bugs
- Feel free to add a plugin for a new provider and open a PR

---

## Credits

- **Design** inspired by the black theme of [monochrome](https://github.com/monochrome-music/monochrome)
- The **whole idea** comes from the existence of the [hifi](https://github.com/binimum/hifi-api) API

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
> - This tool serves as an interface for personal, non-commercial use.
>
> MusicShack assumes no responsibility for any misuse or legal violations.
