package registry

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
)

type EtcdRegistry struct {
	client *clientv3.Client
	sess   *concurrency.Session
}

func NewRegistry(client *clientv3.Client) *EtcdRegistry {
	sess, err := concurrency.NewSession(client, concurrency.WithTTL(5))
	if err != nil {
		panic(err)
	}

	return &EtcdRegistry{
		client: client,
		sess:   sess,
	}
}

// Register 在 etcd 中注册服务实例。
// 使用 etcd 会话（session）机制确保服务节点崩溃时数据能够自动删除。
// 服务实例数据以 JSON 格式存储，并绑定会话租约。
func (r *EtcdRegistry) Register(ctx context.Context, si ServiceInstance) error {
	// 获取会话租约 ID
	leaseID := r.sess.Lease()

	// 构建服务实例的键
	key := r.instanceKey(si)

	// 将服务实例数据序列化为 JSON
	val, err := json.Marshal(si)
	if err != nil {
		return fmt.Errorf("failed to marshal service instance: %w", err)
	}

	// 将服务实例数据写入 etcd，并绑定会话租约 ID
	_, err = r.client.Put(ctx, key, string(val), clientv3.WithLease(leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// 会话管理会自动处理租约的续约
	return nil
}

// UnRegister 从 etcd 中注销服务实例。
// 根据服务实例的键删除相应的数据。
func (r *EtcdRegistry) UnRegister(ctx context.Context, si ServiceInstance) error {
	// 构建服务实例的键
	key := r.instanceKey(si)

	// 从 etcd 中删除服务实例数据
	_, err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to unregister service: %w", err)
	}
	log.Printf("Service unregistered: %s at %s", si.Name, si.Address)
	return nil
}

// ListServices 列出指定服务的所有实例
func (r *EtcdRegistry) ListServices(ctx context.Context, name string) ([]ServiceInstance, error) {
	key := r.serviceKey(ServiceInstance{Name: name})
	resp, err := r.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get service instances: %w", err)
	}

	var instances []ServiceInstance
	for _, kv := range resp.Kvs {
		var instance ServiceInstance
		if err := json.Unmarshal(kv.Value, &instance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal service instance: %w", err)
		}
		instances = append(instances, instance)
	}

	log.Printf("Listed services for: %s, instances found: %d", name, len(instances))
	return instances, nil
}

// Subscribe 订阅指定服务的变更事件
func (r *EtcdRegistry) Subscribe(name string) (<-chan Event, error) {
	ch := make(chan Event)
	go func() {
		defer close(ch)
		key := r.serviceKey(ServiceInstance{Name: name})
		rch := r.client.Watch(context.Background(), key, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				var si ServiceInstance
				if err := json.Unmarshal(ev.Kv.Value, &si); err != nil {
					log.Printf("Failed to unmarshal service instance: %v", err)
					continue
				}
				eventType := "PUT"
				if ev.Type == clientv3.EventTypeDelete {
					eventType = "DELETE"
				}
				log.Printf("Service event: %s - %s at %s", eventType, si.Name, si.Address)
				ch <- Event{
					Type:     eventType,
					Instance: si,
				}
			}
		}
	}()
	return ch, nil
}

func (r *EtcdRegistry) SubscribeAll() (<-chan Event, error) {
	ch := make(chan Event)
	go func() {
		defer close(ch)
		rch := r.client.Watch(context.Background(), "/services/", clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				var si ServiceInstance
				if err := json.Unmarshal(ev.Kv.Value, &si); err != nil {
					continue
				}
				eventType := "PUT"
				if ev.Type == clientv3.EventTypeDelete {
					eventType = "DELETE"
				}
				// 只发送一次事件类型
				if ev.Type == clientv3.EventTypePut {
					eventType = "REGISTER"
				}
				ch <- Event{
					Type:     eventType,
					Instance: si,
				}
			}
		}
	}()
	return ch, nil
}

// Close 关闭 sess 会话。
// 如果关闭会话失败，返回相应的错误。
func (r *EtcdRegistry) Close() error {
	return r.sess.Close()
}

// instanceKey 返回服务实例在 etcd 中的键。
func (r *EtcdRegistry) instanceKey(si ServiceInstance) string {
	return fmt.Sprintf("/services/%s/%s", si.Name, si.Address)
}

// serviceKey 返回服务在 etcd 中的键。
func (r *EtcdRegistry) serviceKey(si ServiceInstance) string {
	return fmt.Sprintf("/services/%s", si.Name)
}
