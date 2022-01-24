package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"github.com/galaxy-future/cudgx/internal/predict/service"
	"github.com/galaxy-future/cudgx/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// QueryRedundancyByQPS 基于QPS指标数据输出冗余度
func QueryRedundancyByQPS(c *gin.Context) {
	serviceName, clusterName, begin, end, pass := validateQPSQuery(c)
	if !pass {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(fmt.Sprintf("参数错误")))
		return
	}
	if end > time.Now().Unix() {
		end = time.Now().Unix()
	}
	begin = begin / consts.TrimmedSecond * consts.TrimmedSecond
	end = end / consts.TrimmedSecond * consts.TrimmedSecond
	rule, err := service.GetPredictRuleByServiceNameAndClusterName(serviceName, clusterName)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(fmt.Sprintf("获取规则时出错, err: %s", err)))
		return
	}
	benchmark := rule.BenchmarkQps
	if benchmark <= 0 {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("benchmark不能为0"))
		return
	}

	redundancySeries, err := service.QueryRedundancyByQPS(serviceName, clusterName, float64(benchmark), begin, end, consts.TrimmedSecond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MkFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MkSuccessResponse(redundancySeries))
}

// QueryTotalQPS 查询QPS
func QueryTotalQPS(c *gin.Context) {
	serviceName, clusterName, begin, end, pass := validateQPSQuery(c)
	if !pass {
		return
	}
	if end > time.Now().Unix() {
		end = time.Now().Unix()
	}
	begin = begin / consts.TrimmedSecond * consts.TrimmedSecond
	end = end / consts.TrimmedSecond * consts.TrimmedSecond
	redundancySeries, err := service.QueryServiceTotalQPS(serviceName, clusterName, begin, end, consts.TrimmedSecond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MkFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MkSuccessResponse(redundancySeries))
}

// QueryInstanceCountByQPSMetrics 查询机器数
func QueryInstanceCountByQPSMetrics(c *gin.Context) {
	serviceName, clusterName, begin, end, pass := validateQPSQuery(c)
	if !pass {
		return
	}
	if end > time.Now().Unix() {
		end = time.Now().Unix()
	}
	begin = begin / consts.TrimmedSecond * consts.TrimmedSecond
	end = end / consts.TrimmedSecond * consts.TrimmedSecond
	redundancySeries, err := service.QueryInstancesByQPS(serviceName, clusterName, begin, end, consts.TrimmedSecond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MkFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MkSuccessResponse(redundancySeries))
}

// validateQPSQuery 校验参数合法性
func validateQPSQuery(c *gin.Context) (serviceName, clusterName string, begin, end int64, pass bool) {
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
