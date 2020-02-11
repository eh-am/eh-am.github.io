#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test "runs fine, even though MY_VAR is not set" { 
	run "$BATS_TEST_DIRNAME/nounset.sh"
	[ "$status" -eq 0 ]
	[ "$output" == 'MY_VAR value: ' ]
}

@test "fails since 'nounset' is set" {
	run "$BATS_TEST_DIRNAME/nounset_correct.sh"
	[ "$status" -ne 0 ]
	[[ "$output" =~ 'MY_VAR: unbound variable' ]]

}
