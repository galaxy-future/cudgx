package consumer

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"

	"github.com/galaxy-future/cudgx/common/kafka"
	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/common/mod"
	"github.com/galaxy-future/cudgx/common/victoriametrics"
	"github.com/galaxy-future/cudgx/internal/consumer/consts"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	"go.uber.org/zap"
)

type Consumer struct {
	kafkaClient           *kafka.ConsumerClient
	victoriaMetricsWriter *victoriametrics.AsyncWriter
	messageChan           chan interface{}
	config                *Config
}

func NewConsumer(config *Config) (*Consumer, error) {
	messagesCh := make(chan interface{}, 100000)
	consumer := &Consumer{
		config:      config,
		messageChan: messagesCh,
	}

	kafkaClient, err := kafka.NewConsumers(messagesCh, config.Kafka.Brokers, config.Kafka.Topic, config.Kafka.Group, config.Kafka.Consumer)
	if err != nil {
		return nil, err
	}
	consumer.kafkaClient = kafkaClient
	writer, err := victoriametrics.NewWriter(config.VictoriaMetrics, messagesCh, consumer.commit)
	if err != nil {
		return nil, err
	}
	consumer.victoriaMetricsWriter = writer
	return consumer, nil
}

func (consumer *Consumer) Start(ctx context.Context) {
	var wgKafka sync.WaitGroup
	wgKafka.Add(1)
	go func() {
		defer wgKafka.Done()
		consumer.kafkaClient.Start(ctx)
		logger.GetLogger().Info("kafka process exists")
	}()

	var wgWriter sync.WaitGroup
	wgWriter.Add(1)
	go func() {
		defer wgWriter.Done()
		consumer.victoriaMetricsWriter.Init()
	}()
	<-ctx.Done()

	wgKafka.Wait()
	consumer.kafkaClient.Stop()

	close(consumer.messageChan)

	wgWriter.Wait()

}

func (consumer *Consumer) commit(cli *http.Client, messages []interface{}) error {
	wirteRequest, err := toPromePb(messages)
	if err != nil {
		return err
	}
	if wirteRequest == nil || len(wirteRequest.Timeseries) == 0 {
		return nil
	}
	data, err := proto.Marshal(wirteRequest)
	if err != nil {
		return err
	}
	// ex: http://127.0.0.1:8480/insert/0/prometheus/api/v1/write
	httpReq, err := http.NewRequest("POST", consumer.config.VictoriaMetrics.Writer.VmUrl, bytes.NewReader(snappy.Encode(nil, data)))
	if err != nil {
		return err
	}
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	httpReq.Header.Set("Connection", "keep-alive")
	resp, err := cli.Do(httpReq)
	defer func() { _ = resp.Body.Close() }()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		logger.GetLogger().Error("remote write request failed", zap.Int("status_code", resp.StatusCode))
		return fmt.Errorf("remote write request failed,status_code:%d", resp.StatusCode)
	}
	return nil
}

func IsEqual(s, d *mod.MetricsMessage) bool {
	if s.MetricName != d.MetricName {
		return false
	}
	if s.ServiceName != d.ServiceName {
		return false
	}
	if s.ClusterName != d.ClusterName {
		return false
	}
	if s.ServiceHost != d.ServiceHost {
		return false
	}
	if s.ServiceAz != d.ServiceAz {
		return false
	}
	if s.ServiceRegion != d.ServiceRegion {
		return false
	}
	if !reflect.DeepEqual(s.Labels, d.Labels) {
		return false
	}
	return true
}

func toPromePb(messages []interface{}) (*prompb.WriteRequest, error) {
	var metricBatch *mod.MetricBatch
	var items []*prompb.TimeSeries
	var labels []*prompb.Label
	metricMap := make(map[*mod.MetricsMessage][]prompb.Sample)
	eq := false
	for _, m := range messages {
		metricBatch = &mod.MetricBatch{}
		err := proto.Unmarshal(m.([]byte), metricBatch)
		if err != nil {
			logger.GetLogger().Error("Message to PromePb failed", zap.Error(err))
			continue
		}

		for _, metric := range metricBatch.Messages {
			eq = false
			if len(metricMap) == 0 {
				metricMap[metric] = []prompb.Sample{
					{
						Value:     metric.Value,
						Timestamp: metric.Timestamp,
					},
				}
				continue
			}
			for m, _ := range metricMap {
				if IsEqual(m, metric) {
					metricMap[m] = append(metricMap[m], prompb.Sample{
						Value:     metric.Value,
						Timestamp: metric.Timestamp,
					})
					eq = true
					break
				}
			}
			if !eq {
				metricMap[metric] = []prompb.Sample{
					{
						Value:     metric.Value,
						Timestamp: metric.Timestamp,
					},
				}
			}

		}
	}

	for m, ss := range metricMap {
		labels = labels[:0]
		labels = transform(m)
		items = append(items, &prompb.TimeSeries{
			Labels:  labels,
			Samples: ss,
		})
	}
	return &prompb.WriteRequest{
		Timeseries: items,
	}, nil
}

func transform(metric *mod.MetricsMessage) []*prompb.Label {
	var labels []*prompb.Label
	if metric.MetricName != "" {
		labels = append(labels, &prompb.Label{
			Name:  consts.FieldMetricName,
			Value: metric.MetricName,
		})
		labels = append(labels, &prompb.Label{
			Name:  consts.FieldName,
			Value: metric.MetricName,
		})
	}
	if metric.ServiceName != "" {
		labels = append(labels, &prompb.Label{
			Name:  consts.FieldServiceName,
			Value: metric.ServiceName,
		})
	}
	if metric.ClusterName != "" {
		labels = append(labels, &prompb.Label{
			Name:  consts.FieldClusterName,
			Value: metric.ClusterName,
		})
	}
	if metric.ServiceHost != "" {
		labels = append(labels, &prompb.Label{
			Name:  consts.FieldServiceHost,
			Value: metric.ServiceHost,
		})
	}
	if metric.ServiceRegion != "" {
		labels = append(labels, &prompb.Label{
			Name:  consts.FieldServiceRegion,
			Value: metric.ServiceRegion,
		})
	}
	if metric.ServiceAz != "" {
		labels = append(labels, &prompb.Label{
			Name:  consts.FieldServiceAz,
			Value: metric.ServiceAz,
		})
	}
	for key, value := range metric.Labels {
		labels = append(labels, &prompb.Label{
			Name:  key,
			Value: value,
		})
	}
	return labels
}
