package clients

type ExpandAndShrinkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type GetServiceScheduleResponse struct {
	Code int64           `json:"code"`
	Msg  string          `json:"msg"`
	Data ServiceSchedule `json:"data"`
}

type GetServiceClusterInstanceResponse struct {
	Code int64                           `json:"code"`
	Msg  string                          `json:"msg"`
	Data ServiceClusterInstanceCountList `json:"data"`
}

type ServiceClusterInstanceCountList struct {
	ServiceClusterList []*ServiceClusterInstanceCount `json:"service_cluster_list"`
}

type ServiceClusterInstanceCount struct {
	ServiceClusterId   int64  `json:"service_cluster_id"`
	ServiceClusterName string `json:"service_cluster_name"`
	InstanceCount      int    `json:"instance_count"`
}

type ServiceSchedule struct {
	Scheduling         bool   `json:"scheduling"`
	ServiceName        string `json:"service_name"`
	ServiceClusterName string `json:"service_cluster_name"`
}

type ServiceCluster struct {
	ServiceId          int64  `json:"service_id"`
	ServiceName        string `json:"service_name"`
	ClusterNum         int    `json:"cluster_num"`
	Language           string `json:"language"`
	ImageUrl           string `json:"image_url"`
	ServiceClusterId   int64  `json:"service_cluster_id"`
	ServiceClusterName string `json:"service_cluster_name"`
	TmplExpandId       int64  `json:"tmpl_expand_id"`
	TmplExpandName     string `json:"tmpl_expand_name"`
	Description        string `json:"description"`
	AutoDecision       string `json:"auto_decision"`
	TaskTypeStatus     string `json:"task_type_status"`
}

type ServiceByIpResponse struct {
	Code int64              `json:"code"`
	Msg  string             `json:"msg"`
	Data GetServiceByIpData `json:"data"`
}

type GetServiceByIpData struct {
	ServiceName string `json:"service_name"`
	ClusterName string `json:"cluster_name"`
}
