//go:generate ./_local/go-enum -f=$GOFILE --marshal --lower --noprefix  --nocase

package mlog

// LogEventLevel is an enumeration of log event levels that are allowed.
/* ENUM(
Unknown, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, CriticalLevel
)
*/
type LogEventLevel int
