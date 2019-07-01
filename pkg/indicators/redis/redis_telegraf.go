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
package redis

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/influxdata/telegraf/interna/tls"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"umc-agent/pkg/config"
	"umc-agent/pkg/constant"
	"umc-agent/pkg/constant/metric"
	"umc-agent/pkg/indicators"
	"umc-agent/pkg/logger"
	"umc-agent/pkg/transport"
)



func IndicatorRunner() {
	if !config.GlobalConfig.Indicator.Kafka.Enabled {
		logger.Main.Warn("No enabled kafka metrics runner!")
		return
	}
	logger.Main.Info("Starting kafka indicators runner ...")

	/*for true {
		// New redis metric aggregator
		redisAggregator := indicators.NewMetricAggregator("Redis")

		now := time.Now().UnixNano() / 1e6
		redisAggregator.Timestamp = now
		// Do collect.
		redis := Redis2{}
		redis.Gather(redisAggregator)

		// Send to servers.
		transport.DoSendSubmit(constant.Metric, &redisAggregator)

		// Sleep.
		time.Sleep(config.GlobalConfig.Indicator.Redis.Delay)
	}*/

	//Test
	now := time.Now().UnixNano() / 1e6
	redisAggregator := indicators.NewMetricAggregator("Redis")
	redisAggregator.Timestamp = now
	redis := Redis2{}
	redis.Gather(redisAggregator)
	transport.DoSendSubmit(constant.Metric, &redisAggregator)

}




type Redis2 struct {
	Servers  []string
	Password string
	tls.ClientConfig

	clients     []Client
	initialized bool
}

type Client interface {
	Info() *redis.StringCmd
	BaseTags() map[string]string
}

type RedisClient struct {
	client *redis.Client
	tags   map[string]string
}

func (r *RedisClient) Info() *redis.StringCmd {
	return r.client.Info()
}

func (r *RedisClient) BaseTags() map[string]string {
	tags := make(map[string]string)
	for k, v := range r.tags {
		tags[k] = v
	}
	return tags
}

var sampleConfig = `
  ## specify servers via a url matching:
  ##  [protocol://][:password]@address[:port]
  ##  e.g.
  ##    tcp://localhost:6379
  ##    tcp://:password@192.168.99.100
  ##    unix:///var/run/redis.sock
  ##
  ## If no servers are specified, then localhost is used as the host.
  ## If no port is specified, 6379 is used
  servers = ["tcp://localhost:6379"]

  ## specify server password
  # password = "s#cr@t%"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = true
`

func (r *Redis2) SampleConfig() string {
	return sampleConfig
}

func (r *Redis2) Description() string {
	return "Read metrics from one or many redis servers"
}

var Tracking = map[string]string{
	"uptime_in_seconds": "uptime",
	"connected_clients": "clients",
	"role":              "replication_role",
}

func (r *Redis2) init() error {
	if r.initialized {
		return nil
	}

	if len(r.Servers) == 0 {
		r.Servers = []string{"tcp://safecloud-test:6379"}
	}

	r.clients = make([]Client, len(r.Servers))

	for i, serv := range r.Servers {
		if !strings.HasPrefix(serv, "tcp://") && !strings.HasPrefix(serv, "unix://") {
			log.Printf("W! [inputs.redis]: server URL found without scheme; please update your configuration file")
			serv = "tcp://" + serv
		}

		u, err := url.Parse(serv)
		if err != nil {
			return fmt.Errorf("Unable to parse to address %q: %v", serv, err)
		}

		password := "zzx!@#$%"
		if u.User != nil {
			pw, ok := u.User.Password()
			if ok {
				password = pw
			}
		}
		if len(r.Password) > 0 {
			password = r.Password
		}

		var address string
		if u.Scheme == "unix" {
			address = u.Path
		} else {
			address = u.Host
		}

		tlsConfig, err := r.ClientConfig.TLSConfig()
		if err != nil {
			return err
		}

		client := redis.NewClient(
			&redis.Options{
				Addr:      address,
				Password:  password,
				Network:   u.Scheme,
				PoolSize:  1,
				TLSConfig: tlsConfig,
			},
		)

		tags := map[string]string{}
		if u.Scheme == "unix" {
			tags["socket"] = u.Path
		} else {
			tags["server"] = u.Hostname()
			tags["port"] = u.Port()
		}

		r.clients[i] = &RedisClient{
			client: client,
			tags:   tags,
		}
	}

	r.initialized = true
	return nil
}

// Reads stats from all configured servers accumulates stats.
// Returns one of the errors encountered while gather stats (if any).
func (r *Redis2) Gather(redisAggregator *indicators.MetricAggregator) error {
	if !r.initialized {
		err := r.init()
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup

	for _, client := range r.clients {
		wg.Add(1)
		go func(client Client) {
			defer wg.Done()
			//acc.AddError(r.gatherServer(client, acc))
			r.gatherServer(client,redisAggregator)
		}(client)
	}

	wg.Wait()
	return nil
}

func (r *Redis2) gatherServer(client Client,redisAggregator *indicators.MetricAggregator) error {
	info, err := client.Info().Result()
	if err != nil {
		return err
	}

	rdr := strings.NewReader(info)
	return gatherInfoOutput(rdr, client.BaseTags(),redisAggregator)
}

// gatherInfoOutput gathers
func gatherInfoOutput(
	rdr io.Reader,
	//acc telegraf.Accumulator,
	tags map[string]string,
	redisAggregator *indicators.MetricAggregator,
) error {
	var section string
	var keyspace_hits, keyspace_misses int64

	server := tags["server"]+":"+tags["port"]

	scanner := bufio.NewScanner(rdr)
	fields := make(map[string]interface{})
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			if len(line) > 2 {
				section = line[2:]
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		name := string(parts[0])

		if section == "Server" {
			if name != "lru_clock" && name != "uptime_in_seconds" && name != "redis_version" {
				continue
			}
		}

		if strings.HasPrefix(name, "master_replid") {
			continue
		}

		if name == "mem_allocator" {
			continue
		}

		if strings.HasSuffix(name, "_human") {
			continue
		}

		metri, ok := Tracking[name]
		if !ok {
			if section == "Keyspace" {
				kline := strings.TrimSpace(string(parts[1]))
				gatherKeyspaceLine(name, kline, tags)
				continue
			}
			metri = name
		}

		val := strings.TrimSpace(parts[1])

		// Try parsing as int
		if ival, err := strconv.ParseInt(val, 10, 64); err == nil {
			switch name {
			case "keyspace_hits":
				keyspace_hits = ival
			case "keyspace_misses":
				keyspace_misses = ival
			case "rdb_last_save_time":
				// influxdb can't calculate this, so we have to do it
				fields["rdb_last_save_time_elapsed"] = time.Now().Unix() - ival
			}
			fields[metri] = ival
			redisAggregator.NewMetric(metri,float64(ival)).ATag(metric.SERVER,server)
			continue
		}

		// Try parsing as a float
		if fval, err := strconv.ParseFloat(val, 64); err == nil {
			fields[metri] = fval
			redisAggregator.NewMetric(metri,fval).ATag(metric.SERVER,server)
			continue
		}

		// Treat it as a string

		if name == "role" {
			tags["replication_role"] = val
			continue
		}
		fields[metri] = val
	}
	var keyspace_hitrate float64 = 0.0
	if keyspace_hits != 0 || keyspace_misses != 0 {
		keyspace_hitrate = float64(keyspace_hits) / float64(keyspace_hits+keyspace_misses)
	}
	fields["keyspace_hitrate"] = keyspace_hitrate
	//acc.AddFields("redis", fields, tags)

	//fmt.Println(common.ToJSONString(fields))
	//fmt.Println(common.ToJSONString(tags))

	return nil
}

// Parse the special Keyspace line at end of redis stats
// This is a special line that looks something like:
//     db0:keys=2,expires=0,avg_ttl=0
// And there is one for each db on the redis instance
func gatherKeyspaceLine(
	name string,
	line string,
	//acc telegraf.Accumulator,
	global_tags map[string]string,
) {
	if strings.Contains(line, "keys=") {
		fields := make(map[string]interface{})
		tags := make(map[string]string)
		for k, v := range global_tags {
			tags[k] = v
		}
		tags["database"] = name
		dbparts := strings.Split(line, ",")
		for _, dbp := range dbparts {
			kv := strings.Split(dbp, "=")
			ival, err := strconv.ParseInt(kv[1], 10, 64)
			if err == nil {
				fields[kv[0]] = ival
			}
		}
		//acc.AddFields("redis_keyspace", fields, tags)
		//fmt.Println(common.ToJSONString(fields))
		//fmt.Println(common.ToJSONString(tags))
	}
}

/*func init() {
	inputs.Add("redis", func() telegraf.Input {
		return &Redis2{}
	})
}
*/