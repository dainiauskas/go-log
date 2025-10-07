package log

import (
	"os"
	"sync"
	"testing"
)

// TestConcurrentLogging spawns many goroutines that write logs concurrently.
// Run with -race to detect data races.
func TestConcurrentLogging(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-log-concurrency-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	// Do not remove the dir immediately; some OSes may keep file handles open.
	// defer os.RemoveAll(dir)

	if err := Init(dir, 1, true); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Reduce console noise during test
	SetLogToConsole(false)

	var wg sync.WaitGroup
	nGoroutines := 50
	perG := 200
	wg.Add(nGoroutines)

	for i := 0; i < nGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				Info("gor=%d info %d", id, j)
				Error("gor=%d error %d", id, j)
				Query("SELECT * FROM t WHERE id=?", id)
				Debug("gor=%d debug %d", id, j)
				Trace("gor=%d trace %d", id, j)
			}
		}(i)
	}

	wg.Wait()

	// Disable further logging to avoid race with test cleanup if any.
	SetLogDisable()

	// Attempt best-effort cleanup
	_ = os.RemoveAll(dir)
}
