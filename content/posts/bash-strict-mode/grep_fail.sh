#!/usr/bin/env bash

set -e

status_code=$(grep non_existant_word /dev/null)
echo "Hello world"
