package gsd

import (
	"crypto/tls"
	"encoding/base64"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
)

type HttpsService struct {
	name string
}

func NewHttpsService() *HttpsService {
	return &HttpsService{name: "HTTPS"}
}

func (s *HttpsService) Name() string {
	return s.name
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

	userAgents := []string{
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.0.6.01001)",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0",
		"Mozilla/5.0 (X11; U; Linux x86_64; de; rv:1.9.2.8) Gecko/20100723 Ubuntu/10.04 (lucid) Firefox/3.6.8",
		"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.10.289 Version/12.01",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) AppleWebKit/600.3.18 (KHTML, like Gecko) Version/8.0.3 Safari/600.3.18",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B466 Safari/600.1.4",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/600.3.18 (KHTML, like Gecko) Version/7.1.3 Safari/537.85.12",
		"Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		"Mozilla/6.0 (iPhone; U; CPU like Mac OS X; en) AppleWebKit/420+ (KHTML, like Gecko) Version/3.0 Mobile/1A543a Safari/419.3",
	}

	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

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
