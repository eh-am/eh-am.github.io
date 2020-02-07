#!/usr/bin/env bats

load '../../../node_modules/bats-support/load'
load '../../../node_modules/bats-assert/load'

@test 'runs fine, since unset does not affect "$@"' {
	run "$BATS_TEST_DIRNAME/nounset_quirk_1.sh"
	[ "$status" -eq 0 ]
	[ "$output" == "Received args: ''" ]
}

@test "validates input manually" {
	run "$BATS_TEST_DIRNAME/nounset_quirk_1_correct.sh"
	[ "$status" -eq 1 ]
	[ "$output" == "No arguments supplied" ]
}
