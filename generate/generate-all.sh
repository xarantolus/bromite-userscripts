#!/usr/bin/env bash 
set -euo pipefail 

echo "::group::Generating IDCAC script"
cd idcac
bash ./generate.sh
cd ..
echo "::endgroup::"


echo "::group::Generating cosmetic filter script"
cd cosmetic
bash ./generate.sh
cd ..
echo "::endgroup::"
