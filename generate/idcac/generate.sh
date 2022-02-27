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

SCRIPT_PATH="../../block/idcac.user.js"

go run main.go -base extension -output "$SCRIPT_PATH"

rm -rf "extension" "extension.zip" || true

# Now if we only changed the lines with version info in them (there are two lines), we reset the file
CHANGED_LINE_COUNT="$(git diff -U0 $SCRIPT_PATH | grep '^[+]' | grep -Ev '^(--- a/|\+\+\+ b/)' | wc -l)"
echo "$CHANGED_LINE_COUNT lines changed"

if [[ $CHANGED_LINE_COUNT -lt 3 ]]; then
    git checkout main "$SCRIPT_PATH"
fi
