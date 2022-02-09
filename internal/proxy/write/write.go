package write

import (
	"github.com/prometheus/prometheus/prompb"
)
type RemoteWriter interface {
	WriteMessage(request prompb.WriteRequest) error
}

//CudgXGatewayWriter 接收RemoteWrite消息，写入到gateway完成分发存储逻辑
type CudgXGatewayWriter struct {

}
