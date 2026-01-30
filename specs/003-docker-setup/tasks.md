# Tasks: Docker Development Environment Setup

**Input**: Design documents from `/Users/aoshima/dev/github/aoshimash/optel-workout/specs/003-docker-setup/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`  
**Tests**: No new automated tests requested; include smoke-check tasks where relevant.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare repository structure and baseline documentation updates for this feature.

- [x] T001 Create `.devcontainer/` directory at repo root
- [x] T002 Create `specs/003-docker-setup/tasks.md` (this file) and keep task IDs stable as work progresses
- [x] T003 [P] Add a short note to `README.md` clarifying: devcontainer for development, docker compose for smoke-testing

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the target Docker layout and remove legacy dev-container assumptions from docs/tooling.

- [x] T004 Update `docs/DEVELOPMENT.md` to describe DevContainer as the primary development workflow (remove guidance implying Docker Compose is the only dev toolchain)
- [x] T005 [P] Update `api/README.md` to align with the new DevContainer-first workflow and the single `api/Dockerfile`
- [x] T006 Update `api/Makefile` to run generate/lint/test/run natively (tools must be on PATH; DevContainer provides them automatically). Remove hard dependency on `docker compose exec dev` and ensure commands work when tools are installed locally or via DevContainer.

**Checkpoint**: Documentation and tooling no longer assume a Docker Compose “dev tool container”.

---

## Phase 3: User Story 1 - VS Code/Cursor DevContainer Integration (Priority: P1) 🎯 MVP

**Goal**: Open the repo in VS Code/Cursor and get a working, consistent dev environment automatically.

**Independent Test**:

- Open in VS Code/Cursor → reopen in DevContainer → run repository generate/lint/test workflows successfully.

### Implementation for User Story 1

- [x] T007 [US1] Add `.devcontainer/devcontainer.json` for a Go 1.25+ toolchain and repo workspace mounting
- [x] T008 [P] [US1] Add `.devcontainer/README.md` describing how to use DevContainer in VS Code/Cursor
- [x] T009 [P] [US1] Add `.devcontainer/extensions.json` with recommended extensions (Go, Docker, YAML, etc.)
- [x] T010 [US1] Ensure required dev tools are installed in DevContainer (Go 1.25+, oapi-codegen, golangci-lint) via devcontainer features and/or post-create hook, and verify they are available on PATH
- [x] T011 [US1] Add port forwarding configuration in `.devcontainer/devcontainer.json` (API port) for local runs
- [x] T012 [US1] Validate DevContainer workflow by running generate/lint/test and documenting any required environment variables in `.devcontainer/README.md`

**Checkpoint**: User Story 1 is independently usable for local development.

---

## Phase 4: User Story 2 - Production Docker Image Distribution (Priority: P2)

**Goal**: Distribute a production-ready container image built from a single Dockerfile at `api/Dockerfile`.

**Independent Test**:

- Build an image from `api/Dockerfile`
- Run container with port exposed
- Verify the API is reachable via HTTP on the configured port

### Implementation for User Story 2

- [x] T013 [US2] Replace `api/Dockerfile` contents with the production build (multi-stage, minimal runtime, non-root, static binary)
- [x] T013a [US2] Verify FR-016: production image builds successfully and runs without requiring source code access after build (test by building image, removing source directory, and running container)
- [x] T014 [US2] Remove legacy `api/Dockerfile.prod`
- [x] T015 [US2] Update `api/.dockerignore` if needed to keep the production build context minimal and deterministic
- [x] T016 [P] [US2] Update `docs/DEVELOPMENT.md` and `README.md` with the production image build/run instructions (without reintroducing a dev Dockerfile)
- [x] T017 [US2] Add a smoke-check section to `docs/DEVELOPMENT.md` for verifying the built image responds to HTTP requests

**Checkpoint**: A single `api/Dockerfile` builds a minimal production image and runs successfully.

---

## Phase 5: User Story 3 - Local Testing with Docker Compose (Priority: P3)

**Goal**: Run production-like local smoke-testing using `docker-compose.yml`.

**Independent Test**:

- `docker compose up -d` starts the API service
- API is reachable via HTTP on the configured host port
- `docker compose logs` shows application logs
- `docker compose down` stops cleanly

### Implementation for User Story 3

- [x] T018 [US3] Update `docker-compose.yml` to build the API service with `context: ./api` and `dockerfile: Dockerfile` (single production Dockerfile, per FR-012)
- [x] T019 [US3] Rename the service from `dev` to `api` (or similar) in `docker-compose.yml` to reflect smoke-testing purpose
- [x] T020 [US3] Remove development-only volume mounts and “keep-alive” commands from `docker-compose.yml` (production-like runtime)
- [x] T021 [US3] Ensure `docker-compose.yml` supports configurable host port and environment variables (document defaults in `docs/DEVELOPMENT.md`)
- [x] T022 [US3] Add a short smoke-test walkthrough to `README.md` (compose up/logs/down)

**Checkpoint**: User Story 3 is independently usable for local smoke-testing.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Consistency, cleanup, and verification across all stories.

- [x] T023 [P] Remove/adjust any stale references to `Dockerfile.prod` across docs (`README.md`, `docs/DEVELOPMENT.md`, `api/README.md`)
- [x] T024 [P] Remove/adjust any stale references to a Compose-based dev tool container in docs and make targets
- [x] T025 Run a full verification pass:
  - DevContainer: reopen and run generate/lint/test
  - Production image: build and run, verify HTTP response
  - Docker Compose: up/logs/down, verify HTTP response
- [x] T026 [P] Update `specs/003-docker-setup/quickstart.md` if the implementation details diverged from the planned steps

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup completion
- **User Story 1 (Phase 3)**: Can start after Setup; recommended before other stories
- **User Story 2 (Phase 4)**: Depends on Foundational (Makefile/docs alignment) and clarifications already captured in spec
- **User Story 3 (Phase 5)**: Depends on User Story 2 (compose must build from the single production Dockerfile)
- **Polish (Phase 6)**: After desired user stories are complete

### User Story Dependencies

- **US1 (P1)**: Independent (DevContainer)
- **US2 (P2)**: Independent distribution artifact; should not depend on US1
- **US3 (P3)**: Depends on US2 (compose uses the single production Dockerfile)

### Parallel Opportunities

- Documentation updates across `README.md`, `docs/DEVELOPMENT.md`, `api/README.md` can run in parallel ([P]).
- DevContainer editor configuration files inside `.devcontainer/` can be created in parallel ([P]).

---

## Parallel Example: User Story 1

```bash
Task: "Add .devcontainer/extensions.json with recommended extensions"
Task: "Add .devcontainer/README.md describing DevContainer usage"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 (Setup)
2. Implement Phase 3 (US1) and validate the DevContainer workflow end-to-end
3. Stop and validate before touching Dockerfile consolidation

### Incremental Delivery

1. US1: DevContainer development workflow
2. US2: Single production Dockerfile in `api/Dockerfile` + remove legacy `Dockerfile.prod`
3. US3: Production-like `docker-compose.yml` for smoke-testing

