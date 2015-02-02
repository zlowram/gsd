package gsd

import (
	"bufio"
	"net"
	"strconv"
	"time"
)

type SSHService struct {
	name string
}

func NewSSHService() *SSHService {
	return &SSHService{name: "SSH"}
}

func (s *SSHService) Name() string {
	return s.name
}

func (s *SSHService) GetBanner(ip string, port int) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	conn, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(port), 5*time.Second)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}
	now := time.Now()
	conn.SetDeadline(now.Add(5 * time.Second))

	resp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	banner.Content = string(resp)

	return banner
}
