package common

import "fmt"

type Msg interface {
	GetText() string
}

type SuccessMsg struct {
	text string
}

func NewSuccessMsg(text string, a ...any) *SuccessMsg {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}

	return &SuccessMsg{text: text}
}

func (m *SuccessMsg) GetText() string {
	return m.text
}

type InfoMsg struct {
	text string
}

func NewInfoMsg(text string, a ...any) *InfoMsg {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}

	return &InfoMsg{text: text}
}

func (m *InfoMsg) GetText() string {
	return m.text
}

type WarnMsg struct {
	text string
}

func NewWarnMsg(text string, a ...any) *WarnMsg {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}

	return &WarnMsg{text: text}
}

func (m *WarnMsg) GetText() string {
	return m.text
}

type ErrMsg struct {
	text string
}

func NewErrMsg(text string, a ...any) *ErrMsg {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}

	return &ErrMsg{text: text}
}

func (m *ErrMsg) GetText() string {
	return m.text
}

type RegularMsg struct {
	text string
}

func NewRegularMsg(text string, a ...any) *RegularMsg {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}

	return &RegularMsg{text: text}
}

func (m *RegularMsg) GetText() string {
	return m.text
}
