package read

import "github.com/prometheus/prometheus/prompb"

func formatFilter(query prompb.Query)(string,error) {
	//检查是否指定了服务名称，是否指定了指标名称

	//queries 长度应该是1
	//query.StartTimestampMs,query.EndTimestampMs,ka开始时间和结束时间，都是以毫秒为单位


	//query.Matchers
	//所有筛选条件，__name__为指标名称，serviceName为服务名称，


	//拼写SQL，
	//select timestamp * 1000 as t ,serviceName , metricName , func(value) from metrics.metircs_gf where timestamp >= startTimestmapMs/1000 and timestamp < endTimestamp /1000
	//and serviceName = 'someservice' and metricName = 'somemetric' group by t , serviceName,metricName order by t, metricName , serviceName;

	//可以分为selelct ,feilter， group , order区分模块分别实现。

}
