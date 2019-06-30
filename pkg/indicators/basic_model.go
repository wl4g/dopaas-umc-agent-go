/**
 * Copyright 2017 ~ 2025 the original author or authors[983708408@qq.com].
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
package indicators

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"umc-agent/pkg/config"
)

// Deprecated: Metric wrapper has been changed to uniform aggregation metrics,
// see: `indicators. New Metric Aggregator`.
type TotalStat struct {
	Meta      MetaInfo               `json:"meta"`
	Mem       *mem.VirtualMemoryStat `json:"memInfo"`
	Cpu       []float64              `json:"cpu"`
	DiskStats []DiskStat             `json:"diskInfos"`
	NetStats  []NetworkStat          `json:"netInfos"`
}

// Deprecated: Metric wrapper has been changed to uniform aggregation metrics,
// see: `indicators. New Metric Aggregator`.
type DiskStat struct {
	PartitionStat disk.PartitionStat `json:"partitionStat"`
	Usage         disk.UsageStat     `json:"usage"`
}

// Deprecated: Metric wrapper has been changed to uniform aggregation metrics,
// see: `indicators. New Metric Aggregator`.
type NetworkStat struct {
	Port      int `json:"port"`
	Up        int `json:"up"`
	Down      int `json:"down"`
	Count     int `json:"count"`
	Estab     int `json:"estab"`
	CloseWait int `json:"closeWait"`
	TimeWait  int `json:"timeWait"`
	Close     int `json:"close"`
	Listen    int `json:"listen"`
	Closing   int `json:"closing"`
}

// Deprecated: Metric wrapper has been changed to uniform aggregation metrics,
// see: `indicators. New Metric Aggregator`.
type MetaInfo struct {
	Id        string `json:"physicalId"`
	Type      string `json:"type"`
	Namespace string `json:"namespace"`
}

// Deprecated: Metric wrapper has been changed to uniform aggregation metrics,
// see: `indicators. New Metric Aggregator`.
// Create meta info
func CreateMeta(metaType string) MetaInfo {
	meta := MetaInfo{
		Id:        config.LocalHardwareAddrId,
		Type:      metaType,
		Namespace: config.GlobalConfig.Indicators.Namespace,
	}
	return meta
}
