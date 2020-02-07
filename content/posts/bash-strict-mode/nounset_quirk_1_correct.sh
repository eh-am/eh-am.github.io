#!/usr/bin/env bash

set -euo pipefail

my_fn() {
	echo "Received args: '$@'"
}

if [ $# -eq 0 ]; then
	echo "No arguments supplied"
	exit 1
else
	my_fn "$@"
fi
