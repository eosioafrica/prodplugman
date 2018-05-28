package consul

import (
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/eosioafrica/prodplugman/ppman"
)

type Service struct {

	Name        string
	TTL         time.Duration
	PPMan 		ppman.NodeosEndpoint
	ConsulAgent *consul.Agent
}

func New(c *consul.Client,addrs string, ttl time.Duration) (*Service, error) {
	s := new(Service)
	s.Name = "webkv"
	s.TTL = ttl

	p := ppman.NodeosEndpoint{

		URL: addrs,
		Port: 8300,
	}

	s.PPMan = p

	ok, err := s.Check()
	if !ok {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	s.ConsulAgent = c.Agent()

	serviceDef := &consul.AgentServiceRegistration{
		Name: s.Name,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}

	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return nil, err
	}
	go s.UpdateTTL(s.Check)

	return s, nil
}