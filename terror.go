package terror

// README:
// Used with error wrapping and eventually be caught and unwrapped.
// Specific use for area that has no context error wrapping

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
)

// baseVersion actual version of terror module
const baseVersion = "v0.0.5"

// ErrKind Kind of error
type ErrKind string

// ErrKindSystem any kind of error that is not caused by the input
const ErrKindSystem ErrKind = "system"

// ErrKindInput error caused by failing sanity check or bad or invalid input
const ErrKindInput ErrKind = "input"

const genericErrorMessage string = "program error occured, please contact admin if error continues"

// maxDepth maximum depth the error/panic will unwrap or traverse
var maxDepth int = 20

// AppVersion holds the version of the caller app
var AppVersion string = "v0.0.0"

// Error is the custom error type
type Error struct {
	IsPanic  bool              // is it a panic
	File     string            // which file caused error
	FuncName string            // which function caused error
	Line     int               // which line caused error
	Message  string            // friendly message to the user or error log storage
	Err      error             // actual that is refered to
	ErrKind  ErrKind           // kind of error
	Meta     map[string]string // any additional information that is useful in debugging error (backend only, do not expose this to user)
}

// SetVersion so caller can set the correct version, and echo correct version
func SetVersion(v string) {
	AppVersion = v
}

// Error mimic golang errors.Error
func (e *Error) Error() string {
	return e.Message
}

// Unwrap the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// new constructor for Error with full parameters support
func new(err error, file, funcName string, line int, message string, errKind ErrKind, kvs ...string) *Error {
	meta := map[string]string{}
	if len(kvs)%2 == 0 {
		prev := ""
		for i, val := range kvs {
			if i%2 == 0 {
				meta[val] = ""
			} else {
				meta[prev] = val
			}
			prev = val
		}
	} else {
		meta["kvNotEven"] = "Number of KVs not even"
		log.Println("ERROR: Number of KVs not even")
	}

	// if friendly message is not included, then it will use err.Error()
	if len(message) == 0 && err != nil {
		message = err.Error()
	} else if len(message) == 0 {
		message = genericErrorMessage
	}

	return &Error{
		File:     file,
		FuncName: funcName,
		Line:     line,
		Message:  message,
		Err:      err,
		ErrKind:  errKind,
		Meta:     meta,
	}
}

// New returns a new Error, zero length message will use generic message
func New(err error, friendlyMessage string, kvs ...string) *Error {
	if err == nil {
		return nil
	}

	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	return new(err, file, funcName, line, friendlyMessage, ErrKindSystem, kvs...)
}

// NewPanic returns a new Error, zero length message will use generic message
func NewPanic(err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		IsPanic: true,
		Err:     err,
		Message: err.Error(),
	}
}

// Echo will walk through error stack and echo output to the screen
func Echo(err error) string {
	if err == nil {
		return "Error is nil (Echo)"
	}

	i := 0
	j := 0
	var xErr *Error
	errLines := []string{}
	verrLines := []string{}
	g := err
	for {
		if errors.As(g, &xErr) {
			if xErr.IsPanic {
				return EchoPanic(xErr)
			}
			i++
			errLines = append(errLines, fmt.Sprintf("  %d > \033[1;34m%s\033[0m[%s:%d] %v", i, xErr.FuncName, xErr.File, xErr.Line, xErr))
			verrLines = append(verrLines, fmt.Sprintf("  %d > %s[%s:%d] %v", i, xErr.FuncName, xErr.File, xErr.Line, xErr))
			g = xErr.Unwrap()

		} else if xe, ok := g.(*Error); ok {
			if xErr.IsPanic {
				return EchoPanic(xErr)
			}
			i++
			errLines = append(errLines, fmt.Sprintf("  %d > \033[1;34m%s\033[0m[%s:%d] %v", i, xe.FuncName, xe.File, xe.Line, xe))
			verrLines = append(verrLines, fmt.Sprintf("  %d > %s[%s:%d] %v", i, xe.FuncName, xe.File, xe.Line, xe))
			g = xe.Unwrap()

		} else {
			i++
			errLines = append(errLines, fmt.Sprintf("  %d > %+v", i, g))
			verrLines = append(verrLines, fmt.Sprintf("  %d > %+v", i, g))
			break
		}

		j++

		// gone too deep, stop
		if j > maxDepth {
			errLines = append(errLines, fmt.Sprintf("stop >%d deep", maxDepth))
			verrLines = append(verrLines, fmt.Sprintf("stop >%d deep", maxDepth))
			break
		}
	}

	// reverse
	errLines = StringSliceReverse(errLines)
	verrLines = StringSliceReverse(verrLines)

	out := fmt.Sprintf("\033[1;31mERROR\033[0m ver: %s  \n%+v", AppVersion, strings.Join(errLines, "\n"))
	msg := strings.SplitAfterN(verrLines[0], " > ", 2)[1]
	vout := fmt.Sprintf("ERROR  %s\n%+v", msg, strings.Join(verrLines, "\n"))

	log.Println(out)
	return vout
}

// EchoPanic will walk through panic stack and echo output to the screen
func EchoPanic(err *Error) string {
	if err == nil {
		return "Error is nil (EchoPanic)"
	}

	lines := []string{}
	vlines := []string{} // vanilla line, without ascii colour
	// need to adjust i depth, depending on the project
	i := 2
	j := 0
	for j < maxDepth {
		pc, fn, line, _ := runtime.Caller(i)
		if line == 0 {
			break
		}
		lines = append(lines, fmt.Sprintf("  %d > \033[1;34m%s\033[0m[%s:%d]", j, runtime.FuncForPC(pc).Name(), fn, line))
		vlines = append(vlines, fmt.Sprintf("  %d > %s[%s:%d]", j, runtime.FuncForPC(pc).Name(), fn, line))
		i++
		j++
	}
	if j >= maxDepth {
		lines = append(lines, fmt.Sprintf("  %d > exceeded max depth", j))
		vlines = append(vlines, fmt.Sprintf("  %d > exceeded max depth", j))
	}

	out := fmt.Sprintf("\033[1;31mPANIC\033[0m ver: %s  %s \n%+v", AppVersion, err.Message, strings.Join(lines, "\n"))
	vout := fmt.Sprintf("PANIC  %s\n%+v", err.Message, strings.Join(vlines, "\n"))

	log.Println(out)
	return vout
}
