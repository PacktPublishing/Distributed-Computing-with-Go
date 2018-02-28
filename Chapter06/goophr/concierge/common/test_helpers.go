package common

import "testing"

const checkMark = "\u2713"
const ballotX = "\u2717"
const prefix = "\t\t - "

// OnTestSuccess is for pass.
func OnTestSuccess(t *testing.T, msg string) {
	t.Log(prefix+msg, checkMark)
}

//OnTestError is for fail.
func OnTestError(t *testing.T, msg string) {
	t.Error(prefix+msg, ballotX)
}

//OnTestUnexpectedError is for unexpected fail.
func OnTestUnexpectedError(t *testing.T, err error) {
	OnTestError(t, "Unexpected Error:\n"+err.Error())
}
