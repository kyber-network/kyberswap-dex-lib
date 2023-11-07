package balancerv1

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	poolpkg "github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool"
	utils "github.com/KyberNetwork/kyberswap-dex-lib/pkg/util/bignumber"
)

func TestPoolSimulator_CalcAmountOut(t *testing.T) {
	testCases := []struct {
		name              string
		poolSimulator     PoolSimulator
		tokenAmountIn     poolpkg.TokenAmount
		tokenOut          string
		expectedAmountOut *big.Int
		expectedError     error
	}{
		{
			name: "it should return correct amountOut",
			poolSimulator: PoolSimulator{
				records: map[string]Record{
					"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": {
						Bound:   true,
						Balance: utils.NewBig("181453339134494385762"),
						Denorm:  utils.NewBig("25000000000000000000"),
					},
					"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599": {
						Bound:   true,
						Balance: utils.NewBig("982184296"),
						Denorm:  utils.NewBig("25000000000000000000"),
					},
				},
				publicSwap: true,
				swapFee:    utils.NewBig("4000000000000000"),
			},
			tokenAmountIn:     poolpkg.TokenAmount{Token: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", Amount: utils.NewBig("81275824825923290")},
			tokenOut:          "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
			expectedAmountOut: utils.NewBig("437981"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.poolSimulator.CalcAmountOut(tc.tokenAmountIn, tc.tokenOut)

			assert.ErrorIs(t, tc.expectedError, err)
			if tc.expectedAmountOut != nil {
				assert.Equal(t, 0, tc.expectedAmountOut.Cmp(result.TokenAmountOut.Amount))
			}
		})
	}
}

func TestPoolSimulator_UpdateBalance(t *testing.T) {
	testCases := []struct {
		name               string
		poolSimulator      PoolSimulator
		params             poolpkg.UpdateBalanceParams
		expectedBalanceIn  *big.Int
		expectedBalanceOut *big.Int
	}{
		{
			name: "it should return correct amountOut",
			poolSimulator: PoolSimulator{
				records: map[string]Record{
					"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": {
						Bound:   true,
						Balance: utils.NewBig("181453339134494385762"),
						Denorm:  utils.NewBig("25000000000000000000"),
					},
					"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599": {
						Bound:   true,
						Balance: utils.NewBig("982184296"),
						Denorm:  utils.NewBig("25000000000000000000"),
					},
				},
				publicSwap: true,
				swapFee:    utils.NewBig("4000000000000000"),
			},
			params: poolpkg.UpdateBalanceParams{
				TokenAmountIn:  poolpkg.TokenAmount{Token: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", Amount: utils.NewBig("81275824825923290")},
				TokenAmountOut: poolpkg.TokenAmount{Token: "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", Amount: utils.NewBig("437981")},
			},
			expectedBalanceIn:  utils.NewBig("181534614959320309052"),
			expectedBalanceOut: utils.NewBig("981746315"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.poolSimulator.UpdateBalance(tc.params)

			assert.Equal(t, 0, tc.expectedBalanceIn.Cmp(tc.poolSimulator.records[tc.params.TokenAmountIn.Token].Balance))
			assert.Equal(t, 0, tc.expectedBalanceOut.Cmp(tc.poolSimulator.records[tc.params.TokenAmountOut.Token].Balance))
		})
	}
}
