package errcode

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"runtime"
	"time"
)

func code2Status(code2 Code) *Status {
	return &Status{
		s: &PBStatus{
			Code:    code.Code(code2),
			Message: code2.Message(),
		},
		ctx: context.TODO(),
	}
}

func newError(code Code, message string) *Status {
	st := code2Status(code)
	st.s.Message = message
	return st.withStackEntries(message, 3)
}

// Error new status with code and message
func Error(code Code, message string) *Status {
	return newError(code, message)
}

// Errorf new status with code and message
func Errorf(code Code, format string, args ...interface{}) *Status {
	return newError(code, fmt.Sprintf(format, args...))
}

var _ Codes = &Status{}

// Status statusError is an alias of a status proto
// implement Codes
type Status struct {
	s   *PBStatus
	ctx context.Context
}

func (s *Status) WithContext(ctx context.Context) Codes {
	return &Status{
		s:   s.s,
		ctx: ctx,
	}
}

func (s *Status) WithCancel() (codes Codes, cancel context.CancelFunc) {
	baseCtx := s.Context()
	ctx, cancel := context.WithCancel(baseCtx)
	return s.WithContext(ctx), cancel
}

func (s *Status) WithDeadline(d time.Time) (Codes, context.CancelFunc) {
	baseCtx := s.Context()
	ctx, cancel := context.WithDeadline(baseCtx, d)
	return s.WithContext(ctx), cancel
}

func (s *Status) WithTimeout(timeout time.Duration) (Codes, context.CancelFunc) {
	baseCtx := s.Context()
	ctx, cancel := context.WithTimeout(baseCtx, timeout)
	return s.WithContext(ctx), cancel
}

func (s *Status) WithValue(key, val interface{}) Codes {
	baseCtx := s.Context()
	ctx := context.WithValue(baseCtx, key, val)
	return s.WithContext(ctx)
}

func (s *Status) Context() context.Context {
	if s.ctx == nil {
		return context.TODO()
	}
	return s.ctx
}

// Error implement error
func (s *Status) Error() string {
	if s == nil || s.s == nil {
		return ""
	}
	return s.s.Message
}

// Code return error code
func (s *Status) Code() int {
	if s == nil || s.s == nil {
		return OK.Code()
	}
	return int(s.s.Code)
}

// Message return error message for developer
func (s *Status) Message() string {
	return Code(s.Code()).Message()
}

// Details return error details
func (s *Status) Details() []interface{} {
	if s == nil || s.s == nil {
		return nil
	}
	details := make([]interface{}, 0, len(s.s.Details))
	for _, any := range s.s.Details {
		debugInfo := &errdetails.DebugInfo{}
		if err := any.UnmarshalTo(debugInfo); err != nil {
			fmt.Printf("unmarshal failed: %v", err)
			continue
		}
		details = append(details, debugInfo)
	}
	return details
}

func (s *Status) HttpCode() int {
	return Code(s.Code()).HttpCode()
}

func (s *Status) StackEntries() (details []*errdetails.DebugInfo) {
	if s == nil || s.s == nil {
		return nil
	}
	for _, any := range s.s.Details {
		debugInfo := &errdetails.DebugInfo{}
		if err := any.UnmarshalTo(debugInfo); err != nil {
			fmt.Printf("unmarshal failed: %v", err)
			continue
		}
		details = append(details, debugInfo)
	}
	return details
}

// WithDetails WithDetails
func (s *Status) WithDetails(pbs ...proto.Message) (*Status, error) {
	for _, pb := range pbs {
		anyMsg, err := anypb.New(pb)
		if err != nil {
			return s, err
		}
		s.s.Details = append(s.s.Details, anyMsg)
	}
	return s, nil
}

// Equal for compatible.
func (s *Status) Equal(err error) bool {
	return EqualError(s, err)
}

// Proto return origin protobuf message
func (s *Status) Proto() *PBStatus {
	return s.s
}

// calldepth 表示跳过的代码深度。数字加一表示高一层
func (s *Status) withStackEntries(detail string, calldepth int) *Status {
	var stackEntries []string
	for ; ; calldepth += 1 {
		pc, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			//fmt.Println(fmt.Printf("calldepth: %v, !no", calldepth))
			break
		}
		_func := runtime.FuncForPC(pc)
		function := ""
		if _func != nil {
			function = _func.Name()
		}
		stackEntries = append(stackEntries, fmt.Sprintf("%s:%d %s", file, line, function))
	}

	debugInfo := &errdetails.DebugInfo{
		StackEntries: stackEntries,
		Detail:       detail,
	}
	_, _ = s.WithDetails(debugInfo)
	return s
}

func (s *Status) WithStackEntries(msg string) *Status {
	return s.withStackEntries(msg, 2)
}

func (s *Status) MergeStackEntries(rhs Codes) *Status {
	var buf []proto.Message
	for _, s := range rhs.StackEntries() {
		s.Detail = fmt.Sprintf(":%v|%v|%s", rhs.Code(), rhs.Message(), s.Detail)
		buf = append(buf, s)
	}
	_, _ = s.WithDetails(buf...)
	return s
}

// FromCode create status from ecode
func FromCode(code2 Code) *Status {
	st := &Status{s: &PBStatus{Code: code.Code(code2)}}
	return st.withStackEntries("", 2)
}

// WrapCodes create status from Codes
// 这个是为了对之前的代码做一些 workaround。在之前代码我们使用 *Status 作为
func WrapCodes(codes Codes) *Status {
	if st, ok := codes.(*Status); ok {
		return st
	} else {
		st := &Status{s: &PBStatus{Code: code.Code(codes.Code()), Message: codes.Error()}}
		return st.withStackEntries("", 2)
	}
}

// FromError create status from error
// Pay attention to the difference with Cause()
func FromError(e error, code2 Code) Codes {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	st := &Status{s: &PBStatus{
		Code:    code.Code(code2),
		Message: e.Error(),
	}}
	return st.withStackEntries("", 2)
}

// FromProto new status from grpc detail
func FromProto(pbMsg proto.Message) Codes {
	if msg, ok := pbMsg.(*PBStatus); ok {
		if msg.Message == "" {
			// NOTE: if message is empty convert to pure Code, will get message from config center.
			return Code(msg.Code)
		}
		return &Status{s: msg}
	}
	return newError(Internal, fmt.Sprintf("invalid proto message get %v", pbMsg))
}
