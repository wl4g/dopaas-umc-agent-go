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

// Global configuration.
var GlobalConfig GlobalProperties

// Local hardware addr ID.
var LocalHardwareAddrId = ""

// Init global config properties.
func InitGlobalConfig(path string) {
	// Create default config.
	GlobalConfig := createDefault()

	conf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Read config '%s' error! %s", path, err)
		panic(err)
		return
	}

	err = yaml.Unmarshal(conf, GlobalConfig)
	if err != nil {
		fmt.Printf("Unmarshal config '%s' error! %s", path, err)
		panic(err)
		return
	}

	// Post properties.
	afterPropertiesSet()
}

// Create default config.
func createDefault() *GlobalProperties {
	globalConfig := &GlobalProperties{
		Logging: LoggingProperties{
			LogItems: map[string]LogItemProperties{
				constant.DefaultLogMain: {
					FileName: constant.DefaultLogDir + constant.DefaultLogMain + ".log",
					Level:    constant.DefaultLogLevel,
					Policy: PolicyProperties{
						RetentionDays: constant.DefaultLogRetentionDays,
						MaxBackups:    constant.DefaultLogMaxBackups,
						MaxSize:       constant.DefaultLogMaxSize,
					},
				},
				constant.DefaultLogReceive: {
					FileName: constant.DefaultLogDir + constant.DefaultLogReceive + ".log",
					Level:    constant.DefaultLogLevel,
					Policy: PolicyProperties{
						RetentionDays: constant.DefaultLogRetentionDays,
						MaxBackups:    constant.DefaultLogMaxBackups,
						MaxSize:       constant.DefaultLogMaxSize,
					},
				},
			},
		},
		Launcher: LauncherProperties{
			Http: HttpLauncherProperties{
				ServerGateway: constant.DefaultHttpServerGateway,
			},
			Kafka: KafkaLauncherProperties{
				Enabled:          false,
				BootstrapServers: constant.DefaultLauncherKafkaServers,
				MetricTopic:      constant.DefaultLauncherKafkaMetricTopic,
				ReceiveTopic:     constant.DefaultLauncherKafkaReceiveTopic,
				Ack:              constant.DefaultLauncherKafkaAck,
				Timeout:          constant.DefaultLauncherKafkaTimeout,
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
			Docker: DockerIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Mesos: MesosIndicatorProperties{
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
			Consul: ConsulIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Kafka: KafkaIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Emq: EmqIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			RabbitMQ: RabbitMQIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			RocketMQ: RocketMQIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Redis: RedisIndicatorProperties{
				Enabled:    false,
				Delay:      constant.DefaultIndicatorsDelay,
				Servers:    constant.DefaultRedisIndicatorsServers,
				Password:   constant.DefaultRedisIndicatorsPassword,
				Properties: constant.DefaultRedisIndicatorsProperties,
			},
			Memcached: MemcachedIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			ElasticSearch: ElasticSearchIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Mongodb: MongodbIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			MySQL: MySQLIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			PostgreSQL: PostgreSQLIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			OpenTSDB: OpenTSDBIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
			Cassandra: CassandraIndicatorProperties{
				Enabled: false,
				Delay:   constant.DefaultIndicatorsDelay,
			},
		},
	}

	return globalConfig
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
