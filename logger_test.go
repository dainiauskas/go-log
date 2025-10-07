package log

import (
	"testing"
)

// TestGormPrintEmptyArgs ensures Gorm.Print does not panic when called with no arguments.
func TestGormPrintEmptyArgs(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Gorm.Print panicked with empty args: %v", r)
		}
	}()

	var g Gorm
	// Call with zero args
	g.Print()
}

// TestGormPrintShortArgs ensures Gorm.Print does not panic when called with short arg slices.
func TestGormPrintShortArgs(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Gorm.Print panicked with short args: %v", r)
		}
	}()

	var g Gorm
	g.Print("sql")
	g.Print("sql", "one")
	g.Print("log")
	g.Print("log", "one")
}

// TestLoggerPrintlnPrintf ensures Logger.Println and Logger.Printf wrapper do not panic.
func TestLoggerPrintlnPrintf(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Logger methods panicked: %v", r)
		}
	}()

	var l Logger
	l.Println("a", "b", 1)
	l.Printf("%s %d", "x", 2)
}
