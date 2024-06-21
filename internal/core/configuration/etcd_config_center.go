package configuration

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfigCenter struct {
	client *clientv3.Client
}

type EtcdConfigCenterOptions struct {
	EtcdAddresses []string
	Username      string
	Password      string
	DialTimeout   time.Duration
	LogLevel      string
	TLSCertFile   string
	TLSKeyFile    string
	TLSCaFile     string
	RetryAttempts int
	RetryInterval time.Duration
}

func DefaultEtcdConfigCenterOptions() EtcdConfigCenterOptions {
	return EtcdConfigCenterOptions{
		EtcdAddresses: []string{"localhost:2379"},
		DialTimeout:   5 * time.Second,
		LogLevel:      "info",
		RetryAttempts: 3,
		RetryInterval: 1 * time.Second,
	}
}

func NewEtcdConfigCenter(opts EtcdConfigCenterOptions) (*EtcdConfigCenter, error) {
	clientConfig := clientv3.Config{
		Endpoints:   opts.EtcdAddresses,
		DialTimeout: opts.DialTimeout,
	}

	if opts.Username != "" && opts.Password != "" {
		clientConfig.Username = opts.Username
		clientConfig.Password = opts.Password
	}

	if opts.TLSCertFile != "" && opts.TLSKeyFile != "" && opts.TLSCaFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(opts.TLSCertFile, opts.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load key pair: %w", err)
		}
		caCertPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile(opts.TLSCaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert file: %w", err)
		}
		caCertPool.AppendCertsFromPEM(caCert)

		clientConfig.TLS = &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			RootCAs:      caCertPool,
		}
	}

	client, err := clientv3.New(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	return &EtcdConfigCenter{client: client}, nil
}

func (c *EtcdConfigCenter) PutConfig(ctx context.Context, key string, value string) error {
	_, err := c.client.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to put config: %w", err)
	}
	log.Printf("Config put: %s = %s", key, value)
	return nil
}

func (c *EtcdConfigCenter) GetConfig(ctx context.Context, key string) (string, error) {
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("config not found for key: %s", key)
	}
	value := string(resp.Kvs[0].Value)
	log.Printf("Config get: %s = %s", key, value)
	return value, nil
}

func (c *EtcdConfigCenter) DeleteConfig(ctx context.Context, key string) error {
	_, err := c.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}
	log.Printf("Config deleted: %s", key)
	return nil
}

func (c *EtcdConfigCenter) WatchConfig(ctx context.Context, key string) (<-chan string, error) {
	ch := make(chan string)
	go func() {
		defer close(ch)
		rch := c.client.Watch(ctx, key)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == clientv3.EventTypePut {
					ch <- string(ev.Kv.Value)
					log.Printf("Config updated: %s = %s", key, string(ev.Kv.Value))
				} else if ev.Type == clientv3.EventTypeDelete {
					ch <- ""
					log.Printf("Config deleted: %s", key)
				}
			}
		}
	}()
	return ch, nil
}
