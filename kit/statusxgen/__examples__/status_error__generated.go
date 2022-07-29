// This is a generated source file. DO NOT EDIT
// Version: 0.0.1
// Source: examples/status_error__generated.go
// Date: Jul 26 00:01:39

package examples

import (
	"github.com/iotexproject/Bumblebee/kit/statusx"
)

var _ statusx.Error = (*StatusError)(nil)

func (v StatusError) StatusErr() *statusx.StatusErr {
	return &statusx.StatusErr{
		Key:       v.Key(),
		Code:      v.Code(),
		Msg:       v.Msg(),
		CanBeTalk: v.CanBeTalk(),
	}
}

func (v StatusError) Unwrap() error {
	return v.StatusErr()
}

func (v StatusError) Error() string {
	return v.StatusErr().Error()
}

func (v StatusError) StatusCode() int {
	return statusx.StatusCodeFromCode(int(v))
}

func (v StatusError) Code() int {
	if with, ok := (interface{})(v).(statusx.ServiceCode); ok {
		return with.ServiceCode() + int(v)
	}
	return int(v)

}

func (v StatusError) Key() string {
	switch v {
	case Unauthorized:
		return "Unauthorized"
	case InternalServerError:
		return "InternalServerError"
	}
	return "UNKNOWN"
}

func (v StatusError) Msg() string {
	switch v {
	case Unauthorized:
		return "Unauthorized"
	case InternalServerError:
		return "InternalServerError 内部错误"
	}
	return "-"
}

func (v StatusError) CanBeTalk() bool {
	switch v {
	case Unauthorized:
		return true
	case InternalServerError:
		return false
	}
	return false
}
