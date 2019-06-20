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

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
	"umc-agent/pkg/config"
	"umc-agent/pkg/launcher"
	"umc-agent/pkg/monitor/physical"
	"umc-agent/pkg/monitor/virtual"
)

func TotalThread() {
	for true {
		var result Total

		result.Id = PhysicalId
		result.Type = "total"

		v, _ := mem.VirtualMemory()
		result.Mem = v

		p, _ := cpu.Percent(0, false)
		result.Cpu = p

		disks := physical.GetDisks()
		result.DiskInfos = disks

		n := physical.GetNetInfo()
		result.NetInfos = n

		dockerInfo := virtual.GetDocker()
		result.DockerInfos = dockerInfo

		launcher.DoSendSubmit("total", result)
		time.Sleep(config.GlobalPropertiesObj.PhysicalPropertiesObj.Delay * time.Millisecond)
	}
}
