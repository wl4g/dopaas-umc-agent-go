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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var commandPath = "./netCommand.txt"
var command string

//var sumCommand = "ss -n sport == 22|awk '{sumup += $3};{sumdo += $4};END {print sumup,sumdo}'"

func getNet(port string) string {
	if command == "" {
		command = ReadAll(commandPath)
	}
	command2 := strings.Replace(command, "#{port}", port, -1)
	s, _ := exec_shell(command2)
	fmt.Println(s)
	return s
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func exec_shell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	checkErr(err)
	return out.String(), err
}

func getDocker() [] DockerInfo{
	var command string  = "docker stats --no-stream --format \"{\\\"containerId\\\":\\\"{{.ID}}\\\",\\\"name\\\":\\\"{{.Name}}\\\",\\\"cpuPerc\\\":\\\"{{.CPUPerc}}\\\",\\\"memUsage\\\":\\\"{{.MemUsage}}\\\",\\\"memPerc\\\":\\\"{{.MemPerc}}\\\",\\\"netIO\\\":\\\"{{.NetIO}}\\\",\\\"blockIO\\\":\\\"{{.BlockIO}}\\\",\\\"PIDs\\\":\\\"{{.PIDs}}\\\"}\""
	s, _ := exec_shell(command)
	var dockerInfos [] DockerInfo
	if (s != "") {
		s = strings.ReplaceAll(s, "\n", ",")
		s = strings.TrimSuffix(s,",")
		s = "[" + s + "]"
		json.Unmarshal([]byte(s), &dockerInfos)
	}
	return dockerInfos;
}

type DockerInfo struct {
	ContainerId   string        `json:"containerId"`
	Name string                 `json:"name"`
	CpuPerc string                 `json:"cpuPerc"`
	MemUsage string                 `json:"memUsage"`
	MemPerc string                 `json:"memPerc"`
	NetIO string                 `json:"netIO"`
	BlockIO string                 `json:"blockIO"`
	PIDs string                 `json:"PIDs"`
}

//错误处理函数
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func ReadAll(filePth string) string {
	f, err := os.Open(filePth)
	if err != nil {
		panic(err)
	}
	s, _ := ioutil.ReadAll(f)
	return string(s)
}
