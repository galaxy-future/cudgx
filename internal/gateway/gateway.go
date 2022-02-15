package gateway

import (
	"errors"
	"fmt"
	"strings"

	"github.com/galaxy-future/cudgx/common/kafka"
	"github.com/galaxy-future/cudgx/common/mod"
	wrapmod "github.com/galaxy-future/cudgx/internal/gateway/mod"
)

var ErrWrongParams = errors.New("wrong params")

type KafkaClient struct {
	producer *kafka.ProducerClient
	topic    string
}

func (gateway *Gateway) GetMonitoringStorageConfigEntry(serviceName, metricName string) *StorageEntryConfig {

	resultEntry := gateway.GetConfig().MonitoringRoute.Default
	for _, entry := range gateway.GetConfig().MonitoringRoute.Entries {
		if entry.ServicePrefix == serviceName {
			resultEntry = entry
			break
		}
		if strings.Contains(serviceName, entry.ServicePrefix+".") && len(entry.ServicePrefix) > len(resultEntry.ServicePrefix) {
			resultEntry = entry
		}
	}
	return resultEntry

}

func (gateway *Gateway) GetStreamingStorageConfigEntry(serviceName, metricName string) *StorageEntryConfig {

	resultEntry := gateway.GetConfig().StreamingRoute.Default
	for _, entry := range gateway.GetConfig().StreamingRoute.Entries {
		if entry.ServicePrefix == serviceName {
			resultEntry = entry
			break
		}
		if strings.Contains(serviceName, entry.ServicePrefix+".") && len(entry.ServicePrefix) > len(resultEntry.ServicePrefix) {
			resultEntry = entry
		}
	}
	return resultEntry

}

func (g *KafkaClient) SendMessage(serviceName, metricName string, data []byte) {
	g.producer.SendMessage(g.topic, metricName, data)
}

func (gateway *Gateway) GetMonitoringWriter(serviceName, metricName string) (*KafkaClient, error) {
	configEntry := gateway.GetMonitoringStorageConfigEntry(serviceName, metricName)
	writer, exists := gateway.monitoringClient[configEntry.ServicePrefix]
	if exists {
		return writer, nil
	}

	client, err := kafka.NewProducer(configEntry.Brokers, gateway.GetConfig().Producer)
	if err != nil {
		return nil, err
	}

	writer = &KafkaClient{
		producer: client,
		topic:    configEntry.Topic,
	}

	gateway.monitoringClient[configEntry.ServicePrefix] = writer

	return writer, nil
}

func (gateway *Gateway) GetStreamingWriter(serviceName, metricName string) (*KafkaClient, error) {
	configEntry := gateway.GetStreamingStorageConfigEntry(serviceName, metricName)
	writer, exists := gateway.streamingClient[configEntry.ServicePrefix]
	if exists {
		return writer, nil
	}

	client, err := kafka.NewProducer(configEntry.Brokers, gateway.GetConfig().Producer)
	if err != nil {
		return nil, err
	}

	writer = &KafkaClient{
		producer: client,
		topic:    configEntry.Topic,
	}

	gateway.streamingClient[configEntry.ServicePrefix] = writer

	return writer, nil
}

func (gateway *Gateway) WrapStreamingMessage(streamingBatch *mod.StreamingBatch) (batch *wrapmod.StreamingRuleBatch, err error) {
	if gateway.ruleManager != nil {
		return gateway.ruleManager.WrapStreamingMessage(streamingBatch)
	} else {
		return nil, fmt.Errorf(" RuleManager is not ready! ")
	}
}
