# Feature Specification: Docker Development Environment Setup

**Feature Branch**: `003-docker-setup`  
**Created**: 2026-01-25  
**Status**: Draft  
**Input**: User description: "ローカル開発用にdevcontainer環境を、配布用にdocker imageを、ローカルでの動作確認用にdocker composeを用意する"

## Clarifications

### Session 2026-01-25

- Q: Where should the single production Dockerfile live? → A: `api/Dockerfile`

## User Scenarios & Testing *(mandatory)*

### User Story 1 - VS Code/Cursor DevContainer Integration (Priority: P1)

Developers using VS Code or Cursor can open the project and have a fully configured development environment automatically set up with all required tools, dependencies, and extensions.

**Why this priority**: DevContainer provides the best developer experience for local development. It eliminates manual setup steps, ensures consistency across team members, and integrates seamlessly with IDE features like IntelliSense, debugging, and terminal access.

**Independent Test**: Can be fully tested by opening the project in VS Code/Cursor with the Dev Containers extension installed. The container should start automatically, required development tools for this repository should be available, and the developer should be able to run the repository's standard generate/lint/test workflow without any additional setup.

**Acceptance Scenarios**:

1. **Given** a developer has VS Code/Cursor with Dev Containers extension installed, **When** they open the project folder, **Then** they are prompted to reopen in container, and after doing so, all development tools are available in the integrated terminal
2. **Given** the devcontainer is running, **When** the developer runs the repository's code generation workflow, **Then** code generation completes successfully without errors
3. **Given** the devcontainer is running, **When** the developer runs the repository's lint workflow, **Then** linting executes and reports any code quality issues
4. **Given** the devcontainer is running, **When** the developer runs the repository's test workflow, **Then** all tests execute successfully
5. **Given** the devcontainer is running, **When** the developer opens a Go file, **Then** IntelliSense and code completion work correctly

---

### User Story 2 - Production Docker Image Distribution (Priority: P2)

The project can be distributed as a Docker image that can be deployed to any container orchestration platform (Kubernetes, Docker Swarm, cloud container services) or run standalone.

**Why this priority**: Production images enable deployment and distribution of the application. While not needed for initial development, having a production-ready image is essential for real-world usage and CI/CD pipelines.

**Independent Test**: Can be fully tested by building the production image and running it in a clean environment. The container should start successfully, respond to health checks (when implemented), and serve API requests without requiring any development tools or source code.

**Acceptance Scenarios**:

1. **Given** a developer has Docker installed, **When** they build the production image using the provided Dockerfile, **Then** the build completes successfully and creates a minimal, secure image (note: there is only one Dockerfile for production use, no separate development Dockerfile)
2. **Given** the production image is built, **When** it is run with the container port exposed to the host, **Then** the container starts and the API is accessible via HTTP on the configured host port
3. **Given** the production image is running, **When** a client sends an HTTP request, **Then** the API responds correctly with appropriate status codes and data
4. **Given** the production image, **When** it is inspected, **Then** it uses a non-root user and has a minimal attack surface (e.g., minimal runtime with no shell/package manager)
5. **Given** the production image, **When** it is pushed to a container registry, **Then** other users can pull and run it without requiring source code or build tools

---

### User Story 3 - Local Testing with Docker Compose (Priority: P3)

Developers can quickly start the entire application stack locally using Docker Compose for testing, integration testing, or demonstration purposes without needing to set up the development environment.

**Why this priority**: Docker Compose provides a convenient way to test the application in a production-like environment locally. It's useful for integration testing, demonstrating the application, and verifying that all components work together correctly.

**Independent Test**: Can be fully tested by running `docker compose up` from the project root. The application should start, be accessible via HTTP, and all services should be properly configured and connected.

**Clarified intent**: Docker Compose is for local smoke-testing in a production-like way (using the production container build), while iterative development happens in devcontainer.

**Acceptance Scenarios**:

1. **Given** a developer has Docker and Docker Compose installed, **When** they run `docker compose up` from the project root, **Then** all services defined in docker-compose.yml start successfully
2. **Given** docker compose services are running, **When** the developer accesses `http://localhost:8080/api/v1/workouts`, **Then** the API responds correctly
3. **Given** docker compose services are running, **When** the developer runs `docker compose logs`, **Then** they can see application logs for debugging
4. **Given** docker compose services are running, **When** the developer runs `docker compose down`, **Then** all containers stop and are removed cleanly
5. **Given** docker compose is configured, **When** future services (e.g., PostgreSQL) are added, **Then** they can be added to the same compose file and started together

### Edge Cases

- What happens when port 8080 is already in use on the host machine?
- How does the system handle Docker daemon not running?
- What happens when the devcontainer fails to build or start?
- How does the system handle missing or outdated Docker/Compose versions?
- What happens when the production image is run without required environment variables?
- How does the system handle network connectivity issues when pulling base images?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a devcontainer configuration (`.devcontainer/devcontainer.json`) that automatically sets up a development environment with all required development tools for this repository
- **FR-002**: DevContainer MUST mount the project source code as a volume so changes are immediately reflected
- **FR-003**: DevContainer MUST configure VS Code/Cursor extensions (Go extension, recommended extensions) for optimal development experience
- **FR-004**: DevContainer MUST set up the working directory and environment variables correctly
- **FR-005**: System MUST provide a Dockerfile at `api/Dockerfile` that builds a minimal, secure container image suitable for production distribution
- **FR-006**: System MUST NOT include a development Dockerfile (development environment is provided by devcontainer, not a separate Dockerfile)
- **FR-007**: Dockerfile MUST use a multi-stage build to minimize final image size
- **FR-008**: Dockerfile MUST use a minimal runtime base image for security (e.g., no shell/package manager)
- **FR-009**: Dockerfile MUST run as a non-root user
- **FR-010**: Dockerfile MUST build a static binary for compatibility with minimal runtime images
- **FR-011**: System MUST provide a docker-compose.yml file that defines all services needed for local testing
- **FR-012**: Docker Compose MUST use `api/Dockerfile` for building the application service (e.g., `context: ./api` and `dockerfile: Dockerfile`)
- **FR-013**: Docker Compose MUST expose the API on a configurable port (default 8080)
- **FR-014**: Docker Compose MUST support environment variable configuration
- **FR-015**: All Docker configurations MUST be documented with clear instructions for usage
- **FR-016**: Production image MUST be buildable without requiring source code access after build
- **FR-017**: DevContainer MUST support both development workflow (live code editing) and build workflow (compiling and testing)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can set up a complete development environment by opening the project in VS Code/Cursor with Dev Containers extension, with zero manual configuration steps required
- **SC-002**: Production Docker image builds successfully in under 2 minutes on a standard development machine
- **SC-003**: Production Docker image size is under 20MB (binary and minimal runtime only)
- **SC-004**: Developers can start the entire application stack locally using Docker Compose in under 30 seconds
- **SC-005**: All three environments (devcontainer, production image, docker compose) can run simultaneously on the same machine without port conflicts when properly configured
- **SC-006**: 100% of development tasks (code generation, linting, testing, building) can be performed within the devcontainer without requiring local tool installation
- **SC-007**: Production image starts and responds to API requests within 5 seconds of container start
- **SC-008**: Documentation enables new team members to use any of the three environments (devcontainer, production image, docker compose) within 10 minutes of reading the setup instructions

## Assumptions

- Developers have Docker Desktop (or Docker Engine + Docker Compose) installed
- Developers using devcontainer have VS Code or Cursor with Dev Containers extension installed
- The production image will be distributed via container registries (Docker Hub, GitHub Container Registry, etc.)
- Local testing with Docker Compose does not require persistent data storage initially (in-memory store is sufficient)
- The application runs on Linux containers (amd64 architecture)
- Port 8080 is available for local development (configurable if needed)
- Development tools are only needed in development, not in production
- Development environment is provided exclusively through devcontainer (no separate development Dockerfile needed)
- Existing development Dockerfile (if present) will be removed or deprecated as part of this feature
- The single production Dockerfile is located at `api/Dockerfile`
- Time-based success criteria (build/startup times) assume a standard development machine (e.g., 8+ CPU cores, 16GB+ RAM) and a warm Docker cache unless otherwise stated

## Dependencies

- Docker and Docker Compose must be installed on developer machines
- VS Code/Cursor Dev Containers extension for devcontainer functionality
- Access to container base images from public registries
- A containerized build toolchain suitable for this repository

## Out of Scope

- Separate development Dockerfile (development is handled by devcontainer)
- Kubernetes deployment manifests (future feature)
- Docker Swarm configuration (future feature)
- Multi-architecture builds (ARM, etc.) - initial focus on amd64
- CI/CD pipeline integration (will use these Docker setups but pipeline definition is separate)
- Database containers in docker-compose (Phase 2 feature)
- Health check endpoints implementation (required for production but separate feature)
- Image signing and security scanning automation (future enhancement)
