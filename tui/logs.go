package tui

import (
	"fmt"

	"github.com/deref/exo/kernel/api"
	"github.com/rivo/tview"
)

type EventLog struct {
	*tview.TextView
	title  *tview.TextView
	events []api.Event
}

func NewEventLog() *EventLog {
	view := &EventLog{
		TextView: tview.NewTextView().ScrollToEnd(),
	}
	view.render()
	return view
}

func (view *EventLog) AddEvents(events []api.Event) {
	view.events = append(view.events, events...)
	view.render()
}

func (view *EventLog) render() {
	view.Clear()
	for _, event := range view.events {
		_, _ = fmt.Fprintf(view, "%s %s\n", event.Log, event.Message)
	}
}
