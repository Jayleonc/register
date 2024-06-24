package configuration

import (
	"context"
)

type ConfigCenter interface {
	PutConfig(ctx context.Context, key string, value string) error
	GetConfig(ctx context.Context, key string) (string, error)
	DeleteConfig(ctx context.Context, key string) error
	WatchConfig(ctx context.Context, key string) (<-chan string, error)
	ListConfig(ctx context.Context, prefix string) (map[string]string, error)
}
