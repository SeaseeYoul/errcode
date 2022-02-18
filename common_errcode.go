package errcode

import (
	"google.golang.org/genproto/googleapis/rpc/code"
	"net/http"
)

// All common errcode
var (
	OK               = addCode(code.Code_OK, http.StatusOK, "")                                // #0
	Cancelled        = addCode(code.Code_CANCELLED, 499, "已取消")                                // #1
	Unknown          = addCode(code.Code_UNKNOWN, http.StatusInternalServerError, "未知错误")      // #2
	InvalidArgument  = addCode(code.Code_INVALID_ARGUMENT, http.StatusBadRequest, "非法输入")      // #3
	DeadlineExceeded = addCode(code.Code_DEADLINE_EXCEEDED, http.StatusGatewayTimeout, "超时错误") // #4
	// Some requested entity (e.g., file or directory) was not found.
	//
	// Note to server developers: if a request is denied for an entire class
	// of users, such as gradual feature rollout or undocumented whitelist,
	// `NOT_FOUND` may be used. If a request is denied for some users within
	// a class of users, such as user-based access control, `PERMISSION_DENIED`
	// must be used.
	//
	// HTTP Mapping: 404 Not Found
	NotFound = addCode(code.Code_NOT_FOUND, http.StatusNotFound, "没找到对象") // #5
	// The entity that a client attempted to create (e.g., file or directory)
	// already exists.
	//
	// HTTP Mapping: 409 Conflict
	AlreadyExists = addCode(code.Code_ALREADY_EXISTS, http.StatusConflict, "已经存在") // #6
	// The caller does not have permission to execute the specified
	// operation. `PERMISSION_DENIED` must not be used for rejections
	// caused by exhausting some resource (use `RESOURCE_EXHAUSTED`
	// instead for those errors). `PERMISSION_DENIED` must not be
	// used if the caller can not be identified (use `UNAUTHENTICATED`
	// instead for those errors). This error code does not imply the
	// request is valid or the requested entity exists or satisfies
	// other pre-conditions.
	//
	// HTTP Mapping: 403 Forbidden
	PermissionDenied = addCode(code.Code_PERMISSION_DENIED, http.StatusForbidden, "权限错误") // #7
	// The request does not have valid authentication credentials for the
	// operation.
	//
	// HTTP Mapping: 401 Unauthorized
	Unauthenticated = addCode(code.Code_UNAUTHENTICATED, http.StatusUnauthorized, "未通过身份验证") // #16
	// Some resource has been exhausted, perhaps a per-user quota, or
	// perhaps the entire file system is out of space.
	//
	// HTTP Mapping: 429 Too Many Requests
	ResourceExhausted = addCode(code.Code_RESOURCE_EXHAUSTED, http.StatusTooManyRequests, "资源耗尽") // #8

	FailedPrecondition = addCode(code.Code_FAILED_PRECONDITION, http.StatusBadRequest, "非预期状态") // #9
	// The operation was aborted, typically due to a concurrency issue such as
	// a sequencer check failure or transaction abort.
	//
	// See the guidelines above for deciding between `FAILED_PRECONDITION`,
	// `ABORTED`, and `UNAVAILABLE`.
	//
	// HTTP Mapping: 409 Conflict
	Aborted = addCode(code.Code_ABORTED, http.StatusConflict, "访问拒绝") // #10
	// The operation was attempted past the valid range.  E.g., seeking or
	// reading past end-of-file.
	//
	// Unlike `INVALID_ARGUMENT`, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate `INVALID_ARGUMENT` if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// `OUT_OF_RANGE` if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between `FAILED_PRECONDITION` and
	// `OUT_OF_RANGE`.  We recommend using `OUT_OF_RANGE` (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an `OUT_OF_RANGE` error to detect when
	// they are done.
	//
	// HTTP Mapping: 400 Bad Request
	OutOfRange = addCode(code.Code_OUT_OF_RANGE, http.StatusBadRequest, "超出范围") // #11
	// The operation is not implemented or is not supported/enabled in this
	// service.
	//
	// HTTP Mapping: 501 Not Implemented
	Unimplemented = addCode(code.Code_UNIMPLEMENTED, http.StatusNotImplemented, "没有实现") // #12
	// Internal errors.  This means that some invariants expected by the
	// underlying system have been broken.  This error code is reserved
	// for serious errors.
	//
	// HTTP Mapping: 500 Internal Server Error
	Internal = addCode(code.Code_INTERNAL, http.StatusInternalServerError, "内部错误") // #13
	// The service is currently unavailable.  This is most likely a
	// transient condition, which can be corrected by retrying with
	// a backoff.
	//
	// See the guidelines above for deciding between `FAILED_PRECONDITION`,
	// `ABORTED`, and `UNAVAILABLE`.
	//
	// HTTP Mapping: 503 Service Unavailable
	Unavailable = addCode(code.Code_UNAVAILABLE, http.StatusServiceUnavailable, "不可用") // #14
	// Unrecoverable data loss or corruption.
	//
	// HTTP Mapping: 500 Internal Server Error
	DataLoss = addCode(code.Code_DATA_LOSS, http.StatusInternalServerError, "数据丢失") // #15

)
