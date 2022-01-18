package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/galaxy-future/cudgx/common/logger"
	"go.uber.org/zap"
)

var piServer = flag.String("gf.cudgx.sample.benchmark.sever-address", "http://localhost:8090/pi", "server bind address default(http://localhost:8090)")
var concurrency = flag.Int("gf.cudgx.sample.benchmark.concurrency", 5, "concurrency count for benchmark")
var base = flag.Int("gf.cudgx.sample.benchmark.base", 150, "base request count")
var revolution = flag.Int("gf.cudgx.sample.benchmark.revolution", 120, "revolution of time change / seconds")

var (
	httpclient *http.Client
)

type benchmarkAlgo struct {
	requestSigs chan struct{}
}

func (benchmark *benchmarkAlgo) generate() {
	count := *base + int(math.Sin(float64(time.Now().Unix())/float64(*revolution)*math.Pi)*float64(*base/2))
	for i := 0; i < count; i++ {
		benchmark.requestSigs <- struct{}{}
	}
}

func main() {
	flag.Parse()

	httpclient = &http.Client{
		Timeout: 1000 * time.Millisecond,
	}
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	requestSigs := make(chan struct{}, 1000)

	benchmark := benchmarkAlgo{requestSigs: requestSigs}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTimer(1 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				benchmark.generate()
				ticker.Reset(time.Second)
			}
		}
	}()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		Looper:
			for {
				select {
				case <-ctx.Done():
					break Looper
				case <-requestSigs:
					request()
				}
			}
		}()
	}

	//wait
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	//wait a signal
	<-c

	cancel()
	wg.Wait()

}

func request() {

	path := "/low"
	switch rand.Intn(3) {
	case 0:
		path = "/low"
	case 1:
		path = "/middle"
	case 2:
		path = "/high"
	}
	resp, err := httpclient.Get(fmt.Sprintf("%s%s", *piServer, path))
	if err != nil {
		logger.GetLogger().Error("request failed", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Error("request failed", zap.Error(err))
		return
	}

	if resp.StatusCode/100 != 2 && resp.StatusCode/100 != 3 {
		logger.GetLogger().Error("request failed", zap.Error(fmt.Errorf(string(respData))))
		return
	}

}
