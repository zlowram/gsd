package gsd

import (
	"crypto/tls"
	"encoding/base64"
	"net"
	"net/http"
	"net/http/httputil"
)

type HttpsService struct {
	name   string
	header *http.Header
}

func NewHttpsService() *HttpsService {
	return &HttpsService{name: "HTTPS", header: &http.Header{}}
}

func (s *HttpsService) Name() string {
	return s.name
}

func (s *HttpsService) SetHeader(key string, value string) {
	s.header.Add(key, value)
}

func (s *HttpsService) GetBanner(ip string, port string) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	// Connect and make a GET request
	tr := &http.Transport{
		Dial: (&GsdDialer{
			Dialer: &net.Dialer{
				Timeout:   connTimeout,
				KeepAlive: 0,
			},
		}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &http.Client{Timeout: connTimeout, Transport: tr}

	req, err := http.NewRequest("GET", "http://"+ip+":"+port, nil)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	req.Header = *s.header

	res, err := c.Do(req)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}
	defer res.Body.Close()

	// Check if the connection is encrypted, get the certificate and base64 it
	if res.TLS == nil {
		banner.Error = "Non encrypted HTTP connection"
		return banner
	}

	rawCert := res.TLS.PeerCertificates[0].Raw
	b64Cert := base64.StdEncoding.EncodeToString(rawCert)

	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	banner.Content = "-----BEGIN CERTIFICATE-----\n" + b64Cert +
		"-----END CERTIFICATE-----\n\n" + string(dump)

	return banner
}
