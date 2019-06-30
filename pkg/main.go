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
package main

import (
	"flag"
	"fmt"
	"sync"
	"umc-agent/pkg/config"
	"umc-agent/pkg/constant"
	"umc-agent/pkg/indicators/cassandra"
	"umc-agent/pkg/indicators/consul"
	"umc-agent/pkg/indicators/docker"
	"umc-agent/pkg/indicators/elasticsearch"
	"umc-agent/pkg/indicators/emq"
	"umc-agent/pkg/indicators/etcd"
	"umc-agent/pkg/indicators/kafka"
	"umc-agent/pkg/indicators/memcached"
	"umc-agent/pkg/indicators/mesos"
	"umc-agent/pkg/indicators/mongodb"
	"umc-agent/pkg/indicators/mysql"
	"umc-agent/pkg/indicators/opentsdb"
	"umc-agent/pkg/indicators/physical"
	"umc-agent/pkg/indicators/postgresql"
	"umc-agent/pkg/indicators/rabbitmq"
	"umc-agent/pkg/indicators/redis"
	"umc-agent/pkg/indicators/rocketmq"
	"umc-agent/pkg/indicators/zookeeper"
	"umc-agent/pkg/logger"
	"umc-agent/pkg/transport"
)

var (
	wg = &sync.WaitGroup{}
)

func init() {
	confPath := constant.DefaultConfigPath

	// Command config path
	flag.StringVar(&confPath, "c", constant.DefaultConfigPath, "Umc agent config path.")
	flag.Parse()
	//flag.Usage()
	fmt.Printf("Initialize config path for - '%s'\n", confPath)

	// Init global config.
	config.InitGlobalConfig(confPath)

	// Init zap logger.
	logger.InitZapLogger()

	// Testing if necessary.
	testingIfNecessary()

	// Init kafka launcher.(if necessary)
	transport.InitKafkaTransportIfNecessary()
}

func main() {
	startCollectorRunners(wg)
	wg.Wait()
}

// Starting indicator runners all
func startCollectorRunners(wg *sync.WaitGroup) {
	wg.Add(1)
	go physical.IndicatorRunner()

	go docker.IndicatorRunner()
	go mesos.IndicatorRunner()

	go zookeeper.IndicatorRunner()
	go etcd.IndicatorRunner()
	go consul.IndicatorRunner()

	go kafka.IndicatorRunner()
	go emq.IndicatorRunner()
	go rabbitmq.IndicatorRunner()
	go rocketmq.IndicatorRunner()
	go redis.IndicatorRunner()
	go memcached.IndicatorRunner()

	go elasticsearch.IndicatorRunner()
	go mongodb.IndicatorRunner()
	go mysql.IndicatorRunner()
	go postgresql.IndicatorRunner()
	go opentsdb.IndicatorRunner()
	go cassandra.IndicatorRunner()
}

// Testing
func testingIfNecessary() {
	//var aggregator = indicators2.NewMetricAggregator("Kafka")
	//aggregator.NewMetric("kafka_partition_current_offset", 10.12).ATag("topic", "testTopic1").ATag("partition", "1")
	//
	//fmt.Println(aggregator.ToJSONString())
	//fmt.Println(aggregator.String())
	//fmt.Println(aggregator.ToProtoBufArray())
	//
	//os.Exit(0)
}
