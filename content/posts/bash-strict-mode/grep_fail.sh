#!/usr/bin/env bash

set -e
echo xd

status_code=$(grep non_existant_word /dev/null)
echo "Hello world ${status_code}"
