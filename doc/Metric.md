###Host
| Metric                      | 异常触发条件 |
| --------------------------- | ------------ |
| host.cpu                    | >=70%        |
| host.mem.total              |              |
| host.mem.free               | <1G          |
| host.mem.usedPercent        | >=86%        |
| host.mem.used               |              |
| host.mem.cached             |              |
| host.mem.buffers            |              |
| host.disk.total             |              |
| host.disk.free              |              |
| host.disk.used              |              |
| host.disk.usedPercent       | >=80%        |
| host.disk.inodesTotal       |              |
| host.disk.inodesUsed        |              |
| host.disk.inodesFree        |              |
| host.disk.inodesUsedPercent | >=80%        |
| host.net.up                 |              |
| host.net.down               |              |
| host.net.count              |              |
| host.net.estab              |              |
| host.net.closeWait          | >512         |
| host.net.timeWait           |              |
| host.net.close              |              |
| host.net.listen             |              |
| host.net.closing            |              |

###Kafka
| Metric                                           | 异常触发条件 |
| ------------------------------------------------ | ------------ |
| kafka_brokers                                    | <3           |
| kafka_topic_partition_current_offset             |              |
| kafka_topic_partition_in_sync_replica            |              |
| kafka_topic_partition_leader                     |              |
| kafka_topic_partition_leader_is_preferred        |              |
| kafka_topic_partition_oldest_offset              |              |
| kafka_topic_partition_replicas                   |              |
| kafka_topic_partition_under_replicated_partition |              |
| kafka_topic_partitions                           |              |
| kafka_consumergroup_members                      |              |
| kafka_consumergroup_current_offset               |              |
| kafka_consumergroup_lag                          | >10          |
| kafka_consumergroup_current_offset_sum           |              |
| kafka_consumergroup_lag_sum                      | >10          |

###Zookeeper
| Metric                                  | 异常触发条件           |
| --------------------------------------- | ---------------------- |
| zookeeper.zk.avg.latency                |                        |
| zookeeper.zk.max.latency                | >1000 最大响应延迟(ms) |
| zookeeper.zk.min.latency                |                        |
| zookeeper.zk.packets.received           |                        |
| zookeeper.zk.packets.sent               |                        |
| zookeeper.zk.num.alive.connections      | >500 活跃连接数        |
| zookeeper.zk.outstanding.requests       | >5堆积请求数           |
| zookeeper.zk.znode.count                |                        |
| zookeeper.zk.watch.count                |                        |
| zookeeper.zk.ephemerals.count           |                        |
| zookeeper.zk.approximate.data.size      |                        |
| zookeeper.zk.open.file.descriptor.count |                        |
| zookeeper.zk.max.file.descriptor.count  |                        |
| zookeeper.zk.followers                  |                        |
| zookeeper.zk.synced.followers           |                        |
| zookeeper.zk.pending.syncs              |                        |

###Redis

| Metric                                |          异常触发条件                          |
| ------------------------------------- | -------------------------------------------- |
| redis.uptime                          |                                              |
| redis.lru.clock                       |                                              |
| redis.clients                         | >100连接数                                    |
| redis.client.recent.max.input.buffer  |                                              |
| redis.client.recent.max.output.buffer |                                              |
| redis.blocked.clients                 |                                              |
| redis.used.memory                     | >1024x1024x1024内存使用(B)==1G                |
| redis.used.memory.peak                | >1024x1024x1024x1.2内存使用的峰值大小(B)==1G   |
| redis.used.memory.overhead            |                                              |
| redis.used.memory.startup             |                                              |
| redis.used.memory.dataset             |                                              |
| redis.allocator.allocated             |                                              |
| redis.allocator.active                |                                              |
| redis.allocator.resident              |                                              |
| redis.total.system.memory             |                                              |
| redis.used.memory.lua                 |                                              |
| redis.used.memory.scripts             |                                              |
| redis.number.of.cached.scripts        |                                              |
| redis.maxmemory                       |                                              |
| redis.allocator.frag.ratio            |                                              |
| redis.allocator.frag.bytes            |                                              |
| redis.allocator.rss.ratio             |                                              |
| redis.allocator.rss.bytes             |                                              |
| redis.rss.overhead.ratio              |                                              |
| redis.rss.overhead.bytes              |                                              |
| redis.mem.fragmentation.ratio         |                                              |
| redis.mem.fragmentation.bytes         |                                              |
| redis.mem.not.counted.for.evict       |                                              |
| redis.mem.replication.backlog         |                                              |
| redis.mem.clients.slaves              |                                              |
| redis.mem.clients.normal              |                                              |
| redis.mem.aof.buffer                  |                                              |
| redis.active.defrag.running           |                                              |
| redis.lazyfree.pending.objects        |                                              |
| redis.loading                         |                                              |
| redis.rdb.changes.since.last.save     |                                              |
| redis.rdb.bgsave.in.progress          |                                              |
| redis.rdb.last.save.time              |                                              |
| redis.rdb.last.bgsave.time.sec        | >                                            |
| redis.rdb.current.bgsave.time.sec     | >                                            |
| redis.rdb.last.cow.size               |                                              |
| redis.aof.enabled                     |                                              |
| redis.aof.rewrite.in.progress         |                                              |
| redis.aof.rewrite.scheduled           |                                              |
| redis.aof.last.rewrite.time.sec       |                                              |
| redis.aof.current.rewrite.time.sec    |                                              |
| redis.aof.last.cow.size               |                                              |
| redis.aof.current.size                |                                              |
| redis.aof.base.size                   |                                              |
| redis.aof.pending.rewrite             |                                              |
| redis.aof.buffer.length               |                                              |
| redis.aof.rewrite.buffer.length       |                                              |
| redis.aof.pending.bio.fsync           |                                              |
| redis.aof.delayed.fsync               |                                              |
| redis.total.connections.received      |                                              |
| redis.total.commands.processed        |                                              |
| redis.instantaneous.ops.per.sec       |                                              |
| redis.total.net.input.bytes           |                                              |
| redis.total.net.output.bytes          |                                              |
| redis.instantaneous.input.kbps        | >1024网络入口kps                             |
| redis.instantaneous.output.kbps       | >1024网络出口kps                             |
| redis.rejected.connections            |                                              |
| redis.sync.full                       |                                              |
| redis.sync.partial.ok                 |                                              |
| redis.sync.partial.err                |                                              |
| redis.expired.keys                    |                                              |
| redis.expired.stale.perc              |                                              |
| redis.expired.time.cap.reached.count  |                                              |
| redis.evicted.keys                    |                                              |
| redis.keyspace.hits                   |                                              |
| redis.keyspace.misses                 |                                              |
| redis.pubsub.channels                 |                                              |
| redis.pubsub.patterns                 |                                              |
| redis.latest.fork.usec                |                                              |
| redis.migrate.cached.sockets          |                                              |
| redis.slave.expires.tracked.keys      |                                              |
| redis.active.defrag.hits              |                                              |
| redis.active.defrag.misses            |                                              |
| redis.active.defrag.key.hits          |                                              |
| redis.active.defrag.key.misses        |                                              |
| redis.connected.slaves                |                                              |
| redis.master.repl.offset              |                                              |
| redis.second.repl.offset              |                                              |
| redis.repl.backlog.active             |                                              |
| redis.repl.backlog.size               |                                              |
| redis.repl.backlog.first.byte.offset  |                                              |
| redis.repl.backlog.histlen            |                                              |
| redis.used.cpu.sys                    |                                              |
| redis.used.cpu.user                   |                                              |
| redis.used.cpu.sys.children           |                                              |
| redis.used.cpu.user.children          |                                              |
| redis.cluster.enabled                 |                                              |

###docker

| Metric           |      |
| ---------------- | ---- |
| docker.cpu.perc  | >80  |
| docker.mem.usage |      |
| docker.mem.perc  | >80  |
| docker.net.in    |      |
| docker.net.out   |      |
| docker.block.in  |      |
| docker.block.out |      |


