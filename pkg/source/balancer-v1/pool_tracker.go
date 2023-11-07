package balancerv1

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"time"

	"github.com/KyberNetwork/ethrpc"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool"
)

type (
	PoolTracker struct {
		ethrpcClient *ethrpc.Client
	}

	PoolData struct {
		Tokens       []string
		SwapFee      *big.Int
		Records      map[string]Record
		IsPublicSwap bool
		BlockNumber  uint64
	}
)

func NewPoolTracker(
	ethrpcClient *ethrpc.Client,
) (*PoolTracker, error) {
	return &PoolTracker{
		ethrpcClient: ethrpcClient,
	}, nil
}

func (t *PoolTracker) GetNewPoolState(
	ctx context.Context,
	p entity.Pool,
	params pool.GetNewPoolStateParams,
) (entity.Pool, error) {
	startTime := time.Now()
	logger.WithFields(logger.Fields{"pool_id": p.Address}).Info("Started getting new pool state")
	defer func() {
		logger.
			WithFields(
				logger.Fields{
					"pool_id":     p.Address,
					"duration_ms": time.Since(startTime).Milliseconds(),
				},
			).
			Info("Finished getting new pool state")
	}()

	poolData, err := t.getPoolData(ctx, p.Address)
	if err != nil {
		return p, err
	}

	p, err = t.updatePool(p, poolData)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (t *PoolTracker) getPoolData(ctx context.Context, address string) (PoolData, error) {
	var (
		tokenAddresses []common.Address
		swapFee        *big.Int
		isPublicSwap   bool
		blockNumber    uint64
	)

	getPoolRequest := t.ethrpcClient.NewRequest().SetContext(ctx)

	getPoolRequest.AddCall(&ethrpc.Call{
		ABI:    bPoolABI,
		Target: address,
		Method: bPoolMethodGetCurrentTokens,
		Params: nil,
	}, []interface{}{&tokenAddresses})
	getPoolRequest.AddCall(&ethrpc.Call{
		ABI:    bPoolABI,
		Target: address,
		Method: bPoolMethodGetSwapFee,
		Params: nil,
	}, []interface{}{&swapFee})
	getPoolRequest.AddCall(&ethrpc.Call{
		ABI:    bPoolABI,
		Target: address,
		Method: bPoolMethodIsPublicSwap,
		Params: nil,
	}, []interface{}{&isPublicSwap})

	resp, err := getPoolRequest.TryBlockAndAggregate()
	if err != nil {
		return PoolData{}, err
	}

	blockNumber = resp.BlockNumber.Uint64()

	tokensLen := len(tokenAddresses)
	boundList := make([]bool, tokensLen)
	balanceList := make([]*big.Int, tokensLen)
	denormList := make([]*big.Int, tokensLen)

	getPoolRecordsRequest := t.ethrpcClient.NewRequest().SetContext(ctx).SetBlockNumber(resp.BlockNumber)
	for i, token := range tokenAddresses {
		getPoolRecordsRequest.AddCall(&ethrpc.Call{
			ABI:    bPoolABI,
			Target: address,
			Method: bPoolMethodIsBound,
			Params: []interface{}{token},
		}, []interface{}{&boundList[i]})
		getPoolRecordsRequest.AddCall(&ethrpc.Call{
			ABI:    bPoolABI,
			Target: address,
			Method: bPoolMethodGetBalance,
			Params: []interface{}{token},
		}, []interface{}{&balanceList[i]})
		getPoolRecordsRequest.AddCall(&ethrpc.Call{
			ABI:    bPoolABI,
			Target: address,
			Method: bPoolMethodGetDenormalizedWeight,
			Params: []interface{}{token},
		}, []interface{}{&denormList[i]})
	}

	resp, err = getPoolRecordsRequest.TryBlockAndAggregate()
	if err != nil {
		return PoolData{}, err
	}

	tokens := make([]string, 0, tokensLen)
	records := make(map[string]Record, tokensLen)
	for i, token := range tokenAddresses {
		tokenAddressStr := strings.ToLower(token.String())

		records[tokenAddressStr] = Record{
			Bound:   boundList[i],
			Balance: balanceList[i],
			Denorm:  balanceList[i],
		}
		tokens = append(tokens, tokenAddressStr)
	}

	return PoolData{
		Tokens:       tokens,
		SwapFee:      swapFee,
		IsPublicSwap: isPublicSwap,
		Records:      records,
		BlockNumber:  blockNumber,
	}, nil
}

func (t *PoolTracker) updatePool(p entity.Pool, poolData PoolData) (entity.Pool, error) {
	extra := PoolExtra{
		Records:    poolData.Records,
		PublicSwap: poolData.IsPublicSwap,
		SwapFee:    poolData.SwapFee,
	}
	extraBytes, err := json.Marshal(extra)
	if err != nil {
		return p, err
	}

	poolTokens := make([]*entity.PoolToken, 0, len(poolData.Tokens))
	reserves := make([]string, 0, len(poolData.Tokens))
	for _, token := range poolData.Tokens {
		poolTokens = append(poolTokens, &entity.PoolToken{Address: token, Swappable: true})
		reserves = append(reserves, poolData.Records[token].Balance.String())
	}

	p.Tokens = poolTokens
	p.Reserves = reserves
	p.Extra = string(extraBytes)
	p.BlockNumber = poolData.BlockNumber
	p.Timestamp = time.Now().Unix()

	return p, nil
}
