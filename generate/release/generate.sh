#!/usr/bin/env bash 
set -euo pipefail 

go run main.go -output "release.md" ../../block/*.user.js
