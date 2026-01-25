# DevContainer for OPTel Training

This DevContainer provides a consistent development environment for the OPTel Training API.

## Prerequisites

- VS Code or Cursor with the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed
- Docker Desktop (or Docker Engine) running

## Getting Started

1. **Open the repository** in VS Code/Cursor
2. When prompted, click **"Reopen in Container"** (or use Command Palette: `Dev Containers: Reopen in Container`)
3. Wait for the container to build and start (first time may take a few minutes)
4. Once the container is ready, you'll have:
   - Go 1.25+ compiler
   - oapi-codegen (OpenAPI code generator)
   - golangci-lint (Go linter)
   - All tools available on PATH

## Development Workflow

Inside the DevContainer, you can run all development commands natively:

```bash
cd api
make generate  # Generate code from OpenAPI spec
make lint      # Run linter
make test      # Run tests
make run       # Start the API server
```

The API server will be available at `http://localhost:8080` (port forwarding is configured automatically).

## Environment Variables

No special environment variables are required for basic development. The DevContainer sets up:
- `PATH` includes Go bin directory (`$(go env GOPATH)/bin`)
- Go workspace is configured at `/workspace`

## GitHub CLI (`gh`)

This DevContainer includes GitHub CLI (`gh`).

### Recommended: pass a host token into the container

If you want `gh` and `git push` to work non-interactively (including from automation), set `GH_TOKEN` on your host **before** starting Cursor/VS Code, then rebuild the container.

- **Host**: set `GH_TOKEN` (a GitHub Personal Access Token with appropriate repo scopes)
- **DevContainer**: `GH_TOKEN` is injected via `${localEnv:GH_TOKEN}` (see `.devcontainer/devcontainer.json`)

After rebuild, inside the container:

```bash
echo "GH_TOKEN_len=${#GH_TOKEN}"
gh api user --jq '.login'
```

### Alternative: interactive login (manual)

You can also authenticate inside the container:

```bash
gh auth login
```

## Troubleshooting

- **Container fails to start**: Check Docker is running and has enough resources allocated
- **Tools not found**: Run `source ~/.bashrc` or restart the terminal after container creation
- **Port conflicts**: Change the port forwarding in `devcontainer.json` if port 8080 is already in use
