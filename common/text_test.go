package common

import (
	"strings"
	"testing"
)

func TestInfoText(t *testing.T) {
	text := InfoText("Info text")

	if !strings.Contains(text, "Info text") {
		t.Errorf("Expected 'Info text', got '%s'", text)
	}
}

func TestSuccessText(t *testing.T) {
	text := SuccessText("Success text")

	if !strings.Contains(text, "Success text") {
		t.Errorf("Expected 'Success text', got '%s'", text)
	}
}

func TestWarningText(t *testing.T) {
	text := WarningText("Warning text")

	if !strings.Contains(text, "Warning text") {
		t.Errorf("Expected 'Warning text', got '%s'", text)
	}
}

func TestErrorText(t *testing.T) {
	text := ErrorText("Error text")

	if !strings.Contains(text, "Error text") {
		t.Errorf("Expected 'Error text', got '%s'", text)
	}
}
