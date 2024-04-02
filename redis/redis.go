package redis

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Redis struct {
	Store     store.Db
	repl_info ReplicationInfo
}

type ReplicationInfo struct {
	role string
	// "master" - if the instance is replica of noone
	// "slave"  - if the intsance is a replica of some master instances

	connected_slaves int32
	// No of connected Replicas

	master_replid string
	// The replication ID of the Redis server.

	master_repl_offset int32
	// The server's current replication offset

	second_repl_offset int32
	// The offset up to which replication IDs are accepted

	repl_backlog_active bool
	// Flag indicating replication backlog is active

	repl_backlog_size int64
	// Total size in bytes of the replication backlog buffer

	repl_backlog_first_byte_offset int32
	// The master offset of the replication backlog buffer

	repl_backlog_hlisten int32
	// Size in bytes of the data in the replication backlog buffer

}

func NewRedisMaster() *Redis {
	return &Redis{
		Store: *store.NewDb(),
		repl_info: ReplicationInfo{
			role: "master",
		},
	}
}

func NewRedisSlave() *Redis {
	return &Redis{
		Store: *store.NewDb(),
		repl_info: ReplicationInfo{
			role: "slave",
		},
	}
}

func (redis *Redis) GetInfo() resp.Value {
	s := fmt.Sprintf("role:%s", redis.repl_info.role)

	return resp.Value{
		Typ:  resp.BulkStringType,
		Bulk: s,
	}
}
