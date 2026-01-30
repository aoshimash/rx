# Tasks: OpenAPI Documentation Tool Integration

**Input**: Design documents from `/specs/006-openapi-docs/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: No automated tests required for this feature (infrastructure-only change, manual verification).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

This feature modifies infrastructure files at the repository root:
- `docker-compose.yml` - Docker Compose configuration
- `.env.example` - Environment variable template
- `docs/DEVELOPMENT.md` - Developer documentation

---

## Phase 1: Setup (Docker Compose Configuration)

**Purpose**: Add Scalar API Reference service to Docker Compose

- [x] T001 Add `api-docs` service to `docker-compose.yml` with Scalar API Reference image

**Checkpoint**: Docker Compose configuration is ready for documentation server

---

## Phase 2: User Story 1 - View API Documentation Locally (Priority: P1) 🎯 MVP

**Goal**: Developers can view beautiful, interactive API documentation in their local browser

**Independent Test**: Start the documentation server with `docker compose up -d api-docs`, open `http://localhost:8081` in browser, and verify all API endpoints are displayed with correct request/response schemas

### Implementation for User Story 1

- [x] T002 [US1] Verify `api-docs` service starts successfully with `docker compose up -d api-docs`
- [x] T003 [US1] Verify documentation is accessible at `http://localhost:8081`
- [x] T004 [US1] Verify all API endpoints from `api/openapi/openapi.yaml` are displayed correctly
- [x] T005 [US1] Verify endpoint details (request parameters, request body, response schemas) are shown

**Checkpoint**: At this point, User Story 1 should be fully functional - developers can view API documentation locally

---

## Phase 3: User Story 2 - Try API Requests from Documentation (Priority: P2)

**Goal**: Developers can send test requests directly from the API documentation interface

**Independent Test**: Start both API server and documentation server, use the "Try it" feature on any endpoint, and verify the request is sent and response is displayed

### Implementation for User Story 2

- [x] T006 [US2] Start API server alongside documentation server with `docker compose up -d api api-docs`
- [x] T007 [US2] Verify "Try it" feature is available on endpoints in the documentation UI
- [x] T008 [US2] Verify test request can be sent to a running API endpoint (e.g., GET /api/v1/workouts)
- [x] T009 [US2] Verify response is displayed correctly in the documentation interface

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - developers can view docs and test API requests

---

## Phase 4: Polish & Documentation

**Purpose**: Update developer documentation and finalize the feature

- [x] T010 Add "API Documentation" section to `docs/DEVELOPMENT.md`
- [x] T011 [P] Document startup command `docker compose up -d api-docs` in `docs/DEVELOPMENT.md`
- [x] T012 [P] Document access URL `http://localhost:8081` in `docs/DEVELOPMENT.md`
- [x] T013 Document "Try it" feature usage with running API server in `docs/DEVELOPMENT.md`
- [x] T014 Run quickstart.md validation (manual test of all documented commands)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **User Story 1 (Phase 2)**: Depends on Phase 1 completion
- **User Story 2 (Phase 3)**: Depends on Phase 2 completion (needs working documentation first)
- **Polish (Phase 4)**: Can start after Phase 1, but should complete after Phase 3

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on Phase 1 (Setup) - Can be completed independently
- **User Story 2 (P2)**: Depends on User Story 1 - Requires documentation server to be working first

### Within Each Phase

- T001 is a single task (Phase 1)
- T002-T005 are verification steps that must run sequentially
- T006-T009 are verification steps that must run sequentially
- T010-T013 documentation updates: T011-T012 can run in parallel after T010

### Parallel Opportunities

- **Phase 4**: T011 and T012 can run in parallel after T010 creates the section creates the section

---

## Parallel Example: Phase 4 Documentation

```bash
# After T010 creates the section, launch T011 and T012 together:
Task: "Document startup command in docs/DEVELOPMENT.md"
Task: "Document access URL in docs/DEVELOPMENT.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: User Story 1 (T002-T005)
3. **STOP and VALIDATE**: Test documentation server works independently
4. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup → Docker Compose configured
2. Add User Story 1 → Test independently → Developers can view docs (MVP!)
3. Add User Story 2 → Test independently → Developers can test API from docs
4. Complete Polish → Documentation updated, feature complete

### Task Breakdown by File

| File | Tasks |
|------|-------|
| `docker-compose.yml` | T001 |
| `docs/DEVELOPMENT.md` | T010, T011, T012, T013 |
| Manual verification | T002-T009, T014 |

---

## Docker Compose Service Configuration

Reference configuration from `contracts/README.md`:

```yaml
api-docs:
  image: scalarapi/api-reference:latest
  container_name: optel-workout-api-docs
  ports:
    - "${API_DOCS_PORT:-8081}:80"
  volumes:
    - ./api/openapi/openapi.yaml:/app/public/openapi.yaml:ro
  healthcheck:
    test: ["CMD", "wget", "-q", "--spider", "http://localhost:80/health"]
    interval: 10s
    timeout: 5s
    retries: 3
  restart: unless-stopped
```

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- This feature is infrastructure-only - no application code changes
- Manual verification is used instead of automated tests (appropriate for Docker Compose configuration)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
