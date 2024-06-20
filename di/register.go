package di

import (
	"github.com/Jayleonc/register/internal/core/registry"
	"github.com/Jayleonc/register/internal/core/resolver"
	registry2 "github.com/Jayleonc/register/registry"
	"github.com/Jayleonc/register/sdk"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func provideEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	})
}

func provideRegistry(client *clientv3.Client) *registry.EtcdRegistry {
	return registry.NewRegistry(client)
}

func provideResolver(client *clientv3.Client) resolver.Resolver {
	return resolver.NewEtcdResolver(client)
}

func provideRegisterServer(client *clientv3.Client) *registry2.Server {
	return registry2.NewServer("my-service", registry2.WithRegistry(client), sdk.WithResolver(resolver.NewEtcdResolver(client)))
}
