# Rx Philosophy

## Why Rx Exists

### Why Spreadsheets?

Many serious lifters — especially powerlifters — track their training in spreadsheets. This is not because they lack better tools; it is because spreadsheets are the best tool available. Existing training apps force a fixed schema: sets × reps × weight. But real training programs need fields those apps do not offer — tempo prescriptions (3-1-2-0), RPE/RIR targets, band/chain accommodations, pause lengths, breathing cues, and whatever else a coach or methodology demands. Spreadsheets let each lifter define their own schema. That freedom is the killer feature.

### Why Not Spreadsheets?

Spreadsheets are accessible — Google Sheets has an API, sharing a link is trivial. The problem is not access but **interoperability**:

- **No shared schema** — Every lifter's spreadsheet has a different column layout, naming convention, and unit format. A tool that works with one person's sheet cannot work with another's without custom adaptation.
- **Implicit semantics** — "Column D is RPE" is knowledge that lives in the user's head, not in the data. No external tool can reliably interpret the data without being told what each column means.
- **No separation of plan and execution** — Prescribed weights and actual weights live in the same cells or adjacent columns. Comparing intent vs. outcome across weeks requires manual effort.

The fundamental tension: **flexibility demands less structure, but external utility demands more.**

### Rx's Approach: Minimal Viable Structure

Rx resolves this tension by imposing the minimum structure required for interoperability while preserving maximum field flexibility:

1. **Structured core** — A small set of fields that nearly every lifter shares (exercise, sets, reps, load) are first-class, typed, and queryable. This is what makes the data useful to external tools.
2. **Flexible metadata** — Everything else lives in open-ended metadata fields that each user defines for their own methodology. Rx stores it, serves it, but does not interpret it.
3. **API as the contract** — Both humans (via Web/Mobile) and machines (via API/CLI) interact through the same structured interface. Unlike spreadsheets, the schema is explicit and consistent across all users.

The result: a system where a powerlifter's custom tracking fields are preserved exactly as they want them, while an AI agent can still query "show me all squat sessions where RPE exceeded the plan" — something a spreadsheet cannot support without per-user custom integration.

## Plan Driven Training

Rx follows **Plan Driven Training** — a philosophy where training is guided by a pre-defined plan, not ad-hoc decisions at the gym. Just as Test Driven Development starts with tests and Domain Driven Design starts with domain models, Plan Driven Training starts with a Program that defines the training structure before any weight is lifted.

## What Rx Is

Rx is a **training data backend** that supports Plan Driven Training. The core loop is:

1. **Plan** — An external tool (AI agent, script, or human) creates a training program and registers it via API or CLI
2. **Execute** — User records actual training sessions (via Web, Mobile, API, or CLI)
3. **Review** — An external tool (AI agent, script, or human) reads the data and analyzes it

**Rx itself provides no AI features, no planning logic, and no analysis.** It is a structured data store that makes it easy for external tools to read and write training data.

The Web UI exists to make data entry as frictionless as possible. It is a client of the same API available to external tools.

## Core Principles

### 1. Planning and Analysis are External

- **Planning and analysis happen outside Rx** — humans, AI agents (e.g., Claude Code), or scripts use the API/CLI to create plans and interpret results
- **Rx stores and serves data without interpreting it** — no recommendations, no scoring, no judgment
- **The backend is a tool, not an actor** — it responds to requests; it never initiates

### 2. API-First, Multi-Client

- **API and CLI are required** — AI agents and automation tools interact via API or CLI
- **Web and Mobile are first-class clients** — Full-featured interfaces, not samples or demos
- **No features that only exist in UI** — Every action available in Web/Mobile must also be accessible via API

### 3. No Opinionated Health Logic

- **No health scores, wellness indices, or motivation messages**
- **No recommendations or comparisons** — the backend stores and retrieves, AI analyzes
- **No gamification** — no streaks, badges, or achievements
- Raw data in, raw data out

### 4. Domain-Driven Schema-First Development

- **Domain models define business rules** — `internal/domain/` contains validation and invariants
- **OpenAPI spec defines the API contract** — single source of truth for HTTP API
- **Code generation from OpenAPI** — Go types and server stubs are generated
- **Handlers bridge the gap** — convert between OpenAPI types and domain models
- **Keep them synchronized** — domain models and OpenAPI specs must stay in sync

## What the Backend Does

### Prohibited Features

The backend MUST NOT:

1. Calculate "health scores" or "wellness indices"
2. Provide "motivation" or "encouragement" messages
3. Make recommendations about training intensity
4. Compare users or create leaderboards
5. Implement gamification (streaks, badges, achievements)

### Permitted Features

The backend MAY:

1. Store and retrieve Programs, Plans, and Logs
2. Provide CRUD operations for all resources
3. Expose time-series telemetry data for external analysis
4. Implement filtering and aggregation queries
5. Serve data to AI Agents for analysis (via API or CLI)

### Interface Requirements

All of the following are first-class interfaces:

- **Web** — Full-featured browser client
- **Mobile** — Full-featured native mobile client
- **REST API** — For AI agents and automation
- **CLI** — For scripting and local automation

Every feature MUST be accessible via API (and ideally CLI).

## Terminology

| Term | Description |
|------|-------------|
| Program | A reusable training template containing multiple sessions (created by humans, AI, or import) |
| ProgramSession | A named workout day within a Program, containing exercise prescriptions |
| ProgramGroup | An optional organizational grouping of sessions (e.g., Block, Week), max 2 levels deep |
| Plan | The user's execution queue — an ordered list of upcoming sessions to perform (singleton per user) |
| PlanSession | A concrete single-workout prescription within a Plan, optionally derived from a Program session |
| Log | A record of one training session actually performed, optionally linked to a Program |
| Exercise | A specific movement or activity |
| Rep | A single repetition of an exercise |
| Set | A group of repetitions |
| RPE | Rate of Perceived Exertion (1–10 scale) |
| 1RM | One-rep max: maximum weight for a single repetition |
| load_kg | Weight used for an exercise, in kilograms |
