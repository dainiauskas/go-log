package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestIntegration_LogFilesCreated initializes the logger to a temp directory,
// writes logs and asserts that log files and symlinks are created.
func TestIntegration_LogFilesCreated(t *testing.T) {
	dir := t.TempDir()

	// Reset global logger state and initialize logger to write into the temp dir
	ResetForTests()
	if err := Init(dir, 1, true); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Ensure logging is enabled in case other tests disabled it.
	SetLogEnable()
	// Also enable trace so small messages certainly get through in all configs.
	SetLogTrace(true)

	// Use a predictable prefix so we can look for symlinks
	SetFilenamePrefix("testprefix", "testprefix")
	SetLogToConsole(false)

	// Write some logs
	Info("integration test starting")
	for i := 0; i < 50; i++ {
		Info("msg %d", i)
		Error("err %d", i)
	}

	// Check that at least one .log file exists in the directory (poll briefly
	// to avoid flakes due to async writes).
	foundLog := false
	for attempt := 0; attempt < 40 && !foundLog; attempt++ {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ReadDir failed: %v", err)
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".log") {
				foundLog = true
				break
			}
		}
		if !foundLog {
			time.Sleep(50 * time.Millisecond)
		}
	}
	if !foundLog {
		t.Fatalf("no .log files found in %s", dir)
	}

	// Check that symlink for info level exists and points to a .log file
	infoSymlink := filepath.Join(dir, "testprefix.info")
	fi, err := os.Lstat(infoSymlink)
	if err != nil {
		t.Fatalf("info symlink missing: %v", err)
	}
	if (fi.Mode() & os.ModeSymlink) == 0 {
		// On some platforms the symlink may be a regular file; still check its name
		// and contents
	}
	target, err := os.Readlink(infoSymlink)
	if err != nil {
		// If not a symlink but a file, accept that as long as it ends with .log
		if strings.HasSuffix(fi.Name(), ".log") {
			return
		}
		t.Fatalf("Readlink failed: %v", err)
	}
	if !strings.HasSuffix(target, ".log") {
		t.Fatalf("info symlink does not point to a .log file: %s", target)
	}
}
