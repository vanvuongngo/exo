package tui

import (
	"github.com/deref/exo/kernel/api"
	"github.com/rivo/tview"
)

type ProcessList struct {
	*tview.Flex
	title     *tview.TextView
	processes []api.ProcessDescription
	fill      *tview.Box
}

func NewProcessList() *ProcessList {
	view := &ProcessList{
		Flex: tview.NewFlex(),
		title: tview.NewTextView().
			SetText("Processes"),
		fill: tview.NewBox(),
	}
	view.render()
	return view
}

func (view *ProcessList) SetData(processes []api.ProcessDescription) {
	view.processes = processes
	view.render()
}

func (view *ProcessList) render() {
	view.Flex.
		Clear().
		SetDirection(tview.FlexRow).
		AddItem(view.title, 1, 0, false)
	if len(view.processes) == 0 {
		item := tview.NewTextView().SetText("(none)")
		view.Flex.AddItem(item, 1, 0, false)
	}
	for _, process := range view.processes {
		item := tview.NewTextView().SetText("  " + process.Name)
		view.Flex.AddItem(item, 1, 0, false)
	}
	view.Flex.AddItem(view.fill, 0, 1, false)
}
