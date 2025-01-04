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

func TestQuestion(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte(" "))
		w.Write([]byte("j"))
		w.Write([]byte(" "))
		w.Write([]byte(" "))
		w.Write([]byte("j"))
		w.Write([]byte("k"))
		w.Write([]byte("j"))
		w.Write([]byte(" "))
		w.Write([]byte("enter"))
	}()

	answers, err, exited := c.Question("test?", []string{"a", "b", "c"}, nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if exited {
		t.Errorf("expected exited to be false, got true")
	}

	expectedAnswers := []string{"a", "c"}

	if len(answers) != len(expectedAnswers) {
		t.Errorf("expected %d answers, got %d", len(expectedAnswers), len(answers))
	}

	for i, answer := range answers {
		if answer != expectedAnswers[i] {
			t.Errorf("expected answer '%s', got '%s'", expectedAnswers[i], answer)
		}
	}
}

func TestQuestionArgumentLengthMismatch(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte(" "))
		w.Write([]byte("j"))
		w.Write([]byte(" "))
		w.Write([]byte(" "))
		w.Write([]byte("j"))
		w.Write([]byte("k"))
		w.Write([]byte("j"))
		w.Write([]byte(" "))
		w.Write([]byte("enter"))
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Question arguments mismatch did not panic")
			return
		}
	}()

	_, err, exited := c.Question("test?", []string{"a", "b", "c"}, []bool{true, false})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if !exited {
		t.Errorf("expected exited to be true, got false")
	}
}

func TestSelection(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("j"))
		w.Write([]byte("k"))
		w.Write([]byte("j"))
		w.Write([]byte("enter"))
	}()

	answer, err, exited := c.Selection("test?", []string{"a", "b", "c"})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if exited {
		t.Errorf("expected exited to be false, got true")
	}

	expectedAnswer := "b"

	if answer != expectedAnswer {
		t.Errorf("expected answer '%s', got '%s'", expectedAnswer, answer)
	}
}

func TestInput(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))
	inputString := "test&*!@#123"

	go func() {
		defer w.Close()

		w.Write([]byte(inputString))
		w.Write([]byte("enter"))
	}()

	input := c.Input("test?", "")

	if input != inputString {
		t.Errorf("expected input '%s', got '%s'", inputString, input)
	}
}

func TestInputWithDefault(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))
	inputString := "test&*!@#123"

	go func() {
		defer w.Close()

		w.Write([]byte("enter"))
	}()

	input := c.Input("test?", inputString)

	if input != inputString {
		t.Errorf("expected input '%s', got '%s'", inputString, input)
	}
}

func TestInputBackspace(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("backspace"))
		w.Write([]byte("enter"))
	}()

	input := c.Input("test?", "a")

	if input != "" {
		t.Errorf("expected input '', got '%s'", input)
	}
}

func TestInputDelete(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("left"))
		w.Write([]byte("delete"))
		w.Write([]byte("enter"))
	}()

	input := c.Input("test?", "a")

	if input != "" {
		t.Errorf("expected input '', got '%s'", input)
	}
}

func TestInputLeftRight(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("left"))
		w.Write([]byte("left"))
		w.Write([]byte("backspace"))
		w.Write([]byte("right"))
		w.Write([]byte("right"))
		w.Write([]byte("backspace"))
		w.Write([]byte("enter"))
	}()

	input := c.Input("test?", "abcde")

	if input != "abd" {
		t.Errorf("expected input 'abd', got '%s'", input)
	}
}

func TestInputCtrlC(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("ctrl+c"))
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected os.Exit to be called, but it was not")
		}
	}()

	input := c.Input("test?", "")

	if input != "" {
		t.Errorf("expected input '', got '%s'", input)
	}
}

func TestInputInsertBetween(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("left"))
		w.Write([]byte("a"))
		w.Write([]byte("right"))
		w.Write([]byte("c"))
		w.Write([]byte("enter"))
	}()

	input := c.Input("test?", "b")

	if input != "abc" {
		t.Errorf("expected input 'abc', got '%s'", input)
	}
}

func TestInputIllegalButtons(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("up"))
		w.Write([]byte("down"))
		w.Write([]byte("tab"))

		time.Sleep(1 * time.Second)

		w.Write([]byte("enter"))
	}()

	input := c.Input("test?", "")

	if input != "" {
		t.Errorf("expected input '', got '%s'", input)
	}
}

func TestConfirm(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("y"))
		w.Write([]byte("enter"))
	}()

	confirmed := c.Confirm("test?")
	expectedConfirmed := true

	if confirmed != expectedConfirmed {
		t.Errorf("expected confirmed '%t', got '%t'", expectedConfirmed, confirmed)
	}
}

func TestConfirmAction(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("y"))
		w.Write([]byte("backspace"))
		w.Write([]byte("n"))
		w.Write([]byte("left"))
		w.Write([]byte("n"))
		w.Write([]byte("delete"))
		w.Write([]byte("right"))
		w.Write([]byte("up"))
		w.Write([]byte("down"))
		w.Write([]byte("tab"))

		time.Sleep(1 * time.Second)

		w.Write([]byte("enter"))
	}()

	c.Confirm("test?")
}

func TestConfirmCtrlC(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("ctrl+c"))
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected os.Exit to be called, but it was not")
		}
	}()

	input := c.Input("test?", "")

	if input != "" {
		t.Errorf("expected input '', got '%s'", input)
	}
}

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

func TestQuitQuestion(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("ctrl+c"))
	}()

	answers, err, exited := c.Question("test?", []string{"a", "b", "c"}, nil)

	if answers != nil {
		t.Errorf("expected answers to be nil, got %v", answers)
	}

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !exited {
		t.Errorf("expected exited to be true, got false")
	}
}

func TestQuitLoading(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := TestingCLI("loading_quit", WithIn(r), WithOut(output))

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

func TestQuitSelection(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("ctrl+c"))
	}()

	answer, err, exited := c.Selection("test?", []string{"a", "b", "c"})

	if answer != "" {
		t.Errorf("expected answer to be empty, got %s", answer)
	}

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !exited {
		t.Errorf("expected exited to be true, got false")
	}
}

func TestLoadingDone(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := TestingCLI("loading_done", WithIn(r), WithOut(output))

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
