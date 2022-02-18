package query

import (
	"fmt"

	"github.com/galaxy-future/cudgx/common/victoriametrics"
	"github.com/galaxy-future/cudgx/internal/clients"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
)

var Reader *victoriametrics.Reader

//AverageMetric 查询服务/集群的平均Metric值
//Deprecated: Use QueryAverageMetricByVM
func AverageMetric(serviceName, clusterName, metricName string, begin, end int64) (samples []ClusterSample, err error) {
	client := clients.ClickhouseRdCli
	sqlContent := fmt.Sprintf(`select timestamp ,clusterName , sum(value)/ count( distinct(serviceHost) ) 
			from %s.%s 
			where timestamp >= %d and  timestamp < %d and metricName = '%s'  and serviceName = '%s' and clusterName = '%s'
			group by timestamp ,serviceName, clusterName, metricName
			order by timestamp `, client.Database, client.Table, begin, end, metricName, serviceName, clusterName)

	return queryClusterSamples(sqlContent)

}

//TotalMetric 查询集群Metric
//Deprecated: Use QueryTotalMetricByVM
func TotalMetric(serviceName, clusterName, metricName string, begin, end int64) (samples []ClusterSample, err error) {

	sqlContent := fmt.Sprintf(`select timestamp ,clusterName , sum(value)
			from %s.%s 
			where timestamp >= %d and  timestamp < %d and metricName = '%s'  and serviceName = '%s' and clusterName = '%s' 
			group by timestamp , serviceName, clusterName, metricName
			order by timestamp `, clients.ClickhouseRdCli.Database, clients.ClickhouseRdCli.Table,
		begin, end, metricName, serviceName, clusterName)

	return queryClusterSamples(sqlContent)
}

//InstanceCountByMetric 查询服务节点数量
//Deprecated: Use QueryInstanceCountByVM
func InstanceCountByMetric(serviceName, clusterName, metricName string, begin, end int64) (samples []ClusterSample, err error) {
	client := clients.ClickhouseRdCli
	sqlContent := fmt.Sprintf(`select timestamp ,clusterName , count( distinct(serviceHost)) 
			from %s.%s 
			where timestamp >= %d and  timestamp < %d and metricName = '%s'  and serviceName = '%s' and clusterName = '%s'
			group by timestamp , serviceName, clusterName, metricName
			order by timestamp `, client.Database, client.Table, begin, end, metricName, serviceName, clusterName)

	return queryClusterSamples(sqlContent)
}

//AverageMetricByVM 查询服务/集群的平均Metric值
func AverageMetricByVM(serviceName, clusterName, metricName string, begin, end int64) (samples []ClusterSample, err error) {
	promeQL := fmt.Sprintf("sum(%s{serviceName='%s',clusterName='%s'})/count(%s{serviceName='%s',clusterName='%s'}) by(metricName,serviceName,clusterName)", metricName, serviceName, clusterName, metricName, serviceName, clusterName)
	res, err := Reader.QueryRange(promeQL, begin, end, consts.StepDuration)
	if err != nil {
		return nil, err
	}
	return convertSamples(res), nil
}

//TotalMetricByVM 查询集群Metric
func TotalMetricByVM(serviceName, clusterName, metricName string, begin, end int64) (samples []ClusterSample, err error) {
	promeQL := fmt.Sprintf("sum(%s{serviceName='%s',clusterName='%s'}) by(metricName,serviceName,clusterName)", metricName, serviceName, clusterName)
	res, err := Reader.QueryRange(promeQL, begin, end, consts.StepDuration)
	if err != nil {
		return nil, err
	}
	return convertSamples(res), nil
}

//InstanceCountByVM 查询服务节点数量
func InstanceCountByVM(serviceName, clusterName, metricName string, begin, end int64) (samples []ClusterSample, err error) {
	promeQL := fmt.Sprintf("count(%s{serviceName='%s',clusterName='%s'}) by(metricName,serviceName,clusterName)", metricName, serviceName, clusterName)
	res, err := Reader.QueryRange(promeQL, begin, end, consts.StepDuration)
	if err != nil {
		return nil, err
	}
	return convertSamples(res), nil
}

// convertSamples
func convertSamples(res *victoriametrics.Response) []ClusterSample {
	var samples []ClusterSample
	for _, sampleStream := range res.Data {
		clusterName, ok := sampleStream.Metric["clusterName"]
		if !ok {
			continue
		}
		for _, value := range sampleStream.Values {
			samples = append(samples, ClusterSample{
				Timestamp:   int64(value.Timestamp) / 1000,
				Value:       float64(value.Value),
				ClusterName: string(clusterName),
			})
		}
	}
	return samples
}
