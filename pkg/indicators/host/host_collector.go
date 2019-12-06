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
package host

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/wl4g/super-devops-umc-agent/pkg/common"
	"github.com/wl4g/super-devops-umc-agent/pkg/config"
	"github.com/wl4g/super-devops-umc-agent/pkg/constant"
	"github.com/wl4g/super-devops-umc-agent/pkg/constant/metric"
	"github.com/wl4g/super-devops-umc-agent/pkg/indicators"
	"github.com/wl4g/super-devops-umc-agent/pkg/logger"
	"github.com/wl4g/super-devops-umc-agent/pkg/transport"
	"strconv"
	"strings"
	"time"
)

// Host indicators runner
func IndicatorRunner() {
	if !config.GlobalConfig.Indicator.Host.Enabled {
		logger.Main.Debug("No enabled host metrics runner!")
		return
	}
	logger.Main.Info("Starting host indicators runner ...")

	// Loop monitor
	for true {
		aggregator := indicators.NewMetricAggregator("Host")

		// Do host metric collect.
		handlePhysicalMetricCollect(aggregator)

		// Send to servers.
		transport.SendMetrics(aggregator)
		time.Sleep(config.GlobalConfig.Indicator.Host.Delay * time.Millisecond)
	}
}

func handlePhysicalMetricCollect(physicalAggregator *indicators.MetricAggregator) {
	//cpu
	gatherCpu(physicalAggregator)

	//mem
	gatherMem(physicalAggregator)

	//disk
	gatherDisk(physicalAggregator)

	//net
	gatherNet(physicalAggregator)
}

//cpu
func gatherCpu(physicalAggregator *indicators.MetricAggregator) {
	p, _ := cpu.Percent(0, false)
	physicalAggregator.NewMetric(metric.PHYSICAL_CPU, p[0])
}

//mem
func gatherMem(physicalAggregator *indicators.MetricAggregator) {
	v, _ := mem.VirtualMemory()
	physicalAggregator.NewMetric(metric.PHYSICAL_MEM_TOTAL, float64(v.Total))
	physicalAggregator.NewMetric(metric.PHYSICAL_MEM_FREE, float64(v.Free))
	physicalAggregator.NewMetric(metric.PHYSICAL_MEM_USED_PERCENT, float64(v.UsedPercent))
	physicalAggregator.NewMetric(metric.PHYSICAL_MEM_USED, float64(v.Used))
	physicalAggregator.NewMetric(metric.PHYSICAL_MEM_CACHE, float64(v.Cached))
	physicalAggregator.NewMetric(metric.PHYSICAL_MEM_BUFFERS, float64(v.Buffers))
}

// disk
func gatherDisk(physicalAggregator *indicators.MetricAggregator) {
	partitionStats, _ := disk.Partitions(false)
	for _, value := range partitionStats {
		mountpoint := value.Mountpoint
		usageStat, _ := disk.Usage(mountpoint)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_TOTAL, float64(usageStat.Total)).ATag(constant.TagDiskDevice, value.Device)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_FREE, float64(usageStat.Free)).ATag(constant.TagDiskDevice, value.Device)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_USED, float64(usageStat.Used)).ATag(constant.TagDiskDevice, value.Device)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_USED_PERCENT, float64(usageStat.UsedPercent)).ATag(constant.TagDiskDevice, value.Device)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_INODES_TOTAL, float64(usageStat.InodesTotal)).ATag(constant.TagDiskDevice, value.Device)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_INODES_USED, float64(usageStat.InodesUsed)).ATag(constant.TagDiskDevice, value.Device)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_INODES_FREE, float64(usageStat.InodesFree)).ATag(constant.TagDiskDevice, value.Device)
		physicalAggregator.NewMetric(metric.PHYSICAL_DISK_INODES_USED_PERCENT, float64(usageStat.InodesUsedPercent)).ATag(constant.TagDiskDevice, value.Device)
	}
}

// Network stats info
func gatherNet(physicalAggregator *indicators.MetricAggregator) {
	ports := strings.Split(config.GlobalConfig.Indicator.Host.NetPorts, ",")
	for _, p := range ports {
		re := common.GetNetworkInterfaces(p)
		res := strings.Split(re, " ")
		if len(res) == 9 {
			up, _ := strconv.ParseFloat(res[0], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_UP, up).ATag(constant.TagNetPort, p)
			down, _ := strconv.ParseFloat(res[2], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_DOWN, down).ATag(constant.TagNetPort, p)
			count, _ := strconv.ParseFloat(res[2], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_COUNT, count).ATag(constant.TagNetPort, p)
			estab, _ := strconv.ParseFloat(res[3], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_ESTAB, estab).ATag(constant.TagNetPort, p)
			closeWait, _ := strconv.ParseFloat(res[4], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_CLOSE_WAIT, closeWait).ATag(constant.TagNetPort, p)
			timeWait, _ := strconv.ParseFloat(res[5], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_TIME_WAIT, timeWait).ATag(constant.TagNetPort, p)
			close, _ := strconv.ParseFloat(res[6], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_CLOSE, close).ATag(constant.TagNetPort, p)
			listen, _ := strconv.ParseFloat(res[7], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_LISTEN, listen).ATag(constant.TagNetPort, p)
			closing, _ := strconv.ParseFloat(res[8], 64)
			physicalAggregator.NewMetric(metric.PHYSICAL_NET_CLOSING, closing).ATag(constant.TagNetPort, p)
		}
	}
}
