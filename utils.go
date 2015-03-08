package gsd

import (
	"log"
	"net"

	"golang.org/x/net/proxy"
)

type GsdDialer struct {
	*net.Dialer
}

func (d *GsdDialer) Dial(network, address string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   d.Timeout,
		KeepAlive: d.KeepAlive,
	}

	if proxyHost != "" {
		// Check if proxy is there
		conn, err := net.Dial(network, proxyHost)
		if err != nil {
			log.Fatal(err)
		}

		// Proxy is there, connect
		proxyDialer, err := proxy.SOCKS5(network, proxyHost, proxyAuth, dialer)
		if err != nil {
			log.Fatal(err)
		}
		conn, err = proxyDialer.Dial(network, address)
		return conn, err
	} else {
		conn, err := dialer.Dial(network, address)
		return conn, err
	}
}
