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
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"umc-agent/pkg/config"
	"umc-agent/pkg/constant/metric"
	"umc-agent/pkg/indicators"
	"umc-agent/pkg/logger"
	"umc-agent/pkg/transport"
)

const (
	clientID = "kafka_exporter"
)

var (
	kafkaClient sarama.Client
)

// Original see: https://github.com/danielqsj/kafka_exporter
func IndicatorRunner() {
	if !config.GlobalConfig.Indicator.Kafka.Enabled {
		logger.Main.Debug("No enabled kafka metrics runner!")
		return
	}
	logger.Main.Info("Starting kafka indicators runner ...")

	for true {
		// New kafka metric aggregator
		aggregator := indicators.NewMetricAggregator("Kafka")

		// Do collect brokers metrics.
		handleKafkaMetricCollect(aggregator)

		// Send to servers.
		transport.SendMetrics(aggregator)

		// Sleep.
		time.Sleep(config.GlobalConfig.Indicator.Kafka.Delay)
	}
}

// New kafka connect client.
func newConnectClient() sarama.Client {
	opts := kafkaOpts{}
	uris := strings.Split(config.GlobalConfig.Indicator.Kafka.Servers, ",")
	for _, v := range uris {
		opts.uri = append(opts.uri, v)
	}
	opts.useSASL = config.GlobalConfig.Indicator.Kafka.UseSASL
	opts.useSASLHandshake = config.GlobalConfig.Indicator.Kafka.UseSASLHandshake
	opts.useTLS = config.GlobalConfig.Indicator.Kafka.UseTLS
	opts.tlsInsecureSkipTLSVerify = config.GlobalConfig.Indicator.Kafka.TlsInsecureSkipTLSVerify
	opts.kafkaVersion = sarama.V1_0_0_0.String()
	opts.useZooKeeperLag = config.GlobalConfig.Indicator.Kafka.UseZooKeeperLag
	opts.metadataRefreshInterval = config.GlobalConfig.Indicator.Kafka.MetadataRefreshInterval

	// Build configN start
	configN := sarama.NewConfig()
	configN.ClientID = clientID
	kafkaVersion, err := sarama.ParseKafkaVersion(opts.kafkaVersion)
	if err != nil {
		panic(err)
	}
	configN.Version = kafkaVersion
	if opts.useSASL {
		configN.Net.SASL.Enable = true
		configN.Net.SASL.Handshake = config.GlobalConfig.Indicator.Kafka.UseSASLHandshake
		if config.GlobalConfig.Indicator.Kafka.SaslUsername != "" {
			configN.Net.SASL.User = config.GlobalConfig.Indicator.Kafka.SaslUsername
		}
		if config.GlobalConfig.Indicator.Kafka.SaslPassword != "" {
			configN.Net.SASL.Password = config.GlobalConfig.Indicator.Kafka.SaslPassword
		}
	}
	if opts.useTLS {
		configN.Net.TLS.Enable = true
		configN.Net.TLS.Config = &tls.Config{
			RootCAs:            x509.NewCertPool(),
			InsecureSkipVerify: config.GlobalConfig.Indicator.Kafka.TlsInsecureSkipTLSVerify,
		}
		if config.GlobalConfig.Indicator.Kafka.TlsCAFile != "" {
			if ca, err := ioutil.ReadFile(config.GlobalConfig.Indicator.Kafka.TlsCAFile); err == nil {
				configN.Net.TLS.Config.RootCAs.AppendCertsFromPEM(ca)
			} else {
				logger.Main.Error("", zap.Error(err))
			}
		}
		canReadCertAndKey, err := canReadCertAndKey(config.GlobalConfig.Indicator.Kafka.TlsCertFile, config.GlobalConfig.Indicator.Kafka.TlsKeyFile)
		if err != nil {
			logger.Main.Error("", zap.Error(err))
		}
		if canReadCertAndKey {
			cert, err := tls.LoadX509KeyPair(config.GlobalConfig.Indicator.Kafka.TlsCertFile, config.GlobalConfig.Indicator.Kafka.TlsKeyFile)
			if err == nil {
				configN.Net.TLS.Config.Certificates = []tls.Certificate{cert}
			} else {
				logger.Main.Error("", zap.Error(err))
			}
		}
	}
	//if opts.useZooKeeperLag {
	//	zookeeperClient, err := kazoo.NewKazoo(opts.uriZookeeper, nil)
	//}
	interval, err := time.ParseDuration(opts.metadataRefreshInterval)
	if err != nil {
		logger.Main.Error("Cannot parse metadata refresh interval", zap.Error(err))
		panic(err)
	}
	configN.Metadata.RefreshFrequency = interval
	client, err := sarama.NewClient(opts.uri, configN)
	return client
}

// Do kafka metrics collect.
func handleKafkaMetricCollect(kafkaAggregator *indicators.MetricAggregator) {
	// Connect servers.
	client := newConnectClient()
	defer client.Close()

	var mu sync.Mutex
	// Metric for brokers count
	kafkaAggregator.NewMetric(metric.KafkaBrokersMetric, float64(len(client.Brokers())))

	// Refresh Metadata
	client.RefreshMetadata()

	topics, err := client.Topics()

	var wg = sync.WaitGroup{}
	offset := make(map[string]map[int32]int64)

	getTopicMetrics := func(topic string) {
		defer wg.Done()

		partitions, err := client.Partitions(topic)
		if err != nil {
			logger.Main.Error("Cannot get partitions of topic ", zap.String("topic", topic), zap.Error(err))
			panic(err)
		}

		//Partitions count
		kafkaAggregator.NewMetric(metric.KafkaTopicPartitionsMetric, float64(len(partitions))).ATag(metric.DefaultKafkaTopicTag, topic)

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
					float64(broker.ID())).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
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
					float64(currentOffset)).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
			}

			oldestOffset, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
			if err != nil {
				logger.Main.Error("Cannot get oldest offset of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionOldestOffsetMetric,
					float64(oldestOffset)).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
			}

			replicas, err := client.Replicas(topic, partition)
			if err != nil {
				logger.Main.Error("Cannot get replicas of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionReplicasMetric,
					float64(len(replicas))).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
			}

			inSyncReplicas, err := client.InSyncReplicas(topic, partition)
			if err != nil {
				logger.Main.Error("Cannot get in-sync replicas of topic",
					zap.Int32("partition", partition),
					zap.String("topic", topic),
					zap.Error(err))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionInSyncReplicaMetric,
					float64(len(inSyncReplicas))).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
			}

			if broker != nil && replicas != nil && len(replicas) > 0 && broker.ID() == replicas[0] {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionLeaderIsPreferredMetric,
					float64(1)).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionLeaderIsPreferredMetric,
					float64(0)).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
			}

			if replicas != nil && inSyncReplicas != nil && len(inSyncReplicas) < len(replicas) {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionUnderReplicatedPartitionMetric,
					float64(1)).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
			} else {
				kafkaAggregator.NewMetric(metric.KafkaTopicPartitionUnderReplicatedPartitionMetric,
					float64(0)).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
					strconv.FormatInt(int64(partition), 10))
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

			kafkaAggregator.NewMetric(metric.KafkaConsumerGroupMembers,
				float64(len(group.Members))).ATag(metric.DefaultKafkaGroupIdTag, group.GroupId)

			if offsetFetchResponse, err := broker.FetchOffset(&offsetFetchRequest); err != nil {
				logger.Main.Error("Cannot get offset of group",
					zap.String("groudId", group.GroupId), zap.Error(err))
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
									zap.Int32("partition", partition), zap.Error(err))
								continue
							}
							currentOffset := offsetFetchResponseBlock.Offset
							currentOffsetSum += currentOffset

							kafkaAggregator.NewMetric(metric.KafkaConsumerGroupCurrentOffset,
								float64(currentOffset)).ATag(metric.DefaultKafkaGroupIdTag,
								group.GroupId).ATag(metric.DefaultKafkaTopicTag, topic).ATag(metric.DefaultKafkaPartitionTag,
								strconv.FormatInt(int64(partition), 10))

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
								kafkaAggregator.NewMetric(metric.KafkaConsumerGroupLag,
									float64(lag)).ATag(metric.DefaultKafkaGroupIdTag, group.GroupId).ATag(metric.DefaultKafkaTopicTag,
									topic).ATag(metric.DefaultKafkaPartitionTag, strconv.FormatInt(int64(partition), 10))
							} else {
								logger.Main.Error("No offset of topic", zap.Int32("partition", partition),
									zap.String("topic", topic), zap.Error(err))
							}
							mu.Unlock()
						}

						kafkaAggregator.NewMetric(metric.KafkaConsumerGroupCurrentOffsetSum,
							float64(currentOffsetSum)).ATag(metric.DefaultKafkaGroupIdTag, group.GroupId).ATag(metric.DefaultKafkaTopicTag, topic)
						kafkaAggregator.NewMetric(metric.KafkaConsumerGroupLagSum,
							float64(lagSum)).ATag(metric.DefaultKafkaGroupIdTag, group.GroupId).ATag(metric.DefaultKafkaTopicTag, topic)
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

// Can readCertAndKey returns true if the certificate and key files already exists,
// otherwise returns false. If lost one of cert and key, returns error.
func canReadCertAndKey(certPath, keyPath string) (bool, error) {
	certReadable := canReadFile(certPath)
	keyReadable := canReadFile(keyPath)

	if !certReadable && !keyReadable {
		return false, nil
	}

	if !certReadable {
		logger.Main.Error("Error reading , certificate and key must be supplied as a pair",
			zap.String("certPath", certPath))
	}

	if !keyReadable {
		logger.Main.Error("Error reading , certificate and key must be supplied as a pair",
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
