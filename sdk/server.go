package sdk

import (
	"context"
	"encoding/json"
	"github.com/Jayleonc/register/internal/core/registry"
	"github.com/Jayleonc/register/internal/core/resolver"
	"net"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Server struct {
	name            string
	registry        *registry.EtcdRegistry
	registryTimeout time.Duration
	listener        net.Listener
	resolver        resolver.Resolver
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

func WithListener(listener net.Listener) Option {
	return func(s *Server) {
		s.listener = listener
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

func (s *Server) Register(interfaceBuilder *InterfaceBuilder) error {
	// 使用已经传入的 listener

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
			Address: s.listener.Addr().String(), // 使用传入的 listener
			Metadata: map[string]string{
				"interfaces": string(interfaceData),
			},
		}
		err = s.registry.Register(ctx, serviceInstance)
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
			Address: s.listener.Addr().String(),
		})
		return err
	}

	_ = s.listener.Close()
	return nil
}

func (s *Server) SubscribeAllServices() (<-chan registry.Event, error) {
	return s.registry.SubscribeAll()
}
