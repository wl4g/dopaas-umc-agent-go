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
	// Transport constants.
	// -------------------------------

	// Default http server gateway
	DefaultHttpServerGateway = "http://localhost:14046/umc/basic"

	// Default transport kafka bootstrap servers.
	DefaultTransportKafkaServers = "localhost:9092"

	// Default transport kafka metric topic.
	DefaultTransportKafkaMetricTopic = "umc_agent_metrics"

	// Default transport kafka receive topic.
	DefaultTransportKafkaReceiveTopic = "umc_agent_receive"

	// Default transport kafka producer ack.
	DefaultTransportKafkaAck = 1

	// Default transport kafka producer timeout.
	DefaultTransportKafkaTimeout = 10 * time.Second
)
