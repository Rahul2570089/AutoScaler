package aws

import (
	"strings"
)

type ErrorType string

const (
	ErrorTypeCapacity  ErrorType = "capacity"
	ErrorTypeThrottle  ErrorType = "throttle"
	ErrorTypeLimit     ErrorType = "limit"
	ErrorTypeFatal     ErrorType = "fatal"
	ErrorTypeTransient ErrorType = "transient"
)

func ClassifyError(err error) ErrorType {
	if err == nil {
		return ErrorTypeTransient
	}

	if strings.Contains(err.Error(), "InsufficientInstanceCapacity") {
		return ErrorTypeCapacity
	}
	if strings.Contains(err.Error(), "RequestLimitExceeded") {
		return ErrorTypeThrottle
	}
	if strings.Contains(err.Error(), "InstanceLimitExceeded") || strings.Contains(err.Error(), "VcpuLimitExceeded") {
		return ErrorTypeLimit
	}
	if strings.Contains(err.Error(), "InvalidParameter") || strings.Contains(err.Error(), "UnauthorizedOperation") {
		return ErrorTypeFatal
	}

	return ErrorTypeTransient
}
