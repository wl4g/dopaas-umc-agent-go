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
package zookeeper

import (
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/logging"
	"umc-agent/pkg/transport"
)

func IndicatorRunner() {
	if !config.GlobalConfig.Indicators.Docker.Enabled {
		logging.MainLogger.Warn("No enabled zookeeper metrics runner!")
		return
	}
	logging.MainLogger.Info("Starting zookeeper indicators runner ...")

	// Loop monitor
	for true {
		result := getZookeeperStats()
		result.Meta = config.CreateMeta("zookeeper")

		transport.DoSendSubmit(result.Meta.Type, result)
		time.Sleep(config.GlobalConfig.Indicators.Zookeeper.Delay * time.Millisecond)
	}
}

func getZookeeperStats() Zookeeper {
	var result Zookeeper

	comm := strings.Split(config.GlobalConfig.Indicators.Zookeeper.Command, ",")
	props := strings.Split(config.GlobalConfig.Indicators.Zookeeper.Properties, ",")

	var infoSum string
	for _, command := range comm {
		info, _ := getZkInfo(command)
		infoSum = infoSum + info
	}
	infos := wrap(infoSum, props)
	result.Properties = infos
	return result

}

func getZkInfo(command string) (string, error) {
	servers := config.GlobalConfig.Indicators.Zookeeper.Servers
	logging.MainLogger.Info("Execution zookeeper stat", zap.String("servers", servers))

	conn, err := net.Dial("tcp", servers)
	if err != nil {
		logging.MainLogger.Error("Execution connect zookeeper failed",
			zap.String("servers", servers), zap.Error(err))
		return "", err
	}

	cmd := command
	// Clean trim space
	cmd = strings.TrimSpace(cmd)

	// for test: Console input
	//reader := bufio.NewReader(os.Stdin)

	// Write commands
	conn.Write([]byte(cmd))

	// Read response
	buf := make([]byte, 1024)
	cnt, err := conn.Read(buf)

	if err != nil {
		logging.MainLogger.Error("Got zk metric failed!", zap.Error(err))
		if conn != nil {
			conn.Close()
		}
		return "", err
	}
	respStat := string(buf[0:cnt])
	conn.Close()

	logging.MainLogger.Debug("Zookeeper server response", zap.String("respStat", respStat))
	return respStat, nil
}

func wrap(info string, property []string) map[string]string {
	var mapInfo map[string]string
	mapInfo = make(map[string]string)
	infos := strings.Split(info, "\n")
	for _, line := range infos {

		i := strings.Split(line, "\t")
		if len(i) != 2 {
			continue
		}
		s1 := i[0]
		if !common.StringsContains(property, s1) {
			continue
		}
		s2 := i[1]
		s2 = strings.TrimSpace(s2)
		mapInfo[s1] = s2
	}
	return mapInfo
}
