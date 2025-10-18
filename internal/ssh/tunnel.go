package ssh

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// NewTunnelDialer creates a new SSH tunnel dialer
func NewTunnelDialer(host string, port int, user, pass string) (func(network, addr string) (net.Conn, error), func() error, error) {
	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: replace with known_hosts for production
		Timeout:         5 * time.Second,
	}
	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), cfg)
	if err != nil {
		return nil, nil, err
	}
	return func(network, addr string) (net.Conn, error) { return c.Dial("tcp", addr) }, c.Close, nil
}
