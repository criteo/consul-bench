package bench

import (
	"net"
	"strconv"

	consul "github.com/hashicorp/consul/api"
)

type Config struct {
	ConsulHost string
	HTTPPort   int
	RPCPort    int
	ACLToken   string
}

type Bench struct {
	cfg        Config
	httpClient *consul.Client
	stats      chan Stat
}

func New(cfg Config) (*Bench, error) {
	c, err := consul.NewClient(&consul.Config{
		Address: net.JoinHostPort(cfg.ConsulHost, strconv.Itoa(cfg.HTTPPort)),
		Token:   cfg.ACLToken,
	})
	if err != nil {
		return nil, err
	}

	return &Bench{
		cfg:        cfg,
		httpClient: c,
		stats:      make(chan Stat),
	}, nil
}
