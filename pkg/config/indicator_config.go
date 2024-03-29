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
	"time"
)

const (
	// Used for metric filtering checks.
	// See: ./pkg/indicators/metric_builder.go#NewMetric()
	MetricExcludeFieldName = "MetricExclude"
)

// ----------------------
// Indicator properties.
// ----------------------

type IndicatorProperties struct {
	Namespace     string                           `yaml:"namespace"`
	Netcard       string                           `yaml:"netcard"`
	Physical      PhysicalIndicatorProperties      `yaml:"host"`
	Docker        DockerIndicatorProperties        `yaml:"docker"`
	Mesos         MesosIndicatorProperties         `yaml:"mesos"`
	Zookeeper     ZookeeperIndicatorProperties     `yaml:"zookeeper"`
	Etcd          EtcdIndicatorProperties          `yaml:"etcd"`
	Consul        ConsulIndicatorProperties        `yaml:"consul"`
	Kafka         KafkaIndicatorProperties         `yaml:"kafka"`
	Emq           EmqIndicatorProperties           `yaml:"emq"`
	RabbitMQ      RabbitMQIndicatorProperties      `yaml:"rabbitmq"`
	RocketMQ      RocketMQIndicatorProperties      `yaml:"rocketmq"`
	Redis         RedisIndicatorProperties         `yaml:"redis"`
	Memcached     MemcachedIndicatorProperties     `yaml:"memcached"`
	ElasticSearch ElasticSearchIndicatorProperties `yaml:"elasticsearch"`
	Mongodb       MongodbIndicatorProperties       `yaml:"mongodb"`
	MySQL         MySQLIndicatorProperties         `yaml:"mysql"`
	PostgreSQL    PostgreSQLIndicatorProperties    `yaml:"postgresql"`
	OpenTSDB      OpenTSDBIndicatorProperties      `yaml:"opentsdb"`
	Cassandra     CassandraIndicatorProperties     `yaml:"cassandra"`
}

// ------- Infrastructure ------

// Indicator host properties.
type PhysicalIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	NetPorts      string        `yaml:"range-port"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// ---------- Virtual ----------

// Indicator docker properties.
type DockerIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator mesos properties.
type MesosIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// --------- Discovery --------

// Indicator zookeeper properties.
type ZookeeperIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	Servers       string        `yaml:"servers"`
	Command       string        `yaml:"command"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator etcd properties.
type EtcdIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator consul properties.
type ConsulIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// ----------- MQ -------------

// Indicator kafka properties.
type KafkaIndicatorProperties struct {
	Enabled                  bool          `yaml:"enabled"`
	Delay                    time.Duration `yaml:"delay"`
	Servers                  string        `yaml:"servers"`
	UseSASL                  bool          `yaml:"useSASL"`
	UseSASLHandshake         bool          `yaml:"useSASLHandshake"`
	SaslUsername             string        `yaml:"saslUsername"`
	SaslPassword             string        `yaml:"saslPassword"`
	UseTLS                   bool          `yaml:"useTLS"`
	TlsCAFile                string        `yaml:"tlsCAFile"`
	TlsCertFile              string        `yaml:"tlsCertFile"`
	TlsKeyFile               string        `yaml:"tlsKeyFile"`
	TlsInsecureSkipTLSVerify bool          `yaml:"tlsInsecureSkipTLSVerify"`
	KafkaVersion             string        `yaml:"kafkaVersion"`
	UseZooKeeperLag          bool          `yaml:"useZooKeeperLag"`
	UriZookeeper             []string      `yaml:"uriZookeeper"`
	Labels                   string        `yaml:"labels"`
	MetadataRefreshInterval  string        `yaml:"metadataRefreshInterval"`
	MetricExclude            string        `yaml:"metric-exclude"`
}

// Indicator emq properties.
type EmqIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator rabbitMQ properties.
type RabbitMQIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator rocketMQ properties.
type RocketMQIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// ----------- Cache -----------

// Indicator redis properties.
type RedisIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	Servers       string        `yaml:"servers"`
	Password      string        `yaml:"password"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator memcached properties.
type MemcachedIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// ------------ DB -------------

// Indicator elastic-search properties.
type ElasticSearchIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator Mongodb properties.
type MongodbIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator MySQL properties.
type MySQLIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator PostgreSQL properties.
type PostgreSQLIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator OpenTSDB properties.
type OpenTSDBIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}

// Indicator Cassandra properties.
type CassandraIndicatorProperties struct {
	Enabled       bool          `yaml:"enabled"`
	Delay         time.Duration `yaml:"delay"`
	MetricExclude string        `yaml:"metric-exclude"`
}
