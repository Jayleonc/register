package registry

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type ConfigCenter struct {
	client *clientv3.Client
}

func NewConfigCenter(client *clientv3.Client) *ConfigCenter {
	return &ConfigCenter{client: client}
}

func (cc *ConfigCenter) PutConfig(ctx context.Context, key, value string) error {
	_, err := cc.client.Put(ctx, key, value)
	if err != nil {
		return err
	}
	log.Printf("Config put: %s = %s", key, value)
	return nil
}

func (cc *ConfigCenter) GetConfig(ctx context.Context, key string) (string, error) {
	resp, err := cc.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	value := string(resp.Kvs[0].Value)
	log.Printf("Config get: %s = %s", key, value)
	return value, nil
}

func (cc *ConfigCenter) WatchConfig(ctx context.Context, key string) (<-chan string, error) {
	ch := make(chan string)
	go func() {
		rch := cc.client.Watch(ctx, key)
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
		close(ch)
	}()
	return ch, nil
}
