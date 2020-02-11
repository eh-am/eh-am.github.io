#!/usr/bin/env bash

set -e

if grep non_existant_word /dev/null; then
	echo "Hello world"
else
	echo "Does not exist"
fi
