#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test "returns exit code of xargs" {
	run "$BATS_TEST_DIRNAME/pipefail_quirk_1.sh"
	[ "$status" -eq 123 ]
}

@test "returns exit code of cat" {
	run "$BATS_TEST_DIRNAME/pipefail_quirk_1_correct.sh"
	[ "$status" -eq 1 ]
}
