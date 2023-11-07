package balancerv1

import (
	"errors"
	"math/big"

	"github.com/KyberNetwork/blockchain-toolkit/integer"
)

var (
	ErrDivZero         = errors.New("ERR_DIV_ZERO")
	ErrDivInternal     = errors.New("ERR_DIV_INTERNAL")
	ErrSubUnderflow    = errors.New("ERR_SUB_UNDERFLOW")
	ErrMulOverflow     = errors.New("ERR_MUL_OVERFLOW")
	ErrAddOverFlow     = errors.New("ERR_ADD_OVERFLOW")
	ErrBPowBaseTooLow  = errors.New("ERR_BPOW_BASE_TOO_LOW")
	ErrBPowBaseTooHigh = errors.New("ERR_BPOW_BASE_TOO_HIGH")
)

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L20
func btoi(a *big.Int) *big.Int {
	return new(big.Int).Div(a, BONE)
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L27
func bfloor(a *big.Int) *big.Int {
	return new(big.Int).Mul(btoi(a), BONE)
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L34
func badd(a *big.Int, b *big.Int) (*big.Int, error) {
	c := new(big.Int).Add(a, b)

	if c.Cmp(a) < 0 {
		return nil, ErrAddOverFlow
	}

	return c, nil
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L43
func bsub(a *big.Int, b *big.Int) (*big.Int, error) {
	c, flag := bsubSign(a, b)

	if flag {
		return nil, ErrSubUnderflow
	}

	return c, nil
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L52
func bsubSign(a *big.Int, b *big.Int) (*big.Int, bool) {
	if a.Cmp(b) >= 0 {
		return new(big.Int).Sub(a, b), false
	}

	return new(big.Int).Sub(b, a), true
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L63
func bmul(a *big.Int, b *big.Int) (*big.Int, error) {
	c0 := new(big.Int).Mul(a, b)

	if a.Cmp(integer.Zero()) != 0 && new(big.Int).Div(c0, a).Cmp(b) != 0 {
		return nil, ErrMulOverflow
	}

	c1 := new(big.Int).Add(c0, new(big.Int).Div(BONE, integer.Two()))

	if c1.Cmp(c0) < 0 {
		return nil, ErrMulOverflow
	}

	c2 := new(big.Int).Div(c1, BONE)

	return c2, nil
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L75
func bdiv(a *big.Int, b *big.Int) (*big.Int, error) {
	if b.Cmp(integer.Zero()) == 0 {
		return nil, ErrDivZero
	}

	c0 := new(big.Int).Mul(a, BONE)

	if a.Cmp(integer.Zero()) != 0 && new(big.Int).Div(c0, a).Cmp(BONE) != 0 {
		return nil, ErrDivInternal
	}

	c1 := new(big.Int).Add(c0, new(big.Int).Div(b, integer.Two()))

	if c1.Cmp(c0) < 0 {
		return nil, ErrDivInternal
	}

	c2 := new(big.Int).Div(c1, b)

	return c2, nil
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L89
func bpowi(a *big.Int, n *big.Int) (*big.Int, error) {
	var (
		z   *big.Int
		err error
	)

	if new(big.Int).Mod(n, integer.Two()).Cmp(integer.Zero()) != 0 {
		z = a
	} else {
		z = BONE
	}

	for n = new(big.Int).Div(n, integer.Two()); n.Cmp(integer.Zero()) != 0; n = new(big.Int).Div(n, integer.Two()) {
		a, err = bmul(a, a)
		if err != nil {
			return nil, err
		}

		if new(big.Int).Mod(n, integer.Two()).Cmp(integer.Zero()) != 0 {
			z, err = bmul(z, a)
			if err != nil {
				return nil, err
			}
		}
	}

	return z, nil
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L128C14-L128C24
func bpowApprox(base *big.Int, exp *big.Int, precision *big.Int) (*big.Int, error) {
	a := new(big.Int).Set(exp)
	x, xneg := bsubSign(base, BONE)
	term := new(big.Int).Set(BONE)
	sum := new(big.Int).Set(term)
	negative := false

	for i := integer.One(); term.Cmp(precision) >= 0; i = new(big.Int).Add(i, integer.One()) {
		bigK := new(big.Int).Mul(i, BONE)

		bsubBigKAndBone, err := bsub(bigK, BONE)
		if err != nil {
			return nil, err
		}

		c, cneg := bsubSign(a, bsubBigKAndBone)

		bmulCAndX, err := bmul(c, x)
		if err != nil {
			return nil, err
		}

		term, err := bmul(term, bmulCAndX)
		if err != nil {
			return nil, err
		}

		term, err = bdiv(term, bigK)
		if err != nil {
			return nil, err
		}

		if term.Cmp(integer.Zero()) == 0 {
			break
		}

		if xneg {
			negative = !negative
		}

		if cneg {
			negative = !negative
		}

		if negative {
			sum, err = bsub(sum, term)
			if err != nil {
				return nil, err
			}
		} else {
			sum, err = badd(sum, term)
			if err != nil {
				return nil, err
			}
		}
	}

	return sum, nil
}

// https://github.com/balancer/balancer-core/blob/f4ed5d65362a8d6cec21662fb6eae233b0babc1f/contracts/BNum.sol#L108
func bpow(base *big.Int, exp *big.Int) (*big.Int, error) {
	if base.Cmp(MIN_BPOW_BASE) < 0 {
		return nil, ErrBPowBaseTooLow
	}

	if base.Cmp(MAX_BPOW_BASE) > 0 {
		return nil, ErrBPowBaseTooHigh
	}

	whole := bfloor(exp)
	remain, err := bsub(exp, whole)
	if err != nil {
		return nil, err
	}

	wholePow, err := bpowi(base, btoi(whole))
	if err != nil {
		return nil, err
	}

	if remain.Cmp(integer.Zero()) == 0 {
		return wholePow, nil
	}

	partialResult, err := bpowApprox(base, remain, BPOW_PRECISION)
	if err != nil {
		return nil, err
	}

	return bmul(wholePow, partialResult)
}
