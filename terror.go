package terror

// README:
// Used with error wrapping and eventually be caught and unwrapped.
// Specific use for area that has no context error wrapping

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
)

// baseVersion actual version of terror module
const baseVersion string = "v2.0.6"

// ErrKind Kind of error
type ErrKind string

// ErrKindSystem any kind of error that is not caused by the input
const ErrKindSystem ErrKind = "system"

// ErrKindInput error caused by failing sanity check or bad or invalid input
const ErrKindInput ErrKind = "input"

// ErrMessageGeneric basic message to be used for the error
const ErrMessageGeneric string = "program error occured, please contact admin if error continues"

// MaxDepth maximum depth the error/panic will unwrap or traverse, can change to suit needs
var MaxDepth int = 20

// AppVersion holds the version of the caller app
var AppVersion string = "v0.0.0"

// ErrLevel type
type ErrLevel int

// Various error level
const (
	ErrLevelWarn ErrLevel = iota + 1
	ErrLevelError
	ErrLevelPanic
)

// funcExec is type executing function for certain level
type funcExec func(Meta, error)

// various function to run for certain error level
var (
	funcDoWarn  *funcExec // error that is not require action
	funcDoError *funcExec // error that require action
	funcDoPanic *funcExec // error that require action but a panic
)

// Meta data type
type Meta map[string]string

// TError is the custom error type
type TError struct {
	Level    ErrLevel // error level
	File     string   // which file caused error
	FuncName string   // which function caused error
	Line     int      // which line caused error
	Message  string   // friendly message to the user or error log storage
	Err      error    // actual that is refered to
	ErrKind  ErrKind  // kind of error
	Meta     Meta     // any additional information that is useful in debugging error (backend only, do not expose this to user. flat, so each meta layer could overwrite accidentally)
}

// SetVersion so caller can set the correct version, and echo correct version
func SetVersion(v string) {
	AppVersion = v
}

// Error mimic golang errors.Error
func (e *TError) Error() string {
	return e.Err.Error()
}

// Unwrap the underlying error
func (e *TError) Unwrap() error {
	return e.Err
}

// KVs sets the key-value for the error to add extra level of log
func (e *TError) KVs(kvs ...string) *TError {
	if e.Meta == nil {
		e.Meta = Meta{}
	}

	if len(kvs)%2 == 0 {
		prev := ""
		for i, val := range kvs {
			if i%2 == 0 {
				e.Meta[val] = ""
			} else {
				e.Meta[prev] = val
			}
			prev = val
		}
	} else {
		e.Meta["kvNotEven"] = "Number of KVs not even"
		log.Println("ERROR: Number of KVs not even")
	}
	return e
}

// new constructor for Error with full parameters support
func new(err error, file, funcName string, line int, message string, errKind ErrKind, errLevel ErrLevel, kvs ...string) *TError {
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
		message = ErrMessageGeneric
	}

	return &TError{
		Level:    errLevel,
		File:     file,
		FuncName: funcName,
		Line:     line,
		Message:  message,
		Err:      err,
		ErrKind:  errKind,
		Meta:     meta,
	}
}

// SetCallbackWarn set callback function when .Warn() called
func SetCallbackWarn(callback funcExec) {
	funcDoWarn = &callback
}

// SetCallbackError set callback function when .Error() called
func SetCallbackError(callback funcExec) {
	funcDoError = &callback
}

// SetCallbackPanic set callback function when .Panic() called
func SetCallbackPanic(callback funcExec) {
	funcDoPanic = &callback
}

// Error returns TError with level error
func Error(err error, friendlyMessage ...string) *TError {
	// deals with accidental error == nil
	if err == nil {
		err = fmt.Errorf("Error is nil (Error)")
	}

	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	// if friendly message passed, set message
	// if not, use child friendly message
	msg := err.Error()
	if len(friendlyMessage) > 0 {
		msg = strings.Join(friendlyMessage, ". ")

	} else {
		var bErr *TError
		if errors.As(err, &bErr) {
			msg = bErr.Message
		}
	}

	return new(err, file, funcName, line, msg, ErrKindSystem, ErrLevelError)
}

// Panic returns a new TError with level Panic
func Panic(err error, friendlyMessage ...string) *TError {
	// deals with accidental error == nil
	if err == nil {
		err = fmt.Errorf("Error is nil (Panic)")
	}

	// if friendly message passed, set message
	// if not, use child friendly message
	msg := err.Error()
	if len(friendlyMessage) > 0 {
		msg = strings.Join(friendlyMessage, ". ")

	} else {
		var bErr *TError
		if errors.As(err, &bErr) {
			msg = bErr.Message
		}
	}

	return &TError{
		Level:   ErrLevelPanic,
		Err:     err,
		Message: msg,
	}
}

// Warn returns a new TError with level Warn
func Warn(err error, friendlyMessage ...string) *TError {
	// deals with accidental error == nil
	if err == nil {
		err = fmt.Errorf("Error is nil (Warn)")
	}

	// if friendly message passed, set message
	// if not, use child friendly message
	msg := err.Error()
	if len(friendlyMessage) > 0 {
		msg = strings.Join(friendlyMessage, ". ")

	} else {
		var bErr *TError
		if errors.As(err, &bErr) {
			msg = bErr.Message
		}
	}

	return &TError{
		Level:   ErrLevelWarn,
		Err:     err,
		Message: msg,
	}
}

// GetLevel will find out if error is TError and their level. 0 is not a TError
func GetLevel(err error) ErrLevel {
	e, ok := err.(*TError)
	if !ok {
		return 0
	}
	return e.Level
}

// Echo will walk through error stack and echo output to the screen
func Echo(err error, noEchos ...bool) string {
	if err == nil {
		return "Error is nil (Echo)"
	}
	var noEcho bool
	if len(noEchos) > 0 && noEchos[0] {
		noEcho = true
	}

	level := 0

	i := 0
	j := 0
	var xErr *TError
	errLines := []string{}
	verrLines := []string{}
	metaData := Meta{}
	g := err
	for {
		if errors.As(g, &xErr) {
			level = int(xErr.Level)
			for k, v := range xErr.Meta {
				metaData[k] = v
			}
			if xErr.Level == ErrLevelPanic {
				return echoPanic(xErr, metaData, noEcho)
			}
			i++
			errLines = append(errLines, fmt.Sprintf("  %d > \033[1;34m%s\033[0m[%s:%d] %v", i, xErr.FuncName, xErr.File, xErr.Line, xErr))
			verrLines = append(verrLines, fmt.Sprintf("  %d > %s[%s:%d] %v", i, xErr.FuncName, xErr.File, xErr.Line, xErr))
			g = xErr.Unwrap()

		} else if xe, ok := g.(*TError); ok {
			level = int(xe.Level)
			for k, v := range xErr.Meta {
				metaData[k] = v
			}
			if xErr.Level == ErrLevelPanic {
				return echoPanic(xErr, metaData, noEcho)
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
		if j > MaxDepth {
			errLines = append(errLines, fmt.Sprintf("stop >%d deep", MaxDepth))
			verrLines = append(verrLines, fmt.Sprintf("stop >%d deep", MaxDepth))
			break
		}
	}

	// reverse
	errLines = StringSliceReverse(errLines)
	verrLines = StringSliceReverse(verrLines)

	out := ""
	if level == int(ErrLevelWarn) {
		// Yellow WARN text
		out = fmt.Sprintf("\033[1;33mWARN\033[0m ver: %s  \n%+v", AppVersion, strings.Join(errLines, "\n"))
	} else {
		// Red ERROR text
		out = fmt.Sprintf("\033[1;31mERROR\033[0m ver: %s  \n%+v", AppVersion, strings.Join(errLines, "\n"))
	}
	msg := strings.SplitAfterN(verrLines[0], " > ", 2)[1]
	vout := fmt.Sprintf("ERROR  %s\n%+v", msg, strings.Join(verrLines, "\n"))

	if !noEcho {
		log.Println(out)
	}

	// recover from panic from funcDoWarn, funcDoError
	defer func() {
		if rec := recover(); rec != nil {
			message := "terror funcDoError panicked"
			if level == int(ErrLevelWarn) {
				message = "terror funcDoWarn panicked"
			}
			strStack := string(debug.Stack())

			var err error
			switch v := rec.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("func error")
			}

			log.Printf("%s recovered: %s. %s\n", message, err.Error(), strStack)
		}
	}()

	// execute function
	if level == int(ErrLevelWarn) {
		if funcDoWarn != nil {
			(*funcDoWarn)(metaData, err)
		}
	} else {
		if funcDoError != nil {
			(*funcDoError)(metaData, err)
		}
	}

	return vout
}

// echoPanic will walk through panic stack and echo output to the screen
func echoPanic(err *TError, metaData Meta, noEcho bool) string {
	if err == nil {
		return "Error is nil (echoPanic)"
	}

	lines := []string{}
	vlines := []string{} // vanilla line, without ascii colour
	// need to adjust i depth, depending on the project
	i := 2
	j := 0
	for j < MaxDepth {
		pc, fn, line, _ := runtime.Caller(i)
		if line == 0 {
			break
		}
		lines = append(lines, fmt.Sprintf("  %d > \033[1;34m%s\033[0m[%s:%d]", j, runtime.FuncForPC(pc).Name(), fn, line))
		vlines = append(vlines, fmt.Sprintf("  %d > %s[%s:%d]", j, runtime.FuncForPC(pc).Name(), fn, line))
		i++
		j++
	}
	if j >= MaxDepth {
		lines = append(lines, fmt.Sprintf("  %d > exceeded max depth", j))
		vlines = append(vlines, fmt.Sprintf("  %d > exceeded max depth", j))
	}

	// Red background and White Blinking PANIC text
	// Note: panic has no true origin err source message, so get from err.Err.Error()
	msg := err.Err.Error()
	if msg != err.Message {
		msg += ". " + err.Message
	}
	out := fmt.Sprintf("\033[5;41;37mPANIC\033[0m ver: %s  %s \n%+v", AppVersion, msg, strings.Join(lines, "\n"))
	vout := fmt.Sprintf("PANIC  %s\n%+v", err.Message, strings.Join(vlines, "\n"))

	if !noEcho {
		log.Println(out)
	}

	// recover from panic from funcDoPanic
	defer func() {
		if rec := recover(); rec != nil {
			message := "terror funcDoPanic panick-panick"
			strStack := string(debug.Stack())

			var err error
			switch v := rec.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("funcDoPanic error")
			}

			log.Printf("%s recovered: %s. %s\n", message, err.Error(), strStack)
		}
	}()

	// exec function
	if funcDoPanic != nil {
		(*funcDoPanic)(metaData, err)
	}

	return vout
}
