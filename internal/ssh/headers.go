package ssh

const (
	NameHeader Header = iota
	HostHeader
	PortHeader
	UserHeader
	IdentityHeader
)

const headersLen = 5

var HeadersStringList = [headersLen]string{"Name", "Host", "Port", "User", "Identity"}

type Header uint8

func (h Header) String() string {
	return HeadersStringList[h]
}
