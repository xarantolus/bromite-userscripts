#!/usr/bin/env bash
set -euo pipefail

SCRIPTS_PATH="../../block"

# Download top 1M domains
wget -q "http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip" -O "top1m.zip"
unzip -o "top1m.zip" -d "top1m"

# Run the normal generator
go run .

# For backwards compatibility, we copy the english scripts

cd cosmetic-outputs
cp -f cosmetic-en.user.js cosmetic.user.js
cp -f cosmetic-en-lite.user.js cosmetic-lite.user.js

mv ./* "../$SCRIPTS_PATH"
