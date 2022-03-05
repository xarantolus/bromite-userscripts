#!/usr/bin/env bash 
set -euo pipefail 

echo "Installing dependencies..."

apt-get install -y unzip wget jq || true

go install "github.com/xarantolus/jsonextract/cmd/jsonx"

echo "Downloading Firefox extension"

XPI_URL="$(jsonx https://addons.mozilla.org/en/firefox/addon/i-dont-care-about-cookies/ id created hash permissions url | jq -r .url | grep .xpi)"

wget --no-check-certificate -O "extension.zip" "$XPI_URL"

rm -rf "extension" || true

mkdir -p "extension"
unzip -ou "extension.zip" -d "extension"

SCRIPT_PATH="../../block/idcac.user.js"

go run main.go -base extension -output "$SCRIPT_PATH"

rm -rf "extension" "extension.zip" || true
