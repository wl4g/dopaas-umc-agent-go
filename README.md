UMC - A lightweight, fully open source unified monitoring and alerting platform based on SpringBoot Cloud, and one of the core subsystems in the Super DevOps platform.

[中文文档](README_CN.md).

## Background
Although there are currently popular open source monitoring and alerting solutions, such as: zabbix, prometheus + grafana.
- Zabbix is still relatively more traditional, and it is more inclined to monitor the static service. In the current dynamic service world of docker, it seems that you are not in the heart;

- The grafana interface design is very good, and it also provides a more flexible Dashboard configuration, but there are still some problems if it is directly used in the enterprise production environment, such as: no support for multi-tenancy (the author is currently responsible
  Some Internet of Things paas/saas cloud platform, multi-tenant is a must, it seems that there is no better HA solution;

- The Prom ecosystem has a relatively wide coverage and a variety of xx_exporter plugins, but there are still some problems in large-scale deployment. The most prominent problem is
  Prom does not have a mature distributed cluster solution (only federated clusters, deploying multiple Prom instances, manually configuring partition control targets), and there is a bottleneck in natural stand-alone performance.
  When a Prom instance hangs, you can only manually re-segment the cluster, and the data between the Prom instances cannot be shared (of course, you can use a third-party plug-in to wrap a layer.
  Such as: Thanos, but this architecture is too complicated, production environment failure is difficult to check), then the pull model is not working well on compatible multi-data sources, such as we have
  The scene needs to specify the exact time, and for example, some of our data is from the log or from Kafka, there is no ready-made solution, and the custom spirit
  Live alarm rules or subsequent integration of spark/ML, etc. to achieve no-threshold intelligent alarms, etc., secondary development may not have Java convenience, and the pull model needs to be authenticated.
  There will be no push model easy (because each exporter needs to leak a port, and pull only needs to leak a port on the umc server), pull model
  In large-scale monitoring, the thread pool that pulls the expoter indicator at regular time will inevitably impose a heavy burden on the server process.

- The UMC self-developed push model does not need to maintain the thread pool required by the pull model, the server has a small burden, and direct load balancing can be deployed. The supported third-party service monitoring types
  It can also be rich, and only one client of the umc-agent executable file can not be installed more easily, and dynamic hot update (such as update collection interval) is also supported in configuration update.
  The backend is based on springcloud, there is no need to worry about scalability, secondary development is also very easy, the boss's perspective: it is easy to find a slightly better java to maintain!

[VS Prometheus](VS_PROMUS_CN.md)

## Quick start

#### Development environment installation
```
cd ${PROJECT_HOME}
go run pkg/main.go -c resources/etc/umc-agent.yml
```

#### Production environment installation
- [Windows](scripts/build.bat)
- [Linux](scripts/build.sh)
