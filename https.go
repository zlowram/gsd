package gsd

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
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

func (s *HttpsService) GetBanner(ip string, port int) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Timeout: 5 * time.Second, Transport: tr}
	res, err := c.Get("https://" + ip + ":" + strconv.Itoa(port))
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	banner.Content = string(dump)

	return banner
}
