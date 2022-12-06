package log

import "testing"

func Test_Init(t *testing.T) {
	Init("./log", 1, true)

	Info("Test Info")
}
