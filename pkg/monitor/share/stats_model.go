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
package share

import "umc-agent/pkg/config"

type MetricInfo struct {
	StatMetrics []StatMetric `json:"statMetrics"`
	Timestamp   int64        `json:"timestamp"`
}

type StatMetric struct {
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
	Value  float64           `json:"value"`
}

// Create stat information.
func NewStatMetric(metric string, value float64) StatMetric {
	var statsInfo StatMetric
	statsInfo.Tags = make(map[string]string)

	statsInfo.Tags["instance"] = config.LocalHardwareAddrId
	statsInfo.Metric = config.GlobalConfig.Indicators.Namespace + "_" + metric
	statsInfo.Value = value
	return statsInfo
}

// Tags appender.
func (statInfo StatMetric) AppendTag(key string, value string) StatMetric {
	statInfo.Tags[key] = value
	return statInfo
}
