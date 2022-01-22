package errcode

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	_messages    = map[int]string{}
	_mxMessages  = &sync.RWMutex{}
	_httpCodes   = map[int]int{}
	_mxHttpCodes = &sync.RWMutex{}
	_codes       = map[int]struct{}{} // register codes.
)

func RegisterMessages(cm map[int]string) {
	_mxMessages.Lock()
	defer _mxMessages.Unlock()
	for k, v := range cm {
		_messages[k] = v
	}
}

func RegisterMessage(code int, message string) {
	_mxMessages.Lock()
	defer _mxMessages.Unlock()
	_messages[code] = message
}

func RegisterHttpCode(code int, httpCode int) {
	_mxHttpCodes.Lock()
	defer _mxHttpCodes.Unlock()
	_httpCodes[code] = httpCode
}

// New a errcode.Codes by int value.
// NOTE: errcode must unique in global, the New will check repeat and then panic.
func New(e int) Code {
	if e > 0 {
		panic("business ecode must less than zero")
	}
	return add(e)
}

func add(e int) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Int(e)
}

func RegisterCode(c code.Code, httpCode int, message string) Code {
	if c > 0 {
		panic("business ecode must less than zero")
	}
	return addCode(c, httpCode, message)
}

func addCode(c code.Code, httpCode int, message string) Code {
	i := int(c)
	RegisterMessage(i, message)
	RegisterHttpCode(i, httpCode)
	return add(i)
}

func addInt(c int, httpCode int, message string) Code {
	RegisterMessage(c, message)
	RegisterHttpCode(c, httpCode)
	return add(c)
}

// Codes errcode error interface which has a code & message.
type Codes interface {
	// Error sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	// Details get error detail,it may be nil.
	Details() []interface{}

	HttpCode() int

	// StackEntries return nil if it is a code
	StackEntries() (details []*errdetails.DebugInfo)

	Context() context.Context

	WithContext(ctx context.Context) Codes
	WithCancel() (codes Codes, cancel context.CancelFunc)
	WithDeadline(d time.Time) (Codes, context.CancelFunc)
	WithTimeout(timeout time.Duration) (Codes, context.CancelFunc)
	WithValue(key, val interface{}) Codes
}

// A Code is an int error code spec.
type Code int

func (e Code) toStatus() (codes *Status) {
	return code2Status(e)
}

func (e Code) WithContext(ctx context.Context) Codes {
	return e.toStatus().WithContext(ctx)
}

func (e Code) WithCancel() (codes Codes, cancel context.CancelFunc) {
	return e.toStatus().WithCancel()
}

func (e Code) WithDeadline(d time.Time) (Codes, context.CancelFunc) {
	return e.toStatus().WithDeadline(d)
}

func (e Code) WithTimeout(timeout time.Duration) (Codes, context.CancelFunc) {
	return e.toStatus().WithTimeout(timeout)
}

func (e Code) WithValue(key, val interface{}) Codes {
	return e.toStatus().WithValue(key, val)
}

func (e Code) Context() context.Context {
	return context.TODO()
}

//
func (e Code) Error() string {
	{
		_mxMessages.RLock()
		defer _mxMessages.RUnlock()
		if msg, ok := _messages[e.Code()]; ok {
			return msg
		}
	}
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

// Message return error message
func (e Code) Message() string {
	return e.Error()
}

// Details return details.
func (e Code) Details() []interface{} { return nil }

func (e Code) HttpCode() int {
	_mxHttpCodes.RLock()
	defer _mxHttpCodes.RUnlock()
	if httpCode, ok := _httpCodes[e.Code()]; ok {
		return httpCode
	} else {
		return http.StatusOK
	}
}

func (e Code) StackEntries() (details []*errdetails.DebugInfo) {
	return nil
}

// Int parse code int to error.
func Int(i int) Code { return Code(i) }

// String parse code string to error.
func String(e string) Code {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return InvalidArgument
	}
	return Code(i)
}

// Cause cause from error to ecode.
func Cause(e error) Codes {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return String(e.Error())
}

//safeCode if c == nil, may use OK instead
func safeCode(c Codes) Codes {
	if CheckIsNil(c) {
		return OK
	}
	return c
}

// Equal equal a and b by code int.
func Equal(a, b Codes) bool {
	return safeCode(a).Code() == safeCode(b).Code()
}

// EqualError equal error
func EqualError(code Codes, err error) bool {
	return Cause(err).Code() == code.Code()
}

// CheckOk Check whether to be OK
func CheckOk(c Codes) bool {
	if c == nil {
		return true
	}
	return Equal(c, OK)
}

// CheckError Check whether to be Error
func CheckError(c Codes) bool {
	return !CheckOk(c)
}

// IsOk is an alias of CheckOk
func IsOk(c Codes) bool {
	return CheckOk(c)
}

// IsError is an alias of CheckError
func IsError(c Codes) bool {
	return CheckError(c)
}
