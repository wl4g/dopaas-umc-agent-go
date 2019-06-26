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
	"fmt"
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/logging"
	"umc-agent/pkg/transport"
)

func IndicatorRunner() {
	if !config.GlobalConfig.Indicators.Docker.Enabled {
		logging.MainLogger.Warn("No enabled redis metrics runner!")
		return
	}
	logging.MainLogger.Info("Starting redis indicators runner ...")

	// Loop monitor
	for true {
		result := getRedisStats()
		result.Meta = config.CreateMeta("redis")

		transport.DoSendSubmit(result.Meta.Type, result)
		time.Sleep(config.GlobalConfig.Indicators.Redis.Delay * time.Millisecond)
	}
}

func getRedisStats() Redis {
	var redis Redis

	servers := config.GlobalConfig.Indicators.Redis.Servers
	serversArr := strings.Split(servers, ",")

	propsText := config.GlobalConfig.Indicators.Redis.Properties
	propsTextArr := strings.Split(propsText, ",")

	for _, addr := range serversArr {
		re := getRedisInfo(strings.Split(addr, ":"))
		props := wrap(re, propsTextArr)
		var redisInfo Info
		redisInfo.Port = addr
		redisInfo.Properties = props
		redis.RedisInfos = append(redis.RedisInfos, redisInfo)
	}
	return redis
}

func getRedisInfo(parts []string) string {
	redisPwd := config.GlobalConfig.Indicators.Redis.Password
	// e.g: redis-cli -h localhost -p 6380 -a '123456' cluster info
	// e.g: redis-cli -h localhost -p 6380 -a '123456' info all
	var cmd1 = fmt.Sprintf("redis-cli -h %s -p %s -a %s cluster info", parts[0], parts[1], redisPwd)
	var cmd2 = fmt.Sprintf("redis-cli -h %s -p %s -a %s info all", parts[0], parts[1], redisPwd)
	ret1, _ := common.ExecShell(cmd1)
	ret2, _ := common.ExecShell(cmd2)
	return ret1 + "\n" + ret2
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
