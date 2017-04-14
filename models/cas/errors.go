package cas

import (
	"errors"
	"fmt"
)

// CasErrorCode type declaration
type CasErrorCode int

const (
	// ERROR_CODE_INVALID_REQUEST "not all of the required request parameters were present"
	ERROR_CODE_INVALID_REQUEST CasErrorCode = 1 + iota
	// ERROR_CODE_INVALID_TICKET_SPEC "failure to meet the requirements of validation specification"
	ERROR_CODE_INVALID_TICKET_SPEC
	// ERROR_CODE_INVALID_TICKET "the ticket provided was not valid, or the ticket did not come from an initial login and renew was set on validation."
	ERROR_CODE_INVALID_TICKET
	// INVALID_SERVICE "the ticket provided was valid, but the service specified did not match the service associated with the ticket."
	ERROR_CODE_INVALID_SERVICE
	// INVALID_USERNAME "the uid (username) was invalid or removed"
	ERROR_CODE_INVALID_USERNAME
	// ERROR_CODE_INTERNAL_ERROR "an internal error occurred during ticket validation"
	ERROR_CODE_INTERNAL_ERROR
	// ERROR_CODE_UNAUTHORIZED_SERVICE_PROXY "the service is not authorized to perform proxy authentication"
	ERROR_CODE_UNAUTHORIZED_SERVICE_PROXY
	// ERROR_CODE_INVALID_PROXY_CALLBACK "The proxy callback specified is invalid. The credentials specified for proxy authentication do not meet the security requirements"
	ERROR_CODE_INVALID_PROXY_CALLBACK
)

// CasErrorCodes contains all error codes in string format
var CasErrorCodes = [...]string{
	"INVALID_REQUEST",
	"INVALID_TICKET_SPEC",
	"INVALID_TICKET",
	"INVALID_SERVICE",
	"INVALID_USERNAME",
	"INTERNAL_ERROR",
	"UNAUTHORIZED_SERVICE_PROXY",
	"INVALID_PROXY_CALLBACK",
}

func (casErrorCode CasErrorCode) String() string {
	return CasErrorCodes[casErrorCode-1]
}

// CasError contains CAS error information
type CasError struct {
	InnerError error
	Code       CasErrorCode
}

func (c *CasError) Error() string {
	return fmt.Sprintf("%s: %s", c.InnerError, c.Code)
}

func NewCasError(message string, code CasErrorCode) *CasError {
	return &CasError{
		InnerError: errors.New(message),
		Code:       code,
	}
}
