//
// Copyright 2017 ~ 2025 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package virtual

import (
	"encoding/json"
	"strings"
	"umc-agent/pkg/common"
)

func GetDocker() []DockerInfo {
	var command string = "docker stats --no-stream --format \"{\\\"containerId\\\":\\\"{{.ID}}\\\",\\\"name\\\":\\\"{{.Name}}\\\",\\\"cpuPerc\\\":\\\"{{.CPUPerc}}\\\",\\\"memUsage\\\":\\\"{{.MemUsage}}\\\",\\\"memPerc\\\":\\\"{{.MemPerc}}\\\",\\\"netIO\\\":\\\"{{.NetIO}}\\\",\\\"blockIO\\\":\\\"{{.BlockIO}}\\\",\\\"PIDs\\\":\\\"{{.PIDs}}\\\"}\""
	s, _ := common.ExecShell(command)
	var dockerInfos []DockerInfo
	if s != "" {
		s = strings.ReplaceAll(s, "\n", ",")
		s = strings.TrimSuffix(s, ",")
		s = "[" + s + "]"
		json.Unmarshal([]byte(s), &dockerInfos)
	}
	return dockerInfos
}
