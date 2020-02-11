
#!/usr/bin/env bash

set -eo pipefail

contents=$(cat non_existing_file)
curl -qs "$contents"
