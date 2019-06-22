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

// -------------------------------
// Basic constants.
// -------------------------------

// Default profile path.
var DefaultConfigPath = "/etc/umc-agent.yml"

// -------------------------------
// Log constants.
// -------------------------------

var DefaultLogFilename = "./log/umc-agent.log"

// Default log level
var DefaultLogLevel = "INFO"

// Default log retention days.
var DefaultLogRetentionDays = 7

// Default log max backup numbers.
var DefaultLogMaxBackups = 30

// Default log max size(MB).
var DefaultLogMaxSize = 128

// -------------------------------
// Launcher constants.
// -------------------------------

// Default http server gateway
var DefaultHttpServerGateway = "http://localhost:14046/umc/basic"

// Default launcher kafka bootstrap servers.
var DefaultLauncherKafkaServers = "localhost:9092"

// Default launcher kafka topic.
var DefaultLauncherKafkaTopic = "__devops_umc_agent_metric_"

// Default launcher kafka partitions.
var DefaultLauncherKafkaPartitions int32 = 10

// -------------------------------
// Indicators basic constants.
// -------------------------------

// Default namespace.
var DefaultNamespace = "defaultCloud"

// Default netcard name.
var DefaultNetcard = "eth0"

// Default indicator delay.
var DefaultIndicatorsDelay time.Duration = 60000

// Default network indicator ports range.
var DefaultNetIndicatorsNetPorts = "22,8080"

// -------------------------------
// Redis indicators constants.
// -------------------------------

// Default watch redis ports.
var DefaultRedisIndicatorsPorts = "6379,6380"

// Default watch redis ports.
var DefaultRedisIndicatorsPassword = "redis"

// Default watch redis indicators properties.
var DefaultRedisIndicatorsProperties = "connected_clients,used_memory,used_memory_peak,used_memory_overhead,used_memory_dataset,used_memory_dataset_perc,aof_current_size,aof_buffer_length,total_connections_received,total_commands_processed,instantaneous_ops_per_sec,instantaneous_input_kbps,instantaneous_output_kbps,repl_backlog_size,repl_backlog_histlen,used_cpu_sys,used_cpu_user,used_cpu_sys_children,used_cpu_user_children,cluster_state,cluster_slots_assigned,cluster_slots_ok,cluster_slots_pfail,cluster_slots_fail,cluster_known_nodes,cluster_size,cluster_current_epoch,cluster_my_epoch"

// -------------------------------
// Zk indicators constants.
// -------------------------------

// Default watch zk servers.
var DefaultZkIndicatorsServers = "localhost:2181"

// Default watch zk commands(e.g. mntr,conf).
var DefaultZkIndicatorsCommands = "mntr"

// Default watch redis indicators properties.
var DefaultZkIndicatorsProperties = "zk_avg_latency,zk_max_latency,zk_min_latency,zk_packets_received,zk_packets_sent,zk_num_alive_connections,zk_outstanding_requests,zk_server_state,zk_znode_count,zk_watch_count,zk_ephemerals_count,zk_approximate_data_size,zk_open_file_descriptor_count,zk_max_file_descriptor_count,zk_fsync_threshold_exceed_count"
