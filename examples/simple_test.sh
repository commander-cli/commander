#!/usr/bin/env bash

@test "heyho" {
    export ENV_TEST="heyho"

    run echo "hello world"

    assert $output "test"
    assert $line[0] "test"
    assert $output ./example.out
    assert $status SUCCESS
    assert $status 127
}