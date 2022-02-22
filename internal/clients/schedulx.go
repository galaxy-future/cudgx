package clients

import (
	"encoding/json"
	"fmt"
	"golang.org/x/sync/singleflight"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	lru "github.com/hashicorp/golang-lru"
	"go.uber.org/zap"
)

type LRUCache struct {
	*lru.Cache
}

var (
	cache = newLRUCache(1000)
	_     = cleanCachePer3Mins()
	sf    singleflight.Group
)

func newLRUCache(size int) *LRUCache {
	l, err := lru.New(size)
	if err != nil {
		return nil
	}
	return &LRUCache{Cache: l}
}

func NewSchedulxClient(serverAddress string) *Client {
	return &Client{
		ServerAddress: serverAddress,
		HttpClient: &http.Client{
			Timeout:   5000 * time.Millisecond,
			Transport: XclientRoundTripper{r: http.DefaultTransport},
		},
	}
}

type XclientRoundTripper struct {
	r http.RoundTripper
}

func (x XclientRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	token, err := authXClient()
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", "Bearer: "+token)
	return x.r.RoundTrip(r)
}

type SchedulxResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

// CanServiceSchedule 判断该服务集群是否可以调度
func CanServiceSchedule(serviceName, clusterName string) (bool, error) {
	if err := validateNames(serviceName, clusterName); err != nil {
		return false, err
	}
	resp, err := schedulxClient.HttpClient.Get(fmt.Sprintf("%s/api/v1/schedulx/service/scheduling?service_name=%s&service_cluster_name=%s", schedulxClient.ServerAddress, serviceName, clusterName))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var response GetServiceScheduleResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return false, err
	}
	if response.Code != http.StatusOK {
		err = fmt.Errorf("http code:%v | msg:%v", response.Code, response.Msg)
		return false, err
	}
	return !response.Data.Scheduling, nil
}

// GetServiceInstanceCount 获取该服务集群运行中的实例数
func GetServiceInstanceCount(serviceName, clusterName string) (int, error) {
	if err := validateNames(serviceName, clusterName); err != nil {
		return 0, err
	}
	resp, err := schedulxClient.HttpClient.Get(fmt.Sprintf("%s/api/v1/schedulx/instance/count?service_name=%s&service_cluster_name=%s", schedulxClient.ServerAddress, serviceName, clusterName))
	if err != nil {
		return 0, err
	}
	var instanceCount int
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var response GetServiceClusterInstanceResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return 0, err
	}
	if response.Code != http.StatusOK {
		err = fmt.Errorf("http code:%v | msg:%v", response.Code, response.Msg)
		return 0, err
	}
	for _, sc := range response.Data.ServiceClusterList {
		instanceCount += sc.InstanceCount
	}
	return instanceCount, nil
}

// ExpandService 扩容服务集群
func ExpandService(serviceName, clusterName string, count int) error {
	if err := validateParams(serviceName, clusterName, count); err != nil {
		return err
	}
	resp, err := schedulxClient.HttpClient.Get(fmt.Sprintf("%s/api/v1/schedulx/service/expand?service_name=%s&service_cluster=%s&count=%d&exec_type=auto", schedulxClient.ServerAddress, serviceName, clusterName, count))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response ExpandAndShrinkResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return err
	}
	if response.Code != http.StatusOK {
		err = fmt.Errorf("http code:%v | msg:%v", response.Code, response.Msg)
		return err
	}
	logger.GetLogger().Info(consts.SchedulxExpandSuccess, zap.String("service_name", serviceName), zap.String("service_cluster", clusterName), zap.Int("count", count))
	return nil
}

// ShrinkService 缩容服务集群
func ShrinkService(serviceName, clusterName string, count int) error {
	if err := validateParams(serviceName, clusterName, count); err != nil {
		return err
	}
	resp, err := schedulxClient.HttpClient.Get(fmt.Sprintf("%s/api/v1/schedulx/service/shrink?service_name=%s&service_cluster=%s&count=%d&exec_type=auto", schedulxClient.ServerAddress, serviceName, clusterName, count))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response ExpandAndShrinkResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return err
	}
	if response.Code != http.StatusOK {
		err = fmt.Errorf("http code:%v | msg:%v", response.Code, response.Msg)
		return err
	}
	logger.GetLogger().Info(consts.SchedulxShrinkSuccess, zap.String("service_name", serviceName), zap.String("service_cluster", clusterName), zap.Int("count", count))
	return nil
}

// validateParams 参数校验
func validateParams(serviceName, clusterName string, instanceCount int) error {
	if err := validateNames(serviceName, clusterName); err != nil {
		return err
	}
	if instanceCount <= 0 {
		return fmt.Errorf("实例数应大于0")
	}
	return nil
}

func validateNames(serviceName, clusterName string) error {
	if serviceName == "" {
		return fmt.Errorf("服务名称不能为空")
	}
	if clusterName == "" {
		return fmt.Errorf("集群名称不能为空")
	}
	return nil
}

type localServiceIpCache struct {
	data   map[string]GetServiceByIpData
	expire time.Duration
}

// doGetServiceByIp 通过 ip 获取服务名称.
func doGetServiceByIp(ip string) (GetServiceByIpData, error) {
	resp, err := schedulxClient.HttpClient.Get(fmt.Sprintf("%s/api/v1/schedulx/instance/service?ip_inner=%s", schedulxClient.ServerAddress, ip))
	if err != nil {
		return GetServiceByIpData{}, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GetServiceByIpData{}, err

	}
	var response ServiceByIpResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return GetServiceByIpData{}, err
	}
	if response.Code != http.StatusOK {
		err = fmt.Errorf("http code:%v | msg:%v", response.Code, response.Msg)
		return GetServiceByIpData{}, err
	}
	return response.Data, nil
}

func GetServiceByIp(ip string) (GetServiceByIpData, error) {
	srv, ok := cache.Get(ip)
	if ok {
		d, _ := srv.(GetServiceByIpData)
		return d, nil
	}

	data, err, _ := sf.Do(ip, func() (interface{}, error) {
		res, err := doGetServiceByIp(ip)
		if err != nil {
			return nil, err
		}
		cache.Add(ip, res)
		return res, nil
	})
	if err != nil {
		return GetServiceByIpData{}, err
	}
	d, _ := data.(GetServiceByIpData)
	return d, nil
}

func cleanCachePer3Mins() interface{} {
	go func() {
		var l sync.Mutex
		for {
			l.Lock()
			func() {
				defer l.Unlock()
				cache.Cache, _ = lru.New(1000)
			}()
			time.Sleep(time.Minute * 3)
		}
	}()
	return nil
}
