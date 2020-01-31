#!/usr/bin/env bash

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

find "$DIR/../content" -type f -name "*.bats" -exec \
	npm run bats --prefix "$DIR/.." {} +

