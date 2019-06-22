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
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/constant"
	"umc-agent/pkg/monitor/share"
)

// ---------------------
// Global properties
// ---------------------
type GlobalProperties struct {
	Logging    LoggingProperties    `yaml:"logging"`
	Launcher   LauncherProperties   `yaml:"launcher"`
	Indicators IndicatorsProperties `yaml:"indicators"`
}

// ---------------------
// Logging properties
// ---------------------
type LoggingProperties struct {
	FileName string           `yaml:"file"`
	Level    string           `yaml:"level"`
	Policy   PolicyProperties `yaml:"policy"`
}

// Logging archive policy
type PolicyProperties struct {
	RetentionDays int `yaml:"retention-days"`
	MaxBackups    int `yaml:"max-backups"`
	MaxSize       int `yaml:"max-size"`
}

// ---------------------
// Launcher properties.
// ---------------------
type LauncherProperties struct {
	Http  HttpLauncherProperties  `yaml:"http"`
	Kafka KafkaLauncherProperties `yaml:"kafka"`
}

type HttpLauncherProperties struct {
	ServerGateway string `yaml:"server-gateway"`
}

type KafkaLauncherProperties struct {
	Enabled          bool   `yaml:"enabled"`
	BootstrapServers string `yaml:"bootstrap.servers"`
	Topic            string `yaml:"topic"`
	Partitions       int32  `yaml:"partitions"`
}

// ----------------------
// Indicators properties.
// ----------------------

type IndicatorsProperties struct {
	Namespace string                       `yaml:"namespace"`
	Netcard   string                       `yaml:"netcard"`
	Physical  PhysicalIndicatorProperties  `yaml:"physical"`
	Virtual   VirtualIndicatorProperties   `yaml:"virtual"`
	Redis     RedisIndicatorProperties     `yaml:"redis"`
	Zookeeper ZookeeperIndicatorProperties `yaml:"zookeeper"`
	Kafka     KafkaIndicatorProperties     `yaml:"kafka"`
	Etcd      EtcdIndicatorProperties      `yaml:"etcd"`
	Emq       EmqIndicatorProperties       `yaml:"emq"`
	Consul    ConsulIndicatorProperties    `yaml:"consul"`
}

// Indicators physical properties.
type PhysicalIndicatorProperties struct {
	Enabled  bool          `yaml:"enabled"`
	Delay    time.Duration `yaml:"delay"`
	NetPorts string        `yaml:"range-port"`
}

// Indicators virtual properties.
type VirtualIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators redis properties.
type RedisIndicatorProperties struct {
	Enabled    bool          `yaml:"enabled"`
	Delay      time.Duration `yaml:"delay"`
	Ports      string        `yaml:"ports"`
	Password   string        `yaml:"password"`
	Properties string        `yaml:"properties"`
}

// Indicators zookeeper properties.
type ZookeeperIndicatorProperties struct {
	Enabled    bool          `yaml:"enabled"`
	Delay      time.Duration `yaml:"delay"`
	Servers    string        `yaml:"servers"`
	Command    string        `yaml:"command"`
	Properties string        `yaml:"properties"`
}

// Indicators kafka properties.
type KafkaIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators etcd properties.
type EtcdIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators emq properties.
type EmqIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators consul properties.
type ConsulIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Global configuration.
var GlobalConfig GlobalProperties

// Local hardware addr ID.
var LocalHardwareAddrId = ""

// Init global config properties.
func InitGlobalConfig(path string) {
	// Set defaults
	setDefaults()

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Read config '%s' error! %s", path, err)
	}

	err = yaml.Unmarshal(yamlFile, &GlobalConfig)
	if err != nil {
		fmt.Printf("Unmarshal config '%s' error! %s", path, err)
	}

	// Post properties.
	afterPropertiesSet()
}

// Set defaults
func setDefaults() {
	globalConfig := &GlobalProperties{
		Logging: LoggingProperties{
			FileName: constant.DefaultLogFilename,
			Level:    constant.DefaultLogLevel,
			Policy: PolicyProperties{
				RetentionDays: constant.DefaultLogRetentionDays,
				MaxBackups:    constant.DefaultLogMaxBackups,
				MaxSize:       constant.DefaultLogMaxSize,
			},
		},
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
			Namespace: constant.DefaultNamespace,
			Netcard:   constant.DefaultNetcard,
			Physical: PhysicalIndicatorProperties{
				Enabled:  true,
				Delay:    constant.DefaultIndicatorsDelay,
				NetPorts: constant.DefaultNetIndicatorsNetPorts,
			},
			Virtual: VirtualIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Redis: RedisIndicatorProperties{
				Enabled:    false,
				Delay:      constant.DefaultIndicatorsDelay,
				Ports:      constant.DefaultRedisIndicatorsPorts,
				Password:   constant.DefaultRedisIndicatorsPassword,
				Properties: constant.DefaultRedisIndicatorsProperties,
			},
			Kafka: KafkaIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Zookeeper: ZookeeperIndicatorProperties{
				Enabled:    false,
				Delay:      constant.DefaultIndicatorsDelay,
				Servers:    constant.DefaultZkIndicatorsServers,
				Command:    constant.DefaultZkIndicatorsCommands,
				Properties: constant.DefaultZkIndicatorsProperties,
			},
			Etcd: EtcdIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Emq: EmqIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Consul: ConsulIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
		},
	}
	GlobalConfig = *globalConfig
}

// Properties settings after initialization
func afterPropertiesSet() {
	// Environmental variable priority
	var netcard = os.Getenv("indicators.netcard")
	if !common.IsEmpty(netcard) {
		GlobalConfig.Indicators.Netcard = netcard
	}

	// Got local hardware addr
	LocalHardwareAddrId = common.GetHardwareAddr(GlobalConfig.Indicators.Netcard)
	if LocalHardwareAddrId == "" || len(LocalHardwareAddrId) <= 0 {
		panic("net found ip,Please check the net conf")
	}

}

// Create meta info
func CreateMeta(metaType string) share.MetaInfo {
	meta := share.MetaInfo{
		Id:        LocalHardwareAddrId,
		Type:      metaType,
		Namespace: GlobalConfig.Indicators.Namespace,
	}
	return meta
}
