UMC一个轻量级基于SpringBoot Cloud开发完全开源的统一监控告警平台，也是Super DevOps平台中的核心子系统之一

![Build Status](https://api.travis-ci.org/wl4g/super-devops-umc-agent.svg?branch=master)

English version goes [here](README.md).

## 背景
虽然，目前有流行的开源监控告警解决方案，如：zabbix、prometheus+grafana，
- zabbix还是相对比较传统，更偏向于静态服务的监控，在当下这docker横行的动态服务世界里貌似显得你不从心；

- grafana界面设计非常棒，也提供了较灵活性的Dashboard配置，但若直接用于企业生产环境还存在一些问题，如：无支持多租户(作者目前负责
的某物联网paas/saas云平台，多租户是必须阿)，好像也没有较好HA方案；

- Prom生态相对覆盖面广，各种各样的xx_exporter插件，但在大规模部署时还存在一些问题，最突出的问题就是
Prom并无成熟的分布式集群方案（只有联邦集群，部署多个Prom实例，手动配置划分控目标），自然单机性能存在瓶颈，
当某Prom实例挂掉之后只能手动重新划分集群，且Prom实例之间数据不能共享（当然可再使用第三方插件包装一层，
如：Thanos，但这样架构太复杂，生产环境故障难排查），再就是pull模型在兼容多数据源上显得力不从心，比如我们有
场景需要指定精确的时间，还有比如我们有些数据是从日志来的或是从Kafka来的，这些都没有现成的方案，另外自定义灵
活的告警规则或后续整合spark/ML等实现无阀值智能告警等，二次开发可能没有java便捷，同时pull模型在需要实现认证
时会没有push模型容易(因为每个exporter都需暴漏一个端口，而pull只需在umc服务端暴漏一个端口即可)，pull模型
在大规模监控时，定时拉取expoter指标的线程池必然对服务端进程造成较大负担.

- UMC自研push模型无需维护pull模型所需线程池，服务端负担小，直接负载均衡部署即可，支持的第三方服务监控种类
也可丰富，且只有一个umc-agent可执行文件的客户端，安装起来不能更简单，在配置更新方面也支持动态热更新（如更新采集时间间隔），
后端基于springcloud，扩展性就不用担心了，二次开发也非常容易，老板视角：很容易找到稍好点的java就能维护起来！

[VS Prometheus](VS_PROMUS_CN.md)

## 快速开始

#### 开发环境安装
```
cd ${PROJECT_HOME}
go run pkg/main.go -c resources/etc/umc-agent.yml
```

#### 生产环境安装
- [Windows](scripts/build.bat)
- [Linux](scripts/build.sh)