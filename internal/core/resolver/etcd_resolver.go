package resolver

import (
	"Jayleonc/register/internal/core/registry"
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdResolver 使用 etcd 作为服务解析器的实现
type EtcdResolver struct {
	client *clientv3.Client
}

// NewEtcdResolver 创建一个新的 EtcdResolver 实例
func NewEtcdResolver(client *clientv3.Client) *EtcdResolver {
	return &EtcdResolver{client: client}
}

// Resolve 从 etcd 中解析服务实例的位置
func (r *EtcdResolver) Resolve(ctx context.Context, name string) ([]registry.ServiceInstance, error) {
	key := fmt.Sprintf("/services/%s", name)
	resp, err := r.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get service instances: %w", err)
	}

	var instances []registry.ServiceInstance
	for _, kv := range resp.Kvs {
		var instance registry.ServiceInstance
		if err := json.Unmarshal(kv.Value, &instance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal service instance: %w", err)
		}
		instances = append(instances, instance)
	}

	return instances, nil
}
