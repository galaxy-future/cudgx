package gateway

import (
	"encoding/json"
	"github.com/galaxy-future/cudgx/common/kafka"
	"github.com/galaxy-future/cudgx/internal/clients"
	"github.com/galaxy-future/cudgx/internal/gateway/rule"
	"io/ioutil"
	"os"
	"strings"
)

var g *Gateway

type Gateway struct {
	entriesConfig    *Config
	monitoringClient map[string]*KafkaClient
	streamingClient  map[string]*KafkaClient
	ruleManager      *rule.Manager
}

type Config struct {
	MonitoringRoute *MessageRouteConfig   `json:"monitoring_route"`
	StreamingRoute  *MessageRouteConfig   `json:"streaming_route"`
	Producer        *kafka.ProducerConfig `json:"producer"`
	Database        *rule.MysqlOption     `json:"database"`
}

type MessageRouteConfig struct {
	Entries []*StorageEntryConfig `json:"entries"`
	Default *StorageEntryConfig   `json:"default"`
}

type StorageEntryConfig struct {
	ServicePrefix string   `json:"service_prefix"`
	Brokers       []string `json:"brokers"`
	Topic         string   `json:"topic"`
}

type DatabaseConfig struct {
	Dsn            string `json:"dsn"`
	RefreshSeconds int    `json:"refresh_seconds"`
}

func GetGateway() *Gateway {
	return g
}

func Init(configFilename string) (err error) {
	g, err = NewFromConfigFile(configFilename)
	clients.InitializeBridgxClient(g.Xclient.BridgxServerAddress)
	clients.InitializeSchedulxClient(g.Xclient.SchedulxServerAddress)
	return
}

func NewFromConfigFile(fileName string) (*Gateway, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var entriesConfig Config
	err = json.Unmarshal(data, &entriesConfig)
	if err != nil {
		return nil, err
	}

	var ruleManager *rule.Manager
	if entriesConfig.Database != nil {
		var err error
		ruleManager, err = rule.NewRuleManager(entriesConfig.Database)
		if err != nil {
			return nil, err
		}
	}

	return &Gateway{
		entriesConfig:    &entriesConfig,
		monitoringClient: make(map[string]*KafkaClient),
		streamingClient:  make(map[string]*KafkaClient),
		ruleManager:      ruleManager,
	}, nil
}

func (gateway *Gateway) GetConfig() *Config {
	return gateway.entriesConfig
}

func uniqueKey(serviceName, metricName string) string {
	return strings.Join([]string{serviceName, metricName}, "-")
}
