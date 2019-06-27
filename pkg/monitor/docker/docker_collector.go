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
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/logger"
	"umc-agent/pkg/transport"
)

// Docker indicators runner
func IndicatorRunner() {
	if !config.GlobalConfig.Indicators.Docker.Enabled {
		logger.Main.Warn("No enabled docker metrics runner!")
		return
	}
	logger.Main.Info("Starting docker indicators runner ...")

	// Loop monitor
	for true {
		var result DockerStatInfo
		result.Meta = config.CreateMeta("docker")

		stats := getDockerStats()
		logger.Main.Info(common.ToJSONString(stats))
		result.DockerStats = stats

		transport.DoSendSubmit("docker", result)
		time.Sleep(config.GlobalConfig.Indicators.Docker.Delay * time.Millisecond)
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
