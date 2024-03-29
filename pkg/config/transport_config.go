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
package config

import "time"

// ---------------------
// Transport properties.
// ---------------------

type TransportProperties struct {
	Http  HttpTransportProperties  `yaml:"http"`
	Kafka KafkaTransportProperties `yaml:"kafka"`
}

type HttpTransportProperties struct {
	ServerGateway string `yaml:"server-gateway"`
}

type KafkaTransportProperties struct {
	Enabled      bool          `yaml:"enabled"`
	Servers      string        `yaml:"servers"`
	Ack          int16         `yaml:"ack"`
	Timeout      time.Duration `yaml:"timeout"`
	MetricTopic  string        `yaml:"metric-topic"`
	ReceiveTopic string        `yaml:"receive-topic"`
}
