package log

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	//Init("./log", 10, 2, 2, false)

	benchmark(b, 1, 1, 100000)
	benchmark(b, 1, 100000, 1)
	benchmark(b, 2, 100000, 1)
	benchmark(b, 4, 100000, 1)
	benchmark(b, 8, 100000, 1)
}

func benchmark(b *testing.B, nProcs, nGoroutines, nWrites int) {
	runtime.GOMAXPROCS(nProcs)
	bn := fmt.Sprintf("%d Procs %d Goroutines each makes %d writes", nProcs, nGoroutines, nWrites)

	b.Run(bn, func(b *testing.B) {
		var wg sync.WaitGroup
		for i := 0; i < nGoroutines; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < nWrites; j++ {
					Info("Failed to find player! uid=%d plid=%d cmd=%s xxx=%d", 1234, 678942, "getplayer", 102020101)
				}
				wg.Add(-1)
			}()
		}
		wg.Wait()
	})
}

func TestInit(t *testing.T) {
	err := Init("/log", 10, 2, 2, false)
	if err == nil {
		t.Error("Expected error but return nil")
	}
	err = Init("./log", 100001, 2, 2, false)
	if err == nil {
		t.Error("Expected error but return nil")
	}
	err = Init("./log", 10, -1, 2, false)
	if err == nil {
		t.Error("Expected error but return nil")
	}

	err = Init("./log", 10, 2, 2, false)
	if err != nil {
		t.Errorf("Wanted error nil, got: %v", err)
	}
}

func TestLogger(t *testing.T) {
	t.Run("SetLogTrace", func(t *testing.T) {
		SetLogTrace(true)
		if gConf.logTrace() != true {
			t.Error("Log Trace wanted true, got false")
		}
	})

	// Testing Through flag
	t.Run("SetLogThrough", func(t *testing.T) {
		SetLogThrough(true)
		if gConf.logThrough() != true {
			t.Error("Log Through wanted true, got false")
		}
	})

	t.Run("SetLogFunctionName", func(t *testing.T) {
		SetLogFunctionName(true)
		if gConf.logFuncName() != true {
			t.Error("Log FunctionName wanted true, got false")
		}
	})

	t.Run("SetLogFilenameLineNum", func(t *testing.T) {
		SetLogFilenameLineNum(true)
		if gConf.logFilenameLineNum() != true {
			t.Error("Log FilenameLineNum wanted true, got false")
		}
	})

	t.Run("SetLogToConsole", func(t *testing.T) {
		SetLogToConsole(true)
		if gConf.logToConsole() != true {
			t.Error("Log LogToConsole wanted true, got false")
		}
	})

	t.Run("SetLogToConsole", func(t *testing.T) {
		on := gConf.isEnabled()
		if on != true {
			t.Error("Log isEnabled() wanted true, got false")
		}
	})

	t.Run("SetLogUserName", func(t *testing.T) {
		SetLogUserName("admin")
		if infoUserName != "admin" {
			t.Errorf("Log SetLogUserName() wanted admin, got %s", infoUserName)
		}
	})

	t.Run("SetLogDisable", func(t *testing.T) {
		SetLogDisable()
		if gConf.enabled != false {
			t.Errorf("Log SetLogDisable() wanted false, got true")
		}
	})

	t.Run("SetLogEnable", func(t *testing.T) {
		SetLogEnable()
		if gConf.enabled != true {
			t.Errorf("Log SetLogEnable() wanted true, got false")
		}
	})

	t.Run("Output", func(t *testing.T) {
		Console("Test %s", "console")
		Trace("Test %s", "trace")
		Info("Test %s", "info")
		Update("Test %s", "update")
		Warn("Test %s", "warn")
		Error("Test %s", "error")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		Panic("Test %s", "panic")
	})

	t.Run("setMaxSize", func(t *testing.T) {
		gConf.setMaxSize(0)
		if gConf.maxsize != 9222246136947933183 {
			t.Errorf("setMaxSize() wanted 9222246136947933183 got %v", gConf.maxsize)
		}

		gConf.setMaxSize(2)
		if gConf.maxsize != 2097152 {
			t.Errorf("setMaxSize() wanted 2097152 got %v", gConf.maxsize)
		}
	})
}
