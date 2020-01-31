#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test "runs fine even though file does not exist" {
	run "$BATS_TEST_DIRNAME/errexit.sh"
	[ "$status" -eq 0 ]
}

@test "fail since file does not exist" {
	run "$BATS_TEST_DIRNAME/errexit2.sh"
	[ "$status" -ne 0 ]
	[ "$output" == "cat: /tmp/i_do_not_exist: No such file or directory" ]
}
