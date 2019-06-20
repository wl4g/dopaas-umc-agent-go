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
package launcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
)

// 提交数据
func DoSendSubmit(ty string, v interface{}) {
	data := common.ToJSONString(v)

	if config.GlobalConfig.Provider == "kafka" {
		doSendKafka(ty, data)
	} else {
		doSendHttp(ty, data)
	}
}

// Monitor data to HTTP gateway
func doSendHttp(ty string, data string) {
	request, _ := http.NewRequest("POST", config.GlobalConfig.ServerGateway+"/"+ty, strings.NewReader(data))
	//json
	request.Header.Set("Content-Type", "application/json")
	//post数据并接收http响应
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("post data error:%v\n", err)
	} else {
		fmt.Println("post a data successful.")
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("response data:%v\n", string(respBody))
	}
}

// Monitor data to Kafka gateway
func doSendKafka(ty string, data string) {
	doProducer(data)
}
