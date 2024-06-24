package di

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func InitEtcdClient() *clientv3.Client {
	etcdAddresses := viper.GetStringSlice("etcd.addresses")
	if len(etcdAddresses) == 0 {
		log.Fatalf("No etcd addresses found in configuration")
	}

	// 初始化 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddresses,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return cli
}
