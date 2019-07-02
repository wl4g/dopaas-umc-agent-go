
## UMC (Unified Monitoring Center) Self-research VS Prometheus solution selection path

### Prometheus
#### Prepare environment:

A total of 4 hosts are used to collect kafka metric indicators:
```
prometheus-server(10.0.0.26, 10.0.0.57, 10.0.0.160)
kafka-exporter(10.0.0.12)
kafka(10.0.0.12)
```

- prometheus-server Configuration (10.0.0.26, 10.0.0.57, consistent configuration, pull the same kafka-exporter instance)
```
global:
  scrape_interval: 15s 
  external_labels:
    monitor: 'codelab-monitor'
scrape_configs:
  - job_name: 'mykafka'
    scrape_interval: 5s
    static_configs:
      - targets: ['10.0.0.26:9308']
```

- prometheus-server Configuration (10.0.0.160, federated configuration)
```
global:
  scrape_interval: 15s 
  external_labels:
    monitor: 'codelab-monitor'
scrape_configs:
  - job_name: 'federate'
    scrape_interval: 15s
    honor_labels: true
    metrics_path: '/federate'
    params:
      'match[]':
        - '{job="mykafka"}'
        - '{__name__=~"job:.*"}'
    static_configs:
      - targets:
        - '10.0.0.57:9090'
        - '10.0.0.26:9090'
```

#### Federal deploy verification
```
First start the kafka-exporter process of 12, and then start the prometheus-server process of 26, 57, 160 respectively.
Phenomenon: From the kafka-exporter log, you can see the requests from the two machines of 26, 57 (5 seconds each), and then separately
       Open the three prometheus-server consoles, and the data that is queried is the same and there is no duplication.
Conclusion: 26, 57 two machines are each pulled and stored in their respective databases.
```

#### HA experiment
```
Step 1: Only turn on 26, 160. Phenomenon: The data of the two machines is normally acquired and the data is consistent.
Step 2: Turn on 57 again. Phenomenon: All data during the shutdown period is lost.
```

#### Control case

|   Dimension     |   Research                                             |           Prometheus             |
| :--------------: | :---------------------------------------------------: | :-----------------------------:  |
|      HA        |       ✅                                               | Temporarily no mature plan        |
|Dynamic configuration|Support remote updates, such as: collection interval|Restart after updating the configuration|
|Alarm rule update| ✅Management side updates take effect immediately|After updating the configuration, there seems to be a problem with reload. Both the old and new rules take effect?|
|     存储      | Built in:openTSDB,derby，Provide store base classes for easy extension|Add the remote_read module to extend other data sources|
|JVM health    | ✅Integration:spring boot admin                  |           Not support                   |
|Distributed tracking| ✅Integration:spring cloud sleuth, zipkin        |           Not support                    |
| Complexity   |      Lightweight                                 |           Heavy                        |
|              |                                                  |                                        |
|              |                                                  |                                        |

