package gsd

import (
	"net/http"
	"net/http/httputil"
	"time"
)

type HttpService struct {
	name string
}

func NewHttpService() *HttpService {
	return &HttpService{name: "HTTP"}
}

func (s *HttpService) Name() string {
	return s.name
}

func (s *HttpService) GetBanner(ip string, port string) Banner {
	banner := Banner{
		Ip:      ip,
		Port:    port,
		Service: s.Name(),
	}

	c := &http.Client{Timeout: 5 * time.Second}
	res, err := c.Get("http://" + ip + ":" + port)
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
