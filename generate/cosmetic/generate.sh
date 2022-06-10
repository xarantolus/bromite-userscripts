#!/usr/bin/env bash
set -euo pipefail

SCRIPT_PATH="../../block/cosmetic.user.js"
LITE_SCRIPT_PATH="../../block/cosmetic-lite.user.js"

# Download top 1M domains
wget -q "http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip" -O "top1m.zip"
unzip -o "top1m.zip" -d "top1m"

# Run the normal generator
go run main.go -input "filter-lists.txt" -output "$SCRIPT_PATH"

# Run the generator with top domains
go run main.go -input "filter-lists.txt" -top top1m/*.csv -output "$LITE_SCRIPT_PATH"
