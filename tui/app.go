package tui

import (
	"context"
	"sync"
	"time"

	exo "github.com/deref/exo/kernel/api"
	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	t      *tview.Application
	client exo.Kernel

	shutdown chan struct{}

	hstack      *tview.Flex
	processList *ProcessList
	eventLog    *EventLog
	statusBar   *StatusBar

	showProcesses bool
	cursor        string
}

func NewApp(client exo.Kernel) *App {
	return &App{
		t:      tview.NewApplication(),
		client: client,
	}
}

func (app *App) Run(ctx context.Context) error {

	app.processList = NewProcessList()
	app.eventLog = NewEventLog()
	app.statusBar = NewStatusBar()

	app.hstack = tview.NewFlex()

	vstack := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.hstack, 0, 1, false).
		AddItem(app.statusBar, 1, 0, false)

	app.t.SetRoot(vstack, true)

	app.shutdown = make(chan struct{})
	go func() {
		delay := 0
		for {
			select {
			case <-time.After(time.Duration(delay) * time.Millisecond):
				var wg sync.WaitGroup
				wg.Add(2)
				go func() {
					app.t.QueueUpdateDraw(app.pollProcesses(ctx))
					wg.Done()
				}()
				go func() {
					app.t.QueueUpdateDraw(app.pollEvents(ctx))
					wg.Done()
				}()
				wg.Wait()
			case <-app.shutdown:
				return
			}
			delay = 250
		}
	}()

	app.t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				app.exit()
			case 'p':
				app.showProcesses = !app.showProcesses
			}
		case tcell.KeyEsc:
			app.exit()
		}

		app.render()
		return event
	})

	return app.t.Run()
}

func (app *App) render() {
	app.hstack.Clear()
	if app.showProcesses {
		app.hstack.AddItem(app.processList, 16, 0, false)
	}
	app.hstack.AddItem(app.eventLog, 0, 1, true)
}

func (app *App) poll(ctx context.Context) {
	app.pollProcesses(ctx)
	app.pollEvents(ctx)
	app.render()
}

func (app *App) pollProcesses(ctx context.Context) func() {
	processes, err := app.client.DescribeProcesses(ctx, &exo.DescribeProcessesInput{})
	return func() {
		if err != nil {
			app.reportError(err)
		}
		app.processList.SetData(processes.Processes)
		app.render()
	}
}

func (app *App) pollEvents(ctx context.Context) func() {
	events, err := app.client.GetEvents(ctx, &exo.GetEventsInput{
		After: app.cursor,
	})
	return func() {
		if err != nil {
			app.reportError(err)
			return
		}
		app.cursor = events.Cursor
		app.eventLog.AddEvents(events.Events)
		app.render()
	}
}

func (app *App) reportError(err error) {
	app.statusBar.SetErr(err)
}

func (app *App) exit() {
	close(app.shutdown)
	app.t.Stop()
}
