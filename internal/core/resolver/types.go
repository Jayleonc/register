package resolver

import (
	"context"
	"git.daochat.cn/service/registry/internal/core/registry"
)

// Resolver 接口定义了解析服务实例位置的方法
type Resolver interface {
	Resolve(ctx context.Context, name string) ([]registry.ServiceInstance, error)
}
