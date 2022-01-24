package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"github.com/galaxy-future/cudgx/internal/predict/service"
	"github.com/galaxy-future/cudgx/internal/request"
	"github.com/galaxy-future/cudgx/internal/response"
	"github.com/gin-gonic/gin"
)

// GetPredictRule 获取扩缩容规则
func GetPredictRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("未指定规则id"))
		return
	}
	predictRule, err := service.GetPredictRuleById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(predictRule))
}

// CreatePredictRule 创建扩缩容规则
func CreatePredictRule(c *gin.Context) {
	req := request.CreatePredictRuleRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(response.ParamError))
		return
	}
	if strings.ToLower(req.MetricName) != consts.QPSMetricsName &&
		strings.ToLower(req.MetricName) != consts.LatencySectionFactorMetricsName {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(response.MetricNameError))
		return
	}
	err := service.CreatePredictRule(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(nil))
}

// UpdatePredictRule 更新扩缩容规则
func UpdatePredictRule(c *gin.Context) {
	req := request.UpdatePredictRuleRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(response.ParamError))
		return
	}
	if strings.ToLower(req.MetricName) != consts.QPSMetricsName &&
		strings.ToLower(req.MetricName) != consts.LatencySectionFactorMetricsName {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(response.MetricNameError))
		return
	}
	err := service.UpdatePredictRuleById(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(nil))
}

// BatchDeletePredictRule 批量删除扩缩容规则
func BatchDeletePredictRule(c *gin.Context) {
	req := request.BatchDeletePredictRuleRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(response.ParamError))
		return
	}
	err := service.DeletePredictRuleById(&req)
	if err != nil {
		response.MkFailedResponse(err.Error())
		return
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(nil))
}

// ListPredictRules 获取扩缩容列表
func ListPredictRules(c *gin.Context) {
	serviceName := c.Query("service_name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("服务名称不能为空"))
		return
	}
	clusterName := c.Query("cluster_name")
	pageNumber, pageSize, err := getPager(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(response.ParamError))
		return
	}
	predictRules, total, err := service.ListPredictRules(serviceName, clusterName, pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(err.Error()))
		return
	}
	pager := response.Pager{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      total,
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(&response.ListPredictRuleResponse{
		PredictRuleList: predictRules,
		Pager:           pager,
	}))
}

// EnablePredictRule 启用扩缩容规则
func EnablePredictRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("未指定规则id"))
		return
	}
	err = service.UpdatePredictRuleStatus(id, consts.RuleStatusEnable)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(nil))
}

// DisablePredictRule 禁用扩缩容规则
func DisablePredictRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse("未指定规则id"))
		return
	}
	err = service.UpdatePredictRuleStatus(id, consts.RuleStatusDisable)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MkFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.MkSuccessResponse(nil))
}

func getPager(c *gin.Context) (pageNumber int, pageSize int, err error) {
	pageNumber, err = strconv.Atoi(c.Query("page_number"))
	if err != nil {
		return 0, 0, err
	}
	if pageNumber < 1 {
		pageNumber = 1
	}
	pageSize, err = strconv.Atoi(c.Query("page_size"))
	if err != nil {
		return 0, 0, err
	}
	if pageSize < 1 || pageSize > 20 {
		pageSize = 20
	}
	return pageNumber, pageSize, nil
}
