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

import "time"

const (
	// -------------------------------
	// Launcher constants.
	// -------------------------------

	// Default http server gateway
	DefaultHttpServerGateway = "http://localhost:14046/umc/basic"

	// Default launcher kafka bootstrap servers.
	DefaultLauncherKafkaServers = "localhost:9092"

	// Default launcher kafka metric topic.
	DefaultLauncherKafkaMetricTopic = "__umc_agent_metric_"

	// Default launcher kafka receive topic.
	DefaultLauncherKafkaReceiveTopic = "__umc_agent_metric_receive"

	// Default launcher kafka producer ack.
	DefaultLauncherKafkaAck = 1

	// Default launcher kafka producer timeout.
	DefaultLauncherKafkaTimeout = 10 * time.Second
)
