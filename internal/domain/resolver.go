package domain

// Resolver 定义解析器接口
type Resolver interface {
	Resolve(appID string) ([]ServiceInstance, error)
}
