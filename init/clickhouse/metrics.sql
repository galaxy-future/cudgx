
CREATE DATABASE IF NOT EXISTS cudgx ;
use cudgx;
CREATE TABLE  IF NOT EXISTS cudgx.metrics_gf_test_local
(
    `metricName` LowCardinality(String),
    `serviceName` LowCardinality(String),
    `serviceRegion` LowCardinality(String),
    `serviceAz` LowCardinality(String),
    `serviceHost` LowCardinality(String),
    `labelKeys` Array(LowCardinality(String)),
    `labelValues` Array(LowCardinality(String)),
    `timestamp` Int64,
    `value` Float64
)
    ENGINE = MergeTree
PARTITION BY toYYYYMMDD(toDateTime(timestamp))
ORDER BY (serviceName, clusterName, metricName, timestamp)
SETTINGS index_granularity = 8192;

CREATE TABLE  IF NOT EXISTS cudgx.metrics_gf_test
(
    `metricName` LowCardinality(String),
    `serviceName` LowCardinality(String),
    `serviceRegion` LowCardinality(String),
    `serviceAz` LowCardinality(String),
    `serviceHost` LowCardinality(String),
    `labelKeys` Array(LowCardinality(String)),
    `labelValues` Array(LowCardinality(String)),
    `timestamp` Int64,
    `value` Float64
)
ENGINE = Distributed('metrics_gf', 'cudgx', 'metrics_gf_test_local');


