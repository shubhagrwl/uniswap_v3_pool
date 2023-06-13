package posts

import (
	"time"
)

const (
	TABLE_NAME            = "pool_logs"
	COLUMN_ID             = "id"
	COLUMN_POOL_ADDRESS   = "pool_address"
	COLUMN_TXN_ID         = "txn_id"
	COLUMN_BLOCK_NUMBER   = "block_number"
	COLUMN_TOKEN0_BALANCE = "token0_balance"
	COLUMN_TOKEN1_BALANCE = "token1_balance"
	COLUMN_TOKEN0_DELTA   = "token0_delta"
	COLUMN_TOKEN1_DELTA   = "token1_delta"
	COLUMN_TICK           = "tick"
	COLUMN_CREATED_AT     = "created_at"
)

type Logs struct {
	Id            int       `json:"id"`
	PoolAddress   string    `json:"pool_address"`
	TxnId         string    `json:"txn_id"`
	BlockNumber   uint64    `json:"block_number"`
	Token0Balance int64     `json:"token0_balance"`
	Token1Balance int64     `json:"token1_balance"`
	Token0Delta   int64     `json:"token0_delta"`
	Token1Delta   int64     `json:"token1_delta"`
	Tick          int64     `json:"tick"`
	CreatedAt     time.Time `json:"created_at"`
}
