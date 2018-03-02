package comon

import (
	"github.com/onrik/ethrpc"
	"math/big"
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
	GasPrice     big.Int
	UpTime       int
	Err          error
}
