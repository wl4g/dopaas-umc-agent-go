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

// Initialize global properties.
func InitGlobalProperties(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.MainLogger.Info("yamlFile.Get err - ", zap.Error(err))
	}
	err = yaml.Unmarshal(yamlFile, &GlobalConfig)
	if err != nil {
		log.MainLogger.Info("Unmarshal - ", zap.Error(err))
	}

	// Set Default
	if GlobalConfig.ServerGateway == "" {
		GlobalConfig.ServerGateway = constant.CONF_DEFAULT_SERVER_URI
	}
	if GlobalConfig.PhysicalPropertiesObj.Net == "" {
		GlobalConfig.PhysicalPropertiesObj.Net = constant.CONF_DEFAULT_NETCARD
	}
	if GlobalConfig.PhysicalPropertiesObj.Delay == 0 {
		GlobalConfig.PhysicalPropertiesObj.Delay = constant.CONF_DEFAULT_DELAY
	}
	if GlobalConfig.KafkaProducerPropertiesObj.bootstrapServers == "" {
		GlobalConfig.KafkaProducerPropertiesObj.bootstrapServers = constant.CONF_DEFAULT_KAFKA_URL
	}
	if GlobalConfig.KafkaProducerPropertiesObj.Topic == "" {
		GlobalConfig.KafkaProducerPropertiesObj.Topic = constant.CONF_DEFAULT_KAFKA_TOPIC
	}
	if GlobalConfig.Provider == "" {
		GlobalConfig.Provider = constant.CONF_DEFAULT_KAFKA_TOPIC
	}

}
