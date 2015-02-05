package gsd

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type TCPTLSService struct {
	name string
}

func NewTCPTLSService() *TCPTLSService {
	return &TCPTLSService{name: "Generic TCP TLS"}
}

func (s *TCPTLSService) Name() string {
	return s.name
}

func (s *TCPTLSService) GetBanner(ip string, port int) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	// Connect
	dialer := &net.Dialer{Timeout: CONN_TIMEOUT * time.Second}
	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.DialWithDialer(dialer, "tcp", ip+":"+strconv.Itoa(port), config)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	// Get the server certificate and base64 it
	state := conn.ConnectionState()
	rawCert := state.PeerCertificates[0].Raw
	b64Cert := base64.StdEncoding.EncodeToString(rawCert)
	banner.Content = "---- BEGIN CERTIFICATE ----\n" + b64Cert +
		"---- END CERTIFICATE ----\n"

	// Wait to receive content
	now := time.Now()
	conn.SetReadDeadline(now.Add(READ_TIMEOUT * time.Second))

	buff := bytes.NewBuffer(nil)
	r := bufio.NewReader(conn)

	_, err = io.Copy(buff, r)
	if err == nil || buff.Len() > 0 {
		// Got something
		banner.Content += buff.String()
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
	conn.SetReadDeadline(now.Add(READ_TIMEOUT * time.Second))
	_, err = io.Copy(buff, r)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	banner.Content += buff.String()
	return banner
}
