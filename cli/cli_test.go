package cli

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDoneMsg(t *testing.T) {
	d := DoneMsg("test")

	if d.text != "test" {
		t.Errorf("Expected %s, got %s", "test", d.text)
	}
}

func TestErrorMsg(t *testing.T) {
	e := ErrMsg("test")

	if e.text != "test" {
		t.Errorf("Expected %s, got %s", "test", e.text)
	}
}

func TestNew(t *testing.T) {
	c := New()

	if c.out != os.Stdout {
		t.Errorf("Expected %v, got %v", os.Stdout, c.out)
	}
}

func TestWithIn(t *testing.T) {
	r := bytes.NewBuffer(nil)
	c := New(WithIn(r))

	if c.in != r {
		t.Errorf("Expected %v, got %v", r, c.in)
	}
}

func TestWithOut(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := New(WithOut(w))

	if c.out != w {
		t.Errorf("Expected %v, got %v", w, c.out)
	}
}

func TestClearTerminal(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := New(WithOut(w))
	c.ClearTerminal()

	if w.String() != "\033[H\033[2J" {
		t.Errorf("Expected %s, got %s", "\033[H\033[2J", w.String())
	}
}

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

	answers, err, exited := c.Question("test?", []string{"a", "b", "c"})

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

	go func() {
		defer w.Close()

		w.Write([]byte("test &*!@#123"))
		w.Write([]byte("enter"))
	}()

	input := c.Input("test?")
	expectedInput := "test &*!@#123"

	if input != expectedInput {
		t.Errorf("expected input '%s', got '%s'", expectedInput, input)
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

func TestLoading(t *testing.T) {
	r, w := io.Pipe()

	output := &bytes.Buffer{}
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("test"))
	}()

	start := time.Now()

	c.Loading("test", func() tea.Msg {
		time.Sleep(200 * time.Millisecond)
		return DoneMsg("test")
	})

	elapsed := time.Since(start)

	if elapsed < 200*time.Millisecond || elapsed > 400*time.Millisecond {
		t.Errorf("expected loading to take about 200 milliseconds, but it took %v", elapsed)
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

	answers, err, exited := c.Question("test?", []string{"a", "b", "c"})

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

// func TestQuitLoading(t *testing.T) {
// 	r, w := io.Pipe()

// 	output := &bytes.Buffer{}
// 	c := New(WithIn(r), WithOut(output))

// 	go func() {
// 		defer w.Close()

// 		w.Write([]byte("ctrl+c"))
// 	}()

// 	c.Loading("test", func() tea.Msg {
// 		time.Sleep(1 * time.Second)

// 		return DoneMsg("test")
// 	})

// 	if output.String() != "" {
// 		t.Errorf("expected no output, got %s", output.String())
// 	}
// }

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
	c := New(WithIn(r), WithOut(output))

	go func() {
		defer w.Close()

		w.Write([]byte("test"))
	}()

	l := c.Loading("test", func() tea.Msg {
		return DoneMsg("test")
	})

	if !l.done {
		t.Errorf("expected done to be true, got false")
	}

	if l.errorText != "" {
		t.Errorf("expected errorText to be empty, got %s", l.errorText)
	}
}
