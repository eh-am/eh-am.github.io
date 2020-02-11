#!/usr/bin/env bash

set -eo pipefail

non_existent_cmd | another_non_existent_cmd | cat
echo "Hello"
