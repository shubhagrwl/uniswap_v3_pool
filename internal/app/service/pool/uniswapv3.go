package pool

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"uniswapper/internal/app/constants"
	posts "uniswapper/internal/app/db/dto/pool"
	"uniswapper/internal/app/db/repository/pool"
	"uniswapper/internal/app/service/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IUniswapV3Pool interface {
	RunUniswapV3Pool(ctx context.Context)
}

type UniswapV3Pool struct {
	client           *ethclient.Client
	addresses        []common.Address
	PoolLogsDBClient pool.IPoolLogsRepository
}

func NewUniswapV3Pool(ctx context.Context, poolLogsDBClient pool.IPoolLogsRepository) IUniswapV3Pool {
	log := logger.Logger(ctx)
	client, err := ethclient.Dial(constants.Config.PoolConfig.INFURA_MAINNET)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var poolAddresses []string

	// Use json.Unmarshal to convert the string to a []string
	err = json.Unmarshal([]byte(constants.Config.PoolConfig.POOL_ADDRESSES), &poolAddresses)
	if err != nil || len(poolAddresses) == 0 {
		log.Fatalf("Error while reading pool addresses")
	}

	var addresses []common.Address
	for _, addr := range poolAddresses {
		addresses = append(addresses, common.HexToAddress(addr))
	}

	return &UniswapV3Pool{client: client, addresses: addresses, PoolLogsDBClient: poolLogsDBClient}
}

func (u UniswapV3Pool) RunUniswapV3Pool(ctx context.Context) {
	log := logger.Logger(ctx)
	query := ethereum.FilterQuery{
		Addresses: u.addresses,
	}

	logs := make(chan types.Log)
	sub, err := u.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	uniswapPoolABI := `[
		{
		  "inputs": [],
		  "stateMutability": "nonpayable",
		  "type": "constructor"
		},
		{
		"anonymous": false,
		"inputs": [
		  {
			"indexed": true,
			"internalType": "uint128",
			"name": "token0Balance",
			"type": "uint128"
		  },
		  {
			"indexed": true,
			"internalType": "uint128",
			"name": "token1Balance",
			"type": "uint128"
		  },
		  {
			"indexed": false,
			"internalType": "int24",
			"name": "tick",
			"type": "int24"
		  },
		  {
			"indexed": false,
			"internalType": "uint256",
			"name": "token0Delta",
			"type": "uint256"
		  },
		  {
			"indexed": false,
			"internalType": "uint256",
			"name": "token1Delta",
			"type": "uint256"
		  }
		],
		"name": "PoolUpdate",
		"type": "event"
	  },
		{
		  "anonymous": false,
		  "inputs": [
			{
			  "indexed": true,
			  "internalType": "uint24",
			  "name": "fee",
			  "type": "uint24"
			},
			{
			  "indexed": true,
			  "internalType": "int24",
			  "name": "tickSpacing",
			  "type": "int24"
			}
		  ],
		  "name": "FeeAmountEnabled",
		  "type": "event"
		},
		{
		  "anonymous": false,
		  "inputs": [
			{
			  "indexed": true,
			  "internalType": "address",
			  "name": "oldOwner",
			  "type": "address"
			},
			{
			  "indexed": true,
			  "internalType": "address",
			  "name": "newOwner",
			  "type": "address"
			}
		  ],
		  "name": "OwnerChanged",
		  "type": "event"
		},
		{
		  "anonymous": false,
		  "inputs": [
			{
			  "indexed": true,
			  "internalType": "address",
			  "name": "token0",
			  "type": "address"
			},
			{
			  "indexed": true,
			  "internalType": "address",
			  "name": "token1",
			  "type": "address"
			},
			{
			  "indexed": true,
			  "internalType": "uint24",
			  "name": "fee",
			  "type": "uint24"
			},
			{
			  "indexed": false,
			  "internalType": "int24",
			  "name": "tickSpacing",
			  "type": "int24"
			},
			{
			  "indexed": false,
			  "internalType": "address",
			  "name": "pool",
			  "type": "address"
			}
		  ],
		  "name": "PoolCreated",
		  "type": "event"
		},
		{
		  "inputs": [
			{
			  "internalType": "address",
			  "name": "tokenA",
			  "type": "address"
			},
			{
			  "internalType": "address",
			  "name": "tokenB",
			  "type": "address"
			},
			{
			  "internalType": "uint24",
			  "name": "fee",
			  "type": "uint24"
			}
		  ],
		  "name": "createPool",
		  "outputs": [
			{
			  "internalType": "address",
			  "name": "pool",
			  "type": "address"
			}
		  ],
		  "stateMutability": "nonpayable",
		  "type": "function"
		},
		{
		  "inputs": [
			{
			  "internalType": "uint24",
			  "name": "fee",
			  "type": "uint24"
			},
			{
			  "internalType": "int24",
			  "name": "tickSpacing",
			  "type": "int24"
			}
		  ],
		  "name": "enableFeeAmount",
		  "outputs": [],
		  "stateMutability": "nonpayable",
		  "type": "function"
		},
		{
		  "inputs": [
			{
			  "internalType": "uint24",
			  "name": "",
			  "type": "uint24"
			}
		  ],
		  "name": "feeAmountTickSpacing",
		  "outputs": [
			{
			  "internalType": "int24",
			  "name": "",
			  "type": "int24"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		},
		{
		  "inputs": [
			{
			  "internalType": "address",
			  "name": "",
			  "type": "address"
			},
			{
			  "internalType": "address",
			  "name": "",
			  "type": "address"
			},
			{
			  "internalType": "uint24",
			  "name": "",
			  "type": "uint24"
			}
		  ],
		  "name": "getPool",
		  "outputs": [
			{
			  "internalType": "address",
			  "name": "",
			  "type": "address"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		},
		{
		  "inputs": [],
		  "name": "owner",
		  "outputs": [
			{
			  "internalType": "address",
			  "name": "",
			  "type": "address"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		},
		{
		  "inputs": [],
		  "name": "parameters",
		  "outputs": [
			{
			  "internalType": "address",
			  "name": "factory",
			  "type": "address"
			},
			{
			  "internalType": "address",
			  "name": "token0",
			  "type": "address"
			},
			{
			  "internalType": "address",
			  "name": "token1",
			  "type": "address"
			},
			{
			  "internalType": "uint24",
			  "name": "fee",
			  "type": "uint24"
			},
			{
			  "internalType": "int24",
			  "name": "tickSpacing",
			  "type": "int24"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		},
		{
		  "inputs": [
			{
			  "internalType": "address",
			  "name": "_owner",
			  "type": "address"
			}
		  ],
		  "name": "setOwner",
		  "outputs": [],
		  "stateMutability": "nonpayable",
		  "type": "function"
		}
	  ]`
	contractAbi, err := abi.JSON(strings.NewReader(uniswapPoolABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}
	var count uint64 = 0

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatalf("Received subscription error: %v", err)
			case vLog := <-logs:
				count++
				if count%12 == 0 {
					event := struct {
						Token0Balance *big.Int
						Token1Balance *big.Int
						Tick          *big.Int
						Token0Delta   *big.Int
						Token1Delta   *big.Int
					}{}

					err := contractAbi.UnpackIntoInterface(&event, "PoolUpdate", vLog.Data)
					if err != nil {
						log.Fatalf("Failed to unpack data into PoolUpdate event: %v", err)
					}

					var token0Bal, token1Bal *big.Int
					if len(vLog.Topics) > 1 {
						token0Bal = vLog.Topics[1].Big()
					}

					if len(vLog.Topics) > 2 {
						token1Bal = vLog.Topics[2].Big()
					}

					logs := posts.Logs{
						PoolAddress:   vLog.Address.String(),
						TxnId:         vLog.TxHash.String(),
						BlockNumber:   vLog.BlockNumber,
						Token0Balance: token0Bal.Int64(),
						Token1Balance: token1Bal.Int64(),
						Token0Delta:   event.Token0Delta.Int64(),
						Token1Delta:   event.Token1Delta.Int64(),
						Tick:          event.Tick.Int64(),
					}

					log.Info("Get Block Info", logs)

					if err := u.PoolLogsDBClient.StorePoolLogs(ctx, logs); err != nil {
						log.Error("error while storing logs")
					}
				}
			}
		}
	}()

	// Keep the program running indefinitely
	select {}
}
