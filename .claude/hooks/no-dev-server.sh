#!/bin/sh
# PreToolUse hook (Bash): block dev server launch.
# Dev servers (task run, go run, pnpm dev) occupy ports and conflict
# with the user's own dev server. Use `go build` or `task check` instead.

[ -z "$CLAUDE_TOOL_INPUT" ] && exit 0

command=$(printf '%s' "$CLAUDE_TOOL_INPUT" | sed -n 's/.*"command"[[:space:]]*:[[:space:]]*"\(.*\)".*/\1/p')

case "$command" in
  *"task run"*|*"task dev"*|*"go run"*|*"pnpm dev"*|*"pnpm start"*|*"npm run dev"*|*"npm start"*)
    cat <<'MSG'
❌ BLOCKED: Do not start dev servers (task run, go run, pnpm dev).
They occupy ports and conflict with the user's own dev server.
For verification, use:
  - API: `go build ./...` or `task check`
  - Web: `pnpm build` or `pnpm check`
MSG
    exit 2
    ;;
esac

exit 0
