package errormsg

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// StatusOK                   = 200 // RFC 7231, 6.3.1
// StatusCreated              = 201 // RFC 7231, 6.3.2
// StatusAccepted             = 202 // RFC 7231, 6.3.3
// StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
// StatusNoContent            = 204 // RFC 7231, 6.3.5
// StatusResetContent         = 205 // RFC 7231, 6.3.6
// StatusPartialContent       = 206 // RFC 7233, 4.1
// StatusMultiStatus          = 207 // RFC 4918, 11.1
// StatusAlreadyReported      = 208 // RFC 5842, 7.1
// StatusIMUsed               = 226 // RFC 3229, 10.4.1
// StatusBadRequest                   = 400 // RFC 7231, 6.5.1
// StatusUnauthorized                 = 401 // RFC 7235, 3.1
// StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
// StatusForbidden                    = 403 // RFC 7231, 6.5.3
// StatusNotFound                     = 404 // RFC 7231, 6.5.4
// StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
// StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
// StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
// StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
// StatusConflict                     = 409 // RFC 7231, 6.5.8
// StatusGone                         = 410 // RFC 7231, 6.5.9
// StatusLengthRequired               = 411 // RFC 7231, 6.5.10
// StatusPreconditionFailed           = 412 // RFC 7232, 4.2
// StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
// StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
// StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
// StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
// StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
// StatusTeapot                       = 418 // RFC 7168, 2.3.3
// StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
// StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
// StatusLocked                       = 423 // RFC 4918, 11.3
// StatusFailedDependency             = 424 // RFC 4918, 11.4
// StatusTooEarly                     = 425 // RFC 8470, 5.2.
// StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
// StatusPreconditionRequired         = 428 // RFC 6585, 3
// StatusTooManyRequests              = 429 // RFC 6585, 4
// StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
// StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

// StatusInternalServerError           = 500 // RFC 7231, 6.6.1
// StatusNotImplemented                = 501 // RFC 7231, 6.6.2
// StatusBadGateway                    = 502 // RFC 7231, 6.6.3
// StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
// StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
// StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
// StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
// StatusInsufficientStorage           = 507 // RFC 4918, 11.5
// StatusLoopDetected                  = 508 // RFC 5842, 7.2
// StatusNotExtended                   = 510 // RFC 2774, 7
// StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6

var errorCodeMessage sync.Map

func validHTTPStatusCode(code int) (ok bool) {
	if code >= 200 && code <= 208 {
		ok = true
	} else if code >= 400 && code <= 418 {
		ok = true
	} else if code >= 421 && code <= 431 {
		ok = true
	} else if code >= 500 && code <= 511 {
		ok = true
	}
	return
}

func getMessage(code int) (msg []string, ok bool) {
	var v interface{}
	if v, ok = errorCodeMessage.Load(code); ok {
		msg, ok = v.([]string)
	}
	return
}

// AddErrorMessage ?????????????????????
func AddErrorMessage(code int, message ...string) {
	errorCodeMessage.Store(code, message)
}

// AddErrorMessages ?????????????????????
func AddErrorMessages(messages map[int][]string) {
	for k, v := range messages {
		errorCodeMessage.Store(k, v)
	}
}

// Message ???????????????????????????????????????
func Message(code int, errs ...error) (msg string) {
	if d, ok := getMessage(code); ok && len(d) > 1 {
		msg = strings.Join(d[1:], ":")
	} else if code != 0 {
		msg = "fail"
	} else {
		msg = "success"
	}
	if len(errs) > 0 && errs[0] != nil {
		msg = msg + "," + errs[0].Error()
	}
	return
}

// HTTPStatus ???????????????HTTP?????????
func HTTPStatus(code int) (status int) {
	if code == 0 {
		status = 200
	} else if d, ok := getMessage(code); ok && len(d) > 0 {
		status, _ = strconv.Atoi(d[0])
		if !validHTTPStatusCode(status) {
			status = 400
		}
	} else {
		status = 400
	}
	return
}

// ErrorResponse error response
type ErrorResponse struct {
	Code    int    `json:"code"`    // ?????????
	Message string `json:"message"` // ????????????
	details []string
}

// NewError init new error response
func NewError(code int) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: Message(code),
	}
}

// Error implement error interface
func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%d:%s", er.Code, er.Message)
}

// StatusCode http status code
func (er *ErrorResponse) StatusCode() int {
	return HTTPStatus(er.Code)
}

// WithDetails add more error info
func (er *ErrorResponse) WithDetails(details ...string) *ErrorResponse {
	newError := *er
	newError.details = make([]string, 0, len(details))
	newError.details = append(newError.details, details...)
	return &newError
}

func init() {
	errorCodeMessage.Store(0, []string{"200", "success"})
}
