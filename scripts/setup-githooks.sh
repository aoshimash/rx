#!/bin/bash
# Setup Git hooks by linking githooks/ to .git/hooks/
# This script sets up pre-commit hooks for all developers

set -e

REPO_ROOT="$(git rev-parse --show-toplevel)"
GITHOOKS_DIR="$REPO_ROOT/githooks"
GIT_HOOKS_DIR="$REPO_ROOT/.git/hooks"

if [ ! -d "$GITHOOKS_DIR" ]; then
    echo "Error: githooks/ directory not found"
    exit 1
fi

echo "Setting up Git hooks..."
echo ""

# Method 1: Use git config core.hooksPath (recommended for Git 2.9+)
if git config core.hooksPath > /dev/null 2>&1; then
    CURRENT_HOOKS_PATH="$(git config core.hooksPath)"
    if [ "$CURRENT_HOOKS_PATH" != "$GITHOOKS_DIR" ]; then
        echo "Updating core.hooksPath to: $GITHOOKS_DIR"
        git config core.hooksPath "$GITHOOKS_DIR"
    else
        echo "✓ core.hooksPath already set to: $GITHOOKS_DIR"
    fi
else
    echo "Setting core.hooksPath to: $GITHOOKS_DIR"
    git config core.hooksPath "$GITHOOKS_DIR"
fi

# Method 2: Create symlinks for each hook (fallback)
# This ensures hooks work even if core.hooksPath is not supported
for hook in "$GITHOOKS_DIR"/*; do
    if [ -f "$hook" ] && [ -x "$hook" ]; then
        hook_name="$(basename "$hook")"
        target="$GIT_HOOKS_DIR/$hook_name"
        
        if [ -L "$target" ] && [ "$(readlink "$target")" = "$hook" ]; then
            echo "✓ Symlink already exists: $hook_name"
        elif [ -f "$target" ] && [ ! -L "$target" ]; then
            echo "⚠ Warning: $hook_name already exists as a regular file"
            echo "  Backing up to ${target}.backup"
            mv "$target" "${target}.backup"
            ln -s "$hook" "$target"
            echo "✓ Created symlink: $hook_name"
        else
            ln -sf "$hook" "$target"
            echo "✓ Created symlink: $hook_name"
        fi
    fi
done

echo ""
echo "✅ Git hooks setup complete!"
echo ""
echo "Pre-commit hook will now run automatically on 'git commit'."
echo "To test, run: git commit --allow-empty -m 'test pre-commit hook'"
