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
	"fmt"
	"github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"umc-agent/pkg/log"
)

var commandPath = "./pkg/resources/net.port.sh.txt"
var command string

// Get network interfaces.
// e.g. var sumCommand = "ss -n sport == 22|awk '{sumup += $3};{sumdo += $4};END {print sumup,sumdo}'"
func GetNetworkInterfaces(port string) string {
	if command == "" {
		command = ReadFileToString(commandPath)
	}
	cmd := strings.Replace(command, "#{port}", port, -1)
	s, _ := ExecShell(cmd)

	log.MainLogger.Debug("Exec complete", zap.String("result", s))
	return s
}

// Get hardware information such as network card as
// physical host identification.
func GetPhysicalId(netcard string) string {
	var physicalId string
	nets, _ := net.Interfaces()
	var found = false
	for _, value := range nets {
		if strings.EqualFold(netcard, value.Name) {
			hardwareAddr := value.HardwareAddr
			log.MainLogger.Info("Found network information", zap.String("hardwareAddr", hardwareAddr))

			physicalId = hardwareAddr
			reg := regexp.MustCompile(`(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})(\.(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})){3}`)
			for _, addr := range value.Addrs {
				add := addr.Addr
				if len(reg.FindAllString(add, -1)) > 0 {
					fmt.Println("found ip " + add)
					//id = add+" "+id
					a := strings.Split(add, "/")
					if len(a) >= 2 {
						physicalId = a[0]
					} else {
						physicalId = add
					}
					found = true
					break
				}
			}
		}
	}
	if !found {
		panic("net found ip,Please check the net conf")
	}
	return physicalId
}
