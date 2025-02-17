#!/usr/bin/env bash

# Run tests in a pseudo-terminal to prevent bubbletea from panicking
output=$(script -q -c "go test -v ./...; echo \$?" /dev/null)
echo "$output" | head -n -1
exit $(echo "$output" | tail -n 1 | tr -d '\r')
