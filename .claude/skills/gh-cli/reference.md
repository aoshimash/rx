# gh CLI Reference

Keep `SKILL.md` short. Put “rare but useful” commands here.

## Use `--repo` explicitly (safest in automation)

- `gh pr view 123 --repo OWNER/REPO`
- `gh issue view 456 --repo OWNER/REPO`

## Templates used by this repo

- PR template file: `.github/pull_request_template.md`
  - **Do NOT** copy and append to the template (causes duplicate sections)
  - **Do** write complete body following the template structure with `cat >"$tmp"`
- Issue template: `.github/ISSUE_TEMPLATE/work-item.md` (name: "Work Item")
  - Recommended: use `--template "Work Item"` for an initial draft, or use `--body-file` for fully scripted creation

## Temp file one-liners

PR body (write complete body, do not copy template):
- `tmp="$(mktemp)"; cat >"$tmp" <<'EOF'
## Summary
- <what>
## Context
- Link issue(s): N/A
- Why this change: <reason>
## Changes
- <change>
## Test plan
- [x] <test>
## Checklist
- [ ] CI is green
- [x] I added/updated tests (or explained why not)
- [x] I updated docs (or explained why not)
- [x] I checked for backwards compatibility and migrations (if applicable)
EOF
gh pr create --title "Add X to Y" --body-file "$tmp"; rm -f "$tmp"`

Comment body:
- `tmp="$(mktemp)"; cat >"$tmp" <<'EOF'
...
EOF
gh pr comment 123 --body-file "$tmp"
rm -f "$tmp"`

## Devcontainer auth / push pitfalls

If `gh auth status` is OK but `git push` fails (credential helper mismatch), do a one-off push without changing git config:
- `git -c credential.helper= -c credential.helper='!gh auth git-credential' push origin HEAD`

## PR data via JSON

Examples (requires `--jq` support; bundled with `gh`):

- PR basics:
  - `gh pr view 123 --json number,title,url,state,author,baseRefName,headRefName`
- Review state:
  - `gh pr view 123 --json reviewDecision,reviews`
- Changed files:
  - `gh pr view 123 --json files --jq '.files[].path'`

## GitHub Actions / workflow runs

Common:
- List runs for the default repo:
  - `gh run list`
- List runs for a specific branch:
  - `gh run list --branch "$(git branch --show-current)"`
- Watch a run until completion:
  - `gh run watch <run-id>`
- View logs (by default opens pager):
  - `gh run view <run-id> --log`

Tip: start from a PR, then locate the related run via checks:
- `gh pr checks 123`

## `gh api` escape hatch (when `gh` subcommands are insufficient)

Guidelines:
- Prefer stable REST endpoints.
- Capture the response and include it in the transcript.

Examples:
- Get PR details (REST):
  - `gh api repos/OWNER/REPO/pulls/123`
- List PR comments:
  - `gh api repos/OWNER/REPO/issues/123/comments`

## Releases

- List releases:
  - `gh release list`
- View a release:
  - `gh release view <tag>`

