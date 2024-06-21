package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"git.daochat.cn/service/registry/internal/core/resolver"
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

		// 应接入负债均衡器
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
