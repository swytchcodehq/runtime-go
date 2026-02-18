package swytchcode

import "fmt"

// SwytchcodeError is returned when exec fails: spawn error, non-zero exit, or invalid JSON output.
type SwytchcodeError struct {
	Message string
	Cause   any
}

func (e *SwytchcodeError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// IsSwytchcodeError reports whether err is a *SwytchcodeError.
func IsSwytchcodeError(err error) bool {
	_, ok := err.(*SwytchcodeError)
	return ok
}
