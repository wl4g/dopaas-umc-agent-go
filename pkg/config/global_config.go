//
// Copyright 2017 ~ 2025 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/log"
)

var conf Conf

type Conf struct {
	ServerUri          string    `yaml:"server-uri"`
	PostMode           string    `yaml:"post-mode"`
	TogetherOrSeparate string    `yaml:"together-or-separate"`
	Physical           Physical  `yaml:"physical"`
	KafkaConf          KafkaConf `yaml:"kafka"`
}

type Physical struct {
	Delay      time.Duration `yaml:"delay"`
	Net        string        `yaml:"net"`
	GatherPort string        `yaml:"gather-port"`
}

type KafkaConf struct {
	Url        string `yaml:"url"`
	Topic      string `yaml:"topic"`
	Partitions int32  `yaml:""`
}

func getConf(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.MainLogger.Info("yamlFile.Get err - ", zap.Error(err))
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.MainLogger.Info("Unmarshal - ", zap.Error(err))
	}

	// Set Default
	if conf.ServerUri == "" {
		conf.ServerUri = common.CONF_DEFAULT_SERVER_URI
	}
	if conf.Physical.Net == "" {
		conf.Physical.Net = common.CONF_DEFAULT_NETCARD
	}
	if conf.Physical.Delay == 0 {
		conf.Physical.Delay = common.CONF_DEFAULT_DELAY
	}
	if conf.KafkaConf.Url == "" {
		conf.KafkaConf.Url = common.CONF_DEFAULT_KAFKA_URL
	}
	if conf.KafkaConf.Topic == "" {
		conf.KafkaConf.Topic = common.CONF_DEFAULT_KAFKA_TOPIC
	}
	if conf.PostMode == "" {
		conf.PostMode = common.CONF_DEFAULT_KAFKA_TOPIC
	}
	if conf.TogetherOrSeparate == "" {
		conf.TogetherOrSeparate = common.CONF_DFEFAULT_SUBMIT_MODE
	}

}
