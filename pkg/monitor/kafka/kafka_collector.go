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
package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Shopify/sarama"
	"strconv"
	"sync"
	"umc-agent/pkg/config"
	"umc-agent/pkg/constant/metric"
	"umc-agent/pkg/monitor"
	"umc-agent/pkg/transport"

	//"github.com/krallistic/kazoo-go"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"time"
	"umc-agent/pkg/logger"
)

const (
	clientID = "kafka_exporter"
)

func IndicatorRunner() {
	if !config.GlobalConfig.Indicators.Kafka.Enabled {
		logger.Main.Warn("No enabled kafka metrics runner!")
		return
	}
	logger.Main.Info("Starting kafka indicators runner ...")

	for true {
		// New metric aggregator
		kafkaAggregator := monitor.NewMetricAggregate("Kafka")

		// Do collect brokers metrics.
		doKafkaMetricCollect(kafkaAggregator)

		// Send to servers.
		transport.DoSendSubmit(metric.KafkaMeta, &kafkaAggregator)

		// Sleep.
		time.Sleep(config.GlobalConfig.Indicators.Kafka.Delay)
	}
}

func newConnectClient() sarama.Client {
	//TODO read config
	opts := kafkaOpts{}
	opts.uri = append(opts.uri, "localhost:9092")
	opts.useSASL = false
	opts.useSASLHandshake = true
	opts.useTLS = false
	opts.tlsInsecureSkipTLSVerify = false
	opts.kafkaVersion = sarama.V1_0_0_0.String()
	opts.useZooKeeperLag = false
	//opts.uriZookeeper = append(opts.uriZookeeper, "localhost:2181")
	opts.metadataRefreshInterval = "30s"

	//build config start
	config := sarama.NewConfig()
	config.ClientID = clientID
	kafkaVersion, err := sarama.ParseKafkaVersion(opts.kafkaVersion)
	if err != nil {
		panic(err)
	}
	config.Version = kafkaVersion
	if opts.useSASL {
		config.Net.SASL.Enable = true
		config.Net.SASL.Handshake = opts.useSASLHandshake

		if opts.saslUsername != "" {
			config.Net.SASL.User = opts.saslUsername
		}

		if opts.saslPassword != "" {
			config.Net.SASL.Password = opts.saslPassword
		}
	}
	if opts.useTLS {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			RootCAs:            x509.NewCertPool(),
			InsecureSkipVerify: opts.tlsInsecureSkipTLSVerify,
		}
		if opts.tlsCAFile != "" {
			if ca, err := ioutil.ReadFile(opts.tlsCAFile); err == nil {
				config.Net.TLS.Config.RootCAs.AppendCertsFromPEM(ca)
			} else {
				logger.Main.Error("", zap.Error(err))
			}
		}
		canReadCertAndKey, err := CanReadCertAndKey(opts.tlsCertFile, opts.tlsKeyFile)
		if err != nil {
			logger.Main.Error("", zap.Error(err))
		}
		if canReadCertAndKey {
			cert, err := tls.LoadX509KeyPair(opts.tlsCertFile, opts.tlsKeyFile)
			if err == nil {
				config.Net.TLS.Config.Certificates = []tls.Certificate{cert}
			} else {
				logger.Main.Error("", zap.Error(err))
			}
		}
	}
	/*if opts.useZooKeeperLag {
		zookeeperClient, err := kazoo.NewKazoo(opts.uriZookeeper, nil)
	}*/
	interval, err := time.ParseDuration(opts.metadataRefreshInterval)
	if err != nil {
		logger.Main.Error("Cannot parse metadata refresh interval", zap.Error(err))
		panic(err)
	}
	config.Metadata.RefreshFrequency = interval
	client, err := sarama.NewClient(opts.uri, config)
	return client
}

func doKafkaMetricCollect(kafkaAggregator *monitor.MetricAggregate) {
	// Connect servers.
	client := newConnectClient()

	var mu sync.Mutex
	now := time.Now().UnixNano() / 1e6
	kafkaAggregator.Timestamp = now

	//brokers count
	kafkaAggregator.NewMetric(metric.KafkaBrokersMetric, float64(len(client.Brokers())))

	//Refresh Metadata
	client.RefreshMetadata()

	//get topics
	topics, err := client.Topics()

	//wg
	var wg = sync.WaitGroup{}
	offset := make(map[string]map[int32]int64)
	//getTopicMetrics
	getTopicMetrics := func(topic string) {
		defer wg.Done()

		partitions, err := client.Partitions(topic)
		if err != nil {
			logger.Main.Error("Cannot get partitions of topic ", zap.String("topic", topic), zap.Error(err))
			panic(err)
		}

		//Partitions count
		kafkaAggregator.NewMetric(metric.KafkaTopicPartitionsMetric, float64(len(partitions))).ATag(metric.Topic, topic)

		mu.Lock()
		offset[topic] = make(map[int32]int64, len(partitions))
		mu.Unlock()
		for _, partition := range partitions {
			broker, err := client.Leader(topic, partition)
			if err != nil {
				logger.Main.Error("Cannot get leader of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionLeaderMetric,
					float64(broker.ID())).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			}
			currentOffset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				logger.Main.Error("Cannot get current offset of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				mu.Lock()
				offset[topic][partition] = currentOffset
				mu.Unlock()
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionCurrentOffsetMetric,
					float64(currentOffset)).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			}

			oldestOffset, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
			if err != nil {
				logger.Main.Error("Cannot get oldest offset of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionOldestOffsetMetric,
					float64(oldestOffset)).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			}

			replicas, err := client.Replicas(topic, partition)
			if err != nil {
				logger.Main.Error("Cannot get replicas of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionReplicasMetric,
					float64(len(replicas))).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			}

			inSyncReplicas, err := client.InSyncReplicas(topic, partition)
			if err != nil {
				logger.Main.Error("Cannot get in-sync replicas of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionInSyncReplicaMetric,
					float64(len(inSyncReplicas))).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			}

			if broker != nil && replicas != nil && len(replicas) > 0 && broker.ID() == replicas[0] {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionLeaderIsPreferredMetric,
					float64(1)).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionLeaderIsPreferredMetric,
					float64(0)).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			}

			if replicas != nil && inSyncReplicas != nil && len(inSyncReplicas) < len(replicas) {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionUnderReplicatedPartitionMetric,
					float64(1)).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionUnderReplicatedPartitionMetric,
					float64(0)).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
			}
		}
	}

	for _, topic := range topics {
		wg.Add(1)
		go getTopicMetrics(topic)
	}
	wg.Wait()

	//////////////////////////////////
	getConsumerGroupMetrics := func(broker *sarama.Broker) {
		defer wg.Done()
		if err := broker.Open(client.Config()); err != nil && err != sarama.ErrAlreadyConnected {
			logger.Main.Error("Cannot connect to broker",
				zap.Int32("brokerID", broker.ID()),
				zap.Error(err))
			panic(err)
		}
		defer broker.Close()
		//groups, err := broker.ListGroups(&sarama.ListGroupsRequest{})
		if err != nil {
			logger.Main.Error("Cannot get consumer group",
				zap.Error(err))
			panic(err)
		}
		groupIds := make([]string, 0)
		describeGroups, err := broker.DescribeGroups(&sarama.DescribeGroupsRequest{Groups: groupIds})
		if err != nil {
			logger.Main.Error("Cannot get describe groups",
				zap.Error(err))
			panic(err)
		}
		for _, group := range describeGroups.Groups {
			offsetFetchRequest := sarama.OffsetFetchRequest{ConsumerGroup: group.GroupId, Version: 1}
			for topic, partitions := range offset {
				for partition := range partitions {
					offsetFetchRequest.AddPartition(topic, partition)
				}
			}
			kafkaAggregator.NewMetric(metric.KafkaConsumerGroupMembers, float64(len(group.Members))).ATag(metric.GroupId, group.GroupId)

			if offsetFetchResponse, err := broker.FetchOffset(&offsetFetchRequest); err != nil {
				logger.Main.Error("Cannot get offset of group",
					zap.String("groudId", group.GroupId),
					zap.Error(err))
			} else {
				for topic, partitions := range offsetFetchResponse.Blocks {
					// If the topic is not consumed by that consumer group, skip it
					topicConsumed := false
					for _, offsetFetchResponseBlock := range partitions {
						// Kafka will return -1 if there is no offset associated with a topic-partition under that consumer group
						if offsetFetchResponseBlock.Offset != -1 {
							topicConsumed = true
							break
						}
					}
					if topicConsumed {
						var currentOffsetSum int64
						var lagSum int64
						for partition, offsetFetchResponseBlock := range partitions {
							err := offsetFetchResponseBlock.Err
							if err != sarama.ErrNoError {
								logger.Main.Error("Error for  partition",
									zap.Int32("partition", partition),
									zap.Error(err))
								continue
							}
							currentOffset := offsetFetchResponseBlock.Offset
							currentOffsetSum += currentOffset

							kafkaAggregator.NewMetric(metric.KafkaConsumerGroupCurrentOffset, float64(currentOffset)).ATag(metric.GroupId, group.GroupId).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))

							mu.Lock()
							if offset, ok := offset[topic][partition]; ok {
								// If the topic is consumed by that consumer group, but no offset associated with the partition
								// forcing lag to -1 to be able to alert on that
								var lag int64
								if offsetFetchResponseBlock.Offset == -1 {
									lag = -1
								} else {
									lag = offset - offsetFetchResponseBlock.Offset
									lagSum += lag
								}
								kafkaAggregator.NewMetric(metric.KafkaConsumerGroupLag, float64(lag)).ATag(metric.GroupId, group.GroupId).ATag(metric.Topic, topic).ATag(metric.Partition, strconv.FormatInt(int64(partition), 10))
							} else {
								logger.Main.Error("No offset of topic",
									zap.Int32("partition", partition),
									zap.String("topic", topic),
									zap.Error(err))
							}
							mu.Unlock()
						}

						kafkaAggregator.NewMetric(metric.KafkaConsumerGroupCurrentOffsetSum, float64(currentOffsetSum)).ATag(metric.GroupId, group.GroupId).ATag(metric.Topic, topic)

						kafkaAggregator.NewMetric(metric.KafkaConsumerGroupLagSum, float64(lagSum)).ATag(metric.GroupId, group.GroupId).ATag(metric.Topic, topic)
					}
				}
			}
		}
	}

	if len(client.Brokers()) > 0 {
		for _, broker := range client.Brokers() {
			wg.Add(1)
			go getConsumerGroupMetrics(broker)
		}
		wg.Wait()
	} else {
		logger.Main.Error("No valid broker, cannot get consumer group metrics")
	}

}

// CanReadCertAndKey returns true if the certificate and key files already exists,
// otherwise returns false. If lost one of cert and key, returns error.
func CanReadCertAndKey(certPath, keyPath string) (bool, error) {
	certReadable := canReadFile(certPath)
	keyReadable := canReadFile(keyPath)

	if certReadable == false && keyReadable == false {
		return false, nil
	}

	if certReadable == false {
		logger.Main.Error("error reading , certificate and key must be supplied as a pair",
			zap.String("certPath", certPath))
	}

	if keyReadable == false {
		logger.Main.Error("error reading , certificate and key must be supplied as a pair",
			zap.String("keyPath", keyPath))
	}
	return true, nil
}

// If the file represented by path exists and
// readable, returns true otherwise returns false.
func canReadFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}
