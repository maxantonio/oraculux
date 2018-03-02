package comon

import (
	"github.com/onrik/ethrpc"
)

type ServerInfo struct {
	Server       string
	Sincing      *ethrpc.Syncing
	Block        *ethrpc.Block
	Pending      *PoolStatus
	BlockNumber  int
	Peers        int
	IsMining     bool
	Transactions int
	LocalPending int
	Ping         string
	Latency      string
	HashRate     int
	Uncles       int
	GasPrice     int64
	UpTime       int
	Err          error
}
