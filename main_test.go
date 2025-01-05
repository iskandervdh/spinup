package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Error("Error creating pipe for main function", err)
		return
	}

	stdOut := os.Stdout
	os.Stdout = w

	main()

	w.Close()
	out, err := io.ReadAll(r)

	if err != nil {
		t.Error("Error reading output", err)
		return
	}

	if !(strings.Contains(string(out), "Expected")) {
		t.Errorf("Expected command message not printed, got: '%s'", string(out))
		return
	}

	os.Stdout = stdOut
}
