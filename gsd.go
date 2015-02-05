package gsd

const (
	READ_TIMEOUT = 5
	CONN_TIMEOUT = 5
)

type Gsd struct {
	Ips      []string
	Ports    []int
	Services []Service
}

func NewGsd(ips []string, ports []int) *Gsd {
	return &Gsd{Ips: ips, Ports: ports}
}

func (g *Gsd) AddServices(services []Service) {
	g.Services = appendSlice(g.Services, services)
}

func (g *Gsd) Run() []Banner {
	b := make(chan Banner)
	for _, i := range g.Ips {
		for _, p := range g.Ports {
			for _, s := range g.Services {
				go func(s Service, i string, p int) {
					b <- s.GetBanner(i, p)
				}(s, i, p)
			}
		}
	}
	banners := make([]Banner, 0)
	for i := 0; i < len(g.Ips)*len(g.Ports)*len(g.Services); i++ {
		banners = append(banners, <-b)
	}
	return banners
}

func appendSlice(a []Service, b []Service) []Service {
	for _, i := range b {
		a = append(a, i)
	}
	return a
}

type Service interface {
	Name() string
	GetBanner(string, int) Banner
}

type Banner struct {
	Ip      string
	Port    int
	Service string
	Content string
	Error   string
}
