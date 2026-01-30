# AGENTS.md

This file provides guidance for AI coding agents working on this repository.

## Project Overview

**OPTel (Open Physical Telemetry)** - An observability stack for the human physical layer.

This repository (`optel-workout`) is the **Workout Management** component. It manages "Workouts" (physical exertion records) as an agent-native telemetry backend.

## Key Principles

1. **"Dumb Backend"** - No business logic for "health." Strictly stores and retrieves telemetry data.
2. **Domain-Driven Schema-First** - Domain models define business logic, OpenAPI spec defines API contract. Code is generated from OpenAPI spec.

For details, see `.claude/skills/optel-philosophy/`.

## Project Structure

```
optel-workout/
├── api/                  # REST API (Go)
├── mcp/                  # MCP Server (runs on user's local machine)
├── frontend/             # Frontend (future)
├── infra/                # Terraform/Helm (future)
├── docs/                 # Documentation
└── .claude/skills/       # AI agent skills
```

## Skills Reference

| Skill | Description |
|-------|-------------|
| [optel-philosophy](.claude/skills/optel-philosophy/) | Core philosophy and constraints |
| [optel-domain](.claude/skills/optel-domain/) | Domain models (Workout, Program, Telemetry) |
| [optel-go-standards](.claude/skills/optel-go-standards/) | Go coding standards |

## Quick Reference

- **Language**: Go 1.25+
- **HTTP Server**: chi
- **OpenAPI**: oapi-codegen (Domain-Driven Schema-First)
- **Linter**: golangci-lint (strict)
- **Logging**: log/slog
- **Testing**: standard testing package

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [Development Guide](docs/DEVELOPMENT.md)

## AI Agent Command Execution

When executing commands that require host resources (network access, Docker socket access, file system writes), AI agents **MUST** use `required_permissions: ["all"]` parameter.

### Commands Requiring `required_permissions: ["all"]`

| Command | Reason |
|---------|--------|
| `aqua install` | Requires network access (GitHub API) and file system writes (`~/.local/share/aquaproj-aqua/`) |
| `docker compose up/down` | Requires Docker socket access |
| `docker compose pull` | Requires network access and Docker socket access |
| `go get`, `go install` | Requires network access (if downloading packages) |

### Example Usage

```python
# Correct: Using required_permissions
run_terminal_cmd(
    command="aqua install",
    required_permissions=["all"]
)

# Correct: Docker Compose commands
run_terminal_cmd(
    command="docker compose up -d postgres",
    required_permissions=["all"]
)
```

### Important Notes

**Sandbox Limitations:**
- Cursor's sandbox may block certain DNS resolution methods (`ping`, `nslookup`) even with `required_permissions: ["all"]`
- Some commands may fail with "no such host" or "operation not permitted" errors even when network access is available via `curl`
- Docker socket access may be restricted even with `required_permissions: ["all"]`

**Workaround:**
- If commands fail with permission errors, instruct the user to run them manually
- Configuration files (`aqua.yaml`, `docker-compose.yml`) are correctly set up and will work when executed manually
- The user can run `aqua install` and `docker compose up -d postgres` directly in their terminal

**Note:** These commands should trigger a user confirmation prompt before execution. If the prompt doesn't appear, it may be a Cursor configuration issue, and manual execution is recommended.

### Troubleshooting: Permission Prompts Not Appearing

If `required_permissions: ["all"]` doesn't trigger a confirmation prompt:

1. **Check Cursor Settings:**
   - Open Cursor Settings (`Cmd + ,` on macOS, `Ctrl + ,` on Windows/Linux)
   - Search for "sandbox" or "permissions"
   - Ensure sandbox is enabled and permission prompts are enabled

2. **Check Cursor Configuration File:**
   - macOS: `~/Library/Application Support/Cursor/User/settings.json`
   - Windows: `%APPDATA%\Cursor\User\settings.json`
   - Linux: `~/.config/Cursor/User/settings.json`
   - Look for `cursor.sandbox.*` or `cursor.permissions.*` settings
   - Note: These settings may not be visible in the UI and may need to be added manually

3. **Manual Execution:**
   - Run commands directly in your terminal
   - Install tools: `aqua install`
   - Start PostgreSQL: `docker compose up -d postgres`

4. **Alternative:**
   - AI agent will ask for explicit confirmation before running commands
   - You can approve by responding "yes" or "proceed"

## Pre-commit Hook Enforcement

**CRITICAL**: AI agents **MUST NOT** use `git commit --no-verify` to skip pre-commit hooks.

### Pre-commit Checks

Before committing, the following checks are automatically executed:
1. Code formatting (`task format`)
2. Linting (`task lint`)
3. Tests (`task test` with race detection)

### AI Agent Commit Workflow

1. **Before committing**: Ensure all checks pass locally
   ```bash
   cd api
   task check  # Runs format + lint + test
   ```

2. **If checks fail**: Fix errors before committing
   - Format errors: Run `go fmt ./...`
   - Lint errors: Fix according to golangci-lint output
   - Test failures: Fix failing tests

3. **Commit**: Use standard `git commit` (pre-commit hook will run automatically)
   - ❌ **DO NOT** use `git commit --no-verify`
   - ✅ Use `git commit -m "message"` (hook runs automatically)

4. **If hook fails**: The commit will be aborted. Fix errors and try again.

### Setup

Run the setup script to install pre-commit hooks:
```bash
./scripts/setup-githooks.sh
```

This configures Git to use hooks from `githooks/` directory, ensuring all developers (including AI agents) run the same checks.
