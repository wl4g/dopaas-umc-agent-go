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
	"github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"umc-agent/pkg/common"
	"umc-agent/pkg/constant"
)

const (
	// Used for metric filtering checks.
	// See: ./pkg/indicators/metric_builder.go#NewMetric()
	IndicatorsFiledName = "Indicator"
)

// ---------------------
// Global properties
// ---------------------
type GlobalProperties struct {
	Logging   LoggingProperties   `yaml:"logging"`
	Transport TransportProperties `yaml:"transport"`
	Indicator IndicatorProperties `yaml:"indicator"`
}

var (
	// Global config.
	GlobalConfig GlobalProperties

	// Global config buffer.
	_globalConfigBuffer []byte

	// Local hardware addr ID.
	LocalHardwareAddrId = ""
)

// Init global config properties.
func InitGlobalConfig(path string) {
	// Create default config.
	GlobalConfig = *createDefault()

	conf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Read config '%s' error! %s", path, err)
		panic(err)
		return
	}

	err = yaml.Unmarshal(conf, &GlobalConfig)
	if err != nil {
		fmt.Printf("Unmarshal config '%s' error! %s", path, err)
		panic(err)
		return
	}

	// Post properties.
	afterPropertiesSet(&GlobalConfig)
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
		Transport: TransportProperties{
			Http: HttpTransportProperties{
				ServerGateway: constant.DefaultHttpServerGateway,
			},
			Kafka: KafkaTransportProperties{
				Enabled:      false,
				Servers:      constant.DefaultTransportKafkaServers,
				MetricTopic:  constant.DefaultTransportKafkaMetricTopic,
				ReceiveTopic: constant.DefaultTransportKafkaReceiveTopic,
				Ack:          constant.DefaultTransportKafkaAck,
				Timeout:      constant.DefaultTransportKafkaTimeout,
			},
		},
		Indicator: IndicatorProperties{
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
				Enabled:       false,
				Delay:         constant.DefaultIndicatorsDelay,
				Servers:       constant.DefaultZkIndicatorsServers,
				Command:       constant.DefaultZkIndicatorsCommands,
				MetricExcludeRegex: constant.DefaultZkIndicatorsMetricExcludeRegex,
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
				Enabled:       false,
				Delay:         constant.DefaultIndicatorsDelay,
				Servers:      constant.DefaultKafkaIndicatorsServers,
				MetricExcludeRegex: constant.DefaultKafkaIndicatorsMetricExcludeRegex,
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
				Properties: constant.DefaultRedisIndicatorsMetricExcludeRegex,
				MetricExcludeRegex: constant.DefaultRedisIndicatorsMetricExcludeRegex,
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

// MetricExcludeRegex settings after initialization
func afterPropertiesSet(globalConfig *GlobalProperties) {
	// Environment variable priority.
	var netcard = os.Getenv("indicators.netcard")
	if !common.IsEmpty(netcard) {
		globalConfig.Indicator.Netcard = netcard
	}

	// Local hardware addr.
	LocalHardwareAddrId = common.GetHardwareAddr(globalConfig.Indicator.Netcard)
	if LocalHardwareAddrId == "" || len(LocalHardwareAddrId) <= 0 {
		panic("Failed to find network hardware info, please check the net config!")
	}

	// To config json buffer(see: #GetConfig()).
	var buffer, _ = jsoniter.MarshalToString(GlobalConfig)
	_globalConfigBuffer = []byte(buffer)
}

// Get config value.
func GetConfig(path ...interface{}) jsoniter.Any {
	var value = jsoniter.Get(_globalConfigBuffer, path...)
	return value
}
