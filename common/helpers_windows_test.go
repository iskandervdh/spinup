package common

import (
	"testing"
)

func TestIsWindows(t *testing.T) {
	if IsWindows() != true {
		t.Error("Expected false, got true")
	}
}

func TestIsMacOS(t *testing.T) {
	if IsMacOS() != false {
		t.Error("Expected false, got true")
	}
}
