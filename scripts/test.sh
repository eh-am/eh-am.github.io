#!/usr/bin/env bash

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

test_bash() {
	# Bats need the files to have execute permission
	find "$DIR/../content" -type f -name "*.sh" -exec \
		chmod +x {} +

	find "$DIR/../content" -type f -name "*.bats" -exec \
		 npm run bats --prefix "$DIR/.." {} +
}

test_bash
