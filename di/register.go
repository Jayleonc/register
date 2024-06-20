package di

import (
	"Jayleonc/register/internal/core/registry"
	"Jayleonc/register/internal/core/resolver"
	"Jayleonc/register/sdk"
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

func provideRegisterServer(client *clientv3.Client) *sdk.Server {
	return sdk.NewServer("my-service", sdk.WithRegistry(client), sdk.WithResolver(resolver.NewEtcdResolver(client)))
}
