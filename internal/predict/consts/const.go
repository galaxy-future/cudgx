package consts

import "time"

const QPSMetricsName = "qps"
const LatencySectionFactorMetricsName = "qps_section_factor"
const DefaultPredictQueryCount = 60
const DefaultPredictMinCount = 50
const MinSkipFactor = 0.1
const DefaultRuleConcurrency = 10

const DefaultTrimmedSecond = 1
const TrimmedSecond = 5
const StepDuration = time.Second * 1

const (
	SchedulxExpandSuccess = "调用schedulx扩容接口成功："
	SchedulxShrinkSuccess = "调用schedulx缩容接口成功："
)

const (
	XClientUsername = "cudgx"
	XClientPassword = "Zpvo3nNPahZIXA1"
)

const (
	RuleStatusEnable  = "enable"
	RuleStatusDisable = "disable"
)

const (
	MetricNameRedundancy    = "redundancy"
	MetricNameLoad          = "load"
	MetricNameInstanceCount = "insanceCount"
)
