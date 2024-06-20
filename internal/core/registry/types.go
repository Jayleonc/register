package registry

import (
	"context"
	"io"
)

// ServiceInstance 表示一个服务实例
type ServiceInstance struct {
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Event 表示服务实例的变更事件
type Event struct {
	Type     string
	Instance ServiceInstance
	Value    string
	Key      string
}

// EtcdRegistry 定义注册中心的接口
type Registry interface {
	Register(ctx context.Context, si ServiceInstance) error
	UnRegister(ctx context.Context, si ServiceInstance) error
	ListServices(ctx context.Context, name string) ([]ServiceInstance, error)
	Subscribe(name string) (<-chan Event, error)
	SubscribeAll() (<-chan Event, error)
	io.Closer
}
