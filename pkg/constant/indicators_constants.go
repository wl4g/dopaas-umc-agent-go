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
package constant

import (
	"time"
)

const (
	// =============================
	// ------ Basic indicators -----
	// =============================

	// Default namespace.
	DefaultNamespace = "defaultCloud"

	// Default netcard name.
	DefaultNetcard = "eth0"

	// Default indicator delay.
	DefaultIndicatorsDelay time.Duration = 60000

	// Default network indicator ports range.
	DefaultNetIndicatorsNetPorts = "22,8080"

	//Tag disk device
	TagDiskDevice = "device"

	//Tag net port
	TagNetPort = "port"

	// =============================
	// --------- Virtual -----------
	// =============================

	// Docker indicators constants.
	TagDockerName = "name"

	// Mesos indicators constants.

	// =============================
	// --------- Discovery ---------
	// =============================

	// Zk indicators constants.

	// Default watch zk servers.
	DefaultZkIndicatorsServers = "localhost:2181"

	// Default watch zk commands(e.g. mntr,conf).
	DefaultZkIndicatorsCommands = "mntr"

	// Etcd indicators constants.

	// Consul indicators constants.

	// =============================
	// ----------- MQ --------------
	// =============================

	// Kafka indicators constants.

	// Default indicators kafka broker servers.
	DefaultKafkaIndicatorsServers = "localhost:9092"

	// EMQ indicators constants.

	// RabbitMQ indicators constants.

	// RocketMQ indicators constants.

	// =============================
	// ----------- Cache -----------
	// =============================

	// Redis indicators constants.

	// Default watch redis servers.
	DefaultRedisIndicatorsServers = "localhost:6379,localhost:6380"

	// Default watch redis ports.
	DefaultRedisIndicatorsPassword = "redis"

	// Memcached indicators constants.

	// =============================
	// ------------ DB -------------
	// =============================

	// ElasticSearch indicators constants.

	// Mongodb indicators constants.

	// MySQL indicators constants.

	// PostgreSQL indicators constants.

	// OpenTSDB indicators constants.

	// Cassandra indicators constants.

)
