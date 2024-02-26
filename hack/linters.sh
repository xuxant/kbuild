#!/bin/bash



RED='\033[0;31m'
GREEN='\033[0;32m'
RESET='\033[0m'

echo "Running linters..."

gofmt () {
  files=$(find . -name "*.go" | grep -vE '\.\/fs\/|\.\/vendor\/' | xargs gofmt -l -s)
  if [[ $files ]]; then
    echo "Go Fmt errors in files:"
    echo $files
    diff=$(find . -name "*.go" | grep -vE '\.\/fs\/|\.\/vendor\/' | xargs gofmt -w -d)
    echo $diff
    exit 1
  fi
}

gofmt