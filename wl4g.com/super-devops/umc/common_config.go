package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

var conf Conf

type Conf struct {
	ServerUri string   `yaml:"server-uri"`
	PostMode string		`yaml:"post-mode"`
	Physical Physical   `yaml:"physical"`
	KafkaConf KafkaConf `yaml:"kafka"`
}

type Physical struct {
	Delay    time.Duration `yaml:"delay"`
	Net      string `yaml:"net"`
	GatherPort string `yaml:"gather-port"`
}

type KafkaConf struct {
	Url string `yaml:"url"`
	Topic string `yaml:"topic"`
}

func getConf(path string)  {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	//Set Default
	if conf.ServerUri=="" {
		conf.ServerUri = CONF_DEFAULT_SERVER_URI
	}
	if conf.Physical.Net=="" {
		conf.Physical.Net = CONF_DEFAULT_NETCARD
	}
	if conf.Physical.Delay==0 {
		conf.Physical.Delay=CONF_DEFAULT_DELAY
	}
	if conf.KafkaConf.Url=="" {
		conf.KafkaConf.Url=CONF_DEFAULT_KAFKA_URL
	}
	if conf.KafkaConf.Topic=="" {
		conf.KafkaConf.Topic=CONF_DEFAULT_KAFKA_TOPIC
	}
	if conf.PostMode=="" {
		conf.PostMode=CONF_DEFAULT_KAFKA_TOPIC
	}

}

/*func main() {
	getConf("./conf.yml")
	fmt.Println(conf)
}*/