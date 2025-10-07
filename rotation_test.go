package log

import (
	"io/fs"
	"os"
	"strings"
	"testing"
)

// TestSizeRotation writes enough data with a small size threshold and expects
// multiple log files to be created.
func TestSizeRotation(t *testing.T) {
	dir := t.TempDir()
	ResetForTests()
	if err := Init(dir, 1, true); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	SetLogEnable()
	SetLogToConsole(false)
	SetFilenamePrefix("rot", "rot")

	// Set small threshold (1KB) to trigger rotation quickly
	SetMaxFileSizeBytes(1024)

	// Write many messages to exceed the threshold multiple times
	for i := 0; i < 500; i++ {
		Info("rotation test message %d: %s", i, strings.Repeat("x", 40))
	}

	// Count .log files
	var count int
	fs.WalkDir(os.DirFS(dir), ".", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".log") {
			count++
		}
		return nil
	})

	if count < 2 {
		t.Fatalf("expected multiple log files created by rotation, got %d", count)
	}
}
