package routine

import "github.com/timandy/routine/bytesconv"

// StackError an error type contains stack info.
type StackError interface {
	// Message data when panic is raised.
	Message() any

	// StackTrace stack when this instance is created.
	StackTrace() string

	// Error contains Message and StackTrace.
	Error() string
}

// NewStackError create a new instance.
func NewStackError(message any) StackError {
	return &stackError{message: message, stackTrace: bytesconv.String(traceStack())}
}
