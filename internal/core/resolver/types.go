package resolver

import (
	"Jayleonc/register/internal/core/registry"
	"context"
)

// Resolver 接口定义了解析服务实例位置的方法
type Resolver interface {
	Resolve(ctx context.Context, name string) ([]registry.ServiceInstance, error)
}
