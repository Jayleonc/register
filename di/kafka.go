package di

import (
	"Jayleonc/gateway/internal/events"
	"Jayleonc/gateway/internal/events/demo"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafkaSaramaClient() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

// RegisterConsumers 注册 DemoConsumer
func RegisterConsumers(demo *demo.Consumer) []events.Consumer {
	return []events.Consumer{demo}
}
