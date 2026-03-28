# Roadmap

## Plan Driven Training

Rx is built on the principle of **Plan Driven Training** — training guided by a pre-defined plan, not ad-hoc decisions at the gym.

Most training apps are log-centric: you go to the gym and record what you did. Rx inverts this. You define a Program first, convert it into concrete Plans, then log your execution against those Plans. The plan always comes first.

This philosophy mirrors engineering methodologies like Test Driven Development and Domain Driven Design — define the expected outcome first, then execute against it.

## Development Phases

### Phase 1: Core API
- OpenAPI specification for Workout resource
- Go module with clean architecture
- In-memory store for rapid prototyping

### Phase 2: Persistence
- PostgreSQL support
- Database migrations
- Docker Compose for local development

### Phase 3: Program Entity
- Program / Plan / Log lifecycle
- Recursive tree structure
- Convert Program to Plan

### Phase 4: MCP Server
- MCP Server in `mcp/` directory
- Tools for AI agents to query and record data
- Runs on user's local machine, communicates with Rx API via HTTP

### Phase 5: Program Import / Export
- JSON import/export for Programs in the Rx app
- Static site hosting curated Program JSON files
- Published on official homepage to solve cold-start problem for new users
- Users can browse, download, and import Programs directly

## Web Frontend
- Next.js (App Router) with TypeScript
- Plan editor, dashboard
- Ongoing UX improvements

## Future

### Mobile App
- React Native + Expo
- Gym-side training logging, plan viewing
- Offline support (cached plans, sync on reconnect, WiFi-only video upload)

### Authentication
- OAuth2 / JWT with proper login UI
- Multi-user support

### AI Integration
1. External agents via MCP Server (Phase 4)
2. Web-embedded AI chat / analysis
3. AI-powered recommendations

### Infrastructure
- Horizontal scaling (Kubernetes)
- Caching layer (Redis)
- Observability (OpenTelemetry, Prometheus, structured logging)
