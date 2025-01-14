package common

import (
	"testing"
)

func TestIsWindows(t *testing.T) {
	if IsWindows() != false {
		t.Error("Expected false, got true")
	}
}

func TestIsMacOS(t *testing.T) {
	if IsMacOS() != true {
		t.Error("Expected false, got true")
	}
}
