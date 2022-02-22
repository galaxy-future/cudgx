package victoriametrics

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/common/types"
	"go.uber.org/zap"
)

type CommitFunc func(client *http.Client, messages []interface{}) error

type AsyncWriter struct {
	Client *http.Client `json:"client"`
	VmUrl  string       `json:"vm_url"`
	// MessagesCh 写入缓存
	FlushDuration types.Duration `json:"flush_duration"`
	RetryCount    int            `json:"retry_count"`
	Backoff       types.Duration `json:"backoff"`
	BatchSize     int            `json:"batch_size"`
	Concurrency   int            `json:"concurrency"`
	commit        CommitFunc
	messagesCh    <-chan interface{}
}

// NewWriter 新建VictoriaMetrics AsyncWriter
func NewWriter(config *Config, messagesCh <-chan interface{}, commit CommitFunc) (*AsyncWriter, error) {
	cli := createHTTPClient(config)
	return &AsyncWriter{
		Client:        cli,
		VmUrl:         config.Writer.VmUrl,
		FlushDuration: config.Writer.FlushDuration,
		RetryCount:    config.Writer.RetryCount,
		Backoff:       config.Writer.Backoff,
		BatchSize:     config.Writer.BatchSize,
		Concurrency:   config.Writer.Concurrency,
		commit:        commit,
		messagesCh:    messagesCh,
	}, nil
}

func createHTTPClient(cfg *Config) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     time.Duration(90) * time.Second,
		},
	}
	return client
}

type sender struct {
	writer            *AsyncWriter
	messages          []interface{}
	currentConnection int
}

func (writer *AsyncWriter) Init() {
	var wg sync.WaitGroup
	for i := 0; i < writer.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := sender{
				writer:   writer,
				messages: make([]interface{}, 0, writer.BatchSize),
			}
			s.start()
		}()
	}
	wg.Wait()
}

func (s *sender) start() {
	ticker := time.NewTicker(s.writer.FlushDuration.Duration)
	for {
		select {
		case <-ticker.C:
			s.flushWithRetry(s.messages)
			s.messages = s.messages[:0]
		case item, ok := <-s.writer.messagesCh:
			if ok {
				s.messages = append(s.messages, item)
				if len(s.messages) >= s.writer.BatchSize {
					s.flushWithRetry(s.messages)
					s.messages = s.messages[:0]
				}
			} else {
				s.flushWithRetry(s.messages)
				logger.GetLogger().Error("Failed to get the message from MessageCh")
				return
			}
		}
	}
}

func (s *sender) flushWithRetry(messages []interface{}) {
	err := s.flush(messages)
	if err == nil {
		return
	}
	logger.GetLogger().Error("Failed to write metrics", zap.Error(err))
	time.Sleep(s.writer.Backoff.Duration)
	for i := 0; i < s.writer.RetryCount; i++ {
		err := s.flush(messages)
		if err == nil {
			return
		}
		logger.GetLogger().Error("Failed to write metrics", zap.Error(err))
		time.Sleep(s.writer.Backoff.Duration)
	}
}

func (s *sender) flush(messages []interface{}) error {
	return s.writer.commit(s.writer.Client, messages)
}
