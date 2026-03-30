#!/bin/sh
# PostToolUse hook: remind to check docs when key code areas change,
# and remind to run `task generate` after OpenAPI spec changes.
# Receives CLAUDE_FILE_PATHS (newline-separated list of modified files).
# Outputs a reminder to stdout (injected into Claude's context).

[ -z "$CLAUDE_FILE_PATHS" ] && exit 0

reminders=""

while IFS= read -r filepath; do
  case "$filepath" in
    *api/internal/domain/*)
      reminders="${reminders}
- domain layer changed → verify docs/DOMAIN_MODEL.md is still accurate"
      ;;
    *openapi/openapi.yaml)
      reminders="${reminders}
- OpenAPI spec changed → you MUST run \`task generate\` (from api/) before committing
- OpenAPI spec changed → verify docs/DOMAIN_MODEL.md if schemas changed"
      ;;
    *web/app/*/page.tsx)
      reminders="${reminders}
- route structure changed → verify docs/WEB_UI_DESIGN.md screen structure is still accurate"
      ;;
    *api/internal/handler/*.go)
      reminders="${reminders}
- handler changed → verify AGENTS.md error flow table if error handling changed"
      ;;
  esac
done <<EOF
$CLAUDE_FILE_PATHS
EOF

# Deduplicate reminders
if [ -n "$reminders" ]; then
  unique=$(printf '%s\n' "$reminders" | sort -u | sed '/^$/d')
  cat <<REMINDER
⚠️ Doc sync check: code in a doc-mapped area was modified.
$unique
If the change is structural (new entities, renamed fields, new routes), update the corresponding doc before committing.
REMINDER
fi

exit 0
