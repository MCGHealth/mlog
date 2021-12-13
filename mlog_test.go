package mlog_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/MCGhealth/mlog"
	"github.com/stretchr/testify/assert"
)

var iobuffer bytes.Buffer

func TestInit(t *testing.T) {
	type testCase struct {
		c string
		l mlog.LogEventLevel
		w io.Writer
		e bool
		m string
	}
	testcases := []testCase{
		{c: "BAD_LEVEL", l: mlog.UnknownLevel, w: &iobuffer, e: true, m: "log level is invalid"},
		{c: mlog.DEBUG, l: mlog.DebugLevel, w: nil, e: true, m: "arg `w` cannot be nil"},
		{c: mlog.DEBUG, l: mlog.DebugLevel, w: &iobuffer, e: false},
		{c: mlog.INFO, l: mlog.InfoLevel, w: &iobuffer, e: false},
		{c: mlog.WARN, l: mlog.WarnLevel, w: &iobuffer, e: false},
		{c: mlog.ERROR, l: mlog.ErrorLevel, w: &iobuffer, e: false},
		{c: mlog.CRITICAL, l: mlog.CriticalLevel, w: &iobuffer, e: false},
	}

	for _, tc := range testcases {
		mlog.Reset()
		os.Setenv(mlog.MLOG_LOG_LEVEL, tc.c)
		err := mlog.Initialize(tc.w, tc.l)
		if tc.e {
			assert.Errorf(t, err, tc.m)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, mlog.CurrentLevel(), tc.l)
		}
	}
}

func TestSetLogLevel(t *testing.T) {
	testcases := map[mlog.LogEventLevel]bool{
		mlog.UnknownLevel:  true,
		mlog.DebugLevel:    false,
		mlog.InfoLevel:     false,
		mlog.WarnLevel:     false,
		mlog.ErrorLevel:    false,
		mlog.CriticalLevel: false,
	}

	err := mlog.Initialize(&iobuffer, mlog.DebugLevel)
	assert.NoError(t, err)

	for lvl, expectErr := range testcases {
		err := mlog.SetLogLevel(lvl)
		if expectErr {
			assert.Errorf(t, err, "log level `UnknownLevel` is invalid")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, mlog.CurrentLevel(), lvl)
		}
	}
}

func TestInternals(t *testing.T) {
	t.Run("Test Debug", testDebug)
	t.Run("Test Info", testInfo)
	t.Run("Test Warn", testWarn)
	t.Run("Test Error", testError)
	t.Run("Test Critical", testCritical)

	t.Run("Test Debugf", testDebugf)
	t.Run("Test Infof", testInfof)
	t.Run("Test Warnf", testWarnf)
	t.Run("Test Errorf", testErrorf)
	t.Run("Test Criticalf", testCriticalf)
}

func testDebug(t *testing.T) {
	mlog.SetLogLevel(mlog.DebugLevel)
	le := "debug log entry"
	mlog.Debug(le)
	assertLogEntry(t, mlog.DebugLevel, &le, nil)
}

func testInfo(t *testing.T) {
	mlog.SetLogLevel(mlog.InfoLevel)
	le := "info log entry"
	mlog.Info(le)
	assertLogEntry(t, mlog.InfoLevel, &le, nil)
}

func testWarn(t *testing.T) {
	mlog.SetLogLevel(mlog.WarnLevel)
	le := "warn log entry"
	mlog.Warn(le)
	assertLogEntry(t, mlog.WarnLevel, &le, nil)
}

func testError(t *testing.T) {
	tests := map[string]string{
		"error 1": "",
		"error 2": "error message",
	}

	for e, m := range tests {
		err := errors.New(e)

		if len(m) == 0 {
			mlog.Error(err, nil)
			assertLogEntry(t, mlog.ErrorLevel, nil, err)
		} else {
			mlog.Error(err, &m)
			assertLogEntry(t, mlog.ErrorLevel, &m, err)
		}
	}
}

func testCritical(t *testing.T) {
	tests := map[string]string{
		"error 1": "",
		"error 2": "critical error message",
	}

	for e, m := range tests {
		err := errors.New(e)

		if len(m) == 0 {
			mlog.Critical(err, nil)
			assertLogEntry(t, mlog.CriticalLevel, nil, err)
		} else {
			mlog.Critical(err, &m)
			assertLogEntry(t, mlog.CriticalLevel, &m, err)
		}
	}
}

func testDebugf(t *testing.T) {
	mlog.SetLogLevel(mlog.DebugLevel)
	mlog.Debugf("debug %s", "test")
	exp := "debug test"
	assertLogEntry(t, mlog.DebugLevel, &exp, nil)
}

func testInfof(t *testing.T) {
	mlog.SetLogLevel(mlog.InfoLevel)
	mlog.Infof("info %s", "test")
	exp := "info test"
	assertLogEntry(t, mlog.InfoLevel, &exp, nil)
}

func testWarnf(t *testing.T) {
	mlog.SetLogLevel(mlog.WarnLevel)
	mlog.Warnf("warn %s", "test")
	exp := "warn test"
	assertLogEntry(t, mlog.WarnLevel, &exp, nil)
}

func testErrorf(t *testing.T) {
	mlog.SetLogLevel(mlog.ErrorLevel)
	err := errors.New("test error")
	mlog.Errorf(err, "error %s", "test")
	exp := "error test"
	assertLogEntry(t, mlog.ErrorLevel, &exp, err)
}

func testCriticalf(t *testing.T) {
	mlog.SetLogLevel(mlog.CriticalLevel)
	err := errors.New("critical error")
	mlog.Criticalf(err, "critical %s", "test")
	exp := "critical test"
	assertLogEntry(t, mlog.CriticalLevel, &exp, err)
}

func assertLogEntry(t *testing.T, level mlog.LogEventLevel, msg *string, err error) {
	entry := iobuffer.String()
	assert.True(t, strings.Contains(entry, "golog"), "the expected text `%s` was not found: entry text: `%s`", "golog", entry)

	var prefix string
	switch level {
	case mlog.DebugLevel:
		prefix = "DEBUG"
	case mlog.InfoLevel:
		prefix = "INFO"
	case mlog.WarnLevel:
		prefix = "WARN"
	case mlog.ErrorLevel:
		prefix = "ERROR"
	case mlog.CriticalLevel:
		prefix = "CRITICAL"
	}

	assert.True(t, strings.Contains(entry, fmt.Sprintf("%s:", prefix)), "the expected text `%s` was not found: entry text: `%s`", prefix, entry)

	if msg != nil {
		assert.True(t, strings.Contains(entry, *msg), "the expected text `%s` was not found: entry text: `%s`", *msg, entry)
	}

	if err != nil {
		assert.True(t, strings.Contains(entry, err.Error()), "the expected text `%s` was not found: entry text: `%s`", err.Error(), entry)
	}

	iobuffer.Reset()
}
