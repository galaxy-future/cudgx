package query

import (
	"fmt"

	"github.com/galaxy-future/cudgx/internal/clients"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
)

//InstanceCountByQPS 查询服务节点数量
//Deprecated: Use InstanceCountByMetric
func InstanceCountByQPS(serviceName, clusterName string, begin, end int64) (samples []ClusterSample, err error) {
	client := clients.ClickhouseRdCli
	sqlContent := fmt.Sprintf(`select timestamp ,clusterName , count( distinct(serviceHost)) 
			from %s.%s 
			where timestamp >= %d and  timestamp < %d and metricName = '%s'  and serviceName = '%s' and clusterName = '%s'
			group by timestamp , serviceName, clusterName, metricName
			order by timestamp `, client.Database, client.Table, begin, end, consts.QPSMetricsName, serviceName, clusterName)

	return queryClusterSamples(sqlContent)
}
