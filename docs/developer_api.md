# 开发者API手册
* [Api格式说明- response](#----)
* [指标名称](#-----)
* [一、Predict-扩缩容规则模块](#--)
    + [1.创建单个扩缩容规则](#1-----)
    + [2.更新单个扩缩容规则](#2-----)
    + [3.查询单个扩缩容规则](#3-------)
    + [4.查询(分页)扩缩容规则列表](#4-------)
    + [5.批量删除扩缩容规则](#5-----)
    + [6.启用单个扩缩容规则](#6-----)
    + [7.禁用单个扩缩容规则](#7-------)
* [二 、指标查询](#--)
    + [1.冗余度](#1----------)
    + [2.实例数量](#2-----------)
    + [3.指标数据](#3-----------)
## Api格式说明- response
```
#返回-success
{
    "status": "success",
    "message": "",
    "data": "XXX"
}
#返回-success-分页查询
{
    "status": "success",
    "message": "",
    "data": {
        "xxx_list": [XXX],
        "pager": {
            "page_number": 1,
            "page_size": 20,
            "total": 100
        }
    }
}
#返回-failed
{
    "status": "failed",
    "message": "XXX",
    "data": ""
}
```

## 指标名称

- QPS 指标名称为 qps
- qps_metrics : 指标名称为 qps_section_factor

## 一、Predict-扩缩容规则模块

### 1.创建单个扩缩容规则 POST /api/v1/cudgx/predict/rule/create

请求参数：

| 字段                 | 类型     | 必填  | 描述      | 示例                      |
|--------------------|--------|-----|---------|-------------------------|
| name               | string | 是   | 扩缩容规则名称 | "test_predict_rule"     |
| service_name       | string | 是   | 服务名称    | "test_service"          |
| cluster_name       | string | 是   | 关联集群名称  | "test_cluster"          |
| metric_name        | string | 是   | 度量指标名称  | "qps"                   |
| benchmark_qps      | int    | 是   | 单机QPS   | 300                     |
| min_redundancy     | Int    | 是   | 最小冗余度   | 100（表示100%）             |
| max_redundancy     | int    | 是   | 最大冗余度   | 300（表示300%）             |
| min_instance_count | int    | 是   | 最小机器数   | 2                       |
| max_instance_count | int    | 是   | 最大机器数   | 5                       |
| execute_ratio      | int    | 是   | 扩缩容比例   | 30（表示30%）               |
| status             | string | 是   | 状态      | enable/disable（表示启用/禁用） |

返回： Api格式说明- response

### 2.更新单个扩缩容规则 POST /api/v1/cudgx/predict/rule/update

请求参数：

| 字段                 | 类型     | 必填  | 描述      | 示例                      |
|--------------------|--------|-----|---------|-------------------------|
| id                 | int64  | 是   | 扩缩容规则ID | 1                       |
| name               | string | 是   | 扩缩容规则名称 | "test_predict_rule"     |
| service_name       | string | 是   | 服务名称    | "test_service"          |
| cluster_name       | string | 是   | 关联集群名称  | "test_cluster"          |
| metric_name        | string | 是   | 度量指标名称  | "qps"                   |
| benchmark_qps      | int    | 是   | 单机QPS   | 300                     |
| min_redundancy     | Int    | 是   | 最小冗余度   | 100（表示100%）             |
| max_redundancy     | int    | 是   | 最大冗余度   | 300（表示300%）             |
| min_instance_count | int    | 是   | 最小机器数   | 2                       |
| max_instance_count | int    | 是   | 最大机器数   | 5                       |
| execute_ratio      | int    | 是   | 扩缩容比例   | 30（表示30%）               |
| status             | string | 是   | 状态      | enable/disable（表示启用/禁用） |

返回： Api格式说明- response

### 3.查询单个扩缩容规则 GET /api/v1/cudgx/predict/rule/:id

返回：

| 字段                 | 类型     | 必填  | 描述      | 示例                      |
|--------------------|--------|-----|---------|-------------------------|
| id                 | int64  | 是   | 扩缩容规则ID | 1                       |
| name               | string | 是   | 扩缩容规则名称 | "test_predict_rule"     |
| service_name       | string | 是   | 服务名称    | "test_service"          |
| cluster_name       | string | 是   | 关联集群名称  | "test_cluster"          |
| metric_name        | string | 是   | 度量指标名称  | "qps"                   |
| benchmark_qps      | int    | 是   | 单机QPS   | 300                     |
| min_redundancy     | Int    | 是   | 最小冗余度   | 100（表示100%）             |
| max_redundancy     | int    | 是   | 最大冗余度   | 300（表示300%）             |
| min_instance_count | int    | 是   | 最小机器数   | 2                       |
| max_instance_count | int    | 是   | 最大机器数   | 5                       |
| execute_ratio      | int    | 是   | 扩缩容比例   | 30（表示30%）               |
| status             | string | 是   | 状态      | enable/disable（表示启用/禁用） |
| created_time       | int64  | 是   | 创建时间    | 1639711726              |

### 4.查询(分页)扩缩容规则列表 GET /api/v1/cudgx/predict/rule/list?service_name=test&cluster_name=test&page_number=1&page_size=20

返回：

| 字段                 | 类型     | 必填  | 描述      | 示例                      |
|--------------------|--------|-----|---------|-------------------------|
| id                 | int64  | 是   | 扩缩容规则ID | 1                       |
| name               | string | 是   | 扩缩容规则名称 | "test_predict_rule"     |
| service_name       | string | 是   | 服务名称    | "test_service"          |
| cluster_name       | string | 是   | 关联集群名称  | "test_cluster"          |
| metric_name        | string | 是   | 度量指标名称  | "qps"                   |
| benchmark_qps      | int    | 是   | 单机QPS   | 300                     |
| min_redundancy     | Int    | 是   | 最小冗余度   | 100（表示100%）             |
| max_redundancy     | int    | 是   | 最大冗余度   | 300（表示300%）             |
| min_instance_count | int    | 是   | 最小机器数   | 2                       |
| max_instance_count | int    | 是   | 最大机器数   | 5                       |
| execute_ratio      | int    | 是   | 扩缩容比例   | 30（表示30%）               |
| status             | string | 是   | 状态      | enable/disable（表示启用/禁用） |
| created_time       | int64  | 是   | 创建时间    | 1639711726              |

分页格式：Api格式说明- response
### 5.批量删除扩缩容规则 POST /api/v1/cudgx/predict/rule/batch/delete

请求参数：

| 字段  | 类型      | 必填  | 描述           | 示例      |
|-----|---------|-----|--------------|---------|
| ids | []int64 | 是   | 扩缩容规则ID list | [1,2,3] |

返回： Api格式说明- response

### 6.启用单个扩缩容规则 POST /api/v1/cudgx/predict/rule/:id/enable

返回： Api格式说明- response

### 7.禁用单个扩缩容规则 POST /api/v1/cudgx/predict/rule/:id/disable

返回： Api格式说明- response

## 二 指标查询

### 1.冗余度 GET /api/v1/query/metric/redundancy/:metric_name?service_name=gf.cudgx.pi&cluster_name=default&begin=1640695000&end=1640695149

请求参数：

| 字段           | 二级字段       | 类型        | 描述          | 示例                         |
|--------------|------------|-----------|-------------|----------------------------|
| service_name |            | string    | 服务名称        | "test_service"             |
| metric_name  |            | string    | 度量指标名称      | "qps"/"qps_section_factor" |
| clusters     |            | []object  | 冗余度信息       |                            |
|              | cluster    | string    | 集群名称        | "default"                  |
|              | timestamps | []int64   | 时间戳（每5秒一个点） | 1639711726                 |
|              | values     | []float64 | 时间戳对应指标值    | 100                        |

返回Data字段为，具体请查看 Api格式说明- response ：

| 字段           | 二级字段       | 类型        | 描述          | 示例                         |
|--------------|------------|-----------|-------------|----------------------------|
| service_name |            | string    | 服务名称        | "test_service"             |
| metric_name  |            | string    | 度量指标名称      | "qps"/"qps_section_factor" |
| clusters     |            | []object  | 冗余度信息       |                            |
|              | cluster    | string    | 集群名称        | "default"                  |
|              | timestamps | []int64   | 时间戳（每5秒一个点） | 1639711726                 |
|              | values     | []float64 | 时间戳对应冗余度    | 1.12                       |

### 2.实例数量 GET /api/v1/query/metric/instance_count/:metric_name?service_name=gf.cudgx.pi&cluster_name=default&begin=1640695000&end=1640695149

请求字段：

| 字段           | 类型     | 必填  | 描述     | 示例                         |
|--------------|--------|-----|--------|----------------------------|
| service_name | string | 是   | 服务名称   | "test_service"             |
| metric_name  | string | 是   | 度量指标名称 | "qps"/"qps_section_factor" |
| cluster_name | string | 是   | 集群名称   | "default"                  |
| begin        | int64  | 是   | 开始时间戳  | 1640695000                 |
| end          | int64  | 是   | 结束时间戳  | 1640695149                 |

返回Data字段为，具体请查看 Api格式说明- response ：

| 字段           | 二级字段       | 类型        | 描述          | 示例                         |
|--------------|------------|-----------|-------------|----------------------------|
| service_name |            | string    | 服务名称        | "test_service"             |
| metric_name  |            | string    | 度量指标名称      | "qps"/"qps_section_factor" |
| clusters     |            | []object  | 冗余度信息       |                            |
|              | cluster    | string    | 集群名称        | "default"                  |
|              | timestamps | []int64   | 时间戳（每5秒一个点） | 1639711726                 |
|              | values     | []float64 | 时间戳对机器数     | 1.12                       |

### 3.指标数据 GET /api/v1/query/metric/load/:metric_name?service_name=gf.cudgx.pi&cluster_name=default&begin=1640695000&end=1640695149

请求字段：

| 字段           | 类型     | 必填  | 描述     | 示例                         |
|--------------|--------|-----|--------|----------------------------|
| service_name | string | 是   | 服务名称   | "test_service"             |
| metric_name  | string | 是   | 度量指标名称 | "qps"/"qps_section_factor" |
| cluster_name | string | 是   | 集群名称   | "default"                  |
| begin        | int64  | 是   | 开始时间戳  | 1640695000                 |
| end          | int64  | 是   | 结束时间戳  | 1640695149                 |

返回Data字段为，具体请查看 Api格式说明- response ：

| 字段           | 二级字段       | 类型        | 描述          | 示例                         |
|--------------|------------|-----------|-------------|----------------------------|
| service_name |            | string    | 服务名称        | "test_service"             |
| metric_name  |            | string    | 度量指标名称      | "qps"/"qps_section_factor" |
| clusters     |            | []object  | 冗余度信息       |                            |
|              | cluster    | string    | 集群名称        | "default"                  |
|              | timestamps | []int64   | 时间戳（每5秒一个点） | 1639711726                 |
|              | values     | []float64 | 时间戳对应指标值    | 100                        |
