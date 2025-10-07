//go:build test
// +build test

package log

import "time"

// gMaxFileSizeBytes controls rotation by size when > 0. Default 0 (disabled).
// The variable itself is declared in logger.go; tests may override it via
// SetMaxFileSizeBytes.

// SetMaxFileSizeBytes configures the file size threshold (in bytes) that triggers
// rotation when the current logfile grows beyond this value. A value of 0
// disables size-based rotation.
func SetMaxFileSizeBytes(n int64) {
	gMaxFileSizeBytes = n
}

// ResetForTests closes any open logger files and resets internal state. This is
// intended for use by tests to ensure a clean environment between test cases.
func ResetForTests() {
	for i := range gLoggers {
		gLoggers[i].lock.Lock()
		if gLoggers[i].file != nil {
			_ = gLoggers[i].file.Close()
			gLoggers[i].file = nil
		}
		gLoggers[i].day = 0
		gLoggers[i].size = 0
		gLoggers[i].purged = time.Time{}
		gLoggers[i].lock.Unlock()
	}
	// Reset config to defaults used on package init
	gConf = config{
		logPath:  "./log/",
		logflags: flagLogFilenameLineNum | flagLogThrough,
		maxdays:  30,
		enabled:  true,
	}
	isSymlink = map[string]bool{}
}
