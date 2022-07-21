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
Removing hosts: $hosts

Whitelist: '$whitelist'
	"

	# Imagine we are passing those two parameters
	# To another command
}

cat non_existent_whitelist_file | remove_hosts
