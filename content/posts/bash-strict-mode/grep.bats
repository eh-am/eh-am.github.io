#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test "fails since grep returns non 0" {
	run "$BATS_TEST_DIRNAME/grep_fail.sh"
	[ "$status" -ne 0 ]
}

@test "runs fine since grep is in a if statement" {
	run "$BATS_TEST_DIRNAME/grep_correct.sh"
	[ "$status" -eq 0 ]
	[ "$output" == "Does not exist" ]
}
