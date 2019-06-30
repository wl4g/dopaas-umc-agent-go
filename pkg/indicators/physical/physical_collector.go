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
package physical

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/indicators"
	"umc-agent/pkg/logger"
	"umc-agent/pkg/transport"
)

// Physical indicators runner
func IndicatorRunner() {
	if !config.GlobalConfig.Indicators.Physical.Enabled {
		logger.Main.Warn("No enabled physical metrics runner!")
		return
	}
	logger.Main.Info("Starting physical indicators runner ...")

	// Loop monitor
	for true {
		var result indicators.TotalStat
		result.Meta = indicators.CreateMeta("physical")

		p, _ := cpu.Percent(0, false)
		result.Cpu = p

		v, _ := mem.VirtualMemory()
		result.Mem = v

		result.DiskStats = getDiskStatsInfo()
		result.NetStats = getNetworkStatsInfo()

		transport.DoSendSubmit("physical", result)
		time.Sleep(config.GlobalConfig.Indicators.Docker.Delay * time.Millisecond)
	}
}

// Disks stats info
func getDiskStatsInfo() []indicators.DiskStat {
	partitionStats, _ := disk.Partitions(false)
	var disks []indicators.DiskStat
	for _, value := range partitionStats {
		var disk1 indicators.DiskStat
		mountpoint := value.Mountpoint
		usageStat, _ := disk.Usage(mountpoint)
		disk1.PartitionStat = value
		disk1.Usage = *usageStat
		disks = append(disks, disk1)
	}
	return disks
}

// Network stats info
func getNetworkStatsInfo() []indicators.NetworkStat {
	ports := strings.Split(config.GlobalConfig.Indicators.Physical.NetPorts, ",")
	//n, _ := net.IOCounters(true)
	//fmt.Println(n)
	//te, _ := net.Interfaces()
	//fmt.Println(te)
	var n []indicators.NetworkStat
	for _, p := range ports {
		re := common.GetNetworkInterfaces(p)
		res := strings.Split(re, " ")
		if len(res) == 9 {
			var netinfo indicators.NetworkStat
			netinfo.Port, _ = strconv.Atoi(p)
			netinfo.Up, _ = strconv.Atoi(res[0])
			netinfo.Down, _ = strconv.Atoi(res[1])
			netinfo.Count, _ = strconv.Atoi(res[2])
			netinfo.Estab, _ = strconv.Atoi(res[3])
			netinfo.CloseWait, _ = strconv.Atoi(res[4])
			netinfo.TimeWait, _ = strconv.Atoi(res[5])
			netinfo.Close, _ = strconv.Atoi(res[6])
			netinfo.Listen, _ = strconv.Atoi(res[7])
			netinfo.Closing, _ = strconv.Atoi(res[8])
			n = append(n, netinfo)
		}
	}
	return n
}