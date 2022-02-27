#!/usr/bin/env bash 
set -euo pipefail 

SCRIPT_PATH="../../block/cosmetic.user.js"

go run main.go -input "filter-lists.txt" -output "$SCRIPT_PATH"
