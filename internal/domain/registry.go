package domain

// ServiceInstance 定义服务实例信息
type ServiceInstance struct {
	InstanceID string            `json:"instanceId"`
	ServiceURL string            `json:"serviceUrl"`
	Metadata   map[string]string `json:"metadata"`
}

// Registry 定义注册中心接口
type Registry interface {
	Register(instance ServiceInstance) error
	Unregister(appID string, instanceID string) error
	Discover(appID string) ([]ServiceInstance, error)
	HealthCheck() error
}

// HealthChecker 定义健康检查接口
type HealthChecker interface {
	Check(instance ServiceInstance) bool
}
