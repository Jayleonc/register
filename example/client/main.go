package main

import (
	"context"
	"git.daochat.cn/service/registry/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}

	sdkClient, err := registry.NewClient(
		registry.ClientWithResolver(registry.NewEtcdResolver(client)),
	)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx := context.Background()
	httpClient, err := sdkClient.Dial(ctx, "user_service")
	if err != nil {
		log.Fatalf("Failed to dial service: %v", err)
	}

	// 使用 httpClient 进行服务调用
	resp, err := httpClient.Get("http://user_service/health")
	if err != nil {
		log.Fatalf("Failed to call service: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	log.Println("Service response status:", resp.Status)
}
