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

// ----------------------
// Indicators properties.
// ----------------------

type IndicatorsProperties struct {
	Namespace     string                           `yaml:"namespace"`
	Netcard       string                           `yaml:"netcard"`
	Physical      PhysicalIndicatorProperties      `yaml:"physical"`
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

// Indicators physical properties.
type PhysicalIndicatorProperties struct {
	Enabled  bool          `yaml:"enabled"`
	Delay    time.Duration `yaml:"delay"`
	NetPorts string        `yaml:"range-port"`
}

// ---------- Virtual ----------

// Indicators docker properties.
type DockerIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators mesos properties.
type MesosIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// --------- Discovery --------

// Indicators zookeeper properties.
type ZookeeperIndicatorProperties struct {
	Enabled    bool          `yaml:"enabled"`
	Delay      time.Duration `yaml:"delay"`
	Servers    string        `yaml:"servers"`
	Command    string        `yaml:"command"`
	Properties string        `yaml:"properties"`
}

// Indicators etcd properties.
type EtcdIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators consul properties.
type ConsulIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// ----------- MQ -------------

// Indicators kafka properties.
type KafkaIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators emq properties.
type EmqIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators rabbitMQ properties.
type RabbitMQIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators rocketMQ properties.
type RocketMQIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// ----------- Cache -----------

// Indicators redis properties.
type RedisIndicatorProperties struct {
	Enabled    bool          `yaml:"enabled"`
	Delay      time.Duration `yaml:"delay"`
	Servers    string        `yaml:"servers"`
	Password   string        `yaml:"password"`
	Properties string        `yaml:"properties"`
}

// Indicators memcached properties.
type MemcachedIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// ------------ DB -------------

// Indicators elastic-search properties.
type ElasticSearchIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators Mongodb properties.
type MongodbIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators MySQL properties.
type MySQLIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators PostgreSQL properties.
type PostgreSQLIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators OpenTSDB properties.
type OpenTSDBIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}

// Indicators Cassandra properties.
type CassandraIndicatorProperties struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
}
