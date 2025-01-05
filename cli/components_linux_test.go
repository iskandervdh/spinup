package cli

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/iskandervdh/spinup/common"
)

func TestLoading(t *testing.T) {
	c := TestingCLI("loading")

	start := time.Now()

	loadingMessage := c.Loading("test", func() common.Msg {
		time.Sleep(200 * time.Millisecond)
		return common.NewSuccessMsg(fmt.Sprintf("Completed loading '%s'", "test"))
	})

	elapsed := time.Since(start)

	if elapsed < 200*time.Millisecond || elapsed > 400*time.Millisecond {
		t.Errorf("expected loading to take about 200 milliseconds, but it took %v", elapsed)
	}

	if strings.Contains(loadingMessage.GetText(), "Completed loading 'test'") {
		t.Errorf("expected loading message to contain 'Completed loading 'test'', got '%s'", loadingMessage.GetText())
	}
}

func TestQuitLoading(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := TestingCLI("loading_quit", WithIn(r), WithOut(output), WithErr(output))

	go func() {
		defer w.Close()

		w.Write([]byte("ctrl+c"))
	}()

	c.Loading("test", func() common.Msg {
		time.Sleep(1 * time.Second)

		return common.NewSuccessMsg("test")
	})

	if output.String() != "" {
		t.Errorf("expected no output, got %s", output.String())
	}
}

func TestLoadingDone(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := TestingCLI("loading_done", WithIn(r), WithOut(output), WithErr(output))

	go func() {
		defer w.Close()

		w.Write([]byte("test"))
	}()

	loadingMessage := c.Loading("test", func() common.Msg {
		return common.NewSuccessMsg("test")
	})

	if loadingMessage == nil {
		t.Errorf("expected done to be true, got false")
		return
	}

	if loadingMessage.GetText() != "" {
		t.Errorf("expected errorText to be empty, got %s", loadingMessage.GetText())
	}
}

func TestLoadingError(t *testing.T) {
	c := TestingCLI("loading_error")

	loadingMessage := c.Loading("test", func() common.Msg {
		return common.NewErrMsg("test error")
	})

	if loadingMessage == nil {
		t.Errorf("expected done to be true, got false")
		return
	}

	if strings.Contains(loadingMessage.GetText(), "test error") {
		t.Errorf("expected errorText to be 'test error', got %s", loadingMessage.GetText())
	}
}
