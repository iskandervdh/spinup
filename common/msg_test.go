package common

import (
	"testing"
)

func TestSuccessMsg(t *testing.T) {
	m := NewSuccessMsg("success")

	if m == nil {
		t.Errorf("Expected %T, got nil", SuccessMsg{})
	}

	if m.GetText() != "success" {
		t.Errorf("Expected %s, got %s", "success", m.GetText())
	}

	m = NewSuccessMsg("success %s", "test")

	if m.GetText() != "success test" {
		t.Errorf("Expected %s, got %s", "success test", m.GetText())
	}
}

func TestInfoMsg(t *testing.T) {
	m := NewInfoMsg("info")

	if m == nil {
		t.Errorf("Expected %T, got nil", InfoMsg{})
	}

	if m.GetText() != "info" {
		t.Errorf("Expected %s, got %s", "info", m.GetText())
	}

	m = NewInfoMsg("info %s", "test")

	if m.GetText() != "info test" {
		t.Errorf("Expected %s, got %s", "info test", m.GetText())
	}
}

func TestWarnMsg(t *testing.T) {
	m := NewWarnMsg("warn")

	if m == nil {
		t.Errorf("Expected %T, got nil", WarnMsg{})
	}

	if m.GetText() != "warn" {
		t.Errorf("Expected %s, got %s", "warn", m.GetText())
	}

	m = NewWarnMsg("warn %s", "test")

	if m.GetText() != "warn test" {
		t.Errorf("Expected %s, got %s", "warn test", m.GetText())
	}
}

func TestErrorMsg(t *testing.T) {
	m := NewErrMsg("error")

	if m == nil {
		t.Errorf("Expected %T, got nil", ErrMsg{})
	}

	if m.GetText() != "error" {
		t.Errorf("Expected %s, got %s", "error", m.GetText())
	}

	m = NewErrMsg("error %s", "test")

	if m.GetText() != "error test" {
		t.Errorf("Expected %s, got %s", "error test", m.GetText())
	}
}

func TestRegularMsg(t *testing.T) {
	m := NewRegularMsg("regular")

	if m == nil {
		t.Errorf("Expected %T, got nil", RegularMsg{})
	}

	if m.GetText() != "regular" {
		t.Errorf("Expected %s, got %s", "regular", m.GetText())
	}

	m = NewRegularMsg("regular %s", "test")

	if m.GetText() != "regular test" {
		t.Errorf("Expected %s, got %s", "regular test", m.GetText())
	}
}
