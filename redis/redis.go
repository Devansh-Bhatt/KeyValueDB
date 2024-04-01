package redis

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Redis struct {
	store     store.Db
	repl_info ReplicationInfostruct
}

type ReplicationInfostruct struct {
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

func NewRedis() *Redis {
	return &Redis{
		store: *store.NewDb(),
		repl_info: ReplicationInfostruct{
			role: "master",
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
