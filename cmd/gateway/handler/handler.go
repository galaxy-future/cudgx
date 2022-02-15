package handler

import (
	"github.com/galaxy-future/cudgx/internal/clients"
	"io/ioutil"

	"github.com/galaxy-future/cudgx/common/mod"
	"github.com/galaxy-future/cudgx/internal/gateway"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

// HandlerMonitoringMessageBatch monitoring指标数据处理
func HandlerMonitoringMessageBatch(c *gin.Context) {
	serviceName := c.Param("service")
	metricName := c.Param("metric")

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"message": "read request body failed"})
		return
	}

	writer, err := gateway.GetGateway().GetMonitoringWriter(serviceName, metricName)
	if err != nil {
		c.JSON(500, gin.H{"message": "get kafka client failed "})
		return
	}
	writer.SendMessage(serviceName, metricName, data)

	c.JSON(200, gin.H{"message": "success"})
}

// HandlerPing gateway探活接口
func HandlerPing(c *gin.Context) {
	c.JSON(200, mod.GatewayPingResult{
		Status: mod.GatewayStatusSuccess,
		Module: mod.GatewayModuleName,
	})
}

// HandlerStreamingMessageBatch streaming指标数据处理
func HandlerStreamingMessageBatch(c *gin.Context) {
	serviceName := c.Param("service")
	metricName := c.Param("metric")

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"message": "read request body failed"})
		return
	}
	var streamingBatch mod.StreamingBatch
	err = proto.Unmarshal(data, &streamingBatch)
	if err != nil {
		c.JSON(400, gin.H{"message": "unmarshal messages failed", "error": err.Error()})
		return
	}
	wrapStreamingBatch, err := gateway.GetGateway().WrapStreamingMessage(&streamingBatch)
	if err != nil {
		c.JSON(500, gin.H{"message": "wrap messages failed", "error": err.Error()})
		return
	}
	data, err = proto.Marshal(wrapStreamingBatch)
	if err != nil {
		c.JSON(500, gin.H{"message": "marshal messages failed", "error": err.Error()})
		return
	}
	writer, err := gateway.GetGateway().GetStreamingWriter(serviceName, metricName)
	if err != nil {
		c.JSON(500, gin.H{"message": "get kafka client failed "})
		return
	}
	writer.SendMessage(serviceName, metricName, data)

	c.JSON(200, gin.H{"message": "success"})
}

func RemoteWrite(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"message": "read body failed", "error": err.Error()})
		return
	}
	reqBuf, err := snappy.Decode(nil, body)
	if err != nil {
		c.JSON(400, gin.H{"message": "read body failed", "error": err.Error()})
		return
	}
	req := &prompb.WriteRequest{}
	err = req.Unmarshal(reqBuf)
	if err != nil {
		c.JSON(400, gin.H{"message": "read unmarshal request failed", "error": err.Error()})
		return
	}
	ip := c.ClientIP()
	service, err := clients.GetServiceByIp(ip)
	if err != nil {
		c.JSON(400, gin.H{"message": "wrong ip", "error": err.Error()})
		return
	}
	msgs := make([]*mod.StreamingMessage, 0, len(req.Timeseries))
	writer, err := gateway.GetGateway().GetStreamingWriter(service.ServiceName, "")
	if err != nil {
		c.JSON(500, gin.H{"message": "get kafka client failed ", "error": err.Error()})
		return
	}
	data := &mod.StreamingBatch{
		ServiceName: service.ServiceName,
		MetricName:  "",
		Messages:    msgs,
	}
	bData, err := proto.Marshal(data)
	if err != nil {
		c.JSON(400, gin.H{"message": "StreamingBatch marshal failed", "error": err.Error()})
		return
	}
	writer.SendMessage(service.ServiceName, "", bData)
	c.JSON(200, gin.H{"message": "success"})
}
