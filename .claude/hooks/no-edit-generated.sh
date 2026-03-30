#!/bin/sh
# PreToolUse hook (Edit): block direct edits to generated files.
# server.gen.go is generated from openapi.yaml via `task generate`.
# Direct edits will be overwritten on next generation.

[ -z "$CLAUDE_FILE_PATHS" ] && exit 0

while IFS= read -r filepath; do
  case "$filepath" in
    *server.gen.go|*_gen.go)
      cat <<'MSG'
❌ BLOCKED: Do not edit generated files directly.
server.gen.go is auto-generated from openapi/openapi.yaml.
Instead: edit openapi.yaml, then run `task generate` from api/.
MSG
      exit 2
      ;;
  esac
done <<EOF
$CLAUDE_FILE_PATHS
EOF

exit 0
