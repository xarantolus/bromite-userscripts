#!/usr/bin/env bash 
set -euo pipefail 

SCRIPT_PATH="../../block/cosmetic.user.js"

go run main.go -input "filter-lists.txt" -output "$SCRIPT_PATH"

CHANGED_LINE_COUNT="$(git diff -U0 $SCRIPT_PATH | grep '^[+]' | grep -Ev '^(--- a/|\+\+\+ b/)' | wc -l)"
echo "$CHANGED_LINE_COUNT lines changed"

if [[ $CHANGED_LINE_COUNT -lt 3 ]]; then
    git checkout main "$SCRIPT_PATH"
fi
