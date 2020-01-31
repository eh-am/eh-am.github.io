#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test "prints hello" {
	run "$BATS_TEST_DIRNAME/errexit_or.sh"
	[ "$status" -eq 0 ]
	[ "${lines[1]}" == 'Hello' ]
}
