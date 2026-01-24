#!/bin/bash
# Source this script to add the bin directory to your PATH
# Usage: source setup_env.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$SCRIPT_DIR/bin"

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    export PATH="$BIN_DIR:$PATH"
    echo "Added $BIN_DIR to PATH"
else
    echo "$BIN_DIR is already in PATH"
fi
