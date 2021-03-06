package gsd

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"time"
)

type TCPService struct {
	name string
}

func NewTCPService() *TCPService {
	return &TCPService{name: "Generic TCP"}
}

func (s *TCPService) Name() string {
	return s.name
}

func (s *TCPService) GetBanner(ip string, port string) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	d := &GsdDialer{
		&net.Dialer{
			Timeout:   connTimeout,
			KeepAlive: 0,
		},
	}

	conn, err := d.Dial("tcp", ip+":"+port)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}
	defer conn.Close()

	// Wait to receive content
	now := time.Now()
	conn.SetReadDeadline(now.Add(readTimeout))

	buff := bytes.NewBuffer(nil)
	r := bufio.NewReader(conn)

	_, err = io.Copy(buff, r)
	if err == nil || buff.Len() > 0 {
		// Got something
		banner.Content = buff.String()
		return banner
	}

	// Timeout! Lets fuzz with something random
	fuzz := make([]byte, 128)
	_, err = rand.Read(fuzz)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}
	_, err = fmt.Fprintf(conn, "%s\n\n", fuzz)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	// Wait to receive content again
	now = time.Now()
	conn.SetReadDeadline(now.Add(readTimeout))
	_, err = io.Copy(buff, r)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	banner.Content = buff.String()
	return banner
}
