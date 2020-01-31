#!/usr/bin/env bash

IFS=$'\t'
manifest='
apiVersion: v1
kind: ConfigMap
metadata:
  name: whitelist
data:
  whitelist:
  	hardcoded_ip
	$REPLACE_ME
'

generate_whitelist > whitelist

contents="$(cat whitelist)"
replaced="$(echo $manifest | sed "s@\$REPLACE_ME@$contents@g")"

echo $replaced | kubectl apply --dry-run -o yaml -f -
