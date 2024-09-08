#!/bin/sh

echo "Setting up Git hooks..."

# Ensure .git/hooks exists
mkdir -p .git/hooks

# Link the pre-commit hook
ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit

# Link the commit-msg hook
ln -sf ../../.githooks/commit-msg .git/hooks/commit-msg

# Make the hooks executable
chmod +x .githooks/commit-msg
chmod +x .githooks/pre-commit

echo "Git hooks setup completed."

