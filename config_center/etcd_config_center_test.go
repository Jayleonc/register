package config_center

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewConfigCenter(t *testing.T) {

	ctx := context.Background()

	client, err := NewClient(
		WithEtcdAddress("localhost:2379"),
		//WithCredentials("user", "password"),
		WithDialTimeout(10*time.Second),
		//WithLogLevel("debug"),
		//WithTLS("cert.pem", "key.pem", "ca.pem"),
		//WithRetry(5, 2*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create config center client: %v", err)
	}

	err = client.PutConfig(ctx, "user:status:1", "active")
	if err != nil {
		log.Fatalf("Failed to put config: %v", err)
	}

	value, err := client.GetConfig(ctx, "user:status:1")
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}
	fmt.Println("Config value:", value)
}
