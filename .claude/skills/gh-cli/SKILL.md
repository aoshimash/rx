---
name: gh-cli
description: Operate GitHub via the gh CLI with low-context, repeatable workflows (PRs, issues, reviews, checks). Use when the user mentions gh, GitHub CLI, pull requests, issues, reviews, CI checks, or wants to automate GitHub operations without the GitHub MCP server.
---

# gh CLI Workflow

## Scope

This skill standardizes common GitHub operations using the `gh` CLI:
- Pull requests (create/view/diff/status/review/merge)
- Issues (list/view/comment/create)
- Checks / CI status (via PR status and runs)

Prefer `gh` over large MCP toolsets when the task fits a known workflow.

## Safety rails (always)

1. **Confirm auth**:
   - `gh auth status`
2. **Confirm repo context** (avoid acting on the wrong repo):
   - `gh repo view`
   - If needed: `gh repo set-default OWNER/REPO`
3. **Confirm branch before any commit-related action**:
   - `git branch --show-current`
4. **Never pass long text via CLI flags**:
   - Prefer `--body-file` with a temporary file.
5. **Do not push / merge unless explicitly requested**.

## Temp file pattern (recommended)

Use this pattern to avoid fragile quoting and "argument too long" style issues:

1. Create a temp file and write the body.
2. Run the `gh` command with `--body-file "$tmp"`.
3. Delete the temp file.

Example:
- `tmp="$(mktemp)"; cat >"$tmp" <<'EOF'
... body ...
EOF
gh <command> --body-file "$tmp"
rm -f "$tmp"`

## Pull request workflows

### Create a PR

1. Ensure the branch is ready (tests/lints if required by the repo).
2. Choose a PR title **without Conventional Commit prefixes**:
   - Good: `Add workout validation for empty telemetry`
   - Bad: `feat(workout): add workout validation for empty telemetry`
2. Create PR using a temp file (structured body):
   - `tmp="$(mktemp)"; cp .github/pull_request_template.md "$tmp"; cat >>"$tmp" <<'EOF'

<!-- Fill in the sections above; add extra context below if needed. -->
EOF
gh pr create --base main --head HEAD --title "..." --body-file "$tmp"
rm -f "$tmp"`
3. Confirm:
   - `gh pr view --web`

### Inspect a PR (local branch or by number)

- Status/checks:
  - `gh pr status`
  - `gh pr checks <number>`
- Diff:
  - `gh pr diff <number>`
- Files list:
  - `gh pr view <number> --json files --jq '.files[].path'`

### Add a PR comment (non-review)

- `tmp="$(mktemp)"; cat >"$tmp" <<'EOF'
...
EOF
gh pr comment <number> --body-file "$tmp"
rm -f "$tmp"`

### Merge a PR (only when explicitly requested)

Pick one merge strategy (follow repo policy):
- `gh pr merge <number> --merge`
- `gh pr merge <number> --squash`
- `gh pr merge <number> --rebase`

If CI is required, wait until checks pass:
- `gh pr checks <number> --watch`

## Issue workflows

### List / view

- `gh issue list`
- `gh issue view <number>`

### Comment

- `tmp="$(mktemp)"; cat >"$tmp" <<'EOF'
...
EOF
gh issue comment <number> --body-file "$tmp"
rm -f "$tmp"`

### Create

- Create using a temp file (structured body):
  - `tmp="$(mktemp)"; cat >"$tmp" <<'EOF'
## Context
- ...

## Goal
- ...

## Acceptance criteria
- [ ] ...

## Out of scope
- ...

## Notes / links
- ...
EOF
gh issue create --title "..." --body-file "$tmp"
rm -f "$tmp"`

## CI / checks (quick)

When you have a PR number:
- `gh pr checks <number>`
- `gh pr checks <number> --watch`

If you need deeper workflow/run inspection, see [reference.md](reference.md).

## Output expectations

When executing `gh` commands, capture:
- The exact command run
- The key outputs (PR URL, issue number, failing check name)
- Any errors verbatim

## Troubleshooting

- If `gh` says not authenticated: run `gh auth login` (interactive) or fix token/SSH setup.
- If acting on the wrong repo: run `gh repo set-default OWNER/REPO` or pass `--repo OWNER/REPO`.
- If JSON queries fail: verify `--json` fields and `--jq` expression.

