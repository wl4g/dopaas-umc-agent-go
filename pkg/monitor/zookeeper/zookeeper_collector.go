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
package zookeeper

import (
	"fmt"
	"net"
	"strings"
	"time"
	"umc-agent/pkg/common"
	"umc-agent/pkg/config"
	"umc-agent/pkg/launcher"
	"umc-agent/pkg/monitor/physical"
)

func ZookeeperIndicatorsRunner()  {
	for true {
		zookeeper := getByconf()
		launcher.DoSendSubmit(zookeeper.Type, zookeeper)
		time.Sleep(config.GlobalConfig.Indicators.Zookeeper.Delay * time.Millisecond)
	}
}


func getByconf() Zookeeper{

	var zookeeper Zookeeper
	zookeeper.Id = physical.PhysicalIndicatorId
	zookeeper.Type = "zookeeper"

	comm := strings.Split(config.GlobalConfig.Indicators.Zookeeper.Command,",")
	property := strings.Split(config.GlobalConfig.Indicators.Zookeeper.Properties,",")

	var infoSum string
	for _, command :=  range comm{
		info,_ := getZookeeperInfo(command)
		infoSum  = infoSum + info
	}
	infos := AnalysisZk(infoSum,property)
	zookeeper.Properties = infos
	return zookeeper

}

func getZookeeperInfo(comm string) (string,error) {
	conn, err := net.Dial("tcp", config.GlobalConfig.Indicators.Zookeeper.Host)
	if err != nil {
		fmt.Println("客户端建立连接失败")
		return "",err
	}

	//返回一个拥有 默认size 的reader，接收客户端输入
	//reader := bufio.NewReader(os.Stdin)
	//缓存 conn 中的数据
	buf := make([]byte, 1024)
	fmt.Println("请输入客户端请求数据...")
	//客户端输入
	input := comm
	//去除输入两端空格
	input = strings.TrimSpace(input)
	//客户端请求数据写入 conn，并传输
	conn.Write([]byte(input))
	//服务器端返回的数据写入空buf
	cnt, err := conn.Read(buf)

	if err != nil {
		fmt.Printf("客户端读取数据失败 %s\n", err)
		if conn!=nil {
			conn.Close()
		}
		return "",err
	}
	//回显服务器端回传的信息
	//fmt.Print("服务器端回复" + string(buf[0:cnt]))
	conn.Close();
	return string(buf[0:cnt]),nil
}

func AnalysisZk(info string,property []string) map[string]string {
	var mapInfo map[string]string
	mapInfo = make(map[string]string)
	infos := strings.Split(info,"\n")
	for _, line :=  range infos{

		i := strings.Split(line,"\t")
		if len(i)!=2 {
			continue
		}
		s1 := i[0]
		if(!common.StringsContains(property,s1)){
			continue
		}
		//fmt.Println(s1)
		s2 := i[1]
		s2 = strings.TrimSpace(s2)
		//fmt.Println(s2)
		mapInfo [ s1 ] = s2
	}
	return mapInfo
}