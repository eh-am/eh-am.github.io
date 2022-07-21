#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test "runs fine since 'pipefail' is not set" { 
	run "$BATS_TEST_DIRNAME/pipefail_first.sh"
	[ "$status" -eq 0 ]
	[ "${lines[2]}" == 'Hello' ]
}

@test "fails since 'pipefail' is set" {
	run "$BATS_TEST_DIRNAME/pipefail_first_correct.sh"
	[ "$status" -ne 0 ]
	[ "${lines[2]}" != 'Hello' ]
}
