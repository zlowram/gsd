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

func (s *TCPTLSService) GetBanner(ip string, port string) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	// Connect
	dialer := &GsdDialer{
		Dialer: &net.Dialer{
			Timeout:   connTimeout,
			KeepAlive: 0,
		},
	}
	config := &tls.Config{InsecureSkipVerify: true}
	dconn, err := dialer.Dial("tcp", ip+":"+port)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}
	defer dconn.Close()

	conn := tls.Client(dconn, config)
	defer conn.Close()

	// Check if the connection is encrypted, get the certificate and base64 it
	state := conn.ConnectionState()
	if state.PeerCertificates == nil {
		banner.Error = "Non encrypted HTTP connection"
		return banner
	}
	rawCert := state.PeerCertificates[0].Raw
	b64Cert := base64.StdEncoding.EncodeToString(rawCert)
	banner.Content = "-----BEGIN CERTIFICATE-----\n" + b64Cert +
		"-----END CERTIFICATE-----\n"

	// Wait to receive content
	now := time.Now()
	conn.SetReadDeadline(now.Add(readTimeout))

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
	conn.SetReadDeadline(now.Add(readTimeout))
	_, err = io.Copy(buff, r)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	banner.Content += buff.String()
	return banner
}
