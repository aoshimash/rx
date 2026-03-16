# Rx Web Frontend

Web application for Rx - training program planning, workout management, and data visualization.

## Technology Stack

| Category | Technology |
|----------|------------|
| Framework | Next.js (App Router) |
| Language | TypeScript (no JavaScript) |
| Linter/Formatter | Biome |
| Styling | Tailwind CSS |

## Prerequisites

All tools are managed via [aqua](https://aquaproj.github.io/). Install aqua first, then run from the repository root:

```bash
aqua install
```

This installs Node.js, pnpm, and Biome.

## Getting Started

### Installation

```bash
cd web
pnpm install
```

### Development

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Linting and Formatting

```bash
pnpm check       # Run lint and format check
pnpm check:fix   # Fix lint and format issues
```

### Build

```bash
pnpm build
pnpm start
```

## Project Structure

```
web/
├── app/                    # Next.js App Router
│   ├── layout.tsx          # Root layout
│   └── page.tsx            # Home page
├── public/                 # Static assets
├── biome.json              # Biome configuration
├── next.config.ts          # Next.js configuration
├── package.json            # Dependencies and scripts
└── tsconfig.json           # TypeScript configuration
```

## Configuration

### Environment Variables

Create `.env.local` for local development:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Conventions

See [Frontend Architecture](../docs/FRONTEND_ARCHITECTURE.md) for coding standards.

## Related Documentation

- [Frontend Architecture](../docs/FRONTEND_ARCHITECTURE.md) - Architecture decisions
- [API Documentation](../api/README.md) - Backend API reference
