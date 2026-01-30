# Quickstart: PostgreSQL導入

**Feature**: 005-postgresql-setup  
**Date**: 2026-01-30

## Prerequisites

- Docker & Docker Compose
- Go 1.25+
- make

## Quick Setup

### 1. Start PostgreSQL

```bash
# Start PostgreSQL container
docker compose up -d postgres

# Verify it's running
docker compose ps
# Expected: optel-workout-postgres running, healthy
```

### 2. Run Migrations

```bash
cd api

# Run all migrations
make migrate

# Check migration status
make migrate-status
```

### 3. Start API Server

```bash
cd api

# With PostgreSQL (default)
make run

# Or with in-memory store (for quick testing)
STORAGE_TYPE=memory make run
```

### 4. Verify Connection

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","database":"connected"}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | optel | Database user |
| `DB_PASSWORD` | optel | Database password |
| `DB_NAME` | optel_workout | Database name |
| `DB_SSLMODE` | disable | SSL mode |
| `DB_MAX_CONNS` | 10 | Max pool connections |
| `STORAGE_TYPE` | postgres | Storage type (postgres/memory) |

## Common Commands

```bash
# Database
make migrate           # Run migrations
make migrate-down      # Rollback all migrations
make migrate-down-1    # Rollback last migration
make migrate-status    # Check migration status
make migrate-create NAME=xxx  # Create new migration

# Development
make run              # Start server
make test             # Run tests (uses Testcontainers)
make test-integration # Run integration tests only

# Docker
docker compose up -d postgres   # Start PostgreSQL only
docker compose down             # Stop all services
docker compose logs postgres    # View PostgreSQL logs
```

## Troubleshooting

### Connection Refused

```bash
# Check if PostgreSQL is running
docker compose ps

# Check logs
docker compose logs postgres

# Verify port is accessible
nc -zv localhost 5432
```

### Migration Failed

```bash
# Check migration status
make migrate-status

# Force to specific version (use with caution)
make migrate-force VERSION=1

# View migration files
ls api/migrations/
```

### Test Database Issues

Tests use Testcontainers which automatically:
1. Pulls PostgreSQL 17 image
2. Starts a temporary container
3. Runs migrations
4. Executes tests
5. Cleans up container

If tests fail with Docker errors:
```bash
# Ensure Docker is running
docker info

# Clean up stale containers
docker container prune -f
```

## Next Steps

1. **Verify CRUD operations**: Test all endpoints with PostgreSQL backend
2. **Check data persistence**: Restart server and verify data persists
3. **Run integration tests**: `make test-integration`
4. **Monitor performance**: Check connection pool usage in logs
