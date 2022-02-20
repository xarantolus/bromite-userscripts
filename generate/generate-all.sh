#!/usr/bin/env bash 
set -euo pipefail 

cd idcac
bash ./generate.sh
cd ..


git config --global user.name 'github-actions'
git config --global user.email '41898282+github-actions[bot]@users.noreply.github.com'

# Add regenerated script files
git add ../block

# This fails if there are no changes
git commit -m"Automatically update script(s)" && git push || true
