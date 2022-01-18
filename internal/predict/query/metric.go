package query

import (
	"fmt"
	"github.com/galaxy-future/cudgx/internal/clients"
)

//AverageMetric 查询服务/集群的平均Metric值
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
func InstanceCountByMetric(serviceName, clusterName, metricName string, begin, end int64) (samples []ClusterSample, err error) {
	client := clients.ClickhouseRdCli
	sqlContent := fmt.Sprintf(`select timestamp ,clusterName , count( distinct(serviceHost)) 
			from %s.%s 
			where timestamp >= %d and  timestamp < %d and metricName = '%s'  and serviceName = '%s' and clusterName = '%s'
			group by timestamp , serviceName, clusterName, metricName
			order by timestamp `, client.Database, client.Table, begin, end, metricName, serviceName, clusterName)

	return queryClusterSamples(sqlContent)
}
