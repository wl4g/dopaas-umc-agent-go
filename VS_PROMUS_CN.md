
## UMC(统一监控中心) 自研 VS Prometheus 方案选型之路

### Prometheus
#### 准备环境:

共使用4台主机，以采集kafka metric指标:
```
prometheus-server(10.0.0.26, 10.0.0.57, 10.0.0.160)
kafka-exporter(10.0.0.12)
kafka(10.0.0.12)
```

- prometheus-server 配置(10.0.0.26, 10.0.0.57，配置一致，pull同一个kafka-exporter实例)
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

- prometheus-server 配置(10.0.0.160，联邦配置)
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

#### 联邦验证
```
首先启动12的kafka-exporter进程，然后再分别启动26, 57, 160的prometheus-server进程，
现象: 从kafka-exporter日志能看到分别来自26, 57两台机器的请求(分别每间隔5秒)，再分别
      打开三个prometheus-server控制台，查询到的数据也一样且都没有重复。
结论: 26, 57两台机器各自pull, 存储到各自的数据库, 
```

#### HA验证
```
步骤1: 只开启26, 160，现象: 两台机器数据正常获取, 且数据一致。
步骤2: 再开启57, 现象: 57停机期间的数据全部丢失。
```

#### 对照情况

|     对比     |       自研       |           prometheus            |
| :----------: | :--------------: | :-----------------------------: |
|      HA      |        ✅         |        暂时无成熟方案         |
|   动态配置   |                  |          修改配置文件后重启      |
| altermanager |        ✅         |        reload后有问题     |
|     集中     |                  |                                 |
|     存储     | 使用tsdb,derby等 | 可通过remote_read使用其他数据看 |
|    代码量    |      轻量级      |             比较重              |
|              |                  |                                 |
|              |                  |                                 |
|              |                  |                                 |





