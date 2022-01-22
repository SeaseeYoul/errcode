package errcode_test

import (
	"errcode"
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	defer func() {
		errStr := recover()
		errcode.ExpectEQ(t, "ecode: -1 already exist", errStr, "New duplicate ecode should cause panic")
	}()
	var _ error = errcode.New(-1)
	var _ error = errcode.New(-2)
	var _ error = errcode.New(-1)
}

func TestNew2(t *testing.T) {
	defer func() {
		errStr := recover()
		errcode.ExpectEQ(t, "business ecode must less than zero", errStr, "New should check business ecode must greater than zero")
	}()
	var _ error = errcode.New(5555)
}

func TestErrMessage(t *testing.T) {
	e1 := errcode.New(-65535)
	errcode.ExpectEQ(t, "-65535", e1.Error(), "ecode message should be `-65535`")
	errcode.ExpectEQ(t, "-65535", e1.Message(), "unregistered ecode message should be ecode number")

	errcode.RegisterMessages(map[int]string{-65535: "testErr"})
	errcode.ExpectEQ(t, "testErr", e1.Message(), "registered ecode message should be `testErr`")
}

func TestCause(t *testing.T) {
	e1 := errcode.New(-65534)
	var err error = e1
	e2 := errcode.Cause(err)

	errcode.ExpectEQ(t, -65534, e2.Code())
	errcode.ExpectEQ(t, "-65534", e2.Message())

	e3 := errcode.Cause(nil)
	errcode.ExpectEQ(t, e3.Code(), errcode.OK.Code())
	errcode.ExpectEQ(t, e3.Message(), errcode.OK.Message())

	e4 := errcode.Cause(errors.New("123"))
	errcode.ExpectEQ(t, 123, e4.Code())
	errcode.ExpectEQ(t, "123", e4.Message())
}

func TestInt(t *testing.T) {
	e1 := errcode.Int(123456)
	errcode.ExpectEQ(t, 123456, e1.Code())
	errcode.ExpectEQ(t, "123456", e1.Error())
	errcode.ExpectEQ(t, "123456", e1.Message())
	errcode.ExpectEQ(t, 0, len(e1.Details()))
}

func TestString(t *testing.T) {
	eStr := errcode.String("123")
	errcode.ExpectEQ(t, 123, eStr.Code())
	errcode.ExpectEQ(t, "123", eStr.Message())

	eStr = errcode.String("test")
	t.Logf("code: %v vs. error code: %v", eStr.Code(), errcode.InvalidArgument.Code())
	errcode.ExpectEQ(t, eStr.Code(), errcode.InvalidArgument.Code())
	errcode.ExpectEQ(t, eStr.Error(), errcode.InvalidArgument.Error())
	errcode.ExpectEQ(t, eStr.Message(), errcode.InvalidArgument.Message())

	eStr = errcode.String("")
	t.Logf("code: %v vs. error code: %v", eStr.Code(), errcode.OK.Code())
	errcode.ExpectEQ(t, eStr.Code(), errcode.OK.Code())
	errcode.ExpectEQ(t, eStr.Error(), errcode.OK.Error())
	errcode.ExpectEQ(t, eStr.Message(), errcode.OK.Message())

}

func TestEqualError(t *testing.T) {
	errcode.ExpectTrue(t, errcode.EqualError(errcode.OK, errcode.OK))
	errcode.ExpectFalse(t, errcode.EqualError(errcode.OK, errcode.Cancelled))
}

func TestCheckOk(t *testing.T) {
	var status *errcode.Status = nil
	errcode.ExpectTrue(t, errcode.CheckOk(status), "nil pointer status should return as OK")
	var codes errcode.Codes = status
	errcode.ExpectTrue(t, errcode.CheckOk(codes), "nil pointer codes interface should return as OK")
}

func TestMustCheckOk(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Logf("error: %v", err)
		}
	}()
	errcode.Errorf(errcode.Internal, "test error")
}
