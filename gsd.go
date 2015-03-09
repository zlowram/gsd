package gsd

import (
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

const (
	readTimeout = 5 * time.Second
	connTimeout = 5 * time.Second
)

var (
	ProxyRetryTime = 5 * time.Second
	proxyHost      string
	proxyAuth      *proxy.Auth
)

type Gsd struct {
	Ips      []string
	Ports    []string
	Services []Service
}

func NewGsd(ips []string, ports []string) *Gsd {
	return &Gsd{Ips: ips, Ports: ports}
}

func (g *Gsd) AddServices(services []Service) {
	g.Services = append(g.Services, services...)
}

func (g *Gsd) SetProxy(host string, auth *proxy.Auth) {
	proxyHost = host
	proxyAuth = auth
}

func (g *Gsd) Run(maxConn int) chan Banner {
	b := make(chan Banner)
	go g.iterateServices(b, maxConn)
	return b
}

func (g *Gsd) iterateServices(b chan<- Banner, maxConn int) {
	var wg sync.WaitGroup

	c := make(chan int, maxConn)
	for _, i := range g.Ips {
		for _, p := range g.Ports {
			for _, s := range g.Services {
				c <- 0
				wg.Add(1)
				go func(s Service, i string, p string) {
					defer wg.Done()
					b <- s.GetBanner(i, p)
					<-c
				}(s, i, p)
			}
		}
	}

	wg.Wait()
	close(b)
}

type Service interface {
	Name() string
	GetBanner(string, string) Banner
}

type Banner struct {
	Ip      string
	Port    string
	Service string
	Content string
	Error   string
}
