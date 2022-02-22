package consumer

import (
	"encoding/json"

	"github.com/galaxy-future/cudgx/common/kafka"
	"github.com/galaxy-future/cudgx/common/victoriametrics"
)

func LoadConfig(data []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(data, &config)
	return &config, err
}

//Config 是consumer的配置
type Config struct {
	//Kafka 配置
	Kafka *KafkaConfig `json:"kafka"`
	//VictoriaMetrics 连接配置
	VictoriaMetrics *victoriametrics.Config `json:"victoria_metrics"`
}

//KafkaConfig 消费程序用到kafka的配置
type KafkaConfig struct {
	//Brokers kafka brokers
	Brokers []string `json:"brokers"`
	//Group kafka 消费Group
	Group string `json:"group"`
	//Topic 消费Topic
	Topic string `json:"topic"`
	//Consumer consumer配置
	Consumer *kafka.ConsumerConfig `json:"consumer"`
}
