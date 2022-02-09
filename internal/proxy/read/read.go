package read

import (
	"github.com/prometheus/prometheus/prompb"
)
type RemoteRead interface {
	Query(request prompb.ReadRequest)(prompb.ReadResponse,error)
}

type CudgXRemoteReader struct {

}

func (cudgx * CudgXRemoteReader) Query(request prompb.ReadRequest)(prompb.ReadResponse,error) {

}

