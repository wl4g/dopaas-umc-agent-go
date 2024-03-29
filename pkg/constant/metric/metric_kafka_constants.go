/**
 * Copyright 2017 ~ 2025 the original author or authors.
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
package metric

const (
	//
	// Tags.
	//

	DefaultKafkaTopicTag = "topic"

	DefaultKafkaGroupIdTag = "groupid"

	DefaultKafkaPartitionTag = "partition"

	//
	// Metrics.
	//

	KafkaBrokersMetric = "kafka_brokers"

	KafkaTopicPartitionCurrentOffsetMetric = "kafka_topic_partition_current_offset"

	KafkaTopicPartitionInSyncReplicaMetric = "kafka_topic_partition_in_sync_replica"

	KafkaTopicPartitionLeaderMetric = "kafka_topic_partition_leader"

	KafkaTopicPartitionLeaderIsPreferredMetric = "kafka_topic_partition_leader_is_preferred"

	KafkaTopicPartitionOldestOffsetMetric = "kafka_topic_partition_oldest_offset"

	KafkaTopicPartitionReplicasMetric = "kafka_topic_partition_replicas"

	KafkaTopicPartitionUnderReplicatedPartitionMetric = "kafka_topic_partition_under_replicated_partition"

	KafkaTopicPartitionsMetric = "kafka_topic_partitions"

	KafkaConsumerGroupMembers = "kafka_consumergroup_members"

	KafkaConsumerGroupCurrentOffset = "kafka_consumergroup_current_offset"

	KafkaConsumerGroupLag = "kafka_consumergroup_lag"

	KafkaConsumerGroupCurrentOffsetSum = "kafka_consumergroup_current_offset_sum"

	KafkaConsumerGroupLagSum = "kafka_consumergroup_lag_sum"
)
