package api_errors

import (
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApiError struct {
	Domain   string            `json:"domain"`
	Reason   string            `json:"reason"`
	Metadata map[string]string `json:"metadata"`
	s        *status.Status
}

// New returns an error object for the code, message.
func New(code codes.Code, domain, reason, message string) *ApiError {
	return &ApiError{
		s:      status.New(code, message),
		Domain: domain,
		Reason: reason,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code codes.Code, domain, reason, format string, a ...interface{}) *ApiError {
	return New(code, domain, reason, fmt.Sprintf(format, a...))
}

// GRPCStatus returns the Status represented by se.
func (e *ApiError) GRPCStatus() *status.Status {
	s, err := e.s.WithDetails(&errdetails.ErrorInfo{
		Domain:   e.Domain,
		Reason:   e.Reason,
		Metadata: e.Metadata,
	})
	if err != nil {
		return e.s
	}
	return s
}

func (e ApiError) Error() string {
	return fmt.Sprintf("error: domain = %s reason = %s metadata = %v", e.Domain, e.Reason, e.Metadata)
}

func (e ApiError) Code() uint32 {
	return uint32(e.s.Code())
}

func (e ApiError) Data() interface{} {
	return e.Metadata
}

func (e ApiError) Message() string {
	return e.Error()
}

// Code returns the code for a particular error.
// It supports wrapped errors.
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if se := FromError(err); err != nil {
		return se.s.Code()
	}
	return codes.Unknown
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *ApiError) WithMetadata(md map[string]string) *ApiError {
	err := *e
	err.Metadata = md
	return &err
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *ApiError {
	if err == nil {
		return nil
	}
	if target := new(ApiError); errors.As(err, &target) {
		return target
	}
	gs, ok := status.FromError(err)
	if ok {
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				return New(
					gs.Code(),
					d.Domain,
					d.Reason,
					gs.Message(),
				).WithMetadata(d.Metadata)
			}
		}
	}
	return New(gs.Code(), "", "", err.Error())
}

// Is matches each error in the chain with the target value.
func (e *ApiError) Is(err error) bool {
	if target := new(ApiError); errors.As(err, &target) {
		return target.Domain == e.Domain && target.Reason == e.Reason
	}
	return false
}
