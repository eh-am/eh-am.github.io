#!/usr/bin/env bash

set -e

function random_bytes {
	echo $(head -c "$1" /dev/random | base64)
}

random_bytes 10
random_bytes 50
random_bytes 
random_bytes 5
