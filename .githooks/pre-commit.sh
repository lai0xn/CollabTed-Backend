#!/bin/sh

echo "Running gofmt..."
gofmt -w .
if [ $? -ne 0 ]; then
  echo "gofmt failed. Please fix formatting issues."
  exit 1
fi

echo "Running golint..."
golint ./...
if [ $? -ne 0 ]; then
  echo "golint failed. Please fix linting issues."
  exit 1
fi

echo "Pre-commit checks passed."