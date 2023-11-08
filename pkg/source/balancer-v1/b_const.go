package balancerv1

import (
	"math/big"

	"github.com/KyberNetwork/blockchain-toolkit/integer"
)

var (
	BONE = integer.TenPow(18)

	MIN_BPOW_BASE  = integer.One()
	MAX_BPOW_BASE  = new(big.Int).Sub(new(big.Int).Mul(integer.Two(), BONE), integer.One())
	BPOW_PRECISION = new(big.Int).Div(BONE, integer.TenPow(10))

	MAX_IN_RATIO = new(big.Int).Div(BONE, integer.Two())
)
