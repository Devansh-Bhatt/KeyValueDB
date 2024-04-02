package redis

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
	"github.com/codecrafters-io/redis-starter-go/util"
)

type Redis struct {
	Store     store.Db
	repl_info ReplicationInfo
}

type RedisSlave struct {
	Redis
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
			role:               "master",
			master_replid:      util.Randomalphanumericgenerator(40),
			master_repl_offset: 0,
		},
	}
}

func NewRedisSlave() *RedisSlave {
	return &RedisSlave{
		Redis{Store: *store.NewDb(),
			repl_info: ReplicationInfo{
				role:               "slave",
				master_replid:      util.Randomalphanumericgenerator(40),
				master_repl_offset: 0,
			}},
	}
}

func (rs *RedisSlave) ConnectMaster(Master string, port string) (net.Conn, error) {
	MasterAddr := fmt.Sprintf("%s:%s", Master, port)

	conn, err := net.Dial("tcp", MasterAddr)

	if err != nil {
		fmt.Println("couldnt not connect to master")
	}

	return conn, nil
}

func (redis *Redis) GetInfo() resp.Value {
	s := fmt.Sprintf("role:%s\n master_replid:%s\n master_repl_offset:%d", redis.repl_info.role, redis.repl_info.master_replid, redis.repl_info.master_repl_offset)

	return resp.Value{
		Typ:  resp.BulkStringType,
		Bulk: s,
	}
}
