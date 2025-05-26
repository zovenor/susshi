package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zovenor/susshi/internal/config"
	"github.com/zovenor/susshi/internal/keys"
	"github.com/zovenor/susshi/internal/ssh"
)

type App struct {
	tApp  *tview.Application
	cfg   *config.Config
	hosts []*ssh.Host
}

func New(cfg *config.Config) (*App, error) {
	hosts, err := ssh.LoadHosts(cfg.SSHConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load ssh config: %s", err)
	}

	tApp := tview.NewApplication()
	tApp.SetInputCapture(keys.VimKeysCapture(tApp))
	return &App{
		tApp:  tApp,
		cfg:   cfg,
		hosts: hosts,
	}, nil
}

func (a *App) Run() error {
	mainFlex := tview.NewFlex()
	mainFlex.SetDirection(tview.FlexRow)

	if !a.cfg.HideIcon {
		iconBox := tview.NewTextView()
		iconBox.SetTextAlign(tview.AlignCenter).SetBorderPadding(1, 0, 0, 0)
		iconBox.SetText(a.IconASCII()).SetBackgroundColor(tcell.ColorDefault)
		iconBox.SetDynamicColors(true)
		mainFlex.AddItem(iconBox, 0, 1, true)
	}

	table := tview.NewTable()
	table.SetBorderAttributes(tcell.AttrBold)
	table.SetBackgroundColor(tcell.ColorDefault).SetTitle("Hosts").SetBorder(true)

	// Headers
	for i, header := range ssh.HeadersStringList {
		table.SetCell(0, i, &tview.TableCell{
			Text:          header,
			Color:         tcell.ColorOrangeRed,
			Attributes:    tcell.AttrBold,
			NotSelectable: true,
			Expansion:     1,
		})
	}
	table.SetFixed(1, 0)

	for i, host := range a.hosts {
		for header, value := range host.Iter() {
			cell := tview.NewTableCell(value)
			table.SetCell(i+1, int(header), cell)
		}
	}
	table.SetSelectable(true, false)
	table.SetSelectedFunc(func(row, column int) {
		host := a.hosts[row-1]
		a.tApp.Stop()
		ssh.Connect(host)
	})
	mainFlex.AddItem(table, 0, 1, true)

	if err := a.tApp.SetRoot(mainFlex, true).SetFocus(table).Run(); err != nil {
		return err
	}
	return nil
}

func (a *App) IconASCII() string {
	return iconTview
}
