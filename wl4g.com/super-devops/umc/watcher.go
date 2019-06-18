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
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//physicalId
var physicalId string = "UNKNOWN"

var confPath string = CONF_DEFAULT_FILENAME

//初始化
func init() {
	//get conf path
	flag.StringVar(&confPath, "p", CONF_DEFAULT_FILENAME, "conf path")
	flag.Parse()
	//flag.Usage()//usage
	MainLogger.Info("confPath=" + confPath)

	//0618,配置改成用对象接收
	getConf(confPath)

	//init
	getId()//获取ip信息,作为Physical的标示

	//init kafka
	initKafka()
}

//主函数
func main() {

	if(conf.TogetherOrSeparate=="together"){
		go totalThread()
	}else{
		go memThread()
		go cpuThread()
		go diskThread()
		go netThread()
		go dockerThread()
	}
	for true {
		time.Sleep(100000 * time.Millisecond)
	}

	//for Test
	//go memThread()
	//go cpuThread()
	//go diskThread()
	//go netThread()
	//go dockerThread()
	//go totalThread()

	//close kafka (Meaningless)
	//defer producer.Close()
}

//mem
func memThread() {
	for true {
		var result Mem
		v, _ := mem.VirtualMemory()
		fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
		//fmt.Println(v)
		result.Id = physicalId
		result.Type = "mem"
		result.Mem = v
		fmt.Println("result = " + toJsonString(result))
		post("mem", result)
		time.Sleep(conf.Physical.Delay * time.Millisecond)
	}
}

//cpu
func cpuThread() {
	for true {
		var result Cpu
		p, _ := cpu.Percent(0, false)
		//p, _ := cpu.Times(true)
		fmt.Println(p)
		/*pa, _ := cpu.Percent(10000* time.Millisecond, true)
		fmt.Println(pa)*/
		result.Id = physicalId
		result.Type = "cpu"
		result.Cpu = p
		post("cpu", result)
		time.Sleep(conf.Physical.Delay * time.Millisecond)
	}
}

//disk
func diskThread() {
	for true {
		var result Disk
		disks := getDisks()
		fmt.Println(disks)
		result.Id = physicalId
		result.Type = "disk"
		result.Disks = disks
		post("disk", result)
		time.Sleep(conf.Physical.Delay * time.Millisecond)
	}
}

func getDisks() []DiskInfo {
	partitionStats, _ := disk.Partitions(false)
	var disks []DiskInfo
	for _, value := range partitionStats {
		var disk1 DiskInfo
		mountpoint := value.Mountpoint
		usageStat, _ := disk.Usage(mountpoint)
		disk1.PartitionStat = value
		disk1.Usage = *usageStat
		disks = append(disks, disk1)
	}
	return disks
}

//net
func netThread() {
	for true {
		var result NetInfos
		n := getNetInfo()
		result.Id = physicalId
		result.Type = "net"
		result.NetInfo = n
		post("net", result)
		time.Sleep(conf.Physical.Delay * time.Millisecond)
	}
}

func getNetInfo() []NetInfo {
	ports := strings.Split(conf.Physical.GatherPort, ",")
	//n, _ := net.IOCounters(true)
	//fmt.Println(n)
	//te, _ := net.Interfaces()
	//fmt.Println(te)
	var n []NetInfo
	for _, p := range ports {
		re := getNet(p)
		res := strings.Split(re, " ")
		if len(res) == 9 {
			var netinfo NetInfo
			netinfo.Port, _ = strconv.Atoi(p)
			netinfo.Up, _ = strconv.Atoi(res[0])
			netinfo.Down, _ = strconv.Atoi(res[1])
			netinfo.Count, _ = strconv.Atoi(res[2])
			netinfo.Estab, _ = strconv.Atoi(res[3])
			netinfo.CloseWait, _ = strconv.Atoi(res[4])
			netinfo.TimeWait, _ = strconv.Atoi(res[5])
			netinfo.Close, _ = strconv.Atoi(res[6])
			netinfo.Listen, _ = strconv.Atoi(res[7])
			netinfo.Closing, _ = strconv.Atoi(res[8])
			n = append(n, netinfo)
		}
	}
	return n
}

func dockerThread()  {
	for true {
		dockerInfo := getDocker()
		MainLogger.Info(toJsonString(dockerInfo))
		var result Docker
		result.Id = physicalId
		result.Type = "docker"
		result.DockerInfos = dockerInfo
		post("docker",result)
		time.Sleep(conf.Physical.Delay * time.Millisecond)
	}
}

func totalThread()  {
	for true {
		var result Total

		result.Id = physicalId
		result.Type = "total"

		v, _ := mem.VirtualMemory()
		result.Mem = v

		p, _ := cpu.Percent(0, false)
		result.Cpu = p

		disks := getDisks()
		result.Disks = disks

		n := getNetInfo()
		result.NetInfo = n

		dockerInfo := getDocker()
		result.DockerInfos = dockerInfo

		post("total",result)
		time.Sleep(conf.Physical.Delay * time.Millisecond)
	}
}




//提交数据
func post(ty string, v interface{}) {
	data := toJsonString(v)

	if conf.PostMode=="kafka" {
		postKafka(ty,data)
	}else{
		postHttp(ty,data)
	}
}

//post by http
func postHttp(ty string, data string) {
	request, _ := http.NewRequest("POST", conf.ServerUri+"/"+ty, strings.NewReader(data))
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

//post by kafka
func postKafka(ty string, data string) {
	send(data)
}



type Total struct {
	Id   string                 `json:"physicalId"`
	Type string                 `json:"type"`
	Mem  *mem.VirtualMemoryStat `json:"memInfo"`
	Cpu  []float64 `json:"cpu"`
	Disks []DiskInfo `json:"diskInfos"`
	NetInfo []NetInfo `json:"netInfos"`
	DockerInfos  []DockerInfo `json:"dockerInfos"`
}

type Mem struct {
	Id   string                 `json:"physicalId"`
	Type string                 `json:"type"`
	Mem  *mem.VirtualMemoryStat `json:"memInfo"`
}

type Cpu struct {
	Id   string    `json:"physicalId"`
	Type string    `json:"type"`
	Cpu  []float64 `json:"cpu"`
}

type Disk struct {
	Id    string     `json:"physicalId"`
	Type  string     `json:"type"`
	Disks []DiskInfo `json:"diskInfos"`
}

type DiskInfo struct {
	PartitionStat disk.PartitionStat `json:"partitionStat"`
	Usage         disk.UsageStat     `json:"usage"`
}

type NetInfos struct {
	Id      string    `json:"physicalId"`
	Type    string    `json:"type"`
	NetInfo []NetInfo `json:"netInfos"`
}

type NetInfo struct {
	Port      int `json:"port"`
	Up        int `json:"up"`
	Down      int `json:"down"`
	Count     int `json:"count"`
	Estab     int `json:"estab"`
	CloseWait int `json:"closeWait"`
	TimeWait  int `json:"timeWait"`
	Close     int `json:"close"`
	Listen    int `json:"listen"`
	Closing   int `json:"closing"`
}

type Docker struct {
	Id   string    `json:"physicalId"`
	Type string    `json:"type"`
	DockerInfos  []DockerInfo `json:"dockerInfos"`
}



func getId() {
	nets, _ := net.Interfaces()
	var found bool = false
	for _, value := range nets {
		if strings.EqualFold(conf.Physical.Net, value.Name) {
			hardwareAddr := value.HardwareAddr
			fmt.Println("found net card:" + hardwareAddr)
			physicalId = hardwareAddr
			reg := regexp.MustCompile(`(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})(\.(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})){3}`)
			for _, addr := range value.Addrs {
				add := addr.Addr
				if len(reg.FindAllString(add, -1)) > 0 {
					fmt.Println("found ip " + add)
					found = true
					//id = add+" "+id
					physicalId = add
					break
				}
			}
		}
	}
	if !found {
		panic("net found ip,Please check the net conf")
	}
	/*str := "10.0.0.26"
	matched, err := regexp.MatchString("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}", str)
	fmt.Println(matched, err)*/
}
