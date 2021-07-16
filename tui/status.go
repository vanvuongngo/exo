package tui

import "github.com/rivo/tview"

type StatusBar struct {
	*tview.Flex
	message *tview.TextView
	err     error
}

func NewStatusBar() *StatusBar {
	view := &StatusBar{}
	help := tview.NewTextView().
		SetText("  [q]uit    toggle [p]rocess list")
	view.message = tview.NewTextView().
		SetTextAlign(tview.AlignRight)
	view.Flex = tview.NewFlex().
		AddItem(help, 0, 1, false).
		AddItem(view.message, 30, 0, false)
	view.render()
	return view
}

func (view *StatusBar) SetErr(err error) {
	view.err = err
}

func (view *StatusBar) render() {
	msg := "status: ok"
	if view.err != nil {
		msg = view.err.Error()
	}
	view.message.SetText(msg + "  ")
}
