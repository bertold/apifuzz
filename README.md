# API Fuzzing Examples

This repo created for the [LASCON '2023](https://www.lascon.org)
conference talk shows simple APIs with automated tests including
unit and fuzz tests.

The reverse service implementation is based on the
[Getting started with fuzzing](https://go.dev/doc/tutorial/fuzz)
tutorial.

## Compile

You can compile the app using `go build`.

## Unit Tests

Execute them using the standard `go test` command.

## Run The Fuzzing Tests

For example, this command runs the reverse server fuzzing
tests for 10 seconds

`go test -v --skip Test -fuzz=FuzzReverseServer --fuzztime 10s`

## Illustrating Buggy Implementation

In [reverseserver.go](reverseserver.go), change the `Reverse`
function to use `Reverse_Buggy`.
The fuzzing tests will quickly find an issue.
