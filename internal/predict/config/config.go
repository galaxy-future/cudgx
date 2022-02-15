package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/galaxy-future/cudgx/common/clickhouse"
	"github.com/galaxy-future/cudgx/common/types"
	"github.com/galaxy-future/cudgx/common/victoriametrics"
)

//Config 预测器配置
type Config struct {
	//Predict 预测参数配置
	Predict *Param `json:"param"`
	//Clickhouse 连接配置
	Clickhouse *clickhouse.Config `json:"clickhouse"`
	//Mysql 连接配置
	Database *Database `json:"database"`
	//xclient 连接配置
	Xclient *Xclient `json:"xclient"`
	//VictoriaMetrics 连接配置
	VictoriaMetrics *victoriametrics.Config `json:"victoria_metrics"`
}

//Xclient bridgx/schedulx连接配置
type Xclient struct {
	BridgxServerAddress   string `json:"bridgx_server_address"`
	SchedulxServerAddress string `json:"schedulx_server_address"`
}

//Param 是Predict过程中使用到的多个可调参数
type Param struct {
	//SamplesQueryCount 判定过程中，需要查询的Sample数量
	RunDuration types.Duration `json:"run_duration"`
	//RuleConcurrency 并行运行规则数量
	RuleConcurrency int `json:"rule_concurrency"`
	//MinimalSampleCount 参与判断中最少的指标点数，
	MinimalSampleCount int `json:"minimal_sample_count"`
	//LookbackDuration 回查多久
	LookbackDuration types.Duration `json:"lookback_duration"`
	//MetricSendDuration z指标传输所需时间，在这段时间内的指标是不准确的
	MetricSendDuration types.Duration `json:"metric_send_duration"`
}

//LoadConfig 从文件中加载配置
func LoadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("can not open configure file : %v ", configFile)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read configure file failed, err : %v ", err)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

type Database struct {
	Dsn string `json:"dsn"`
}
