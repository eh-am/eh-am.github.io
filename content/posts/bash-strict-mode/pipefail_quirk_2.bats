#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test "runs fine even though file does not exist" {
	run "$BATS_TEST_DIRNAME/pipefail_quirk_2.sh"
	[ "$status" -ne 0 ]
	[ "${lines[2]}" == "Whitelist: ''" ]
}

@test "fails since we verify file presence" {
	run "$BATS_TEST_DIRNAME/pipefail_quirk_2_correct.sh"
	[ "$status" -eq 1 ]
	[ "$output" == "Whitelist file does not exist" ]
}
