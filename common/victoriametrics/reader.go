package victoriametrics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	json "github.com/json-iterator/go"
)

type Reader struct {
	Client *http.Client `json:"client"`
	VmUrl  string       `json:"vm_url"`
}

// NewReader 新建VictoriaMetrics Reader
func NewReader(config *Config) *Reader {
	cli := createHTTPClient(config)
	return &Reader{
		Client: cli,
		VmUrl:  config.Reader.VmUrl,
	}
}

// QueryRange Prometheus query_range API
func (r Reader) QueryRange(query string, start, end int64, step time.Duration) (*Response, error) {
	u := url.Values{}
	u.Set("query", query)
	u.Set("start", strconv.FormatInt(start, 10))
	u.Set("end", strconv.FormatInt(end, 10))
	u.Set("step", strconv.FormatFloat(step.Seconds(), 'f', -1, 64))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/query_range", r.VmUrl), strings.NewReader(u.Encode()))
	if err != nil {
		return nil, nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code:%d", resp.StatusCode)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result RawMsg
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(result.Data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
