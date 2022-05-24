#!/usr/bin/env bash
set -euo pipefail

PREV_DIR="$(pwd)"

cd "$(mktemp -d)"
git clone "https://github.com/AdguardTeam/Scriptlets.git"
cd Scriptlets
npm install
npm run build
CORELIB_PATH="$(pwd)/dist/scriptlets.corelibs.json"
cd "$PREV_DIR"

SCRIPT_PATH="../../block/cosmetic.user.js"

go run main.go -input "filter-lists.txt" -scriptlets "$CORELIB_PATH" -output "$SCRIPT_PATH"
