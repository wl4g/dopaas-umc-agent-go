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
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"umc-agent/pkg/config"
	"umc-agent/pkg/indicators"
	"umc-agent/pkg/logger"
)

// Send metrics to http gateway
func doHttpSend(aggregator *indicators.MetricAggregator) {
	if !config.GlobalConfig.Transport.Kafka.Enabled {
		request, err := http.NewRequest("POST", config.GlobalConfig.Transport.Http.ServerGateway,
			bytes.NewReader(aggregator.ToProtoBufArray()))
		if err != nil {
			panic(fmt.Sprintf("Create post request failed! %s", err))
		}
		//request.Header.Set("Content-Type", "application/json")
		//request.Header.Set("Accept", "application/json, text/plain, */*")

		//request.Header.Set("Connection", "keep-alive")

		// Do request.
		httpClient := &http.Client{Timeout: 30000}
		resp, err := httpClient.Do(request)
		defer resp.Body.Close()
		if err != nil {
			logger.Main.Error("Post failed", zap.Error(err))
		} else {
			ret, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Receive.Error("Failed to get response body", zap.Error(err))
			} else {
				logger.Receive.Info("Receive response message", zap.String("data", string(ret)))
			}
		}
	}

}
