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
package main

import (
	"flag"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/launcher"
	"umc-agent/pkg/log"
	"umc-agent/pkg/monitor/physical"
	"umc-agent/pkg/monitor/share"
	"umc-agent/pkg/monitor/virtual"
)

var confPath string = common.CONF_DEFAULT_FILENAME

//初始化
func init() {
	//get conf path
	flag.StringVar(&confPath, "p", common.CONF_DEFAULT_FILENAME, "conf path")
	flag.Parse()
	//flag.Usage()//usage
	log.MainLogger.Info("confPath=" + confPath)

	//0618,配置改成用对象接收
	config.GetConf(confPath)

	//init
	share.PhysicalId = common.GetPhysicalId() //获取ip信息,作为Physical的标示

	//init kafka
	launcher.BuildKafkaProducer()
}

//主函数
func main() {
	if config.GlobalPropertiesObj.Batch == true {
		go share.TotalThread()
	} else {
		go physical.MemThread()
		go physical.CpuThread()
		go physical.DiskThread()
		go physical.NetThread()
		go virtual.DockerThread()
	}
	for true {
		time.Sleep(100000 * time.Millisecond)
	}

	//for Test
	//memThread()
	//go cpuThread()
	//go diskThread()
	//go netThread()
	//go dockerThread()
	//go totalThread()

	//close kafka (Meaningless)
	//defer producer.Close()
}
