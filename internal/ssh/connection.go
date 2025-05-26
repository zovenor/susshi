package ssh

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func Connect(host *Host) {
	user := "root"
	identity := ""
	port := ""
	if h := host.Get(UserHeader); h != "" {
		user = h
	}
	if i := host.Get(IdentityHeader); i != "" {
		identity = "-i " + i
	}
	if p := host.Get(PortHeader); p != "" {
		port = "-p " + p
	}
	args := strings.Split(fmt.Sprintf("%s@%s %s %s", user, host.Get(HostHeader), identity, port), " ")
	args = slices.DeleteFunc(args, func(arg string) bool {
		return arg == ""
	})
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
