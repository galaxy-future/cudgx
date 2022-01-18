package query

import (
	"fmt"

	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/clients"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"go.uber.org/zap"
)

type ClusterSample struct {
	Timestamp   int64
	Value       float64
	ClusterName string
}

//AverageQPS 查询服务/集群的平均QPS
//Deprecated: Use AverageMetric
func AverageQPS(serviceName, clusterName string, begin, end int64) (samples []ClusterSample, err error) {
	client := clients.ClickhouseRdCli
	sqlContent := fmt.Sprintf(`select timestamp ,clusterName , sum(value)/ count( distinct(serviceHost) ) 
			from %s.%s 
			where timestamp >= %d and  timestamp < %d and metricName = '%s'  and serviceName = '%s' and clusterName = '%s'
			group by timestamp ,serviceName, clusterName, metricName
			order by timestamp `, client.Database, client.Table, begin, end, consts.QPSMetricsName, serviceName, clusterName)

	return queryClusterSamples(sqlContent)

}

//TotalQPS 查询集群 QPS
//Deprecated: Use TotalMetric
func TotalQPS(serviceName, clusterName string, begin, end int64) (samples []ClusterSample, err error) {

	sqlContent := fmt.Sprintf(`select timestamp ,clusterName , sum(value)
			from %s.%s 
			where timestamp >= %d and  timestamp < %d and metricName = '%s'  and serviceName = '%s' and clusterName = '%s' 
			group by timestamp , serviceName, clusterName, metricName
			order by timestamp `, clients.ClickhouseRdCli.Database, clients.ClickhouseRdCli.Table, begin, end, consts.QPSMetricsName, serviceName, clusterName)

	return queryClusterSamples(sqlContent)
}

func queryClusterSamples(sql string) (samples []ClusterSample, err error) {
	rows, err := clients.ClickhouseRdCli.Client.Query(sql)
	if err != nil {
		logger.GetLogger().Error("failed to query qps metrics", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var timestamp int64
		var value float64
		var clusterName string
		err := rows.Scan(&timestamp, &clusterName, &value)
		if err != nil {
			logger.GetLogger().Error("failed to query qps metrics", zap.Error(err))
			return nil, err
		}
		samples = append(samples, ClusterSample{
			Timestamp:   timestamp,
			Value:       value,
			ClusterName: clusterName,
		})
	}
	return
}
