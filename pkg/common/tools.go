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
package common

import (
	"github.com/shirou/gopsutil/net"
	"github.com/wl4g/super-devops-umc-agent/pkg/constant"
	"regexp"
	"strings"
)

var (
	// Network state commands.
	_netPortTcpStateCmd = constant.DefaultNetPortStateCmd
)

// Get network interfaces.
// e.g. var sumCommand = "ss -n sport == 22|awk '{sumup += $3};{sumdo += $4};END {print sumup,sumdo}'"
func GetNetworkInterfaces(port string) string {
	cmd := strings.Replace(_netPortTcpStateCmd, "#{port}", port, -1)
	s, _ := ExecShell(cmd)

	//fmt.Printf("Execution completed for - '%s'", s)
	return s
}

// Get hardware information such as network card as
// host host identification.
func GetHardwareAddr(netcard string) string {
	var hardwareId = ""
	nets, _ := net.Interfaces()
ok:
	for _, netStat := range nets {
		if strings.EqualFold(netcard, netStat.Name) {
			//fmt.Printf("Found network interfaces for - '%s'", netStat.HardwareAddr)

			reg := regexp.MustCompile(`(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})(\.(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})){3}`)
			for _, addr := range netStat.Addrs {
				_addr := addr.Addr
				if len(reg.FindAllString(_addr, -1)) > 0 {
					a := strings.Split(_addr, "/")
					if len(a) >= 2 {
						hardwareId = a[0]
					} else {
						hardwareId = _addr
					}
					break ok
				}
			}
		}
	}
	return hardwareId
}
