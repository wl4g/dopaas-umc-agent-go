//
// Copyright 2017 ~ 2025 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package launcher

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"umc-agent/pkg/config"
	"umc-agent/pkg/log"
)

var kafkaProducer sarama.SyncProducer

func buildKafkaProducer() {
	producerConfig := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll

	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner

	// 是否等待成功和失败后的响应
	producerConfig.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	var err error
	kafkaProducer, err = sarama.NewSyncProducer([]string{config.KafkaConf.Url}, producerConfig)
	if err != nil {
		panic(err)
	}
}

func doKafkaProducer(text string) {
	//构建发送的消息，
	msg := &sarama.ProducerMessage{
		Topic:     config.KafkaConf.Topic,
		Value:     sarama.ByteEncoder(text),
		Partition: int32(config.KafkaConf.Partitions), //
		//Key:        sarama.StringEncoder("key"),//
	}

	partition, offset, err := kafkaProducer.SendMessage(msg)

	if err != nil {
		log.MainLogger.Error("Send message Fail", zap.Error(err))
	}

	log.MainLogger.Info("Send message Success - ",
		zap.Int32("Partition", partition),
		zap.Int64("offset", offset))

}
