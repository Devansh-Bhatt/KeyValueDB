package repl

import (
	"encoding/hex"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const (
	hexcontent = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
)

type ReplicationInfo struct {
	Role string
	// "master" - if the instance is replica of noone
	// "slave"  - if the intsance is a replica of some master instances

	Connected_slaves int32
	// No of connected Replicas

	Master_replid string
	// The replication ID of the Redis server.

	Master_repl_offset int32
	// The server's current replication offset

	Second_repl_offset int32
	// The offset up to which replication IDs are accepted

	Repl_backlog_active bool
	// Flag indicating replication backlog is active

	Repl_backlog_size int64
	// Total size in bytes of the replication backlog buffer

	Repl_backlog_first_byte_offset int32
	// The master offset of the replication backlog buffer

	Repl_backlog_hlisten int32
	// Size in bytes of the data in the replication backlog buffer

}

func (Ri *ReplicationInfo) GetInfo() resp.Value {
	s := fmt.Sprintf("role:%s\n master_replid:%s\n master_repl_offset:%d", Ri.Role, Ri.Master_replid, Ri.Master_repl_offset)

	return resp.Value{
		Typ:  resp.BulkStringType,
		Bulk: s,
	}
}

func (Ri *ReplicationInfo) FullResync() resp.Value {
	HexContent := hexcontent
	Bin, err := hex.DecodeString(HexContent)

	if err != nil {
		fmt.Println(err)
	}
	return resp.Value{
		Typ:   resp.RDBType,
		Bytes: Bin,
	}
}
