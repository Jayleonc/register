package configuration

import (
	"context"
	"fmt"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfigCenter struct {
	client *clientv3.Client
}

func NewEtcdConfigCenter(client *clientv3.Client) *EtcdConfigCenter {
	return &EtcdConfigCenter{client: client}
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
