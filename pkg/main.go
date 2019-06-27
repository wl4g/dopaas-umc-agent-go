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
	"umc-agent/pkg/logger"
	"umc-agent/pkg/monitor/cassandra"
	"umc-agent/pkg/monitor/consul"
	"umc-agent/pkg/monitor/docker"
	"umc-agent/pkg/monitor/elasticsearch"
	"umc-agent/pkg/monitor/emq"
	"umc-agent/pkg/monitor/etcd"
	"umc-agent/pkg/monitor/kafka"
	"umc-agent/pkg/monitor/memcached"
	"umc-agent/pkg/monitor/mesos"
	"umc-agent/pkg/monitor/mongodb"
	"umc-agent/pkg/monitor/mysqld"
	"umc-agent/pkg/monitor/opentsdb"
	"umc-agent/pkg/monitor/physical"
	"umc-agent/pkg/monitor/postgresql"
	"umc-agent/pkg/monitor/rabbitmq"
	"umc-agent/pkg/monitor/redis"
	"umc-agent/pkg/monitor/rocketmq"
	"umc-agent/pkg/monitor/zookeeper"
	"umc-agent/pkg/transport"
)

var (
	wg = &sync.WaitGroup{}
)

func init() {
	var confPath = constant.DefaultConfigPath

	// Command config path
	flag.StringVar(&confPath, "c", constant.DefaultConfigPath, "Config must is required!")
	flag.Parse()
	//flag.Usage()
	fmt.Printf("Initialize config path for - '%s'\n", confPath)

	// Init global config.
	config.InitGlobalConfig(confPath)

	// Init zap logger.
	logger.InitZapLogger()

	// Init kafka launcher.(if necessary)
	transport.InitKafkaLauncherIfNecessary()
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
	go mysqld.IndicatorRunner()
	go postgresql.IndicatorRunner()
	go opentsdb.IndicatorRunner()
	go cassandra.IndicatorRunner()
}
