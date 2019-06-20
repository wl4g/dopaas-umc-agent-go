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
package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"umc-agent/pkg/constant"
	"umc-agent/pkg/log"
)

// Global config properties.
var GlobalConfig GlobalProperties

type GlobalProperties struct {
	Launcher   LauncherProperties
	Indicators IndicatorsProperties
}

// ---------------------
// Launcher properties.
// ---------------------

type LauncherProperties struct {
	Http  HttpLauncherProperties
	Kafka KafkaLauncherProperties
}

type HttpLauncherProperties struct {
	ServerGateway string `yaml:"launcher.http.server-gateway"`
}

type KafkaLauncherProperties struct {
	Enabled          bool   `yaml:"launcher.kafka.enabled"`
	BootstrapServers string `yaml:"launcher.kafka.bootstrap.servers"`
	Topic            string `yaml:"launcher.kafka.topic"`
	Partitions       int32  `yaml:"launcher.kafka.partitions"`
}

// ----------------------
// Indicators properties.
// ----------------------

type IndicatorsProperties struct {
	Netcard   string `yaml:"netcard"`
	Physical  PhysicalIndicatorProperties
	Virtual   VirtualIndicatorProperties
	Redis     RedisIndicatorProperties
	Zookeeper ZookeeperIndicatorProperties
	Kafka     KafkaIndicatorProperties
	Etcd      EtcdIndicatorProperties
	Emq       EmqIndicatorProperties
}

// Indicators physical properties.
type PhysicalIndicatorProperties struct {
	Delay     time.Duration `yaml:"indicators.physical.delay"`
	RangePort string        `yaml:"indicators.physical.range-port"`
}

// Indicators virtual properties.
type VirtualIndicatorProperties struct {
	Delay time.Duration `yaml:"indicators.virtual.delay"`
}

// Indicators redis properties.
type RedisIndicatorProperties struct {
	Delay time.Duration `yaml:"indicators.redis.delay"`
}

// Indicators zookeeper properties.
type ZookeeperIndicatorProperties struct {
	Delay time.Duration `yaml:"indicators.zookeeper.delay"`
}

// Indicators kafka properties.
type KafkaIndicatorProperties struct {
	Delay time.Duration `yaml:"indicators.kafka.delay"`
}

// Indicators etcd properties.
type EtcdIndicatorProperties struct {
	Delay time.Duration `yaml:"indicators.etcd.delay"`
}

// Indicators emq properties.
type EmqIndicatorProperties struct {
	Delay time.Duration `yaml:"indicators.emq.delay"`
}

// Initialize global config properties.
func InitGlobalConfig(path string) {
	// Set defaults
	setDefaults()

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.MainLogger.Info("yamlFile.Get err - ", zap.Error(err))
	}

	err = yaml.Unmarshal(yamlFile, &GlobalConfig)
	if err != nil {
		log.MainLogger.Error("Unmarshal error", zap.Error(err))
	}
}

// Set defaults
func setDefaults() {
	globalConfig := &GlobalProperties{
		Launcher: LauncherProperties{
			Http: HttpLauncherProperties{
				ServerGateway: constant.DefaultHttpServerGateway,
			},
			Kafka: KafkaLauncherProperties{
				Enabled:          false,
				BootstrapServers: constant.DefaultLauncherKafkaServers,
				Topic:            constant.DefaultLauncherKafkaTopic,
				Partitions:       constant.DefaultLauncherKafkaPartitions,
			},
		},
		Indicators: IndicatorsProperties{
			Netcard: constant.DefaultNetcard,
			Physical: PhysicalIndicatorProperties{
				Delay:     constant.DefaultIndicatorsDelay,
				RangePort: constant.DefaultNetIndicatorPortRange,
			},
			Virtual: VirtualIndicatorProperties{
				Delay: constant.DefaultIndicatorsDelay,
			},
			Redis: RedisIndicatorProperties{
				Delay: constant.DefaultIndicatorsDelay,
			},
			Kafka: KafkaIndicatorProperties{
				Delay: constant.DefaultIndicatorsDelay,
			},
			Zookeeper: ZookeeperIndicatorProperties{
				Delay: constant.DefaultIndicatorsDelay,
			},
			Etcd: EtcdIndicatorProperties{
				Delay: constant.DefaultIndicatorsDelay,
			},
			Emq: EmqIndicatorProperties{
				Delay: constant.DefaultIndicatorsDelay,
			},
		},
	}
	GlobalConfig = *globalConfig
}
