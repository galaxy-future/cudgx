package proxy

import (
	"github.com/galaxy-future/cudgx/internal/proxy/read"
	"github.com/galaxy-future/cudgx/internal/proxy/write"
)

type CudgxProxy struct {
	writer * write.CudgXGatewayWriter
	reader * read.CudgXRemoteReader
}

type Proxyconfig struct {
}

var proxy * CudgxProxy


func GetProxy() *CudgxProxy{
	return proxy
}


func  initCudgxProxy(*Proxyconfig) *CudgxProxy {

}





