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
package main

import (
	"flag"
	"go.uber.org/zap"
	"time"
	"umc-agent/pkg/config"
	"umc-agent/pkg/constant"
	"umc-agent/pkg/launcher"
	"umc-agent/pkg/logging"
	"umc-agent/pkg/monitor/physical"
	"umc-agent/pkg/monitor/redis"
	"umc-agent/pkg/monitor/virtual"
	"umc-agent/pkg/monitor/zookeeper"
)

func init() {
	var confPath = constant.DefaultConfigPath

	// Command config path
	flag.StringVar(&confPath, "c", constant.DefaultConfigPath, "Config must is required!")
	flag.Parse()
	//flag.Usage()

	logging.MainLogger.Info("Initialize config file", zap.String("confPath", confPath))

	// Init global config.
	config.InitGlobalConfig(confPath)

	// Init kafka launcher.(if necessary)
	launcher.InitKafkaLauncherIfNecessary()
}

func main() {
	go physical.IndicatorsRunner()
	go virtual.IndicatorsRunner()
	go redis.IndicatorsRunner()
	go zookeeper.IndicatorsRunner()

	for true {
		time.Sleep(100000 * time.Millisecond)
	}
}
