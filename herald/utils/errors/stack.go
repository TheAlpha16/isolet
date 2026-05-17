package errors

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// GetStackTrace returns a human-readable stack trace
func GetStackTrace() string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder
	sb.WriteString("Traceback (most recent call last):\n")

	for {
		frame, more := frames.Next()
		fmt.Fprintf(&sb,
			"\t-> %s:%d (%s)\n",
			frame.File, frame.Line,
			path.Base(frame.Function),
		)
		if !more {
			break
		}
	}
	return sb.String()
}

// RecoverPanic can be deferred to capture and log a panic
func RecoverPanic(msg string) func() {
	return func() {
		if r := recover(); r != nil {
			fmt.Printf("%s PANIC: %v\n%s", msg, r, GetStackTrace())
		}
	}
}
