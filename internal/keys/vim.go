package keys

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func VimKeysCapture(app *tview.Application) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'h':
			return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
		case 'l':
			return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
		case 'q':
			app.Stop()
			return nil
		}
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
			return nil
		}
		return event
	}
}
