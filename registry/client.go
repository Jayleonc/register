package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Jayleonc/register/internal/core/resolver"
)

type ClientOption func(c *Client)

type Client struct {
	insecure bool
	timeout  time.Duration
	resolver resolver.Resolver
	*http.Client
}

func NewClient(opts ...ClientOption) (*Client, error) {
	res := &Client{
		Client: &http.Client{},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func ClientInsecure() ClientOption {
	return func(c *Client) {
		c.insecure = true
	}
}

func ClientWithResolver(res resolver.Resolver) ClientOption {
	return func(c *Client) {
		c.resolver = res
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.Client = client
	}
}

func (c *Client) Dial(ctx context.Context, service string) (*http.Client, error) {
	if c.resolver != nil {
		ctx, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()
		instances, err := c.resolver.Resolve(ctx, service)
		if err != nil {
			return nil, err
		}

		if len(instances) == 0 {
			return nil, fmt.Errorf("no instances found for service: %s", service)
		}

		// 这里选择一个实例作为示例，可以根据实际需求实现负载均衡策略
		selectedInstance := instances[0]
		c.Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   selectedInstance.Address,
			}),
		}
	}

	return c.Client, nil
}

func (c *Client) GetServiceInterfaces(ctx context.Context, service string) ([]ServiceInterface, error) {
	if c.resolver == nil {
		return nil, fmt.Errorf("resolver is not set")
	}

	instances, err := c.resolver.Resolve(ctx, service)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances found for service: %s", service)
	}

	// 选择第一个实例，获取其接口描述
	selectedInstance := instances[0]
	var interfaces []ServiceInterface
	err = json.Unmarshal([]byte(selectedInstance.Metadata["interfaces"]), &interfaces)
	if err != nil {
		return nil, err
	}

	return interfaces, nil
}
