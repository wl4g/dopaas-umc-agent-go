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
package docker

import (
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/constant"
	"umc-agent/pkg/constant/metric"
	"umc-agent/pkg/indicators"
	"umc-agent/pkg/logger"
	"umc-agent/pkg/transport"
)

// Docker indicators runner
func IndicatorRunner() {
	if !config.GlobalConfig.Indicator.Docker.Enabled {
		logger.Main.Debug("No enabled docker metrics runner!")
		return
	}
	logger.Main.Info("Starting docker indicators runner ...")

	// Loop monitor
	for true {
		aggregator := indicators.NewMetricAggregator("Docker")

		// Do collect.
		handleDockerMetricCollect(aggregator)

		// Send to servers.
		transport.SendMetrics(aggregator)
		time.Sleep(config.GlobalConfig.Indicator.Docker.Delay * time.Millisecond)
	}
}

func handleDockerMetricCollect(aggregator *indicators.MetricAggregator) {
	stats := getDockerStats()
	logger.Main.Info(common.ToJSONString(stats))
	for _, stat := range stats {
		stat.CpuPerc = strings.ReplaceAll(stat.CpuPerc, "%", "")
		cpuPerc, _ := strconv.ParseFloat(stat.CpuPerc, 64)
		aggregator.NewMetric(metric.DOCKER_CPU_PERC, cpuPerc).ATag(constant.TagDockerName, stat.Name)

		memUsage, _ := duel(stat.MemUsage)
		aggregator.NewMetric(metric.DOCKER_MEM_USAGE, memUsage).ATag(constant.TagDockerName, stat.Name)

		stat.MemPerc = strings.ReplaceAll(stat.MemPerc, "%", "")
		memPerc, _ := strconv.ParseFloat(stat.MemPerc, 64)
		aggregator.NewMetric(metric.DOCKER_MEM_PERC, memPerc).ATag(constant.TagDockerName, stat.Name)

		netI, netO := duel(stat.NetIO)
		aggregator.NewMetric(metric.DOCKER_NET_IN, netI).ATag(constant.TagDockerName, stat.Name)
		aggregator.NewMetric(metric.DOCKER_NET_OUT, netO).ATag(constant.TagDockerName, stat.Name)

		blockI, blockO := duel(stat.BlockIO)
		aggregator.NewMetric(metric.DOCKER_BLOCK_IN, blockI).ATag(constant.TagDockerName, stat.Name)
		aggregator.NewMetric(metric.DOCKER_BLOCK_OUT, blockO).ATag(constant.TagDockerName, stat.Name)
	}
}

// Docker stats info
func getDockerStats() []DockerStat {
	var cmd = "docker stats --no-stream --format \"{\\\"containerId\\\":\\\"{{.ID}}\\\",\\\"name\\\":\\\"{{.Name}}\\\",\\\"cpuPerc\\\":\\\"{{.CPUPerc}}\\\",\\\"memUsage\\\":\\\"{{.MemUsage}}\\\",\\\"memPerc\\\":\\\"{{.MemPerc}}\\\",\\\"netIO\\\":\\\"{{.NetIO}}\\\",\\\"blockIO\\\":\\\"{{.BlockIO}}\\\",\\\"PIDs\\\":\\\"{{.PIDs}}\\\"}\""
	logger.Main.Info("Execution docker stat", zap.String("cmd", cmd))
	s, _ := common.ExecShell(cmd)
	var dockerStats []DockerStat
	if s != "" {
		s = strings.ReplaceAll(s, "\n", ",")
		s = strings.TrimSuffix(s, ",")
		s = "[" + s + "]"
		json.Unmarshal([]byte(s), &dockerStats)
	}
	return dockerStats
}

func duel(str string) (float64, float64) {
	strs := strings.Split(str, "/")
	if len(strs) != 2 {
		return 0, 0
	}
	strs[0] = strings.TrimSpace(strs[0])
	strs[1] = strings.TrimSpace(strs[1])
	return turn2Byte(strs[0]), turn2Byte(strs[1])
}

func turn2Byte(str string) float64 {
	var result float64 = 0
	strings.Contains(str, "TB")
	if strings.Contains(str, "TB") || strings.Contains(str, "TiB") {
		strings.ReplaceAll(str, "TB", "")
		strings.ReplaceAll(str, "TiB", "")
		result, _ = strconv.ParseFloat(str, 64)
		result = result * 1024 * 1024 * 1024 * 1024
	} else if strings.Contains(str, "GB") || strings.Contains(str, "GiB") {
		strings.ReplaceAll(str, "GB", "")
		strings.ReplaceAll(str, "GiB", "")
		result, _ = strconv.ParseFloat(str, 64)
		result = result * 1024 * 1024 * 1024
	} else if strings.Contains(str, "MB") || strings.Contains(str, "MiB") {
		strings.ReplaceAll(str, "MB", "")
		strings.ReplaceAll(str, "MiB", "")
		result, _ = strconv.ParseFloat(str, 64)
		result = result * 1024 * 1024
	} else if strings.Contains(str, "KB") || strings.Contains(str, "KiB") || strings.Contains(str, "KB") {
		strings.ReplaceAll(str, "KB", "")
		strings.ReplaceAll(str, "KiB", "")
		strings.ReplaceAll(str, "KB", "")
		result, _ = strconv.ParseFloat(str, 64)
		result = result * 1024
	} else if strings.Contains(str, "B") {
		strings.ReplaceAll(str, "B", "")
		result, _ = strconv.ParseFloat(str, 64)
	} else {
		logger.Main.Error("Can not turn Byte")
		result = 0
	}
	return result
}
