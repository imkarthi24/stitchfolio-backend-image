# Deployment (GitHub Actions)

## Stack

- **Backend**: Go app (see `Dockerfile`).
- **Nginx**: Reverse proxy in front of the backend (see `nginx/default.conf`).
- **Orchestration**: Docker Compose (`docker-compose.yml`) runs both services.

## Workflows

| Branch | Workflow       | Env file                     | Compose project    | Nginx on host |
|--------|----------------|------------------------------|--------------------|---------------|
| `dev`  | Deploy Dev     | `/root/stitchfolio_env/dev.env`  | `stitchfolio-dev`  | port **9001** |
| `main` | Deploy Main    | `/root/stitchfolio_env/prod.env` | `stitchfolio-prod` | port **80**   |

- **Trigger**: Push or merge to `dev` or `main` (or run manually via “Run workflow”).
- **Steps**: Checkout → SSH → pull branch → write `.env` for compose → build → run migrations → `docker compose up -d` (backend + nginx).

## GitHub secrets

- **`SSH_PRIVATE_KEY`** (required): Private key that can SSH as `root@31.97.202.6`. Add in **Settings → Secrets and variables → Actions** (or at org level).

## Server setup

On `31.97.202.6`:

1. **Env files** (already in place):
   - `/root/stitchfolio_env/dev.env`
   - `/root/stitchfolio_env/prod.env`

2. **Deploy directories**: The workflow creates `/stitchfolio` and clones the repo into `/stitchfolio/backend-dev` or `/stitchfolio/backend-prod` on first run if those paths don’t exist. For a **private** repo, the server must be able to clone from GitHub (e.g. add the server’s SSH key as a deploy key in the repo). You can still create the dirs and clone manually if you prefer.

3. **Docker**: Docker (and Docker Compose v2) must be installed. The user used by the workflow (e.g. `root`) must be able to run `docker` and `docker compose` without sudo.

4. **Ports**: Dev is exposed on host port **9001** (nginx). Prod is on **80** (nginx). The backend is not exposed on the host; only nginx is.

5. **Git (for auto-clone)**: If you rely on the workflow to create and clone the repo, ensure `git` is installed on the server. For private repos, configure GitHub clone access (e.g. deploy key: `ssh-keygen` on server, add `~/.ssh/id_rsa.pub` to the repo’s **Deploy keys**).

## Repo layout

- `Dockerfile` – backend image.
- `docker-compose.yml` – backend + nginx; reads `ENV_FILE_PATH`, `CONFIG_FILE`, `NGINX_HOST_PORT` from `.env` (the workflow writes `.env` on deploy).
- `nginx/default.conf` – nginx reverse proxy to `backend:9000`.
- `.env.example` – example for local or reference; real `.env` is written by the workflow on the server (do not commit `.env`).
