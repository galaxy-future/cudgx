package handler

import (
	"fmt"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"github.com/galaxy-future/cudgx/internal/predict/service"
	"github.com/galaxy-future/cudgx/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

//QueryRedundancyMetricName 获取服务冗余度衡量指标，当前支持只支持单指标
func QueryRedundancyMetricName(c *gin.Context) {
	serviceName := c.Query("service_name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("服务名称不能为空"))
		return
	}
	clusterName := c.Query("cluster_name")
	if clusterName == "" {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("集群名称不能为空"))
		return
	}
	rule, err := service.GetPredictRuleByServiceNameAndClusterName(serviceName, clusterName)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(fmt.Sprintf("获取规则时出错, err: %s", err)))
		return
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(rule.MetricName))
	return
}

// QueryRedundancy 基于QPS指标数据输出冗余度
func QueryRedundancy(c *gin.Context) {
	serviceName, clusterName, metricName, begin, end, pass := validateMetricQuery(c)
	if !pass {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(fmt.Sprintf("参数错误")))
		return
	}
	rule, err := service.GetPredictRuleByServiceNameAndClusterName(serviceName, clusterName)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(fmt.Sprintf("获取规则时出错, err: %s", err)))
		return
	}
	if rule.MetricName != metricName {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(
			fmt.Sprintf("服务和所在集群冗余度评估方式指标名称不同, metricName: %s", rule.MetricName)))
		return
	}
	benchmark := rule.BenchmarkQps
	if benchmark <= 0 {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("benchmark不能为0"))
		return
	}

	redundancySeries, err := service.QueryRedundancy(serviceName, clusterName, metricName, float64(benchmark), begin, end, consts.TrimmedSecond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MkFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MkSuccessResponse(redundancySeries))
}

// QueryTotalMetric 查询QPS
func QueryTotalMetric(c *gin.Context) {
	serviceName, clusterName, metricName, begin, end, pass := validateMetricQuery(c)
	if !pass {
		return
	}
	redundancySeries, err := service.QueryServiceTotalMetric(serviceName, clusterName, metricName, begin, end, consts.TrimmedSecond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MkFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MkSuccessResponse(redundancySeries))
}

// QueryInstanceCountByMetrics 查询机器数
func QueryInstanceCountByMetrics(c *gin.Context) {
	serviceName, clusterName, metricName, begin, end, pass := validateMetricQuery(c)
	if !pass {
		return
	}
	redundancySeries, err := service.QueryInstancesByMetric(serviceName, clusterName, metricName, begin, end, consts.TrimmedSecond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MkFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MkSuccessResponse(redundancySeries))
}

// validateMetricQuery 校验参数合法性
func validateMetricQuery(c *gin.Context) (serviceName, clusterName, metricName string, begin, end int64, pass bool) {
	metricName = c.Param("metric_name")
	if metricName == "" {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("metric不能为空"))
		return
	}
	serviceName = c.Query("service_name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("服务名称不能为空"))
		return
	}
	clusterName = c.Query("cluster_name")
	if clusterName == "" {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("集群名称不能为空"))
		return
	}
	begin = cast.ToInt64(c.Query("begin"))
	if begin == 0 {
		begin = time.Now().Add(-5 * time.Minute).Unix()
	}

	end = cast.ToInt64(c.Query("end"))
	if end == 0 {
		end = time.Now().Unix()
	}

	if end <= begin {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("开始时间不能大于结束时间"))
		return
	}
	pass = true
	return
}
