package gsd

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type HttpService struct {
	name   string
	header *http.Header
}

func NewHttpService() *HttpService {
	return &HttpService{name: "HTTP", header: &http.Header{}}
}

func (s *HttpService) Name() string {
	return s.name
}

func (s *HttpService) SetHeader(key, value string) {
	s.header.Add(key, value)
}

func (s *HttpService) GetBanner(ip string, port string) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	tr := &http.Transport{
		Dial: (&GsdDialer{
			Dialer: &net.Dialer{
				Timeout:   connTimeout,
				KeepAlive: 0,
			},
		}).Dial,
	}

	c := &http.Client{Timeout: 5 * time.Second, Transport: tr}

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

	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		banner.Error = err.Error()
		return banner
	}

	banner.Content = string(dump)

	return banner
}
