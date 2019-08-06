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
package metric

const (
	/** metric -- cpu */
	PHYSICAL_CPU = "host.cpu"

	/** metric -- mem : total */
	PHYSICAL_MEM_TOTAL = "host.mem.total"
	/** metric -- mem : free */
	PHYSICAL_MEM_FREE = "host.mem.free"
	/** metric -- mem : used percent */
	PHYSICAL_MEM_USED_PERCENT = "host.mem.usedPercent"
	/** metric -- mem : used */
	PHYSICAL_MEM_USED = "host.mem.used"
	/** metric -- mem cached */
	PHYSICAL_MEM_CACHE = "host.mem.cached"
	/** metric -- mem buffers */
	PHYSICAL_MEM_BUFFERS = "host.mem.buffers"

	/** metric -- disk : total */
	PHYSICAL_DISK_TOTAL = "host.disk.total"
	/** metric -- disk : free */
	PHYSICAL_DISK_FREE = "host.disk.free"
	/** metric -- disk : used */
	PHYSICAL_DISK_USED = "host.disk.used"
	/** metric -- disk : used Percent */
	PHYSICAL_DISK_USED_PERCENT = "host.disk.usedPercent"
	/** metric -- disk : inodes Physical */
	PHYSICAL_DISK_INODES_TOTAL = "host.disk.inodesTotal"
	/** metric -- disk : inodes Used */
	PHYSICAL_DISK_INODES_USED = "host.disk.inodesUsed"
	/** metric -- disk : inodes Free */
	PHYSICAL_DISK_INODES_FREE = "host.disk.inodesFree"
	/** metric -- disk : inodes Used Percent */
	PHYSICAL_DISK_INODES_USED_PERCENT = "host.disk.inodesUsedPercent"

	/** metric -- net : up */
	PHYSICAL_NET_UP = "host.net.up"
	/** metric -- net : down */
	PHYSICAL_NET_DOWN = "host.net.down"
	/** metric -- net : count */
	PHYSICAL_NET_COUNT = "host.net.count"
	/** metric -- net : estab */
	PHYSICAL_NET_ESTAB = "host.net.estab"
	/** metric -- net : closeWait */
	PHYSICAL_NET_CLOSE_WAIT = "host.net.closeWait"
	/** metric -- net : timeWait */
	PHYSICAL_NET_TIME_WAIT = "host.net.timeWait"
	/** metric -- net : close */
	PHYSICAL_NET_CLOSE = "host.net.close"
	/** metric -- net : listen */
	PHYSICAL_NET_LISTEN = "host.net.listen"
	/** metric -- net : closing */
	PHYSICAL_NET_CLOSING = "host.net.closing"
)
