package cli

import (
	"strings"
	"testing"

	"github.com/iskandervdh/spinup/common"
)

type textMethod func(string) string

func TestTextMethods(t *testing.T) {
	textMethods := []textMethod{common.InfoText, common.ErrorText, common.SuccessText, common.WarningText}

	for _, textMethod := range textMethods {
		text := textMethod("test")

		if !strings.Contains(text, "test") {
			t.Errorf("Expected text to contain 'test', got %s", text)
		}
	}
}

func TestPrints(t *testing.T) {
	c := New()

	c.InfoPrint("info")
	c.SuccessPrint("success")
	c.WarningPrint("warning")
	c.ErrorPrint("error")

	c.InfoPrintf("info %s", "printf")
	c.SuccessPrintf("success %s", "printf")
	c.WarningPrintf("warning %s", "printf")
	c.ErrorPrintf("error %s", "printf")
}
