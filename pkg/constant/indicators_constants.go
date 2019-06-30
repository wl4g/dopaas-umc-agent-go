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
	// ------ Basic indicators -----

	// Default namespace.
	DefaultNamespace = "defaultCloud"

	// Default netcard name.
	DefaultNetcard = "eth0"

	// Default indicator delay.
	DefaultIndicatorsDelay time.Duration = 60000

	// Default network indicator ports range.
	DefaultNetIndicatorsNetPorts = "22,8080"

	// --------- Virtual -----------

	// Docker indicators constants.

	// Mesos indicators constants.

	// --------- Discovery ---------

	// Zk indicators constants.

	// Default watch zk servers.
	DefaultZkIndicatorsServers = "localhost:2181"

	// Default watch zk commands(e.g. mntr,conf).
	DefaultZkIndicatorsCommands = "mntr"

	// Default watch zk indicators metric filters.
	DefaultZkIndicatorsMetricFilters = "zk_avg_latency,zk_max_latency,zk_min_latency,zk_packets_received,zk_packets_sent,zk_num_alive_connections,zk_outstanding_requests,zk_server_state,zk_znode_count,zk_watch_count,zk_ephemerals_count,zk_approximate_data_size,zk_open_file_descriptor_count,zk_max_file_descriptor_count,zk_fsync_threshold_exceed_count"

	// Etcd indicators constants.

	// Consul indicators constants.

	// ----------- MQ -------------

	// Kafka indicators constants.
	// Default watch redis indicators metric filters.
	DefaultKafkaIndicatorsMetricFilters = "kafka_topic_partition_oldest_offset,kafka_brokers,kafka_partition_current_offset"

	// EMQ indicators constants.

	// RabbitMQ indicators constants.

	// RocketMQ indicators constants.

	// ----------- Cache -----------

	// Redis indicators constants.

	// Default watch redis servers.
	DefaultRedisIndicatorsServers = "localhost:6379,localhost:6380"

	// Default watch redis ports.
	DefaultRedisIndicatorsPassword = "redis"

	// Default watch redis indicators metric filters.
	DefaultRedisIndicatorsMetricFilters = "connected_clients,used_memory,used_memory_peak,used_memory_overhead,used_memory_dataset,used_memory_dataset_perc,aof_current_size,aof_buffer_length,total_connections_received,total_commands_processed,instantaneous_ops_per_sec,instantaneous_input_kbps,instantaneous_output_kbps,repl_backlog_size,repl_backlog_histlen,used_cpu_sys,used_cpu_user,used_cpu_sys_children,used_cpu_user_children,cluster_state,cluster_slots_assigned,cluster_slots_ok,cluster_slots_pfail,cluster_slots_fail,cluster_known_nodes,cluster_size,cluster_current_epoch,cluster_my_epoch"

	// Memcached indicators constants.

	// ------------ DB -------------

	// ElasticSearch indicators constants.

	// Mongodb indicators constants.

	// MySQL indicators constants.

	// PostgreSQL indicators constants.

	// OpenTSDB indicators constants.

	// Cassandra indicators constants.

)
