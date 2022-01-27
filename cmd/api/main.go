package main

import (
	"context"
	"flag"
	"net"

	"github.com/galaxy-future/cudgx/cmd/api/handler"
	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/predict"
	"github.com/galaxy-future/cudgx/internal/predict/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	configFile = flag.String("gf.cudgx.api.config", "conf/api.json", "api configure file")
	serverBind = flag.String("gf.cudgx.api.bind", "0.0.0.0:19003", "server bind address default(0.0.0.0:19003)")
)

func main() {
	flag.Parse()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger.GetLogger())

	theConfig, err := config.LoadConfig(*configFile)
	if err != nil {
		panic("load theConfig file falied : " + err.Error())
	}

	if err := predict.InitializeByConfig(theConfig); err != nil {
		panic(err)
	}

	go predict.StartRedundancyKeeper(context.Background())

	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())
	r.GET("/", func(context *gin.Context) {
		context.String(200, "success")
	})
	r.GET("/ping", func(context *gin.Context) {
		context.String(200, "cudgx/api-service is running")
	})
	redundancyGroup := r.Group("/api/v1/query/redundancy")
	{
		redundancyGroup.GET("/qps_average", handler.QueryRedundancyByQPS)
		redundancyGroup.GET("/instance_count", handler.QueryInstanceCountByQPSMetrics)
		redundancyGroup.GET("/qps_total", handler.QueryTotalQPS)
	}
	metricGroup := r.Group("/api/v1/query/metric")
	{
		metricGroup.GET("/redundancy/:metric_name", handler.QueryRedundancy)
		metricGroup.GET("/redundancy", handler.QueryRedundancyByServiceNameAndClusterName)
		metricGroup.GET("/instance_count/:metric_name", handler.QueryInstanceCountByMetrics)
		metricGroup.GET("/load/:metric_name", handler.QueryTotalMetric)
	}

	predictApiV1 := r.Group("/api/v1/cudgx/predict")
	rulePath := predictApiV1.Group("/rule")
	{
		rulePath.GET("/:id", handler.GetPredictRule)
		rulePath.GET("/info", handler.GetPredictRuleInfo)
		rulePath.POST("/create", handler.CreatePredictRule)
		rulePath.POST("/update", handler.UpdatePredictRule)
		rulePath.POST("/batch/delete", handler.BatchDeletePredictRule)
		rulePath.GET("/list", handler.ListPredictRules)
		rulePath.POST("/:id/enable", handler.EnablePredictRule)
		rulePath.POST("/:id/disable", handler.DisablePredictRule)
	}

	l, err := net.Listen("tcp", *serverBind)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic(err)
	}
	err = r.RunListener(l)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic("server start failed")
	}
}
