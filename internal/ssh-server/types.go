package sshserver

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/logging/prettyPrints"
	"github.com/zovenor/susshi/pkg/files"
	"github.com/zovenor/tea-models/models"
	"github.com/zovenor/tea-models/models/base"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SSHServer struct {
	Address  string `json:"address"`
	Port     uint32 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
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
		Password: password,
		Username: username,
	}
}

func (s *SSHServer) SetName(name string) {
	s.Name = name
}

func (s *SSHServer) GetName() string {
	var name string
	if s.Name == "" {
		name += fmt.Sprintf("%v:%v", s.Address, s.Port)
	} else {
		name += s.Name
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
			ssh.Password(s.Password),
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
		listItemView.AddItem(server.GetName(), server)
	}
	return listItemView, nil
}

func (servers *Servers) ImportFromFile(filePath string) error {
	files.CreateFolder(filePath)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, servers)
	if err != nil {
		return err
	}
	return nil
}
