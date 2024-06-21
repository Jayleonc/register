package config_center

import (
	"context"
	"github.com/Jayleonc/register/internal/core/configuration"
)

type ConfigCenter interface {
	PutConfig(ctx context.Context, key string, value string) error
	GetConfig(ctx context.Context, key string) (string, error)
	DeleteConfig(ctx context.Context, key string) error
	WatchConfig(ctx context.Context, key string) (<-chan string, error)
}

type Client struct {
	configCenter ConfigCenter
}

type ClientOption func(*configuration.EtcdConfigCenterOptions)

func NewClient(opts ...ClientOption) (*Client, error) {
	options := configuration.DefaultEtcdConfigCenterOptions()
	for _, opt := range opts {
		opt(&options)
	}

	configCenter, err := configuration.NewEtcdConfigCenter(options)
	if err != nil {
		return nil, err
	}

	return &Client{
		configCenter: configCenter,
	}, nil
}

func (c *Client) PutConfig(ctx context.Context, key string, value string) error {
	return c.configCenter.PutConfig(ctx, key, value)
}

func (c *Client) GetConfig(ctx context.Context, key string) (string, error) {
	return c.configCenter.GetConfig(ctx, key)
}

func (c *Client) DeleteConfig(ctx context.Context, key string) error {
	return c.configCenter.DeleteConfig(ctx, key)
}

func (c *Client) WatchConfig(ctx context.Context, key string) (<-chan string, error) {
	return c.configCenter.WatchConfig(ctx, key)
}
