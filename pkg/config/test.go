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
package config

import (
	"fmt"
)

type Config struct {
	App     string
	Port    int      `default:"8000"`
	IsDebug bool     `env:"DEBUG"`
	Hosts   []string `slice_sep:","`
	//Timeout time.Duration

	Redis struct {
		Version string `sep:""` // no sep between `CONFIG` and `REDIS`
		Host    string
		Port    int
	}

	MySQL struct {
		Version string `default:"5.7"`
		Host    string
		Port    int
	}
}

func main() {
	cfg := new(Config)
	err := Fill(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Home:", cfg.App)
	fmt.Println("Port:", cfg.Port)
	fmt.Println("IsDebug:", cfg.IsDebug)
	fmt.Println("Hosts:", cfg.Hosts, len(cfg.Hosts))
	//fmt.Println("Duration:", cfg.Timeout)
	fmt.Println("Redis_Version:", cfg.Redis.Version)
	fmt.Println("Redis_Host:", cfg.Redis.Host)
	fmt.Println("Redis_Port:", cfg.Redis.Port)
	fmt.Println("MySQL_Version:", cfg.MySQL.Version)
	fmt.Println("MySQL_Name:", cfg.MySQL.Host)
	fmt.Println("MySQL_port:", cfg.MySQL.Port)
}
