package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	metricGo "github.com/galaxy-future/metrics-go"
	"github.com/galaxy-future/metrics-go/aggregate"
	"github.com/galaxy-future/metrics-go/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewLatencySectorFactorBuilder() *LatencySectorFactorBuilder {
	return &LatencySectorFactorBuilder{}
}

type LatencySectorFactorBuilder struct {
}

func (builder *LatencySectorFactorBuilder) Build() aggregate.Function {
	return &LatencySectorFactor{}
}

type LatencySectorFactor struct {
	section1QPS int
	section2QPS int
	section3QPS int
	section4QPS int
	section5QPS int
}

func (latency *LatencySectorFactor) Accumulate(value float64) {
	if value < 10 {
		latency.section1QPS++
	} else if value < 50 {
		latency.section2QPS++
	} else if value < 100 {
		latency.section3QPS++
	} else if value < 500 {
		latency.section4QPS++
	} else {
		latency.section5QPS++
	}

}

func (latency *LatencySectorFactor) Aggregate() float64 {
	return float64(latency.section1QPS)*0.1 + float64(latency.section2QPS)*0.5 +
		float64(latency.section3QPS)*2.0 + float64(latency.section4QPS)*4.0 +
		float64(latency.section5QPS)*8
}

func (latency *LatencySectorFactor) Values() []float64 {
	return []float64{latency.Aggregate()}
}

var (
	serverBind     = flag.String("gf.cudgx.sample.pi.bind", "0.0.0.0:8090", "server bind address default(0.0.0.0:8090)")
	goRoutineCount = flag.Int("gf.cudgx.sample.pi.count", 5000, "go routine count to calc pi")
)

var (
	latencyMin           types.Metrics
	latencyMax           types.Metrics
	latencySectionFactor types.Metrics
	latency              types.Metrics
	qps                  types.Metrics
)

func main() {
	flag.Parse()
	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.GET("/pi", HandlerCalcPiLow)
	r.GET("/pi/low", HandlerCalcPiLow)
	r.GET("/pi/middle", HandlerCalcPiMiddle)
	r.GET("/pi/high", HandlerCalcPiHigh)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "success")
	})

	l, err := net.Listen("tcp", *serverBind)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic("server listen failed ")
	}

	initMetrics()

	err = r.RunListener(l)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic("server start failed ")
	}
}

func HandlerCalcPiLow(c *gin.Context) {
	begin := time.Now()
	c.String(200, fmt.Sprintf("%v", pi(*goRoutineCount/5)))

	cost := time.Now().Sub(begin).Milliseconds()
	latencySectionFactor.With().Value(float64(time.Now().Sub(begin).Microseconds()) / 1000.0)
	latencyMin.With().Value(float64(cost))
	latencyMax.With().Value(float64(cost))
	latency.With().Value(float64(cost))
	qps.With().Value(1)
}
func HandlerCalcPiMiddle(c *gin.Context) {
	begin := time.Now()
	c.String(200, fmt.Sprintf("%v", pi(*goRoutineCount/2)))

	cost := time.Now().Sub(begin).Milliseconds()
	latencySectionFactor.With().Value(float64(time.Now().Sub(begin).Microseconds()) / 1000.0)
	latencyMin.With().Value(float64(cost))
	latencyMax.With().Value(float64(cost))
	latency.With().Value(float64(cost))
	qps.With().Value(1)
}
func HandlerCalcPiHigh(c *gin.Context) {
	begin := time.Now()
	c.String(200, fmt.Sprintf("%v", pi(*goRoutineCount)))

	cost := time.Now().Sub(begin).Milliseconds()
	latencySectionFactor.With().Value(float64(time.Now().Sub(begin).Microseconds()) / 1000.0)
	latencyMin.With().Value(float64(cost))
	latencyMax.With().Value(float64(cost))
	latency.With().Value(float64(cost))
	qps.With().Value(1)
}

// pi launches n goroutines to compute an
// approximation of pi.
func pi(n int) float64 {
	f := 0.0
	for k := 0; k < n; k++ {
		f += 4 * math.Pow(-1, float64(k)) / (2*float64(k) + 1)
	}
	return f
}

func initMetrics() {
	latencyMin = metricGo.NewMonitoringMetric("latencyMin", []string{}, aggregate.NewMinBuilder())
	latencyMax = metricGo.NewMonitoringMetric("latencyMax", []string{}, aggregate.NewMaxBuilder())
	latencySectionFactor = metricGo.NewMonitoringMetric(consts.LatencySectionFactorMetricsName,
		[]string{}, NewLatencySectorFactorBuilder())

	latency = metricGo.NewStreamingMetric("latency", []string{})
	qps = metricGo.NewMonitoringMetric("qps", []string{}, aggregate.NewCountBuilder())
}
