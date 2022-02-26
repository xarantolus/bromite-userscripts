#!/usr/bin/env bash 
set -euo pipefail 

echo "Installing dependencies..."

apt-get install -y unzip wget || true

go install mvdan.cc/xurls/v2/cmd/xurls@latest


echo "Downloading Firefox extension"

XPI_URL="$(curl -L https://addons.mozilla.org/en/firefox/addon/i-dont-care-about-cookies/ | xurls | grep downloads | grep .xpi)"

wget --no-check-certificate -O "extension.zip" "$XPI_URL"

rm -rf "extension" || true

mkdir -p "extension"
unzip -ou "extension.zip" -d "extension"

go run main.go -base extension -output "../../block/idcac.user.js"

rm -rf "extension" "extension.zip" || true
