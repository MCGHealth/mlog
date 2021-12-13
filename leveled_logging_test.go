package mlog_test

import (
	"strings"
	"testing"

	"github.com/MCGhealth/mlog"
	"github.com/stretchr/testify/assert"
)

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

func TestDebugLoggingLevel(t *testing.T) {
	mlog.Initialize(&iobuffer, mlog.DebugLevel)
	iobuffer.Reset()
	mlog.SetLogLevel(mlog.DebugLevel)

	logDebug()
	logInfo()
	logWarn()

	le := iobuffer.String()

	assert.True(t, strings.Contains(le, "DEBUG"))
	assert.True(t, strings.Contains(le, "INFO"))
	assert.True(t, strings.Contains(le, "WARN"))
}

func TestInfoLoggingLevel(t *testing.T) {
	mlog.Initialize(&iobuffer, mlog.InfoLevel)
	iobuffer.Reset()
	mlog.SetLogLevel(mlog.InfoLevel)

	logDebug()
	logInfo()
	logWarn()

	le := iobuffer.String()

	assert.False(t, strings.Contains(le, "DEBUG"))
	assert.True(t, strings.Contains(le, "INFO"))
	assert.True(t, strings.Contains(le, "WARN"))
}

func TestWarnLoggingLevel(t *testing.T) {
	mlog.Initialize(&iobuffer, mlog.WarnLevel)
	iobuffer.Reset()
	mlog.SetLogLevel(mlog.WarnLevel)

	logDebug()
	logInfo()
	logWarn()

	le := iobuffer.String()

	assert.False(t, strings.Contains(le, "DEBUG"))
	assert.False(t, strings.Contains(le, "INFO"))
	assert.True(t, strings.Contains(le, "WARN"))
}

func TestErrorLoggingLevel(t *testing.T) {
	mlog.Initialize(&iobuffer, mlog.ErrorLevel)
	iobuffer.Reset()
	mlog.SetLogLevel(mlog.ErrorLevel)

	logDebug()
	logInfo()
	logWarn()

	le := iobuffer.String()

	assert.False(t, strings.Contains(le, "DEBUG"))
	assert.False(t, strings.Contains(le, "INFO"))
	assert.False(t, strings.Contains(le, "WARN"))
}
