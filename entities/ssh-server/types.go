package sshserver

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/logging/prettyPrints"
	"github.com/zovenor/tea-models/models"
	"github.com/zovenor/tea-models/models/base"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SSHServer struct {
	Address  string
	Port     uint32
	Username string

	password string
	name     string
	tags     []string

	sshClient *ssh.Client
}

type Servers []*SSHServer

type CommandStatus int

const (
	DoneCommandStatus CommandStatus = iota
	WaitingCommandStatus
	WriteCommandStatus
)

// Methods for SSHServer

func New(
	address string,
	port uint32,
	username string,
	password string,
) *SSHServer {
	return &SSHServer{
		Address:  address,
		Port:     port,
		password: password,
		Username: username,
	}
}

func (s *SSHServer) SetName(name string) {
	s.name = name
}

func (s *SSHServer) Name() string {
	var name string
	if s.name == "" {
		name += fmt.Sprintf("%v:%v", s.Address, s.Port)
	} else {
		name += s.name
	}
	for _, tag := range s.tags {
		name += fmt.Sprintf(" %v", tag)
	}
	return name
}

func (s *SSHServer) AddTag(tag string) {
	s.tags = append(s.tags, tag)
}

func (s *SSHServer) Connect() error {
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Address, s.Port), config)
	if err != nil {
		return err
	}
	s.sshClient = client
	return nil
}

func (s *SSHServer) NewSession() (*ssh.Session, error) {
	return s.sshClient.NewSession()
}

// Methods for Servers

func (servers *Servers) AddServer(address string, port uint32, username string, password string) *SSHServer {
	newServer := New(address, port, username, password)
	*servers = append(*servers, newServer)
	return newServer
}

func (servers *Servers) ServerList() []*SSHServer {
	return *servers
}

func serversUpdateFunc(lim *models.ListItemsModel, msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case base.ForwardKey:
			fd := int(os.Stdout.Fd())

			// Get the terminal size
			width, height, err := terminal.GetSize(fd)
			if err != nil {
				fmt.Printf("error getting terminal size: %v\n", err)
				return lim, nil
			}
			serverInterface, err := lim.GetCurrentItem()
			if err != nil {
				lim.SetError(err)
				return lim, nil
			}
			server, ok := serverInterface.GetValue().(*SSHServer)
			if !ok {
				lim.SetError(fmt.Errorf("can not convert interface to SSHServer"))
			}
			err = server.Connect()
			if err != nil {
				lim.SetError(err)
				return lim, nil
			}
			session, err := server.NewSession()
			if err != nil {
				lim.SetError(err)
				return lim, nil
			}
			modes := ssh.TerminalModes{
				ssh.ECHO: 1,
				// ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_ISPEED: 28800,
				// ssh.TTY_OP_OSPEED: 14400,
				ssh.TTY_OP_OSPEED: 28800,
			}
			if err := session.RequestPty("xterm", height, width, modes); err != nil {
				lim.SetError(fmt.Errorf("failed to allocate a pseudo-terminal: %v", err))
				return lim, nil
			}
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			session.Stdin = os.Stdin
			prettyPrints.ClearTerminal()
			err = session.Shell()
			if err != nil {
				lim.SetError(err)
				return lim, nil
			}
			session.Wait()
			if err := session.Close(); err != nil {
				lim.SetError(err)
			}
			prettyPrints.ClearTerminal()
			return lim.Update(base.BackKey)
		}
	}

	return nil, nil
}

func (servers *Servers) ToTeaModel(parent tea.Model, parentPath string) (tea.Model, error) {
	updateF := serversUpdateFunc
	listViewConf := models.ListItemsConf{
		Name:           "Servers",
		FindMode:       true,
		MaxItemsInPage: 20,
		Parent:         parent,
		ParentPath:     parentPath,
		UpdateF:        &updateF,
	}
	listItemView, err := models.NewListItemsModel(listViewConf)
	if err != nil {
		return nil, err
	}
	for _, server := range servers.ServerList() {
		listItemView.AddItem(server.Name(), server)
	}
	return listItemView, nil
}
