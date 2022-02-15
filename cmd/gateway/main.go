package main

import (
	"flag"
	"net"

	"github.com/galaxy-future/cudgx/cmd/gateway/handler"
	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/gateway"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	configFile = flag.String("gf.cudgx.gateway.config", "conf/gateway.json", "gateway configure file")
	serverBind = flag.String("gf.cudgx.gateway.bind", "0.0.0.0:8080", "server bind address default(0.0.0.0:8080)")
)

func main() {
	flag.Parse()
	defer logger.GetLogger().Sync()

	err := gateway.Init(*configFile)
	if err != nil {
		panic("load config file failed : " + err.Error())
	}

	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "success")
	})
	r.GET("/ping", handler.HandlerPing)
	r.POST("/v1/monitoring/:service/:metric", handler.HandlerMonitoringMessageBatch)
	r.POST("/v1/streaming/:service/:metric", handler.HandlerStreamingMessageBatch)
	r.POST("/v1/prom/remote/write", handler.RemoteWrite)

	l, err := net.Listen("tcp", *serverBind)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic("server listen failed ")
	}

	err = r.RunListener(l)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic("server start failed ")
	}
}
