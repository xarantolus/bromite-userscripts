#!/usr/bin/env bash 
set -euo pipefail 

cd idcac
bash ./generate.sh
cd ..

cd cosmetic
bash ./generate.sh
cd ..
