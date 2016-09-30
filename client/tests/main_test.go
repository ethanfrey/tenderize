package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	StartNode()
	os.Exit(m.Run())
}
