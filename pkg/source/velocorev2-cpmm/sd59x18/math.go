package sd59x18

import (
	"math/big"
)

func Log2(x SD59x18) (SD59x18, error) {
	xBI := new(big.Int).Set(x)

	if xBI.Cmp(zeroBI) < 0 {
		return nil, ErrMathSD59x18LogInputTooSmall
	}

	var sign *big.Int
	if xBI.Cmp(uUnit) >= 0 {
		sign = big.NewInt(1)
	} else {
		sign = big.NewInt(-1)
		xBI = new(big.Int).Div(uUnitSquared, xBI)
	}

	n := msb(new(big.Int).Div(xBI, uUnit))

	resultBI := new(big.Int).Mul(n, uUnit)

	y := new(big.Int).Rsh(xBI, uint(n.Uint64()))

	if y.Cmp(uUnit) == 0 {
		return new(big.Int).Mul(resultBI, sign), nil
	}

	doubleUnit := new(big.Int).Mul(uUnit, big.NewInt(2))
	for delta := new(big.Int).Set(uHalfUnit); delta.Cmp(zeroBI) > 0; delta = new(big.Int).Rsh(delta, 1) { // TODO: does this change vallue of uHalfUnit?
		y = new(big.Int).Div(new(big.Int).Mul(y, y), uUnit)

		if y.Cmp(doubleUnit) >= 0 {
			resultBI = new(big.Int).Add(resultBI, delta)

			y = new(big.Int).Rsh(y, 1)
		}
	}
	resultBI = new(big.Int).Mul(resultBI, sign)
	return resultBI, nil
}

func msb(x *big.Int) *big.Int {
	l := 0
	r := 256
	twoBI := big.NewInt(2)
	for r-l > 1 {
		m := (l + r) >> 1
		twoPowM := new(big.Int).Exp(twoBI, big.NewInt(int64(m)), nil)
		if twoPowM.Cmp(x) <= 0 {
			l = m
		} else {
			r = m
		}
	}
	return big.NewInt(int64(l))
}

func Exp2(x SD59x18) (SD59x18, error) {
	xBI := new(big.Int).Set(x)
	if xBI.Cmp(zeroBI) < 0 {
		magicNbr, _ := new(big.Int).SetString("-59794705707972522261", 10)
		if xBI.Cmp(magicNbr) < 0 {
			return zeroBI, nil
		}
		xBI := new(big.Int).Mul(xBI, big.NewInt(-1))
		result, err := Exp2(xBI)
		if err != nil {
			return nil, err
		}
		return new(big.Int).Div(uUnitSquared, result), nil
	}

	if xBI.Cmp(uExp2MaxInput) > 0 {
		return nil, ErrMathSD59x18Exp2InputTooBig
	}

	xType192x64 := new(big.Int).Div(new(big.Int).Lsh(xBI, 64), uUnit)
	return exp2(xType192x64), nil
}

func exp2(x *big.Int) *big.Int {
	// Start from 0.5 in the 192.64-bit fixed-point format.
	result, _ := new(big.Int).SetString("800000000000000000000000000000000000000000000000", 16)

	const n = 8

	v := [n]string{
		"FF00000000000000",
		"FF000000000000",
		"FF0000000000",
		"FF00000000",
		"FF000000",
		"FF0000",
		"FF00",
		"FF",
	}
	type pair struct {
		X string
		Y string
	}
	pairs := [n][n]pair{
		{
			{"8000000000000000", "16A09E667F3BCC909"},
			{"4000000000000000", "1306FE0A31B7152DF"},
			{"2000000000000000", "1172B83C7D517ADCE"},
			{"1000000000000000", "10B5586CF9890F62A"},
			{"800000000000000", "1059B0D31585743AE"},
			{"400000000000000", "102C9A3E778060EE7"},
			{"200000000000000", "10163DA9FB33356D8"},
			{"100000000000000", "100B1AFA5ABCBED61"},
		},
		{
			{"80000000000000", "10058C86DA1C09EA2"},
			{"40000000000000", "1002C605E2E8CEC50"},
			{"20000000000000", "100162F3904051FA1"},
			{"10000000000000", "1000B175EFFDC76BA"},
			{"8000000000000", "100058BA01FB9F96D"},
			{"4000000000000", "10002C5CC37DA9492"},
			{"2000000000000", "1000162E525EE0547"},
			{"1000000000000", "10000B17255775C04"},
		},
		{
			{"800000000000", "1000058B91B5BC9AE"},
			{"400000000000", "100002C5C89D5EC6D"},
			{"200000000000", "10000162E43F4F831"},
			{"100000000000", "100000B1721BCFC9A"},
			{"80000000000", "10000058B90CF1E6E"},
			{"40000000000", "1000002C5C863B73F"},
			{"20000000000", "100000162E430E5A2"},
			{"10000000000", "1000000B172183551"},
		},
		{
			{"8000000000", "100000058B90C0B49"},
			{"4000000000", "10000002C5C8601CC"},
			{"2000000000", "1000000162E42FFF0"},
			{"1000000000", "10000000B17217FBB"},
			{"800000000", "1000000058B90BFCE"},
			{"400000000", "100000002C5C85FE3"},
			{"200000000", "10000000162E42FF1"},
			{"100000000", "100000000B17217F8"},
		},
		{
			{"80000000", "10000000058B90BFC"},
			{"40000000", "1000000002C5C85FE"},
			{"20000000", "100000000162E42FF"},
			{"10000000", "1000000000B17217F"},
			{"8000000", "100000000058B90C0"},
			{"4000000", "10000000002C5C860"},
			{"2000000", "1000000000162E430"},
			{"1000000", "10000000000B17218"},
		},
		{
			{"800000", "1000000000058B90C"},
			{"400000", "100000000002C5C86"},
			{"200000", "10000000000162E43"},
			{"100000", "100000000000B1721"},
			{"80000", "10000000000058B91"},
			{"40000", "1000000000002C5C8"},
			{"20000", "100000000000162E4"},
			{"10000", "1000000000000B172"},
		},
		{
			{"8000", "100000000000058B9"},
			{"4000", "10000000000002C5D"},
			{"2000", "1000000000000162E"},
			{"1000", "10000000000000B17"},
			{"800", "1000000000000058C"},
			{"400", "100000000000002C6"},
			{"200", "10000000000000163"},
			{"100", "100000000000000B1"},
		},
		{
			{"80", "10000000000000059"},
			{"40", "1000000000000002C"},
			{"20", "10000000000000016"},
			{"10", "1000000000000000B"},
			{"8", "10000000000000006"},
			{"4", "10000000000000003"},
			{"2", "10000000000000001"},
			{"1", "10000000000000001"},
		},
	}

	for i := 0; i < n; i++ {
		vi, ok := new(big.Int).SetString(v[i], 16)
		if !ok {
			// TODO: handle error
			panic("failed to parse hex string")
		}

		if new(big.Int).And(x, vi).Cmp(zeroBI) <= 0 {
			continue
		}

		for j := 0; j < n; j++ {
			xj, ok := new(big.Int).SetString(pairs[i][j].X, 16)
			if !ok {
				panic("failed to parse hex string xj")
				// TODO: handle error
			}

			if new(big.Int).And(x, xj).Cmp(zeroBI) <= 0 {
				continue
			}

			yj, ok := new(big.Int).SetString(pairs[i][j].Y, 16)
			if !ok {
				panic("failed to parse hex string yj")
				// TODO: handle err
			}

			result = new(big.Int).Rsh(new(big.Int).Mul(result, yj), 64)
		}
	}

	result = new(big.Int).Mul(result, uUnit)
	result = new(big.Int).Rsh(
		result,
		uint(new(big.Int).Sub(
			big.NewInt(191),
			new(big.Int).Rsh(x, 64),
		).Uint64()),
	)

	return result
}

func Pow(x SD59x18, y SD59x18) (SD59x18, error) {
	var (
		xBI *big.Int = x
		yBI *big.Int = y
	)

	if xBI.Cmp(zeroBI) == 0 {
		ret := Zero()
		if yBI.Cmp(zeroBI) == 0 {
			ret = uUnit
		}
		return ret, nil
	}

	if xBI.Cmp(uUnit) == 0 {
		return uUnit, nil
	}

	if yBI.Cmp(zeroBI) == 0 {
		return uUnit, nil
	}

	if yBI.Cmp(uUnit) == 0 {
		return xBI, nil
	}

	a, err := Log2(x)
	if err != nil {
		return nil, err
	}
	a, err = Mul(a, y)
	return exp2(a), nil
}

func Mul(x SD59x18, y SD59x18) (SD59x18, error) {
	var (
		xBI *big.Int = x
		yBI *big.Int = y
	)

	if xBI.Cmp(uMinSD59x18) == 0 || yBI.Cmp(uMinSD59x18) == 0 {
		return nil, ErrMathSD59x18MulInputTooSmall
	}

	xAbs := new(big.Int).Abs(xBI)
	yAbs := new(big.Int).Abs(yBI)

	resultAbs, err := mulDiv18(xAbs, yAbs)
	if err != nil {
		return nil, err
	}

	if resultAbs.Cmp(uMaxSD59x18) > 0 {
		return nil, ErrMathSD59x18MulOverflow
	}

	sameSign := xBI.Sign() == yBI.Sign()
	result := resultAbs
	if !sameSign {
		result = new(big.Int).Neg(resultAbs)
	}
	return result, nil
}

func Div(x SD59x18, y SD59x18) (SD59x18, error) {
	var (
		xBI = new(big.Int).Set(x)
		yBI = new(big.Int).Set(y)
	)

	if xBI.Cmp(uMinSD59x18) == 0 || yBI.Cmp(uMinSD59x18) == 0 {
		return nil, ErrMathSD59x18DivInputTooSmall
	}

	var (
		xAbs = new(big.Int).Abs(xBI)
		yAbs = new(big.Int).Abs(yBI)
	)

	resultAbs, err := mulDiv(xAbs, uUnit, yAbs)
	if err != nil {
		return nil, err
	}

	if resultAbs.Cmp(uMaxSD59x18) > 0 {
		return nil, ErrMathSD59x18DivOverflow
	}

	sameSign := xBI.Sign() == yBI.Sign()
	result := resultAbs
	if !sameSign {
		result = new(big.Int).Neg(resultAbs)
	}

	return result, nil
}
