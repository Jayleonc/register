package main

import (
	"context"
	"github.com/Jayleonc/register/internal/core/resolver"
	"github.com/Jayleonc/register/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

func main() {
	// 创建 etcd 客户端和注册中心实例
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}

	// 创建 SDK 客户端
	sdkClient, err := registry.NewClient(
		registry.ClientWithResolver(resolver.NewEtcdResolver(client)),
	)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// 获取日志服务的接口描述
	ctx := context.Background()
	interfaces, err := sdkClient.GetServiceInterfaces(ctx, "user_service")
	if err != nil {
		log.Fatalf("Failed to get service interfaces: %v", err)
	}

	for _, iface := range interfaces {
		log.Printf("  interface: %s %s", iface.Method, iface.Path)
		for _, param := range iface.Params {
			log.Printf("Param: field:%s type:%s", param.Name, param.Type)
		}
		for _, ret := range iface.Returns {
			log.Printf("Return: field:%s type:%s", ret.Name, ret.Type)
		}
	}
}
