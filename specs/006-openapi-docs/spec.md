# Feature Specification: OpenAPI Documentation Tool Integration

**Feature Branch**: `006-openapi-docs`  
**Created**: 2026-01-30  
**Status**: Draft  
**Input**: User description: "OpenAPIのドキュメンテーションツールを導入し、ローカルでAPIドキュメントが確認できるようにする。将来的にはGitHub Pagesで公開する予定だが、今は不要。ドキュメンテーションツールはこれを使う https://github.com/scalar/scalar"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View API Documentation Locally (Priority: P1)

As a developer working on the OPTel Training API, I want to view beautiful, interactive API documentation in my local browser so that I can understand available endpoints, request/response formats, and test API calls directly from the documentation.

**Why this priority**: This is the core functionality requested. Without the ability to view documentation locally, the feature has no value. This enables developers to quickly understand and interact with the API during development.

**Independent Test**: Can be fully tested by starting the documentation server and opening the browser to verify all endpoints are displayed with correct request/response schemas, and delivers immediate developer productivity improvement.

**Acceptance Scenarios**:

1. **Given** the developer has cloned the repository and installed dependencies, **When** they run the documentation server command, **Then** the API documentation is accessible at a local URL (e.g., http://localhost:8080/docs)
2. **Given** the documentation server is running, **When** the developer opens the documentation URL in a browser, **Then** they see a modern, interactive API reference with all endpoints from the OpenAPI specification
3. **Given** the documentation is displayed, **When** the developer clicks on an endpoint, **Then** they see detailed information including request parameters, request body schema, response codes, and response schemas

---

### User Story 2 - Try API Requests from Documentation (Priority: P2)

As a developer, I want to send test requests directly from the API documentation interface so that I can quickly verify API behavior without switching to external tools like curl or Postman.

**Why this priority**: This enhances developer experience significantly by providing an integrated testing capability, reducing context switching and improving productivity.

**Independent Test**: Can be tested by making a test request from the documentation UI to a running API server and verifying the response is displayed correctly.

**Acceptance Scenarios**:

1. **Given** the documentation is displayed and the API server is running, **When** the developer fills in request parameters and clicks "Send", **Then** the actual API request is made and the response is displayed in the documentation interface
2. **Given** the developer is viewing an endpoint that requires a request body, **When** they view the request example, **Then** they see a pre-populated example based on the OpenAPI schema that they can modify

---

### Edge Cases

- What happens when the OpenAPI specification file is invalid or missing?
  - The documentation tool should display a clear error message indicating the issue
- What happens when the developer tries to access documentation before starting the server?
  - Standard browser "connection refused" behavior; documentation should include clear startup instructions
- How does the system handle OpenAPI specification updates?
  - Documentation should reflect the latest specification when the server is restarted (or automatically if hot-reload is supported)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a command to start a local documentation server that serves the API reference via Docker Compose (extending the existing `docker-compose.yml`)
- **FR-002**: System MUST use Scalar (https://github.com/scalar/scalar) as the documentation rendering tool
- **FR-003**: System MUST render documentation from the OpenAPI specification file at `api/openapi.yaml`
- **FR-004**: Documentation MUST display all API endpoints with their HTTP methods, paths, and descriptions
- **FR-005**: Documentation MUST show request parameters (path, query, header) with their types and descriptions
- **FR-006**: Documentation MUST show request body schemas with examples
- **FR-007**: Documentation MUST show response schemas for all documented status codes
- **FR-008**: Documentation MUST provide an interactive "Try it" feature to send requests to a running API server
- **FR-009**: System MUST document the startup command in the project's development documentation
- **FR-010**: System MUST NOT require external hosting or deployment for local usage (all resources served locally)

### Key Entities

- **OpenAPI Specification**: The existing YAML file at `api/openapi.yaml` defining the API contract
- **Documentation Server**: A local server that serves the Scalar-rendered API reference
- **API Reference UI**: The interactive web interface displayed to developers

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can view complete API documentation within 30 seconds of running the startup command
- **SC-002**: 100% of documented API endpoints are visible and correctly rendered in the documentation
- **SC-003**: Developers can successfully send test requests to all endpoints from the documentation interface
- **SC-004**: Documentation startup requires a single command with no additional configuration for basic usage
- **SC-005**: Documentation loads and is fully interactive within 3 seconds on a standard development machine

## Clarifications

### Session 2026-01-30

- Q: ドキュメンテーションサーバーの起動方法は？ → A: Docker Compose に追加（既存の docker-compose.yml を拡張）
- Q: OpenAPI 仕様ファイルの場所は？ → A: `api/openapi.yaml`

## Assumptions

- The existing OpenAPI specification file in the repository is valid and up-to-date
- Developers have Docker available for running the documentation server (standard development environment assumption)
- The API server can be run locally for testing the "Try it" feature
- Future GitHub Pages deployment is out of scope for this feature; the architecture should not preclude it but does not need to explicitly support it
