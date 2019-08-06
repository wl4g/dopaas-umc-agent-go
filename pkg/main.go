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
	"github.com/wl4g/super-devops-umc-agent/pkg/config"
	"github.com/wl4g/super-devops-umc-agent/pkg/constant"
	"github.com/wl4g/super-devops-umc-agent/pkg/indicators/redis"
	"github.com/wl4g/super-devops-umc-agent/pkg/logger"
	"github.com/wl4g/super-devops-umc-agent/pkg/transport"
	"sync"
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

	//test
	testAlarm()
}

func main() {
	startIndicatorRunners(wg)
	wg.Wait()
}

// Starting indicator runners all.
func startIndicatorRunners(wg *sync.WaitGroup) {
	wg.Add(1)
	/*go physical.IndicatorRunner()

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
	go cassandra.IndicatorRunner()*/

	go redis.IndicatorRunner()
}

// Testing
func testingIfNecessary() {
	/*var aggregator = indicators.NewMetricAggregator("Kafka")
	aggregator.NewMetric("kafka_partition_current_offset", 10.12).ATag("topic", "testTopic1").ATag("partition", "1")

	aggregator.NewMetric("kafka_partition_current_offset", 10.12).ATag("topic", "testTopic1").ATag("partition", "2")

	fmt.Println(aggregator.ToJSONString())
	fmt.Println(aggregator.String())
	fmt.Println(aggregator.ToProtoBufArray())

	os.Exit(0)*/
}

func testAlarm() {
	/*aggregator := indicators.NewMetricAggregator("Test")
	aggregator.NewMetric("hwjtest.cpu", float64(79)).ATag("tag1","1").ATag("tag2","2")
	transport.SendMetrics(aggregator)
	os.Exit(0)*/
}
