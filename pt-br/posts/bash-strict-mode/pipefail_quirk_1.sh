#!/usr/bin/env bash

set -eo pipefail

cat non_existing_file | xargs curl -qs
