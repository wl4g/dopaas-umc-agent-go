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
	"strconv"
	"strings"
	"time"
	"umc-agent/pkg/config"
	"umc-agent/pkg/constant"
	"umc-agent/pkg/indicators"
	"umc-agent/pkg/logger"
	"umc-agent/pkg/transport"
)

func IndicatorRunner() {
	if !config.GlobalConfig.Indicator.Zookeeper.Enabled {
		logger.Main.Warn("No enabled zookeeper metrics runner!")
		return
	}
	logger.Main.Info("Starting zookeeper indicators runner ...")

	// Loop monitor
	for true {
		zookeeperAggregator := indicators.NewMetricAggregator("Zookeeper")
		now := time.Now().UnixNano() / 1e6
		zookeeperAggregator.Timestamp = now
		//gather
		Gather(zookeeperAggregator)
		// Send to servers.
		transport.DoSendSubmit(constant.Metric, &zookeeperAggregator)
		time.Sleep(config.GlobalConfig.Indicator.Zookeeper.Delay * time.Millisecond)
	}
}

func Gather(zookeeperAggregator *indicators.MetricAggregator)  {
	comm := strings.Split(config.GlobalConfig.Indicator.Zookeeper.Command, ",")
	var infoSum string
	for _, command := range comm {
		info, _ := getZkInfo(command)
		infoSum = infoSum + info
	}
	wrap(infoSum, zookeeperAggregator)
}


func wrap(info string,zookeeperAggregator *indicators.MetricAggregator) {
	infos := strings.Split(info, "\n")
	for _, line := range infos {

		i := strings.Split(line, "\t")
		if len(i) != 2 {
			continue
		}
		s1 := i[0]
		s2 := i[1]
		s2 = strings.TrimSpace(s2)

		//TODO
		if fval, err := strconv.ParseFloat(s2, 64); err == nil {
			key := "zookeeper." + strings.ReplaceAll(s1,"_",".");
			zookeeperAggregator.NewMetric(key,fval)
		}
	}
}

func getZkInfo(command string) (string, error) {
	servers := config.GlobalConfig.Indicator.Zookeeper.Servers
	logger.Main.Info("Execution zookeeper stat", zap.String("servers", servers))

	conn, err := net.Dial("tcp", servers)
	if err != nil {
		logger.Main.Error("Execution connect zookeeper failed",
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
		logger.Main.Error("Got zk metric failed!", zap.Error(err))
		if conn != nil {
			conn.Close()
		}
		return "", err
	}
	respStat := string(buf[0:cnt])
	conn.Close()

	logger.Main.Debug("Zookeeper server response", zap.String("respStat", respStat))
	return respStat, nil
}


