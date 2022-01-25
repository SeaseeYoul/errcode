package errcode

import (
	"testing"
)

func TestNew(t *testing.T) {
	defer func() {
		errStr := recover()
		ExpectEQ(t, "ecode: -1 already exist", errStr, "New duplicate ecode should cause panic")
	}()
	var _ error = New(-1)
	var _ error = New(-2)
	var _ error = New(-1)
	Error(0, "this is tesr")
}
