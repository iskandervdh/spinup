package common

import (
	"testing"
)

func TestIsWindows(t *testing.T) {
	if IsWindows() != false {
		t.Error("Expected false, got true")
	}
}

func TestAppInstalled(t *testing.T) {
	AppInstalled()
}
