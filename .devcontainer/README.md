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

## Troubleshooting

- **Container fails to start**: Check Docker is running and has enough resources allocated
- **Tools not found**: Run `source ~/.bashrc` or restart the terminal after container creation
- **Port conflicts**: Change the port forwarding in `devcontainer.json` if port 8080 is already in use
