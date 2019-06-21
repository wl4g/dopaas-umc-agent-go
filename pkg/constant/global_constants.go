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

//
// Basic constants.
//

// Default profile path.
var DefaultConfigPath = "/etc/umc-agent.yml"

//
// Log constants.
//

var DefaultLogBaseDir = "./log/"
var DefaultLogMainFilename = DefaultLogBaseDir + "main.log"
var DefaultLogHttpFilename = DefaultLogBaseDir + "http.log"
var DefaultLogGatewayFilename = DefaultLogBaseDir + "gateway.log"

//
// Launcher constants.
//

// Default http server gateway
var DefaultHttpServerGateway = "http://localhost:14046/umc/basic"

// Default launcher kafka bootstrap servers.
var DefaultLauncherKafkaServers = "localhost:9092"

// Default launcher kafka topic.
var DefaultLauncherKafkaTopic = "__devops_umc_agent_metric_"

// Default launcher kafka partitions.
var DefaultLauncherKafkaPartitions int32 = 10

//
// Indicators constants.
//

// Default netcard name.
var DefaultNetcard = "eth0"

// Default indicator delay.
var DefaultIndicatorsDelay time.Duration = 60000

// Default network indicator ports range.
var DefaultNetIndicatorPortRange = "22,6380"
