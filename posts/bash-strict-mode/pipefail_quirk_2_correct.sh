#!/usr/bin/env bash

set -eo pipefail

function all_hosts() {
	echo 'host-1
host-2
host-a
host-b'
}


function remove_hosts() {
	hosts=$(all_hosts | tr '\n' ' ')
	whitelist="$1"
	echo "
	Removing hosts:
	$hosts

	Whitelist:
	'$whitelist'
	"

	# Imagine we are passing those two parameters
	# To another command
}

readonly local whitelist_file="non_existent_whitelist_file"

if [ ! -f "$whitelist_file" ]; then
	echo "Whitelist file does not exist"
	exit 1
fi

cat "$whitelist_file" | remove_hosts
