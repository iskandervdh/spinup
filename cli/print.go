package cli

import (
	"fmt"

	"github.com/iskandervdh/spinup/common"
)

func (c *CLI) InfoPrint(a ...any) {
	fmt.Fprint(c.out, common.InfoText(fmt.Sprint(a...)))
}

func (c *CLI) InfoPrintf(format string, a ...any) {
	fmt.Fprint(c.out, common.InfoText(fmt.Sprintf(format, a...)))
}

func (c *CLI) SuccessPrint(a ...any) {
	fmt.Fprint(c.out, common.SuccessText(fmt.Sprint(a...)))
}

func (c *CLI) SuccessPrintf(format string, a ...any) {
	fmt.Fprint(c.out, common.SuccessText(fmt.Sprintf(format, a...)))
}

func (c *CLI) WarningPrint(a ...any) {
	fmt.Fprint(c.out, common.WarningText(fmt.Sprint(a...)))
}

func (c *CLI) WarningPrintf(format string, a ...any) {
	fmt.Fprint(c.out, common.WarningText(fmt.Sprintf(format, a...)))
}

func (c *CLI) ErrorPrint(a ...any) {
	fmt.Fprint(c.out, common.ErrorText(fmt.Sprint(a...)))
}

func (c *CLI) ErrorPrintf(format string, a ...any) {
	fmt.Fprint(c.out, common.ErrorText(fmt.Sprintf(format, a...)))
}

func (c *CLI) MsgPrint(msg common.Msg) {
	if msg == nil {
		return
	}

	switch msg.(type) {
	case *common.ErrMsg:
		c.ErrorPrint(msg.GetText())
	case *common.WarnMsg:
		c.WarningPrint(msg.GetText())
	case *common.InfoMsg:
		c.InfoPrint(msg.GetText())
	case *common.SuccessMsg:
		c.SuccessPrint(msg.GetText())
	default:
		fmt.Print(msg.GetText())
	}
}
