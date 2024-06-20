package sdk

import (
	"Jayleonc/register/internal/core/registry"
	"Jayleonc/register/internal/core/resolver"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Resolver interface {
	Resolve(ctx context.Context, name string) ([]registry.ServiceInstance, error)
}

type etcdResolverDecorator struct {
	internalResolver resolver.Resolver
}

func NewEtcdResolver(client *clientv3.Client) Resolver {
	return &etcdResolverDecorator{
		internalResolver: resolver.NewEtcdResolver(client),
	}
}

func (r *etcdResolverDecorator) Resolve(ctx context.Context, name string) ([]registry.ServiceInstance, error) {
	return r.internalResolver.Resolve(ctx, name)
}
