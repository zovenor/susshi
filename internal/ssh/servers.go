package ssh

import (
	"iter"
	"strconv"
	"strings"

	"github.com/mikkeloscar/sshconfig"
)

type Host [headersLen]string

func (h Host) Get(header Header) string {
	return h[header]
}

func (h *Host) Iter() iter.Seq2[Header, string] {
	return func(next func(Header, string) bool) {
		for i := range headersLen {
			if !next(Header(i), h[i]) {
				return
			}
		}
	}
}

func LoadHosts(fp string) ([]*Host, error) {
	parsedHosts, err := sshconfig.Parse(fp)
	if err != nil {
		return nil, err
	}

	var hosts []*Host
	for _, ph := range parsedHosts {
		if ph.HostName == "" {
			continue
		}
		if ph.Port == 0 {
			ph.Port = 22
		}

		var host Host
		host[NameHeader] = strings.Join(ph.Host, " ")
		host[HostHeader] = ph.HostName
		host[PortHeader] = strconv.Itoa(ph.Port)
		host[UserHeader] = ph.User
		host[IdentityHeader] = ph.IdentityFile
		hosts = append(hosts, &host)
	}

	return hosts, nil
}
