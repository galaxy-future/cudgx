package service

import (
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"github.com/galaxy-future/cudgx/internal/predict/query"
)

//QueryRedundancy 基于QPS查询系统冗余度
func QueryRedundancy(serviceName, clusterName, metricName string, benchmark float64, begin, end int64, trimmedSecond int64) (*RedundancySeries, error) {
	samples, err := query.AverageMetric(serviceName, clusterName, metricName, begin, end)
	if err != nil {
		return nil, err
	}
	clusters := samples2ClusterSeries(samples, trimmedSecond)
	for _, cluster := range clusters {
		for i := range cluster.Values {
			cluster.Values[i] = benchmark / cluster.Values[i]
		}
	}
	series := &RedundancySeries{
		ServiceName: serviceName,
		MetricName:  consts.QPSMetricsName,
	}

	for _, cluster := range clusters {
		series.Clusters = append(series.Clusters, cluster)
	}
	return series, nil
}

//QueryServiceTotalMetric 基于QPS查询系统冗余度
func QueryServiceTotalMetric(serviceName, clusterName, metricName string, begin, end int64, trimmedSecond int64) (*RedundancySeries, error) {
	samples, err := query.TotalMetric(serviceName, clusterName, metricName, begin, end)
	if err != nil {
		return nil, err
	}

	clusters := samples2ClusterSeries(samples, trimmedSecond)

	series := &RedundancySeries{
		ServiceName: serviceName,
		MetricName:  consts.QPSMetricsName,
	}

	for _, cluster := range clusters {
		series.Clusters = append(series.Clusters, cluster)
	}
	return series, nil
}

//QueryInstancesByMetric 基于QPS查询服务节点个数
func QueryInstancesByMetric(serviceName, clusterName, metricName string, begin, end int64, trimmedSecond int64) (*RedundancySeries, error) {
	samples, err := query.InstanceCountByMetric(serviceName, clusterName, metricName, begin, end)
	if err != nil {
		return nil, err
	}

	clusters := samples2ClusterSeries(samples, trimmedSecond)

	series := &RedundancySeries{
		ServiceName: serviceName,
		MetricName:  consts.QPSMetricsName,
	}

	for _, cluster := range clusters {
		series.Clusters = append(series.Clusters, cluster)
	}
	return series, nil
}
