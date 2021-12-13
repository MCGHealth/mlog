package mlog

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"log"
)

var (
	mtx      sync.RWMutex
	ilog     *log.Logger
	loglevel LogEventLevel
	isInit   bool = false
)

func reset() {
	isInit = false
	ilog = nil
}

// Intialize setups up the logger with the writer to be used
// to write out the log entries.
// Examples:
// * use the stdout within a container:
func Initialize(w io.Writer, level LogEventLevel) error {
	mtx.Lock()
	defer mtx.Unlock()

	if isInit {
		return nil
	}

	if w == nil {
		return errors.New("an `io.writer` cannot be nil")
	}

	if level == UnknownLevel {
		return errors.New("log level is invalid")
	}

	loglevel = level
	ilog = log.New(w, "golog ", log.Ldate|log.Ltime|log.LUTC)
	Infof("internal logging set to level %s", loglevel)
	isInit = true
	return nil
}

// CurrentLevel returns the current logging level.
func CurrentLevel() LogEventLevel {
	mtx.RLock()
	defer mtx.RUnlock()
	return loglevel
}

// SetLogLevel allows the logging level to be updated after
// initialization.
func SetLogLevel(newLevel LogEventLevel) error {
	if newLevel == UnknownLevel {
		return errors.New("log level `UnknownLevel` is invalid")
	}
	mtx.Lock()
	defer mtx.Unlock()
	loglevel = newLevel
	return nil
}

// Debug writes out an internal debug log message.
func Debug(msg string) {
	if loglevel > DebugLevel {
		return
	}
	ilog.Printf("DEBUG: %v", msg)
}

// Info writes out an internal info log message.
func Info(msg string) {
	if loglevel > InfoLevel {
		return
	}
	ilog.Printf("INFO: %v", msg)
}

// Warn writes out an internal warning log message.
func Warn(msg string) {
	if loglevel > WarnLevel {
		return
	}
	ilog.Printf("WARN: %v", msg)
}

// Error writes out an internal error log message.
func Error(err error, msg *string) {
	if msg != nil {
		ilog.Printf("ERROR: %v - %v", err, *msg)
	} else {
		ilog.Printf("ERROR: %v", err)
	}
}

// Critical writes out a formatted internal critica log message.
func Critical(err error, msg *string) {
	if msg != nil {
		ilog.Printf("CRITICAL: %v - %v", err, *msg)
	} else {
		ilog.Printf("CRITICAL: %v", err)
	}
}

// Debugf writes out a formatted internal debug log message.
func Debugf(msg string, args ...interface{}) {
	if loglevel > DebugLevel {
		return
	}
	fmsg := fmt.Sprintf(msg, args...)
	ilog.Printf("DEBUG: %v", fmsg)
}

// Infof writes out a formatted internal info log message.
func Infof(msg string, args ...interface{}) {
	if loglevel > InfoLevel {
		return
	}
	fmsg := fmt.Sprintf(msg, args...)
	ilog.Printf("INFO: %v", fmsg)
}

// Warnf writes out a formatted internal warn log message.
func Warnf(msg string, args ...interface{}) {
	if loglevel > WarnLevel {
		return
	}
	fmsg := fmt.Sprintf(msg, args...)
	ilog.Printf("WARN: %v", fmsg)
}

// Errorf writes out a formatted internal error log message.
func Errorf(err error, msg string, args ...interface{}) {
	fmsg := fmt.Sprintf(msg, args...)
	ilog.Printf("ERROR: - %s; error: %v", fmsg, err)
}

// Criticalf writes out a formatted internal critical log message.
func Criticalf(err error, msg string, args ...interface{}) {
	fmsg := fmt.Sprintf(msg, args...)
	ilog.Printf("CRITICAL: - %s; error: %v", fmsg, err)
}
