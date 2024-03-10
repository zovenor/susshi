package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/zovenor/logging/prettyPrints"
	sshserver "github.com/zovenor/susshi/entities/ssh-server"
	"github.com/zovenor/tea-models/models"
)

type App struct {
	servers *sshserver.Servers
}

func New(servers *sshserver.Servers) *App {
	return &App{
		servers: servers,
	}
}

func (app *App) Run() error {
	prettyPrints.ClearTerminal()

	listViewConf := models.ListItemsConf{
		Name:           "Susshi",
		FindMode:       true,
		MaxItemsInPage: 20,
	}
	listView, err := models.NewListItemsModel(listViewConf)
	if err != nil {
		return err
	}
	serversView, err := app.servers.ToTeaModel(listView, listViewConf.Name)
	if err != nil {
		return err
	}
	listView.AddItem("Servers", serversView)

	p := tea.NewProgram(listView)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
