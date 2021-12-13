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
	"github.com/stretchr/testify/require"
)

var iobuffer bytes.Buffer

const prefix = "test prefix"

func initlog(t *testing.T, lvl, pfx string) {
	os.Setenv(mlog.MLOG_LOG_LEVEL, lvl)
	os.Setenv(mlog.MLOG_PREFIX, pfx)
	mlog.Reset()
	iobuffer.Reset()
	err := mlog.Initialize(&iobuffer)
	require.NoError(t, err)
}

func TestInit(t *testing.T) {
	type testCase struct {
		pfx    string
		cfg    string
		lvl    mlog.LogEventLevel
		writer io.Writer
		retErr bool
		errmsg string
	}
	testcases := []testCase{
		{pfx: "", cfg: mlog.DEBUG, lvl: mlog.DebugLevel, writer: &iobuffer, retErr: true, errmsg: "env var `MLOG_PREFIX` is missing"},
		{pfx: prefix, cfg: "BAD_LEVEL", lvl: mlog.UnknownLevel, writer: &iobuffer, retErr: true, errmsg: "log level is invalid"},
		{pfx: prefix, cfg: mlog.DEBUG, lvl: mlog.DebugLevel, writer: nil, retErr: true, errmsg: "arg `w` cannot be nil"},
		{pfx: prefix, cfg: mlog.DEBUG, lvl: mlog.DebugLevel, writer: &iobuffer, retErr: false},
		{pfx: prefix, cfg: mlog.INFO, lvl: mlog.InfoLevel, writer: &iobuffer, retErr: false},
		{pfx: prefix, cfg: mlog.WARN, lvl: mlog.WarnLevel, writer: &iobuffer, retErr: false},
		{pfx: prefix, cfg: mlog.ERROR, lvl: mlog.ErrorLevel, writer: &iobuffer, retErr: false},
		{pfx: prefix, cfg: mlog.CRITICAL, lvl: mlog.CriticalLevel, writer: &iobuffer, retErr: false},
	}

	for _, tc := range testcases {
		os.Setenv(mlog.MLOG_LOG_LEVEL, tc.cfg)
		os.Setenv(mlog.MLOG_PREFIX, tc.pfx)
		mlog.Reset()
		iobuffer.Reset()
		err := mlog.Initialize(tc.writer)
		if tc.retErr {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
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

	initlog(t, mlog.DEBUG, "test")
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

func TestDebug(t *testing.T) {
	initlog(t, mlog.DEBUG, prefix)
	le := "debug log entry"
	mlog.Debug(le)
	assertLogEntry(t, mlog.DebugLevel, &le, nil)
}

func TestInfo(t *testing.T) {
	initlog(t, mlog.INFO, prefix)
	le := "info log entry"
	mlog.Info(le)
	assertLogEntry(t, mlog.InfoLevel, &le, nil)
}

func TestWarn(t *testing.T) {
	initlog(t, mlog.WARN, prefix)
	le := "warn log entry"
	mlog.Warn(le)
	assertLogEntry(t, mlog.WarnLevel, &le, nil)
}

func TestError(t *testing.T) {
	initlog(t, mlog.ERROR, prefix)
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

func TestCritical(t *testing.T) {
	initlog(t, mlog.CRITICAL, prefix)
	tests := map[string]string{
		"critical 1": "",
		"critical 2": "critical message",
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

func TestDebugf(t *testing.T) {
	initlog(t, mlog.DEBUG, prefix)
	mlog.Debugf("debug %s", "test")
	exp := "debug test"
	assertLogEntry(t, mlog.DebugLevel, &exp, nil)
}

func TestInfof(t *testing.T) {
	initlog(t, mlog.INFO,prefix)
	mlog.Infof("info %s", "test")
	exp := "info test"
	assertLogEntry(t, mlog.InfoLevel, &exp, nil)
}

func TestWarnf(t *testing.T) {
	initlog(t, mlog.WARN, prefix)
	mlog.Warnf("warn %s", "test")
	exp := "warn test"
	assertLogEntry(t, mlog.WarnLevel, &exp, nil)
}

func TestErrorf(t *testing.T) {
	initlog(t, mlog.ERROR, prefix)
	mlog.SetLogLevel(mlog.ErrorLevel)
	err := errors.New("test error")
	mlog.Errorf(err, "error %s", "test")
	exp := "error test"
	assertLogEntry(t, mlog.ErrorLevel, &exp, err)
}

func TestCriticalf(t *testing.T) {
	initlog(t, mlog.CRITICAL, prefix)
	mlog.SetLogLevel(mlog.CriticalLevel)
	err := errors.New("critical error")
	mlog.Criticalf(err, "critical %s", "test")
	exp := "critical test"
	assertLogEntry(t, mlog.CriticalLevel, &exp, err)
}

func assertLogEntry(t *testing.T, level mlog.LogEventLevel, msg *string, err error) {
	entry := iobuffer.String()
	require.Containsf(t, entry, prefix, "the expected text `%s` was not found: entry text: `%s`", prefix, entry)

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
}

func logDebug() {
	mlog.Debug("Debug entry")
	mlog.Debugf("Debug entry - %d", 1)
}

func logInfo() {
	mlog.Info("Info entry")
	mlog.Infof("Info entry - %d", 2)
}

func logWarn() {
	mlog.Warn("Warn entry")
	mlog.Warnf("Warn entry - %d", 3)
}

func logErr() {
	err := errors.New("test error")
	msg := "Error entry"
	mlog.Error(err, nil)
	mlog.Error(err, &msg)
	mlog.Errorf(err, "Error entry - %d", 4)
}

func logCrit() {
	err := errors.New("test critical")
	msg := "critical entry"
	mlog.Critical(err, nil)
	mlog.Critical(err, &msg)
	mlog.Criticalf(err, "critical entry - %d", 4)
}

func logEntries() {
	logDebug()
	logInfo()
	logWarn()
	logErr()
	logCrit()
}

func TestDebugLoggingLevel(t *testing.T) {
	initlog(t, mlog.DEBUG, "test")
	logEntries()
	le := iobuffer.String()
	assert.Contains(t, le, "DEBUG: ")
	assert.Contains(t, le, "INFO: ")
	assert.Contains(t, le, "WARN: ")
	assert.Contains(t, le, "ERROR: ")
	assert.Contains(t, le, "CRITICAL: ")
}

func TestInfoLoggingLevel(t *testing.T) {
	initlog(t, mlog.INFO, "test")
	logEntries()
	le := iobuffer.String()
	assert.NotContains(t, le, "DEBUG: ")
	assert.Contains(t, le, "INFO: ")
	assert.Contains(t, le, "WARN: ")
	assert.Contains(t, le, "ERROR: ")
	assert.Contains(t, le, "CRITICAL: ")
}

func TestWarnLoggingLevel(t *testing.T) {
	initlog(t, mlog.WARN, "test")
	logEntries()
	le := iobuffer.String()
	assert.NotContains(t, le, "DEBUG: ")
	assert.NotContains(t, le, "INFO: ")
	assert.Contains(t, le, "WARN: ")
	assert.Contains(t, le, "ERROR: ")
	assert.Contains(t, le, "CRITICAL: ")
}

func TestErrorLoggingLevel(t *testing.T) {
	initlog(t, mlog.ERROR, "test")
	logEntries()
	le := iobuffer.String()
	assert.NotContains(t, le, "DEBUG: ")
	assert.NotContains(t, le, "INFO: ")
	assert.NotContains(t, le, "WARN: ")
	assert.Contains(t, le, "ERROR: ")
	assert.Contains(t, le, "CRITICAL: ")
}

func TestCriticalLoggingLevel(t *testing.T) {
	initlog(t, mlog.CRITICAL, "test")
	logEntries()
	le := iobuffer.String()
	assert.NotContains(t, le, "DEBUG: ")
	assert.NotContains(t, le, "INFO: ")
	assert.NotContains(t, le, "WARN: ")
	assert.NotContains(t, le, "ERROR: ")
	assert.Contains(t, le, "CRITICAL: ")
}
