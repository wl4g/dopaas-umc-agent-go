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
	"time"
)

//
// Log constants definitions.
//

var LOG_BASE_DIR string = "./log/"
var LOG_MAIN_FILENAME string = LOG_BASE_DIR + "main.log"
var LOG_HTTP_FILENAME string = LOG_BASE_DIR + "http.log"
var LOG_LMDB_FILENAME string = LOG_BASE_DIR + "lmdb.log"
var LOG_GATEWAY_FILENAME string = LOG_BASE_DIR + "gateway.log"
var LOG_REDIS_FILENAME string = LOG_BASE_DIR + "redis.log"
var LOG_KAFKA_FILENAME string = LOG_BASE_DIR + "kafka.log"

//
// Config constants definitions.
//

// Default profile path.
var CONF_DEFAULT_FILENAME string = "/etc/umc-agent.yml"

// Default netcard name.
var CONF_DEFAULT_NETCARD string = "eth0"

// Default scanning frequency (in milliseconds).
var CONF_DEFAULT_DELAY time.Duration = 10000

// Default collection ports range.
var CONF_DEFAULT_COLLECT_PARTS string = "22,6380"

// Default local umc-agent instance ID.
var CONF_DEFAULT_INSTANCEID string = "UNKNOW"
