#!/usr/bin/env bash

set -euo pipefail

my_fn() {
	echo "Received args: '$@'"
}

my_fn "$@"
