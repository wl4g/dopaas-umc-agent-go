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
package transport

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strings"
	"sync"
	"umc-agent/pkg/config"
	"umc-agent/pkg/logger"
)

var (
	wg sync.WaitGroup
)

var kafkaProducer sarama.SyncProducer
var kafkaConsumer sarama.Consumer

// Init kafka producer launcher(if necessary)
func InitKafkaLauncherIfNecessary() {
	var kafkaProperties = config.GlobalConfig.Launcher.Kafka

	// Check kafka producer enabled?
	if !kafkaProperties.Enabled {
		logger.Main.Warn("No enabled kafka launcher!")
		return
	}
	logger.Main.Info("Kafka launcher Starting...")

	// Create producer
	createKafkaProducer(kafkaProperties)

	// Create consumer
	createKafkaConsumer(kafkaProperties)
}

// create kafkaProducer
func createKafkaProducer(kafkaProperties config.KafkaLauncherProperties) {
	// Configuration
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.RequiredAcks(kafkaProperties.Ack)
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaConfig.Producer.Timeout = kafkaProperties.Timeout
	kafkaConfig.Producer.Return.Successes = true

	// Create syncProducer
	var err error
	kafkaProducer, err = sarama.NewSyncProducer(strings.Split(kafkaProperties.BootstrapServers, ","), kafkaConfig)
	if err != nil {
		panic(err)
	}
}

// Create kafkaConsumer, See: https://github.com/Shopify/sarama/blob/master/examples/consumergroup/main.go
func createKafkaConsumer(kafkaProperties config.KafkaLauncherProperties) {
	// Configuration
	kafkaConfig := sarama.NewConfig()
	// Grouping cannot be specified here, otherwise shared subscriptions will be used,
	// and the configuration of server-side downloads will not be consumed by all agents.
	//kafkaConfig.Consumer.Group = ""
	kafkaConfig.Consumer.Return.Errors = true

	// Create consumer
	var err error
	kafkaConsumer, err = sarama.NewConsumer(strings.Split(kafkaProperties.BootstrapServers, ","), kafkaConfig)
	if err != nil {
		logger.Receive.Error("Failed to start consumer: %s", zap.Error(err))
		return
	}
	defer kafkaConsumer.Close()

	// Partitions all
	partitions, err := kafkaConsumer.Partitions(kafkaProperties.ReceiveTopic)
	if err != nil {
		logger.Receive.Error("Failed to get the list of partitions: ", zap.Error(err))
		return
	}

	for partition := range partitions {
		pc, err := kafkaConsumer.ConsumePartition(kafkaProperties.ReceiveTopic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logger.Receive.Error("Failed to start consumer",
				zap.Int("partition", partition),
				zap.Error(err))
			return
		}
		//defer pc.AsyncClose()
		wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				var key = string(msg.Key)
				var data = string(msg.Value)

				logger.Receive.Info("Receive consumer message",
					zap.Int32("partition", msg.Partition),
					zap.Int64("offset", msg.Offset),
					zap.String("key", key),
					zap.String("data", data))
			}
		}(pc)
	}
	wg.Wait()

	logger.Receive.Info("Finished kafka consumer")
}

// Do production
func doProducerSend(key string, data string) {
	if !config.GlobalConfig.Launcher.Kafka.Enabled {
		panic("No enabled kafka launcher!")
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: config.GlobalConfig.Launcher.Kafka.MetricTopic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := kafkaProducer.SendMessage(msg)
	if err != nil {
		logger.Main.Error("Sent failed", zap.Error(err))
	} else {
		logger.Main.Info("Sent completed",
			zap.String("key", key),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
			zap.String("data", data))
	}

}
