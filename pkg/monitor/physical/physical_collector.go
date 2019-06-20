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
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/launcher"
	"umc-agent/pkg/monitor/share"
)

// Memory indicators runner
func MemIndicatorsRunner() {
	for true {
		var result share.Total
		v, _ := mem.VirtualMemory()
		fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
		//fmt.Println(v)
		result.Id = share.PhysicalId
		result.Type = "mem"
		result.Mem = v
		fmt.Println("result = " + common.ToJSONString(result))
		launcher.DoSendSubmit("mem", result)
		time.Sleep(config.GlobalPropertiesObj.PhysicalPropertiesObj.Delay * time.Millisecond)
	}
}

// CPU indicators runner
func CpuIndicatorsRunner() {
	for true {
		var result share.Total
		p, _ := cpu.Percent(0, false)
		//p, _ := cpu.Times(true)
		fmt.Println(p)
		/*pa, _ := cpu.Percent(10000* time.Millisecond, true)
		fmt.Println(pa)*/
		result.Id = share.PhysicalId
		result.Type = "cpu"
		result.Cpu = p
		launcher.DoSendSubmit("cpu", result)
		time.Sleep(config.GlobalPropertiesObj.PhysicalPropertiesObj.Delay * time.Millisecond)
	}
}

// Disk indicators runner
func DiskIndicatorsRunner() {
	for true {
		var result share.Total
		disks := GetDiskStatsInfo()
		fmt.Println(disks)
		result.Id = share.PhysicalId
		result.Type = "disk"
		result.DiskInfos = disks
		launcher.DoSendSubmit("disk", result)
		time.Sleep(config.GlobalPropertiesObj.PhysicalPropertiesObj.Delay * time.Millisecond)
	}
}

// Disks stats info
func GetDiskStatsInfo() []share.DiskInfo {
	partitionStats, _ := disk.Partitions(false)
	var disks []share.DiskInfo
	for _, value := range partitionStats {
		var disk1 share.DiskInfo
		mountpoint := value.Mountpoint
		usageStat, _ := disk.Usage(mountpoint)
		disk1.PartitionStat = value
		disk1.Usage = *usageStat
		disks = append(disks, disk1)
	}
	return disks
}

// Network indicators runner
func NetIndicatorsRunner() {
	for true {
		var result share.Total
		n := GetNetworkStatsInfo()
		result.Id = share.PhysicalId
		result.Type = "net"
		result.NetInfos = n
		launcher.DoSendSubmit("net", result)
		time.Sleep(config.GlobalPropertiesObj.PhysicalPropertiesObj.Delay * time.Millisecond)
	}
}

// Network stats info
func GetNetworkStatsInfo() []share.NetInfo {
	ports := strings.Split(config.GlobalPropertiesObj.PhysicalPropertiesObj.GatherPort, ",")
	//n, _ := net.IOCounters(true)
	//fmt.Println(n)
	//te, _ := net.Interfaces()
	//fmt.Println(te)
	var n []share.NetInfo
	for _, p := range ports {
		re := common.GetNetworkInterfaces(p)
		res := strings.Split(re, " ")
		if len(res) == 9 {
			var netinfo share.NetInfo
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
