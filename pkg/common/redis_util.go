/**
 * Copyright 2017 ~ 2025 the original author or authors.
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
package common

import "strings"

func getRedisInfo(port int) string {

	var redisPwd string = ""
	//example: redis-cli -p 6380  -a 'zzx!@#$%' cluster info
	var comm1 string = "redis-cli -p "+string(port)+"  -a '"+redisPwd+"' cluster info"
	//example: redis-cli -p 6380  -a 'zzx!@#$%' info all
	var comm2 string = "redis-cli -p "+string(port)+"  -a '"+redisPwd+"' info all"
	re1,_ := ExecShell(comm1)
	re2,_ := ExecShell(comm2)
	return re1+re2
}

func getByconf(property []string)  {
	var port int = 8080
	re := getRedisInfo(port)
	AnalysisRedis(re,property)
}

func AnalysisRedis(info string,property []string) map[string]string {
	var mapInfo map[string]string
	infos := strings.Split(info,"\n")
	for _, line :=  range infos{
		i := strings.Index(line,":")
		s1 := line[0:i]
		if(!StringsContains(property,s1)){
			continue
		}
		//fmt.Println(s1)
		s2 := line[(i+1):len(line)]
		s2 = strings.TrimSpace(s2)
		//fmt.Println(s2)
		mapInfo [ s1 ] = s2
	}
	return mapInfo
}
