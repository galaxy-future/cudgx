package victoriametrics

import "github.com/galaxy-future/cudgx/common/types"

type Config struct {
	Writer Write  `json:"writer"`
	Reader Reader `json:"reader"`
}

type Write struct {
	// Remote write url
	VmUrl         string         `json:"vm_url"`
	FlushDuration types.Duration `json:"flush_duration"`
	RetryCount    int            `json:"retry_count"`
	Backoff       types.Duration `json:"backoff"`
	BatchSize     int            `json:"batch_size"`
	Concurrency   int            `json:"concurrency"`
}

type Read struct {
	// Prometheus APIs base url
	VmUrl string `json:"vm_url"`
}
