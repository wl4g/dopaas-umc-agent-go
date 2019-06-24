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
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
)

// 提交数据
func DoSendSubmit(ty string, v interface{}) {
	data := common.ToJSONString(v)

	// Kafka launcher takes precedence over HTTP launcher.
	if config.GlobalConfig.Launcher.Kafka.Enabled == true {
		doSendKafka(ty, data)
	} else {
		doSendHttp(ty, data)
	}
}

// Send indicators-data to Kafka
func doSendKafka(ty string, data string) {
	doProducer(ty, data)
}

type Launcher interface {
	doSend(ty string, data string)
}
