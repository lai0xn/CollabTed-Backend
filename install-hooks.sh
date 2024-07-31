#!/bin/sh

echo "Setting up Git hooks..."

# Ensure .git/hooks exists
mkdir -p .git/hooks

# Link the pre-commit hook
ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit

echo "Git hooks setup completed."
