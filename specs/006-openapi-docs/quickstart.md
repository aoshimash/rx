# Quickstart: OpenAPI Documentation Tool Integration

**Feature Branch**: `006-openapi-docs`  
**Date**: 2026-01-30

## Prerequisites

- Docker and Docker Compose installed
- Repository cloned locally

## Quick Start

### 1. Start the Documentation Server

```bash
docker compose up -d api-docs
```

### 2. Open API Documentation

Open your browser and navigate to:

```
http://localhost:8081
```

You should see the Scalar API Reference with all OPTel Training API endpoints.

### 3. Test API Endpoints (Optional)

To use the "Try it" feature:

1. Start the API server:
   ```bash
   docker compose up -d api
   ```

2. In the documentation UI, click on any endpoint
3. Click "Send" to make a test request
4. View the response in the documentation interface

## Common Commands

| Command | Description |
|---------|-------------|
| `docker compose up -d api-docs` | Start documentation server |
| `docker compose stop api-docs` | Stop documentation server |
| `docker compose logs api-docs` | View documentation server logs |
| `docker compose up -d api api-docs` | Start both API and documentation |

## Configuration

### Custom Port

To use a different port, set the `API_DOCS_PORT` environment variable:

```bash
API_DOCS_PORT=3000 docker compose up -d api-docs
```

Then access at `http://localhost:3000`

### Using .env File

Add to `.env`:

```
API_DOCS_PORT=8081
```

## Troubleshooting

### Documentation shows "No OpenAPI document found"

Ensure the OpenAPI specification exists:
```bash
ls api/openapi/openapi.yaml
```

### "Try it" requests fail

1. Ensure the API server is running:
   ```bash
   docker compose ps api
   ```

2. Check API server logs:
   ```bash
   docker compose logs api
   ```

### Container health check fails

Check container logs:
```bash
docker compose logs api-docs
```

## Development Workflow

1. Edit `api/openapi/openapi.yaml`
2. Restart documentation server to see changes:
   ```bash
   docker compose restart api-docs
   ```

Or for automatic updates (development mode):
```bash
docker compose up api-docs  # foreground mode, restart manually as needed
```
