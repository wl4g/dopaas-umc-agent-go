/**
 * Copyright 2017 ~ 2025 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package monitor

import (
	"strings"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
)

// -----------------
// NewMetric operation.
// -----------------

// NewMetric aggregate info.
type MetricAggregate struct {
	indicatorName string
	Metrics       []Metric `json:"metrics"`
	Timestamp     int64    `json:"timestamp"`
}

// Specific stat metric.
type Metric struct {
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
	Value  float64           `json:"value"`
}

// Create stat metric info.
func NewMetric(metricName string, value float64) *Metric {
	var _metric = new(Metric)
	_metric.Tags = make(map[string]string)
	_metric.Tags["instance"] = config.LocalHardwareAddrId
	_metric.Metric = config.GlobalConfig.Indicators.Namespace + "_" + metricName
	_metric.Value = value
	return _metric
}

// NewMetric tags appender.
func (_this Metric) ATag(key string, value string) Metric {
	_this.Tags[key] = value
	return _this
}

// -----------------------
// NewMetric indicatorName operation.
// -----------------------

// New metric aggregator. Indicator names must be consistent with the global configuration,
// pay attention to case-sensitive (e.g., optional values: Redis, Kafka, Emq...),
// See: `./pkg/config/indicators_config.go#IndicatorsProperties` for names of members
// See: `./pkg/monitor/metric_config.go#NewMetric` for `config.GetConfig(...)`
func NewMetricAggregate(group string) *MetricAggregate {
	return &MetricAggregate{indicatorName: group}
}

// NewMetric aggregate to json string.
func (_this *MetricAggregate) ToJsonString() string {
	return common.ToJSONString(_this)
}

// Create metrics with the creator.
func (_this *MetricAggregate) NewMetric(metricName string, value float64) *Metric {
	// Check necessary.(Note: that the project configuration
	// structure must correspond to this.)
	var props = config.GetConfig("Indicators", _this.indicatorName, "Properties")

	// Create metricName.
	var _metric = NewMetric(metricName, value)

	// Enabled only if the configuration metricName name is included.
	if common.StringsContains(strings.Split(props.ToString(), ","), metricName) {
		__metric := *_metric // Be sure to clone?
		_this.Metrics = append(_this.Metrics, __metric)
	}

	return _metric
}
