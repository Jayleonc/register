package registry

import (
	"context"
	"errors"
	"github.com/Jayleonc/register/internal/core/registry"
	"net/http"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Server struct {
	name            string
	address         string
	registry        *registry.EtcdRegistry
	registryTimeout time.Duration
	*http.Server
}

type Option func(*Server)

// MustWithAddress 用于注册服务的地址
// 如果没有调用这个方法，就要调用一次 UpdateAddress 方法
func MustWithAddress(address string) Option {
	return func(s *Server) {
		s.address = address
	}
}

func WithRegistry(client *clientv3.Client) Option {
	return func(s *Server) {
		s.registry = registry.NewRegistry(client)
	}
}

func WithRegistryTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.registryTimeout = timeout
	}
}

func WithHTTPServer(httpServer *http.Server) Option {
	return func(s *Server) {
		s.Server = httpServer
	}
}

func NewServer(name string, opts ...Option) *Server {
	server := &Server{
		name:            name,
		registryTimeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Register() error {
	if s.registry != nil {
		if s.address == "" {
			return errors.New("address is empty")
		}
		ctx, cancel := context.WithTimeout(context.Background(), s.registryTimeout)
		defer cancel()
		serviceInstance := registry.ServiceInstance{
			Name:     s.name,
			Address:  s.address,
			Metadata: map[string]string{},
		}
		err := s.registry.Register(ctx, serviceInstance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Close() error {
	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.registryTimeout)
		defer cancel()
		err := s.registry.UnRegister(ctx, registry.ServiceInstance{
			Name:    s.name,
			Address: s.address,
		})
		return err
	}

	return nil
}

func (s *Server) SubscribeAllServices() (<-chan registry.Event, error) {
	return s.registry.SubscribeAll()
}

func (s *Server) UpdateAddress(address string) {
	s.address = address
}
