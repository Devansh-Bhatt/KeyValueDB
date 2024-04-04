package redis

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
	"github.com/codecrafters-io/redis-starter-go/util"
)

type Redis struct {
	Store     store.Db
	Repl_info ReplicationInfo
	IsSlave   bool
}

type ReplicationInfo struct {
	role string
	// "master" - if the instance is replica of noone
	// "slave"  - if the intsance is a replica of some master instances

	connected_slaves int32
	// No of connected Replicas

	Master_replid string
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

func NewRedisServer(isSlave bool) *Redis {
	var role string
	if isSlave {
		role = "slave"
	} else {
		role = "master"
	}

	return &Redis{
		Store: *store.NewDb(),
		Repl_info: ReplicationInfo{
			role:               role,
			Master_replid:      util.Randomalphanumericgenerator(40),
			master_repl_offset: 0,
		},
	}
}

func (redis *Redis) GetInfo() resp.Value {
	s := fmt.Sprintf("role:%s\n master_replid:%s\n master_repl_offset:%d", redis.Repl_info.role, redis.Repl_info.Master_replid, redis.Repl_info.master_repl_offset)

	return resp.Value{
		Typ:  resp.BulkStringType,
		Bulk: s,
	}
}
