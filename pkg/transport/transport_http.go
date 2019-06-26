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
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"umc-agent/pkg/config"
	"umc-agent/pkg/logger"
)

// Send indicators to http gateway
func doPostSend(key string, data string) {
	request, _ := http.NewRequest("POST", config.GlobalConfig.Launcher.Http.ServerGateway+"/"+key, strings.NewReader(data))
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Main.Error("Post failed", zap.Error(err))
	} else {
		ret, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Main.Error("Failed to get response body", zap.Error(err))
		} else {
			logger.Main.Info("Receive response message", zap.String("data", string(ret)))
		}
	}
}
