package victoriametrics

import (
	"fmt"

	json "github.com/json-iterator/go"
	"github.com/prometheus/common/model"
)

type Response struct {
	Data model.Matrix
}

func (res *Response) UnmarshalJSON(b []byte) error {
	v := struct {
		ResultType model.ValueType `json:"resultType"`
		Result     json.RawMessage `json:"result"`
	}{}

	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v.ResultType != model.ValMatrix {
		return fmt.Errorf("unexpected value type %q", v.ResultType)
	}
	var matrix model.Matrix
	err = json.Unmarshal(v.Result, &matrix)
	res.Data = matrix

	return err
}

type RawMsg struct {
	Data json.RawMessage `json:"data"`
}
