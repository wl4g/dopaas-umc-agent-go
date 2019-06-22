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
package redis

import (
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/launcher"
	"umc-agent/pkg/logging"
)

func IndicatorsRunner() {
	if !config.GlobalConfig.Indicators.Virtual.Enabled {
		logging.MainLogger.Warn("No enabled redis metrics runner!")
		return
	}
	logging.MainLogger.Info("Starting redis indicators runner ...")

	// Loop monitor
	for true {
		result := getRedisStats()
		result.Meta = config.CreateMeta("redis")

		launcher.DoSendSubmit(result.Meta.Type, result)
		time.Sleep(config.GlobalConfig.Indicators.Redis.Delay * time.Millisecond)
	}
}

func getRedisStats() Redis {
	var redis Redis

	portText := config.GlobalConfig.Indicators.Redis.Ports
	portTextArr := strings.Split(portText, ",")

	propsText := config.GlobalConfig.Indicators.Redis.Properties
	propsTextArr := strings.Split(propsText, ",")

	for _, port := range portTextArr {
		re := getRedisInfo(port)
		props := wrap(re, propsTextArr)
		var redisInfo Info
		redisInfo.Port = port
		redisInfo.Properties = props
		redis.RedisInfos = append(redis.RedisInfos, redisInfo)
	}
	return redis
}

func getRedisInfo(port string) string {
	redisPwd := config.GlobalConfig.Indicators.Redis.Password
	//example: redis-cli -p 6380  -a 'zzx!@#$%' cluster info
	var comm1 string = "redis-cli -p " + port + "  -a '" + redisPwd + "' cluster info"
	//example: redis-cli -p 6380  -a 'zzx!@#$%' info all
	var comm2 string = "redis-cli -p " + port + "  -a '" + redisPwd + "' info all"
	re1, _ := common.ExecShell(comm1)
	re2, _ := common.ExecShell(comm2)
	return re1 + "\n" + re2
}

func wrap(info string, property []string) map[string]string {
	var mapInfo map[string]string
	mapInfo = make(map[string]string)
	infos := strings.Split(info, "\n")
	for _, line := range infos {
		i := strings.Index(line, ":")
		if i <= 0 {
			continue
		}
		s1 := line[0:i]
		if !common.StringsContains(property, s1) {
			continue
		}
		//fmt.Println(s1)
		s2 := line[(i + 1):len(line)]
		s2 = strings.TrimSpace(s2)
		//fmt.Println(s2)
		mapInfo[s1] = s2
	}
	return mapInfo
}
