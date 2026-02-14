# Deployment (GitHub Actions)

## Stack

- **Backend**: Go app (see `Dockerfile`).
- **Nginx**: Reverse proxy in front of the backend (see `nginx/default.conf`).
- **Orchestration**: Docker Compose (`docker-compose.yml`) runs both services.

## Workflows

| Branch | Workflow       | Env file                     | Compose project    | Nginx on host |
|--------|----------------|------------------------------|--------------------|---------------|
| `dev`  | Deploy Dev     | `/stitchfolio/env/.dev.env`  | `stitchfolio-dev`  | port **9001** |
| `main` | Deploy Main    | `/stitchfolio/env/.prod.env` | `stitchfolio-prod` | port **80**   |

- **Trigger**: Push or merge to `dev` or `main` (or run manually via “Run workflow”).
- **Steps**: Checkout → SSH → pull branch → write `.env` for compose → build → run migrations → `docker compose up -d` (backend + nginx).

## GitHub secrets

- **`SSH_PRIVATE_KEY`** (required): Private key that can SSH as `root@72.61.254.64`. Add in **Settings → Secrets and variables → Actions** (or at org level).

## Server setup

On `72.61.254.64`:

1. **Env files** (already in place):
   - `/stitchfolio/env/.dev.env`
   - `/stitchfolio/env/.prod.env`

2. **Clone the repo in two directories** (dev and prod run as separate compose stacks):

   ```bash
   mkdir -p /stitchfolio
   git clone <repo-url> /stitchfolio/backend-dev
   git clone <repo-url> /stitchfolio/backend-prod
   cd /stitchfolio/backend-dev && git checkout dev
   cd /stitchfolio/backend-prod && git checkout main
   ```

3. **Docker**: Docker (and Docker Compose v2) must be installed. The user used by the workflow (e.g. `root`) must be able to run `docker` and `docker compose` without sudo.

4. **Ports**: Dev is exposed on host port **9001** (nginx). Prod is on **80** (nginx). The backend is not exposed on the host; only nginx is.

## Repo layout

- `Dockerfile` – backend image.
- `docker-compose.yml` – backend + nginx; reads `ENV_FILE_PATH`, `CONFIG_FILE`, `NGINX_HOST_PORT` from `.env` (the workflow writes `.env` on deploy).
- `nginx/default.conf` – nginx reverse proxy to `backend:9000`.
- `.env.example` – example for local or reference; real `.env` is written by the workflow on the server (do not commit `.env`).
