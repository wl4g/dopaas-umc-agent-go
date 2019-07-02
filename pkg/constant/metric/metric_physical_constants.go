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
	PHYSICAL_CPU = "physical.cpu";

	/** metric -- mem : total */
	PHYSICAL_MEM_TOTAL = "physical.mem.total";
	/** metric -- mem : free */
	PHYSICAL_MEM_FREE = "physical.mem.free";
	/** metric -- mem : used percent */
	PHYSICAL_MEM_USED_PERCENT = "physical.mem.usedPercent";
	/** metric -- mem : used */
	PHYSICAL_MEM_USED = "physical.mem.used";
	/** metric -- mem cached */
	PHYSICAL_MEM_CACHE = "physical.mem.cached";
	/** metric -- mem buffers */
	PHYSICAL_MEM_BUFFERS = "physical.mem.buffers";

	/** metric -- disk : total */
	PHYSICAL_DISK_TOTAL = "physical.disk.total";
	/** metric -- disk : free */
	PHYSICAL_DISK_FREE = "physical.disk.free";
	/** metric -- disk : used */
	PHYSICAL_DISK_USED = "physical.disk.used";
	/** metric -- disk : used Percent */
	PHYSICAL_DISK_USED_PERCENT = "physical.disk.usedPercent";
	/** metric -- disk : inodes Physical */
	PHYSICAL_DISK_INODES_TOTAL = "physical.disk.inodesTotal";
	/** metric -- disk : inodes Used */
	PHYSICAL_DISK_INODES_USED = "physical.disk.inodesUsed";
	/** metric -- disk : inodes Free */
	PHYSICAL_DISK_INODES_FREE = "physical.disk.inodesFree";
	/** metric -- disk : inodes Used Percent */
	PHYSICAL_DISK_INODES_USED_PERCENT = "physical.disk.inodesUsedPercent";

	/** metric -- net : up */
	PHYSICAL_NET_UP = "physical.net.up";
	/** metric -- net : down */
	PHYSICAL_NET_DOWN = "physical.net.down";
	/** metric -- net : count */
	PHYSICAL_NET_COUNT = "physical.net.count";
	/** metric -- net : estab */
	PHYSICAL_NET_ESTAB = "physical.net.estab";
	/** metric -- net : closeWait */
	PHYSICAL_NET_CLOSE_WAIT = "physical.net.closeWait";
	/** metric -- net : timeWait */
	PHYSICAL_NET_TIME_WAIT = "physical.net.timeWait";
	/** metric -- net : close */
	PHYSICAL_NET_CLOSE = "physical.net.close";
	/** metric -- net : listen */
	PHYSICAL_NET_LISTEN = "physical.net.listen";
	/** metric -- net : closing */
	PHYSICAL_NET_CLOSING = "physical.net.closing";
)
