#!/usr/bin/env bash

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

REAL_DIR="$(realpath $DIR/..)"
ls -d $REAL_DIR/**/*.sh | entr "$DIR/test.sh"
