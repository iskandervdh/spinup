package cli

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
	"github.com/iskandervdh/spinup/core"
)

func TestingConfigDir(testName string) string {
	return path.Join(os.TempDir(), common.ProgramName, testName)
}

func TestingCore(testName string) *core.Core {
	// Remove old tmp config dir
	testingConfigDir := TestingConfigDir(testName)
	err := os.RemoveAll(testingConfigDir)

	if err != nil {
		panic(err)
	}

	// Mock msgChan to prevent blocking during testing
	msgChan := new(chan common.Msg)
	*msgChan = make(chan common.Msg)

	cfg := config.NewTesting(testingConfigDir)
	c := core.New(core.WithConfig(cfg), core.WithMsgChan(msgChan))

	// Mock init to prevent errors during testing
	c.Init()

	return c
}

func TestingCLI(testName string, options ...func(*CLI)) *CLI {
	return New(append([]func(*CLI){WithCore(TestingCore(testName))}, options...)...)
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

func TestWithErr(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := New(WithErr(w))

	if c.err != w {
		t.Errorf("Expected %v, got %v", w, c.err)
	}
}

func TestWithCore(t *testing.T) {
	co := core.New()
	c := New(WithCore(co))

	if c.core != co {
		t.Errorf("Expected %v, got %v", co, c.core)
	}
}

func TestMsgPrint(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := New(WithOut(w), WithErr(w))

	c.MsgPrint(common.NewSuccessMsg("test"))

	if w.String() != "test\n" {
		t.Errorf("Expected %s, got %s", "test\n", w.String())
	}
}

func TestSendMsg(t *testing.T) {
	c := New()

	c.sendMsg(common.NewSuccessMsg("test"))

	msg := <-*c.msgChan

	if msg.GetText() != "test" {
		t.Errorf("Expected %s, got %s", "test", msg.GetText())
	}
}

func TestHelpMsg(t *testing.T) {
	c := TestingCLI("send_help_msg")

	os.Args = []string{"spinup", "--help"}
	c.Handle()
}

func TestClearTerminal(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := New(WithOut(w), WithErr(w))
	c.ClearTerminal()

	if w.String() != "\033[H\033[2J" {
		t.Errorf("Expected %s, got %s", "\033[H\033[2J", w.String())
	}
}

func TestCLIHandleUnknownSubcommand(t *testing.T) {
	c := TestingCLI("handle")

	// Test handle without any arguments
	os.Args = []string{"spinup", "handle"}
	c.Handle()
}

// func TestCLIHandleNoArgs(t *testing.T) {
// 	r, w := io.Pipe()

// 	os.Args = []string{"spinup"}

// 	output := &bytes.Buffer{}
// 	c := TestingCLI("handle_no_args", WithIn(r), WithOut(output), WithErr(output))

// 	go func() {
// 		defer w.Close()

// 		w.Write([]byte("ctrl+c"))
// 	}()

// 	c.Handle()
// }

func TestCLIHandleInit(t *testing.T) {
	c := TestingCLI("handle_init")

	os.Args = []string{"spinup", "init"}
	c.Handle()
}

func TestCLIHandleVersion(t *testing.T) {
	c := TestingCLI("handle_version")

	os.Args = []string{"spinup", "-v"}
	c.Handle()
}

func TestCLIHandleCommand(*testing.T) {
	c := TestingCLI("handle_command")

	os.Args = []string{"spinup", "c"}
	c.Handle()

	c = TestingCLI("handle_command")
	os.Args = []string{"spinup", "c", "ls"}
	c.Handle()

	c = TestingCLI("handle_command")
	os.Args = []string{"spinup", "c", "add", "test"}
	c.Handle()

	c = TestingCLI("handle_command")
	os.Args = []string{"spinup", "c", "add", "test", "echo test"}
	c.Handle()

	c = TestingCLI("handle_command")
	os.Args = []string{"spinup", "c", "rm", "test"}
	c.Handle()

	c = TestingCLI("handle_command")
	os.Args = []string{"spinup", "c", "test"}
	c.Handle()
}

func TestCLIHandleProject(t *testing.T) {
	c := TestingCLI("handle_project")

	os.Args = []string{"spinup", "p"}
	c.Handle()

	c = TestingCLI("handle_project")
	os.Args = []string{"spinup", "p", "ls"}
	c.Handle()

	c = TestingCLI("handle_project")
	os.Args = []string{"spinup", "p", "add", "test", "echo test"}
	c.Handle()

	// c = TestingCLI("handle_project")
	// os.Args = []string{"spinup", "p", "rm", "test"}
	// c.Handle()

	c = TestingCLI("handle_project")
	os.Args = []string{"spinup", "p", "test"}
	c.Handle()
}

func TestCLIHandleVariable(t *testing.T) {
	c := TestingCLI("handle_variable")

	os.Args = []string{"spinup", "v"}
	c.Handle()

	c = TestingCLI("handle_variable")
	os.Args = []string{"spinup", "v", "ls"}
	c.Handle()

	c = TestingCLI("handle_variable")
	os.Args = []string{"spinup", "v", "ls", "test"}
	c.Handle()

	c = TestingCLI("handle_variable")
	os.Args = []string{"spinup", "v", "add", "test"}
	c.Handle()

	c = TestingCLI("handle_variable")
	os.Args = []string{"spinup", "v", "add", "test", "echo test"}
	c.Handle()

	c = TestingCLI("handle_variable")
	os.Args = []string{"spinup", "v", "rm", "test"}
	c.Handle()

	c = TestingCLI("handle_variable")
	os.Args = []string{"spinup", "v", "test"}
	c.Handle()
}

func TestCLIHandleRun(t *testing.T) {
	c := TestingCLI("handle_run")

	os.Args = []string{"spinup", "run", "test"}
	c.Handle()

	c = TestingCLI("handle_run")
	os.Args = []string{"spinup", "run", "test", "echo test"}
	c.Handle()

	c = TestingCLI("handle_run")
	os.Args = []string{"spinup", "run"}
	c.Handle()
}
