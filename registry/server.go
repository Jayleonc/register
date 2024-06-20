package sdk

import (
	"Jayleonc/register/internal/core/registry"
	"Jayleonc/register/internal/core/resolver"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Server struct {
	name            string
	registry        *registry.EtcdRegistry
	registryTimeout time.Duration
	listener        net.Listener
	resolver        resolver.Resolver
	*http.Server
}

type Option func(*Server)

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

func WithResolver(res resolver.Resolver) Option {
	return func(s *Server) {
		s.resolver = res
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
		Server:          &http.Server{},
		registryTimeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Start(addr string, interfaceBuilder *InterfaceBuilder) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = listener

	// 构建接口信息
	interfaces := interfaceBuilder.GetInterfaces()
	interfaceData, err := json.Marshal(interfaces)
	if err != nil {
		return err
	}

	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.registryTimeout)
		defer cancel()
		serviceInstance := registry.ServiceInstance{
			Name:    s.name,
			Address: s.listener.Addr().String(),
			Metadata: map[string]string{
				"interfaces": string(interfaceData),
			},
		}
		err = s.registry.Register(ctx, serviceInstance)
		if err != nil {
			return err
		}
		defer func() {
			_ = s.registry.Close()
		}()
	}

	s.Addr = addr
	err = s.Serve(listener)
	return err
}

func (s *Server) Stop() error {
	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.registryTimeout)
		defer cancel()
		err := s.registry.UnRegister(ctx, registry.ServiceInstance{
			Name:    s.name,
			Address: s.listener.Addr().String(),
		})
		return err
	}

	_ = s.Shutdown(context.Background())
	_ = s.listener.Close()
	return nil
}

func (s *Server) SubscribeAllServices() (<-chan registry.Event, error) {
	return s.registry.SubscribeAll()
}
