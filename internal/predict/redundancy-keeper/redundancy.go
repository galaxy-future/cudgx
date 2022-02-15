package redundancy_keeper

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/clients"
	"github.com/galaxy-future/cudgx/internal/predict/config"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"github.com/galaxy-future/cudgx/internal/predict/model"
	"github.com/galaxy-future/cudgx/internal/predict/service"
	"go.uber.org/zap"
)

var (
	redundancyKeeper *ScheduleXRedundancyKeeper
)

//ScheduleXRedundancyKeeper 负责保持服务的冗余度
type ScheduleXRedundancyKeeper struct {
	ScheduleDuration time.Duration
	concurrencyLock  chan struct{}
	//MinimalSampleCount 参与判断中最少的指标点数，
	MinimalSampleCount int `json:"minimal_sample_count"`
	//LookbackDuration 回查多久
	LookbackDuration time.Duration `json:"lookback_duration"`
	//MetricSendDuration 指标传输所需时间，在这段时间内的指标是不准确的
	MetricSendDuration time.Duration `json:"metric_send_duration"`
}

func InitRedundancyKeeper(param *config.Param) {
	redundancyKeeper = &ScheduleXRedundancyKeeper{
		ScheduleDuration:   param.RunDuration.Duration,
		concurrencyLock:    make(chan struct{}, param.RuleConcurrency),
		MinimalSampleCount: param.MinimalSampleCount,
		LookbackDuration:   param.LookbackDuration.Duration,
		MetricSendDuration: param.MetricSendDuration.Duration,
	}
}

func Start(ctx context.Context) {
	ticker := time.NewTicker(redundancyKeeper.ScheduleDuration)
	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			err := redundancyKeeper.schedule()
			if err != nil {
				logger.GetLogger().Error("failed schedule rules", zap.Error(err))
			}
		}
	}
}

func (keeper *ScheduleXRedundancyKeeper) schedule() error {
	rules, err := model.ListAllPredictRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Status != consts.RuleStatusEnable {
			continue
		}
		keeper.concurrencyLock <- struct{}{}
		go func(theRule *model.PredictRule) {
			defer func() {
				<-keeper.concurrencyLock
			}()
			err := scheduleRule(theRule)
			if err != nil {
				logger.GetLogger().Error("failed to schedule service", zap.String("service", theRule.ServiceName), zap.String("cluster", theRule.ClusterName), zap.Error(err))
			}
		}(rule)
	}
	return nil
}

func scheduleRule(rule *model.PredictRule) error {
	const lookbackDuration = time.Minute
	const metricsSendDuration = 5 * time.Second
	const minSampleCount = lookbackDuration - 30*time.Second
	serviceName := rule.ServiceName
	clusterName := rule.ClusterName
	metricName := rule.MetricName
	benchmark := rule.BenchmarkQps

	series, err := service.QueryRedundancy(serviceName, clusterName, metricName, float64(benchmark), time.Now().Add(-1*lookbackDuration).Unix(), time.Now().Add(-1*metricsSendDuration).Unix(), consts.DefaultTrimmedSecond)
	if err != nil {
		return err
	}

	canSchedule, err := clients.CanServiceSchedule(serviceName, clusterName)
	if err != nil {
		return fmt.Errorf("query service schedule failed , %w", err)
	}
	if !canSchedule {
		return nil
	}

	currentCount, err := clients.GetServiceInstanceCount(serviceName, clusterName)
	if err != nil {
		return fmt.Errorf("query service instance count failed , %w", err)
	}

	for _, cluster := range series.Clusters {
		if cluster.ClusterName != clusterName {
			continue
		}
		// 没有足够的采集点
		if len(cluster.Values) < int(minSampleCount.Seconds()) {
			continue
		}
		sort.Float64s(cluster.Values)

		// 取中间数
		redundancy := cluster.Values[len(cluster.Values)/2]

		//不需要调度
		if int(redundancy*100) < rule.MaxRedundancy && int(redundancy*100) > rule.MinRedundancy {
			continue
		}

		//取冗余度的中间数
		midRedundancy := float64((rule.MaxRedundancy+rule.MinRedundancy)/2) / 100.0

		expectCount := int(midRedundancy / redundancy * float64(currentCount))

		diff := expectCount - currentCount

		countToChange := int(math.Ceil(float64(diff*rule.ExecuteRatio) / 100.0))

		if countToChange == 0 {
			continue
		}
		if countToChange > 0 {
			if currentCount+countToChange > rule.MaxInstanceCount {
				countToChange = rule.MaxInstanceCount - currentCount
			}
			if countToChange > 30 {
				countToChange = 30
			}
			err := clients.ExpandService(serviceName, clusterName, countToChange)
			if err != nil {
				return fmt.Errorf("expand service failed , %w", err)
			}
		} else {
			countToChange = int(math.Abs(float64(countToChange)))
			if currentCount-countToChange < rule.MinInstanceCount {
				countToChange = currentCount - rule.MinInstanceCount
			}
			if countToChange > 30 {
				countToChange = 30
			}
			err := clients.ShrinkService(serviceName, clusterName, countToChange)
			if err != nil {
				return fmt.Errorf("shrink service failed , %w", err)
			}
		}
	}
	return nil
}
