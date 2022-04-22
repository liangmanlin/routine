package routine

import "fmt"

type stackError struct {
	message    any
	stackTrace string
}

func (se *stackError) Message() any {
	return se.message
}

func (se *stackError) StackTrace() string {
	return se.stackTrace
}

func (se *stackError) Error() string {
	s := "StackError"
	if message := fmt.Sprint(se.message); len(message) > 0 {
		s = s + ": " + message
	}
	return s + "\n" + se.stackTrace
}
