package balancerv1

import "math/big"

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BMath.sol#L28
func calcSpotPrice(
	tokenBalanceIn *big.Int,
	tokenWeightIn *big.Int,
	tokenBalanceOut *big.Int,
	tokenWeightOut *big.Int,
	swapFee *big.Int,
) (*big.Int, error) {
	numer, err := bdiv(tokenBalanceIn, tokenWeightIn)
	if err != nil {
		return nil, err
	}

	denom, err := bdiv(tokenBalanceOut, tokenWeightOut)
	if err != nil {
		return nil, err
	}

	ratio, err := bdiv(numer, denom)
	if err != nil {
		return nil, err
	}

	bsubBONEAndSwapFee, err := bsub(BONE, swapFee)
	if err != nil {
		return nil, err
	}

	scale, err := bdiv(BONE, bsubBONEAndSwapFee)
	if err != nil {
		return nil, err
	}

	return bmul(ratio, scale)
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BMath.sol#L55
func calcOutGivenIn(
	tokenBalanceIn *big.Int,
	tokenWeightIn *big.Int,
	tokenBalanceOut *big.Int,
	tokenWeightOut *big.Int,
	tokenAmountIn *big.Int,
	swapFee *big.Int,
) (*big.Int, error) {
	weightRatio, err := bdiv(tokenWeightIn, tokenWeightOut)
	if err != nil {
		return nil, err
	}

	adjustedIn, err := bsub(BONE, swapFee)
	if err != nil {
		return nil, err
	}

	adjustedIn, err = bmul(tokenAmountIn, adjustedIn)
	if err != nil {
		return nil, err
	}

	baddTokenBalanceInAndAdjustedIn, err := badd(tokenBalanceIn, adjustedIn)
	if err != nil {
		return nil, err
	}

	y, err := bdiv(tokenBalanceIn, baddTokenBalanceInAndAdjustedIn)
	if err != nil {
		return nil, err
	}

	foo, err := bpow(y, weightRatio)
	if err != nil {
		return nil, err
	}

	bar, err := bsub(BONE, foo)
	if err != nil {
		return nil, err
	}

	return bmul(tokenBalanceOut, bar)
}
