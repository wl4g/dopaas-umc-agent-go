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

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/log"
)

var GlobalPropertiesObj GlobalProperties

type GlobalProperties struct {
	ServerUri                  string                  `yaml:"server-uri"`
	Provider                   string                  `yaml:"provider"`
	Batch                      bool                    `yaml:"batch"`
	PhysicalPropertiesObj      PhysicalProperties      `yaml:"physical"`
	KafkaProducerPropertiesObj KafkaProducerProperties `yaml:"kafka"`
}

type PhysicalProperties struct {
	Delay      time.Duration `yaml:"delay"`
	Net        string        `yaml:"net"`
	GatherPort string        `yaml:"gather-port"`
}

type KafkaProducerProperties struct {
	Url        string `yaml:"url"`
	Topic      string `yaml:"topic"`
	Partitions int32  `yaml:""`
}

func GetConf(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.MainLogger.Info("yamlFile.Get err - ", zap.Error(err))
	}
	err = yaml.Unmarshal(yamlFile, &GlobalPropertiesObj)
	if err != nil {
		log.MainLogger.Info("Unmarshal - ", zap.Error(err))
	}

	// Set Default
	if GlobalPropertiesObj.ServerUri == "" {
		GlobalPropertiesObj.ServerUri = common.CONF_DEFAULT_SERVER_URI
	}
	if GlobalPropertiesObj.PhysicalPropertiesObj.Net == "" {
		GlobalPropertiesObj.PhysicalPropertiesObj.Net = common.CONF_DEFAULT_NETCARD
	}
	if GlobalPropertiesObj.PhysicalPropertiesObj.Delay == 0 {
		GlobalPropertiesObj.PhysicalPropertiesObj.Delay = common.CONF_DEFAULT_DELAY
	}
	if GlobalPropertiesObj.KafkaProducerPropertiesObj.Url == "" {
		GlobalPropertiesObj.KafkaProducerPropertiesObj.Url = common.CONF_DEFAULT_KAFKA_URL
	}
	if GlobalPropertiesObj.KafkaProducerPropertiesObj.Topic == "" {
		GlobalPropertiesObj.KafkaProducerPropertiesObj.Topic = common.CONF_DEFAULT_KAFKA_TOPIC
	}
	if GlobalPropertiesObj.Provider == "" {
		GlobalPropertiesObj.Provider = common.CONF_DEFAULT_KAFKA_TOPIC
	}


}
