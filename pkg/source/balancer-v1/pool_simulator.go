package balancerv1

import (
	"errors"
	"math/big"

	"github.com/KyberNetwork/logger"

	poolpkg "github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool"
)

var (
	ErrNotBound      = errors.New("ERR_NOT_BOUND")
	ErrSwapNotPublic = errors.New("ERR_SWAP_NOT_PUBLIC")
	ErrMaxInRatio    = errors.New("ERR_MAX_IN_RATIO")
	ErrMathApprox    = errors.New("ERR_MATH_APPROX")
	//ErrBadLimitPrice = errors.New("ERR_BAD_LIMIT_PRICE")
	//ErrLimitOut      = errors.New("ERR_LIMIT_OUT")
	//ErrLimitPrice    = errors.New("ERR_LIMIT_PRICE")
)

type (
	PoolSimulator struct {
		poolpkg.Pool

		records    map[string]Record
		publicSwap bool
		swapFee    *big.Int

		gas Gas
	}

	PoolExtra struct {
		Records    map[string]Record `json:"records"`
		PublicSwap bool              `json:"publicSwap"`
		SwapFee    *big.Int          `json:"swapFee"`
	}

	Record struct {
		Bound   bool     `json:"bound"`
		Denorm  *big.Int `json:"denorm"`
		Balance *big.Int `json:"balance"`
	}

	PoolMeta struct {
		BlockNumber uint64
	}

	Gas struct {
		SwapExactAmountIn int64
	}
)

func (s *PoolSimulator) CalcAmountOut(tokenAmountIn poolpkg.TokenAmount, tokenOut string) (*poolpkg.CalcAmountOutResult, error) {
	amountOut, _, err := s.swapExactAmountIn(tokenAmountIn.Token, tokenAmountIn.Amount, tokenOut, nil, nil)
	if err != nil {
		return nil, err
	}

	return &poolpkg.CalcAmountOutResult{
		TokenAmountOut: &poolpkg.TokenAmount{Token: tokenOut, Amount: amountOut},
		Gas:            s.gas.SwapExactAmountIn,
	}, nil
}

func (s *PoolSimulator) UpdateBalance(params poolpkg.UpdateBalanceParams) {
	inRecord, outRecord := s.records[params.TokenAmountIn.Token], s.records[params.TokenAmountOut.Token]

	newBalanceIn, err := badd(inRecord.Balance, params.TokenAmountIn.Amount)
	if err != nil {
		logger.
			WithFields(logger.Fields{"poolAddress": s.GetAddress(), "err": err}).
			Warn("failed to update balance")
		return
	}

	newBalanceOut, err := bsub(outRecord.Balance, params.TokenAmountOut.Amount)
	if err != nil {
		logger.
			WithFields(logger.Fields{"poolAddress": s.GetAddress(), "err": err}).
			Warn("failed to update balance")
		return
	}

	inRecord.Balance = newBalanceIn
	outRecord.Balance = newBalanceOut

	s.records[params.TokenAmountIn.Token] = inRecord
	s.records[params.TokenAmountOut.Token] = outRecord
}

func (s *PoolSimulator) GetMetaInfo(_ string, _ string) interface{} {
	return PoolMeta{
		BlockNumber: s.Pool.Info.BlockNumber,
	}
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BPool.sol#L423
// NOTE: ignore minAmountOut and maxPrice because they are not necessary for our simulation
func (s *PoolSimulator) swapExactAmountIn(
	tokenIn string,
	tokenAmountIn *big.Int,
	tokenOut string,
	_ *big.Int, // minAmountOut
	_ *big.Int, // maxPrice
) (*big.Int, *big.Int, error) {
	if !s.records[tokenIn].Bound {
		return nil, nil, ErrNotBound
	}

	if !s.records[tokenOut].Bound {
		return nil, nil, ErrNotBound
	}

	if !s.publicSwap {
		return nil, nil, ErrSwapNotPublic
	}

	inRecord, outRecord := s.records[tokenIn], s.records[tokenOut]

	bmulBalanceInAndMaxIn, err := bmul(inRecord.Balance, MAX_IN_RATIO)
	if err != nil {
		return nil, nil, err
	}

	if tokenAmountIn.Cmp(bmulBalanceInAndMaxIn) > 0 {
		return nil, nil, ErrMaxInRatio
	}

	spotPriceBefore, err := calcSpotPrice(
		inRecord.Balance,
		inRecord.Denorm,
		outRecord.Balance,
		outRecord.Denorm,
		s.swapFee,
	)
	if err != nil {
		return nil, nil, err
	}

	//if spotPriceBefore.Cmp(maxPrice) > 0 {
	//	return nil, nil, ErrBadLimitPrice
	//}

	tokenAmountOut, err := calcOutGivenIn(
		inRecord.Balance,
		inRecord.Denorm,
		outRecord.Balance,
		outRecord.Denorm,
		tokenAmountIn,
		s.swapFee,
	)
	if err != nil {
		return nil, nil, err
	}

	//if tokenAmountOut.Cmp(minAmountOut) < 0 {
	//	return nil, nil, ErrLimitOut
	//}

	inRecord.Balance, err = badd(inRecord.Balance, tokenAmountIn)
	if err != nil {
		return nil, nil, err
	}

	outRecord.Balance, err = bsub(outRecord.Balance, tokenAmountOut)
	if err != nil {
		return nil, nil, err
	}

	spotPriceAfter, err := calcSpotPrice(
		inRecord.Balance,
		inRecord.Denorm,
		outRecord.Balance,
		outRecord.Denorm,
		s.swapFee,
	)
	if err != nil {
		return nil, nil, err
	}

	if spotPriceAfter.Cmp(spotPriceBefore) < 0 {
		return nil, nil, ErrMathApprox
	}

	//if spotPriceAfter.Cmp(maxPrice) > 0 {
	//	return nil, nil, ErrLimitPrice
	//}

	bdivTokenAmountInAndOut, err := bdiv(tokenAmountIn, tokenAmountOut)
	if err != nil {
		return nil, nil, err
	}

	if spotPriceBefore.Cmp(bdivTokenAmountInAndOut) > 0 {
		return nil, nil, ErrMathApprox
	}

	return tokenAmountOut, spotPriceAfter, nil
}
